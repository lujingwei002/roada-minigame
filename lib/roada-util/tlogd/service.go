package tlogd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/roada-go/roada"
)

type TLogService struct {
	config    Config
	road      *roada.Road
	service   *roada.Service
	chRoad    chan *roada.Request
	file      *os.File
	timeStart int64
	lineTotal int64
	byteTotal int64
	netWriter *bufio.Writer
	conn      net.Conn
	fileIndex int
}

func newTLogService(road *roada.Road, config *Config) error {
	var self = &TLogService{
		config:    *config,
		road:      road,
		chRoad:    make(chan *roada.Request, 1),
		timeStart: time.Now().Unix(),
		lineTotal: 0,
		byteTotal: 0,
		fileIndex: 0,
	}
	if len(config.TcpAddr) > 0 {
		conn, err := net.Dial("tcp", config.TcpAddr)
		if err != nil {
			return err
		}
		self.conn = conn
		self.netWriter = bufio.NewWriter(conn)
	}
	//创建目录
	if _, err := os.Stat(config.Dir); err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(config.Dir, 0666); err != nil {
			return err
		}
	}
	if _, err := os.Stat(config.BackupDir); err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(config.BackupDir, 0666); err != nil {
			return err
		}
	}
	if err := road.LocalSet("tlog"); err != nil {
		return err
	}
	if err := road.Handle("tlog", self); err != nil {
		return err
	}
	service, err := road.Register(self)
	if err != nil {
		return err
	}
	self.service = service
	if err := self.openFile(); err != nil {
		return err
	}
	go self.loop()
	return nil
}

func (self *TLogService) loop() {
	for {
		select {
		case r := <-self.chRoad:
			{
				self.service.ServeRPC(self, r)
			}
		}
	}
}

func (self *TLogService) ServeRPC(r *roada.Request) {
	self.chRoad <- r
}

func (self *TLogService) Print(r *roada.Request, str string) error {
	if self.config.Console {
		self.writeConsole(str)
	}
	self.writeFile(str)
	self.writeTcp(str)
	return nil
}

func (self *TLogService) writeConsole(str string) {
	log.Println("[tlog]", str)
}

func (self *TLogService) checkRotateFile() bool {
	if self.lineTotal > self.config.LineLimit {
		return true
	}
	now := time.Now().Unix()
	if now-self.timeStart > self.config.TimeLimit*60 {
		return true
	}
	return false
}

func (self *TLogService) tryReconnect() error {
	if len(self.config.TcpAddr) <= 0 {
		return nil
	}
	if self.conn != nil {
		return nil
	}
	conn, err := net.Dial("tcp", self.config.TcpAddr)
	if err != nil {
		return err
	}
	self.conn = conn
	self.netWriter = bufio.NewWriter(conn)
	return nil
}

func (self *TLogService) writeTcp(str string) error {
	if len(self.config.TcpAddr) <= 0 {
		return nil
	}
	for i := 0; i < 2; i++ {
		if err := self.tryReconnect(); err != nil {
			continue
		}
		_, err := self.netWriter.WriteString(str + "\n")
		if err == nil {
			self.netWriter.Flush()
			return nil
		}
		self.conn.Close()
		self.netWriter = nil
		self.conn = nil
	}
	log.Println("discard tlog", str)
	return nil
}

func (self *TLogService) writeFile(str string) error {
	if self.file == nil {
		if err := self.openFile(); err != nil {
			return err
		}
	}
	if self.checkRotateFile() {
		self.closeFile()
		if err := self.openFile(); err != nil {
			return err
		}
	}
	byteWrited, err := self.file.WriteString(str + "\n")
	if err != nil {
		return err
	}
	self.lineTotal++
	self.byteTotal = self.byteTotal + int64(byteWrited)
	return nil
}

func (self *TLogService) rotateFile() error {
	oldFileName := fmt.Sprintf("%s/%s_tlog.log", self.config.Dir, self.road.NodeName)
	backupFileName := ""
	backupFilePath := ""
	now := time.Now()
	for i := 1; i <= 9999; i++ {
		self.fileIndex++
		backupFileName = fmt.Sprintf("%s_tlog_%04d%02d%02d%02d%02d%02d%02d.log",
			self.road.NodeName, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), self.fileIndex)
		backupFilePath = fmt.Sprintf("%s/%s", self.config.BackupDir, backupFileName)
		if _, err := os.Stat(backupFilePath); err == nil {
			continue
		}
		if err := os.Rename(oldFileName, backupFilePath); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("[tlogd] rename file failed")
}

func (self *TLogService) copyFile(srcFile string, dstFile string) error {
	source, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer source.Close()
	destination, err := os.OpenFile(dstFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer destination.Close()
	buf := make([]byte, 1024*1024)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}

func (self *TLogService) openFile() error {
	if self.file != nil {
		self.file.Close()
		self.file = nil
	}
	fileName := fmt.Sprintf("%s/%s_tlog.log", self.config.Dir, self.road.NodeName)
	if _, err := os.Stat(fileName); err == nil {
		if err := self.rotateFile(); err != nil {
			return err
		}
	}
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	self.lineTotal = 0
	self.timeStart = time.Now().Unix()
	self.file = file
	return nil
}

func (self *TLogService) closeFile() {
	if self.file == nil {
		return
	}
	self.file.Close()
	self.file = nil
}

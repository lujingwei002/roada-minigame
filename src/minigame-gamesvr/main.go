package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	_ "net/http/pprof"

	"github.com/roada-go/gat"
	"github.com/roada-go/gat/middleware"
	"github.com/roada-go/gat/serialize/protobuf"
	"github.com/roada-go/roada"
	"github.com/shark/minigame-common/conf"
	"github.com/shark/minigame-common/db"
	"github.com/shark/minigame-common/tlog"
	"github.com/shark/minigame/game"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Minigame"
	app.Author = "lujingwei"
	app.Email = "lujingwei@xx.org"
	app.Description = "Minigame"
	// flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "gate",
			Value: "192.168.1.208:9003",
			Usage: "gate address",
		},
		cli.StringFlag{
			Name:  "local",
			Value: ":9002",
			Usage: "server address",
		},
		cli.StringFlag{
			Name:  "coord",
			Value: ":9000",
			Usage: "coord address",
		},
		cli.StringFlag{
			Name:  "name",
			Value: "gamesvr_1",
			Usage: "server name",
		},
		cli.StringFlag{
			Name:  "pprof",
			Value: "",
			Usage: "pprof address",
		},
	}
	app.Action = runAction
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("[main] startup server error %+v", err)
	}
}

func startPprof(addr string) {
	if addr == "" {
		return
	}
	runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
	runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪
	go func() {
		http.ListenAndServe(addr, nil)
	}()
}

func runAction(args *cli.Context) error {
	localAddr := args.String("local")
	if localAddr == "" {
		return fmt.Errorf("[main] server address cannot empty")
	}
	coordAddr := args.String("coord")
	if coordAddr == "" {
		return fmt.Errorf("[main] coord address cannot empty")
	}
	gateAddr := args.String("gate")
	if gateAddr == "" {
		return fmt.Errorf("[main] gate address cannot empty")
	}
	serverName := args.String("name")
	if serverName == "" {
		return fmt.Errorf("[main] server name cannot empty")
	}
	pprofAddr := args.String("pprof")

	road := roada.Default()
	gate := gat.Default()
	//加载配置文件
	conf.Load()
	if conf.Ini.Basic.Debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
	//启动pprof
	startPprof(pprofAddr)
	//启动RPC
	err := road.Slave(serverName, localAddr, coordAddr)
	if err != nil {
		log.Fatalf("[main] roada.Slave failed, error=%s\n", err.Error())
	}
	if err := tlog.Register(road); err != nil {
		log.Fatalf("[main] tlog.Register failed, error=%s\n", err.Error())
	}
	//启动db服务
	db.Register(road)
	//启动游戏服务
	if err := game.Register(road, gate); err != nil {
		log.Fatalf("[main] game.Register failed, error=%s\n", err.Error())
	}
	//启动网关
	if conf.Ini.Basic.Debug {
		gate.WithDebugMode()
	}
	gate.WithCheckOriginFunc(func(_ *http.Request) bool {
		log.Println("WithCheckOriginFunc")
		return true
	})
	log.Println("fffffffffffffffffffffff")

	gate.WithSerializer(protobuf.NewSerializer())
	gate.WithRSAPrivateKey("priv.key")
	gate.WithIsWebsocket(true)
	gate.WithWSPath("/minigame")
	//gate.UseMiddleware(&middleware.Stat{})
	gate.UseMiddleware(middleware.NewGateStat())
	if conf.Ini.Game.UseTLS {
		gate.WithTSLConfig(conf.Ini.Game.TLSCrt, conf.Ini.Game.TLSKey)
	}

	a := make([]int, 1, 1)
	log.Println(len(a), cap(a), a)
	a1 := append(a, 1, 2, 3)

	a[0] = 22
	log.Println(len(a), cap(a), a)
	log.Println(len(a1), cap(a1), a1)
	if err := gate.Run(gateAddr); err != nil {
		log.Fatalln(err)
	}

	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	select {
	case s := <-sg:
		log.Println("got signal", s)
		gate.Maintain(true)     //设置维护状态
		gate.Kick("先休息一会，很快回来") //踢人下线
		gate.Shutdown()         //关闭网关
		game.Shutdown()         //关闭游戏
		road.Shutdown()         //关闭RPC
	}
	log.Printf("[main] quit")
	return nil
}

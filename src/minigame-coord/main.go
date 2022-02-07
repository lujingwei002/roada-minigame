package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/roada-go/roada"
	"github.com/shark/minigame-common/conf"
	"github.com/shark/minigame-common/db"
	"github.com/shark/minigame-coord/rank"
	"github.com/shark/minigame-coord/router"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Minigame Coord"
	app.Author = "lujingwei"
	app.Email = "lujingwei@xx.org"
	app.Description = "Minigame Coord"
	// flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "restore",
			Value: "true",
			Usage: "restore mode",
		},
		cli.StringFlag{
			Name:  "local",
			Value: ":9000",
			Usage: "server address",
		},
		cli.StringFlag{
			Name:  "http",
			Value: ":9001",
			Usage: "http address",
		},
	}
	app.Action = runAction
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("[main] startup server error %+v", err)
	}
}

func runCoord(localAddr string, restore bool) {
	//初始化roada
	road := roada.Default()
	err := road.Master(localAddr, restore)
	if err != nil {
		log.Fatalf("[main] roada.Master failed, error=%s", err.Error())
	}
	//tlog.Register(road)
	db.Register(road)
	rank.Register()
	//road.Run()
}

func runHttp(httpAddr string) {
	go router.Run(httpAddr)
}

func runAction(args *cli.Context) error {
	localAddr := args.String("local")
	if localAddr == "" {
		return fmt.Errorf("[main] server address cannot empty")
	}
	httpAddr := args.String("http")
	if httpAddr == "" {
		return fmt.Errorf("[main] server address cannot empty")
	}
	restore := args.Bool("restore")
	conf.Load()
	if conf.Ini.Basic.Debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
	runHttp(httpAddr)
	runCoord(localAddr, restore)
	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	select {
	case s := <-sg:
		log.Println("[main] got signal", s)
		road := roada.Default()
		road.Shutdown()
	}
	log.Printf("[main] quit")
	return nil
}

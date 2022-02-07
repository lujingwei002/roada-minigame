package main

import (
	"fmt"
	"log"
	"os"
	"time"

	clients "github.com/roada-go/cli"
	"github.com/roada-go/gat/serialize/protobuf"
	"github.com/shark/testclient/game"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Minigame client"
	app.Author = "lujingwei"
	app.Email = "lujingwei@xx.org"
	app.Description = "Minigame client"
	// flags
	app.Flags = []cli.Flag{}
	app.Action = runAction
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app.Commands = []cli.Command{
		{
			Name: "daxigua",
			//Aliases: []string{"start"},
			Usage:  "Start daxigua test",
			Action: daxiguaAction,
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "gate",
			Value: ":9003",
			Usage: "gate address",
		},
		cli.StringFlag{
			Name:  "count",
			Value: "1",
			Usage: "user count",
		},
	}
	app.Action = runAction
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("[main] startup server error %+v", err)
	}
}

func daxiguaAction(args *cli.Context) error {
	return nil
}

func runAction(args *cli.Context) error {
	gateAddr := args.String("gate")
	if gateAddr == "" {
		return fmt.Errorf("[main] gate address cannot empty")
	}
	userCount := int32(args.Int("count"))
	if userCount == 0 {
		return fmt.Errorf("[main] user count cannot empty")
	}
	client := clients.Default()
	//c.WithDebugMode()
	//消息包类型
	client.WithSerializer(protobuf.NewSerializer())
	//协议加密
	client.WithRSAPublicKey("pub.key")
	client.WithIsWebsocket(true)
	client.WithWSPath("/minigame")
	client.WithServerAdd(gateAddr)
	client.WithTimerPrecision(time.Millisecond)
	game.Register(client, userCount)
	client.Run()
	return nil
}

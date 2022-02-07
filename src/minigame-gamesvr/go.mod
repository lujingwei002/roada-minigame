module github.com/lujingwei002/minigame

require (
	github.com/go-redis/redis/v8 v8.10.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/icattlecoder/godaemon v0.0.0-20140627053409-f0fff2a3c017 // indirect
	github.com/jmoiron/sqlx v1.3.4
	github.com/roada-go/gat v0.0.0
	github.com/roada-go/roada v0.0.0
	github.com/roada-go/util v0.0.0
	github.com/lujingwei002/minigame-common v0.0.0
	github.com/urfave/cli v1.22.5
	google.golang.org/grpc v1.38.0 // indirect
	google.golang.org/protobuf v1.26.0
	gopkg.in/ini.v1 v1.62.0
	honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc
)

replace github.com/roada-go/gat => ../../lib/roada-gate

replace github.com/roada-go/roada => ../../lib/roada

replace github.com/roada-go/util => ../../lib/roada-util

replace github.com/lujingwei002/minigame-common => ../minigame-common

go 1.15

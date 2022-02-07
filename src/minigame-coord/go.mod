module github.com/lujingwei002/minigame-coord

go 1.15

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/gin-gonic/gin v1.7.2
	github.com/roada-go/roada v0.0.0
	github.com/lujingwei002/minigame-common v0.0.0-00010101000000-000000000000
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.5.1
	github.com/urfave/cli v1.22.5
)

replace github.com/roada-go/roada => ../../lib/roada

replace github.com/roada-go/util => ../../lib/roada-util

replace github.com/lujingwei002/minigame-common => ../minigame-common

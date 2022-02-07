module github.com/lujingwei002/minigame-client

require (
	github.com/roada-go/cli v0.0.0
	github.com/roada-go/gat v0.0.0
	github.com/lujingwei002/minigame-common v0.0.0
	github.com/urfave/cli v1.22.5
)

replace github.com/roada-go/gat => ../../lib/roada-gate

replace github.com/roada-go/cli => ../../lib/roada-client

replace github.com/lujingwei002/minigame-common => ../minigame-common

go 1.15

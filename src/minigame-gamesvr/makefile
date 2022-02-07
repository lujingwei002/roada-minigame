all:
	go build -o ../../bin/minigame-gamesvr main.go
	cd ../../runtime/minigame-gamesvr; ../../bin/minigame-gamesvr --pprof :6060

proto:
	protoc --go_out=gamepb gamepb/gamepb.proto

start:
	supervisorctl start minigame	

stop:
	supervisorctl stop minigame		
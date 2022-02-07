package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/roada-go/roada"
	"github.com/shark/minigame-common/conf"
)

func Register(road *roada.Road) error {
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		conf.Ini.MySql.User, conf.Ini.MySql.Password, conf.Ini.MySql.Ip, conf.Ini.MySql.Port, conf.Ini.MySql.Db, conf.Ini.MySql.Charset)
	for i := 1; i <= 5; i++ {
		if err := newDbService(road, addr, i); err != nil {
			panic(err)
		}
	}
	return nil
}

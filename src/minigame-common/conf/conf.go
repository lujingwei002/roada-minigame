package conf

import (
	"log"

	"encoding/json"

	"gopkg.in/ini.v1"
)

var Ini struct {
	MySql struct {
		Ip       string `ini:"ip"`
		Port     int    `ini:"port"`
		User     string `ini:"user"`
		Password string `ini:"password"`
		Db       string `ini:"db"`
		Charset  string `ini:"charset"`
	} `ini:"mysql"`

	Redis struct {
		Ip       string `ini:"ip"`
		Port     int    `ini:"port"`
		Password string `ini:"password"`
		Db       int    `ini:"db"`
	} `ini:"redis"`

	Game struct {
		Id                int32  `ini:"id"`
		Name              string `ini:"name"`
		UrlValidate       bool   `ini:"urlvalidate"`
		UrlValidateSecret string `ini:"urlvalidate_secret"`
		ConfigDir         string `ini:"config_dir"`
		UseTLS            bool   `ini:"use_tls"`
		TLSCrt            string `ini:"tls_crt"`
		TLSKey            string `ini:"tls_key"`
	} `ini:"game"`

	Tlog struct {
		Dir       string `ini:"dir"`
		BackupDir string `ini:"backupdir"`
		TimeLimit int64  `ini:"timelimit"`
		LineLimit int64  `ini:"linelimit"`
		TcpAddr   string `ini:"tcpaddr"`
		Console   bool   `ini:"console"`
	} `ini:"tlog"`

	SharkSdk struct {
		Url    string `ini:"url"`
		Secret string `ini:"secret"`
	} `ini:"sharksdk"`

	Basic struct {
		Debug bool `ini:"debug"`
	} `ini:"basic"`
	Coord struct {
		UseTLS bool   `ini:"use_tls"`
		TLSPem string `ini:"tls_pem"`
		TLSKey string `ini:"tls_key"`
	} `ini:"coord"`
}

func Load() {
	err := ini.MapTo(&Ini, "config.ini")
	if err != nil {
		panic(err)
	}
	str, err := json.MarshalIndent(Ini, "", "\t")
	if err != nil {
		panic(err)
	}
	log.Printf("[config] %s\n", str)
}

func Reload() error {
	return nil
}

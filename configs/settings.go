package configs

import (
	"embed"
	"gopkg.in/ini.v1"
	"log"
	"reflect"
)

// basicSettings 服务的基本配置
type basicSettings struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	StartupMode    string `json:"startup_mode"`
	DisplayName    string `json:"display_name"`
	Description    string `json:"description"`
	SrvDepends     string `json:"srv_depends"`
	SrvProtocol    string `json:"srv_protocol"`
	SrvIP          string `json:"srv_http_ip"`
	SrvPort        string `json:"srv_http_port"`
	GlobalLogPath  string `json:"global_log_path"`
	CronLoggerPath string `json:"cron_logger_path"`
}

// postgresqlSettings 服务所依赖的postgresql连接配置
type postgresqlSettings struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DB       string `json:"db"`

	ConnMaxLifetime int `json:"conn_max_life_time"`
	MaxIdleConn     int `json:"max_idle_conn"`
	MaxOpenConn     int `json:"max_open_conn"`
}

// redisSettings 服务所依赖的redis连接配置
type redisSettings struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       string `json:"db"`

	MaxRetries  int `json:"max_retries"`
	MinIdleConn int `json:"min_idle_conn"`
	PoolSize    int `json:"pool_size"`
}

type Ss struct {
	Base  basicSettings
	DB    postgresqlSettings
	Cache redisSettings
}

var Settings = Load()

// go:embed configs
var configFiles embed.FS

func Load() Ss {
	s := new(Ss)

	data, er := configFiles.ReadFile("configs/dashboard.ini")
	if er != nil {
		log.Fatal("load basic config: ", er)
	}
	cfg, el := ini.Load(data)
	if el != nil {
		log.Fatal("load basic config stream: ", el)
	}

	parse := func(_section *ini.Section, _setting any) {
		bt := reflect.TypeOf(_setting)
		bv := reflect.ValueOf(_setting)

		for i := 0; i < bt.NumField(); i++ {
			field := bt.Field(i).Tag.Get("json")
			key, err := _section.GetKey(field)
			if err != nil {
				log.Fatal("parse basic config stream: ", err)
			}
			switch bt.Field(i).Type.Kind() {
			case reflect.Int:
				v, ev := key.Int64()
				if ev != nil {
					log.Fatal("convert config to int: ", ev)
				}
				bv.Field(i).SetInt(v)
			case reflect.String:
				bv.Field(i).SetString(key.String())
			default:
				log.Fatal("invalid config type")
			}
		}
	}

	parse(cfg.Section(""), &s.Base)
	parse(cfg.Section("db"), &s.DB)
	parse(cfg.Section("cache"), &s.Cache)

	return *s
}

package configs

// basicSettings 服务的基本配置
type basicSettings struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	StartupMode   string `json:"startup_mode"`
	DisplayName   string `json:"display_name"`
	Description   string `json:"description"`
	SrvDepends    string `json:"srv_depends"`
	SrvIP         string `json:"srv_http_ip"`
	SrvPort       string `json:"srv_http_port"`
	GlobalLogPath string `json:"global_log_path"`
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

type Settings struct {
	Basic      basicSettings
	Postgresql postgresqlSettings
	Redis      redisSettings
	// TODO more
}

func init() {

}

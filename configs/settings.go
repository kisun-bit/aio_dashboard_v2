package configs

// BasicSettings 服务的基本配置
type BasicSettings struct {
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

// PostgresqlSettings 服务所依赖的postgresql连接配置
type PostgresqlSettings struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DB       string `json:"db"`

	ConnMaxLifetime int `json:"conn_max_life_time"`
	MaxIdleConn     int `json:"max_idle_conn"`
	MaxOpenConn     int `json:"max_open_conn"`
}

// RedisSettings 服务所依赖的redis连接配置
type RedisSettings struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       string `json:"db"`

	MaxRetries  int `json:"max_retries"`
	MinIdleConn int `json:"min_idle_conn"`
	PoolSize    int `json:"pool_size"`
}

type Setting struct {
	Basic      BasicSettings
	Postgresql PostgresqlSettings
	Redis      RedisSettings
	// TODO more
}

var Settings Setting = loadConfig()

// loadConfig 加载【应用基础配置】、【数据库配置】、【缓存配置】、【消息队列配置】等
func loadConfig() Setting {

}

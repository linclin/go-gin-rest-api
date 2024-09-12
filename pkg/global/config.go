package global

import "log/slog"

// Configuration 系统配置, 配置字段可参见yml注释
// viper内置了mapstructure, yml文件用"-"区分单词, 转为驼峰方便
type Configuration struct {
	System    SystemConfiguration    `mapstructure:"system" json:"system"`
	Logs      LogsConfiguration      `mapstructure:"logs" json:"logs"`
	Mysql     *MysqlConfiguration    `mapstructure:"mysql" json:"mysql"`
	Pgsql     *PgsqlConfiguration    `mapstructure:"pgsql" json:"pgsql"`
	Casbin    CasbinConfiguration    `mapstructure:"casbin" json:"casbin"`
	Jwt       JwtConfiguration       `mapstructure:"jwt" json:"jwt"`
	RateLimit RateLimitConfiguration `mapstructure:"rate-limit" json:"rateLimit"`
}

type SystemConfiguration struct {
	AppName       string `mapstructure:"app-name" json:"appName"`
	RunMode       string `mapstructure:"run-mode" json:"runMode"`
	UrlPathPrefix string `mapstructure:"url-path-prefix" json:"urlPathPrefix"`
	Port          int    `mapstructure:"port" json:"port"`
	BaseApi       string `mapstructure:"base-api" json:"baseApi"`
	Transaction   bool   `mapstructure:"transaction" json:"transaction"`
	InitData      bool   `mapstructure:"init-data" json:"initData"`
}

type LogsConfiguration struct {
	Level      slog.Level `mapstructure:"level" json:"level"`
	Path       string     `mapstructure:"path" json:"path"`
	MaxSize    int        `mapstructure:"max-size" json:"maxSize"`
	MaxBackups int        `mapstructure:"max-backups" json:"maxBackups"`
	MaxAge     int        `mapstructure:"max-age" json:"maxAge"`
	Compress   bool       `mapstructure:"compress" json:"compress"`
}

type MysqlConfiguration struct {
	Username    string `mapstructure:"username" json:"username"`
	Password    string `mapstructure:"password" json:"password"`
	Database    string `mapstructure:"database" json:"database"`
	Host        string `mapstructure:"host" json:"host"`
	Port        int    `mapstructure:"port" json:"port"`
	Query       string `mapstructure:"query" json:"query"`
	LogMode     bool   `mapstructure:"log-mode" json:"logMode"`
	TablePrefix string `mapstructure:"table-prefix" json:"tablePrefix"`
}

type PgsqlConfiguration struct {
	Username    string `mapstructure:"username" json:"username"`
	Password    string `mapstructure:"password" json:"password"`
	Database    string `mapstructure:"database" json:"database"`
	Host        string `mapstructure:"host" json:"host"`
	Port        int    `mapstructure:"port" json:"port"`
	SslMode     string `mapstructure:"ssl-mode" json:"sslMode"`
	LogMode     bool   `mapstructure:"log-mode" json:"logMode"`
	TablePrefix string `mapstructure:"table-prefix" json:"tablePrefix"`
}
type CasbinConfiguration struct {
	ModelPath string `mapstructure:"model-path" json:"modelPath"`
}

type JwtConfiguration struct {
	Timeout    int `mapstructure:"timeout" json:"timeout"`
	MaxRefresh int `mapstructure:"max-refresh" json:"maxRefresh"`
}

type RateLimitConfiguration struct {
	Max int64 `mapstructure:"max" json:"max"`
}

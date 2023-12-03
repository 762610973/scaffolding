package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"scaffolding/pkg/internal"
	"time"
)

var Conf = Config{}

type Config struct {
	System     System     `toml:"system"`
	Zap        Zap        `toml:"zap"`
	Lumberjack Lumberjack `toml:"lumberjack"`
	Storage    Storage    `toml:"storage"`
}

type Storage struct {
	Mysql `toml:"mysql"`
	Redis `toml:"redis"`
}

type Mysql struct {
	Ip       string `toml:"ip"`
	Port     string `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

func (m Mysql) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&Loc=Local", m.User, m.Password, m.Ip, m.Port, m.Database)
}

type Redis struct{}

type System struct {
	Port string `toml:"port"`
	// Mode: debug, release
	Mode        string        `toml:"mode"`
	QuitMaxTime time.Duration `toml:"quit_max_time"`
}

type Zap struct {
	Level        string `toml:"level"`          // 级别
	Format       string `toml:"format"`         // Format: json, console
	EncodeLevel  string `toml:"encode_level"`   // 编码级
	ShowLine     bool   `toml:"show_line"`      // 显示行
	LogInConsole bool   `toml:"log_in_console"` // 输出控制台
	LogInFile    bool   `toml:"log_in_file"`
}

type Lumberjack struct {
	FileName   string `toml:"filename"`
	MaxSize    int    `toml:"max_size"`
	MaxBackups int    `toml:"max_backups"`
	MaxAge     int    `toml:"max_age"`
	Compress   bool   `toml:"compress"`
}

type JWT struct {
	SigningKey  string `toml:"signing_key"`  // jwt签名
	ExpiresTime int64  `toml:"expires_time"` // 过期时间
	BufferTime  int64  `toml:"buffer_time"`  // 缓冲时间
	Issuer      string `toml:"issuer"`       // 签发者
}

func InitConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	// 必须设置TagName
	decodeFn := func(config *mapstructure.DecoderConfig) { config.TagName = "toml" }
	if err := viper.Unmarshal(&Conf, decodeFn); err != nil {
		return err
	}
	viper.OnConfigChange(func(_ fsnotify.Event) {
		var c Config
		if err := viper.Unmarshal(&c, decodeFn); err == nil {
			Conf = c
			internal.SetLogLevel(c.Zap.Level)
		}
	})
	viper.WatchConfig()

	return nil
}

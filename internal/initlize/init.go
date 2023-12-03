package initlize

import (
	"flag"

	cfg "scaffolding/pkg/config"
	zlog "scaffolding/pkg/log"
)

var (
	configPath = flag.String("conf", "./config.toml", "the config file path")
)

func Init() error {
	flag.Parse()
	if err := cfg.InitConfig(*configPath); err != nil {
		return err
	}
	zlog.InitLog()
	return nil
}

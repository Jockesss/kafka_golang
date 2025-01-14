package modules

type LoggerCfg struct {
	Level string `env:"LEVEL" env-required:"true"`
}

package modules

type Kafka struct {
	Brokers            []string `env:"BROKERS" env-required:"true"`
	ApproveDebugNoAuth bool     `env:"APPROVE_NO_AUTH" default:"false"`
	//SaslUser                 string `env:"SASL_USER" env-required:"true"`
	//SaslPassword             string `env:"SASL_PWD" env-required:"true"`
	ProducerBufferMaxMsg int    `env:"PRODUCER_BUFFER" default:"100000"`
	IgnoreTimeoutMs      int    `env:"IGNORE_TIMEOUT_MS" env-required:"true"`
	RetryCount           int    `env:"RETRY_COUNT" env-required:"true"`
	TestTopic            string `env:"TEST_TOPIC" env-required:"true"`
}

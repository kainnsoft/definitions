package rmq

type Config struct {
	Host                  string           `yaml:"host" validate:"nonzero"`
	Port                  int              `yaml:"port" validate:"nonzero"`
	Username              string           `yaml:"username" validate:"nonzero"`
	Password              string           `yaml:"password" validate:"nonzero"`
	VirtualHost           string           `yaml:"virtualHost"`
	ConnectionName        string           `yaml:"connectionName"`
	DefaultQos            int              `yaml:"defaultQos"`
	Qos                   map[string]int   `yaml:"qos"`
	DefaultConsumersCount int              `yaml:"defaultConsumersCount" validate:"nonzero"`
	ConsumersCount        map[string]int   `yaml:"consumersCount"`
	DefaultTimeout        int64            `yaml:"defaultTimeout"`
	Timeouts              map[string]int64 `yaml:"timeouts"`
	LogTrimmedFields      []string         `yaml:"logTrimmedFields"`
	ShortLogQueues        []string         `yaml:"shortLogQueues"`
	Debug                 bool             `yaml:"debug"`
}

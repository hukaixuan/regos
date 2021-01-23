package config

type Config struct {
	Addr  string `toml:"addr"`
	DBNum int    `toml:"db_num"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Addr:  "127.0.0.1:6380",
		DBNum: 16,
	}
}

func LoadConfigFromFile(file string) (*Config, error) {
	return &Config{}, nil
}

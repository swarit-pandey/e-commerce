package config

type Config struct {
	Order OrderConfig `mapstructure:"order"`
}

type OrderConfig struct {
	Metadata Metadata `mapstructure:"metadata"`
	Server   Server   `mapstructure:"server"`
	Broker   Broker   `mapstructure:"broker"`
	Cache    Cache    `mapstructure:"cache"`
	Database Database `mapstructure:"database"`
}

type Metadata struct {
	Enabled bool   `mapstructure:"enabled"`
	Domain  string `mapstructure:"domain"`
	Env     string `mapstructure:"env"`
}

type Publish struct {
	Exchange string `mapstructure:"exchange"`
	Queue    string `mapstructure:"queue"`
}

type Consume struct {
	Exchange string `mapstructure:"exchange"`
	Queue    string `mapstructure:"queue"`
}

type Server struct {
	BasePath string `mapstructure:"basepath"`
	Port     int    `mapstructure:"port"`
	Address  string `mapstructure:"address"`
}

type Broker struct {
	Port    int     `mapstructure:"port"`
	Driver  string  `mapstructure:"driver"`
	Address string  `mapstructure:"address"`
	Publish Publish `mapstructure:"publish"`
	Consume Consume `mapstructure:"consume"`
}

type Cache struct {
	Port   int    `mapstructure:"port"`
	Driver string `mapstructure:"driver"`
}

type Database struct {
	Port        int    `mapstructure:"port"`
	Driver      string `mapstructure:"driver"`
	Address     string `mapstructure:"address"`
	Dialect     string `mapstructure:"dialect"`
	Name        string `mapstructure:"name"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	SSL         string `mapstructure:"ssl"`
	MaxConnPool int    `mapstructure:"maxconnpool"`
}

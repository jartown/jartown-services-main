package config

var Conf Config

type Config struct {
	Env     string `yaml:"env"`
	Debug   bool   `yaml:"debug"`
	Port    int    `yaml:"port"`
	MongoDB struct {
		URI      string `yaml:"uri"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Schema   string `yaml:"schema"`
	} `yaml:"mongodb"`
	JWTMenu struct {
		Key     string `yaml:"key"`
		Timeout int    `yaml:"timeout-secs"`
	} `yaml:"jwt-menu"`
	Url string `yaml:"uri"`
}

const ApplicationName = "jartown-services-main"

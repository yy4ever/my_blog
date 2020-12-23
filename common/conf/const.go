package conf

type Conf struct {
	DbHost               string `yaml:"DbHost"`
	DbPort               int    `yaml:"DbPort"`
	DbDataBase           string `yaml:"DbDataBase"`
	DbUser               string `yaml:"DbUser"`
	DbPassword           string `yaml:"DbPassword"`
	RedisHost            string `yaml:"RedisHost"`
	RedisPwd             string `yaml:"RedisPwd"`
	RedisDb              int    `yaml:"RedisDb"`
	SecretKey            string `yaml:"SecretKey"`
	LogFilePath          string `yaml:"LogFilePath"`
	PORT                 string `yaml:"PORT"`
	LoginRememberSeconds int    `yaml:"LoginRememberSeconds"`
}

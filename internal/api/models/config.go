package config

type Config struct {
	RedisAddr      string `yaml:"redisAddr"`
	RedisPass      string `yaml:"redisPass"`
	Port           string `yaml:"port"`
	ShortUrlLength int    `yaml:"shortUrlLength"`
}

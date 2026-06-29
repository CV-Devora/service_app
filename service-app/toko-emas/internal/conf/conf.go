package conf

// Config is the top-level configuration structure
type Config struct {
	Server Server `yaml:"server"`
	Data   Data   `yaml:"data"`
	Log    Log    `yaml:"log"`
	Auth   Auth   `yaml:"auth"`
}

type Server struct {
	HTTP HTTP `yaml:"http"`
}

type HTTP struct {
	Addr    string `yaml:"addr"`
	Timeout string `yaml:"timeout"`
}

type Data struct {
	Database Database `yaml:"database"`
}

type Database struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

type Log struct {
	Level string `yaml:"level"`
}

type Auth struct {
	JWTSecret      string `yaml:"jwt_secret"`
	AccessTokenTTL string `yaml:"access_token_ttl"`
}

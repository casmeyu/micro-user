package structs

type JwtConfig struct {
	Secret     string
	Expiration int
}

type AppConfig struct {
	Name string
	Ip   string
}

type DbConfig struct {
	User     string
	Password string
	DbName   string
	Ip       string
}

type Config struct {
	App AppConfig
	Db  DbConfig
	Jwt JwtConfig
}

package structs

type JwtConfig struct {
	Secret     string `json:"-"`
	Expiration int    `json:"expiration"`
}

type AppConfig struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

type DbConfig struct {
	User     string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"db"`
	Ip       string `json:"ip"`
}

type Config struct {
	App AppConfig `json:"app"`
	Db  DbConfig  `json:"db"`
	Jwt JwtConfig `json:"jwt"`
}

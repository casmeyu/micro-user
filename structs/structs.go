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

type IError struct {
	Field string
	Tag   string
	Value string
}

type UserLogin struct {
	Username string `json:"username" validate:"required,min=1,max=50"`
	Password string `json:"password" validate:"required,min=1"`
}

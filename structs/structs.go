package structs

import "time"

// Config
type JwtConfig struct {
	Secret     string `json:"-"`
	Expiration int    `json:"expiration"`
}

type AppConfig struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

type DbConfig struct {
	User     string `json:"-"`
	Password string `json:"-"`
	Name     string `json:"db"`
	Ip       string `json:"ip"`
}

type Config struct {
	App AppConfig `json:"app"`
	Db  DbConfig  `json:"db"`
	Jwt JwtConfig `json:"jwt"`
}

// END Config

// Services
//
//	General
type ServiceResponse struct {
	Success bool
	Status  int
	Result  interface{}
	Err     string
}

type UserLogin struct {
	Username string `json:"username" validate:"required,min=1,max=50"`
	Password string `json:"password" validate:"required,min=1"`
}

type PublicUser struct {
	Id             uint      `json:"id"`
	Username       string    `json:"username"`
	LastConnection time.Time `json:"lastConnection"`
}

// Auth
type LoginCredentials struct {
	PublicUser
	Jwt string
}

// END Services

// General
type IError struct {
	Field string
	Tag   string
	Value string
}

// END General

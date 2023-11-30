package config

type Smtp struct {
	Enable   bool
	Host     string
	Port     int32
	Email    string
	Password string
}

var SmtpConfig = new(Smtp)

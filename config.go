package main

type Config struct {
	Round  `ini:"round"`
	Remote `ini:"remote"`
	Checks []*Check
}

func NewConfig() *Config {
	return &Config{Round: Round{MaxPoints: 100}}
}

type Round struct {
	Title     string `ini:"title"`
	Os        string `ini:"os"`
	User      string `ini:"user"`
	Local     string `ini:"local"`
	MaxPoints int    `ini:"maxPoints"`
}

type Remote struct {
	Enable   bool   `ini:"enable"`
	Name     string `ini:"name"`
	Server   string `ini:"server"`
	Password string `ini:"password"`
}

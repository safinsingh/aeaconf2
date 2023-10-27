package main

import "github.com/fatih/color"

type Config struct {
	Round  `ini:"round"`
	Remote `ini:"remote"`
	Checks []*Check
}

func NewConfig() *Config {
	return &Config{Round: Round{MaxPoints: 100}}
}

type Check struct {
	Message string
	Points  int
	// points were left unspecified
	PointsEmpty bool

	Condition
	// separate root hint from condition tree
	Hint string
}

func (c *Check) Debug() string {
	cl := color.New(color.Bold)
	ret := cl.Sprintf("%s (%d Points)%s\n", c.Message, c.Points, formatHint(c.Hint))
	return ret + DebugCondition(c.Condition)
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

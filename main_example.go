package main

import (
	"fmt"

	"github.com/safinsingh/aeaconf2/compat"
	"gopkg.in/ini.v1"
)

type config struct {
	Round  round  `ini:"round"`
	Remote remote `ini:"remote"`
}

type round struct {
	Title     string `ini:"title"`
	Os        string `ini:"os"`
	User      string `ini:"user"`
	Local     string `ini:"local"`
	MaxPoints int    `ini:"maxPoints"`
}

type remote struct {
	Enable   bool   `ini:"enable"`
	Name     string `ini:"name"`
	Server   string `ini:"server"`
	Password string `ini:"password"`
}

func main() {
	headerRaw, checksRaw, err := compat.SeparateConfigFrom("example.acf")
	if err != nil {
		Fatal(STAGE_PRE, fmt.Sprintf("failed to read config file: %s", err.Error()))
	}

	cfg := new(config)
	err = ini.MapTo(cfg, headerRaw)
	if err != nil {
		Fatal(STAGE_INI, fmt.Sprintf("failed to parse ini header: %s", err.Error()))
	}

	exampleFunctionRegistry := getFunctionRegistry()
	ab := DefaultAeaconfBuilder(checksRaw, exampleFunctionRegistry).
		SetLineOffset(CountLines(headerRaw)).
		SetMaxPoints(cfg.Round.MaxPoints)

	checks := ab.GetChecks()
	for _, check := range checks {
		fmt.Println(check.Debug() + "\n")
	}
}

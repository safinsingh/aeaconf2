package aeaconf2_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/safinsingh/aeaconf2"
	"github.com/safinsingh/aeaconf2/compat"
	"gopkg.in/ini.v1"
)

type Config struct {
	Round  `ini:"round"`
	Remote `ini:"remote"`
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

func TestAeaconf(t *testing.T) {
	headerRaw, checksRaw, err := compat.SeparateConfigFrom("example.acf")
	if err != nil {
		t.Error(errors.Wrap(err, "failed to read config file"))
	}

	cfg := new(Config)
	err = ini.MapTo(cfg, headerRaw)
	if err != nil {
		t.Error(errors.Wrap(err, "failed to parse header"))
	}

	exampleFunctionRegistry := getFunctionRegistry()
	ab := aeaconf2.DefaultAeaconfBuilder(checksRaw, exampleFunctionRegistry).
		SetLineOffset(countLines(headerRaw)).
		SetMaxPoints(cfg.Round.MaxPoints)

	_ = ab.GetChecks()
}

func countLines(source []byte) int {
	ret := 0
	for _, b := range source {
		if b == '\n' {
			ret++
		}
	}
	return ret + 1
}

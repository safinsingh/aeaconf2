package aeaconf2_test

import (
	"fmt"
	"testing"
	"unicode/utf8"

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

func GetAeaconf(t *testing.T) []*aeaconf2.Check {
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

	return ab.GetChecks()
}

func TestAeaconf(t *testing.T) {
	_ = GetAeaconf(t)
}

func TestModifyStrings(t *testing.T) {
	checks := GetAeaconf(t)

	for idx := range checks {
		compat.ModifyConditionStrings(
			checks[idx].Condition,
			func(s string) string { return reverse(s) },
		)

		fmt.Println(checks[idx].Debug())
	}
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

// https://stackoverflow.com/questions/1752414/how-to-reverse-a-string-in-go
func reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}

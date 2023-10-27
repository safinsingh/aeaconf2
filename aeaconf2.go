package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"

	"gopkg.in/ini.v1"
)

func init() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

func separateConfig(config []byte) ([]byte, []byte) {
	re := regexp.MustCompile(`(?m)^[[:space:]]*---[[:space:]]*$`)
	loc := re.FindIndex(config)

	return config[:loc[0]], config[loc[1]:]
}

func main() {
	data, _ := os.ReadFile("example.acf")
	headerIni, checksRaw := separateConfig(data)

	config := NewConfig()
	err := ini.MapTo(config, headerIni)
	if err != nil {
		Fatal(STAGE_INI, fmt.Sprintf("failed to parse ini header: %s", err.Error()))
	}

	l := NewLexer(bytes.TrimSpace(checksRaw), CountLines(headerIni))
	p := NewParser(l, funcRegistry)
	config.Checks = p.Checks()
	config.DistributeMaxPoints()

	for _, check := range config.Checks {
		fmt.Println(check.Debug() + "\n")
	}
}

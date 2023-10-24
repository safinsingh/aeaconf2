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

	var config Config
	ini.StrictMapTo(headerIni, config)

	l := NewLexer(bytes.TrimSpace(checksRaw))
	p := NewParser(l)
	config.Checks = p.Checks()

	for _, check := range config.Checks {
		fmt.Println(check.Debug() + "\n")
	}
}

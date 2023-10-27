package main

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"regexp"

	"gopkg.in/ini.v1"
)

func separateConfig(config []byte) ([]byte, []byte) {
	re := regexp.MustCompile(`(?m)^[[:space:]]*---[[:space:]]*$`)
	loc := re.FindIndex(config)

	return config[:loc[0]], config[loc[1]:]
}

func GetConfig(fileName string, funcRegistry map[string]reflect.Type) *Config {
	data, err := os.ReadFile(fileName)
	if err != nil {
		Fatal(STAGE_PRE, fmt.Sprintf("failed to read config file '%s': %s", fileName, err.Error()))
	}
	headerIni, checksRaw := separateConfig(data)

	config := NewConfig()
	err = ini.MapTo(config, headerIni)
	if err != nil {
		Fatal(STAGE_INI, fmt.Sprintf("failed to parse ini header: %s", err.Error()))
	}

	l := NewLexer(bytes.TrimSpace(checksRaw), CountLines(headerIni))
	p := NewParser(l, funcRegistry)
	config.Checks = p.Checks()
	config.DistributeMaxPoints()

	return config
}

func main() {
	config := GetConfig("example.acf", funcRegistry)
	for _, check := range config.Checks {
		fmt.Println(check.Debug() + "\n")
	}
}

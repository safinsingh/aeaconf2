package compat

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/pkg/errors"
)

func SeparateConfig(data []byte) ([]byte, []byte, error) {
	re := regexp.MustCompile(`(?m)^[[:space:]]*---[[:space:]]*$`)
	loc := re.FindIndex(data)
	if loc == nil {
		return nil, nil, errors.New("could not split configuration at separator, expected: '---' between header and checks")
	}

	return bytes.TrimSpace(data[:loc[0]]), bytes.TrimSpace(data[loc[1]:]), nil
}

func SeparateConfigFrom(fileName string) ([]byte, []byte, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("failed to open config file for reading: '%s'", fileName))
	}

	return SeparateConfig(data)
}

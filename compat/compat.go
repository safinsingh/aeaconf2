package compat

import (
	"bytes"
	"os"
	"regexp"

	"github.com/pkg/errors"
)

func SeparateConfigFrom(fileName string) ([]byte, []byte, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to read config file")
	}

	re := regexp.MustCompile(`(?m)^[[:space:]]*---[[:space:]]*$`)
	loc := re.FindIndex(data)

	return bytes.TrimSpace(data[:loc[0]]), bytes.TrimSpace(data[loc[1]:]), nil
}

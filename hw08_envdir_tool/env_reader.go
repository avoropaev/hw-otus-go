package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]EnvValue

func (e Environment) toStrings() []string {
	envs := make([]string, 0, len(e))

	for name, envValue := range e {
		if name == "" {
			continue
		}

		value := ""
		if !envValue.NeedRemove {
			value = envValue.Value
		}

		envs = append(envs, name+"="+value)
	}

	return envs
}

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envs := make(Environment)

	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}

		fileName := string(bytes.Replace([]byte(fileInfo.Name()), []byte("="), []byte{}, -1))

		if fileInfo.Size() == 0 {
			envs[fileName] = EnvValue{NeedRemove: true}
			continue
		}

		file, err := os.OpenFile(dir+"/"+fileName, os.O_RDONLY, 0644)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		scanner.Scan()

		if err := scanner.Err(); err != nil {
			_ = file.Close()
			return nil, err
		}

		value := getValueFromLine(scanner.Text())
		if value == "" {
			envs[fileName] = EnvValue{NeedRemove: true}
			_ = file.Close()
			continue
		}

		envs[fileName] = EnvValue{Value: value}
		_ = file.Close()
	}

	return envs, nil
}

func getValueFromLine(line string) string {
	value := strings.Replace(line, string([]byte{0x00}), "\n", -1)

	for strings.HasSuffix(value, " ") || strings.HasSuffix(value, "\t") {
		value = strings.TrimRight(value, "\t")
		value = strings.TrimRight(value, " ")
	}

	return value
}

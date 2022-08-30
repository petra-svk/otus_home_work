package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("can't read dir %s: %w", dir, err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("dir %s is empty", dir)
	}
	env := make(Environment)
	for _, source := range files {
		if strings.Contains(source.Name(), "=") {
			return nil, fmt.Errorf("file name `%s` contains '='", source.Name())
		}
		file := filepath.Join(dir, source.Name())
		f, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("can't read file %s: %w", file, err)
		}
		defer f.Close()
		fi, err := f.Stat()
		if err != nil {
			return nil, fmt.Errorf("can't stat %s: %w", file, err)
		}
		if fi.Size() == 0 {
			env[source.Name()] = EnvValue{Value: "", NeedRemove: true}
		} else {
			scanner := bufio.NewScanner(f)
			scanner.Split(bufio.ScanBytes)
			var param []byte
			for scanner.Scan() {
				b := scanner.Bytes()
				if string(b) == string('\n') {
					break
				}
				if string(b) == "\x00" {
					param = append(param, byte('\n'))
					continue
				}
				param = append(param, b[0])
			}
			line := string(param)
			line = strings.TrimRight(line, " \t")
			env[source.Name()] = EnvValue{Value: line, NeedRemove: false}
		}
	}
	return env, nil
}

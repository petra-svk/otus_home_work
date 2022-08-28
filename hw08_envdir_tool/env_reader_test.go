package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const dirPath = "testdata/env"

func TestReadDir(t *testing.T) {
	t.Run("common case", func(t *testing.T) {
		gotEnvMap, err := ReadDir(dirPath)
		require.NoError(t, err)
		expectedEnvMap := Environment{
			"BAR":   {Value: "bar", NeedRemove: false},
			"EMPTY": {Value: "", NeedRemove: false},
			"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": {Value: "\"hello\"", NeedRemove: false},
			"UNSET": {Value: "", NeedRemove: true},
		}
		if !assert.ObjectsAreEqual(expectedEnvMap, gotEnvMap) {
			t.Errorf("map is wrong. got: %v, expected: %v", gotEnvMap, expectedEnvMap)
		}
	})

	t.Run("bad symbol '=' in name of file", func(t *testing.T) {
		tDir := os.TempDir()
		file, err := os.CreateTemp(tDir, "PARAM=")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(file.Name())
		_, err = ReadDir(tDir)
		require.Error(t, err)
	})
}

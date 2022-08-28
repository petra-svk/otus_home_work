package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("exec command with set env", func(t *testing.T) {
		command := []string{"/bin/bash", "./testdata/echo.sh", "arg1=1", "arg2=2"}
		EnvMap := Environment{
			"BAR":   {Value: "bar", NeedRemove: false},
			"EMPTY": {Value: "", NeedRemove: false},
			"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": {Value: "\"hello\"", NeedRemove: false},
			"UNSET": {Value: "", NeedRemove: true},
		}
		result := RunCmd(command, EnvMap)
		require.Equal(t, int(0), result)
	})
}

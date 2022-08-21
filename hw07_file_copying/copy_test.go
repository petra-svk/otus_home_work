package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

const fromPath string = "testdata/input.txt"

func TestCopy(t *testing.T) {
	tests := []struct {
		caption string
		limit   int64
		offset  int64
	}{
		{caption: "full copy from one file to another", limit: 0, offset: 0},
		{caption: "copy with limit=10", limit: 10, offset: 0},
		{caption: "copy with limit more then size of file", limit: 10000, offset: 0},
		{caption: "copy with limit and offset", limit: 1000, offset: 100},
		{caption: "copy with big offset", limit: 1000, offset: 6000},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.caption, func(t *testing.T) {
			to, err := os.CreateTemp("", "out.*.txt")
			if err != nil {
				log.Fatal(err)
			}
			defer os.Remove(to.Name())

			offset := tc.offset
			limit := tc.limit
			err = Copy(fromPath, to.Name(), offset, limit)
			require.NoError(t, err)

			comparedFile := fmt.Sprintf("testdata/out_offset%d_limit%d.txt", offset, limit)
			cmd := exec.Command("cmp", to.Name(), comparedFile)
			err = cmd.Run()
			require.NoError(t, err)
		})
	}
}

func TestCopyEdgeCases(t *testing.T) {
	t.Run("offset is more than file's size", func(t *testing.T) {
		to, err := os.CreateTemp("", "out.*.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(to.Name())

		var offset int64 = 10000
		var limit int64
		err = Copy(fromPath, to.Name(), offset, limit)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual error %q", err)
	})

	t.Run("bad file", func(t *testing.T) {
		from := "/dev/null"
		to, err := os.CreateTemp("", "out.*.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(to.Name())

		var offset int64
		var limit int64
		err = Copy(from, to.Name(), offset, limit)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual error %q", err)
	})

	t.Run("bad offset", func(t *testing.T) {
		to, err := os.CreateTemp("", "out.*.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(to.Name())

		var offset int64 = -10
		var limit int64
		err = Copy(fromPath, to.Name(), offset, limit)
		require.ErrorContains(t, err, "limit or offset is less then 0")
	})
}

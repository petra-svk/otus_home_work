package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if limit < 0 || offset < 0 {
		return fmt.Errorf("limit or offset is less then 0: limit=%d, offset=%d", limit, offset)
	}
	source, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("can't open %s: %w", fromPath, err)
	}
	defer source.Close()

	fi, err := source.Stat()
	if err != nil {
		return fmt.Errorf("can't stat %s: %w", fromPath, err)
	}
	sourceSize := fi.Size()
	if sourceSize == 0 || !fi.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	if offset > sourceSize {
		return ErrOffsetExceedsFileSize
	}
	if limit == 0 || sourceSize-offset < limit {
		limit = sourceSize - offset
	}

	dest, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("can't create %s: %w", toPath, err)
	}
	defer dest.Close()
	// create bar
	bar := pb.Full.Start64(limit)
	// create proxy reader
	if offset > 0 {
		source.Seek(offset, io.SeekStart)
	}
	reader := bar.NewProxyReader(source)
	// and copy from reader
	_, err = io.CopyN(dest, reader, limit)
	if err != nil {
		return fmt.Errorf("can't copy from source file %s: %w", fromPath, err)
	}
	bar.Finish()

	return nil
}

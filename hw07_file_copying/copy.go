package main

import (
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrFromOrToPathIsEmpty        = errors.New("from and to paths should be don't empty")
	ErrLimitOrOffsetIsNotPositive = errors.New("limit and offset should be positive numbers")
	ErrUnsupportedFile            = errors.New("unsupported file")
	ErrOffsetExceedsFileSize      = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if err := validateCopyParams(fromPath, toPath, offset, limit); err != nil {
		return err
	}

	fileFrom, err := os.OpenFile(fromPath, os.O_RDONLY, 644)
	if err != nil {
		return err
	}

	defer func(fileFrom *os.File) {
		err := fileFrom.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(fileFrom)

	fileInfo, err := fileFrom.Stat()
	if err != nil {
		return err
	}

	if err := validateFileInfo(fileInfo, offset); err != nil {
		return err
	}

	_, err = fileFrom.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	fileTo, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer func(fileTo *os.File) {
		err := fileTo.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(fileTo)

	if limit == 0 {
		limit = fileInfo.Size()
	}

	barLimit := fileInfo.Size() - offset
	if limit < barLimit {
		barLimit = limit
	}

	bar := pb.New64(barLimit)
	bar.SetRefreshRate(time.Nanosecond)
	bar.Set(pb.Bytes, true)
	bar.Start()

	barReader := bar.NewProxyReader(fileFrom)
	defer bar.Finish()

	_, err = io.CopyN(fileTo, barReader, limit)
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func validateCopyParams(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrFromOrToPathIsEmpty
	}

	if limit < 0 || offset < 0 {
		return ErrLimitOrOffsetIsNotPositive
	}

	return nil
}

func validateFileInfo(fileInfo os.FileInfo, offset int64) error {
	if fileInfo.Size() == 0. {
		return ErrUnsupportedFile
	}

	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	return nil
}

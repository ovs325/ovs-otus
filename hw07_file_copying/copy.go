package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromPath = filepath.Clean(fromPath)
	toPath = filepath.Clean(toPath)
	size, err := IsSrcValid(fromPath)
	if err != nil {
		return err
	}
	if size <= 0 {
		return ErrUnsupportedFile
	}
	if offset > size {
		return ErrOffsetExceedsFileSize
	}
	if offset+limit > size || limit == 0 {
		limit = size - offset
	}
	// Открытие файла-источника
	fromFile, err := os.OpenFile(fromPath, os.O_RDWR, 0o666)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл-источник %s: %s", fromPath, err.Error())
	}
	defer fromFile.Close()
	// Открытие файла-приёмника
	var toFile *os.File
	var dir, name string
	if toPath == fromPath {
		fileName := filepath.Base(toPath)
		dir = toPath[:len(toPath)-len(fileName)]
		name = strings.Split(fileName, ".")[0]
		toFile, err = os.CreateTemp(dir, fmt.Sprintf("%s_", name))
		if err != nil {
			return fmt.Errorf("не удалось создать временный файл %s", err.Error())
		}
	} else {
		toFile, err = os.Create(toPath)
		defer func() {
			_ = toFile.Close()
		}()
		if err != nil {
			return fmt.Errorf("не удалось открыть файл-приемник %s: %s", toPath, err.Error())
		}
	}
	fromFile.Seek(offset, io.SeekStart)
	if err = Copier(toFile, fromFile, offset, limit); err != nil {
		toFile.Close()
		return err
	}
	if toPath == fromPath { // перезаписываем исходный файл
		if err := Renamer(toFile, toPath); err != nil {
			return err
		}
	}
	return nil
}

func Copier(dst, src *os.File, offsetStart, limit int64) (err error) {
	var lenBuf int
	switch {
	case limit <= 10:
		lenBuf = 1
	case limit > 10 && limit <= 1000:
		lenBuf = int(limit) / 10
	case limit > 1000:
		lenBuf = 1024
	}
	buf := make([]byte, lenBuf)
	bar := progressbar.Default(limit / int64(lenBuf))
	count := int64(0)
	offset := int64(0)

	for {
		count++
		bar.Add64(1)
		n, err := src.ReadAt(buf, offset+offsetStart)
		if err != nil {
			if errors.Is(err, io.EOF) {
				_, _ = dst.WriteAt(buf[:n], offset)
				bar.Add64((limit / int64(lenBuf)) - count)
				return nil
			}
			return fmt.Errorf("операция чтения не удалась: %s", err.Error())
		}
		_, err = dst.WriteAt(buf[:n], offset)
		offset += int64(n)
		if err != nil {
			return fmt.Errorf("операция записи не удалась: %s", err.Error())
		}
		if offset >= limit {
			bar.Add64((limit / int64(lenBuf)) - count)
			return nil
		}
	}
}

func IsSrcValid(fromPath string) (size int64, err error) {
	fromStat, err := os.Stat(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return int64(0), fmt.Errorf("файл-источник не существует %s: %s", fromPath, err.Error())
		}
		return int64(0), fmt.Errorf("не удалось обработать статистику файла %s", err.Error())
	}
	if !fromStat.Mode().IsRegular() {
		return int64(0), ErrUnsupportedFile
	}
	return fromStat.Size(), nil
}

func Renamer(file *os.File, toPath string) (err error) {
	tempFileName := file.Name()
	file.Close()
	os.Remove(toPath)
	err = os.Rename(tempFileName, toPath)
	if err != nil {
		return err
	}
	return
}

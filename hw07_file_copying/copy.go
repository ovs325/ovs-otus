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
	bar := progressbar.Default(7)
	fromPath = filepath.Clean(fromPath)
	toPath = filepath.Clean(toPath)
	bar.Add(1)

	size, err := IsSrcValid(fromPath)
	bar.Add(1)
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
	bar.Add(1)
	// Открытие файла-источника
	fromFile, err := os.OpenFile(fromPath, os.O_RDWR, 0o666)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл-источник %s: %s", fromPath, err.Error())
	}
	defer fromFile.Close()
	bar.Add(1)
	fromFile.Seek(offset, io.SeekStart)
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
	bar.Add(1)
	if err = Copier(toFile, fromFile, limit); err != nil {
		toFile.Close()
		return err
	}
	bar.Add(1)
	if toPath == fromPath { // перезаписываем исходный файл
		if err := Renamer(toFile, toPath); err != nil {
			return err
		}
	}
	bar.Add(1)
	return nil
}

func Copier(dst, src *os.File, limit int64) (err error) {
	_, err = io.CopyN(dst, src, limit)
	if err != nil {
		return fmt.Errorf("операция копирования не удалась: %s", err.Error())
	}
	return
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

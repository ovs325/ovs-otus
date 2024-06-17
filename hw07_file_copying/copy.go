package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
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
	fromFile.Seek(offset, io.SeekStart)
	// Открытие файла-приёмника
	var toFile *os.File
	if toPath == fromPath {
		toFile, err = os.CreateTemp("/tmp", "temp_")
		if err != nil {
			return fmt.Errorf("не удалось создать временный файл %s", err.Error())
		}
		defer func() {
			tempFileName := toFile.Name()
			toFile.Close()
			err := os.Remove(tempFileName)
			if err != nil {
				fmt.Printf("не удалось удалить временный файл %s: %s", tempFileName, err.Error())
			}
		}()
	} else {
		toFile, err = os.Create(toPath)
		if err != nil {
			return fmt.Errorf("не удалось открыть файл-приемник %s: %s", toPath, err.Error())
		}
		defer toFile.Close()
	}
	if err = Copier(toFile, fromFile, limit); err != nil {
		return err
	}
	if toPath == fromPath { // перезаписываем исходный файл
		if err = Overwriter(toFile, fromFile, int(limit)); err != nil {
			return err
		}
	}
	return nil
}

func Overwriter(src, dst *os.File, totalProgrBar int) (err error) {
	dst.Truncate(0)
	dst.Seek(int64(0), io.SeekStart)
	src.Seek(int64(0), io.SeekStart)
	bar := pb.StartNew(totalProgrBar)
	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("операция перезаписи исходного файла не удалась: %s", err.Error())
	}
	bar.FinishPrint("File overwritten successfully!")
	return
}

func Copier(dst, src *os.File, limit int64) (err error) {
	bar := pb.StartNew(int(limit))
	_, err = io.CopyN(dst, src, limit)
	if err != nil {
		return fmt.Errorf("операция копирования не удалась: %s", err.Error())
	}
	bar.FinishPrint("File copied successfully!")
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

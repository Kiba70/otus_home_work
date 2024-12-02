package main

import (
	"errors"
	"io"
	"io/fs"
	"os"

	"github.com/cheggaaa/pb/v3" //nolint: depguard
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrIncorrectParameter    = errors.New("incorrect parameter")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	var finfo os.FileInfo

	// На всякий случай
	if offset < 0 || limit < 0 || len(fromPath) == 0 || len(toPath) == 0 {
		return ErrIncorrectParameter
	}

	if !fs.ValidPath(fromPath) || !fs.ValidPath(toPath) {
		return ErrUnsupportedFile
	}

	// Откуда копируем
	fin, err := os.Open(fromPath)
	if err != nil {
		return errors.Join(ErrUnsupportedFile, err)
	}
	defer fin.Close()

	// offset подходит?
	if finfo, err = fin.Stat(); err != nil {
		return err
	}
	if finfo.Size() <= offset {
		return ErrOffsetExceedsFileSize
	}

	// Куда копируем
	fout, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer fout.Close()

	// Встаём на начало чтения
	if offset > 0 {
		if _, err := fin.Seek(offset, io.SeekStart); err != nil {
			return err
		}
	}

	if limit == 0 { // Весь файл
		limit = finfo.Size() - offset
	} else {
		limit = min(finfo.Size()-offset, limit)
	}

	// Делаем прогресc бар
	bar := pb.Default.Start64(limit)
	barReader := bar.NewProxyReader(fin)
	defer bar.Finish()

	// Копируем
	if _, err = io.CopyN(fout, barReader, limit); err != nil {
		return err
	}

	return nil
}

package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Params struct {
	fromPath, toPath string
	offset, limit    int64
}

type Expected struct {
	err      error
	testPath string
}

func TestCopy(t *testing.T) {
	cases := []struct {
		name     string
		params   Params
		expected Expected
	}{
		{
			name: "offset0_limit0",
			params: Params{
				fromPath: "./testdata/input.txt",
				toPath:   "./testdata/out.txt",
				offset:   int64(0),
				limit:    int64(0),
			},
			expected: Expected{testPath: "./testdata/out_offset0_limit0.txt"},
		},
		{
			name: "out_offset0_limit10",
			params: Params{
				fromPath: "./testdata/input.txt",
				toPath:   "./testdata/out.txt",
				offset:   int64(0),
				limit:    int64(10),
			},
			expected: Expected{testPath: "./testdata/out_offset0_limit10.txt"},
		},
		{
			name: "out_offset0_limit1000",
			params: Params{
				fromPath: "./testdata/input.txt",
				toPath:   "./testdata/out.txt",
				offset:   int64(0),
				limit:    int64(1000),
			},
			expected: Expected{testPath: "./testdata/out_offset0_limit1000.txt"},
		},
		{
			name: "out_offset0_limit10000",
			params: Params{
				fromPath: "./testdata/input.txt",
				toPath:   "./testdata/out.txt",
				offset:   int64(0),
				limit:    int64(10000),
			},
			expected: Expected{testPath: "./testdata/out_offset0_limit10000.txt"},
		},
		{
			name: "out_offset100_limit1000",
			params: Params{
				fromPath: "./testdata/input.txt",
				toPath:   "./testdata/out.txt",
				offset:   int64(100),
				limit:    int64(1000),
			},
			expected: Expected{testPath: "./testdata/out_offset100_limit1000.txt"},
		},
		{
			name: "out_offset6000_limit1000",
			params: Params{
				fromPath: "./testdata/input.txt",
				toPath:   "./testdata/out.txt",
				offset:   int64(6000),
				limit:    int64(1000),
			},
			expected: Expected{testPath: "./testdata/out_offset6000_limit1000.txt"},
		},
		{
			name: "ExceedingLimit",
			params: Params{
				fromPath: "./testdata/input.txt",
				toPath:   "./testdata/out.txt",
				offset:   int64(0),
				limit:    int64(1000000000),
			},
			expected: Expected{testPath: "./testdata/out_offset0_limit0.txt"},
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			errRes := Copy(cs.params.fromPath, cs.params.toPath, cs.params.offset, cs.params.limit)

			assert.ErrorIs(t, errRes, cs.expected.err)
			assert.FileExists(t, cs.params.toPath)
			resStat, _ := os.Stat(cs.params.toPath)
			expectedStat, _ := os.Stat(cs.expected.testPath)
			assert.Equal(t, resStat.Size(), expectedStat.Size())
			assert.Equalf(t, GetTxtContent(cs.params.toPath), GetTxtContent(cs.expected.testPath), "Плохая копия")
			if _, err := IsSrcValid(cs.params.toPath); err == nil {
				err := os.Remove(cs.params.toPath)
				assert.NoErrorf(t, err, "Ошибка удаления временного файла результата теста")
			}
		})
	}
}

func TestCopy_Errors(t *testing.T) {
	cases := []struct {
		name     string
		params   Params
		expected Expected
	}{
		{
			name: "ErrOffsetExceedsFileSize",
			params: Params{
				fromPath: "./testdata/input.txt",
				toPath:   "./testdata/out.txt",
				offset:   int64(1000000),
				limit:    int64(100),
			},
			expected: Expected{err: ErrOffsetExceedsFileSize},
		},
		{
			name: "ErrUnsupportedFile",
			params: Params{
				fromPath: "/dev/urandom",
				toPath:   "./testdata/out.txt",
				offset:   int64(100),
				limit:    int64(100),
			},
			expected: Expected{err: ErrUnsupportedFile},
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			errRes := Copy(cs.params.fromPath, cs.params.toPath, cs.params.offset, cs.params.limit)

			assert.ErrorIs(t, errRes, cs.expected.err)
			assert.NoFileExists(t, cs.params.toPath)
			if _, err := IsSrcValid(cs.params.toPath); err == nil {
				err := os.Remove(cs.params.toPath)
				assert.NoErrorf(t, err, "Ошибка удаления временного файла результата теста")
			}
		})
	}
}

func GetTxtContent(path string) (content string) {
	file, _ := os.OpenFile(path, os.O_RDONLY, 0o666)
	defer file.Close()
	res, _ := io.ReadAll(file)
	return string(res)
}

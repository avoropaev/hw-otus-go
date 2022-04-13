package main

import (
	"bufio"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		offset           int64
		limit            int64
		expectedFileName string
	}{
		{
			offset:           0,
			limit:            0,
			expectedFileName: "out_offset0_limit0",
		},
		{
			offset:           0,
			limit:            10,
			expectedFileName: "out_offset0_limit10",
		},
		{
			offset:           0,
			limit:            1000,
			expectedFileName: "out_offset0_limit1000",
		},
		{
			offset:           0,
			limit:            10000,
			expectedFileName: "out_offset0_limit10000",
		},
		{
			offset:           100,
			limit:            1000,
			expectedFileName: "out_offset100_limit1000",
		},
		{
			offset:           6000,
			limit:            1000,
			expectedFileName: "out_offset6000_limit1000",
		},
	}

	for _, val := range tests {
		t.Run(val.expectedFileName, func(t *testing.T) {
			temp, err := os.CreateTemp(os.TempDir(), "test_copy")
			require.NoError(t, err)

			err = Copy("testdata/input.txt", temp.Name(), val.offset, val.limit)
			require.NoError(t, err)

			expectedFile, err := os.OpenFile("testdata/"+val.expectedFileName+".txt", os.O_RDONLY, 644)
			require.NoError(t, err)

			sc1 := bufio.NewScanner(expectedFile)
			sc2 := bufio.NewScanner(temp)

			for {
				sc1Bool := sc1.Scan()
				sc2Bool := sc2.Scan()
				if !sc1Bool && !sc2Bool {
					break
				}

				require.Equal(t, sc1.Text(), sc2.Text())
			}

			err = expectedFile.Close()
			require.NoError(t, err)

			err = temp.Close()
			require.NoError(t, err)
			err = os.Remove(temp.Name())
			require.NoError(t, err)
		})
	}
}

func TestCopyError(t *testing.T) {
	temp, err := os.CreateTemp(os.TempDir(), "test_copy")
	require.NoError(t, err)

	tests := []struct {
		name          string
		from          string
		to            string
		offset        int64
		limit         int64
		expectedError error
	}{
		{
			name:          "from is empty",
			from:          "",
			to:            temp.Name(),
			expectedError: ErrFromOrToPathIsEmpty,
		},
		{
			name:          "to is empty",
			from:          "testdata/size100.txt",
			to:            "",
			expectedError: ErrFromOrToPathIsEmpty,
		},
		{
			name:          "from and to are empty",
			from:          "",
			to:            "",
			expectedError: ErrFromOrToPathIsEmpty,
		},
		{
			name:          "offset is not positive",
			from:          "testdata/size100.txt",
			to:            temp.Name(),
			offset:        -1,
			expectedError: ErrLimitOrOffsetIsNotPositive,
		},
		{
			name:          "limit is not positive",
			from:          "testdata/size100.txt",
			to:            temp.Name(),
			limit:         -1,
			expectedError: ErrLimitOrOffsetIsNotPositive,
		},
		{
			name:          "offset and limit is not positive",
			from:          "testdata/size100.txt",
			to:            temp.Name(),
			offset:        -1,
			limit:         -1,
			expectedError: ErrLimitOrOffsetIsNotPositive,
		},
		{
			name:          "unsupported file",
			from:          "/dev/urandom",
			to:            temp.Name(),
			expectedError: ErrUnsupportedFile,
		},
		{
			name:          "offset exceeds file size",
			from:          "testdata/size100.txt",
			to:            temp.Name(),
			offset:        101,
			expectedError: ErrOffsetExceedsFileSize,
		},
		{
			name:          "offset exceeds file size",
			from:          "testdata/not_exists.txt",
			to:            temp.Name(),
			expectedError: fs.ErrNotExist,
		},
	}

	for _, val := range tests {
		t.Run(val.name, func(t *testing.T) {
			err = Copy(val.from, val.to, val.offset, val.limit)
			require.ErrorIs(t, err, val.expectedError)
		})
	}

	err = temp.Close()
	require.NoError(t, err)
	err = os.Remove(temp.Name())
	require.NoError(t, err)
}

// Copyright 2019 Huawei Technologies Co.,Ltd.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use
// this file except in compliance with the License.  You may obtain a copy of
// the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package obs

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// partSlice Swap Tests

func TestPartSlice_Swap_ShouldSwapElements_WhenValidIndices(t *testing.T) {
	parts := partSlice{
		{PartNumber: 1, ETag: "etag1"},
		{PartNumber: 2, ETag: "etag2"},
		{PartNumber: 3, ETag: "etag3"},
	}

	parts.Swap(0, 2)

	assert.Equal(t, 3, parts[0].PartNumber)
	assert.Equal(t, "etag3", parts[0].ETag)
	assert.Equal(t, 2, parts[1].PartNumber)
	assert.Equal(t, 1, parts[2].PartNumber)
	assert.Equal(t, "etag1", parts[2].ETag)
}

func TestPartSlice_Swap_ShouldSwapWithSameElement_WhenSameIndices(t *testing.T) {
	parts := partSlice{
		{PartNumber: 1, ETag: "etag1"},
		{PartNumber: 2, ETag: "etag2"},
	}

	parts.Swap(1, 1)

	assert.Equal(t, 2, parts[1].PartNumber)
	assert.Equal(t, "etag2", parts[1].ETag)
}

// readerWrapper seek Tests

func TestReaderWrapper_Seek_ShouldSeek_WhenReaderIsStringsReader(t *testing.T) {
	reader := strings.NewReader("Hello, World!")
	rw := &readerWrapper{
		reader:     reader,
		totalCount: -1, // Read all
	}

	newPos, err := rw.seek(7, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(7), newPos)

	buf := make([]byte, 5)
	n, _ := reader.Read(buf)
	assert.Equal(t, "World", string(buf[:n]))
}

func TestReaderWrapper_Seek_ShouldSeek_WhenReaderIsBytesReader(t *testing.T) {
	data := []byte("Hello, World!")
	reader := bytes.NewReader(data)
	rw := &readerWrapper{
		reader:     reader,
		totalCount: -1,
	}

	newPos, err := rw.seek(0, 2) // Seek from end
	assert.NoError(t, err)
	assert.Equal(t, int64(13), newPos)
}

func TestReaderWrapper_Seek_ShouldSeek_WhenReaderIsFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test-seek")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("Hello, World!")
	require.NoError(t, err)
	tempFile.Close()

	reader, err := os.Open(tempFile.Name())
	require.NoError(t, err)
	defer reader.Close()

	rw := &readerWrapper{
		reader:     reader,
		totalCount: -1,
	}

	newPos, err := rw.seek(7, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(7), newPos)
}

func TestReaderWrapper_Seek_ShouldReturnOffset_WhenReaderNotSeekable(t *testing.T) {
	reader := &bytes.Buffer{}
	rw := &readerWrapper{
		reader:     reader,
		totalCount: -1,
	}

	newPos, err := rw.seek(10, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), newPos)
}

// readerWrapper Read Tests

func TestReaderWrapper_Read_ShouldReturnEOF_WhenTotalCountIsZero(t *testing.T) {
	reader := strings.NewReader("Hello, World!")
	rw := &readerWrapper{
		reader:     reader,
		totalCount: 0,
	}

	buf := make([]byte, 100)
	n, err := rw.Read(buf)

	assert.Equal(t, 0, n)
	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
}

func TestReaderWrapper_Read_ShouldReadExactlyTotalCount_WhenBufferLarger(t *testing.T) {
	reader := strings.NewReader("Hello, World!")
	rw := &readerWrapper{
		reader:     reader,
		totalCount: 5, // Only read 5 bytes
	}

	buf := make([]byte, 100)
	n, err := rw.Read(buf)

	assert.Equal(t, 5, n)
	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, "Hello", string(buf[:n]))
}

func TestReaderWrapper_Read_ShouldReadAll_WhenTotalCountIsNegative(t *testing.T) {
	reader := strings.NewReader("Hello, World!")
	rw := &readerWrapper{
		reader:     reader,
		totalCount: -1, // Read all
	}

	buf := make([]byte, 100)
	n, err := rw.Read(buf)

	assert.Equal(t, 13, n)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", string(buf[:n]))
}

func TestReaderWrapper_Read_ShouldReadInChunks_WhenReadingMultipleTimes(t *testing.T) {
	reader := strings.NewReader("Hello, World!")
	rw := &readerWrapper{
		reader:     reader,
		totalCount: 8, // Only read 8 bytes total
	}

	// First read 5 bytes
	buf1 := make([]byte, 5)
	n1, err1 := rw.Read(buf1)
	assert.Equal(t, 5, n1)
	assert.NoError(t, err1)
	assert.Equal(t, "Hello", string(buf1[:n1]))

	// Second read, should get remaining 3 bytes and EOF
	buf2 := make([]byte, 10)
	n2, err2 := rw.Read(buf2)
	assert.Equal(t, 3, n2)
	assert.Error(t, err2)
	assert.Equal(t, io.EOF, err2)
	assert.Equal(t, ", W", string(buf2[:n2]))
}

func TestReaderWrapper_Read_ShouldReturnZero_WhenReaderReturnsZeroAndNoError(t *testing.T) {
	// Create a reader that returns (0, nil) on first read
	reader := &mockZeroReader{}
	rw := &readerWrapper{
		reader:     reader,
		totalCount: -1,
	}

	buf := make([]byte, 100)
	n, err := rw.Read(buf)

	assert.Equal(t, 0, n)
	assert.NoError(t, err)
}

// mockZeroReader simulates a reader that returns (0, nil)
type mockZeroReader struct{}

func (m *mockZeroReader) Read(p []byte) (n int, err error) {
	return 0, nil
}

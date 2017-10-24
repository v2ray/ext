package build

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileMetadata struct {
	name     string
	size     int
	checksum string
}

func (meta *FileMetadata) Checksum() string {
	return meta.checksum
}

func GenerateFileMetadata(file string) (*FileMetadata, error) {
	fileReader, err := os.Open(file)
	if err != nil {
		return nil, newError("failed to open file: ", file).Base(err)
	}
	defer fileReader.Close()

	hasher := sha1.New()
	nBytes, err := io.Copy(hasher, fileReader)
	if err != nil {
		return nil, newError("failed to read file: ", file).Base(err)
	}
	sha1sum := hasher.Sum(nil)
	return &FileMetadata{
		name:     filepath.Base(file),
		size:     int(nBytes),
		checksum: hex.EncodeToString(sha1sum),
	}, nil
}

type FileMetadataWriter struct {
	writer *os.File
}

func NewFileMetadataWriter(file string) (*FileMetadataWriter, error) {
	writer, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, newError("failed to open metadata file: ", file).Base(err)
	}
	return &FileMetadataWriter{
		writer: writer,
	}, nil
}

func (w *FileMetadataWriter) Append(metadata *FileMetadata) {
	fmt.Fprintf(w.writer, "File: %s\n", metadata.name)
	fmt.Fprintf(w.writer, "Size: %d\n", metadata.size)
	fmt.Fprintf(w.writer, "SHA1: %s\n", metadata.checksum)
	fmt.Fprintln(w.writer)
}

func (w *FileMetadataWriter) Close() {
	w.writer.Close()
}

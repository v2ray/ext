package zip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

type ZipWorker struct {
	zipWriter *zip.Writer
	root      string
}

func NewZipWorker(zipFile io.Writer, root string) *ZipWorker {
	return &ZipWorker{
		zipWriter: zip.NewWriter(zipFile),
		root:      root,
	}
}

func (worker *ZipWorker) run() error {
	defer worker.close()
	return filepath.Walk(worker.root, worker.zipAllFiles)
}

func (worker *ZipWorker) zipAllFiles(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	fileReader, err := os.Open(path)
	if err != nil {
		return err
	}
	fileInfo, err := fileReader.Stat()
	if err != nil {
		return err
	}
	fileHeader, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}
	fileHeader.Name = path
	fileHeader.Method = zip.Deflate
	fileWriter, err := worker.zipWriter.CreateHeader(fileHeader)
	if err != nil {
		return err
	}
	_, err = io.Copy(fileWriter, fileReader)
	return err
}

func (worker *ZipWorker) close() {
	worker.zipWriter.Close()
}

func goZipFolder(folder string, file string) error {
	if _, err := os.Stat(file); err == nil {
		os.Remove(file)
	}

	zipFile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	return NewZipWorker(zipFile, folder).run()
}

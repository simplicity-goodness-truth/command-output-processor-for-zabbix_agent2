package file

import (
	"bytes"
	"io"
	"linuxProcessesCounterStandalone/check"
	"os"
)

// Class definition

type File struct {
	FilePath    string
	FileContent []byte
}

// Constructor

func NewFile(filePath string) *File {

	return &File{

		FilePath:    filePath,
		FileContent: readFile(filePath),
	}

}

func NewFileReader(filePath string) io.Reader {

	 return bytes.NewReader(readFile(filePath))
	
}

// Private methods

func readFile(filePath string) []byte {

	var fileContent []byte

	fileContent, err := os.ReadFile(filePath)

	check.Check(err)

	return fileContent

}
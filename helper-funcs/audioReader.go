package helperfuncs

import (
	"fmt"
	"io"
	"os"
)

func readAudioFile(filename string) ([]byte, error) {

	if filename == "" {
		filename = "record.mp3" // default value
	}

	file, err := os.Open(filename)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	limitedReader := io.LimitReader(file, 1024*1024) // 1MB limit
	fileBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return []byte{}, err
	}

	fmt.Printf("File size: %d bytes\n", len(fileBytes))
	return fileBytes, nil
}

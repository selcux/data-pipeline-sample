package internal

import (
	"bufio"
	"os"
	"time"
)

type FileOps struct {
	FilePath string
}

func NewFileOps(filePath string) *FileOps {
	return &FileOps{filePath}
}

func (fo *FileOps) ReadInterval(interval time.Duration, lineCh chan<- string, errorCh chan<- error, closeCh chan<- struct{}) {
	file, err := os.Open(fo.FilePath)
	if err != nil {
		close(lineCh)
		errorCh <- err
		return
	}
	defer file.Close()
	defer close(closeCh)

	scanner := bufio.NewScanner(file)
	for range time.Tick(interval) {
		if scanner.Scan() {
			line := scanner.Text()
			lineCh <- line
		} else {
			err = scanner.Err()
			if err != nil {
				close(lineCh)
				errorCh <-err
			}

			return
		}
	}
}

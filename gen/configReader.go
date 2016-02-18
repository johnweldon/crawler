package gen

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
)

type configFileReader struct {
	fpath string
	out   chan string
}

func (r *configFileReader) Start() <-chan string {
	go func() {
		defer close(r.out)
		fd, err := os.Open(r.fpath)
		if err != nil {
			log.Fatal(err)
		}
		defer fd.Close()
		s := bufio.NewScanner(fd)
		for s.Scan() {
			r.out <- s.Text()
		}
		if err := s.Err(); err != nil {
			log.Fatal(err)
		}
	}()
	return r.out
}

func NewConfigFileReader(configFile string) (URLGenerator, error) {
	var gen URLGenerator

	clean := filepath.Clean(configFile)
	_, err := os.Stat(clean)
	if err != nil {
		return gen, err
	}

	gen = &configFileReader{fpath: clean, out: make(chan string)}
	return gen, nil
}

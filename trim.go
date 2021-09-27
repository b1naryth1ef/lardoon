package lardoon

import (
	"bytes"
	"os"

	"github.com/b1naryth1ef/jambon/tacview"
)

func trimTacView(path string, start, end int) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader, err := tacview.NewParser(file)
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer([]byte{})
	writer := tacview.NewRawWriter(buff)

	err = tacview.TrimRaw(reader, writer, float64(start), float64(end))
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

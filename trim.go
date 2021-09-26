package lardoon

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/b1naryth1ef/jambon/tacview"
)

type CloseBuffer struct {
	*bytes.Buffer
}

func (c *CloseBuffer) Close() error {
	return nil
}

func trimTacView(path string, start, end int) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader, err := tacview.NewReader(file)
	if err != nil {
		return nil, err
	}

	readerDone := make(chan error)
	done := make(chan error)
	timeFrames := make(chan *tacview.TimeFrame)

	go func() {
		defer close(readerDone)
		err := reader.ProcessTimeFrames(1, timeFrames)
		if err != nil {
			readerDone <- err
		}
	}()

	// Stores any objects which where created before the start of our time window,
	// and have not yet been destroyed.
	preStartObjects := make(map[uint64]*tacview.Object)
	var firstTimeFrame *tacview.TimeFrame

	log.Printf("Scanning to offset %d...\n", start)
	for {
		tf, ok := <-timeFrames
		if !ok {
			return nil, io.EOF
		}

		if int(tf.Offset) >= start {
			firstTimeFrame = tf
			break
		}

		for _, object := range tf.Objects {
			existingObject := preStartObjects[object.Id]

			if existingObject == nil && !object.Deleted {
				typeProp := object.Get("Type")
				if typeProp == nil || strings.Contains(typeProp.Value, "Misc") {
					continue
				}
				preStartObjects[object.Id] = object
				continue
			} else if object.Deleted {
				delete(preStartObjects, object.Id)
			}

			for _, newProp := range object.Properties {
				existingObject.Set(newProp.Key, newProp.Value)
			}
		}
	}

	log.Printf("Collected %d active objects for frame 0\n", len(preStartObjects))

	referenceTime := reader.Header.ReferenceTime.Add(time.Second * time.Duration(start))

	// We copy the initial time frame completely
	initialTimeFrame := tacview.NewTimeFrame()
	initialTimeFrame.Offset = 0
	initialTimeFrame.Objects = reader.Header.InitialTimeFrame.Objects

	// We generate a frame 0 that contains any objects which existed at the start
	// point of our time window.
	preTimeFrame := tacview.NewTimeFrame()
	preTimeFrame.Offset = 0
	for _, object := range preStartObjects {
		preTimeFrame.Objects = append(firstTimeFrame.Objects, object)
	}

	header := &tacview.Header{
		FileType:         reader.Header.FileType,
		FileVersion:      reader.Header.FileVersion,
		ReferenceTime:    referenceTime,
		InitialTimeFrame: *initialTimeFrame,
	}

	buff := &CloseBuffer{bytes.NewBuffer([]byte{})}

	writer, err := tacview.NewWriter(buff, header)
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	err = writer.WriteTimeFrame(preTimeFrame)
	if err != nil {
		return nil, err
	}

	err = writer.WriteTimeFrame(firstTimeFrame)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(done)
		defer close(readerDone)

		for {
			tf, ok := <-timeFrames

			if !ok {
				return
			}

			if int(tf.Offset) >= end {
				return
			}

			tf.Offset = tf.Offset - float64(start)

			err := writer.WriteTimeFrame(tf)
			if err != nil {
				done <- err
				return
			}
		}
	}()

	err = <-done
	if err != nil {
		return nil, err
	}

	err = <-readerDone
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

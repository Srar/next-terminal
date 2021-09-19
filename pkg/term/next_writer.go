package term

import (
	"errors"
)

var ErrNextWriterHasBeenClosed = errors.New("the nextWriter has been closed. ")

type NextWriter struct {
	isClosed bool
	buffer   chan []byte
}

func NewNextWriter() *NextWriter {
	return &NextWriter{
		isClosed: false,
		buffer:   make(chan []byte, 16),
	}
}

func (w *NextWriter) Write(p []byte) (bytes int, err error) {
	// fast path
	if w.isClosed {
		return 0, ErrNextWriterHasBeenClosed
	}

	defer func() {
		if recover() != nil {
			bytes = 0
			err = ErrNextWriterHasBeenClosed
		}
	}()

	w.buffer <- p

	return len(p), nil
}

func (w *NextWriter) Read() ([]byte, int, error) {
	// fast path
	if w.isClosed {
		return nil, 0, ErrNextWriterHasBeenClosed
	}

	b, opened := <-w.buffer
	if !opened {
		return nil, 0, ErrNextWriterHasBeenClosed
	}

	return b, len(b), nil
}

func (w *NextWriter) Close() {
	if !w.isClosed {
		w.isClosed = true
		close(w.buffer)
	}
}

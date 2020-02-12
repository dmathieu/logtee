package logtee

import (
	"bufio"
	"io"
	"os"
	"sync"
)

// Handler is an interface which allows setting up generic handlers executed
// whenever a line is received
type Handler interface {
	Handle([]byte) error
}

// Streamer allows configuring, running and stopping a log streamer.
// The handler will be executed for each log line
type Streamer struct {
	input   io.ReadCloser
	handler Handler
	mux     sync.Mutex
	closed  bool
}

// NewStreamer creates a new instance of streamer from a file path
func NewStreamer(path string, handler Handler) (*Streamer, error) {
	input, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &Streamer{
		input:   input,
		handler: handler,
		closed:  false,
	}, nil
}

// Close closes the opened streamer
func (s *Streamer) Close() {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.input.Close()
	s.closed = true
}

func (s *Streamer) isClosed() bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.closed
}

// Run scans every line in the streamer and executes the handler with them
func (s *Streamer) Run() error {
	scanner := bufio.NewScanner(s.input)

	for !s.isClosed() {
		if scanner.Scan() == false {
			scanner = bufio.NewScanner(s.input)
			continue
		}

		err := s.handler.Handle(scanner.Bytes())
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}

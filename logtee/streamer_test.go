package logtee

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testHandler struct {
	content [][]byte
}

func (t *testHandler) Handle(c []byte) error {
	t.content = append(t.content, c)
	return nil
}

func TestNewStreamer(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "streamer")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	s, err := NewStreamer(tmpFile.Name(), &testHandler{})
	defer s.Close()

	assert.NoError(t, err)
	assert.NotNil(t, s)
}

func TestNewStreamerMissingFile(t *testing.T) {
	s, err := NewStreamer("/foo/bar", &testHandler{})

	assert.Error(t, err)
	assert.Equal(t, "open /foo/bar: no such file or directory", err.Error())
	assert.Nil(t, s)
}

func TestRunStreamer(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "streamer")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	handler := &testHandler{}
	s, err := NewStreamer(tmpFile.Name(), handler)

	assert.NoError(t, err)

	go func() {
		time.Sleep(time.Millisecond)
		_, err := tmpFile.Write([]byte("foo bar"))
		assert.NoError(t, err)
		time.Sleep(time.Millisecond)
		s.Close()
	}()

	_, err = tmpFile.Write([]byte("hello world"))
	assert.NoError(t, err)
	err = s.Run()
	assert.NoError(t, err)

	assert.Equal(t, [][]byte{[]byte("hello world"), []byte("foo bar")}, handler.content)
}

package logtee

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testParsedHandler struct {
	content []LogLine
}

func (t *testParsedHandler) Handle(l LogLine) error {
	t.content = append(t.content, l)
	return nil
}

func TestParseLogLines(t *testing.T) {
	parsedHandler := &testParsedHandler{}
	handler := Parser{handler: parsedHandler}

	jamesTime, err := time.Parse(timeLayout, "09/May/2018:16:00:39 +0000")
	assert.NoError(t, err)
	jillTime, err := time.Parse(timeLayout, "09/May/2018:16:00:41 +0000")
	assert.NoError(t, err)
	frankTime, err := time.Parse(timeLayout, "09/May/2018:16:00:42 +0000")
	assert.NoError(t, err)
	maryTime, err := time.Parse(timeLayout, "09/May/2018:16:00:42 +0000")
	assert.NoError(t, err)

	for _, test := range []struct {
		Name    string
		Line    []byte
		LogLine LogLine
		Ok      bool
	}{
		{
			Name: "with an invalid log line",
			Line: []byte("hello world"),
			Ok:   false,
		},
		{
			Name: "with a valid log line to report",
			Line: []byte("127.0.0.1 - james [09/May/2018:16:00:39 +0000] \"GET /report HTTP/1.0\" 200 123"),
			LogLine: LogLine{
				IPAddress: "127.0.0.1",
				User:      "james",
				Time:      jamesTime,
				Method:    "GET",
				Path:      "/report",
				Status:    200,
				Length:    123,
			},
			Ok: true,
		},
		{
			Name: "with a valid log line to get api user",
			Line: []byte("127.0.0.1 - jill [09/May/2018:16:00:41 +0000] \"GET /api/user HTTP/1.0\" 200 234"),
			LogLine: LogLine{
				IPAddress: "127.0.0.1",
				User:      "jill",
				Time:      jillTime,
				Method:    "GET",
				Path:      "/api/user",
				Status:    200,
				Length:    234,
			},
			Ok: true,
		},
		{
			Name: "with a valid log line to post api user",
			Line: []byte("127.0.0.1 - frank [09/May/2018:16:00:42 +0000] \"POST /api/user HTTP/1.0\" 200 34"),
			LogLine: LogLine{
				IPAddress: "127.0.0.1",
				User:      "frank",
				Time:      frankTime,
				Method:    "POST",
				Path:      "/api/user",
				Status:    200,
				Length:    34,
			},
			Ok: true,
		},
		{
			Name: "with a valid log line to post api user and error response",
			Line: []byte("127.0.0.1 - mary [09/May/2018:16:00:42 +0000] \"POST /api/user HTTP/1.0\" 503 12"),
			LogLine: LogLine{
				IPAddress: "127.0.0.1",
				User:      "mary",
				Time:      maryTime,
				Method:    "POST",
				Path:      "/api/user",
				Status:    503,
				Length:    12,
			},
			Ok: true,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			beforeLen := len(parsedHandler.content)
			err := handler.Handle(test.Line)
			afterLen := len(parsedHandler.content)

			if test.Ok {
				assert.NoError(t, err)
				assert.Equal(t, beforeLen+1, afterLen)
				assert.Equal(t, test.LogLine, parsedHandler.content[afterLen-1])
				assert.Equal(t, int64(0), invalidLogLineMetric.Collect(time.Now()))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, beforeLen, afterLen)
				assert.Equal(t, int64(1), invalidLogLineMetric.Collect(time.Now().Add(-time.Second)))
			}
		})
	}
}

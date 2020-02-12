package logtee

import (
	"regexp"
	"strconv"
	"time"

	"github.com/dmathieu/logtee/metrics"
)

var (
	invalidLogLineMetric = metrics.RegistryProvider().Counter("logtee.invalid_line")
	parseRegex           = regexp.MustCompile("^(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}) - ([a-z0-9]+) \\[(.*)\\] \"([A-Z]+) (.*) HTTP/1.0\" ([0-9]+) ([0-9]+)$")
)

const timeLayout = "02/Jan/2006:15:04:05 -0700"

// LogLine represents a single log line after it has been parsed
type LogLine struct {
	IPAddress string
	User      string
	Time      time.Time
	Method    string
	Path      string
	Status    int
	Length    int
}

// ParsedHandler is an interface which allows setting up generic handlers executed whenever a line has been parsed
type ParsedHandler interface {
	Handle(LogLine) error
}

// Parser is an `Handler` which gets the line as a byte array, turns the
// content into a `LogLine` and passes it to its child handler
type Parser struct {
	handler ParsedHandler
}

// NewParser creates a new instance of Parser
func NewParser(h ParsedHandler) Parser {
	return Parser{h}
}

// Handle receives a log line as a byte array, parses it and passes it to the
// child handler
func (p Parser) Handle(c []byte) error {
	m := parseRegex.FindSubmatch(c)
	if m == nil {
		invalidLogLineMetric.Increment(time.Now())
		return nil
	}

	t, err := time.Parse(timeLayout, string(m[3]))
	if err != nil {
		return err
	}

	status, err := strconv.Atoi(string(m[6]))
	if err != nil {
		return err
	}

	length, err := strconv.Atoi(string(m[7]))
	if err != nil {
		return err
	}

	l := LogLine{
		IPAddress: string(m[1]),
		User:      string(m[2]),
		Time:      t,
		Method:    string(m[4]),
		Path:      string(m[5]),
		Status:    status,
		Length:    length,
	}
	p.handler.Handle(l)

	return nil
}

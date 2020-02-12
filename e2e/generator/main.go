package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	logFile           string
	requestsPerSecond = 100000
	methods           = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	statuses          = []int{200, 201, 404, 500, 503}
	paths             = []string{
		"/report",
		"/report/1",
		"/api/user",
		"/api",
	}
)

const timeLayout = "02/Jan/2006:15:04:05 -0700"

func init() {
	rand.Seed(time.Now().Unix())

	flag.StringVar(&logFile, "file", "/tmp/access.log", "Path to the log file")
	flag.Parse()
}

func main() {
	f, err := os.Create(logFile)
	if err != nil {
		log.Fatal(err)
	}

	for {
		for i := 1; i <= requestsPerSecond; i++ {
			_, err := f.Write([]byte(randomLogLine()))
			if err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(time.Second)
	}
}

func randomLogLine() string {
	now := time.Now()

	return fmt.Sprintf(
		"127.0.0.1 - james [%s] \"%s %s HTTP/1.0\" %d %d\n",
		now.Format(timeLayout),
		randomMethod(),
		randomPath(),
		randomHTTPStatus(),
		randomSize(),
	)
}

func randomMethod() string {
	return methods[rand.Intn(len(methods))]
}

func randomPath() string {
	return paths[rand.Intn(len(paths))]
}

func randomHTTPStatus() int {
	return statuses[rand.Intn(len(statuses))]
}

func randomSize() int {
	return rand.Intn(10000) + 10
}

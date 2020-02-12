# Logtee

I have been writing this as an interview exercise.  
Since I've decided not to pursue interviewing with that company and the
exercise output is quite generic, I figured I would make it public.

## Exercise

Consume an actively written-to w3c-formatted HTTP access log
(https://www.w3.org/Daemon/User/Config/Logging.html). It should default to
reading /tmp/access.log and be overrideable

Example log lines:

```
127.0.0.1 - james [09/May/2018:16:00:39 +0000] "GET /report HTTP/1.0" 200 123
127.0.0.1 - jill [09/May/2018:16:00:41 +0000] "GET /api/user HTTP/1.0" 200 234
127.0.0.1 - frank [09/May/2018:16:00:42 +0000] "POST /api/user HTTP/1.0" 200 34
127.0.0.1 - mary [09/May/2018:16:00:42 +0000] "POST /api/user HTTP/1.0" 503 12
```


* Display stats every 10s about the traffic during those 10s: the sections of
  the web site with the most hits, as well as interesting summary statistics on
  the traffic as a whole. A section is defined as being what's before the
  second '/' in the resource section of the log line. For example, the section
  for "/pages/create" is "/pages"
* Make sure a user can keep the app running and monitor the log file continuously
* Whenever total traffic for the past 2 minutes exceeds a certain number on
  average, print or display a message saying that “High traffic generated an
  alert - hits = {value}, triggered at {time}”. The default threshold should be
  10 requests per second, and should be overridable
* Whenever the total traffic drops again below that value on average for the
  past 2 minutes, print or display another message detailing when the alert
  recovered
* Write a test for the alerting logic
* Explain how you’d improve on this application design

## Setup

This project is written in Go. The package name is
`github.com/dmathieu/logtee`.  So it expects to be located at
`$GOPATH/github.com/dmathieu/logtee`.

There is no specific setup required, as there are no dependencies outside of
having Go installed locally.  This project has been tested against Go 1.13.7,
but should be compatible with any recent version.

### Running the tests

```
make test
```

## Usage

There is a single CLI binary to execute:

```
go run cmd/stream/main.go
```

The command can take two optional arguments:

* `-alert-total-requests` int - Total requests threshold after which to send alert (default 10)
* `-file` string - Path to the log file (default "/tmp/access.log")

## The approach taken

The approach taken is to tail through the log file, and execute a generic
handler for each log line received.

I have also setup a very simple metrics package which allows storing metrics
whenever a log line is received. Only counters are currently supported.

Whenever a log line is received, the handler increments several metrics, for
the total number of requests, the section we're tracking, the HTTP response
code etc.

The metrics API then stores that data in batches of seconds, so we can fetch
only the data for a specific window of time.  It then handles taking a snapshot
of all the metrics for a specific time window, which is used to build the
dashboard and alerts.

## Known issues

* Being able to track only metrics as key/count means tracking something which may have a lot of different keys and low counter values isn't going to be efficient.
  * For example, monitoring IP addresses with one metric for each wouldn't be a good idea.

## Next steps

* Only counters are handled. We could want to track the response size, and generate histograms with them.
* Fetching logs locally isn't really something which can be used in a sufficiently large production environment where there should be more than one machine.
  * One approach could be to tail the logs to Kafka, and have a decoupled process which analyzes them and sends them to the monitoring system.
* The alerting system is pull-based. We could move it to a push-system.
  * This works fine right now, as we only display alerts in the CLI dashboard, and collecting the data is quite fast.
  * We couldn't send the alerts to somewhere else (an email for example) with this approach.
* We are not cleaning up old metrics data even when we know we're not going to use it anymore.
  * We could have a background goroutine which truncates all metrics from data older than the longest window we ever look at, to prevent the system from accruing too much memory.

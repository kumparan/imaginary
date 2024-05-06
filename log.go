package main

import (
	"fmt"
	"github.com/kumparan/go-utils"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const formatPattern = `{"remote_ip": "%s", "time": "%s", "method": "%s", "uri": "%s", "protocol": "%s", "status": "%d", "imaginary_bytes_out": %d, "imaginary_duration_in_ms": %d, "latency_human": "%s", "masked_imaginary_uri": "%s"}%s`

var maskedFields = []string{"s3", "text", "image", "font"}

// LogRecord implements an Apache-compatible HTTP logging
type LogRecord struct {
	http.ResponseWriter
	status                int
	responseBytes         int64
	ip                    string
	method, uri, protocol string
	time                  time.Time
	elapsedTime           time.Duration
}

// Log writes a log entry in the passed io.Writer stream
func (r *LogRecord) Log(out io.Writer) {
	go func() {
		timeFormat := r.time.Format(time.RFC3339Nano)

		splited := strings.Split(r.uri, "?")
		if len(splited) <= 0 {
			_, _ = fmt.Fprintf(out, formatPattern, r.ip, timeFormat, r.method, r.uri, r.protocol, r.status, r.responseBytes, r.elapsedTime.Milliseconds(), r.elapsedTime.String(), "", "\n")
			return
		}
		if len(splited) <= 1 {
			_, _ = fmt.Fprintf(out, formatPattern, r.ip, timeFormat, r.method, r.uri, r.protocol, r.status, r.responseBytes, r.elapsedTime.Milliseconds(), r.elapsedTime.String(), splited[0], "\n")
			return
		}

		maskedURI := splited[0]
		queryParam, err := url.ParseQuery(splited[1])
		if err != nil {
			log.WithField("queryParam", splited[1]).Error(err)
			_, _ = fmt.Fprintf(out, formatPattern, r.ip, timeFormat, r.method, r.uri, r.protocol, r.status, r.responseBytes, r.elapsedTime.Milliseconds(), r.elapsedTime.String(), maskedURI, "\n")
		}

		newQueryParam := url.Values{}
		for k, params := range queryParam {
			if utils.Contains(maskedFields, k) {
				newQueryParam.Add(k, "_")
				continue
			}
			for _, v := range params {
				newQueryParam.Add(k, v)
			}
		}
		maskedURI += "?" + newQueryParam.Encode()
		_, _ = fmt.Fprintf(out, formatPattern, r.ip, timeFormat, r.method, r.uri, r.protocol, r.status, r.responseBytes, r.elapsedTime.Milliseconds(), r.elapsedTime.String(), maskedURI, "\n")
	}()
}

// Write acts like a proxy passing the given bytes buffer to the ResponseWritter
// and additionally counting the passed amount of bytes for logging usage.
func (r *LogRecord) Write(p []byte) (int, error) {
	written, err := r.ResponseWriter.Write(p)
	r.responseBytes += int64(written)
	return written, err
}

// WriteHeader calls ResponseWriter.WriteHeader() and sets the status code
func (r *LogRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// LogHandler maps the HTTP handler with a custom io.Writer compatible stream
type LogHandler struct {
	handler  http.Handler
	io       io.Writer
	logLevel string
}

// NewLog creates a new logger
func NewLog(handler http.Handler, io io.Writer, logLevel string) http.Handler {
	return &LogHandler{handler, io, logLevel}
}

// Implements the required method as standard HTTP handler, serving the request.
func (h *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
		clientIP = clientIP[:colon]
	}

	record := &LogRecord{
		ResponseWriter: w,
		ip:             clientIP,
		time:           time.Time{},
		method:         r.Method,
		uri:            r.RequestURI,
		protocol:       r.Proto,
		status:         http.StatusOK,
		elapsedTime:    time.Duration(0),
	}

	startTime := time.Now()
	h.handler.ServeHTTP(record, r)
	finishTime := time.Now()

	record.time = finishTime.UTC()
	record.elapsedTime = finishTime.Sub(startTime)

	switch h.logLevel {
	case "error":
		if record.status >= http.StatusInternalServerError {
			record.Log(h.io)
		}
	case "warning":
		if record.status >= http.StatusBadRequest {
			record.Log(h.io)
		}
	case "info":
		record.Log(h.io)
	}
}

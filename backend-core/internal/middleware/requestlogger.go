package middleware

import (
	"bufio"
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/Wave-95/boards/backend-core/pkg/logger"
	"github.com/google/uuid"
)

const (
	headerNameRequestID     = "X-Request-ID"
	headerNameCorrelationID = "X-Correlation-ID"

	fieldDuration      = "duration"
	fieldBytes         = "bytes"
	fieldRequestID     = "requestID"
	fieldCorrelationID = "correlationID"
)

// LoggerRW is a wrapper around ResponseWriter meant to capture the status code and number of bytes written
type LoggerRW struct {
	http.ResponseWriter
	StatusCode   int
	BytesWritten int
}

// Write overrides the ResponseWriter Write method. It saves the amount of bytes written to the BytesWritten field which
// will be used in the request logger.
func (lrw *LoggerRW) Write(p []byte) (int, error) {
	bytesWritten, err := lrw.ResponseWriter.Write(p)
	lrw.BytesWritten = bytesWritten

	return bytesWritten, err
}

// WriteHeader overrides the ResponseWriter WriteHeader method. It saves the status code to the StatusCode field which
// will be used in the request logger.
func (lrw *LoggerRW) WriteHeader(statusCode int) {
	lrw.StatusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

// Hijack method is added as a method to LoggerRW to support WebSockets.
func (lrw *LoggerRW) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := lrw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}

	return h.Hijack()
}

// RequestLogger is a middleware that creates a context-aware logger for every request and makes it available to downstream
// handlers on the request context.
func RequestLogger(l logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Get request and correlation IDs and append fields to request logger
			// Then set logger to request context
			reqID, corrID := getOrCreateIDs(r)
			requestLogger := l.With(fieldRequestID, reqID, fieldCorrelationID, corrID)
			ctx := r.Context()
			ctx = context.WithValue(ctx, logger.LoggerKey, requestLogger)
			r = r.WithContext(ctx)

			// Wrap rw to make status code and bytes written available
			lrw := &LoggerRW{ResponseWriter: rw, StatusCode: http.StatusOK}
			next.ServeHTTP(lrw, r)

			queryString := ""
			if r.URL.RawQuery != "" {
				queryString = "?" + r.URL.RawQuery
			}
			// Log duration of request
			requestLogger.
				WithoutCaller().
				With(fieldDuration, time.Since(start).Milliseconds(), fieldBytes, lrw.BytesWritten).
				Infof("[%v] %s: %s%s", lrw.StatusCode, r.Method, r.URL.Path, queryString)
		})
	}
}

// Look for existing request ID and correlation ID from request header
// Return existing values or generate new uuids if not found
func getOrCreateIDs(r *http.Request) (reqID string, corrID string) {
	reqID = getRequestID(r)
	corrID = getCorrelationID(r)

	if reqID == "" {
		reqID = uuid.NewString()
	}

	if corrID == "" {
		corrID = uuid.NewString()
	}

	return reqID, corrID
}

// getRequestID grabs the request ID string off the X-Request-ID header
func getRequestID(r *http.Request) string {
	return r.Header.Get(headerNameRequestID)
}

// getCorrelationId grabs the correlation ID string off the X-Correlation-ID header
// The correlation id groups together multiple request ids
func getCorrelationID(r *http.Request) string {
	return r.Header.Get(headerNameCorrelationID)
}

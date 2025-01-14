package log

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	headerXForwardedFor    = "X-Forwarded-For"
	headerXRealIP          = "X-Real-IP"
	headerUberTraceID      = "Uber-Trace-Id"
	contentLengthThreshold = 1 << 16 // ~ 65 kB
)

var ErrBadReq = errors.New("error") // for satisfying platform's logger interface

func LoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrw := &wrappedResponseWriter{
			ResponseWriter: w,
			StatusCode:     200,
		}

		start := time.Now()

		contentLength, _ := strconv.Atoi(r.Header.Get("Content-Length"))

		var (
			reqBody       []byte
			bodyTruncated bool
		)
		if r.Body != nil {
			reqBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))

			if contentLength > contentLengthThreshold {
				reqBody = reqBody[:contentLengthThreshold]
				bodyTruncated = true
			}
		}

		next.ServeHTTP(wrw, r)

		l := FromContext(r.Context())

		l = l.With(
			"method", r.Method,
			"uri", r.RequestURI,
			"status", fmt.Sprintf("%d %s", wrw.StatusCode, wrw.Outcome),
			"userAgent", r.UserAgent(),
			"xForwardedFor", r.Header.Get(headerXForwardedFor),
			"xRealIP", r.Header.Get(headerXRealIP),
			"uberTraceID", r.Header.Get(headerUberTraceID),
			"remoteAddr", r.RemoteAddr,
			"duration", fmt.Sprintf("%.6fs", time.Since(start).Seconds()),
			"requestBody", string(reqBody),
			"requestBodyTruncated", bodyTruncated,
			"requestBodySize", contentLength,
		)

		if wrw.StatusCode >= 400 {
			l.Error("request processed")
		} else {
			l.Info("request processed")
		}
	})
}

func SetLoggerM(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(
			WithLogger(r.Context(), WithTraceID(r.Context(), global)),
		)

		next.ServeHTTP(w, r)
	})
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Outcome    string
}

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
	w.Outcome = http.StatusText(statusCode)
}

type response struct {
	http.ResponseWriter
	StatusCode int
	Outcome    string
}

func (w *response) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
	w.Outcome = http.StatusText(statusCode)
}

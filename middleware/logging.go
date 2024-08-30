package middleware

import (
    "bytes"
    // "fmt"
    "io/ioutil"
    "net/http"
    "os"
    // "time"

    "github.com/sirupsen/logrus"
)

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Setup logging
        logger := logrus.New()
        file, err := os.OpenFile("requests.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err == nil {
            logger.Out = file
        } else {
            logger.Info("Failed to log to file, using default stderr")
        }

        // Capture request body
        var bodyBytes []byte
        if r.Body != nil {
            bodyBytes, _ = ioutil.ReadAll(r.Body)
        }
        r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

        // Log request
        logger.WithFields(logrus.Fields{
            "method": r.Method,
            "url":    r.URL.String(),
            "body":   string(bodyBytes),
        }).Info("Request")

        // Create response recorder
        rw := &responseWriter{ResponseWriter: w, body: &bytes.Buffer{}}

        // Call the next handler
        next.ServeHTTP(rw, r)

        // Log response
        logger.WithFields(logrus.Fields{
            "status": rw.statusCode,
            "body":   rw.body.String(),
        }).Info("Response")
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
    body       *bytes.Buffer
}

func (rw *responseWriter) WriteHeader(statusCode int) {
    rw.statusCode = statusCode
    rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    rw.body.Write(b)
    return rw.ResponseWriter.Write(b)
}

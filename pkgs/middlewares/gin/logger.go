package middleware

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/dukk308/beetool.dev-go-starter/pkgs/base"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/constants"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/logger"
	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

var loggerExcludedPaths = []string{"/health", "/ping", "/alive"}

func isLoggerExcludedPath(path string) bool {
	for _, excluded := range loggerExcludedPaths {
		if path == excluded {
			return true
		}
	}
	return false
}

func Logger(isLogRequest bool, isLogResponse bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		excluded := isLoggerExcludedPath(path)

		start := time.Now()

		if isLogRequest && !excluded {
			log := logger.FromContext(c.Request.Context())
			logRequest(c, log)
		}

		blw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		c.Next()

		statusCode := c.Writer.Status()
		duration := time.Since(start)
		log := logger.FromContext(c.Request.Context())

		if !excluded && statusCode >= 400 {
			logError(c, log, statusCode, duration, blw.body.String())
		} else if !excluded && isLogResponse {
			logResponse(c, log, statusCode, duration, blw.body.String())
		}
	}
}

func logRequest(c *gin.Context, logger logger.Logger) {
	requestBody := ""
	contentType := c.ContentType()

	if c.Request.Body != nil {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err == nil {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			if len(bodyBytes) > 0 {
				if strings.Contains(contentType, "application/json") {
					requestBody = string(bodyBytes)
				} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
					requestBody = string(bodyBytes)
				} else if strings.Contains(contentType, "multipart/form-data") {
					requestBody = "[multipart/form-data]"
				} else if strings.Contains(contentType, "text/") {
					requestBody = string(bodyBytes)
				} else if strings.Contains(contentType, "application/xml") {
					requestBody = string(bodyBytes)
				} else {
					requestBody = fmt.Sprintf("[binary data: %d bytes]", len(bodyBytes))
				}
			}
		}
	}

	logMsg := fmt.Sprintf(
		"[REQUEST] %s %s | Proto: %s | Host: %s | RemoteAddr: %s | UserAgent: %s | ContentType: %s",
		c.Request.Method,
		c.Request.URL.RequestURI(),
		c.Request.Proto,
		c.Request.Host,
		c.Request.RemoteAddr,
		c.Request.UserAgent(),
		contentType,
	)

	if requestBody != "" {
		logMsg = fmt.Sprintf("%s | Body: %s", logMsg, requestBody)
	}

	logger.Info(logMsg)
}

func logResponse(c *gin.Context, log logger.Logger, statusCode int, duration time.Duration, responseBody string) {
	logMsg := fmt.Sprintf(
		"[RESPONSE] %s %s | Status: %d | Duration: %v",
		c.Request.Method,
		c.Request.URL.RequestURI(),
		statusCode,
		duration,
	)
	if responseBody != "" && len(responseBody) < 1000 {
		logMsg = fmt.Sprintf("%s | Body: %s", logMsg, responseBody)
	} else if len(responseBody) >= 1000 {
		logMsg = fmt.Sprintf("%s | Body: [%d bytes]", logMsg, len(responseBody))
	}
	log.Info(logMsg)
}

func logError(c *gin.Context, log logger.Logger, statusCode int, duration time.Duration, responseBody string) {
	logMsg := fmt.Sprintf(
		"[RESPONSE] %s %s | Status: %d | Duration: %v",
		c.Request.Method,
		c.Request.URL.RequestURI(),
		statusCode,
		duration,
	)
	if val, exists := c.Get(constants.ContextKeyError); exists {
		if err, ok := val.(error); ok && err != nil {
			errToLog := err
			if domainErr, ok := base.AsDomainError(err); ok && domainErr.Unwrap() != nil {
				if root := domainErr.RootCause(); root != nil {
					errToLog = root
				}
			}
			logMsg = fmt.Sprintf("%s | Error: %v", logMsg, errToLog)
		}
	}
	if responseBody != "" && len(responseBody) < 1000 {
		logMsg = fmt.Sprintf("%s | Body: %s", logMsg, responseBody)
	} else if len(responseBody) >= 1000 {
		logMsg = fmt.Sprintf("%s | Body: [%d bytes]", logMsg, len(responseBody))
	}
	log.Error(logMsg)
}

package gin_comp

import "github.com/dukk308/golang-clean-arch-starter/pkgs/logger"

type ginLoggerWriter struct {
	logger logger.Logger
}

func (w ginLoggerWriter) Write(p []byte) (n int, err error) {
	msg := string(p)
	if len(msg) > 0 && msg[len(msg)-1] == '\n' {
		msg = msg[:len(msg)-1]
	}
	if len(msg) > 0 {
		w.logger.Debug(msg)
	}
	return len(p), nil
}

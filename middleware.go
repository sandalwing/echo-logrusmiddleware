package logrusmiddleware

import (
	"io"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type Logger struct {
	*logrus.Logger
}

func (l Logger) Level() log.Lvl {
	switch l.Logger.Level {
	case logrus.DebugLevel:
		return log.DEBUG
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel:
		return log.ERROR
	case logrus.InfoLevel:
		return log.INFO
	default:
		l.Panic("Invalid level")
	}

	return log.OFF
}

func (l Logger) SetPrefix(s string) {
	// TODO
}

func (l Logger) Prefix() string {
	// TODO.  Is this even valid?  I'm not sure it can be translated since
	// logrus uses a Formatter interface.  Which seems to me to probably be
	// a better way to do it.
	return ""
}

func (l Logger) SetLevel(lvl log.Lvl) {
	switch lvl {
	case log.DEBUG:
		logrus.SetLevel(logrus.DebugLevel)
	case log.WARN:
		logrus.SetLevel(logrus.WarnLevel)
	case log.ERROR:
		logrus.SetLevel(logrus.ErrorLevel)
	case log.INFO:
		logrus.SetLevel(logrus.InfoLevel)
	default:
		l.Panic("Invalid level")
	}
}

func (l Logger) Output() io.Writer {
	return l.Out
}

func (l Logger) SetOutput(w io.Writer) {
	logrus.SetOutput(w)
}

func (l Logger) Printj(j log.JSON) {
	logrus.WithFields(logrus.Fields(j)).Print()
}

func (l Logger) Debugj(j log.JSON) {
	logrus.WithFields(logrus.Fields(j)).Debug()
}

func (l Logger) Infoj(j log.JSON) {
	logrus.WithFields(logrus.Fields(j)).Info()
}

func (l Logger) Warnj(j log.JSON) {
	logrus.WithFields(logrus.Fields(j)).Warn()
}

func (l Logger) Errorj(j log.JSON) {
	logrus.WithFields(logrus.Fields(j)).Error()
}

func (l Logger) Fatalj(j log.JSON) {
	logrus.WithFields(logrus.Fields(j)).Fatal()
}

func (l Logger) Panicj(j log.JSON) {
	logrus.WithFields(logrus.Fields(j)).Panic()
}

func logrusMiddlewareHandler(c echo.Context, next echo.HandlerFunc) error {
	req := c.Request()
	res := c.Response()
	start := time.Now()
	if err := next(c); err != nil {
		c.Error(err)
	}
	stop := time.Now()

	p := req.URL.Path
	if p == "" {
		p = "/"
	}

	bytesIn := req.Header.Get(echo.HeaderContentLength)
	if bytesIn == "" {
		bytesIn = "0"
	}

	logrus.WithFields(map[string]interface{}{
		"time_rfc3339":  time.Now().Format(time.RFC3339),
		"remote_ip":     c.RealIP(),
		"host":          req.Host,
		"uri":           req.RequestURI,
		"method":        req.Method,
		"path":          p,
		"referer":       req.Referer(),
		"user_agent":    req.UserAgent(),
		"status":        res.Status,
		"latency":       strconv.FormatInt(stop.Sub(start).Nanoseconds()/1000, 10),
		"latency_human": stop.Sub(start).String(),
		"bytes_in":      bytesIn,
		"bytes_out":     strconv.FormatInt(res.Size, 10),
	}).Info("Handled request")

	return nil
}

func logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return logrusMiddlewareHandler(c, next)
	}
}

func Hook() echo.MiddlewareFunc {
	return logger
}

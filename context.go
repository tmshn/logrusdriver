package logrusdriver

import (
	"io"

	"github.com/labstack/echo"
	gommonlog "github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

type ContextWithLogrus struct {
	echo.Context
	logger *logrusLogger
}

func (c *ContextWithLogrus) Logger() echo.Logger {
	return c.logger
}

func wrapContextWithLogrus(c echo.Context, entry *logrus.Entry) echo.Context {
	return &ContextWithLogrus{logger: &logrusLogger{entry}}
}

type logrusLogger struct {
	*logrus.Entry
}

func (l *logrusLogger) Output() io.Writer {
	return l.Logger.Out
}

func (l *logrusLogger) SetOutput(w io.Writer) {
	l.Logger.SetOutput(w)
}

func (l *logrusLogger) Prefix() string {
	return l.Data["prefix"].(string)
}

func (l *logrusLogger) SetPrefix(p string) {
	l.Data["prefix"] = p
}

func (l *logrusLogger) Level() gommonlog.Lvl {
	logrusLevel := l.Logger.GetLevel()
	switch {
	case logrusLevel <= logrus.PanicLevel:
		return gommonlog.OFF + 2 // gommonlog.panicLevel is not exported
	case logrusLevel <= logrus.FatalLevel:
		return gommonlog.OFF + 1 // gommonlog.fatalLevel is not exported
	case logrusLevel <= logrus.ErrorLevel:
		return gommonlog.ERROR
	case logrusLevel <= logrus.WarnLevel:
		return gommonlog.WARN
	case logrusLevel <= logrus.InfoLevel:
		return gommonlog.INFO
	default:
		return gommonlog.DEBUG
	}
}

func (l *logrusLogger) SetLevel(v gommonlog.Lvl) {
	switch {
	case v >= gommonlog.OFF+2:
		l.Logger.SetLevel(logrus.PanicLevel)
	case v >= gommonlog.OFF:
		l.Logger.SetLevel(logrus.FatalLevel)
	case v >= gommonlog.ERROR:
		l.Logger.SetLevel(logrus.ErrorLevel)
	case v >= gommonlog.WARN:
		l.Logger.SetLevel(logrus.WarnLevel)
	case v >= gommonlog.INFO:
		l.Logger.SetLevel(logrus.InfoLevel)
	default:
		l.Logger.SetLevel(logrus.DebugLevel)
	}
}

func (l *logrusLogger) SetHeader(string) {
	// We ignore this
}

func (l *logrusLogger) Printj(j gommonlog.JSON) {
	fields := logrus.Fields(j)
	msg := fields["message"]
	delete(fields, "message")
	l.WithFields(fields).Print(msg)
}

func (l *logrusLogger) Debugj(j gommonlog.JSON) {
	fields := logrus.Fields(j)
	msg := fields["message"]
	delete(fields, "message")
	l.WithFields(fields).Debug(msg)
}

func (l *logrusLogger) Infoj(j gommonlog.JSON) {
	fields := logrus.Fields(j)
	msg := fields["message"]
	delete(fields, "message")
	l.WithFields(fields).Info(msg)
}

func (l *logrusLogger) Warnj(j gommonlog.JSON) {
	fields := logrus.Fields(j)
	msg := fields["message"]
	delete(fields, "message")
	l.WithFields(fields).Warn(msg)
}

func (l *logrusLogger) Errorj(j gommonlog.JSON) {
	fields := logrus.Fields(j)
	msg := fields["message"]
	delete(fields, "message")
	l.WithFields(fields).Error(msg)
}

func (l *logrusLogger) Fatalj(j gommonlog.JSON) {
	fields := logrus.Fields(j)
	msg := fields["message"]
	delete(fields, "message")
	l.WithFields(fields).Fatal(msg)
}

func (l *logrusLogger) Panicj(j gommonlog.JSON) {
	fields := logrus.Fields(j)
	msg := fields["message"]
	delete(fields, "message")
	l.WithFields(fields).Panic(msg)
}

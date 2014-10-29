package log

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInfo
	LevelDebug

	colorBlack   = 30
	colorRed     = 31
	colorGreen   = 32
	colorYellow  = 33
	colorBlue    = 34
	colorMagenta = 35
	colorCyan    = 36
	colorWhite   = 97

	colorReset = "\033[0m"
)

var logLevelColors map[int]string

// Level represents the a severity level as defined for Syslog
type Level int

func (level Level) String() string {
	switch level {
	case LevelEmergency:
		return "emergency"
	case LevelAlert:
		return "alert"
	case LevelCritical:
		return "critical"
	case LevelError:
		return "error"
	case LevelWarning:
		return "warning"
	case LevelNotice:
		return "notice"
	case LevelInfo:
		return "info"
	case LevelDebug:
		return "debug"
	}
	return ""
}

// NowFunc is a type of function that returns a time. Useful for unit testing
type NowFunc func() time.Time

// Logger defines an interface for UNIX-like (Syslog) logging
type Logger interface {
	Emergency(string) error
	Alert(string) error
	Critical(string) error
	Error(string) error
	Warning(string) error
	Notice(string) error
	Info(string) error
	Debug(string) error
}

func init() {
	logLevelColors = getLevelColors()
}

// NilLogger is a nil implementation of the Logger interface
type NilLogger struct{}

func (l NilLogger) Emergency(m string) error { return nil }
func (l NilLogger) Alert(m string) error     { return nil }
func (l NilLogger) Critical(m string) error  { return nil }
func (l NilLogger) Error(m string) error     { return nil }
func (l NilLogger) Warning(m string) error   { return nil }
func (l NilLogger) Notice(m string) error    { return nil }
func (l NilLogger) Info(m string) error      { return nil }
func (l NilLogger) Debug(m string) error     { return nil }

// wLogger implements the Logger interface using a Writer to log to
type wLogger struct {
	writer   io.Writer
	level    int
	prefix   string
	pattern  string
	coloring bool
	nowFunc  NowFunc
}

// New returns a new wLogger, which uses a writer to write the messages
func New(w io.Writer) *wLogger {
	return &wLogger{
		writer:   w,
		level:    LevelDebug,
		pattern:  "{{ color }}{{ time }} {{ prefix }} [{{ level_literal }}] {{ message }}{{ color_reset }}\n",
		coloring: true,
		nowFunc:  time.Now,
	}
}

// SetLevel sets the threshold level for the logger.
// Only messages with a level lower or equal to the threshold level will be written.
// You can use the defined LevelXXX constants to set it.
// Ex: logger.SetLevel(log.LevelDebug)
func (l *wLogger) SetLevel(level int) {
	l.level = level
}

// SetPrefix sets the prefix for the log lines.
// This is helpful to filter log contents.
func (l *wLogger) SetPrefix(prefix string) {
	l.prefix = prefix
}

// SetPattern sets the log line patter for the logger.
// Defined tokens are:
// {{ time }} - the actual time of the logged event
// {{ level_literal }} - the literal representation of the severity level
// {{ level }} - the numeric severity level
// {{ message }} - the message beign logged
// {{ prefix }} - the prefix set to the logger (if any)
// {{ color }} - the terminal escape sequence for the color assigned to the log level
// {{ color_reset }} - the terminal escape sequence for resetting the coloring (foreground and background)
func (l *wLogger) SetPattern(pattern string) {
	l.pattern = pattern
}

// SetColoring
func (l *wLogger) SetColoring(b bool) {
	l.coloring = b
}

// SetNowFunc sets a custom function for getting the log event time
func (l *wLogger) SetNowFunc(nowFunc NowFunc) {
	l.nowFunc = nowFunc
}

func (l *wLogger) log(level int, message string) error {
	if level > l.level {
		return nil
	}

	var colorSeq, colorOff string
	if l.coloring {
		colorSeq = levelColor(level)
		colorOff = colorReset
	} else {
		colorSeq = ""
		colorOff = ""
	}

	line := l.pattern
	line = strings.Replace(line, "{{ time }}", l.now(), -1)
	line = strings.Replace(line, "{{ level }}", strconv.Itoa(level), -1)
	line = strings.Replace(line, "{{ level_literal }}", strings.ToUpper(Level(level).String()), -1)
	line = strings.Replace(line, "{{ prefix }}", l.prefix, -1)
	line = strings.Replace(line, "{{ message }}", message, -1)
	line = strings.Replace(line, "{{ color }}", colorSeq, -1)
	line = strings.Replace(line, "{{ color_reset }}", colorOff, -1)

	_, err := l.writer.Write([]byte(line))
	return err
}

func (l *wLogger) now() string {
	return l.nowFunc().Format(time.RFC3339)
}

func (l wLogger) Emergency(m string) error { return l.log(LevelEmergency, m) }
func (l wLogger) Alert(m string) error     { return l.log(LevelAlert, m) }
func (l wLogger) Critical(m string) error  { return l.log(LevelCritical, m) }
func (l wLogger) Error(m string) error     { return l.log(LevelError, m) }
func (l wLogger) Warning(m string) error   { return l.log(LevelWarning, m) }
func (l wLogger) Notice(m string) error    { return l.log(LevelNotice, m) }
func (l wLogger) Info(m string) error      { return l.log(LevelInfo, m) }
func (l wLogger) Debug(m string) error     { return l.log(LevelDebug, m) }

func getLevelColors() map[int]string {
	return map[int]string{
		LevelEmergency: colorEscape(colorMagenta, true),
		LevelAlert:     colorEscape(colorMagenta, false),
		LevelCritical:  colorEscape(colorRed, false),
		LevelError:     colorEscape(colorRed, false),
		LevelWarning:   colorEscape(colorYellow, true),
		LevelNotice:    colorEscape(colorYellow, false),
		LevelInfo:      colorEscape(colorGreen, false),
		LevelDebug:     colorEscape(colorCyan, false),
	}
}

func colorEscape(color int, bold bool) string {
	if bold {
		return fmt.Sprintf("\033[%d;1m", color)
	} else {
		return fmt.Sprintf("\033[%dm", color)
	}
}

func levelColor(level int) string {
	return logLevelColors[level]
}

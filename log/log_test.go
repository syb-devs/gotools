package log_test

import (
	"bytes"
	"testing"
	"time"

	"bitbucket.org/syb-devs/gotools/log"
)

func now() time.Time {
	return time.Unix(0, 0)
}

var logLevelsTests = []struct {
	prefix    string
	coloring  bool
	pattern   string
	level     int
	threshold int
	message   string
	expected  string
}{
	{
		prefix:    "--TestOne--",
		coloring:  true,
		level:     log.LevelDebug,
		threshold: log.LevelDebug,
		message:   "Testing debug...",
		expected: "[36m1970-01-01T01:00:00+01:00 --TestOne-- [DEBUG] Testing debug...[0m\n",
	},
	{
		prefix:    "",
		coloring:  false,
		level:     log.LevelEmergency,
		threshold: log.LevelEmergency,
		message:   "Testing emergency...",
		expected:  "1970-01-01T01:00:00+01:00  [EMERGENCY] Testing emergency...\n",
	},
	{
		prefix:    "[TABLE_FLIP]",
		coloring:  false,
		level:     log.LevelWarning,
		threshold: log.LevelDebug,
		message:   "(╯°□°）╯︵ ┻━┻",
		expected:  "1970-01-01T01:00:00+01:00 [TABLE_FLIP] [WARNING] (╯°□°）╯︵ ┻━┻\n",
	},
	{
		prefix:    "[CUSTOM_PATTERN]",
		coloring:  true,
		pattern:   "{{ time }} {{ level_literal }}{{ color }}{{ message }}{{ color_reset }}\n",
		level:     log.LevelAlert,
		threshold: log.LevelWarning,
		message:   "Bring back life form. Priority One. All other priorities rescinded.",
		expected: "1970-01-01T01:00:00+01:00 ALERT[35mBring back life form. Priority One. All other priorities rescinded.[0m\n",
	},
	{
		level:     log.LevelInfo,
		threshold: log.LevelWarning,
		message:   "This should not be logged",
		expected:  "",
	},
}

func TestLogLevels(t *testing.T) {
	for _, test := range logLevelsTests {
		w := &bytes.Buffer{}
		l := log.New(w)
		l.SetNowFunc(now)
		l.SetPrefix(test.prefix)
		l.SetColoring(test.coloring)
		l.SetLevel(test.threshold)

		if test.pattern != "" {
			l.SetPattern(test.pattern)
		}

		switch test.level {
		case log.LevelEmergency:
			l.Emergency(test.message)
		case log.LevelAlert:
			l.Alert(test.message)
		case log.LevelCritical:
			l.Critical(test.message)
		case log.LevelError:
			l.Error(test.message)
		case log.LevelWarning:
			l.Warning(test.message)
		case log.LevelNotice:
			l.Notice(test.message)
		case log.LevelInfo:
			l.Info(test.message)
		case log.LevelDebug:
			l.Debug(test.message)
		}
		if read, _ := w.ReadString(0); read != test.expected {
			t.Errorf("expecting \n%v, got \n%v", test.expected, read)
		}
	}
}

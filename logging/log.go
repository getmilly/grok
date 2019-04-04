package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

//Logger ...
type Logger struct {
	applicationName string
}

//LogEntry ...
type LogEntry struct {
	entry *logrus.Entry
}

var (
	logger = NewLogger()
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.SetOutput(os.Stdout)

	logrus.SetLevel(logrus.DebugLevel)
}

//NewLogger creates a new logger.
func NewLogger() *Logger {
	return &Logger{}
}

//LogWithApplication configure logger to use specified app name.
func LogWithApplication(appName string) {
	logger.applicationName = appName
}

//LogWarn logs an event
func LogWarn(message string, args ...interface{}) {
	LogWith(nil).Warn(message, args...)
}

//LogInfo logs an event
func LogInfo(message string, args ...interface{}) {
	LogWith(nil).Info(message, args...)
}

//LogError logs an error
func LogError(message string, args ...interface{}) {
	LogWith(nil).Error(message, args...)
}

//Warn logs an event
func (e *LogEntry) Warn(message string, args ...interface{}) {
	e.entry.Warnf(message, args...)
}

//Info logs an event
func (e *LogEntry) Info(message string, args ...interface{}) {
	e.entry.Infof(message, args...)
}

//Error logs an error
func (e *LogEntry) Error(message string, args ...interface{}) {
	e.entry.Errorf(message, args...)
}

//LogWith ...
func LogWith(data interface{}) *LogEntry {
	fields := buildData(data)
	entry := logrus.WithFields(fields)
	return &LogEntry{
		entry: entry,
	}
}

func buildData(data interface{}) logrus.Fields {
	fields := make(logrus.Fields)

	fields["data"] = data
	fields["hostname"], _ = os.Hostname()
	fields["application_name"] = logger.applicationName

	return fields
}

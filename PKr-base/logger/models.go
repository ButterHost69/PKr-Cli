package logger

import "os"

type LogLevel int
type acceptedLogLevel map[LogLevel]bool

const (
	InfoLevel LogLevel = iota
	DebugLevel
	CriticalLevel
	NoneLevel
)

const (
	WORKSPACE_PKR_DIR  = ".PKr"
	LOGS_PKR_FILE_PATH = WORKSPACE_PKR_DIR + "\\logs.txt"
	USER_LOG_FILE_PATH = "\\logs.txt"
)

var logLevelName = map[LogLevel]string{
	InfoLevel:     "INFO",
	DebugLevel:    "DEBUG",
	CriticalLevel: "CRITICAL",
}

func (ll LogLevel) String() string {
	return logLevelName[ll]
}

// [ ] Find a Way to Avoid Frequent Read and Writes... IDK How ? Maybe its not possible ?

func logToFile(filepath, log string) error {
	// Opens or Creates the Log File
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	defer file.Close()
	if _, err := file.Write([]byte(log)); err != nil {
		return err
	}

	return nil
}
package logger

import (
	"fmt"
)

type UserLogger struct {
	acceptedLevel      LogLevel
	userConfigFilePath string
	printToTerminal    bool
}

func InitUserLogger(userLogFilePath string) *UserLogger {
	ul := UserLogger{
		acceptedLevel:      NoneLevel,
		userConfigFilePath: userLogFilePath,
		printToTerminal:    false,
	}
	return &ul
}

func (ul *UserLogger) SetPrintToTerminal(printToTerminal bool) {
	ul.printToTerminal = printToTerminal
}

func (ul *UserLogger) SetLogLevel(logLevel LogLevel) {
	ul.acceptedLevel = logLevel
}

func (ul *UserLogger) Info(workspace_name string, log any) {
	slog := fmt.Sprintf("[Info] %v\n", log)
	if ul.printToTerminal {
		fmt.Printf("[Info] %v\n", log)
	}

	if ul.acceptedLevel >= InfoLevel {
		if err := logToFile(ul.userConfigFilePath+USER_LOG_FILE_PATH, slog); err != nil {
			fmt.Println("Could Not Log To File.\nError: ", err)
		}
	}
}

func (ul *UserLogger) Debug(workspace_name string, log any) {
	slog := fmt.Sprintf("[Info] %v\n", log)
	if ul.printToTerminal {
		fmt.Printf("[Info] %v\n", log)
	}

	if ul.acceptedLevel >= DebugLevel {
		if err := logToFile(ul.userConfigFilePath+USER_LOG_FILE_PATH, slog); err != nil {
			fmt.Println("Could Not Log To File.\nError: ", err)
		}
	}
}

func (ul *UserLogger) Critical(workspace_name string, log any) {
	slog := fmt.Sprintf("[Info] %v\n", log)
	if ul.printToTerminal {
		fmt.Printf("[Info] %v\n", log)
	}

	if ul.acceptedLevel >= CriticalLevel {
		if err := logToFile(ul.userConfigFilePath+USER_LOG_FILE_PATH, slog); err != nil {
			fmt.Println("Could Not Log To File.\nError: ", err)
		}
	}
}

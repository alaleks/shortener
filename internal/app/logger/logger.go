package logger_test

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	LogsDirpath = "applogs"
	Permissions = 0o755
)

type LogDir struct {
	LogDirectory string
}

func New() *LogDir {
	err := os.Mkdir(LogsDirpath, Permissions)
	if err != nil {
		return nil
	}

	return &LogDir{
		LogDirectory: LogsDirpath,
	}
}

func SetLogFile() *os.File {
	year, month, day := time.Now().Date()
	fileName := fmt.Sprintf("%v-%v-%v.log", day, month.String(), year)
	filePath, _ := os.OpenFile(LogsDirpath+"/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, Permissions)

	return filePath
}

func (l *LogDir) Info() *log.Logger {
	getFilePath := SetLogFile()

	return log.New(getFilePath, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (l *LogDir) Warning() *log.Logger {
	getFilePath := SetLogFile()

	return log.New(getFilePath, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (l *LogDir) Error() *log.Logger {
	getFilePath := SetLogFile()

	return log.New(getFilePath, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (l *LogDir) Fatal() *log.Logger {
	getFilePath := SetLogFile()

	return log.New(getFilePath, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (l *LogDir) FatalWithErr(msg any) *log.Logger {
	getFilePath := SetLogFile()

	return log.New(getFilePath, fmt.Sprintf("FATAL: %v ", msg), log.Ldate|log.Ltime|log.Lshortfile)
}

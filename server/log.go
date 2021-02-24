package server

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

// log level definition
const (
	Fatal = "FATAL"
	Error = "ERROR"
	Warn = "WARN"
	Info = "INFO"
	Debug = "DEBUG"
)

var logger *log.Logger

func LogInit() {
	logOutputFile := fmt.Sprintf("/Users/shakexin/workspace/stock-platform/server.log")
	f, err :=os.OpenFile(logOutputFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Cannot create log file, caught error: %s", err.Error())
		return
	}
	logger = log.New(f, "", log.Ldate|log.Ltime)
	logger.SetOutput(&lumberjack.Logger{
		Filename:   logOutputFile,
		MaxSize:    128,
		MaxAge:     28,
		MaxBackups: 5,
		Compress:   true,
	})
}

func Log(lvl string, msg string, args ...interface{}) {
	msgBody := fmt.Sprintf(msg, args...)
	logMsg := fmt.Sprintf("[%s] [%s]", lvl, msgBody)
	if logger == nil {
		log.Println(logMsg)
		return
	}
	if err := logger.Output(2, logMsg); err != nil {
		log.Printf("Unexpected error happened when trying to write to logger: %s", err.Error())
		log.Println(logMsg)
	}
}

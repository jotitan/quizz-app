package logger

import (
	"log"
	"os"
	"runtime"
	"strings"
	"fmt"
	"sync"
)

/* Manage log */
/* If file is defined, write inside. In any case, write in console */

type ILogger interface {
	Print(message ...interface{})
}

// Default implementation with reel logger
type FileLogger struct {
	logger * log.Logger
}

func NewFileLogger(filename string)(FileLogger,FileLogger) {
	if file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm) ; err == nil {
		return FileLogger{log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lmicroseconds)},FileLogger{log.New(file, "ERROR ", log.Ldate|log.Ltime)}
	}
	return FileLogger{log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.Lmicroseconds)},FileLogger{log.New(os.Stdout, "ERROR ", log.Ldate|log.Ltime)}
}

func (dl FileLogger)Print(message ...interface{}){
	dl.logger.Println(message...)
}

// Logger manage log production
type ComplexLogger struct{
	info ILogger
	error ILogger
	writeConsole bool
}

func InitComplexLogger(infoLogger ILogger, errorLogger ILogger, console bool)*ComplexLogger{
	complexLogger = &ComplexLogger{infoLogger,errorLogger,console}
	return complexLogger
}

// print messages in specific logger
func (cl ComplexLogger) print(levelLogger ILogger,message ...interface{}) {
	data := append([]interface{}{ getInfo()}, message...)
	levelLogger.Print(data...)
	if cl.writeConsole {
		log.Println(data...)
	}
}

func (cl ComplexLogger) Info(message ...interface{}) {
	cl.print(cl.info,message...)
}

func (cl ComplexLogger) Error(message ...interface{}) {
	cl.print(cl.error,message...)
}

// Fatal write fatal info into log
func (cl ComplexLogger) Fatal(message ...interface{}) {
	cl.print(cl.error,message...)
	os.Exit(1)
}

// Logger manage log production
type Logger struct{
	info * log.Logger
	error * log.Logger
	writeConsole bool
}

func getInfo() string {
	_, file, line, _ := runtime.Caller(3)
	return fmt.Sprintf("%s:%d:", file[strings.LastIndex(file, "/")+1:], line)
}

// Info write info into log
func (l Logger) Info(message ...interface{}) {
	l.print(l.info, message...)
}

// Fatal write fatal info into log
func (l Logger) Fatal(message ...interface{}) {
	l.print(l.error, message...)
	os.Exit(1)
}

// Erro write error message into log
func (l Logger) Error(message ...interface{}) {
	l.print(l.error, message...)
}

// pring write message into logger. If console is enabled, write into too
func (l Logger) print(loggerElement * log.Logger, message ...interface{}) {
	data := append([]interface{}{ getInfo()}, message...)
	loggerElement.Println(data...)
	if l.writeConsole {
		log.Println(data...)
	}
}

// InitLogger init the logger which will write messages into filename file and / or console
func InitLogger(filename string, console bool) *Logger {
	out := os.Stdout
	errOut := os.Stdout
	logger = &Logger{writeConsole:false}
	if filename != "" {
		if file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm) ; err == nil {
			out = file
			errOut = file
			logger.writeConsole = console
		}
	}
	logger.info = log.New(out, "INFO ", log.Ldate|log.Ltime|log.Lmicroseconds)
	logger.error = log.New(errOut, "ERROR ", log.Ldate|log.Ltime)
	return logger
}

// singleton of logger
var logger *Logger
var complexLogger *ComplexLogger
var lock = sync.Mutex{}

// GetLogger return the logger or create it if not exist
func GetLogger() *Logger {
	if logger == nil {
		lock.Lock()
		if logger == nil {
			InitLogger("", false)
		}
		lock.Unlock()
	}
	return logger
}

func GetLogger2() *ComplexLogger {
	if complexLogger == nil {
		lock.Lock()
		if complexLogger == nil {
			l1,l2 := NewFileLogger("")
			complexLogger = InitComplexLogger(l1,l2,false)
		}
		lock.Unlock()
	}
	return complexLogger
}
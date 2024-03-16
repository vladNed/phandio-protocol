package logging

import (
	"log"
	"os"
)

type LoggerOptions struct {
	Name  string
	Flags int
}

type ServerLogger struct {
	logger *log.Logger
}

func NewServerLogger(out *os.File, prefix string, flag int) *ServerLogger {
	return &ServerLogger{
		logger: log.New(out, prefix, flag),
	}
}

func (l *ServerLogger) Info(v ...interface{}) {
	l.logger.SetPrefix("[INFO] ")
	l.logger.Println(v...)
}

func (l *ServerLogger) Error(v ...interface{}) {
	l.logger.SetPrefix("[ERROR] ")
	l.logger.Println(v...)
}

func (l *ServerLogger) Fatal(v ...interface{}) {
	l.logger.SetPrefix("[FATAL] ")
	l.logger.Fatalln(v...)
}

func GetLogger(opts *LoggerOptions) *ServerLogger {
	if opts == nil {
		return NewServerLogger(os.Stdout, "", log.LstdFlags)
	}

	return NewServerLogger(os.Stdout, opts.Name, opts.Flags)
}

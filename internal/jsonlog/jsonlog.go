package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

// iota will succefily assign integer 
const(
	LevelInfo Level = iota // has the value 0
	LevelError   // has value 1,
	LevelFatal	// has value 2,
	LevelOff	 // has value 3,
)

func (l Level) String() string{
	switch l{
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// defining custom logger type
type Logger struct {
	out io.Writer
	minLevel Level
	mu sync.Mutex
}

func New(out io.Writer, minLevel Level) *Logger{
	return &Logger{
		out: out,
		minLevel: minLevel,
	}
}

// helper method to print log
func (l *Logger) PrintInfo(message string, properties map[string]string){
	l.print(LevelInfo, message, properties)
}

func (l *Logger) PrintError(err error, properties map[string]string){
	l.print(LevelError, err.Error(), properties)
}

func (l *Logger) PrintFatal(err error, properties map[string]string){
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

func (l *Logger) print(level Level, message string, properties map[string]string) (int, error){

	if level < l.minLevel{
		return 0, nil
	}

	// ananymous struct holding the data for the log 
	aux := struct {
		Level string `json:"level"`
		Time string `json:"time"`
		Message string `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace string `json:"trace,omitempty"`
	}{
		Level: level.String(),
		Time: time.Now().UTC().Format(time.RFC3339),
		Message: message,
		Properties: properties,
	}

	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	// for holding json formate bytes
	var line []byte

	// here we are changing our struct in json 
	line, err := json.Marshal(aux)
	if err != nil{
		line = []byte(LevelError.String() + ": unabel to marshal log message: "+ err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))

}

func (l *Logger) Write(message []byte) (n int, err error){
	return l.print(LevelError, string(message), nil)
}

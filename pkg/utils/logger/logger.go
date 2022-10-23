package logger

import (
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	LtF *logrus.Logger = logrus.New()
	Std *logrus.Logger = logrus.New()

	LogrusFile   *os.File
	LogFiberFile *os.File
	QueryLogFile *os.File
)

func logrusInit() {
	LtF.Formatter = new(logrus.JSONFormatter)
	LtF.Formatter.(*logrus.JSONFormatter).TimestampFormat = time.RFC822 // RFC3339
	// LtF.Formatter.(*logrus.JSONFormatter).PrettyPrint = true
	LtF.Formatter.(*logrus.JSONFormatter).FieldMap = logrus.FieldMap{
		logrus.FieldKeyTime:  "timestamp",
		logrus.FieldKeyLevel: "level",
		logrus.FieldKeyMsg:   "message",
		// logrus.FieldKeyFile:  "file",
		// logrus.FieldKeyLogrusError: "logrus_error",
		// logrus.FieldKeyFunc:        "func",
	}
	l, err := strconv.Atoi(os.Getenv("LOG_DETAIL_LEVEL"))
	if err != nil {
		logrus.Fatalf("ERROR DETAIL LOG LEVEL", err)
	}
	LtF.Level = logrus.Level(l) // Log level @ logrus.Logger.Level
	if LtF.Level > 5 {
		LtF.SetReportCaller(true)
	}
	file, err := os.OpenFile("./logs/detail.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		LtF.Out = file
	} else {
		logrus.Fatalf("error opening file: %v", err)
	}
	LogrusFile = file
}

func logrusStdoutInit() {
	Std.Formatter = new(logrus.TextFormatter)
	Std.Formatter.(*logrus.TextFormatter).FullTimestamp = true
	Std.Formatter.(*logrus.TextFormatter).TimestampFormat = time.RFC822
	Std.Formatter.(*logrus.TextFormatter).FieldMap = logrus.FieldMap{
		logrus.FieldKeyTime:  "timestamp",
		logrus.FieldKeyLevel: "level",
		logrus.FieldKeyMsg:   "message",
	}
	l, err := strconv.Atoi(os.Getenv("LOG_STDOUT_LEVEL"))
	if err != nil {
		logrus.Fatalf("ERROR STDOUT LOG LEVEL", err)
	}
	Std.Level = logrus.Level(l) // Log level @ logrus.Logger.Level
	Std.Out = os.Stdout
}

func logFiberInit() {
	file, err := os.OpenFile("./logs/access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("error opening file: %v", err)
	}
	LogFiberFile = file
}

func logQueryInit() {
	file, err := os.OpenFile("./logs/query.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("error opening file: %v", err)
	}
	QueryLogFile = file
}

func FileClose() {
	if err := LogrusFile.Close(); err != nil {
		logrus.Printf("logrus file close error \n")
	}
	if err := LogFiberFile.Close(); err != nil {
		logrus.Printf("fiber log file close error \n")
	}
	if err := QueryLogFile.Close(); err != nil {
		logrus.Printf("query log file close error \n")
	}
}

func dirInit() {
	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		os.Mkdir("./logs", 0777)
	} else {
		os.Chmod("./logs", 0777)
	}
}

func InitLogger() {
	dirInit()
	logrusInit()
	logrusStdoutInit()
	logFiberInit()
	logQueryInit()
}

/***
* @ logrus.Logger.Level
**	const (
		// PanicLevel level, highest level of severity. Logs and then calls panic with the
		// message passed to Debug, Info, ...
		PanicLevel Level = iota
		// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
		// logging level is set to Panic.
		FatalLevel
		// ErrorLevel level. Logs. Used for errors that should definitely be noted.
		// Commonly used for hooks to send errors to an error tracking service.
		ErrorLevel
		// WarnLevel level. Non-critical entries that deserve eyes.
		WarnLevel
		// InfoLevel level. General operational entries about what's going on inside the
		// application.
		InfoLevel
		// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
		DebugLevel
		// TraceLevel level. Designates finer-grained informational events than the Debug.
		TraceLevel
	)
*/

/***
* @ logrus.Output
**	The following methods have been overwritten(?) to print logs in two formats.
	- func Print(args ...interface{}) {}
	- func Printf(format string, args ...interface{}) {}
	- func Println(args ...interface{}) {}
	- func Panic(args ...interface{}) {}
	- func Panicf(format string, args ...interface{}) {}
	- func Panicln(args ...interface{}) {}
	- func Fatal(args ...interface{}) {}
	- func Fatalf(format string, args ...interface{}) {}
	- func Fatalln(args ...interface{}) {}
	- func Error(args ...interface{}) {}
	- func Errorf(format string, args ...interface{}) {}
	- func Errorln(args ...interface{}) {}
	- func Warn(args ...interface{}) {}
	- func Warnf(format string, args ...interface{}) {}
	- func Warnln(args ...interface{}) {}
	- func Info(args ...interface{}) {}
	- func Infof(format string, args ...interface{}) {}
	- func Infoln(args ...interface{}) {}
	- func Debug(args ...interface{}) {}
	- func Debugf(format string, args ...interface{}) {}
	- func Debugln(args ...interface{}) {}
	- func Trace(args ...interface{}) {}
	- func Tracef(format string, args ...interface{}) {}
	- func Traceln(args ...interface{}) {}
*/

// Standard log print out with user specific info
//
type User struct {
	Information string `json:"information"`
	Port        uint16 `json:"port"`
	Status      bool   `json:"status"`
}

func NewUser() *User {
	return &User{
		Information: "default info",
		Status:      true,
	}
}

// Logrus logging Standard Output Print (Info level)
func (u *User) Print(args ...interface{}) {
	args = append(append([]interface{}{}, u.Information), args...)
	// LtF.Print(args...)
	Std.Print(args...)
}
func (u *User) Printf(format string, args ...interface{}) {
	// LtF.Printf(format, args...)
	Std.Printf(u.Information+" "+format, args...)
}
func (u *User) Println(args ...interface{}) {
	args = append(append([]interface{}{}, u.Information), args...)
	// LtF.Println(args...)
	Std.Println(args...)
}

// Logrus logging Panic
func (u *User) Panic(args ...interface{}) {
	args = append(append([]interface{}{}, u.Information), args...)
	LtF.Panic(args...)
	Std.Panic(args...)
}
func (u *User) Panicf(format string, args ...interface{}) {
	LtF.Panicf(u.Information+" "+format, args...)
	Std.Panicf(u.Information+" "+format, args...)
}
func (u *User) Panicln(args ...interface{}) {
	args = append(append([]interface{}{}, u.Information), args...)
	LtF.Panicln(args...)
	Std.Panicln(args...)
}

// Logrus logging Fatal
func (u *User) Fatal(args ...interface{}) {
	args = append(append([]interface{}{}, u.Information), args...)
	LtF.Fatal(args...)
	Std.Fatal(args...)
}
func (u *User) Fatalf(format string, args ...interface{}) {
	LtF.Fatalf(u.Information+" "+format, args...)
	Std.Fatalf(u.Information+" "+format, args...)
}
func (u *User) Fatalln(args ...interface{}) {
	args = append(append([]interface{}{}, u.Information), args...)
	LtF.Fatalln(args...)
	Std.Fatalln(args...)
}

// Logrus logging Error
func (u *User) Error(args ...interface{}) {
	args = append(append([]interface{}{}, u.Information), args...)
	LtF.Error(args...)
	Std.Error(args...)
}
func (u *User) Errorf(format string, args ...interface{}) {
	LtF.Errorf(u.Information+" "+format, args...)
	Std.Errorf(u.Information+" "+format, args...)
}
func (u *User) Errorln(args ...interface{}) {
	args = append(append([]interface{}{}, u.Information), args...)
	LtF.Errorln(args...)
	Std.Errorln(args...)
}

// Logrus logging Warn
func (u *User) Warn(args ...interface{}) {
	args = append(append([]interface{}{}, u.Information), args...)
	LtF.Warn(args...)
	Std.Warn(args...)
}
func (u *User) Warnf(format string, args ...interface{}) {
	LtF.Warnf(u.Information+" "+format, args...)
	Std.Warnf(u.Information+" "+format, args...)
}
func (u *User) Warnln(args ...interface{}) {
	args = append(append([]interface{}{}, u.Information), args...)
	LtF.Warnln(args...)
	Std.Warnln(args...)
}

// Logrus logging Info
func (u *User) Info(args ...interface{}) {
	args = append(append([]interface{}{}, u.Information), args...)
	LtF.Info(args...)
	Std.Info(args...)
}
func (u *User) Infof(format string, args ...interface{}) {
	LtF.Infof(u.Information+" "+format, args...)
	Std.Infof(u.Information+" "+format, args...)
}
func (u *User) Infoln(args ...interface{}) {
	args = append(append([]interface{}{}, u.Information), args...)
	LtF.Infoln(args...)
	Std.Infoln(args...)
}

// Logrus logging Debug
func (u *User) Debug(args ...interface{}) {
	args = append(append([]interface{}{}, u.Information), args...)
	LtF.Debug(args...)
	Std.Debug(args...)
}
func (u *User) Debugf(format string, args ...interface{}) {
	LtF.Debugf(u.Information+" "+format, args...)
	Std.Debugf(u.Information+" "+format, args...)
}
func (u *User) Debugln(args ...interface{}) {
	args = append(append([]interface{}{}, u.Information), args...)
	LtF.Debugln(args...)
	Std.Debugln(args...)
}

// Logrus logging Trace
func (u *User) Trace(args ...interface{}) {
	args = append(append([]interface{}{}, u.Information), args...)
	LtF.Trace(args...)
	Std.Trace(args...)
}
func (u *User) Tracef(format string, args ...interface{}) {
	LtF.Tracef(u.Information+" "+format, args...)
	Std.Tracef(u.Information+" "+format, args...)
}
func (u *User) Traceln(args ...interface{}) {
	args = append(append([]interface{}{}, u.Information), args...)
	LtF.Traceln(args...)
	Std.Traceln(args...)
}

// Standard log print out
//

// Logrus logging Standard Output Print (Info level)
func Print(args ...interface{}) {
	// LtF.Print(args...)
	Std.Print(args...)
}
func Printf(format string, args ...interface{}) {
	// LtF.Printf(format, args...)
	Std.Printf(format, args...)
}
func Println(args ...interface{}) {
	// LtF.Println(args...)
	Std.Println(args...)
}

// Logrus logging Panic
func Panic(args ...interface{}) {
	LtF.Panic(args...)
	Std.Panic(args...)
}
func Panicf(format string, args ...interface{}) {
	LtF.Panicf(format, args...)
	Std.Panicf(format, args...)
}
func Panicln(args ...interface{}) {
	LtF.Panicln(args...)
	Std.Panicln(args...)
}

// Logrus logging Fatal
func Fatal(args ...interface{}) {
	LtF.Fatal(args...)
	Std.Fatal(args...)
}
func Fatalf(format string, args ...interface{}) {
	LtF.Fatalf(format, args...)
	Std.Fatalf(format, args...)
}
func Fatalln(args ...interface{}) {
	LtF.Fatalln(args...)
	Std.Fatalln(args...)
}

// Logrus logging Error
func Error(args ...interface{}) {
	LtF.Error(args...)
	Std.Error(args...)
}
func Errorf(format string, args ...interface{}) {
	LtF.Errorf(format, args...)
	Std.Errorf(format, args...)
}
func Errorln(args ...interface{}) {
	LtF.Errorln(args...)
	Std.Errorln(args...)
}

// Logrus logging Warn
func Warn(args ...interface{}) {
	LtF.Warn(args...)
	Std.Warn(args...)
}
func Warnf(format string, args ...interface{}) {
	LtF.Warnf(format, args...)
	Std.Warnf(format, args...)
}
func Warnln(args ...interface{}) {
	LtF.Warnln(args...)
	Std.Warnln(args...)
}

// Logrus logging Info
func Info(args ...interface{}) {
	LtF.Info(args...)
	Std.Info(args...)
}
func Infof(format string, args ...interface{}) {
	LtF.Infof(format, args...)
	Std.Infof(format, args...)
}
func Infoln(args ...interface{}) {
	LtF.Infoln(args...)
	Std.Infoln(args...)
}

// Logrus logging Debug
func Debug(args ...interface{}) {
	LtF.Debug(args...)
	Std.Debug(args...)
}
func Debugf(format string, args ...interface{}) {
	LtF.Debugf(format, args...)
	Std.Debugf(format, args...)
}
func Debugln(args ...interface{}) {
	LtF.Debugln(args...)
	Std.Debugln(args...)
}

// Logrus logging Trace
func Trace(args ...interface{}) {
	LtF.Trace(args...)
	Std.Trace(args...)
}
func Tracef(format string, args ...interface{}) {
	LtF.Tracef(format, args...)
	Std.Tracef(format, args...)
}
func Traceln(args ...interface{}) {
	LtF.Traceln(args...)
	Std.Traceln(args...)
}

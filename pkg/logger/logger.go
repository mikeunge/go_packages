package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Logger - basic logger struct
type Logger struct {
	settings *Settings
	*logrus.Logger
}

// Settings - the settings that define the logger
type Settings struct {
	level  string
	format string
	output string
	path   string
}

var (
	logger   *Logger
	settings *Settings
	once     sync.Once
)

// GetInstance - get the currently running instance or create a new one
func GetInstance() (*Logger, error) {
	var err error

	once.Do(func() {
		settings, err := initLogger()
		if err != nil {
			return
		}
		logger, err = createLogger(settings)
	})
	if err != nil {
		return logger, err
	}
	return logger, nil
}

func loadEnvVar(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return val, fmt.Errorf("key [%s] not found", key)
	}
	return val, nil
}

func initLogger() (*Settings, error) {
	err := godotenv.Load()
	if err != nil {
		return settings, err
	}

	// get all the environment keys
	logLevel, err := loadEnvVar("LOG_LEVEL")
	if err != nil {
		fmt.Printf("logLevel is not defined, using debug in this run")
		logLevel = "DEBUG"
	}
	logFormat, err := loadEnvVar("LOG_FORMAT")
	if err != nil {
		fmt.Println("no log format provided, going with plaintext")
		logFormat = "PLAIN"
	}
	logOutput, err := loadEnvVar("LOG_OUTPUT")
	if err != nil {
		fmt.Println("no output defined, printing to terminal")
		logOutput = "TERM"
	}
	logPath, err := loadEnvVar("LOG_PATH")
	if err != nil {
		if logOutput != "TERM" {
			return settings, fmt.Errorf("could not create logger, missing logpath")
		}
		logPath = ""
	}
	return &Settings{logLevel, logFormat, logOutput, logPath}, nil
}

func createLogger(settings *Settings) (*Logger, error) {
	log := logrus.New() // create new logrus instance

	// after gathering the environment info, apply it
	if strings.ToUpper(settings.output) == "FILE" {
		file, err := os.OpenFile(settings.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return logger, err
		}
		log.Out = file
	} else {
		log.Out = os.Stdout
	}
	if strings.ToUpper(settings.format) == "JSON" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{})
	}
	switch strings.ToUpper(settings.level) {
	case "TRACE":
		log.SetLevel(logrus.TraceLevel)
	case "DEBUG":
		log.SetLevel(logrus.DebugLevel)
	case "INFO":
		log.SetLevel(logrus.InfoLevel)
	case "WARN", "WARNING":
		log.SetLevel(logrus.WarnLevel)
	case "ERROR":
		log.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		log.SetLevel(logrus.FatalLevel)
	case "PANIC":
		log.SetLevel(logrus.PanicLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	// return the constructed logger with it's settings
	return &Logger{
		settings: settings,
		Logger:   log,
	}, nil
}

package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

type basicLog struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
	Time  string `json:"time"`
}

func TestGetInstance(t *testing.T) {
	firstInstance, err := GetInstance()
	if err != nil {
		t.Fatal(err)
	}
	secondInstance, err := GetInstance()
	if err != nil {
		t.Fatal(err)
	}
	if firstInstance != secondInstance {
		t.Error("GetInstance returned two different pointers!")
	}
}

func TestLoadEnvVar(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatal(err)
	}

	wrongKey := "MISSING_KEY"
	_, err = loadEnvVar(wrongKey)
	if err == nil {
		t.Errorf("Key [%s] should not exist!", wrongKey)
	}

	correctKey := "LOG_LEVEL"
	_, err = loadEnvVar(correctKey)
	if err != nil {
		t.Errorf("Key [%s] exists but wasn't found!", correctKey)
	}
}

func TestInitLogger(t *testing.T) {
	var settings *Settings
	settings, err := initLogger()
	if err != nil {
		t.Fatal(err)
	}

	// modify these if .env changes
	expectedLevel := "INFO"
	expectedFormat := "JSON"
	expectedOutput := "FILE"
	expectedPath := "test.log"
	if settings.level != expectedLevel {
		t.Errorf("initLogger returned [%s], expected [%s].", settings.level, expectedLevel)
	}
	if settings.format != expectedFormat {
		t.Errorf("initLogger returned [%s], expected [%s].", settings.format, expectedFormat)
	}
	if settings.output != expectedOutput {
		t.Errorf("initLogger returned [%s], expected [%s].", settings.output, expectedOutput)
	}
	if settings.path != expectedPath {
		t.Errorf("initLogger returned [%s], expected [%s].", settings.path, expectedPath)
	}
}

func TestCreateLogger(t *testing.T) {
	var settings *Settings

	settings, err := initLogger()
	if err != nil {
		t.Fatal(err)
	}
	_, err = createLogger(settings)
	if err != nil {
		t.Fatal(err)
	}
}

func catchOutput(log *Logger, msg, level string) (basicLog, error) {
	// Create a basic log sturcture and a byte buffer.
	// We change the output writer from `log` to the buffer.
	// After that, we create a simple log that gets captured,
	//  we unmarshal the buffer into the basicLog struct.
	// If erything worked out, we should be able to compare the message
	//  as well as the level.
	var output basicLog
	var buf bytes.Buffer

	log.SetOutput(&buf)
	switch strings.ToLower(level) {
	case "trace":
		log.Trace(msg)
	case "debug":
		log.Debug(msg)
	case "info":
		log.Info(msg)
	case "warn", "warning":
		log.Warn(msg)
	case "error":
		log.Error(msg)
	default:
		log.Debug(msg)
	}

	err := json.Unmarshal(buf.Bytes(), &output)
	if err != nil {
		return output, err
	}
	return output, nil
}

func TestCreateLoggerOutput(t *testing.T) {
	var settings *Settings

	settings, err := initLogger()
	if err != nil {
		t.Fatal(err)
	}

	settings.output = "TERM"
	settings.level = "TRACE"
	log, err := createLogger(settings)
	if err != nil {
		t.Fatal(err)
	}

	logs := []struct {
		msg string
		lvl string
	}{
		{"This will work", "info"},
		{"This needs to work", "warning"},
		{"This is fine", "trace"},
		{"Is this an error", "error"},
		{"BlaBla idk even more", "debug"},
		{"loremipsum", "debug"},
		{"we are the champions", "error"},
		{"i hope this breaks", "info"},
		{"I waNnA b3 tR4Cer!", "trace"},
		{"nerf bAsTion", "warning"},
		{"rollErcoaster blabla", "info"},
		{"unlock it !!! unlock it !!!", "debug"},
	}

	for _, elem := range logs {
		output, err := catchOutput(log, elem.msg, elem.lvl)
		if err != nil {
			t.Error(err)
		}
		if output.Level != elem.lvl {
			t.Errorf("Level does not match! Received {%s, %s}, expected {%s, %s}.", output.Msg, output.Level, elem.msg, elem.lvl)
		}
		if output.Msg != elem.msg {
			t.Errorf("Message does not match! Received {%s, %s}, expected {%s, %s}.", output.Msg, output.Level, elem.msg, elem.lvl)
		}
	}
}

func TestCreateLoggerLevel(t *testing.T) {
	var settings *Settings

	settings, err := initLogger()
	if err != nil {
		t.Fatal(err)
	}

	// Logrus levels in int
	//		panic: 0
	//		fatal: 1
	//		error: 2
	// 		warn:  3
	// 		info:  4
	//		debug: 5
	// 		trace: 6
	settings.level = "DEBUG"
	log, err := createLogger(settings)
	if err != nil {
		t.Fatal(err)
	}
	if log.Level != 5 {
		t.Errorf("Log level doesn't match, returned [%d], expected [%d].", log.Level, 4)
	}
}

package enscli

import (
	"io"
	"os"
	"testing"

	"github.com/EnsurityTechnologies/logger"
)

type TestCli struct {
	EnsCli
	log logger.Logger
	msg string
}

func (tc *TestCli) statusFunction() bool {
	tc.log.Info("CLI is running fine")
	return true
}

func (tc *TestCli) echoFunction() bool {
	tc.log.Info(tc.msg)
	return true
}

func TestBasic(t *testing.T) {
	fp, err := os.OpenFile("log.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	logOptions := &logger.LoggerOptions{
		Name:   "TestCLI",
		Color:  []logger.ColorOption{logger.AutoColor, logger.ColorOff},
		Output: []io.Writer{logger.DefaultOutput, fp},
	}

	log := logger.New(logOptions)
	tc := &TestCli{
		log: log,
	}
	cli, err := NewEnsCli("test", log.Named("cli"))
	if err != nil {
		t.Fatal("failed to create cli")
	}
	sf := &Function{
		Handler:    tc.statusFunction,
		Title:      "Status function",
		SuccessMsg: "Status returned successfully",
		FailureMsg: "Status failed",
	}
	ef := &Function{
		Handler:    tc.echoFunction,
		Title:      "Echo function",
		SuccessMsg: "Echo returned successfully",
		FailureMsg: "Echo failed",
	}
	cli.AddCommand("status", sf)
	cli.AddCommand("echo", ef)
	o := &Option{
		Type:    StringType,
		Flag:    "m",
		Ptr:     &tc.msg,
		Default: "Test message",
		Usage:   "Display message",
	}
	cli.AddOption(o)
	os.Args = make([]string, 1)
	os.Args[0] = "test"
	cli.Run()
	os.Args = make([]string, 2)
	os.Args[0] = "test"
	os.Args[1] = "test"
	cli.Run()
	os.Args = make([]string, 2)
	os.Args[0] = "test"
	os.Args[1] = "status"
	cli.Run()
	os.Args = make([]string, 4)
	os.Args[0] = "test"
	os.Args[1] = "echo"
	os.Args[2] = "-m"
	os.Args[3] = "CLI Test message"
	cli.Run()
}

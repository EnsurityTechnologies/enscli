package enscli

import (
	"flag"
	"os"
	"runtime"
	"strings"

	"github.com/EnsurityTechnologies/logger"
)

const (
	StringType  string = "string"
	IntType     string = "int"
	Int64Type   string = "int64"
	UIntType    string = "uint"
	UInt64Type  string = "uint64"
	Float64Type string = "float64"
	BoolType    string = "bool"
)

type CommandHandler func() bool

type Function struct {
	Handler    CommandHandler
	Title      string
	SuccessMsg string
	FailureMsg string
}

type Option struct {
	Type    string
	Flag    string
	Ptr     interface{}
	Default interface{}
	Usage   string
}

type EnsCli struct {
	name    string
	version string
	log     logger.Logger
	funcs   map[string]*Function
	options []*Option
}

func NewEnsCli(name string, log logger.Logger) (*EnsCli, error) {
	return &EnsCli{name: name, version: "0.0.1", log: log, funcs: make(map[string]*Function), options: make([]*Option, 0)}, nil
}

func (cli *EnsCli) SetVersion(version string) {
	cli.version = version
}

func (cli *EnsCli) AddCommand(cmd string, f *Function) {
	cmd = strings.ToLower(cmd)
	cli.funcs[cmd] = f
}

func (cli *EnsCli) AddOption(o *Option) {
	o.Type = strings.ToLower(o.Type)
	cli.options = append(cli.options, o)
}

func (cli *EnsCli) showHelp() {
	msg := "Command line helper\n\n"
	if runtime.GOOS == "windows" {
		msg = msg + cli.name + ".exe <cmd>\n\n"
	} else {
		msg = msg + cli.name + " <cmd>\n\n"
	}
	msg = msg + "Use the following commands\n\n"
	for k := range cli.funcs {
		msg = msg + "  " + k + "\n"
	}
	msg = msg + "\nSupported options\n\n"
	for _, o := range cli.options {
		msg = msg + "  " + o.Flag + " <" + o.Type + ">" + "  :  " + o.Usage + "\n"
	}
	cli.log.Info(msg)
}

func (cli *EnsCli) Run() {
	if len(os.Args) < 2 {
		cli.showHelp()
		return
	}
	cmd := strings.ToLower(os.Args[1])
	if cmd == "-h" || cmd == "-help" {
		cli.showHelp()
		return
	}
	if cmd == "-v" {
		cli.log.Info("Tool version : " + cli.version)
	}
	f, ok := cli.funcs[cmd]
	if !ok {
		cli.log.Error("Unsupported command, please check the helper")
		cli.showHelp()
		return
	}
	os.Args = os.Args[1:]
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	for _, o := range cli.options {
		switch o.Type {
		case StringType:
			ptr, ok := o.Ptr.(*string)
			if !ok {
				cli.log.Error("invalid pointer, expected string pointer", "flag", o.Flag)
				return
			}
			v, ok := o.Default.(string)
			if !ok {
				cli.log.Error("invalid default value, expected string", "flag", o.Flag)
				return
			}
			flag.StringVar(ptr, o.Flag, v, o.Usage)
		case IntType:
			ptr, ok := o.Ptr.(*int)
			if !ok {
				cli.log.Error("invalid pointer, expected integer pointer", "flag", o.Flag)
				return
			}
			v, ok := o.Default.(int)
			if !ok {
				cli.log.Error("invalid default value, expected integer", "flag", o.Flag)
				return
			}
			flag.IntVar(ptr, o.Flag, v, o.Usage)
		case Int64Type:
			ptr, ok := o.Ptr.(*int64)
			if !ok {
				cli.log.Error("invalid pointer, expected integer 64 pointer", "flag", o.Flag)
				return
			}
			v, ok := o.Default.(int64)
			if !ok {
				cli.log.Error("invalid default value, expected integer 64", "flag", o.Flag)
				return
			}
			flag.Int64Var(ptr, o.Flag, v, o.Usage)
		case UIntType:
			ptr, ok := o.Ptr.(*uint)
			if !ok {
				cli.log.Error("invalid pointer, expected unsigned integer pointer", "flag", o.Flag)
				return
			}
			v, ok := o.Default.(uint)
			if !ok {
				cli.log.Error("invalid default value, expected unsigned integer", "flag", o.Flag)
				return
			}
			flag.UintVar(ptr, o.Flag, v, o.Usage)
		case UInt64Type:
			ptr, ok := o.Ptr.(*uint64)
			if !ok {
				cli.log.Error("invalid pointer, expected unsigned integer 64 pointer", "flag", o.Flag)
				return
			}
			v, ok := o.Default.(uint64)
			if !ok {
				cli.log.Error("invalid default value, expected unsigned integer 64", "flag", o.Flag)
				return
			}
			flag.Uint64Var(ptr, o.Flag, v, o.Usage)
		case Float64Type:
			ptr, ok := o.Ptr.(*float64)
			if !ok {
				cli.log.Error("invalid pointer, expected float 64 pointer", "flag", o.Flag)
				return
			}
			v, ok := o.Default.(float64)
			if !ok {
				cli.log.Error("invalid default value, expected float 64", "flag", o.Flag)
				return
			}
			flag.Float64Var(ptr, o.Flag, v, o.Usage)
		case BoolType:
			ptr, ok := o.Ptr.(*bool)
			if !ok {
				cli.log.Error("invalid pointer, expected boolean pointer", "flag", o.Flag)
				return
			}
			v, ok := o.Default.(bool)
			if !ok {
				cli.log.Error("invalid default value, expected boolean", "flag", o.Flag)
				return
			}
			flag.BoolVar(ptr, o.Flag, v, o.Usage)
		}
	}
	flag.Parse()
	cli.log.Info("Executing the command : " + f.Title)
	if !f.Handler() {
		cli.log.Error(f.FailureMsg)
	} else {
		cli.log.Info(f.SuccessMsg)
	}
}

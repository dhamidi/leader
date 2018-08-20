package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/gobuffalo/packr"
)

// Revision contains the current revision (git commit hash) of this program.
//
// It is set during the build process using -ldflags.
var Revision string

// Release contains the current release (git tag pointing to Revision) of this program.
//
// It is set during the build process using -ldflags.
var Release string

// Shell defines the common functions that need to be supported by a
// shell in order for leader to support it.
type Shell interface {
	Commandline() (string, int)
	Init() string
	EvalNext(command string, path []rune) string
}

func main() {
	allSignals := make(chan os.Signal, 1)
	signal.Notify(allSignals)
	errorHandler := NewErrorLogger(os.Stderr)
	defer func() {
		v := recover()
		if err, ok := v.(error); ok {
			errorHandler.Print(err)
		}
	}()
	tty, err := NewTTY()
	if err != nil {
		errorHandler.Fatal(err)
	}

	shell := NewShellFromEnv(os.Getenv)
	line, cursor := shell.Commandline()
	shellParser := NewShellParser()
	if shellParser.InQuotedString(line, cursor) {
		os.Exit(3)
	}

	executor := NewShellExecutor("bash", "-c").Attach(tty.File())
	rootKeyMap := NewKeyMap("root")
	context := &Context{
		Files:         NewOnDiskFileSystem(),
		ErrorLogger:   errorHandler,
		CurrentKeyMap: rootKeyMap,
		Executor:      executor,
		Terminal:      tty,
		Shell:         shell,
	}
	go func() {
		<-allSignals
		tty.Restore()
		os.Exit(0)
	}()
	loadConfig := NewLoadConfig(context, os.Getenv("PWD"), os.Getenv("HOME"))
	errorHandler.Must(loadConfig.Execute)
	parseArgs(context, os.Args)
	errorHandler.Must(tty.MakeRaw)
	selectMenuEntry := NewSelectMenuEntry(context)
	errorHandler.Print(selectMenuEntry.Execute())
}

func parseArgs(context *Context, args []string) {
	if len(args) == 1 {
		return
	}

	if args[1] == "init" {
		initShell(context)
		os.Exit(0)
	}

	if args[1] == "list-keys" {
		context.ErrorLogger.Print(NewListKeys(context).Execute())
		os.Exit(0)
	}

	if args[1] == "version" {
		showVersion()
		os.Exit(0)
	}

	if args[1] == "help" {
		showHelp(context)
		os.Exit(0)
	}

	for i := 0; i < len(args); i++ {
		if args[i] == "print" {
			context.Executor = NewPrintingExecutor(context, os.Stdout)
			continue
		}
		if args[i][0] == '@' {
			navigateTo(context, []rune(args[i][1:]))
		}
	}
}

func navigateTo(context *Context, path []rune) {
	for _, keyRune := range path {
		binding := context.CurrentKeyMap.LookupKey(keyRune)
		if binding == UnboundKey {
			os.Exit(1)
		}
		if !binding.HasChildren() {
			if binding.IsBoundToCommand() {
				context.ErrorLogger.Print(binding.Execute())
				os.Exit(0)
			}
			os.Exit(2)
		}
		context.KeyPath = append(context.KeyPath, keyRune)
		context.Navigate(binding.Children())
	}

}

func showVersion() {
	if Release != "" {
		fmt.Printf("%s\n", Release)
	} else if Revision != "" {
		fmt.Printf("%s\n", Revision)
	}

}
func showHelp(context *Context) {
	assets := packr.NewBox("./assets")
	man := exec.Command("man", "-l", "-")
	man.Stdout = os.Stdout
	man.Stdin = bytes.NewBufferString(assets.String("leader.1"))
	man.Run()
}

func initShell(context *Context) {
	fmt.Printf("%s\n", context.Shell.Init())
}

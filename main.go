package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/gobuffalo/packr"
)

func main() {
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

	line, cursor := getCurrentInputLine()
	shellParser := NewShellParser()
	if shellParser.InQuotedString(line, cursor) {
		os.Exit(3)
	}

	executor := NewShellExecutor("bash", "-c").Attach(tty.File())
	rootKeyMap := NewKeyMap("root")
	context := &Context{
		ErrorLogger:   errorHandler,
		CurrentKeyMap: rootKeyMap,
		Executor:      executor,
		Terminal:      tty,
	}
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
		initShell()
		os.Exit(0)
	}

	if args[1] == "list-keys" {
		context.ErrorLogger.Print(NewListKeys(context).Execute())
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
			os.Exit(2)
		}
		context.KeyPath = append(context.KeyPath, keyRune)
		context.Navigate(binding.Children())
	}

}

func showHelp(context *Context) {
	assets := packr.NewBox("./assets")
	man := exec.Command("man", "-l", "-")
	man.Stdout = os.Stdout
	man.Stdin = bytes.NewBufferString(assets.String("leader.1"))
	man.Run()
}

func getCurrentInputLine() (string, int) {
	shellName := filepath.Base(os.Getenv("SHELL"))
	switch shellName {
	case "zsh":
		return getCurrentInputLineZSH()
	case "bash":
		return getCurrentInputLineBash()
	default:
		return "", 0
	}
}

func getCurrentInputLineZSH() (string, int) {
	line := os.Getenv("BUFFER")
	point, _ := strconv.Atoi(os.Getenv("CURSOR"))
	return line, point
}

func getCurrentInputLineBash() (string, int) {
	readlineLine := os.Getenv("READLINE_LINE")
	readlinePoint, _ := strconv.Atoi(os.Getenv("READLINE_POINT"))
	return readlineLine, readlinePoint
}

func initShell() {
	initFiles := packr.NewBox("./assets")
	shellName := filepath.Base(os.Getenv("SHELL"))
	switch shellName {
	case "zsh":
		fmt.Printf("%s\n", initFiles.String("leader.zsh.sh"))
	case "bash":
		fmt.Printf("%s\n", initFiles.String("leader.bash.sh"))
	default:
		fmt.Fprintf(os.Stderr, "Shell %s not supported!\n", shellName)
	}
}

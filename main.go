package main

import "os"

func main() {
	errorHandler := NewErrorLogger(os.Stderr)
	defer func() {
		v := recover()
		if err, ok := v.(error); ok {
			errorHandler.Print(err)
		}
	}()
	tty, err := NewTerminalTTY()
	if err != nil {
		errorHandler.Fatal(err)
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

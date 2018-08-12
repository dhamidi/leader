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
	rootKeyMap := NewKeyMap("root")
	context := &Context{
		CurrentKeyMap: rootKeyMap,
		Executor:      NewShellExecutor("bash", "-c").Attach(tty.File()),
		Terminal:      tty,
	}
	rootKeyMap.Bind('d').Do(NewRunShellCommand(context, "date").Execute).Describe("date")
	errorHandler.Must(tty.MakeRaw)
	selectMenuEntry := NewSelectMenuEntry(context)
	errorHandler.Print(selectMenuEntry.Execute())
}

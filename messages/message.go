package messages

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"io"
	"os"
)

var stdout = colorable.NewColorableStdout()

var C = aurora.NewAurora(isatty.IsTerminal(os.Stdout.Fd()))

func Output() io.Writer {
	return stdout
}

func Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(stdout, a...)
}

func Print(a ...interface{}) (n int, err error) {
	return fmt.Fprint(stdout, a...)
}

func Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(stdout, format, a...)
}

func Printfln(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(stdout, format+"\n", a...)
}

func Error(str string) {
	Printfln("%s: %s", C.Red("Error"), str)
}

func Errorf(format string, a ...interface{}) {
	Printf("%s: ", C.Red("Error"))
	Printfln(format, a...)
}

func Fatal(str string) {
	Printfln("%s: %s", C.Red("Error"), str)
	os.Exit(1)
}

func Fatalf(format string, a ...interface{}) {
	Printf("%s: ", C.Red("Error"))
	Printfln(format, a...)
	os.Exit(1)
}

func Warning(str string) {
	Printfln("%s: %s", C.Magenta("Warning"), str)
}

func Warningf(format string, a ...interface{}) {
	Printf("%s: ", C.Brown("Warning"))
	Printfln(format, a...)
}

package logger

import "fmt"

type Locator struct {
	scope string
}

func DefaultLocator() Locator {
	return Locator{scope: ""}
}

func (l Locator) GetChild(scope string) Locator {
	return Locator{scope: fmt.Sprintf("%s.%s", l.scope, scope)}
}

type Logger struct {
	file string
}

func DefaultLogger() Logger {
	return Logger{}
}

const (
	ERROR   = "ERROR"
	WARNING = "WARNING"
	INFO    = "INFO"
	PANIC   = "PANIC"
)

func (l *Logger) Infof(locator Locator, base string, args ...interface{}) {
	l.printInternal(INFO, locator, base, args...)
}

func (l *Logger) Errorf(locator Locator, base string, args ...interface{}) {
	l.printInternal(ERROR, locator, base, args...)
}

func (l *Logger) Warningf(locator Locator, base string, args ...interface{}) {
	l.printInternal(WARNING, locator, base, args...)
}

func (l *Logger) Panicf(locator Locator, base string, args ...interface{}) {
	panic(l.printInternal(PANIC, locator, base, args...))
}

func (l *Logger) printInternal(prefix string, locator Locator, base string, args ...interface{}) string {
	msg := fmt.Sprintf("[%s] %s: %s",
		prefix, locator.scope, fmt.Sprintf(base, args...))
	fmt.Printf("%s/n", msg)
	return msg
}

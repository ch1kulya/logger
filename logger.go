package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var (
	level    = LevelDebug
	mu       sync.RWMutex
	output   io.Writer = os.Stdout
	errOut   io.Writer = os.Stderr
	exitFunc           = os.Exit
)

var (
	dimStyle     = lipgloss.NewStyle().Faint(true)
	debugStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	fatalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	greenStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	cyanStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	yellowStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	magentaStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	redStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	grayStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	whiteStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
)

func SetLevel(l Level) {
	mu.Lock()
	defer mu.Unlock()
	level = l
}

func SetOutput(w io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	output = w
	errOut = w
}

func getLevel() Level {
	mu.RLock()
	defer mu.RUnlock()
	return level
}

func getTimestamp() string {
	return dimStyle.Render(time.Now().Format("15:04:05"))
}

func Info(format string, v ...any) {
	if getLevel() > LevelInfo {
		return
	}
	msg := fmt.Sprintf(format, v...)
	fmt.Fprintf(output, "%s %s %s\n", getTimestamp(), infoStyle.Width(7).Render("INFO"), msg)
}

func Error(format string, v ...any) {
	if getLevel() > LevelError {
		return
	}
	msg := fmt.Sprintf(format, v...)
	fmt.Fprintf(errOut, "%s %s %s\n", getTimestamp(), errorStyle.Width(7).Render("ERR"), msg)
}

func Warn(format string, v ...any) {
	if getLevel() > LevelWarn {
		return
	}
	msg := fmt.Sprintf(format, v...)
	fmt.Fprintf(output, "%s %s %s\n", getTimestamp(), warnStyle.Width(7).Render("WARN"), msg)
}

func Debug(format string, v ...any) {
	if getLevel() > LevelDebug {
		return
	}
	msg := fmt.Sprintf(format, v...)
	fmt.Fprintf(output, "%s %s %s\n", getTimestamp(), debugStyle.Width(7).Render("DEBUG"), msg)
}

func Fatal(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	fmt.Fprintf(errOut, "%s %s %s\n", getTimestamp(), fatalStyle.Width(7).Render("FATAL"), msg)
	exitFunc(1)
}

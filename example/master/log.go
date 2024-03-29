package main

import (
	"fmt"
	"log"

	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-colorable"
)

var au aurora.Aurora

func init() {
	log.SetOutput(colorable.NewColorableStdout())
}

func loggerInit(colorLogs bool) {
	au = aurora.NewAurora(colorLogs)
}

func serverColor(input string) uint8 {
	o := byte(0)
	for _, c := range input {
		o += byte(c)
	}

	return (((o % 36) * 36) + (o % 6) + 16) % 255
}

func componentColor(input string) aurora.Color {
	switch input {
	case "shutdown":
		fallthrough
	case "startup":
		return aurora.BrightFg | aurora.MagentaFg

	case "server":
		return aurora.BrightFg | aurora.CyanFg

	case "config":
		return aurora.BrightFg | aurora.YellowFg

	case "maintenance":
		return aurora.BrightFg | aurora.GreenFg

	default:
		return aurora.BrightFg | aurora.WhiteFg
	}
}

func LogServer(server string, format string, args ...interface{}) {
	color := serverColor(server)
	s := fmt.Sprintf("[%s]: %s\n", au.Index(color, server), au.Index(color, format))
	log.Printf(s, args...)
}

func LogServerAlert(server string, format string, args ...interface{}) {
	color := serverColor(server)
	s := fmt.Sprintf("[%s]: %s %s\n", au.Index(color, server), au.Red("!"), au.Yellow(format))
	log.Printf(s, args...)
}

func LogComponent(component string, format string, args ...interface{}) {
	color := componentColor(component)
	s := fmt.Sprintf("{%s}: %s\n", au.Colorize(component, color), au.Colorize(format, color))
	log.Printf(s, args...)
}

func LogComponentAlert(component string, format string, args ...interface{}) {
	color := componentColor(component)
	s := fmt.Sprintf("{%s}: %s %s\n", au.Colorize(component, color), au.Red("!"), au.Yellow(format))
	log.Printf(s, args...)
}

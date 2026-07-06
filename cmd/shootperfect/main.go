package main

import (
	"os"

	"github.com/nikhilkumar961987/shootperfect-core/internal/logger"
)

func main() {
	log := logger.New()

	if len(os.Args) < 2 {
		printHelp(log)
		return
	}

	command := os.Args[1]

	switch command {
	case "version":
		log.Info("ShootPerfect Core", "version", "v0.1.0")

	case "analyze":
		log.Info("analyze command not implemented yet")

	case "serve":
		log.Info("serve command not implemented yet")

	default:
		log.Error("unknown command", "command", command)
		printHelp(log)
	}
}

func printHelp(log interface {
	Info(msg string, args ...any)
}) {
	log.Info("ShootPerfect Core")
	log.Info("Usage",
		"version", "shootperfect version",
		"analyze", "shootperfect analyze --session <path>",
		"serve", "shootperfect serve",
	)
}

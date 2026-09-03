package main

import (
	"cyberstrike-ai/internal/config"
	"cyberstrike-ai/internal/logger"
	"cyberstrike-ai/internal/mcp"
	"cyberstrike-ai/internal/security"
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
)

func main() {
	var configPath = flag.String("config", "config.yaml", "Path to the configuration file")
	flag.Parse()

	// Load the configuration.
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logging on stderr in stdio mode to avoid interfering with JSON-RPC communication.
	log := logger.New(cfg.Log.Level, "stderr")

	// Create the MCP server.
	mcpServer := mcp.NewServer(log.Logger)

	// Create the secure tool executor.
	executor := security.NewExecutor(&cfg.Security, mcpServer, log.Logger)

	// Register tools.
	executor.RegisterTools(mcpServer)
	mcp.RegisterExecutionControlTools(mcpServer, nil)

	log.Logger.Info("MCP server (stdio mode) started; waiting for messages...")

	// Run the stdio loop.
	if err := mcpServer.HandleStdio(); err != nil {
		log.Logger.Error("MCP server failed", zap.Error(err))
		os.Exit(1)
	}
}

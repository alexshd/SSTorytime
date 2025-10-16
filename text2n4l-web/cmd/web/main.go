package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"text2n4l-web/internal/web"

	"github.com/arl/statsviz"
)

func main() {
	// Initialize structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("Starting N4L API server")

	mux := http.NewServeMux()

	// Register statsviz at /debug/statsviz/
	if err := statsviz.Register(mux); err != nil {
		logger.Error("Failed to register statsviz", "error", err)
		os.Exit(1)
	}

	// API endpoints
	mux.HandleFunc("/api/convert", web.ConvertHandler)
	mux.HandleFunc("/api/convert/stream", web.StreamingConvertHandler)

	// Profiling endpoints (accessible at /debug/pprof/)
	// For CPU profiling: curl http://localhost:8080/debug/pprof/profile?seconds=30 -o cpu.prof
	// For heap profiling: curl http://localhost:8080/debug/pprof/heap -o heap.prof
	// For goroutine info: curl http://localhost:8080/debug/pprof/goroutine
	mux.Handle("/debug/", http.DefaultServeMux)

	port := ":5050"
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to signal server shutdown
	done := make(chan bool, 1)

	// Start server in a goroutine
	go func() {
		logger.Info("Server listening", "address", "http://localhost"+port)
		printEndpoints()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
		done <- true
	}()

	// Start interactive command handler
	go handleCommands(ctx, server, logger)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Server started successfully. Press Ctrl+C to stop, or use interactive commands:")
	logger.Info("  R - Show routes")
	logger.Info("  L - Show listening addresses")

	select {
	case <-quit:
		logger.Info("Shutdown signal received, shutting down gracefully...")
	case <-ctx.Done():
		logger.Info("Context cancelled, shutting down...")
	}

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	} else {
		logger.Info("Server shutdown gracefully")
	}

	// Wait for server goroutine to finish
	<-done
	logger.Info("Server stopped")
}

// printEndpoints displays all available endpoints
func printEndpoints() {
	fmt.Println("\n📋 Available Endpoints:")
	fmt.Println("  🌐 API Endpoints:")
	fmt.Println("    POST /api/convert        - Convert text to N4L format (buffered)")
	fmt.Println("    POST /api/convert/stream - Convert text to N4L format (streaming)")
	fmt.Println("  📊 Profiling Endpoints:")
	fmt.Println("    GET  /debug/statsviz/    - Real-time profiling dashboard")
	fmt.Println("    GET  /debug/pprof/       - Profiling index")
	fmt.Println("    GET  /debug/pprof/profile - CPU profile")
	fmt.Println("    GET  /debug/pprof/heap   - Heap profile")
	fmt.Println("    GET  /debug/pprof/goroutine - Goroutine info")
	fmt.Println()
}

// handleCommands processes interactive keyboard commands
func handleCommands(ctx context.Context, server *http.Server, logger *slog.Logger) {
	scanner := bufio.NewScanner(os.Stdin)
	logger.Info("Interactive command mode active. Type 'R' for routes, 'L' for addresses, 'Q' to quit")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if scanner.Scan() {
				cmd := strings.ToUpper(strings.TrimSpace(scanner.Text()))

				switch cmd {
				case "R":
					logger.Info("Showing routes...")
					printEndpoints()

				case "L":
					logger.Info("Showing listening addresses...")
					showListeningAddresses(server, logger)

				case "Q":
					logger.Info("Quit command received")
					os.Exit(0)

				case "":
					// Ignore empty input
					continue

				default:
					logger.Info("Unknown command. Available: R (routes), L (addresses), Q (quit)")
				}
			}

			if err := scanner.Err(); err != nil {
				logger.Error("Error reading input", "error", err)
				return
			}
		}
	}
}

// showListeningAddresses displays all network interfaces and ports the server is listening on
func showListeningAddresses(server *http.Server, logger *slog.Logger) {
	host, port, err := net.SplitHostPort(server.Addr)
	if err != nil {
		logger.Error("Failed to parse server address", "error", err)
		return
	}

	fmt.Printf("\n🌐 Server listening on:\n")

	if host == "" || host == "0.0.0.0" {
		// Server is listening on all interfaces
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			logger.Error("Failed to get network interfaces", "error", err)
			fmt.Printf("  http://localhost:%s\n", port)
			return
		}

		found := false
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					fmt.Printf("  http://%s:%s\n", ipnet.IP.String(), port)
					found = true
				}
			}
		}

		if !found {
			fmt.Printf("  http://localhost:%s\n", port)
		}
	} else {
		fmt.Printf("  http://%s:%s\n", host, port)
	}

	fmt.Println()
}

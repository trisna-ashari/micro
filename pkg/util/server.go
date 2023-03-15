package util

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"micro/pkg/logger"
)

// RunGRPCServerWithGracefulShutdown knows how to run and gracefully shutdown the grpc.Server.
func RunGRPCServerWithGracefulShutdown(s *grpc.Server, port string, logStd *logger.Logger) error {
	logStd.Log.Info("GRPC server is starting ...")

	// Initialize listener to predefined port. It uses to start
	// the grpc.Server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverError := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		logStd.Log.Infof("GRPC server is running in port: %v", port)
		serverError <- s.Serve(lis)
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdownListenerChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownListenerChannel, syscall.SIGINT, syscall.SIGTERM)

	// Blocking and waiting for shutdown or error from server.
	select {
	case err := <-serverError:
		if err != nil {
			logStd.Log.Errorf("GRPC server cannot started, err: %v", err)
			return err
		}
	case sig := <-shutdownListenerChannel:
		logStd.Log.Infof("GRPC server shutdown by signal: %v", sig)
		s.GracefulStop()
		logStd.Log.Info("GRPC server was shutting down gracefully")
	}

	return nil
}

// RunHTTPServerWithGracefulShutdown knows how to run and gracefully shutdown the HTTP Handler.
func RunHTTPServerWithGracefulShutdown(handler http.Handler, addr string, shutdownTimeout time.Duration, logStd *logger.Logger) error {
	logStd.Log.Info("HTTP server is starting ...")

	// Create a server
	server := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverError := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		logStd.Log.Infof("HTTP server is running on port: %v", server.Addr)
		serverError <- server.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdownListenerChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownListenerChannel, syscall.SIGINT, syscall.SIGTERM)

	// Blocking and waiting for shutdown or error from server.
	select {
	case err := <-serverError:
		if err != nil {
			logStd.Log.Errorf("Cannot start the HTTP server: %v", err)
			return err
		}
	case sig := <-shutdownListenerChannel:
		logStd.Log.Infof("HTTP server shutdown by signal: %v", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logStd.Log.Infof("HTTP server was shutting down forcibly")
			err = server.Close()
			return err
		}

		logStd.Log.Info("HTTP server was shutting down gracefully")
	}

	return nil
}

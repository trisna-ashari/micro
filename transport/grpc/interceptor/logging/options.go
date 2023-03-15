package logging

import (
	"google.golang.org/grpc"
	"micro/pkg/logger"
)

// UnaryOption return grpc.UnaryServerInfo with UnaryOption.
type UnaryOption func(info *grpc.UnaryServerInfo)

// WithUnaryLogger is a function uses to print log when unary gRPC call is coming.
func WithUnaryLogger(logger *logger.Logger, serviceName string) UnaryOption {
	return func(info *grpc.UnaryServerInfo) {
		logger.Log.Info(serviceName + " - serving incoming gRPC call for method " + info.FullMethod)
	}
}

// StreamOption return grpc.StreamServerInfo with StreamOption.
type StreamOption func(info *grpc.StreamServerInfo)

// WithStreamLogger is a function uses to print log when stream gRPC call is coming.
func WithStreamLogger(logger *logger.Logger, serviceName string) StreamOption {
	return func(info *grpc.StreamServerInfo) {
		logger.Log.Info(serviceName + " - serving incoming gRPC call for method " + info.FullMethod)
	}
}

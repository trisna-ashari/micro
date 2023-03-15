package logging

import (
	"context"
	"micro/transport/grpc/interceptor"

	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns a new unary server interceptor with metadata.
func UnaryServerInterceptor(options ...UnaryOption) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		for _, opt := range options {
			opt(info)
		}

		resp, err := handler(ctx, req)

		return resp, err
	}
}

// StreamServerInterceptor returns a new stream server interceptor with metadata.
func StreamServerInterceptor(options ...StreamOption) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		for _, opt := range options {
			opt(info)
		}

		wrapped := interceptor.WrapServerStream(ss)

		return handler(srv, wrapped)
	}
}

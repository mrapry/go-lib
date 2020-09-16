package middleware

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mrapry/go-lib/config"
	"github.com/mrapry/go-lib/golibhelper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// GRPCBasicAuth function,
// or Unary interceptor
// additional security for our GRPC server
func (m *Middleware) GRPCBasicAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	defer func() {
		m.grpcLog(start, err, info.FullMethod, "GRPC")
	}()

	e := m.validateGrpcAuth(ctx)
	if e != nil {
		return nil, e
	}

	resp, err = handler(ctx, req)
	return
}

// GRPCBasicAuthStream interceptor
func (m *Middleware) GRPCBasicAuthStream(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	start := time.Now()
	defer func() {
		m.grpcLog(start, err, info.FullMethod, "GRPC-STREAM")
	}()

	if err := m.validateGrpcAuth(stream.Context()); err != nil {
		return err
	}

	return handler(srv, stream)
}

// validateGrpcAuth auth from incoming context
func (m *Middleware) validateGrpcAuth(ctx context.Context) error {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "missing context metadata")
	}

	authorizationMap := meta["authorization"]
	if len(authorizationMap) != 1 {
		return grpc.Errorf(codes.Unauthenticated, "Invalid authorization")
	}

	authorization := authorizationMap[0]
	if err := m.Basic(ctx, authorization); err != nil {
		return grpc.Errorf(codes.Unauthenticated, err.Error())
	}

	return nil
}

// Log incoming grpc request
func (m *Middleware) grpcLog(startTime time.Time, err error, fullMethod string, reqType string) {
	end := time.Now()
	var status = "OK"
	statusColor := golibhelper.Green
	if err != nil {
		statusColor = golibhelper.Red
		status = "ERROR"
	}

	fmt.Fprintf(os.Stdout, "%s[%s]%s :%d %v | %s %-5s %s | %13v | %s\n",
		golibhelper.Cyan, reqType, golibhelper.Reset, config.BaseEnv().GRPCPort,
		end.Format("2006/01/02 - 15:04:05"),
		statusColor, status, golibhelper.Reset,
		end.Sub(startTime),
		fullMethod,
	)
}

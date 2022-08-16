// Package interceptors provides middleware used for gRPC request handling.

package interceptors

import (
	"context"
	"fmt"
	"net"

	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Logging - interceptor for logging requests and latency.
func Logging(ctx context.Context, method string, req interface{},
	reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) (err error) {
	start := time.Now()

	err = invoker(ctx, method, req, reply, cc, opts...)

	log.Debug().Fields(map[string]interface{}{"method": method, "latency": fmt.Sprintf("%v", time.Since(start)), "err": err}).Msg("success")

	return err
}

// CheckIP - interceptor validates trusted subnet.
func CheckIP(s *net.IPNet) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

		if s != nil {
			md, ok := metadata.FromIncomingContext(ctx)
			if ok {
				values := md.Get("X-Real-IP")
				if len(values) > 0 {
					log.Debug().Fields(map[string]interface{}{"real-ip": values[0]}).Msgf("CheckIP")
					if !s.Contains(net.ParseIP(values[0])) {
						return nil, status.Errorf(codes.PermissionDenied, `restricted IP address: %s`, values[0])
					}
				}
			}
		}
		return handler(ctx, req)
	}
}

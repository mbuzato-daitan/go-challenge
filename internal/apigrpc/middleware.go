package apigrpc

import (
	context "context"
	"fmt"
	"os"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	if err := mdwAuthorize(ctx); err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func StreamServerInterceptor(srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {

	if err := mdwAuthorize(ss.Context()); err != nil {
		return err
	}

	return handler(srv, ss)
}

func mdwAuthorize(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("grpc.MdwAuthorize: unauthorized")
	}

	tokens := md.Get("Authorization")
	if len(tokens) == 0 {
		return fmt.Errorf("grpc.MdwAuthorize: unauthorized")
	}

	token := tokens[0]
	if token != "Bearer "+os.Getenv("API_KEY") {
		return fmt.Errorf("grpc.MdwAuthorize: unauthorized")
	}

	return nil
}

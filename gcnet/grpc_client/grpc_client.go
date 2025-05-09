package grpc_client

import (
	"context"
	"fmt"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	options *GrpcOption

	ctx    context.Context
	cancel context.CancelFunc

	client *grpc.ClientConn
}

func NewGrpcClient() *GrpcClient {
	s := &GrpcClient{
		options: NewGrpcOption(),
	}

	return s
}

func (s *GrpcClient) Init(ctx context.Context, option ...any) (err error) {
	s.ctx, s.cancel = context.WithCancel(ctx)

	for _, opt := range option {
		opt.(Option)(s.options)
	}

	credentials := grpc.WithTransportCredentials(insecure.NewCredentials())
	linkAddr := fmt.Sprintf("passthrough:%s", s.options.listenAddr)
	//linkAddr := fmt.Sprintf("%s", s.options.listenAddr)
	s.client, err = grpc.NewClient(linkAddr, credentials)
	//s.client, err = grpc.NewClient(linkAddr)
	if err != nil {
		log.Error("grpc dial err", zap.Error(err))
		return
	}

	return nil
}

func (s *GrpcClient) GetClient() *grpc.ClientConn {
	return s.client
}

func (s *GrpcClient) GetCtx() context.Context {
	return s.ctx
}

func (s *GrpcClient) Start() {

}

func (s *GrpcClient) Stop() {
	s.client.Close()
}

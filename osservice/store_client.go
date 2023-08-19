package osservice

import (
	"context"
	"fmt"
	"time"

	"osservice/osproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ Store = (*StoreClient)(nil)

type StoreClient struct {
	timeout time.Duration
	conn    *grpc.ClientConn
	client  osproto.StoreClient
}

func NewStoreClient(addr string, port int, timeout time.Duration) (*StoreClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &StoreClient{
		timeout: timeout,
		conn:    conn,
		client:  osproto.NewStoreClient(conn),
	}, nil
}

func (l *StoreClient) ctx() (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(l.timeout))
}

func (l *StoreClient) Load(key string) []byte {
	ctx, cancel := l.ctx()
	defer cancel()

	data, err := l.client.Load(ctx, wrapperspb.String(key))
	if err != nil {
		return nil
	}
	return data.GetValue()
}

func (l *StoreClient) Save(key string, data []byte) error {
	ctx, cancel := l.ctx()
	defer cancel()

	req := osproto.SaveRequest{
		Key:   wrapperspb.String(key),
		Bytes: wrapperspb.Bytes(data),
	}

	_, err := l.client.Save(ctx, &req)
	return err
}

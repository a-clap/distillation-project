package osservice

import (
	"context"
	"fmt"
	"time"

	"osservice/osproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var (
	_ Time = (*TimeClient)(nil)
)

type TimeClient struct {
	timeout time.Duration
	conn    *grpc.ClientConn
	client  osproto.TimeClient
}

func NewTimeClient(addr string, port int, timeout time.Duration) (*TimeClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &TimeClient{
		timeout: timeout,
		conn:    conn,
		client:  osproto.NewTimeClient(conn),
	}, nil
}

func (t *TimeClient) ctx() (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(t.timeout))
}

func (t *TimeClient) Now() (time.Time, error) {
	ctx, cancel := t.ctx()
	defer cancel()

	ts, err := t.client.Now(ctx, &emptypb.Empty{})
	if err != nil {
		return time.Time{}, err
	}

	return ts.AsTime(), nil
}

func (t *TimeClient) NTP() (bool, error) {
	ctx, cancel := t.ctx()
	defer cancel()

	ntp, err := t.client.NTP(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}

	return ntp.Value, nil
}

func (t *TimeClient) SetNTP(enable bool) error {
	ctx, cancel := t.ctx()
	defer cancel()

	_, err := t.client.SetNTP(ctx, wrapperspb.Bool(enable))
	return err
}

func (t *TimeClient) SetNow(now time.Time) error {
	ctx, cancel := t.ctx()
	defer cancel()

	_, err := t.client.SetNow(ctx, timestamppb.New(now))
	return err
}

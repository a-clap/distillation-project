package osservice

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Time interface {
	Now() (time.Time, error)
	SetNow(now time.Time) error
	NTP() (bool, error)
	SetNTP(enable bool) error
}

func (o *Os) Now(context.Context, *empty.Empty) (*timestamp.Timestamp, error) {
	t, err := o.time.Now()
	if err != nil {
		return nil, err
	}
	return timestamppb.New(t), nil
}

func (o *Os) NTP(context.Context, *empty.Empty) (*wrapperspb.BoolValue, error) {
	t, err := o.time.NTP()
	if err != nil {
		return nil, err
	}
	return wrapperspb.Bool(t), nil
}

func (o *Os) SetNow(_ context.Context, ts *timestamp.Timestamp) (*empty.Empty, error) {
	return &empty.Empty{}, o.time.SetNow(ts.AsTime())
}

func (o *Os) SetNTP(_ context.Context, value *wrapperspb.BoolValue) (*empty.Empty, error) {
	err := o.time.SetNTP(value.GetValue())
	return &empty.Empty{}, err
}

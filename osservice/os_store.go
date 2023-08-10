package osservice

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"osservice/osproto"
)

type Store interface {
	Load(key string) []byte
	Save(key string, data []byte) error
}


func (o *Os) Load(_ context.Context, value *wrappers.StringValue) (*wrapperspb.BytesValue, error) {
	v := o.store.Load(value.GetValue())
	return wrapperspb.Bytes(v), nil
}

func (o *Os) Save(_ context.Context, request *osproto.SaveRequest) (*empty.Empty, error) {
	err := o.store.Save(request.Key.GetValue(), request.Bytes.GetValue())
	return &empty.Empty{}, err
}

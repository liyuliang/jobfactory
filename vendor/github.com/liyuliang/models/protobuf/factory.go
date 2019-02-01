package protobuf

import (
	"strings"
	"errors"
	"github.com/golang/protobuf/proto"
	"time"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func Get(modelName string) (Model, error) {

	l := List()

	for name, creator := range l {

		if strings.ToLower(name) == strings.ToLower(modelName) {
			m := creator()
			return m, nil
		}
	}
	return nil, errors.New("can not get the Model by name :" + modelName)
}

func Marshal(model proto.Message) ([]byte, error) {
	return proto.Marshal(model)
}

func Unmarshal(buf []byte, model proto.Message) error {

	return proto.Unmarshal(buf, model)
}

func Uint64(v uint64) *uint64 {
	return proto.Uint64(v)
}

func Uint32(v uint32) *uint32 {
	return proto.Uint32(v)
}

func Int(v int) *int32 {
	return proto.Int(v)
}

func Int32(v int32) *int32 {
	return proto.Int32(v)
}

func Int64(v int64) *int64 {
	return proto.Int64(v)
}

func Float32(v float32) *float32 {
	return proto.Float32(v)
}

func Float64(v float64) *float64 {
	return proto.Float64(v)
}

func String(v string) *string {
	return proto.String(v)
}

func Bool(v bool) *bool {
	return proto.Bool(v)
}

func Timestamp(v time.Time) (*timestamp.Timestamp, error) {
	return ptypes.TimestampProto(v)
}

package protobuf

import (
	"strings"
	"errors"
)

type Model interface {
	GetId() uint64
	Reset()
	String() string
	ProtoMessage()
	Descriptor() ([]byte, []int)
}

func CheckModel(pfModel Model, data []byte) (err error) {

	err = Unmarshal(data, pfModel)

	s1 := pfModel.String()
	pfModel.Reset()
	s2 := pfModel.String()
	if s1 == s2 {
		err = errors.New("Empty Data In Protobuf Model. ")
	}
	return err
}

func IsModelNameExist(modelName string) bool {

	exist := false
	for name, _ := range List() {

		if strings.ToLower(modelName) == strings.ToLower(name) {
			exist = true
			break
		}
	}
	return exist
}

package worker

import (
	"github.com/liyuliang/models/protobuf"
)

type Model struct {
	Id    uint64         `json:"Id,omitempty"`
	Name  string         `json:"Name"`
	Model protobuf.Model `json:"Model"`
}
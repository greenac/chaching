package chaching_kafka

import (
	"github.com/greenac/chaching/internal/consts"
	"github.com/segmentio/kafka-go"
)

type Topic struct {
	Name              consts.TopicName
	Partitions        int
	ReplicationFactor int
}

type IDialer interface {
	Dial(network string, address string) (*kafka.Conn, error)
}

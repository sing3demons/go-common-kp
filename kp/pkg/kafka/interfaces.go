package kafka

import (
	"context"
	"net"

	"github.com/segmentio/kafka-go"
)


type Reader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	FetchMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
	Stats() kafka.ReaderStats
	Close() error
}

type Writer interface {
	WriteMessages(ctx context.Context, msg ...kafka.Message) error
	Close() error
	Stats() kafka.WriterStats
}

type Connection interface {
	Controller() (broker kafka.Broker, err error)
	CreateTopics(topics ...kafka.TopicConfig) error
	DeleteTopics(topics ...string) error
	RemoteAddr() net.Addr
	ReadPartitions(topics ...string) (partitions []kafka.Partition, err error)
	Close() error
}

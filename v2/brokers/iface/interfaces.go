package iface

import (
	"context"

	"github.com/RichardKnop/machinery/v2/config"
	"github.com/RichardKnop/machinery/v2/tasks"
)

// Broker - a common interface for all brokers
// 定义了各种broker需要实现的操作
type Broker interface {
	GetConfig() *config.Config             // 获取broker的配置
	SetRegisteredTaskNames(names []string) // broker.registeredTaskNames.items
	IsTaskRegistered(name string) bool
	StartConsuming(consumerTag string, concurrency int, p TaskProcessor) (bool, error) // 消费者标签用来标识不同的消费者，例如redis实例，p是TaskProcessor实例，用来处理消费到的任务
	StopConsuming()
	Publish(ctx context.Context, task *tasks.Signature) error
	GetPendingTasks(queue string) ([]*tasks.Signature, error)
	GetDelayedTasks() ([]*tasks.Signature, error)
	AdjustRoutingKey(s *tasks.Signature)
}

// TaskProcessor - can process a delivered task
// This will probably always be a worker instance
type TaskProcessor interface {
	Process(signature *tasks.Signature) error
	CustomQueue() string
	PreConsumeHandler() bool
}

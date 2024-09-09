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
	StopConsuming()                                                                    // 停止从队列中拉取任务，需要着重看一下不同broker是怎么实现停止消费的
	Publish(ctx context.Context, task *tasks.Signature) error                          // 发布任务到队列中
	GetPendingTasks(queue string) ([]*tasks.Signature, error)                          // 获取指定队列中的待处理任务列表
	GetDelayedTasks() ([]*tasks.Signature, error)                                      // 获取延迟任务列表， 计划在未来某个时间执行的任务
	AdjustRoutingKey(s *tasks.Signature)                                               // 决定任务被发送到哪一个队列
}

// TaskProcessor - can process a delivered task
// This will probably always be a worker instance
type TaskProcessor interface {
	Process(signature *tasks.Signature) error
	CustomQueue() string
	PreConsumeHandler() bool
}

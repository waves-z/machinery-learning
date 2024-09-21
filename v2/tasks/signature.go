package tasks

import (
	"fmt"
	"time"

	"github.com/RichardKnop/machinery/v2/utils"

	"github.com/google/uuid"
)

// Arg represents a single argument passed to invocation fo a task
type Arg struct {
	Name  string      `bson:"name"`
	Type  string      `bson:"type"`
	Value interface{} `bson:"value"`
}

// Headers represents the headers which should be used to direct the task
type Headers map[string]interface{}

// Set on Headers implements opentracing.TextMapWriter for trace propagation
func (h Headers) Set(key, val string) {
	h[key] = val
}

// ForeachKey on Headers implements opentracing.TextMapReader for trace propagation.
// It is essentially the same as the opentracing.TextMapReader implementation except
// for the added casting from interface{} to string.
func (h Headers) ForeachKey(handler func(key, val string) error) error {
	for k, v := range h {
		// Skip any non string values
		stringValue, ok := v.(string)
		if !ok {
			continue
		}

		if err := handler(k, stringValue); err != nil {
			return err
		}
	}

	return nil
}

// Signature represents a single task invocation
// 描述和封装一个任务的所有信息
type Signature struct {
	UUID           string //
	Name           string
	RoutingKey     string
	ETA            *time.Time
	GroupUUID      string
	GroupTaskCount int
	Args           []Arg
	Headers        Headers
	Priority       uint8
	Immutable      bool
	RetryCount     int
	RetryTimeout   int
	OnSuccess      []*Signature
	OnError        []*Signature
	ChordCallback  *Signature
	// MessageGroupId for Broker, e.g. SQS
	BrokerMessageGroupId string
	// ReceiptHandle of SQS Message
	SQSReceiptHandle string
	// StopTaskDeletionOnError used with sqs when we want to send failed messages to dlq,
	// and don't want machinery to delete from source queue
	StopTaskDeletionOnError bool
	// IgnoreWhenTaskNotRegistered auto removes the request when there is no handeler available
	// When this is true a task with no handler will be ignored and not placed back in the queue
	IgnoreWhenTaskNotRegistered bool
	//    Name         string                 `json:"name"`         // 任务名称
	//    UUID         string                 `json:"uuid"`         // 唯一标识符，标识任务
	//    RoutingKey   string                 `json:"routing_key"`  // 路由键，决定任务发布的队列
	//    Args         []Arg                  `json:"args"`         // 任务的参数列表
	//    KWArgs       map[string]interface{} `json:"kwargs"`       // 任务的关键字参数
	//    Headers      map[string]interface{} `json:"headers"`      // 任务的头部信息（元数据）
	//    Eta          *time.Time             `json:"eta"`          // 任务的预定执行时间（支持延迟任务）
	//    GroupUUID    string                 `json:"group_uuid"`   // 如果任务属于某个任务组，标识这个任务组
	//    ChordCallback *Signature            `json:"chord_callback"` // 在任务链中，任务完成后的回调
	//    RetryCount   int                    `json:"retry_count"`  // 当前的重试次数
	//    MaxRetries   int                    `json:"max_retries"`  // 最大重试次数
	//    Priority     uint8                  `json:"priority"`     // 任务的优先级
	//    Immutable    bool                   `json:"immutable"`    // 标记任务参数是否可变
	//    OnSuccess    []*Signature           `json:"on_success"`   // 成功时的回调任务列表
	//    OnError      []*Signature           `json:"on_error"`     // 失败时的回调任务列表
	//    OnRetry      []*Signature           `json:"on_retry"`     // 重试时的回调任务列表
	//    OnFailure    []*Signature           `json:"on_failure"`   // 任务失败时的处理回调列表
}

// NewSignature creates a new task signature
func NewSignature(name string, args []Arg) (*Signature, error) {
	signatureID := uuid.New().String()
	return &Signature{
		UUID: fmt.Sprintf("task_%v", signatureID),
		Name: name,
		Args: args,
	}, nil
}

func CopySignatures(signatures ...*Signature) []*Signature {
	var sigs = make([]*Signature, len(signatures))
	for index, signature := range signatures {
		sigs[index] = CopySignature(signature)
	}
	return sigs
}

func CopySignature(signature *Signature) *Signature {
	var sig = new(Signature)
	_ = utils.DeepCopy(sig, signature)
	return sig
}

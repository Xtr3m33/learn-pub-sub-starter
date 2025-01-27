package pubsub

type SimpleQueueType = int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

type AckType = int

const (
	AckTypeAck AckType = iota
	AckTypeNackRequeue
	AckTypeNackDiscard
)

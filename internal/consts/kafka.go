package consts

type TopicName string

func (t TopicName) String() string {
	return string(t)
}

const (
	TopicNameFetch      TopicName = "chachingFetchWorkerMain"
	TopicNameFetchRetry TopicName = "chachingFetchWorkerRetry"
)

func AllTopics() []TopicName {
	return []TopicName{
		TopicNameFetch,
		TopicNameFetchRetry,
	}
}

const (
	NumberOfPartitions int = 2
	ReplicationFactor      = 1
)

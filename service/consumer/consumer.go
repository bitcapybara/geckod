package consumer

import "github.com/bitcapybara/geckod"

type ConsumerInfoManager interface {
	Add(*AddConsumerParams) (*ConsumerInfo, error)
	Get(id uint64) (*ConsumerInfo, error)
	Del(id uint64) error
}

type AddConsumerParams struct {
	ClientId     uint64
	ConsumerName string
	TopicName    string
	SubName      string
	SubType      geckod.SubScriptionType
}

type ConsumerInfo struct {
	Id       uint64
	Name     string
	ClientId string
	Topic    string
	SubName  string
}

type Consumer interface {
}

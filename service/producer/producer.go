package producer

import "github.com/bitcapybara/geckod"

type ProducerInfoManager interface {
	GetOrCreate(client_id uint64, name string, access_mode geckod.ProducerAccessMode) (*ProducerInfo, error)
	Get(id uint64) (*ProducerInfo, error)
	Del(id uint64) (*ProducerInfo, error)
}

type ProducerInfo struct {
	Id uint64
}

type Producer interface {
}

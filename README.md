# Geckod

参考 Pulsar 实现的简单消息队列，开发中～

### 集成
* 与网络层/存储层解耦，以接口形式集成
* 当前版本只实现了消息队列的核心功能，如果需要实现完整的消息队列，需要自行实现网络/存储层

### 集成方需要实现的接口

* 发送消息给消费者

```golang
type ConsumerMessageSender interface {
	Send([]*geckod.RawMessage) error
}
```
* 消息存储接口(可选，默认内存存储)

```golang
type Storage interface {
	Add(*geckod.RawMessage) (uint64, error)
	Get(id uint64) (*geckod.RawMessage, error)
	GetRange(from uint64, to uint64) ([]*geckod.RawMessage, error)
	GetBatch(ids []uint64) ([]*geckod.RawMessage, error)
	GetMore(limit uint64) ([]*geckod.RawMessage, error)
	DelUntil(id uint64) error
}
```

### TODO
- [ ] 单元测试
- [ ] 集成使用示例
- [ ] 客户端
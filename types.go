package geckod

import (
	cmdpb "github.com/bitcapybara/geckod-proto/gen/go/proto/command"
)

type CommandConnect cmdpb.Connect
type ConnectAuthMethod cmdpb.Connect_AuthMethod
type CommandConnected cmdpb.Connected
type CommandProducer cmdpb.Producer
type ProducerAccessMode cmdpb.Producer_AccessMode
type CommandProducerSuccess cmdpb.ProducerSuccess
type CommandSubscribe cmdpb.Subscribe
type SubScriptionType cmdpb.Subscribe_SubscriptionType
type CommandSubscribeSuccess cmdpb.SubscribeSuccess
type CommandUnsubscribe cmdpb.Unsubscribe
type CommandSend cmdpb.Send
type CommandAck cmdpb.Ack
type CommandFlow cmdpb.Flow

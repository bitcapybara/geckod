package geckod

import (
	cmdpb "github.com/bitcapybara/geckod-proto/gen/go/proto/command"
	"github.com/bitcapybara/geckod/errs"
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
type AckType cmdpb.Ack_AckType
type CommandFlow cmdpb.Flow

func (c SubScriptionType) asRaw() cmdpb.Subscribe_SubscriptionType {
	return cmdpb.Subscribe_SubscriptionType(c)
}

func (c SubScriptionType) MatchAckType(ack AckType) error {
	if ack == AckType(cmdpb.Ack_Cumulative) && c == SubScriptionType(cmdpb.Subscribe_Shared) {
		return errs.ErrUnmatchAckType
	}
	return nil
}

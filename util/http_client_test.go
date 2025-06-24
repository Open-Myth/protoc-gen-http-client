package util

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/runtime/protoimpl"
)

func Test_injectRequestPathValue(t *testing.T) {
	type Channel struct {
		state         protoimpl.MessageState
		sizeCache     protoimpl.SizeCache
		unknownFields protoimpl.UnknownFields

		ChannelId   int64  `protobuf:"varint,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
		Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
		Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	}

	reqData := &Channel{
		ChannelId: 199,
	}
	req := httptest.NewRequest(http.MethodGet, "/channels/{channel_id}", nil)
	injectRequestPathValue(req, reqData)
	require.Equal(t, "/channels/199", req.URL.String())

	type Message struct {
		state         protoimpl.MessageState
		sizeCache     protoimpl.SizeCache
		unknownFields protoimpl.UnknownFields

		ChannelId int64    `protobuf:"varint,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
		MessageId int64    `protobuf:"varint,2,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
		UserId    string   `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
		Content   string   `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
		Reaction  []string `protobuf:"bytes,5,rep,name=reaction,proto3" json:"reaction,omitempty"`
	}

	reqData2 := &Message{
		ChannelId: 199,
		MessageId: 288,
	}

	req2 := httptest.NewRequest(http.MethodGet, "/channels/{channel_id}/messages/{message_id}", nil)
	injectRequestPathValue(req2, reqData2)
	require.Equal(t, "/channels/199/messages/288", req2.URL.String())
}

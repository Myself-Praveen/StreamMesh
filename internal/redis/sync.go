package redis

import (
	"context"
	"encoding/json"

	"github.com/Myself-Praveen/StreamMesh/internal/logger"
	"go.uber.org/zap"
)

const TopicSyncChannel = "streammesh:topic_sync"

// TopicSyncMessage is broadcast when a node adds/removes a local subscription
type TopicSyncMessage struct {
	NodeID string `json:"node_id"`
	Topic  string `json:"topic"`
	Action string `json:"action"` // "subscribe" or "unsubscribe"
}

// BroadcastTopicSync publishes a topic sync message to other nodes
func BroadcastTopicSync(ctx context.Context, topic, action string) {
	if Client == nil {
		return
	}

	msg := TopicSyncMessage{
		NodeID: CurrentNode.ID,
		Topic:  topic,
		Action: action,
	}

	data, _ := json.Marshal(msg)
	Client.Publish(ctx, TopicSyncChannel, data)
}

// StartTopicSyncListener listens for topic sync messages from other nodes
func StartTopicSyncListener(ctx context.Context, onSync func(topic, action string)) {
	if Client == nil {
		return
	}

	sub := Client.Subscribe(ctx, TopicSyncChannel)
	go func() {
		defer sub.Close()
		ch := sub.Channel()

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-ch:
				var syncMsg TopicSyncMessage
				if err := json.Unmarshal([]byte(msg.Payload), &syncMsg); err != nil {
					logger.Log.Error("Failed to parse topic sync message", zap.Error(err))
					continue
				}

				// Ignore our own messages
				if syncMsg.NodeID == CurrentNode.ID {
					continue
				}

				if onSync != nil {
					onSync(syncMsg.Topic, syncMsg.Action)
				}
			}
		}
	}()
}

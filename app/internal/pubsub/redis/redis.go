package pubsub_redis

import (
	"app/graph/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type RedisPubSub struct {
	client *redis.Client
	mu     sync.RWMutex
	subs   map[uuid.UUID]map[*redis.PubSub]struct{}
}

func NewRedisPubSub(client *redis.Client) *RedisPubSub {
	return &RedisPubSub{
		client: client,
		subs:   make(map[uuid.UUID]map[*redis.PubSub]struct{}),
	}
}

func (r *RedisPubSub) PublishComment(ctx context.Context, postID uuid.UUID, comment *model.Comment) error {
	payload, err := json.Marshal(comment)
	if err != nil {
		return fmt.Errorf("failed to marshal comment: %w", err)
	}

	channel := r.getChannelName(postID)
	if err := r.client.Publish(ctx, channel, payload).Err(); err != nil {
		return fmt.Errorf("failed to publish comment: %w", err)
	}

	log.Printf("Published comment to channel %s", channel)
	return nil
}

func (r *RedisPubSub) SubscribeOnComments(ctx context.Context, postID uuid.UUID) (<-chan *model.Comment, error) {
	channel := r.getChannelName(postID)
	pubsub := r.client.Subscribe(ctx, channel)

	r.mu.Lock()
	if _, exists := r.subs[postID]; !exists {
		r.subs[postID] = make(map[*redis.PubSub]struct{})
	}
	r.subs[postID][pubsub] = struct{}{}
	r.mu.Unlock()

	if _, err := pubsub.Receive(ctx); err != nil {
		r.removeSubscription(postID, pubsub)
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	commentChan := make(chan *model.Comment)

	go r.listenForMessages(ctx, postID, pubsub, commentChan)

	return commentChan, nil
}

func (r *RedisPubSub) listenForMessages(ctx context.Context, postID uuid.UUID, pubsub *redis.PubSub, commentChan chan<- *model.Comment) {
	defer func() {
		r.removeSubscription(postID, pubsub)
		close(commentChan)
	}()

	redisChan := pubsub.Channel(redis.WithChannelSize(100))

	for {
		select {
		case <-ctx.Done():
			log.Printf("Subscription canceled for post %s", postID)
			return
		case msg, ok := <-redisChan:
			if !ok {
				log.Printf("Redis channel closed for post %s", postID)
				return
			}

			var comment model.Comment
			if err := json.Unmarshal([]byte(msg.Payload), &comment); err != nil {
				log.Printf("Failed to unmarshal comment: %v", err)
				continue
			}

			select {
			case commentChan <- &comment:
			case <-time.After(5 * time.Second):
				log.Println("Timeout sending comment to channel")
			case <-ctx.Done():
				return
			}
		}
	}
}

func (r *RedisPubSub) removeSubscription(postId uuid.UUID, pubsub *redis.PubSub) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if subs, exists := r.subs[postId]; exists {
		delete(subs, pubsub)
		if len(subs) == 0 {
			delete(r.subs, postId)
		}
	}
	pubsub.Close()
}

func (r *RedisPubSub) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var err error
	for _, subs := range r.subs {
		for pubsub := range subs {
			if e := pubsub.Close(); e != nil && err == nil {
				err = e
			}
		}
	}
	r.subs = make(map[uuid.UUID]map[*redis.PubSub]struct{})

	if e := r.client.Close(); e != nil && err == nil {
		err = e
	}

	return err
}

func (r *RedisPubSub) getChannelName(postID uuid.UUID) string {
	return fmt.Sprintf("comments:%s", postID.String())
}

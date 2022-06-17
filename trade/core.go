package trade

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/longbridgeapp/openapi-go/config"
	"github.com/longbridgeapp/openapi-go/log"
	"github.com/longbridgeapp/openapi-protobufs/gen/go/trade"
	protocol "github.com/longbridgeapp/openapi-protocol/go"
	"github.com/longbridgeapp/openapi-protocol/go/client"
	"github.com/pkg/errors"
)

type Core struct {
	client        *client.Client
	url           string
	subscriptions []string
	mu            sync.Mutex
}

func NewCore(url string, otp string) (*Core, error) {
	cl := client.New()
	err := cl.Dial(context.Background(), url, &protocol.Handshake{
		Version:  1,
		Codec:    protocol.CodecProtobuf,
		Platform: protocol.PlatformOpenapi,
	}, client.WithAuthToken(otp))
	if err != nil {
		return nil, err
	}
	cl.Logger.SetLevel(config.GetLogLevelFromEnv())
	core := &Core{client: cl, url: url}
	return core, nil
}

func (c *Core) SetHandler(f func(*PushEvent)) {
	c.client.AfterReconnected(func() {
		if err := c.resubscribe(context.Background()); err != nil {
			log.Error(err)
		}
	})
	c.client.Subscribe(uint32(tradev1.Command_CMD_NOTIFY), parseNotifyFunc(f))
}

func (c *Core) Subscribe(ctx context.Context, topics []string) (subRes *SubResponse, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.doSubscribe(ctx, topics)
}

func (c *Core) doSubscribe(ctx context.Context, topics []string) (subRes *SubResponse, err error) {
	var res *protocol.Packet
	req := &tradev1.Sub{Topics: topics}
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(tradev1.Command_CMD_SUB), Body: req})
	if err != nil {
		return
	}
	var tradeRes tradev1.SubResponse
	if err = res.Unmarshal(&tradeRes); err != nil {
		return
	}
	subRes = &SubResponse{}
	subRes.Current = tradeRes.GetCurrent()
	subRes.Success = tradeRes.GetSuccess()
	subRes.Fail = make([]*SubResponseFail, len(tradeRes.GetFail()))
	for _, f := range tradeRes.GetFail() {
		subRes.Fail = append(subRes.Fail, &SubResponseFail{Topic: f.GetTopic(), Reason: f.GetReason()})
	}
	c.subscriptions = subRes.Current
	return
}

func (c *Core) Unsubscribe(ctx context.Context, topics []string) (unsubRes *UnsubResponse, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var res *protocol.Packet
	req := &tradev1.Unsub{Topics: topics}
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(tradev1.Command_CMD_UNSUB), Body: req})
	if err != nil {
		return
	}
	var tradeRes tradev1.UnsubResponse
	if err = res.Unmarshal(&tradeRes); err != nil {
		return
	}
	unsubRes = &UnsubResponse{}
	unsubRes.Current = tradeRes.GetCurrent()
	c.subscriptions = tradeRes.GetCurrent()
	return
}

func (c *Core) resubscribe(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	res, err := c.doSubscribe(ctx, c.subscriptions)
	if err != nil {
		return errors.Wrap(err, "resubscribe error")
	}
	if len(res.Fail) > 0 {
		log.Errorf("resubscirbe subscription some failed %v", res.Fail)
	}
	return nil
}

func (c *Core) Close() error {
	return c.client.Close(nil)
}

func parseNotifyFunc(f func(*PushEvent)) func(*protocol.Packet) {
	return func(packet *protocol.Packet) {
		var notify tradev1.Notification
		if err := packet.Unmarshal(&notify); err != nil {
			log.Error(err)
			return
		}
		var event PushEvent
		if err := json.Unmarshal(notify.GetData(), &event); err != nil {
			log.Error(err)
		}
		f(&event)
	}
}

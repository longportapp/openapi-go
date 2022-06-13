package trade

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/longbridgeapp/openapi-protobufs/gen/go/trade"
	protocol "github.com/longbridgeapp/openapi-protocol/go"
	"github.com/longbridgeapp/openapi-protocol/go/client"
	"google.golang.org/protobuf/proto"
)

type Core struct {
	client *client.Client
	url    string
	subscriptions []string
	mu sync.Mutex
}

func NewCore(url string, otp string, f func(*PushEvent)) (*Core, error) {
	cl := client.New()
	err := cl.Dial(context.Background(), url, &protocol.Handshake{
		Version:  1,
		Codec:    protocol.CodecProtobuf,
		Platform: protocol.PlatformOpenapi,
	}, client.WithAuthToken(otp))
	if err != nil {
		return nil, err
	}
	core := &Core{client: cl, url: url}
	cl.AfterReconnected(func (){
		if err := core.resubscribe(context.Background()); err != nil {
			log.Error(err)
		}
	})
	cl.Subscribe(uint32(tradev1.Command_CMD_NOTIFY), parseNotifyFunc(f))
	return core, nil
}

func (c *Core) Subscribe(ctx context.Context, symbols []string) (subRes *SubResponse, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.doSubscirbe(ctx, symbols)
	return
}

func (c *Core) doSubscirbe(ctx context.Context, symbols []string) (subRes *SubResponse, err error) {
	var res *protocol.Packet
	req := &tradev1.Sub{Topics: symbols}
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

func (c *Core) Unsubscribe(ctx context.Context, symbols []string) (unsubRes *UnsubResponse, err error){
	c.mu.Lock()
	defer c.mu.Unlock()
	var res *protocol.Packet
	req := &tradev1.Unsub{Topics: symbols}
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(tradev1.Command_CMD_UNSUB), Body: req })
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
	res, err := c.doSubscirbe(ctx, c.subscriptions)
	if err != nil {
		return err
	}
	if len(res.Fail) >0 {
		log.Errorf(res.Fail)
	}
	return nil
}

func parseNotifyFunc(f func(*PushEvent)) func(*protocol.Packet) {
	return func(packet *protocol.Packet) {
		var notify tradev1.Notification
		if err := packet.Unmarshal(&notify); err != nil {
			log.Errorf(err)
			return
		}
		if notify.GetContentType() != 1 {
			log.Error("trade context event content type not json")
			return
		}
		var orderChange PushOrderChanged
	    if err := json.Unmarshal(notify.GetData(), orderChange); err != nil {
			log.Error(err)
		}
		var event PushEvent
		event.orderChanged = &orderChange
		f(&event)
	}
}

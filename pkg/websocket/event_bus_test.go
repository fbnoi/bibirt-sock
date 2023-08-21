package websocket

import (
	"bibirt-sock/pkg/websocket/pb"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestNew(t *testing.T) {
	bus := NewBus()
	if bus == nil {
		t.Log("New EventBus not created!")
		t.Fail()
	}
}

func TestHasCallback(t *testing.T) {
	bus := NewBus()
	bus.Subscribe(&pb.Heartbeat{}, func(proto.Message) {})
	assert.Equal(t, true, bus.HasCallback(&pb.Heartbeat{}))
	assert.Equal(t, false, bus.HasCallback(&pb.Close{}))
}

func TestSubscribe(t *testing.T) {
	bus := NewBus()
	bus.Subscribe(&pb.Heartbeat{}, func(proto.Message) {})
	assert.Equal(t, true, bus.HasCallback(&pb.Heartbeat{}))
}

func TestSubscribeOnce(t *testing.T) {
	bus := NewBus()
	bus.SubscribeOnce(&pb.Heartbeat{}, func(proto.Message) {})
	assert.Equal(t, true, bus.HasCallback(&pb.Heartbeat{}))
}

func TestSubscribeOnceAndManySubscribe(t *testing.T) {
	bus := NewBus()
	event := &pb.Heartbeat{}
	flag := 0
	fn := func(proto.Message) {
		flag += 1
	}
	bus.SubscribeOnce(event, fn)
	bus.Subscribe(event, fn)
	bus.Subscribe(event, fn)
	bus.Publish(event)
	assert.Equal(t, 3, flag)
}

func TestUnsubscribe(t *testing.T) {
	bus := NewBus()
	handler := func(proto.Message) {}
	bus.Subscribe(&pb.Heartbeat{}, handler)
	assert.Nil(t, bus.Unsubscribe(&pb.Heartbeat{}, handler))
	assert.NotNil(t, bus.Unsubscribe(&pb.Heartbeat{}, handler))
}

type handler struct {
	val int
}

func (h *handler) Handle(proto.Message) {
	h.val++
}

func TestUnsubscribeMethod(t *testing.T) {
	bus := NewBus()
	h := &handler{val: 0}
	event := &pb.Heartbeat{}
	bus.Subscribe(event, h.Handle)
	bus.Publish(event)
	assert.Nil(t, bus.Unsubscribe(event, h.Handle))
	assert.NotNil(t, bus.Unsubscribe(event, h.Handle))
	bus.Publish(event)
	bus.WaitAsync()
	assert.Equal(t, 1, h.val)
}

var uniqueBus = NewBus()
var once = sync.Once{}

func doOnce() {
	once.Do(func() {
		uniqueBus.SubscribeAsync(&pb.Heartbeat{}, func(proto.Message) {
		}, false)
	})
}

func BenchmarkSubscribeAsync(b *testing.B) {
	doOnce()
	for i := 0; i < 1000; i++ {
		uniqueBus.Publish(&pb.Heartbeat{})
	}
	uniqueBus.WaitAsync()
}

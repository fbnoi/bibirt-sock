package websocket

import (
	"reflect"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type HandleFunc func(proto.Message)

type BusSubscriber interface {
	Subscribe(typ proto.Message, fn HandleFunc)
	SubscribeAsync(typ proto.Message, fn HandleFunc, transactional bool)
	SubscribeOnce(typ proto.Message, fn HandleFunc)
	SubscribeOnceAsync(typ proto.Message, fn HandleFunc)
	Unsubscribe(typ proto.Message, fn HandleFunc) error
}

type BusPublisher interface {
	Publish(proto.Message)
}

type BusController interface {
	HasCallback(typ proto.Message) bool
	WaitAsync()
}

type Bus interface {
	BusController
	BusSubscriber
	BusPublisher
}

func NewBus() Bus {
	b := &EventBus{
		make(map[string][]*eventHandler),
		sync.RWMutex{},
		sync.WaitGroup{},
	}
	return b
}

type EventBus struct {
	handlers map[string][]*eventHandler
	lock     sync.RWMutex
	wg       sync.WaitGroup
}

type eventHandler struct {
	handleFunc    HandleFunc
	flagOnce      bool
	async         bool
	transactional bool
	sync.Mutex
}

func (bus *EventBus) doSubscribe(m proto.Message, handler *eventHandler) {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	typ := messageUrlType(m)

	bus.handlers[typ] = append(bus.handlers[typ], handler)
}

func (bus *EventBus) Subscribe(m proto.Message, fn HandleFunc) {
	bus.doSubscribe(m, &eventHandler{fn, false, false, false, sync.Mutex{}})
}

func (bus *EventBus) SubscribeAsync(m proto.Message, fn HandleFunc, transactional bool) {
	bus.doSubscribe(m, &eventHandler{fn, false, true, transactional, sync.Mutex{}})
}

func (bus *EventBus) SubscribeOnce(m proto.Message, fn HandleFunc) {
	bus.doSubscribe(m, &eventHandler{fn, true, false, false, sync.Mutex{}})
}

func (bus *EventBus) SubscribeOnceAsync(m proto.Message, fn HandleFunc) {
	bus.doSubscribe(m, &eventHandler{fn, true, true, false, sync.Mutex{}})
}

func (bus *EventBus) Unsubscribe(m proto.Message, fn HandleFunc) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	typ := messageUrlType(m)

	if _, ok := bus.handlers[typ]; ok && len(bus.handlers[typ]) > 0 {
		bus.removeHandler(typ, bus.findHandlerIdx(typ, fn))
		return nil
	}
	return errors.Errorf("EventBus.Unsubscribe: message %v doesn't exist", typ)
}

func (bus *EventBus) Publish(m proto.Message) {
	bus.lock.RLock()
	defer bus.lock.RUnlock()
	typ := messageUrlType(m)

	if handlers, ok := bus.handlers[typ]; ok && 0 < len(handlers) {
		onceIdx := []int{}
		for i, handler := range handlers {
			if handler.flagOnce {
				onceIdx = append(onceIdx, i)
			}
			if !handler.async {
				bus.doPublish(handler, m)
			} else {
				bus.wg.Add(1)
				if handler.transactional {
					bus.lock.Unlock()
					handler.Lock()
					bus.lock.Lock()
				}
				go bus.doPublishAsync(handler, m)
			}
		}
		for _, i := range onceIdx {
			bus.removeHandler(typ, i)
		}
	}
}

func (bus *EventBus) HasCallback(m proto.Message) bool {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	typ := messageUrlType(m)

	_, ok := bus.handlers[typ]
	if ok {
		return len(bus.handlers[typ]) > 0
	}
	return false
}

func (bus *EventBus) WaitAsync() {
	bus.wg.Wait()
}

func (bus *EventBus) doPublish(handler *eventHandler, m proto.Message) {
	handler.handleFunc(m)
}

func (bus *EventBus) doPublishAsync(handler *eventHandler, m proto.Message) {
	defer bus.wg.Done()
	if handler.transactional {
		defer handler.Unlock()
	}
	bus.doPublish(handler, m)
}

func (bus *EventBus) removeHandler(typ string, i int) {
	if _, ok := bus.handlers[typ]; !ok {
		return
	}
	l := len(bus.handlers[typ])

	if !(0 <= i && i < l) {
		return
	}

	copy(bus.handlers[typ][i:], bus.handlers[typ][i+1:])
	bus.handlers[typ][l-1] = nil
	bus.handlers[typ] = bus.handlers[typ][:l-1]
}

func (bus *EventBus) findHandlerIdx(typ string, fn HandleFunc) int {
	if _, ok := bus.handlers[typ]; ok {
		for i, v := range bus.handlers[typ] {
			sf1 := reflect.ValueOf(v.handleFunc)
			sf2 := reflect.ValueOf(fn)
			if sf1.Type() == sf2.Type() && sf1.Pointer() == sf2.Pointer() {
				return i
			}
		}
	}
	return -1
}

func messageUrlType(m proto.Message) string {
	return string(m.ProtoReflect().Descriptor().FullName())
}

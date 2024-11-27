package events

import "sync"

type EventType string

const (
	EventGameCreated EventType = "game_created"
	EventGameStarted EventType = "game_started"
	EventGameEnded   EventType = "game_ended"
	EventTimerUpdate EventType = "timer_update"
	EventTimeout     EventType = "timeout"
)

type Event struct {
	Type    EventType
	Payload map[string]string
}

type EventHandler func(event Event)

type EventBus struct {
	handlers map[EventType][]EventHandler
	mu       sync.Mutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]EventHandler),
	}
}

func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

func (eb *EventBus) Publish(event Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if handlers, exists := eb.handlers[event.Type]; exists {
		for _, handler := range handlers {
			go handler(event)
		}
	}
}

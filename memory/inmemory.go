package memory

import "github.com/my-scratch-agent/domain"

type InMemory struct {
	history []domain.Message
}

func NewInMemory() *InMemory {
	return &InMemory{}
}

func (m *InMemory) Add(msg domain.Message) {
	m.history = append(m.history, msg)
}

func (m *InMemory) GetHistory() []domain.Message {
	return m.history
}

func (m *InMemory) Clear() {
	m.history = nil
}

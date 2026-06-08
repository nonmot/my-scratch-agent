package domain

type Memory interface {
	Add(msg Message)
	GetHistory() []Message
	Clear()
}

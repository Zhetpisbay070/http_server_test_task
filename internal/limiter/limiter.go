package limiter

import "sync"

// Limiter управляет ограничением запросов
type Limiter struct {
	incoming chan struct{} // Ограничение входящих запросов (100)
	outgoing chan struct{} // Ограничение исходящих запросов (4)
	once     sync.Once
}

// NewLimiter создает новый лимитер
func NewLimiter(maxIncoming, maxOutgoing int) *Limiter {
	return &Limiter{
		incoming: make(chan struct{}, maxIncoming),
		outgoing: make(chan struct{}, maxOutgoing),
	}
}

// AcquireIncoming блокирует входящий запрос, если лимит исчерпан
func (l *Limiter) AcquireIncoming() bool {
	select {
	case l.incoming <- struct{}{}:
		return true
	default:
		return false
	}
}

// ReleaseIncoming освобождает слот для входящего запроса
func (l *Limiter) ReleaseIncoming() {
	<-l.incoming
}

// AcquireOutgoing блокирует исходящий запрос, если лимит исчерпан
func (l *Limiter) AcquireOutgoing() {
	l.outgoing <- struct{}{}
}

// ReleaseOutgoing освобождает слот для исходящего запроса
func (l *Limiter) ReleaseOutgoing() {
	<-l.outgoing
}

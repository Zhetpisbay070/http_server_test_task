package limiter

import (
	"testing"
	"time"
)

func TestLimiter_AcquireIncoming(t *testing.T) {
	l := NewLimiter(2, 1) // Входящих: 2, исходящих: 1

	if !l.AcquireIncoming() {
		t.Errorf("Первый входящий запрос не должен быть заблокирован")
	}
	if !l.AcquireIncoming() {
		t.Errorf("Второй входящий запрос не должен быть заблокирован")
	}

	// Третий должен быть заблокирован
	if l.AcquireIncoming() {
		t.Errorf("Третий входящий запрос должен быть заблокирован")
	}

	// Освобождаем слот
	l.ReleaseIncoming()
	if !l.AcquireIncoming() {
		t.Errorf("После освобождения слот должен быть доступен")
	}
}

func TestLimiter_AcquireOutgoing(t *testing.T) {
	l := NewLimiter(2, 1) // Входящих: 2, исходящих: 1

	// Первый поток должен успешно занять слот
	done := make(chan struct{})
	go func() {
		l.AcquireOutgoing()
		close(done)
	}()

	select {
	case <-done:
		// Все нормально, слот занят
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Исходящий запрос не должен быть заблокирован")
	}

	// Второй поток должен заблокироваться
	blocked := make(chan struct{})
	go func() {
		l.AcquireOutgoing()
		close(blocked)
	}()

	select {
	case <-blocked:
		t.Errorf("Второй исходящий запрос должен быть заблокирован")
	case <-time.After(100 * time.Millisecond):
		// Все верно, поток заблокирован
	}

	// Освобождаем слот
	l.ReleaseOutgoing()

	// Теперь заблокированный поток должен пройти
	select {
	case <-blocked:
		// Все нормально
	case <-time.After(100 * time.Millisecond):
		t.Errorf("После освобождения слот должен быть доступен")
	}
}

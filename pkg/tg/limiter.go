package tg

import (
	"sync"
	"time"
)

const (
	GlobalLimit    = 30          // Максимум запросов в секунду к API
	ChatLimit      = 1           // Максимум сообщений в секунду в один чат
	MultiChatLimit = 30          // Максимум сообщений в секунду в разные чаты
	WaitTime       = time.Second // Ждем полную секунду при превышении лимита
)

type Limiter struct {
	mu sync.Mutex

	ChatTimes    map[int64]time.Time
	MessageTimes []time.Time
	ApiTimes     []time.Time
}

func NewLimiter() *Limiter {
	return &Limiter{
		ChatTimes:    make(map[int64]time.Time),
		MessageTimes: make([]time.Time, 0, MultiChatLimit),
		ApiTimes:     make([]time.Time, 0, GlobalLimit),
	}
}

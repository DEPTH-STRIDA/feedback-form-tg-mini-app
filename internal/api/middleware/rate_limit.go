package middleware

import (
	"net/http"
	"nstu/internal/logger"
	"sync"
	"time"
)

// Структура для хранения информации о запросах с IP
type ipLimit struct {
	count    int       // количество запросов
	lastSeen time.Time // время последнего запроса
}

type RateLimiter struct {
	mu      sync.RWMutex
	limits  map[string]*ipLimit
	rate    int           // максимальное количество запросов
	window  time.Duration // временное окно
	cleanup time.Duration // интервал очистки старых записей
}

func NewRateLimiter(rate int, window, cleanup time.Duration) *RateLimiter {
	limiter := &RateLimiter{
		limits:  make(map[string]*ipLimit),
		rate:    rate,
		window:  window,
		cleanup: cleanup,
	}

	// Запускаем очистку старых записей
	go limiter.cleanupLoop()
	return limiter
}

// Очистка старых записей
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanup)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, limit := range rl.limits {
			if now.Sub(limit.lastSeen) > rl.window {
				delete(rl.limits, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		rl.mu.Lock()
		limit, exists := rl.limits[ip]
		now := time.Now()

		if !exists {
			// Первый запрос с этого IP
			rl.limits[ip] = &ipLimit{
				count:    1,
				lastSeen: now,
			}
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		// Проверяем, не истекло ли окно
		if now.Sub(limit.lastSeen) > rl.window {
			// Сбрасываем счетчик
			limit.count = 1
			limit.lastSeen = now
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		// Увеличиваем счетчик
		limit.count++
		limit.lastSeen = now

		// Проверяем превышение лимита
		if limit.count > rl.rate {
			rl.mu.Unlock()
			logger.Log.Warn().
				Str("ip", ip).
				Int("count", limit.count).
				Msg("Превышен лимит запросов")
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		rl.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// IPRateLimiter ограничивает количество запросов от одного IP-адреса
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit // лимит событий в секунду
	b   int        // размер буфера (burst size)
}

// NewIPRateLimiter создает новый IPRateLimiter с заданными параметрами
// r - лимит событий в секунду (например, rate.Limit(5) = 5 событий в секунду)
// b - размер буфера для всплесков нагрузки (burst size)
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// AddIP добавляет новый IP-адрес в лимит
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter

	return limiter
}

// GetLimiter возвращает лимит для заданного IP-адреса
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	limiter, exists := i.ips[ip]

	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}

	i.mu.Unlock()
	return limiter
}

// RateLimitMiddleware создает middleware с заданными параметрами
func RateLimitMiddleware(r rate.Limit, b int) func(http.Handler) http.Handler {
	limiter := NewIPRateLimiter(r, b)

	// Возвращаем middleware
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// получаем лимит для IP-адреса
			l := limiter.GetLimiter(r.RemoteAddr)
			// если лимит превышен, возвращаем ошибку 429 Too Many Requests
			if !l.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			// если лимит не превышен, передаем запрос следующему обработчику
			next.ServeHTTP(w, r)
		})
	}
}

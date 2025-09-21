package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type limiterEntry struct {
	lim  *rate.Limiter
	last time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*limiterEntry
	r        rate.Limit
	b        int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{visitors: make(map[string]*limiterEntry), r: r, b: b}
}

func (rl *RateLimiter) get(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	v, ok := rl.visitors[ip]
	if !ok {
		l := rate.NewLimiter(rl.r, rl.b)
		rl.visitors[ip] = &limiterEntry{lim: l, last: time.Now()}
		return l
	}
	v.last = time.Now()
	return v.lim
}

func (rl *RateLimiter) Cleanup(expire time.Duration) {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, e := range rl.visitors {
			if time.Since(e.last) > expire {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		lim := rl.get(ip)
		if !lim.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}

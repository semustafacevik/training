package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// DONE: Buradaki mapi thread safe hale getirebilirsiniz.
// 102-concurrency egitimindeki mutex orneklerine bakabilirsiniz.
// Ref: https://pmihaylov.com/thread-safety-concerns-go/
// Ref: https://medium.com/@deckarep/the-new-kid-in-town-gos-sync-map-de24a6bf7c2c
var mutex sync.Mutex
var counter = map[string]*Limit{}

type Limit struct {
	count int
	ttl   time.Time
}

type LimitProxy struct {
	key   string
	limit int
	ttl   time.Duration
}

func ResetLimitHandler(c *fiber.Ctx) error {
	// TODO: [DELETE] /limit/:key/* pathine istek atildiginda limiti sifirlayan handleri implement edebilirsiniz.
	// TODO: implement me!
	return nil
}

func NewLimitProxy(key string, limit int, ttl time.Duration) LimitProxy {
	return LimitProxy{
		key:   key,
		limit: limit,
		ttl:   ttl,
	}
}

func (p LimitProxy) Accept(key string) bool {
	return p.key == key
}

func (p LimitProxy) Proxy(c *fiber.Ctx) error {
	path := c.Path()

	if r, ok := counter[path]; ok && r.count >= p.limit {
		if r.ttl.After(time.Now()) {
			c.Response().SetStatusCode(fiber.StatusTooManyRequests)

			fmt.Printf("request limit reached for %s \n", path)

			return fiber.ErrTooManyRequests
		} else {
			fmt.Printf("resetting counter values \n")

			//counter[path] = &Limit{
			//	count: 0,
			//	ttl:   time.Now().Add(p.ttl),
			//}

			// thread safe version
			defineCounter(path, p.ttl)
		}
	} else if !ok {
		//counter[path] = &Limit{
		//	count: 0,
		//	ttl:   time.Now().Add(p.ttl),
		//}

		// thread safe version
		defineCounter(path, p.ttl)
	}

	if err := c.SendString("Go Turkiye - 103 Http Package"); err != nil {
		return err
	}

	//counter[path].count++

	// thread safe version
	incrementCounter(path)

	return nil
}

func defineCounter(path string, ttl time.Duration) {
	mutex.Lock()
	counter[path] = &Limit{count: 0, ttl: time.Now().Add(ttl)}
	mutex.Unlock()
}

func incrementCounter(path string) {
	mutex.Lock()
	counter[path].count++
	mutex.Unlock()
}

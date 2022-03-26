/*
@Time    : 3/26/22 09:51
@Author  : Neil
@File    : stat.go
*/

package statistics

import (
	"log"
	"sync"
)

var BytesAccount *Stat

func init() {
	BytesAccount = &Stat{cache: make(map[string]int64)}
}

type Stat struct {
	mu    sync.Mutex
	cache map[string]int64
}

func (stat *Stat) Add(key string, delta int64) {
	stat.mu.Lock()
	defer stat.mu.Unlock()
	stat.cache[key] += delta
	log.Printf("cache: %v", stat.cache)
	return
}

func (stat *Stat) save() {

}

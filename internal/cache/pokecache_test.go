package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPokeCacheConcurrent(t *testing.T) {
    c := NewPokeCache(time.Second)
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            key := fmt.Sprintf("k%d", n)
            data := []byte{byte(n)}
            c.Add(key, data)
            if v, ok := c.Get(key); !ok {
                t.Errorf("expected value after add for key %s", key)
            } else if len(v) != 1 {
                t.Errorf("unexpected length: %d", len(v))
            }
        }(i)
    }
    wg.Wait()
}

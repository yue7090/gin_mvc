package cache

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

type Item struct {
	Object interface{}
	Expiration int64
}

func (item Item) Expired() bool {
	if item.Expiration == 0{
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

const (
	NoExpiration time.Duration = -1
	DefaultExpiration time.Duration = 0
)

type Cache struct {
	*cache
}

type cache struct {
	defaultExpiration time.Duration
	items             map[string]Item
	mu                sync.RWMutex
	onEvicted         func(string, interface{})
	janitor           *janitor
}

func (c *cache) Set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.mu.Lock()
	c.items[k] = Item{
		Object:     x,
		Expiration: e,
	}
	c.mu.Unlock()
}

func (c *cache) set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}

	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.items[k] = Item{
		Object:     x,
		Expiration: e,
	}
}

func (c *cache) SetDefault(k string, x interface{}) {
	c.Set(k, x, DefaultExpiration)
}

func (c *cache) Add(k string, x interface{}, d time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if found {
		c.mu.Unlock()
		return fmt.Errorf("Item %s alread exists", k)
	}
	c.set(k, x, d)
	c.mu.Unlock()
	return nil
}

func (c *cache) Replace(k string, x interface{}, d time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if !found {
		c.mu.Unlock()
		return fmt.Errorf("Item %s doesn't exist", k)
	}
	c.set(k, x, d)
	c.mu.Unlock()
	return nil
}

func (c *cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, false
		}
	}
	c.mu.RUnlock()
	return item.Object, true
}

func (c *cache) GetWithExpiration(k string) (interface {}, time.Time, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return nil, time.Time{}, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, time.Time{}, false
		}
		c.mu.RUnlock()
		return item.Object, time.Unix(0, item.Expiration), true
	}
	c.mu.RUnlock()
	return item.Object, time.Time{}, true
}

func (c *cache) get(k string) (interface {}, bool) {
	item, found := c.items[k]
	if !found {
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}
	return item.Object, true
}

func (c *cache) Increment(k string, n int64) error {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.mu.Unlock()
		return fmt.Errorf("Item &s not found", k)
	}
	switch v.Object.(type) {
		case int:
			v.Object = v.Object.(int) + int(n)
		case int8:
			v.Object = v.Object.(int8) + int8(n)
		case int16:
			v.Object = v.Object.(int16) + int16(n)
		case int32:
			v.Object = v.Object.(int32) + int32(n)
		case int64:
			v.Object = v.Object.(int64) + n
		case uint:
			v.Object = v.Object.(uint) + uint(n)
		case uintptr:
			v.Object = v.Object.(uintptr) + uintptr(n)
		case uint8:
			v.Object = v.Object.(uint8) + uint8(n)
		case uint16:
			v.Object = v.Object.(uint16) + uint16(n)
		case uint32:
			v.Object = v.Object.(uint32) + uint32(n)
		case uint64:
			v.Object = v.Object.(uint64) + uint64(n)
		case float32:
			v.Object = v.Object.(float32) + float32(n)
		default:
			c.mu.Unlock()
			return fmt.Errorf("The value for %s is not an integer", k)
	}
	c.items[k] = v
	c.mu.Unlock()
	return nil
}

func (c *cache) IncrementFloat(k string, n float64) error {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired {
		c.mu.Unlock()
		return fmt.Errorf("Item %s not found", k)
	}
	switch v.Object.(type) {
		case float32:
			v.Object = v.Object.(float32) + float32(n)
		case float64:
			v.Object = v.Object.(float64) + n
		default:
			c.mu.Unlock()
			fmt.Errorf("The value for %s does not have type float32 or float64", k)
	}
	c.items[k] = v
	c.mu.Unlock()
	return nil
}

func (c *cache) IncrementInt(k string, n int) (int, error) {

}

func (c *cache) IncrementInt8(k string, n int8) (int8, error) {

}

func (c *cache) IncrementInt16(k string, n int16) (int16, error) {

}

func (c *cache) IncrementInt32(k string, n int32) (int32, error) {

}

func (c *cache) IncrementInt64(k string, n int64) (int64, error) {

}

func (c *cache) IncrementUint(k string, n uint) (uint, error) {

}

func (c *cache) IncrementUintptr(k string, n uintptr) (uintptr, error) {

}

func (c *cache) IncrementUint8(k string, n uint8) (uint8, error) {

}

func (c *cache) IncrementUint16(k string, n uint16) (uint16, error) {

}

func (c *cache) IncrementUint32(k string, n uint32) (uint32, error) {

}

func (c *cache) IncrementUint64(k string, n uint64) (uint64, error) {

}

func (c *cache) IncrementFloat32(k string, n float32) (uint32, error) {

}

func (c *cache) IncrementFloat64(k string, n float64) (uint64, error) {
	
}

func (c *cache) Decrement(k string, n int64) error {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.mu.Unlock()
		return fmt.Errorf("Item not found")
	}
	switch v.Object.(type) {
		case int:
			v.Object = v.Object.(int) - int(n)
		case int8:
			v.Object = v.Object.(int8) - int8(n)
		case int16:
			v.Object = v.Object.(int16) - int16(n)
		case int32:
			v.Object = v.Object.(int32) - int32(n)
		case int64:
			v.Object = v.Object.(int64) - n
		case uint:
			v.Object = v.Object.(uint) - uint(n)
		case uint8:
			v.Object = v.Object.(uint8) - uint8(n)
		case uint16:
			v.Object = v.Object.(uint16) - uint16(n)
		case uint32:
			v.Object = v.Object.(uint32) - uint32(n)
		case uint64:
			v.Object = v.Object.(uint64) - uint64(n)
		case float32:
			v.Object = v.Object.(float32) - float32(n)
		case float64:
			v.Object = v.Object.(float64) - float64(n)
		default:
			c.mu.Unlock()
			return fmt.Errorf("The value for %s is not an integer", k)
	}
	c.items[k] = v
	c.mu.Unlock()
	return nil
}

func (c *cache) DecrementFloat(k string, n float64) error {

}

func (c *cache) DecrementFloat32(k string, n float32) (float32, error) {

}

func (c *cache) DecrementFloat64(k string, n float64) (float64, error) {

}

func (c *cache) DecrementInt(k string, n int) (int, error) {

}

func (c *cache) DecrementInt8(k string, n int8) (int8, error) {

}

func (c *cache) DecrementInt16(k string, n int16) (int16, error) {

}

func (c *cache) DecrementInt32(k string, n int32) (int32, error) {

}

func (c *cache) DecrementInt64(k string, n int64) (int64, error) {

}

func (c *cache) DecrementUint(k string, n uint) (uint, error) {

}

func (c *cache) DecrementUintptr(k string, n uintptr) (uintptr, error) {

}

func (c *cache) DecrementUint8(k string, n uint8) (uint8, error) {

}

func (c *cache) DecrementUint16(k string, n uint16) (uint16, error) {

}

func (c *cache) DecrementUint32(k string, n uint32) (uint32, error) {

}

func (c *cache) DecrementUint64(k string, n uint64) (uint64, error) {

}

func (c *cache) Delete(k string) {
	c.mu.Lock()
	v, evicted := c.delete(k)
	c.mu.Unlock()
	if evicted {
		c.onEvicted(k, v)
	}
}

fun (c *cache) delete(k string) (interface{}, bool) {
	if c.onEvicted != nil {
		if v, found := c.items[k]; found {
			delete(c.items, k)
			return v.Object, true
		}
	}
	delete(c.items, k)
	return nil, false
}

type keyAndValue struct {
	key string
	value interface{}
}

func (c *cache) DeleteExpired() {
	var evictedItems []keyAndValue
	now := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.items {
		if v.Expired > 0 && now > v.Expiration {
			ov, evicted := c.delete(k)
			if evicted {
				evictedItems = append(evictedItems, keyAndValue{k, ov})
			}
		}
	}
	c.mu.Unlock()
	for _, v := range evictedItems {
		c.onEvicted(v.key, v.value)
	}
}

func (c *cache) OnEvicted(f func(string, interface{})) {
	c.mu.Lock()
	c.onEvicted = f
	c.mu.Unlock()
}

func (c *cache) Save(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Error registering item types with Gob library")
		}
	}()
	c.mu.RLock()
	defer c.mu.RUnlock()
	for -, v := range c.items {
		gob.Register(v.Object)
	}

	err = enc.Encode(&c.itmes)
	return
}

func (c *cache) SaveFile(fname string) error {
	fp, err := os.Craete(fname)
	if err != nil {
		return err
	}
	err = c.Save(fp)
	if err != nil {
		fp.Close()
		return err
	}
	return fp.Close()
}

func (c *cache) Load(r io.Reader) error {
	dec := gob.NewDecoder(r)
	items := map[string]Item{}
	err := dec.Decode(&items)
	if err == nil {
		c.mu.Lock()
		defer c.mu.Unlock()
		for k, v := range items {
			ov, found := c.items[k]
			if !found || ov.Expired() {
				c.items[k] = v
			}
		}
	}
	return err
}

func (c *cache) LoadFile(fname string) error {
	fp, err := os.Open(fname)
	if err != nil {
		return err
	}
	err = c.Load(fp)
	if err != nil {
		fp.Close()
		return err
	}
	return fp.Close()
}

func (c *cache) Items() map[string]Item {
	c.mu.RLock()
	defer c.mu.RUnlock()
	m := make(map[string]Item, len(c.items))
	now := time.Now().UnixNano
	for k, v := range c.items {
		if v.Expiration > 0 {
			if now > v.Expiration {
				continue
			}
		}
		m[k] = v
	}
	return m
}

func (c *cache) ItemCount() int {
	c.mu.RLock()
	n := len(c.items)
	c.mu.RUnlock()
	return n
}

func (c *cache) Flush() {
	c.mu.Lock()
	c.items = map[string]Item{}
	c.mu.UnLock()
}

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor) Run(c *cache) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			ticker:stop()
			return
		}
	}
}

func stopJanitor(c *Cache) {
	c.janitor.stop <- true
}

func runJanitor(c *cache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	c.janitor = j
	go j.Run(c)
}

func newCache(de time.Duration, m map[string]Item) *cache {
	if de == 0 {
		de = -1
	}
	c := &cache{
		defaultExpiration: de,
		items:			   m,
	}
	return c
}

func newCacheWithjanitor(de time.Duration, ci time.Duration, m map[string]Item) *Cache {
	c := newCache(de, m)
	C := &Cache{c}
	if ci > 0 {
		runJanitor(c, ci)
		runtime.SetFinalizer(C, stipJanitor)
	}
	return C
}

func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]Item)
	return newCacheWithjanitor(defaultExpiration, cleanupInterval, items)
}

func NewFrom(defaultExpiration, cleanupInterval time.Duration, items map[string]Item) *Cache{
	return newCacheWithjanitor(defaultExpiration, cleanupInterval, items)
}
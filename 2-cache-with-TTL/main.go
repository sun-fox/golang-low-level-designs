package main

import (
	"fmt"
	"sync"
	"time"
)

/**
	Cache with TTL:

	Structs:
	1. Cache -> map[key(string)]cacheItem
	2. CacheItem -> value(string), expiry(timestamp)

	Interface:
	1. Set(key,val,ttl) error
	2. Get(key) (val,error)
	3. Clear(key) error

**/

type cacheItem struct{
	value string
	expiresAt time.Time
}
type Cache struct{
	mu sync.RWMutex
	data map[string]cacheItem
}

func NewCache() *Cache{
	cache := &Cache{
		data: make(map[string]cacheItem),
	}

	go cache.StartCleanUpRoutine()
	return cache
}

type KeyValStoreWithTTL interface{
	Set(key, val string, ttl time.Duration)
	Get(key string) (string,error)
	Clear(key string) 
}

func (s *Cache) Set(key, val string, ttl time.Duration){
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = cacheItem{
		value: val,
		expiresAt: time.Now().Add(ttl),
	}
}

func (s *Cache) Get(key string) (string, error){
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, exists := s.data[key]
	if exists{
		if time.Now().Before(val.expiresAt){
			return val.value,nil
		}
		delete(s.data,key)
		return "", fmt.Errorf("Cache Expired") 	
	}
	return "",fmt.Errorf("Key Not Present")
}

func (s *Cache) Clear(key string) {
	delete(s.data, key)
}

func (s *Cache) StartCleanUpRoutine(){
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C{
		s.mu.Lock()
		now := time.Now()
		for k,v := range s.data{
			if now.After(v.expiresAt){
				fmt.Println("Cleared key:",k)
				delete(s.data,k)
			}
		}
		s.mu.Unlock()
	}
}

func main(){
	cache := NewCache()

	cache.Set("session1", "token_abc", 2*time.Second)

	val, err := cache.Get("session1")
	fmt.Println("Immediate Get:", val, err) 

	time.Sleep(3 * time.Second)

	val, err = cache.Get("session1")
	fmt.Println("Get after expiry:", val, err) 
}
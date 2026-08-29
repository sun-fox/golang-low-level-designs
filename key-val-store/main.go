package main

import (
	"fmt"
	"sync"
)

type keyValStore interface {
	Put(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

type memStore struct{
	mu sync.RWMutex
	data map[string]string
}

func NewMemStore() *memStore{
	return &memStore{
		data: make(map[string]string),
	}
}

func (s *memStore) Put(key string, value string) error{
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	return nil
}

func (s *memStore) Get(key string) (string, error){
	s.mu.RLock()
	defer s.mu.RUnlock()
	val,exist := s.data[key]
	if exist{
		return val, nil
	}
	return "", fmt.Errorf("Non existent Key")
}

func (s *memStore) Delete(key string){
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

func main(){
	cache := NewMemStore()
	err := cache.Put("navan","ssingh@navan.com")
	if err!=nil{
		fmt.Println("Failed to insert:",err)
	}
	err = cache.Put("grab","ssingh@grab.com")
	if err!=nil{
		fmt.Println("Failed to insert:",err)
	}
	val, err := cache.Get("navan")
	if err!=nil{
		fmt.Println(err.Error())
	}else{
		fmt.Println("Val:",val)
	}

}
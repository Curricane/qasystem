package session

import (
	"sync"
)

// MemorySession 内存中的session
type MemorySession struct {
	data   map[string]interface{}
	id     string
	rwlock sync.RWMutex
}

// NewMemorySession MemorySession构造函数
func NewMemorySession(id string) *MemorySession {
	s := &MemorySession{
		id:   id,
		data: make(map[string]interface{}, 16),
	}

	return s
}

// Set set key value
func (m *MemorySession) Set(key string, value interface{}) (err error) {

	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	m.data[key] = value
	return
}

// Get 根据key获取value
func (m *MemorySession) Get(key string) (value interface{}, err error) {
	m.rwlock.RLock()
	defer m.rwlock.RUnlock()

	value, ok := m.data[key]
	if !ok {
		err = ErrKeyNotExistInSession
		return
	}

	return
}

// Del 根据key删除
func (m *MemorySession) Del(key string) (err error) {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	// _, ok := m.data[key]
	// if !ok {
	// 	err = ErrKeyNotExistInSession
	// 	return
	// }

	delete(m.data, key)
	return
}

// Save 仅实现接口 do nothing
func (m *MemorySession) Save() (err error) {
	return
}

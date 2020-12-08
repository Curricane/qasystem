package session

import (
	"encoding/json"
	"sync"

	"github.com/garyburd/redigo/redis"
)

const (
	SessionFlagNone = iota
	SessionFlagModify
	SessionFlagLoad
)

type RedisSession struct {
	id  string
	pool       *redis.Pool
	data map[string]interface{} // 先存放session，后面再存到redis
	rwlock     sync.RWMutex
	flag       int
}

func NewRedisSession(id string, pool *redis.Pool) *RedisSession {
	s := &RedisSession{
		id:   id,
		data: make(map[string]interface{}, 16),
		flag: SessionFlagNone,
		pool: pool,
	}

	return s
}

func (r *RedisSession) Set(key string, value interface{}) error {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	r.data[key] = value
	r.flag = SessionFlagModify
	return nil
}

func (r *RedisSession) loadFromRedis() (err error) {
	conn := r.pool.Get()
	reply, err := conn.Do("GET", r.data)
	if err != nil {
		return
	}

	// 得到数据 json格式
	data, err := redis.String(reply, err)
	if err != nil {
		return
	}

	json.Unmarshal([]byte(data), &r.data)
	if err != nil {
		return
	}

	return
}

func (r *RedisSession) Get(key string) (result interface{}, err error) {
	r.rwlock.RLock()
	defer r.rwlock.RUnlock()

	// 实现延迟加载的功能
	if r.flag == SessionFlagNone {
		// 未加载，从redis中读取
		err = r.loadFromRedis()
		if err != nil {
			return
		}
	}

	result, ok := r.data[key]
	if !ok {
		err = ErrKeyNotExistInSession
		return
	}

	return
}
func (r *RedisSession) Del(key string) error {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	r.flag = SessionFlagModify
	delete(r.data, key)
	return nil
}
func (r *RedisSession) Save() (err error) {

	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	if r.flag != SessionFlagModify {
		return
	}

	data, err := json.Marshal(r.data)
	if err != nil {
		return
	}

	conn := r.pool.Get()
	_, err = conn.Do("SET", r.id, string(data))
	if err != nil {
		return
	}
	return
}

func (r *RedisSession) IsModify() bool {
	r.rwlock.RLock()
	defer r.rwlock.RUnlock()
	if r.flag == SessionFlagModify {
		return true
	}

	return false
}

func (r *RedisSession) Id() (id string) {
	r.rwlock.RLock()
	defer r.rwlock.RUnlock()

	return r.id
}
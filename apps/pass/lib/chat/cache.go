package chat

import "b0go/core/tools/cache"

var msgKey = "message:cache:key"

type MsgCache struct {
	// 单机部署使用，多服务部署时，使用redis代替
	cache *cache.Lrucache
	// todo redisCache *redis
}

func NewCache(size int) *MsgCache {
	cache, _ := cache.NewCache(size)
	return &MsgCache{cache: cache}
}

func (m *MsgCache) Set(msg []byte) error {
	return m.cache.Set(msgKey, msg, 60)
}
func (m *MsgCache) Get() []byte {
	msg, err := m.cache.Get(msgKey)
	if err != nil {
		return nil
	}
	return msg.([]byte)
}

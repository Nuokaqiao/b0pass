package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	HistoryMaxItems = 10
	HistoryMaxAge   = 72 * time.Hour
	historyRedisKey = "stonepass:text:history"
)

// HistoryItem 一条传内容记录
type HistoryItem struct {
	Key  string `json:"key"`
	Val  string `json:"val"`
	Ts   int64  `json:"ts"`   // unix 秒
	Time string `json:"time"` // 展示用
}

// HistoryStore 最近内容历史（最多 10 条，且不超过 72 小时）
type HistoryStore interface {
	Add(text string) (*HistoryItem, error)
	List() ([]HistoryItem, error)
	Backend() string // redis | memory
	Close() error
}

func formatTime(ts int64) string {
	return time.Unix(ts, 0).Local().Format("2006/1/2 15:04:05")
}

func newItem(text string) *HistoryItem {
	now := time.Now()
	ts := now.Unix()
	return &HistoryItem{
		Key:  fmt.Sprintf("%d-%d", now.UnixNano(), now.Unix()%1000),
		Val:  text,
		Ts:   ts,
		Time: formatTime(ts),
	}
}

func pruneItems(items []HistoryItem) []HistoryItem {
	cutoff := time.Now().Add(-HistoryMaxAge).Unix()
	out := make([]HistoryItem, 0, len(items))
	for _, it := range items {
		if it.Ts >= cutoff {
			out = append(out, it)
		}
	}
	if len(out) > HistoryMaxItems {
		out = out[len(out)-HistoryMaxItems:]
	}
	return out
}

// ---------- Memory fallback ----------

type MemoryHistory struct {
	mu    sync.Mutex
	items []HistoryItem
}

func NewMemoryHistory() *MemoryHistory {
	return &MemoryHistory{items: make([]HistoryItem, 0, HistoryMaxItems)}
}

func (m *MemoryHistory) Add(text string) (*HistoryItem, error) {
	item := newItem(text)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items = append(m.items, *item)
	m.items = pruneItems(m.items)
	cp := *item
	return &cp, nil
}

func (m *MemoryHistory) List() ([]HistoryItem, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items = pruneItems(m.items)
	out := make([]HistoryItem, len(m.items))
	copy(out, m.items)
	return out, nil
}

func (m *MemoryHistory) Close() error { return nil }

func (m *MemoryHistory) Backend() string { return "memory" }

// ---------- Redis ----------

type RedisHistory struct {
	rdb *redis.Client
}

type RedisOptions struct {
	Addr     string
	Password string
	DB       int
}

// NewHistoryStore RedisAddr 为空或连不上时回退内存，并打日志
func NewHistoryStore(opt RedisOptions) HistoryStore {
	addr := opt.Addr
	if addr == "" {
		log.Println("text history: RedisAddr empty, use memory store")
		return NewMemoryHistory()
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: opt.Password,
		DB:       opt.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Println("text history: redis ping fail, fallback memory:", err)
		_ = rdb.Close()
		return NewMemoryHistory()
	}
	log.Println("text history: redis connected", addr)
	return &RedisHistory{rdb: rdb}
}

func (r *RedisHistory) Add(text string) (*HistoryItem, error) {
	item := newItem(text)
	raw, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	pipe := r.rdb.Pipeline()
	pipe.ZAdd(ctx, historyRedisKey, &redis.Z{
		Score:  float64(item.Ts),
		Member: string(raw),
	})
	cutoff := time.Now().Add(-HistoryMaxAge).Unix()
	pipe.ZRemRangeByScore(ctx, historyRedisKey, "-inf", strconv.FormatInt(cutoff, 10))
	// 只保留最新 10 条
	pipe.ZRemRangeByRank(ctx, historyRedisKey, 0, int64(-(HistoryMaxItems + 1)))
	if _, err := pipe.Exec(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (r *RedisHistory) List() ([]HistoryItem, error) {
	ctx := context.Background()
	cutoff := time.Now().Add(-HistoryMaxAge).Unix()
	_ = r.rdb.ZRemRangeByScore(ctx, historyRedisKey, "-inf", strconv.FormatInt(cutoff, 10)).Err()
	_ = r.rdb.ZRemRangeByRank(ctx, historyRedisKey, 0, int64(-(HistoryMaxItems + 1))).Err()

	members, err := r.rdb.ZRange(ctx, historyRedisKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	out := make([]HistoryItem, 0, len(members))
	for _, m := range members {
		var it HistoryItem
		if err := json.Unmarshal([]byte(m), &it); err != nil {
			continue
		}
		if it.Time == "" && it.Ts > 0 {
			it.Time = formatTime(it.Ts)
		}
		out = append(out, it)
	}
	return out, nil
}

func (r *RedisHistory) Close() error {
	if r.rdb != nil {
		return r.rdb.Close()
	}
	return nil
}

func (r *RedisHistory) Backend() string { return "redis" }

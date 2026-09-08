package files

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const expireMetaName = ".b0pass-expire.json"

// ExpireStore 文件过期时间（unix 秒）。值为 0 或缺失表示不过期。
type ExpireStore struct {
	mu   sync.Mutex
	root string
	data map[string]int64
}

var (
	expireOnce  sync.Once
	expireStore *ExpireStore
)

// GetExpireStore 按共享根目录获取过期存储
func GetExpireStore(root string) *ExpireStore {
	expireOnce.Do(func() {
		expireStore = &ExpireStore{
			root: root,
			data: map[string]int64{},
		}
		expireStore.load()
	})
	return expireStore
}

func NormalizeRelPath(p string) string {
	p = strings.ReplaceAll(p, "\\", "/")
	p = strings.TrimSpace(p)
	if p == "" {
		return "/"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	p = path.Clean(p)
	if p == "." || p == "" {
		return "/"
	}
	return p
}

func (s *ExpireStore) metaPath() string {
	return filepath.Join(s.root, expireMetaName)
}

func (s *ExpireStore) load() {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, err := os.ReadFile(s.metaPath())
	if err != nil {
		s.data = map[string]int64{}
		return
	}
	m := map[string]int64{}
	if err := json.Unmarshal(b, &m); err != nil {
		log.Println("expire meta load fail:", err)
		s.data = map[string]int64{}
		return
	}
	s.data = m
}

func (s *ExpireStore) saveLocked() error {
	if s.data == nil {
		s.data = map[string]int64{}
	}
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.metaPath(), b, 0644)
}

// Get 返回过期 unix 秒；0 表示不过期
func (s *ExpireStore) Get(rel string) int64 {
	rel = NormalizeRelPath(rel)
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.data[rel]
}

// Set 设置过期时间；unix<=0 表示清除（不过期）
func (s *ExpireStore) Set(rel string, unix int64) error {
	rel = NormalizeRelPath(rel)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data == nil {
		s.data = map[string]int64{}
	}
	if unix <= 0 {
		delete(s.data, rel)
	} else {
		s.data[rel] = unix
	}
	return s.saveLocked()
}

// DeleteMeta 删除元数据（文件被手动删除时）
func (s *ExpireStore) DeleteMeta(rel string) {
	_ = s.Set(rel, 0)
}

// PurgeExpired 删除已过期文件并清理元数据
func (s *ExpireStore) PurgeExpired() int {
	now := time.Now().Unix()
	s.mu.Lock()
	expired := make([]string, 0)
	for rel, ts := range s.data {
		if ts > 0 && ts <= now {
			expired = append(expired, rel)
		}
	}
	s.mu.Unlock()

	removed := 0
	for _, rel := range expired {
		abs := filepath.Join(s.root, filepath.FromSlash(strings.TrimPrefix(rel, "/")))
		info, err := os.Stat(abs)
		if err == nil {
			if info.IsDir() {
				err = os.RemoveAll(abs)
			} else {
				err = os.Remove(abs)
			}
			if err != nil {
				log.Println("expire purge remove fail:", abs, err)
				continue
			}
			log.Println("expire purged:", abs)
			removed++
		} else if !os.IsNotExist(err) {
			log.Println("expire purge stat fail:", abs, err)
			continue
		}
		s.DeleteMeta(rel)
	}
	return removed
}

// StartExpireCleaner 后台定期清理过期文件
func StartExpireCleaner(root string, every time.Duration) {
	store := GetExpireStore(root)
	go func() {
		// 启动先清一次
		store.PurgeExpired()
		t := time.NewTicker(every)
		defer t.Stop()
		for range t.C {
			store.PurgeExpired()
		}
	}()
}

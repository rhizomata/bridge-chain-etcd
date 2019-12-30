package kv

// KV ..
type KV interface {
	Close() error
	PutObject(key string, value interface{}) (revision int64, err error)
	Put(key, val string) (revision int64, err error)
	GetOne(key string) (value []byte, err error)
	GetObject(key string, obj interface{}) (err error)
	GetWithPrefix(key string, handler func(key string, value []byte)) (err error)
	GetWithPrefixLimit(key string, limit int64, handler func(key string, value []byte)) (err error)
	DeleteOne(key string) (deleted bool, err error)
	DeleteWithPrefix(key string) (deleted int64, err error)
	Watch(key string, handler func(key string, value []byte)) *Watcher
	WatchWithPrefix(key string, handler func(key string, value []byte)) *Watcher
}

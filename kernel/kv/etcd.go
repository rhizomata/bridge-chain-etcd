package kv

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

// EtcdKV implements KV
type EtcdKV struct {
	etcdUrls []string
	client   *clientv3.Client
}

// Watcher ..
type Watcher struct {
	Key          string
	watchChannel clientv3.WatchChan
	running      bool
	handler      func(key string, value []byte)
}

func (watcher *Watcher) start() {
	for watchResp := range watcher.watchChannel {
		if !watcher.running {
			break
		}
		kv := watchResp.Events[0].Kv
		watcher.handler(watcher.Key, kv.Value)
	}
}

// Stop stop watching
func (watcher *Watcher) Stop() {
	watcher.running = false
}

// New : Create EtcdKV instance
func New(etcdUrls []string) (kv KV, err error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:            etcdUrls,
		DialTimeout:          3 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
	})

	if err != nil {
		log.Fatal("[BC-ERROR] Cannot connect to ETCD: ", etcdUrls, err)
		return nil, err
	}

	conn := client.ActiveConnection()

	log.Println("[BC-INFO] Connecting to ETCD :", conn.Target(), conn.GetState())

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Status(timeoutCtx, etcdUrls[0])

	if err != nil {
		log.Fatal("[BC-ERROR] Cannot connect to ETCD: ", etcdUrls[0], " : ", err)
	}

	etcd := EtcdKV{etcdUrls: etcdUrls, client: client}
	return &etcd, nil
}

// Close : close etcd client
func (etcd *EtcdKV) Close() error {
	return etcd.client.Close()
}

// Put ..
func (etcd *EtcdKV) Put(key, val string) (revision int64, err error) {
	r, err := etcd.put(context.Background(), key, val)
	return r.Header.Revision, err
}

// PutObject ..
func (etcd *EtcdKV) PutObject(key string, value interface{}) (revision int64, err error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		log.Println("[ERROR] Cannot Json marshal Object : ", err)
	}

	return etcd.Put(key, string(bytes))
}

func (etcd *EtcdKV) put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	r, err := etcd.client.Put(ctx, key, val, opts...)
	return r, err
}

// GetOne ..
func (etcd *EtcdKV) GetOne(key string) (value []byte, err error) {
	r, err := etcd.get(context.Background(), key)

	if err != nil {
		return nil, err
	}

	if r.Count > 0 {
		return r.Kvs[0].Value, nil
	}

	return nil, errors.New("No value for " + key)
}

// GetObject ..
func (etcd *EtcdKV) GetObject(key string, obj interface{}) (err error) {
	data, err := etcd.GetOne(key)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, obj)
	return err
}

// GetWithPrefix ..
func (etcd *EtcdKV) GetWithPrefix(key string, handler func(key string, value []byte)) (err error) {
	r, err := etcd.get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, item := range r.Kvs {
		handler(string(item.Key), item.Value)
	}

	return nil
}

// GetWithPrefixLimit ..
func (etcd *EtcdKV) GetWithPrefixLimit(key string, limit int64, handler func(key string, value []byte)) (err error) {
	r, err := etcd.get(context.Background(), key, clientv3.WithPrefix(), clientv3.WithLimit(limit))
	if err != nil {
		return err
	}

	for _, item := range r.Kvs {
		handler(string(item.Key), item.Value)
	}

	return nil
}

func (etcd *EtcdKV) get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	r, err := etcd.client.Get(ctx, key, opts...)
	return r, err
}

// DeleteOne ..
func (etcd *EtcdKV) DeleteOne(key string) (deleted bool, err error) {
	r, err := etcd.delete(context.Background(), key)
	if err != nil {
		return false, err
	}

	if r.Deleted != 1 {
		log.Println("[WARN] One more keys were deleted ", r.Deleted)
		return true, err
	}
	log.Println("[INFO] KV item deleted for ", key)
	return r.Deleted > 0, nil
}

// DeleteWithPrefix ..
func (etcd *EtcdKV) DeleteWithPrefix(key string) (deleted int64, err error) {
	r, err := etcd.delete(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		return 0, err
	}

	return r.Deleted, nil
}

func (etcd *EtcdKV) delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	r, err := etcd.client.Delete(ctx, key, opts...)
	return r, err
}

// Watch ..
func (etcd *EtcdKV) Watch(key string, handler func(key string, value []byte)) *Watcher {
	watchChannel := etcd.client.Watch(context.Background(), key)

	watcher := Watcher{Key: key, watchChannel: watchChannel, running: true, handler: handler}
	go func() {
		watcher.start()
	}()
	return &watcher
}

// WatchWithPrefix ..
func (etcd *EtcdKV) WatchWithPrefix(key string, handler func(key string, value []byte)) *Watcher {
	watchChannel := etcd.client.Watch(context.Background(), key, clientv3.WithPrefix())

	watcher := Watcher{Key: key, watchChannel: watchChannel, running: true, handler: handler}
	go func() {
		watcher.start()
	}()
	return &watcher
}

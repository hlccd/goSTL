package lru

import (
	"container/list"
	"sync"
)

type LRU struct {
	maxBytes int64
	nowBytes int64
	ll       *list.List
	cache    map[string]*list.Element
	onRemove func(key string, value Value)
	mutex    sync.Mutex //并发控制锁
}
type indexes struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onRemove func(string, Value)) *LRU {
	return &LRU{
		maxBytes: maxBytes,
		nowBytes: 0,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
		onRemove: onRemove,
	}
}
func (l *LRU) Add(key string, value Value) {
	l.mutex.Lock()
	if ele, ok := l.cache[key]; ok {
		l.ll.MoveToFront(ele)
		kv := ele.Value.(*indexes)
		//此处是一个替换,即将cache中的value替换为新的value,同时根据实际存储量修改其当前存储的实际大小
		l.nowBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := l.ll.PushFront(&indexes{key, value})
		l.cache[key] = ele
		//此处是一个增加操作,即原本不存在,所以直接插入即可,同时在当前数值范围内增加对应的占用空间
		l.nowBytes += int64(len(key)) + int64(value.Len())
	}
	//添加完成后根据删除掉尾部一部分数据以保证实际使用空间小于设定的空间上限
	for l.maxBytes != 0 && l.maxBytes < l.nowBytes {
		ele := l.ll.Back()
		if ele != nil {
			l.ll.Remove(ele)
			kv := ele.Value.(*indexes)
			//删除最末尾的同时将其占用空间减去
			delete(l.cache, kv.key)
			l.nowBytes -= int64(len(kv.key)) + int64(kv.value.Len())
			if l.onRemove != nil {
				//删除后的回调函数,可用于持久化该部分数据
				l.onRemove(kv.key, kv.value)
			}
		}
	}
	l.mutex.Unlock()
}
func (l *LRU) RemoveOldest() {
	l.mutex.Lock()
	ele := l.ll.Back()
	if ele != nil {
		l.ll.Remove(ele)
		kv := ele.Value.(*indexes)
		//删除最末尾的同时将其占用空间减去
		delete(l.cache, kv.key)
		l.nowBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if l.onRemove != nil {
			//删除后的回调函数,可用于持久化该部分数据
			l.onRemove(kv.key, kv.value)
		}
	}
	l.mutex.Unlock()
}
func (l *LRU) Get(key string) (value Value, ok bool) {
	l.mutex.Lock()
	if ele, ok := l.cache[key]; ok {
		l.ll.MoveToFront(ele)
		kv := ele.Value.(*indexes)
		l.mutex.Unlock()
		return kv.value, true
	}
	l.mutex.Unlock()
	return nil, false
}

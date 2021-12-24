package lru

//最近最少使用
//相对于仅考虑时间因素的FIFO和仅考虑访问频率的LFU，LRU算法可以认为是相对平衡的一种淘汰算法
//LRU认为，如果数据最近被访问过，那么将来被访问的概率也会更高
//LRU 算法的实现非常简单，维护一个队列
//如果某条记录被访问了，则移动到队尾，那么队首则是最近最少访问的数据，淘汰该条记录即可。
import (
	"container/list"
	"sync"
)

//LRU链结构体
type LRU struct {
	maxBytes int64                         //所能承载的最大byte数量
	nowBytes int64                         //当前承载的byte数
	ll       *list.List                    //用于存储的链表
	cache    map[string]*list.Element      //链表元素与key的映射表
	onRemove func(key string, value Value) //删除元素时的执行函数
	mutex    sync.Mutex                    //并发控制锁
}

//索引结构体,
type indexes struct {
	key   string //索引
	value Value  //存储值
}

//存储值的函数
type Value interface {
	Len() int//用于计算该存储值所占用的长度
}

func New(maxBytes int64, onRemove func(string, Value)) *LRU {
	return &LRU{
		maxBytes: maxBytes,
		nowBytes: 0,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
		onRemove: onRemove,
		mutex:    sync.Mutex{},
	}
}
//向该LRU中插入以key为索引的value
//若已经存在则将其放到队尾
//若空间充足且不存在则直接插入队尾
//若空间不足则淘汰队首元素再将其插入队尾
func (l *LRU) Insert(key string, value Value) {
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
//删除队尾元素
//删除后执行创建时传入的删除执行函数
func (l *LRU) Erase() {
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
//从LRU链中寻找以key为索引的value
//找到则返回,否则返回nil和false
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

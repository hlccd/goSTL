package consistentHash

import (
	"fmt"
	"sync"
)

const (
	maxReplicas = 32
)

//素数表
var primes = []uint32{3, 5, 7, 11, 13, 17, 19, 23, 29, 31,
	37, 41, 43, 47, 53, 59, 61, 67, 71, 73,
	79, 83, 89, 97, 101, 103, 107, 109, 113, 127,
	131, 137, 139, 149, 151, 157, 163, 167, 173, 179,
	121, 191, 193, 197, 199, 211, 223, 227, 229, 233,
	239, 241, 251,
}

//一致性hash结构体
type ConsistentHash struct {
	minReplicas int                      //最小的虚拟节点数
	keys        []uint32                 //存储的结点和虚拟结点的集合
	hashMap     map[uint32]interface{}   //hash与结点之间的映射
	nodeMap     map[interface{}][]uint32 //结点所对应的虚拟节点的hash值
	mutex       sync.Mutex               //并发控制锁
}

//hash计算
func hash(id int, v interface{}) (h uint32) {
	prime := primes[(id*id+len(primes))%len(primes)]
	h = uint32(0)
	s := fmt.Sprintf("%d-%v-%d", id*int(prime), v, prime)
	bs := []byte(s)
	for i := range bs {
		h += uint32(bs[i]) * prime
	}
	return h
}

//新建一个一致性hash结构体指针并返回
//传入其设定的最小的虚拟节点数量
func New(minReplicas int) *ConsistentHash {
	if minReplicas > maxReplicas {
		minReplicas = maxReplicas
	}
	if minReplicas < 0 {
		minReplicas = 1
	}
	ch := &ConsistentHash{
		minReplicas: minReplicas,
		keys:        make([]uint32, 0, 0),
		hashMap:     make(map[uint32]interface{}),
		nodeMap:     make(map[interface{}][]uint32),
		mutex:       sync.Mutex{},
	}
	return ch
}

//向一致性hash中插入结点,同时返回每一个结点插入的数量
func (ch *ConsistentHash) Insert(keys ...interface{}) (nums []int) {
	nums = make([]int, 0, len(keys))
	ch.mutex.Lock()
	//遍历所有待插入的结点
	for _, key := range keys {
		num := 0
		//判断结点是否已经存在
		_, exist := ch.nodeMap[key]
		if !exist {
			//结点不存在,开始插入
			for i := 0; i < maxReplicas || num < ch.minReplicas; i++ {
				//生成每个虚拟节点的hash值
				h := uint32(hash(i, key))
				//判断生成的hash值是否存在,存在则不插入
				_, ok := ch.hashMap[h]
				if !ok {
					//hash值不存在,进行插入
					num++
					ch.keys = append(ch.keys, h)
					//同时建立hash值和结点之间的映射关系
					ch.hashMap[h] = key
					ch.nodeMap[key] = append(ch.nodeMap[key], h)
				}
			}
		}
		nums = append(nums, num)
	}
	//对keys进行排序,以方便后续查找
	ch.sort(0, len(ch.keys)-1)
	ch.mutex.Unlock()
	return nums
}

//二分排序
func (ch *ConsistentHash) sort(L, R int) {
	if L >= R {
		return
	}
	l, r, m := L-1, R+1, ch.keys[(L+R)/2]
	for l < r {
		l++
		for ch.keys[l] < m {
			l++
		}
		r--
		for ch.keys[r] > m {
			r--
		}
		if l < r {
			tmp := ch.keys[l]
			ch.keys[l] = ch.keys[r]
			ch.keys[r] = tmp
		}
	}
	ch.sort(L, l-1)
	ch.sort(r+1, R)
}

//从一致性hash中获取以key为索引的下一个结点的索引
func (ch *ConsistentHash) Get(key interface{}) (ans interface{}) {
	if len(ch.keys) == 0 {
		return nil
	}
	ch.mutex.Lock()
	//计算key的hash值
	h := hash(0, key)
	//从现存的所有虚拟结点中找到该hash值对应的下一个虚拟节点对应的结点的索引
	idx := ch.search(h)
	ans = ch.hashMap[ch.keys[idx]]
	ch.mutex.Unlock()
	return ans
}

//二分查找,找到最近的不小于该值的hash值,如果不存在则返回0,即进行取模运算即可
func (ch *ConsistentHash) search(h uint32) (idx uint32) {
	l, m, r := uint32(0), uint32(len(ch.keys)/2), uint32(len(ch.keys))
	for l < r {
		m = (l + r) / 2
		if ch.keys[m] >= h {
			r = m
		} else {
			l = m + 1
		}
	}
	return l % uint32(len(ch.keys))
}

//删除以key为索引的结点
func (ch *ConsistentHash) Erase(key interface{}) {
	ch.mutex.Lock()
	hs, ok := ch.nodeMap[key]
	if ok {
		//该结点存在于一致性hash内
		//删除该结点
		delete(ch.nodeMap, key)
		//删除所有该结点的虚拟节点与该结点的映射关系
		for i := range hs {
			delete(ch.hashMap, hs[i])
		}
		//将待删除的虚拟结点删除即可
		arr := make([]uint32, 0, len(ch.keys)-len(hs))
		for i := range ch.keys {
			h := ch.keys[i]
			flag := true
			for j := range hs {
				if hs[j] == h {
					flag = false
				}
			}
			if flag {
				arr = append(arr, h)
			}
		}
		ch.keys = arr
	}
	ch.mutex.Unlock()
}

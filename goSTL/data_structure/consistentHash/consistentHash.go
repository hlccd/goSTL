package consistentHash

import (
	"fmt"
	"sync"
)

const (
	maxReplicas = 32
)

var primes = []byte{3, 5, 7, 11, 13, 17, 19, 23, 29, 31,
	37, 41, 43, 47, 53, 59, 61, 67, 71, 73,
	79, 83, 89, 97, 101, 103, 107, 109, 113, 127,
	131, 137, 139, 149, 151, 157, 163, 167, 173, 179,
	121, 191, 193, 197, 199, 211, 223, 227, 229, 233,
	239, 241, 251,
}

type ConsistentHash struct {
	minReplicas int
	keys        []uint32                 //存储的结点和虚拟结点的集合
	hashMap     map[uint32]interface{}   //hash与结点之间的映射
	nodeMap     map[interface{}][]uint32 //结点所对应的虚拟节点的hash值
	mutex       sync.Mutex
}

func hash(id int, v interface{}) (h uint32) {
	prime := primes[(id*id+len(primes))%len(primes)]
	h = uint32(0)
	s := fmt.Sprintf("%d-%v-%d", id*int(prime), v, prime)
	bs := []byte(s)
	for i := range bs {
		h += uint32(bs[i] * prime)
	}
	return h
}
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
func (ch *ConsistentHash) Insert(keys ...interface{}) (nums []int) {
	nums = make([]int, 0, len(keys))
	ch.mutex.Lock()
	for _, key := range keys {
		num := 0
		_, exist := ch.nodeMap[key]
		if !exist {
			for i := 0; i < maxReplicas || num < ch.minReplicas; i++ {
				h := uint32(hash(i, key))
				_, ok := ch.hashMap[h]
				if !ok {
					num++
					ch.keys = append(ch.keys, h)
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
func (ch *ConsistentHash) Get(key interface{}) (ans interface{}) {
	if len(ch.keys) == 0 {
		return nil
	}
	ch.mutex.Lock()
	h := hash(0, key)
	idx := ch.search(h)
	ans = ch.hashMap[ch.keys[idx]]
	ch.mutex.Unlock()
	return ans
}
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
func (ch *ConsistentHash) Erase(key interface{}) {
	ch.mutex.Lock()
	hs, ok := ch.nodeMap[key]
	if ok {
		delete(ch.nodeMap, key)
		for i := range hs {
			delete(ch.hashMap, hs[i])
		}
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

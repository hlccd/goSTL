package consistentHash

//@Title		consistentHash
//@Description
//		一致性哈希-consistent hash
//		一致性hash主要用以解决当出现增删结点时需要重新计算hash值的情况
//		同时,利用虚拟结点解决了数据倾斜的问题
//		hash范围为2^32,即一个uint32的全部范围
//		本次实现中不允许出现hash值相同的点,出现时则舍弃
//		虚拟结点最多有32个,最少数量可由用户自己决定,但不得低于1个
//		使用互斥锁实现并发控制
import (
	"fmt"
	"sync"
)

//最大虚拟节点数量
const (
	maxReplicas = 32
)

//素数表
//用以减少hash冲突
var primes = []uint32{3, 5, 7, 11, 13, 17, 19, 23, 29, 31,
	37, 41, 43, 47, 53, 59, 61, 67, 71, 73,
	79, 83, 89, 97, 101, 103, 107, 109, 113, 127,
	131, 137, 139, 149, 151, 157, 163, 167, 173, 179,
	121, 191, 193, 197, 199, 211, 223, 227, 229, 233,
	239, 241, 251,
}

//一致性hash结构体
//该实例了一致性hash在创建时设定的最小虚拟节点数
//同时保存了所有虚拟节点的hash值
//建立了虚拟节点hash值与实际结点之间的映射表
//每个实际结点也可以映射其所有的虚拟节点的hash值
//并发控制锁用以保证线程安全
type ConsistentHash struct {
	minReplicas int                      //最小的虚拟节点数
	keys        []uint32                 //存储的结点和虚拟结点的集合
	hashMap     map[uint32]interface{}   //hash与结点之间的映射
	nodeMap     map[interface{}][]uint32 //结点所对应的虚拟节点的hash值
	mutex       sync.Mutex               //并发控制锁
}

type consistentHasher interface {
	Size() (num int)                         //返回一致性hash中的虚拟节点数量
	Clear()                                  //清空一致性hash的所有结点
	Empty() (b bool)                         //返回该一致性hash是否有结点存在
	Insert(keys ...interface{}) (nums []int) //向该一致性hash中插入一组结点,同时返回其生成的虚拟节点的数量
	Erase(key interface{}) (ok bool)         //删除结点key
	Get(key interface{}) (ans interface{})   //查找key对应的结点
}

//@title    hash
//@description
//		传入一个虚拟节点id和实际结点
//		计算出它的hash值
//		先利用id从素数表中找到对应的素数,然后将id,素数和实际结点转化为[]byte
//		逐层访问并利用素数计算其hash值随后返回
//@receiver		nil
//@param    	id			int				虚拟节点的id
//@param    	v			interface{}		实际结点的key
//@return    	h			uint32			计算得到的hash值
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

//@title    New
//@description
//		新建一个一致性hash结构体指针并返回
//		传入设定的最少虚拟节点数量,不可小于1也不可大于最大虚拟节点数
//@receiver		nil
//@param    	minReplicas	int						hashMap的hash函数集
//@return    	ch			*ConsistentHash			新建的一致性hash指针
func New(minReplicas int) (ch *ConsistentHash) {
	if minReplicas > maxReplicas {
		//超过最大虚拟节点数
		minReplicas = maxReplicas
	}
	if minReplicas < 1 {
		//最少虚拟节点数量不得小于1
		minReplicas = 1
	}
	ch = &ConsistentHash{
		minReplicas: minReplicas,
		keys:        make([]uint32, 0, 0),
		hashMap:     make(map[uint32]interface{}),
		nodeMap:     make(map[interface{}][]uint32),
		mutex:       sync.Mutex{},
	}
	return ch
}

//@title    Size
//@description
//		以一致性hash做接收者
//		返回该容器当前含有的虚拟结点数量
//		如果容器为nil返回0
//@receiver		ch			*ConsistentHash			接受者一致性hash的指针
//@param    	nil
//@return    	num        	int						当前含有的虚拟结点数量
func (ch *ConsistentHash) Size() (num int) {
	if ch == nil {
		return 0
	}
	return len(ch.keys)
}

//@title    Clear
//@description
//		以一致性hash做接收者
//		将该容器中所承载的所有结点清除
//		被清除结点包括实际结点和虚拟节点
//		重建映射表
//@receiver		ch			*ConsistentHash			接受者一致性hash的指针
//@param    	nil
//@return    	nil
func (ch *ConsistentHash) Clear() {
	if ch == nil {
		return
	}
	ch.mutex.Lock()
	//重建vector并扩容到16
	ch.keys = make([]uint32, 0, 0)
	ch.hashMap = make(map[uint32]interface{})
	ch.nodeMap = make(map[interface{}][]uint32)
	ch.mutex.Unlock()
}

//@title    Empty
//@description
//		以一致性hash做接收者
//		判断该一致性hash中是否含有元素
//		如果含有结点则不为空,返回false
//		如果不含有结点则说明为空,返回true
//		如果容器不存在,返回true
//@receiver		ch			*ConsistentHash			接受者一致性hash的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (ch *ConsistentHash) Empty() (b bool) {
	if ch == nil {
		return false
	}
	return len(ch.keys) > 0
}

//@title    Insert
//@description
//		以一致性hash做接收者
//		向一致性hash中插入结点,同时返回每一个结点插入的数量
//		插入每个结点后将其对应的虚拟节点的hash放入映射表内和keys内
//		对keys做排序，利用二分排序算法
//@receiver		ch			*ConsistentHash			接受者一致性hash的指针
//@param    	keys		...interface{}			待插入的结点的集合
//@return    	nums		[]int					插入结点所生成的虚拟节点数量的集合
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

//@title    sort
//@description
//		以一致性hash做接收者
//		二分排序
//		主要用于一致性hash的keys的排序
//		以保证其有序性
//		同时方便后续使用二分查找
//@receiver		ch			*ConsistentHash			接受者一致性hash的指针
//@param    	L			int						左下标
//@param    	R			int						右下标
//@return    	nil
func (ch *ConsistentHash) sort(L, R int) {
	if L >= R {
		//左下标大于右下标，结束
		return
	}
	//找到中间结点，从左右两侧以双指针形式向中间靠近
	l, r, m := L-1, R+1, ch.keys[(L+R)/2]
	for l < r {
		//左侧出现不小于中间结点时停下
		l++
		for ch.keys[l] < m {
			l++
		}
		//右侧出现不大于中间结点时停下
		r--
		for ch.keys[r] > m {
			r--
		}
		if l < r {
			//左节点仍在右结点左方，交换结点的值
			ch.keys[l], ch.keys[r] = ch.keys[r], ch.keys[l]
		}
	}
	//递归排序左右两侧保证去有序性
	ch.sort(L, l-1)
	ch.sort(r+1, R)
}

//@title    Erase
//@description
//		以一致性hash做接收者
//		删除以key为索引的结点
//		先利用结点映射表判断该key是否存在于一致性hash内
//		若存在则从映射表内删除
//		同时找到其所有的虚拟节点的hash值
//		遍历keys进行删除
//@receiver		ch			*ConsistentHash			接受者一致性hash的指针
//@param    	key			interface{}				待删除元素的key
//@return    	b			bool					删除成功?
func (ch *ConsistentHash) Erase(key interface{}) (ok bool) {
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
	return ok
}

//@title    Get
//@description
//		以一致性hash做接收者
//		从一致性hash中获取以key为索引的下一个结点的索引
//		只要一致性hash内存储了结点则必然可以找到
//		将keys看成一个环
//		要找的key的计算出hash后放入环中
//		按顺时针向下找直到遇到第一个虚拟节点
//		寻找过程利用二分，若二分找到的结点为末尾结点的下一个，即为首个虚拟节点
//@receiver		ch			*ConsistentHash			接受者一致性hash的指针
//@param    	key			interface{}				待查找的结点
//@return    	ans			interface{}				找到的对应节点
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

//@title    search
//@description
//		以一致性hash做接收者
//		二分查找
//		找到最近的不小于该值的hash值
//		如果不存在则返回0,即进行取模运算即可
//@receiver		ch			*ConsistentHash			接受者一致性hash的指针
//@param    	h			uint32					待查找结点的hash值
//@return    	idx			uint32					找到的对应虚拟节点的下标
func (ch *ConsistentHash) search(h uint32) (idx uint32) {
	//二分查找，寻找不小于该hash值的下一个值的下标
	l, m, r := uint32(0), uint32(len(ch.keys)/2), uint32(len(ch.keys))
	for l < r {
		m = (l + r) / 2
		if ch.keys[m] >= h {
			r = m
		} else {
			l = m + 1
		}
	}
	//当找到的下标等同于keys的长度时即为0
	return l % uint32(len(ch.keys))
}

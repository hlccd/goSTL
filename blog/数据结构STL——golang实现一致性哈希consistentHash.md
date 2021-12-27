github仓库存储地址：https://github.com/hlccd/goSTL

### 概述

​		一致性哈希（consistent hash），与一致性哈希相对的是不一致性哈希，但常见的所有的哈希几乎都是不一致的，即哈希桶的容量的不固定的，可以根据需求进行扩容和缩容，不一致性哈希可以提高空间的利用率，但相应的，当进行扩容和缩容操作时需要对桶内存储的所有元素重新计算哈希值，这在某些情况是十分麻烦的事情，特别是在分布式存储的环境下，此时每个哈希结点就相当于一个机器，文件分布在哪台机器由哈希算法来决定，这个系统想要加一台机器时就需要停下来等所有文件重新分布一次才能对外提供服务，而当一台机器掉线的时候尽管只掉了一部分数据，但所有数据访问路由都会出问题。这样整个服务就无法平滑的扩缩容，成为了有状态的服务。

​		为了避免分布式存储出现的状态化问题，就需要利用一致性哈希去解决，一致性哈希是固定了桶的容量，即2^32，然后将结点的哈希值计算后放入其中，当要从一致性哈希中寻找机器结点的时候，可以将hash桶视为一个大小为2^32的环状结构，要查找的结点也计算出hash并放入其中，然后沿顺时针寻找直到遇到第一个机器结点即可。

​		而当增删机器结点的时候，只需要去调度该机器节点临近的一些结点即可，避免了对所有节点的全部调度。

​		当机器节点的hash值比较密集的集中在环上的一部分的时候，就有可能出现**数据倾斜**的问题，即有较多存储结点会映射到其中一个机器结点上，造成该节点的负担较大。对于该问题一般采用增加**虚拟结点**的方案去解决，即对于每一个机器节点，可以凭空创造一些结点作为该虚拟节点的等价替代结点并放入hash环内，以此来充盈整个hash环，从而尽可能减少数据倾斜堆的概率。**虚拟节点数量一般采用32**，主要由于hash环的大小为2^32，采用32个虚拟节点可以近似使得转换出的32bit的每一位都可以存在一个虚拟节点从而实现充盈hash环的目的。

### 原理

​		对于一致性hash来说，它需要做的一方面是对机器结点和其对应的虚拟节点的key进行hash计算，然后将计算机插入到hash环上，随后将hash环进行排序；另一方面是计算用于查找的结点的hash值并利用hash环找到其对应的机器节点或机器节点的虚拟节点，然后映射到机器节点并返回。

​		考虑到存储的值仅仅只有机器节点及其虚拟结点，同时机器节点也可以视为其虚拟节点，所以主要需要解决的问题就变成了虚拟节点之间的hash冲突问题。

​		对于虚拟节点的hash冲突来说，由于一致性hash的寻找不应当出现一次找到两个不同的结点，即hash冲突出现时不应该通过拓展的方式去解决。所以可以考虑采用再次hash法，即当一个虚拟节点与之前的结点发送冲突时则废弃该结点，然后重新生成一个虚拟节点在进行计算比较即可。

​		同时，由于虚拟节点可以通过map映射到机器节点上，所以理论上生成的虚拟节点都可以进行直接映射，而没有太多的额外限制，故重新生成虚拟节点放入hash环内一方面可以解决hash冲突，另一方面也不会出现找不到对应的机器节点的情况。

​		除此之外，考虑到再hash也需要一定的时间，所以在进行hash计算的时候，一方面利用**素数做逐位城际累加**，另一方面利用虚拟结点的**id做平方法**去从素数表中找到对应的素数做计算。从而尽可能降低出现hash冲突的概率。

### 实现

##### hash

```go
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
```

​		传入一个虚拟节点id和实际结点，计算出它的hash值，先利用id从素数表中找到对应的素数,然后将id,素数和实际结点转化为[]byte，逐层访问并利用素数计算其hash值随后返回。

```go
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
```

#### 结构体

​		一致性hash结构体，该实例了一致性hash在创建时设定的最小虚拟节点数，同时保存了所有虚拟节点的hash值，建立了虚拟节点hash值与实际结点之间的映射表，每个实际结点也可以映射其所有的虚拟节点的hash值，并发控制锁用以保证线程安全。

```go
type ConsistentHash struct {
	minReplicas int                      //最小的虚拟节点数
	keys        []uint32                 //存储的结点和虚拟结点的集合
	hashMap     map[uint32]interface{}   //hash与结点之间的映射
	nodeMap     map[interface{}][]uint32 //结点所对应的虚拟节点的hash值
	mutex       sync.Mutex               //并发控制锁
}
```

#### 接口

```go
type consistentHasher interface {
	Size() (num int)                         //返回一致性hash中的虚拟节点数量
	Clear()                                  //清空一致性hash的所有结点
	Empty() (b bool)                         //返回该一致性hash是否有结点存在
	Insert(keys ...interface{}) (nums []int) //向该一致性hash中插入一组结点,同时返回其生成的虚拟节点的数量
	Erase(key interface{}) (ok bool)         //删除结点key
	Get(key interface{}) (ans interface{})   //查找key对应的结点
}
```

#### New

​		新建一个一致性hash结构体指针并返回，传入设定的最少虚拟节点数量,不可小于1也不可大于最大虚拟节点数。

```go
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
```

#### Size

​		以一致性hash做接收者，返回该容器当前含有的虚拟结点数量，如果容器为nil返回0。

```go
func (ch *ConsistentHash) Size() (num int) {
	if ch == nil {
		return 0
	}
	return len(ch.keys)
}
```

#### Clear

​		以一致性hash做接收者，将该容器中所承载的所有结点清除，被清除结点包括实际结点和虚拟节点，重建映射表。

```go
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
```

#### Empty

​		以一致性hash做接收者，判断该一致性hash中是否含有元素，如果含有结点则不为空,返回false，如果不含有结点则说明为空,返回true，如果容器不存在,返回true。

```go
func (ch *ConsistentHash) Empty() (b bool) {
	if ch == nil {
		return false
	}
	return len(ch.keys) > 0
}
```

#### Insert

​		以一致性hash做接收者，向一致性hash中插入结点,同时返回每一个结点插入的数量，插入每个结点后将其对应的虚拟节点的hash放入映射表内和keys内，对keys做排序，利用二分排序算法。

```go
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
```

##### sort

​		以一致性hash做接收者,二分排序,主要用于一致性hash的keys的排序,以保证其有序性,同时方便后续使用二分查找。

```go
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
```

#### Erase

​		以一致性hash做接收者，删除以key为索引的结点，先利用结点映射表判断该key是否存在于一致性hash内，若存在则从映射表内删除，同时找到其所有的虚拟节点的hash值，遍历keys进行删除。

```go
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
```

#### Get

​		以一致性hash做接收者，从一致性hash中获取以key为索引的下一个结点的索引，只要一致性hash内存储了结点则必然可以找到，将keys看成一个环，要找的key的计算出hash后放入环中，按顺时针向下找直到遇到第一个虚拟节点，寻找过程利用二分，若二分找到的结点为末尾结点的下一个，即为首个虚拟节点。

```go
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
```

##### search

​		以一致性hash做接收者，二分查找，找到最近的不小于该值的hash值，如果不存在则返回0,即进行取模运算即可。

```go
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
```

### 使用示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/consistentHash"
)

func main() {
	ch:=consistentHash.New(32)
	fmt.Println(ch.Insert("http://localhost:8001","http://localhost:8002","http://localhost:8003"))
	fmt.Println(ch.Size())
	fmt.Println(ch.Get("group"))
	fmt.Println(ch.Get("http://localhost:8002"))
}
```

> [32 32 32]
> 96
> http://localhost:8001
> http://localhost:8002

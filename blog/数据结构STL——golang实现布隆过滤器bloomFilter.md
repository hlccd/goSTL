github仓库存储地址：https://github.com/hlccd/goSTL

### 概述

​		布隆过滤器（bloom filter）,它实际上是一个很长的二进制向量和一系列随机映射函数。布隆过滤器可以用于**寻找**一个元素是否在一个集合中，但由于元素的通过hash映射转化到集合内的，所以存在误差，即可以百分百判断其不存在，但不能确定其一定存在。它的优点是空间效率和查询时间都比一般的算法要好的多，缺点是有一定的**误识别率**和**删除困难**。

### 原理

​		如果想要判断一个元素是不是在一个集合里，除了可以将所有元素保存起来，通过比较确定外，也可以通过将元素转化为一个hash值放入哈希表内，通过hash函数将元素转化为hash值，该方案可以通过**位图**的形式存贮，当判断其是否存在的时候只需要在位图上找到这个点是否为1即可判断其是否存在。极大的提高了空间利用率同时减小了时间开销，但带来了一定的**识别错误率**，识别错误主要由hash冲突导致。

​		相比于其它的数据结构，布隆过滤器在空间和时间方面都有巨大的优势。布隆过滤器存储空间和插入/查询时间都是常熟级别。另外，布隆过滤器不需要存储元素本身，在某些对保密要求非常严格的场合有优势。同时，布隆过滤器可以用来表示全集，这是其他数据结构所不能实现的。布隆过滤器可以表示全集，其它任何数据结构都不能。

​		对于出现hash冲突的情况来说，它是会出现错误识别的，当然，一般来说并不会有很多。对于此类问题，可以通过建立白名单的方式去存储可能误判的元素，即对hash使用开放寻址法，但一般来说并不需要。

​		同时，考虑到识别错误的问题，删除的操作几乎是不应当实现的，所以本次实现中主要实现**添加和查找**操作，当然也包含清空操作。

### 实现

​		bloomFilter布隆过滤器结构体，包含其用于存储的uint64元素切片，选用uint64是为了更多的利用bit位。

```go
type bloomFilteror interface {
	Insert(v interface{})         //向布隆过滤器中插入v
	Check(v interface{}) (b bool) //检查该值是否存在于布隆过滤器中,该校验存在误差
	Clear()                       //清空该布隆过滤器
}
```

#### 接口

```go
type bloomFilteror interface {
	Insert(v interface{})         //向布隆过滤器中插入v
	Check(v interface{}) (b bool) //检查该值是否存在于布隆过滤器中,该校验存在误差
	Clear()                       //清空该布隆过滤器
}
```

#### hash

​		传入一个虚拟节点id和实际结点，计算出它的hash值，逐层访问并利用素数131计算其hash值随后返回。

​		该hash函数可自行传入，也可使用库内函数

```go
type Hash func(v interface{}) (h uint32)

func hash(v interface{}) (h uint32) {
	h = uint32(0)
	s := fmt.Sprintf("131-%v-%v", v,v)
	bs := []byte(s)
	for i := range bs {
		h += uint32(bs[i]) * 131
	}
	return h
}
```

#### New

​		新建一个bloomFilter布隆过滤器容器并返回，初始bloomFilter的切片数组为空。

```go
func New(h Hash) (bf *bloomFilter) {
	if h == nil {
		h = hash
	}
	return &bloomFilter{
		bits: make([]uint64, 0, 0),
		hash: h,
	}
}
```

#### Insert

​		以bloomFilter布隆过滤器容器做接收者，先将待插入的value计算得到其哈希值hash，再向布隆过滤器中第hash位插入一个元素(下标从0开始)，当hash大于当前所能存储的位范围时,需要进行扩增，若要插入的位比冗余的多不足2^16即1024*64时,则新增1024个uint64，否则则直接增加到可以容纳第hash位的位置,以此可以提高冗余量,避免多次增加。

```go
func (bf *bloomFilter) Insert(v interface{}) {
	//bm不存在时直接结束
	if bf == nil {
		return
	}
	//开始插入
	h := bf.hash(v)
	if h/64+1 > uint32(len(bf.bits)) {
		//当前冗余量小于num位,需要扩增
		var tmp []uint64
		//通过冗余扩增减少扩增次数
		if h/64+1 < uint32(len(bf.bits)+1024) {
			//入的位比冗余的多不足2^16即1024*64时,则新增1024个uint64
			tmp = make([]uint64, len(bf.bits)+1024)
		} else {
			//直接增加到可以容纳第num位的位置
			tmp = make([]uint64, h/64+1)
		}
		//将原有元素复制到新增的切片内,并将bm所指向的修改为扩增后的
		copy(tmp, bf.bits)
		bf.bits = tmp
	}
	//将第num位设为1即实现插入
	bf.bits[h/64] ^= 1 << (h % 64)
}
```

#### Check

​		以bloomFilter布隆过滤器容器做接收者，将待查找的值做哈希计算得到哈希值h，检验第h位在位图中是否存在，当h大于当前所能存储的位范围时,直接返回false，否则判断第h为是否为1,为1返回true,否则返回false，利用布隆过滤器做判断存在误差,即返回true可能也不存在,但返回false则必然不存在。

```go
func (bf *bloomFilter) Check(v interface{}) (b bool) {
	//bf不存在时直接返回false并结束
	if bf == nil {
		return false
	}
	h := bf.hash(v)
	//h超出范围,直接返回false并结束
	if h/64+1 > uint32(len(bf.bits)) {
		return false
	}
	//判断第num是否为1,为1返回true,否则为false
	if bf.bits[h/64]&(1<<(h%64)) > 0 {
		return true
	}
	return false
}
```

#### Clear

​		以bloomFilter布隆过滤器容器做接收者，清空整个布隆过滤器。

```go
func (bf *bloomFilter) Clear() {
	if bf == nil {
		return
	}
	bf.bits = make([]uint64, 0, 0)
}
```

### 使用示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/bloomFilter"
)

func hash(v interface{}) uint32 {
	return uint32(v.(int))
}
func main() {
	bf := bloomFilter.New(nil)
	for i := 0; i < 10; i++ {
		bf.Insert(i)
	}
	for i := 0; i < 15; i++ {
		fmt.Println(i,bf.Check(i))
	}
}
```

> 0 true
> 1 true
> 2 true
> 3 true
> 4 true
> 5 true
> 6 true
> 7 true
> 8 true
> 9 true
> 10 false
> 11 false
> 12 false
> 13 false
> 14 false


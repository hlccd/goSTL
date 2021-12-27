github仓库存储地址：https://github.com/hlccd/goSTL

### 概述

​		向量（Vector）是一个封装了动态大小数组的顺序容器。跟任意其它类型容器一样，它能够存放**各种类型的对象**。可以简单的认为，向量是一个能够存放任意类型的动态数组。

​		对于vector的实现,考虑到它是一个线性容器,其中的元素可以通过下标在O(1)的时间复杂度直接获取访问，可以考虑**底层使用数组**的形式实现。由于vector可以动态的添加和删除元素，所以需要对数组进行动态分配以满足需求。

​		对于golang语言来说，官方并未设计泛型，但由于interface可以充当任意类型的特性，可以使用interface类型作为泛型类型。

### 原理

​		对于一个vector来说，它是可以动态的对元素进行添加和修改的，而如果说每一次增删元素都使得数组恰好可以容纳所有元素的话，那它在以后的每一次增删的过程中都需要去**重新分配空间**，并将原有的元素复制到新数组中，如果按照这样来操作的话，每次删减都会需要花费大量的时间去进行，将会极大的降低其性能。

​		为了提高效率，可以选择牺牲一定的空间作为冗余量，即**时空权衡**，通过每次扩容时多分配一定的空间，这样在后续添加时就可以减少扩容次数，而在删除元素时，如果冗余量不多时，可以选择不进行空间分配，仅将后续元素前移一位，同理，在中间插入元素也仅仅将其后移一位，这样可以极大减少的由于容量不足或超出导致的扩容缩容造成的时间开销。

#### 扩容策略

​		对于动态数组扩容来说，由于其添加是线性且连续的，即每一次只会增加一个，不像bitmap一样可能需要一次添加很多的空间才能保证在其范围内，所以它的扩容策略相对于bitmap有所不同。

​		vector的扩容主要由两种方案：

1. 固定扩容：固定的去增加一定量的空间，该方案时候在数组较大时添加空间，可以避免直接翻倍导致的**冗余过量**问题，到由于增加量是固定的，如果需要一次扩容很多量的话就会比较缓慢。
2. 翻倍扩容：将原有容量直接翻倍，该方案适合在数组不太大的适合添加空间，可以提高扩容量，增加扩容效率，但数组太大时使用的话会导致一次性增加太多空间，进而造成空间的浪费。

#### 缩容策略

​		对于动态数组的缩容来说，同扩容类似，由于元素的减少只会是线性且连续的，即每一次只会减少一个，不会存在由于一个元素被删去导致其他已使用的空间也需要被释放。所以vector的缩容方案和vector的扩容方案十分类似，即是它的逆过程：

1. 固定缩容：释放一个固定的空间，该方案适合在当前数组较大的时候进行，可以减缓需要缩小的量，当空间较大的时候如果采取折半缩减的话;将可能在绝大多数时间内都不会被采用，从而造成空间浪费。
2. 折半缩容：将当前容量进行折半，该方案适合数组不太大的适合进行缩容，可以更快的缩小容量，但对于较大的空间来说并不会频繁使用到。

### 实现

​		由于vector底层由动态数组实现，所以其必然包含一个数组指针，同时，利用时空权衡去减少时间开销导致实际使用空间和实际占用空间不一致，所以分别利用len和cap进行存储。同时，为了解决在高并发情况下的数据不一致问题，引入了并发控制锁。

```go
type vector struct {
   data  []interface{} //动态数组
   len   uint64        //当前已用数量
   cap   uint64        //可容纳元素数量
   mutex sync.Mutex    //并发控制锁
}
```

#### 接口

```go
type vectorer interface {
	Iterator() (i *Iterator.Iterator)  //返回一个包含vector所有元素的迭代器
	Sort(Cmp ...comparator.Comparator) //利用比较器对其进行排序
	Size() (num uint64)                //返回vector的长度
	Clear()                            //清空vector
	Empty() (b bool)                   //返回vector是否为空,为空则返回true反之返回false
	PushBack(e interface{})            //向vector末尾插入一个元素
	PopBack()                          //弹出vector末尾元素
	Insert(idx uint64, e interface{})  //向vector第idx的位置插入元素e,同时idx后的其他元素向后退一位
	Erase(idx uint64)                  //删除vector的第idx个元素
	Reverse()                          //逆转vector中的数据顺序
	At(idx uint64) (e interface{})     //返回vector的第idx的元素
	Front() (e interface{})            //返回vector的第一个元素
	Back() (e interface{})             //返回vector的最后一个元素
}
```

#### New

​		创建一个vector容器并初始化，同时返回其指针。

```go
func New() (v *vector) {
	return &vector{
		data:  make([]interface{}, 1, 1),
		len:   0,
		cap:   1,
		mutex: sync.Mutex{},
	}
}
```

#### Iterator

​		以vector为接受器，返回一个包含vector中所有元素的迭代器，同时将vector中暂未使用的空间进行释放。

```go
func (v *vector) Iterator() (i *Iterator.Iterator) {
	if v == nil {
		v = New()
	}
	v.mutex.Lock()
	if v.data == nil {
		//data不存在,新建一个
		v.data = make([]interface{}, 1, 1)
		v.len = 0
		v.cap = 1
	} else if v.len < v.cap {
		//释放未使用的空间
		tmp := make([]interface{}, v.len, v.len)
		copy(tmp, v.data)
		v.data = tmp
	}
	//创建迭代器
	i = Iterator.New(&v.data)
	v.mutex.Unlock()
	return i
}
```

#### Sort

​		以vector为接受器，重载了比较器中的Sort，实现了让vector利用比较器中的Sort进行排序，使用过程中将会释放vector中暂未使用的空间。

```go
func (v *vector) Sort(Cmp ...comparator.Comparator) {
   if v == nil {
      v = New()
   }
   v.mutex.Lock()
   if v.data == nil {
      //data不存在,新建一个
      v.data = make([]interface{}, 1, 1)
      v.len = 0
      v.cap = 1
   } else if v.len < v.cap {
      //释放未使用空间
      tmp := make([]interface{}, v.len, v.len)
      copy(tmp, v.data)
      v.data = tmp
      v.cap = v.len
   }
   //调用比较器的Sort进行排序
   if len(Cmp) == 0 {
      comparator.Sort(&v.data)
   } else {
      comparator.Sort(&v.data, Cmp[0])
   }
   v.mutex.Unlock()
}
```

#### Size

​		返回vector中当前所包含的元素个数，由于其数值必然是非负整数，所以选用了**uint64**。

```go
func (v *vector) Size() (num uint64) {
   if v == nil {
      v = New()
   }
   return v.len
}
```

#### Clear

​		清空了vector中所承载的所有元素。

```go
func (v *vector) Clear() {
   if v == nil {
      v = New()
   }
   v.mutex.Lock()
   //清空data
   v.data = make([]interface{}, 1, 1)
   v.len = 0
   v.cap = 1
   v.mutex.Unlock()
}
```

#### Empty

​		判断vector中是否为空，通过Size()是否为0进行判断。

```go
func (v *vector) Empty() (b bool) {
   if v == nil {
      v = New()
   }
   return v.Size() <= 0
}
```

#### PushBack

​		以vector为接受器，向vector尾部添加一个元素，添加元素会出现两种情况，第一种是还有冗余量，此时**直接覆盖**以len为下标指向的位置即可，另一种情况是没有冗余量了，需要对动态数组进行**扩容**，此时就需要利用扩容策略。

​		扩容策略上文已做描述，可返回参考，该实现过程种将两者进行了结合使用，可参考下方注释。

​		固定扩容值设为2^16，翻倍扩容上限也为2^16。

```go
func (v *vector) PushBack(e interface{}) {
	if v == nil {
		v = New()
	}
	v.mutex.Lock()
	if v.len < v.cap {
		//还有冗余,直接添加
		v.data[v.len] = e
	} else {
		//冗余不足,需要扩容
		if v.cap <= 65536 {
			//容量翻倍
			if v.cap == 0 {
				v.cap = 1
			}
			v.cap *= 2
		} else {
			//容量增加2^16
			v.cap += 65536
		}
		//复制扩容前的元素
		tmp := make([]interface{}, v.cap, v.cap)
		copy(tmp, v.data)
		v.data = tmp
		v.data[v.len] = e
	}
	v.len++
	v.mutex.Unlock()
}
```

#### PopBack

​		以vector向量容器做接收者，弹出容器最后一个元素,同时长度--即可，若容器为空,则不进行弹出，当弹出元素后,可能进行缩容，当容量和实际使用差值超过2^16时,容量直接减去2^16，否则,当实际使用长度是容量的一半时,进行折半缩容。其缩容策略可参考前方解释。

```go
func (v *vector) PopBack() {
   if v == nil {
      v = New()
   }
   if v.Empty() {
      return
   }
   v.mutex.Lock()
   v.len--
   if v.cap-v.len >= 65536 {
      //容量和实际使用差值超过2^16时,容量直接减去2^16
      v.cap -= 65536
      tmp := make([]interface{}, v.cap, v.cap)
      copy(tmp, v.data)
      v.data = tmp
   } else if v.len*2 < v.cap {
      //实际使用长度是容量的一半时,进行折半缩容
      v.cap /= 2
      tmp := make([]interface{}, v.cap, v.cap)
      copy(tmp, v.data)
      v.data = tmp
   }
   v.mutex.Unlock()
}
```

#### Insert

​		向vector中任意位置插入一个元素，由于其位置必然是非负整数，所以选用了**uint64**进行表示，同时，当插入位置**超出了数组范围**时，将在最近的位置进行插入，由于不可能是负数，所以超出范围的情况只有可能是插入下标大于len，此时只需要插入尾部即可，等同于PushBack的操作。

​		对于插入的下标在范围中时，只需要将后续的一些元素后移一位即可。

​		是否需要扩容需要在开始进行判断，扩容策略同上。

```go
func (v *vector) Insert(idx uint64, e interface{}) {
	if v == nil {
		v = New()
	}
	v.mutex.Lock()
	if v.len >= v.cap {
		//冗余不足,进行扩容
		if v.cap <= 65536 {
			//容量翻倍
			if v.cap == 0 {
				v.cap = 1
			}
			v.cap *= 2
		} else {
			//容量增加2^16
			v.cap += 65536
		}
		//复制扩容前的元素
		tmp := make([]interface{}, v.cap, v.cap)
		copy(tmp, v.data)
		v.data = tmp
	}
	//从后往前复制,即将idx后的全部后移一位即可
	var p uint64
	for p = v.len; p > 0 && p > uint64(idx); p-- {
		v.data[p] = v.data[p-1]
	}
	v.data[p] = e
	v.len++
	v.mutex.Unlock()
}
```

#### Erase

​		删除下标为idx的元素，idx不在数组范围内时，则删除距离其最近的元素，即下标不大于0删除首部，不小于len删除尾部。

​		删除后进行缩容判断，缩容操作同上方介绍的缩容策略。

```go
func (v *vector) Erase(idx uint64) {
   if v == nil {
      v = New()
   }
   if v.Empty() {
      return
   }
   v.mutex.Lock()
   for p := idx; p < v.len-1; p++ {
      v.data[p] = v.data[p+1]
   }
   v.len--
   if v.cap-v.len >= 65536 {
      //容量和实际使用差值超过2^16时,容量直接减去2^16
      v.cap -= 65536
      tmp := make([]interface{}, v.cap, v.cap)
      copy(tmp, v.data)
      v.data = tmp
   } else if v.len*2 < v.cap {
      //实际使用长度是容量的一半时,进行折半缩容
      v.cap /= 2
      tmp := make([]interface{}, v.cap, v.cap)
      copy(tmp, v.data)
      v.data = tmp
   }
   v.mutex.Unlock()
}
```

#### Reverse

​		以vector为接受器，将vector中所承载的元素**逆转**，同时，释放掉所有未使用的空间。

```go
func (v *vector) Reverse() {
   if v == nil {
      v = New()
   }
   v.mutex.Lock()
   if v.data == nil {
      //data不存在,新建一个
      v.data = make([]interface{}, 1, 1)
      v.len = 0
      v.cap = 1
   } else if v.len < v.cap {
      //释放未使用的空间
      tmp := make([]interface{}, v.len, v.len)
      copy(tmp, v.data)
      v.data = tmp
      v.cap=v.len
   }
   for i := uint64(0); i < v.len/2; i++ {
      v.data[i], v.data[v.len-i-1] = v.data[v.len-i-1], v.data[i]
   }
   v.mutex.Unlock()
}
```

#### At

​		返回下标为idx的元素，当idx不在数组范围内时返回nil，当idx在范围内时**返回对应元素**。

```go
func (v *vector) At(idx uint64) (e interface{}) {
   if v == nil {
      v=New()
      return nil
   }
   v.mutex.Lock()
   if idx < 0 && idx >= v.Size() {
      v.mutex.Unlock()
      return nil
   }
   if v.Size() > 0 {
      e = v.data[idx]
      v.mutex.Unlock()
      return e
   }
   v.mutex.Unlock()
   return nil
}
```

#### Front

​		以vector为接受器，返回vector所承载的元素中位于**首部**的元素，如果vector为nil或者元素数组为nil或为空，则返回nil。

```go
func (v *vector) Front() (e interface{}) {
   if v == nil {
      v=New()
      return nil
   }
   v.mutex.Lock()
   if v.Size() > 0 {
      e = v.data[0]
      v.mutex.Unlock()
      return e
   }
   v.mutex.Unlock()
   return nil
}
```

#### Back

​		以vector为接受器，返回vector所承载的元素中位于**尾部**的元素，如果vector为nil或者元素数组为nil或为空，则返回nil。

```go
func (v *vector) Back() (e interface{}) {
	if v == nil {
		v=New()
		return nil
	}
	v.mutex.Lock()
	if v.Size() > 0 {
		e = v.data[v.len-1]
		v.mutex.Unlock()
		return e
	}
	v.mutex.Unlock()
	return nil
}
```

### 使用示例

```go
package main

import (
   "fmt"
   "github.com/hlccd/goSTL/data_structure/vector"
   "sync"
)

func main() {
   v:=vector.New()
   wg:=sync.WaitGroup{}
   for i:=0;i<10;i++{
      wg.Add(1)
      go func(num int) {
         v.PushBack(num)
         v.Insert(uint64(num),num)
         wg.Done()
      }(i)
   }
   wg.Wait()
   fmt.Println("随机插入后的结果:")
   for i := v.Iterator(); i.HasNext(); i.Next() {
      fmt.Println(i.Value())
   }
   v.Sort()
   fmt.Println("排序后的结果:")
   for i := v.Iterator(); i.HasNext(); i.Next() {
      fmt.Println(i.Value())
   }
   fmt.Println("随机删除一半的结果:")
   for i:=0;i<5;i++{
      wg.Add(1)
      go func(num int) {
         v.PopBack()
         v.Erase(uint64(i))
         wg.Wait()
      }(i)
   }
   wg.Done()
   for i:=uint64(0); i<v.Size(); i++ {
      fmt.Println(v.At(i))
   }
   fmt.Println("逆转后的结果:")
   v.Reverse()
   for i:=uint64(0); i<v.Size(); i++ {
      fmt.Println(v.At(i))
   }
   fmt.Println("在下标5处插入元素-1的结果:")
   v.Insert(5,-1)
   for i:=uint64(0); i<v.Size(); i++ {
      fmt.Println(v.At(i))
   }

   fmt.Println("首元素:",v.Front())
   fmt.Println("尾元素:",v.Back())
}
```

注：由于过程中的增删过程是并发执行的，所以其结果和下方示例并不完全相同

> 随机插入后的结果:
> 0
> 1
> 0
> 2
> 3
> 5
> 6
> 7
> 4
> 8
> 1
> 9
> 2
> 3
> 4
> 5
> 6
> 9
> 8
> 7
> 排序后的结果:
> 0
> 0
> 1
> 1
> 2
> 2
> 3
> 3
> 4
> 4
> 5
> 5
> 6
> 6
> 7
> 7
> 8
> 8
> 9
> 9
> 随机删除一半的结果:
> 0
> 0
> 1
> 3
> 3
> 5
> 5
> 6
> 6
> 7
> 逆转后的结果:
> 7
> 6
> 6
> 5
> 5
> 3
> 3
> 1
> 0
> 0
> 在下标5处插入元素-1的结果:
> 7
> 6
> 6
> 5
> 5
> -1
> 3
> 3
> 1
> 0
> 0
> 首元素: 7
> 尾元素: 0


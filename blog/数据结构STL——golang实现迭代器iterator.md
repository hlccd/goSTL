github仓库存储地址：https://github.com/hlccd/goSTL

### Iterator

#### 概述

iterator模式：提供一种方法，使之能依次访问容器内的各个元素，而又不暴露该聚合物内部的表述方式。
STL的中心思想是将算法与数据结构分离，彼此独立设计，最后在用iterator将他们结合在一起，获得最大的适配性。

由于golang官方未实现泛型，而interface存在可以替换为任意结构的特性，故而可以**使用interface实现泛型**。

一个迭代器需要包括其存储的元素序列以及内部存储一个该迭代器当前指向的元素下标。

```go
//Iterator迭代器
//包含泛型切片和该迭代器当前指向元素的下标
//可通过下标和泛型切片长度来判断是否可以前移或后移
//当index不小于0时迭代器可前移
//当index小于data的长度时可后移
type Iterator struct {
	data  *[]interface{} 	//该迭代器中存放的元素集合的指针
	index int           	//该迭代器当前指向的元素下标，-1即不存在元素
}
```

对于一个迭代器,它需要执行的主要方法有：返回该首尾迭代器，获取指向某位的迭代器，获取指定位置的元素，判断能否前后移动以及进行前后移动。

```go
type Iteratorer interface {
	Begin() (I *Iterator)       				//将该迭代器设为位于首节点并返回新迭代器
	End() (I *Iterator)         				//将该迭代器设为位于尾节点并返回新迭代器
	Get(idx int) (I *Iterator)  				//将该迭代器设为位于第idx节点并返回该迭代器
	Value() (e interface{})    					//返回该迭代器下标所指元素
	HasNext() (b bool)      				    //判断该迭代器是否可以后移
	Next() (b bool)         				    //将该迭代器后移一位
	HasPre() (b bool)        					//判罚该迭代器是否可以前移
	Pre() (b bool)           					//将该迭代器前移一位
}
```

#### 接口实现

##### New

对于迭代器的初始化，需要传入一个所要承载的元素集合的指针，可以自己选择性传入一个index下标。该函数会返回一个承载了该元素集合的迭代器。

```go
func New(data *[]interface{}, Idx ...int) (i *Iterator) {
	//迭代器下标
	var idx int
	if len(Idx) <= 0 {
		//没有传入下标，则将下标设为0
		idx = 0
	} else {
		//有传入下标，则将传入下标第一个设为迭代器下标
		idx = Idx[0]
	}
	if len((*data)) > 0 {
		//如果元素集合非空，则判断下标是否超过元素集合范围
		if idx >= len((*data)) {
			//如果传入下标超过元素集合范围则寻找最近的下标值
			idx = len((*data)) - 1
		}
	} else {
		//如果元素集合为空则将下标设为-1
		idx = -1
	}
	//新建并返回迭代器
	return &Iterator{
		data:  data,
		index: idx,
	}
}
```

##### Bgein

以原有的迭代器做接受器，返回一个承载相同元素但下标位于首位的迭代器。

```go
func (i *Iterator) Begin() (I *Iterator) {
	if i == nil {
		//迭代器为空，直接结束
		return nil
	}
	if len((*i.data)) == 0 {
		//迭代器元素集合为空，下标设为-1
		i.index = -1
	} else {
		//迭代器元素集合非空，下标设为0
		i.index = 0
	}
	//返回修改后的新指针
	return &Iterator{
		data:  i.data,
		index: i.index,
	}
}
```

##### End

以原有的迭代器做接受器，返回一个承载相同元素但下标位于末尾的迭代器。

```go
func (i *Iterator) End() (I *Iterator) {
	if i == nil {
		//迭代器为空，直接返回
		return nil
	}
	if len((*i.data)) == 0 {
		//元素集合为空，下标设为-1
		i.index = -1
	} else {
		//元素集合非空，下标设为最后一个元素的下标
		i.index = len((*i.data)) - 1
	}
	//返回修改后的该指针
	return &Iterator{
		data:  i.data,
		index: i.index,
	}
}
```

##### Get

​		以原有的迭代器做接受器，返回一个承载相同元素但下标为自己传入的idx的迭代器，当idx不在元素集合的范围内时，下标设为距离idx最近的值，即小于0设为首位，大于元素集合长度则设为尾部，其他情况设为idx。

```go
func (i *Iterator) Get(idx int) (I *Iterator) {
	if i == nil {
		//迭代器为空，直接返回
		return nil
	}
	if idx <= 0 {
		//预设下标超过元素集合范围，将下标设为最近元素的下标，此状态下为首元素下标
		idx = 0
	} else if idx >= len((*i.data))-1 {
		//预设下标超过元素集合范围，将下标设为最近元素的下标，此状态下为尾元素下标
		idx = len((*i.data)) - 1
	}
	if len((*i.data)) > 0 {
		//元素集合非空，迭代器下标设为预设下标
		i.index = idx
	} else {
		//元素集合为空，迭代器下标设为-1
		i.index = -1
	}
	//返回修改后的迭代器指针
	return i
}
```

##### Value

以原有的迭代器做接受器，返回该迭代器当前下标指向的元素

```go
func (i *Iterator) Value() (e interface{}) {
	if i == nil {
		//迭代器为nil，返回nil
		return nil
	}
	if len((*i.data)) == 0 {
		//元素集合为空，返回nil
		return nil
	}
	if i.index <= 0 {
		//下标超过元素集合范围下限，最近元素为首元素
		i.index = 0
	}
	if i.index >= len((*i.data)) {
		//下标超过元素集合范围上限，最近元素为尾元素
		i.index = len((*i.data)) - 1
	}
	//返回下标指向元素
	return (*i.data)[i.index]
}
```

##### HasNext

以原有的迭代器做接受器，判断该迭代器是否可以进行后移，可以则返回true否则返回false。

```go
func (i *Iterator) HasNext() (b bool) {
	if i == nil {
		//迭代器为nil时不能后移
		return false
	}
	if len((*i.data)) == 0 {
		//元素集合为空时不能后移
		return false
	}
	//下标到达元素集合上限时不能后移,否则可以后移
	return i.index < len((*i.data))
}
```

##### Next

​		以原有的迭代器做接受器，将迭代器下标后移，当满足后移条件时进行后移同时返回true，当不满足后移条件时将下标设为尾元素下标同时返回false，当迭代器为nil时返回false。

```go
func (i *Iterator) Next() (b bool) {
	if i == nil {
		//迭代器为nil时返回false
		return false
	}
	if i.HasNext() {
		//满足后移条件时进行后移
		i.index++
		return true
	}
	if len((*i.data)) == 0 {
		//元素集合为空时下标设为-1同时返回false
		i.index = -1
		return false
	}
	//不满足后移条件时将下标设为尾元素下标并返回false
	i.index = len((*i.data)) - 1
	return false
}
```

##### HasPre

以原有的迭代器做接受器，判断该迭代器是否可以进行前移，可以则返回true否则返回false。

```go
func (i *Iterator) End() (I *Iterator) {
	if i == nil {
		//迭代器为空，直接返回
		return nil
	}
	if len((*i.data)) == 0 {
		//元素集合为空，下标设为-1
		i.index = -1
	} else {
		//元素集合非空，下标设为最后一个元素的下标
		i.index = len((*i.data)) - 1
	}
	//返回修改后的该指针
	return &Iterator{
		data:  i.data,
		index: i.index,
	}
}
```

##### Pre

​		以原有的迭代器做接受器，将迭代器下标前移，当满足前移条件时进行前移同时返回true，当不满足前移条件时将下标设为首元素下标同时返回false，当迭代器为nil时返回false，当元素集合为空时下标设为-1同时返回false。

```go
func (i *Iterator) Pre() (b bool) {
	if i == nil {
		//迭代器为nil时返回false
		return false
	}
	if i.HasPre() {
		//满足后移条件时进行前移
		i.index--
		return true
	}
	if len((*i.data)) == 0 {
		//元素集合为空时下标设为-1同时返回false
		i.index = -1
		return false
	}
	//不满足后移条件时将下标设为尾元素下标并返回false
	i.index = 0
	return false
}
```

#### 使用示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/utils/iterator"
)

func main() {
	var arr =make([]interface{},0,0)
	arr=append(arr,5)
	arr=append(arr,3)
	arr=append(arr,2)
	arr=append(arr,4)
	arr=append(arr,1)
	arr=append(arr,4)
	arr=append(arr,3)
	arr=append(arr,1)
	arr=append(arr,5)
	arr=append(arr,2)
	i:=Iterator.New(&arr)
	fmt.Println("begin")
	for i:=i.Begin();i.HasNext();i.Next(){
		fmt.Println(i.Value())
	}
	fmt.Println()
	fmt.Println("end")
	for i:=i.End();i.HasPre();i.Pre(){
		fmt.Println(i.Value())
	}
	fmt.Println()
	fmt.Println("get4")
	for i:=i.Get(4);i.HasNext();i.Next(){
		fmt.Println(i.Value())
	}
	fmt.Println()
}
```

##### 示例结果

> begin
> 5
> 3
> 2
> 4
> 1
> 4
> 3
> 1
> 5
> 2
>
> end
> 2
> 5
> 1
> 3
> 4
> 1
> 4
> 2
> 3
> 5
>
> get4
> 1
> 4
> 3
> 1
> 5
> 2
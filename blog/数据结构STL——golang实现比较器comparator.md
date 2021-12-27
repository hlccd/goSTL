github仓库存储地址：https://github.com/hlccd/goSTL

### Comparator

#### 概述

​		对于某些存在大小比较的数据结构，如果每次都要特定的实现一些大小比较是十分繁琐的，特别是对于一些官方已经设定的类型，如果将基本类型引入数据结构中时需要实现其元素的比较会简单。

​		同时，对于一些常用的函数，比如排序、查找、排序第n个以及寻找上下边界的函数，这些函数需要通过比较器进行配合实现，为了更进一步简化使用，可以在比较器中实现。

#### 定义

对于一个比较器，除开基本类型外，必须传入比较函数，当然，基本数据类型也可以传入自定的比较函数进行覆盖，对于待使用的比较函数，需要传入两个元素a和b（有先后），同时返回一个int，其中**0表示相等，正数表示a>b，负数表示a<b**。

```go
type Comparator func(a, b interface{}) int
```

除比较器外，同时等一一个判断相等的工具——相等器，用于其他数据结构中不需要比较大小只需要判断是否相等的工具。

```
type Equaler func(a, b interface{}) (B bool)
```



#### GetCmp

对于一些基本数据类型，可以预先设定好比较函数以节约使用时必须重写基本数据类型的比较函数的过程。

当然，对于要进行比较的基本类型来说，需要传入一个对象以获取其数据类型从而返回对应的默认比较器。

以下部分**直接复制即可**，其实现仅仅只是判断类型并返回默认比较器，并无需要理解的部分。

```go
func GetCmp(e interface{}) (cmp Comparator) {
	if e==nil{
		return nil
	}
	switch e.(type) {
	case bool:
		return boolCmp
	case int:
		return intCmp
	case int8:
		return int8Cmp
	case uint8:
		return uint8Cmp
	case int16:
		return int16Cmp
	case uint16:
		return uint16Cmp
	case int32:
		return int32Cmp
	case uint32:
		return uint32Cmp
	case int64:
		return int64Cmp
	case uint64:
		return uint64Cmp
	case float32:
		return float32Cmp
	case float64:
		return float64Cmp
	case complex64:
		return complex64Cmp
	case complex128:
		return complex128Cmp
	case string:
		return stringCmp
	}
	return nil
}

```

##### basicCmp

```go
//以下为系统自带类型的默认比较器

func boolCmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(bool) {
		return 1
	} else if b.(bool) {
		return -1
	}
	return 0
}
func intCmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(int) > b.(int) {
		return 1
	} else if a.(int) < b.(int) {
		return -1
	}
	return 0
}
func int8Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(int8) > b.(int8) {
		return 1
	} else if a.(int8) < b.(int8) {
		return -1
	}
	return 0
}
func uint8Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(uint8) > b.(uint8) {
		return 1
	} else if a.(uint8) < b.(uint8) {
		return -1
	}
	return 0
}
func int16Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(int16) > b.(int16) {
		return 1
	} else if a.(int16) < b.(int16) {
		return -1
	}
	return 0
}
func uint16Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(uint16) > b.(uint16) {
		return 1
	} else if a.(uint16) < b.(uint16) {
		return -1
	}
	return 0
}
func int32Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(int32) > b.(int32) {
		return 1
	} else if a.(int32) < b.(int32) {
		return -1
	}
	return 0
}
func uint32Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(uint32) > b.(uint32) {
		return 1
	} else if a.(uint32) < b.(uint32) {
		return -1
	}
	return 0
}
func int64Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(int64) > b.(int64) {
		return 1
	} else if a.(int64) < b.(int64) {
		return -1
	}
	return 0
}
func uint64Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(uint64) > b.(uint64) {
		return 1
	} else if a.(uint64) < b.(uint64) {
		return -1
	}
	return 0
}
func float32Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(float32) > b.(float32) {
		return 1
	} else if a.(float32) < b.(float32) {
		return -1
	}
	return 0
}
func float64Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(float64) > b.(float64) {
		return 1
	} else if a.(float64) < b.(float64) {
		return -1
	}
	return 0
}
func complex64Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if real(a.(complex64)) > real(b.(complex64)) {
		return 1
	} else if real(a.(complex64)) < real(b.(complex64)) {
		return -1
	} else {
		if imag(a.(complex64)) > imag(b.(complex64)) {
			return 1
		} else if imag(a.(complex64)) < imag(b.(complex64)) {
			return -1
		}
	}
	return 0
}
func complex128Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if real(a.(complex128)) > real(b.(complex128)) {
		return 1
	} else if real(a.(complex128)) < real(b.(complex128)) {
		return -1
	} else {
		if imag(a.(complex128)) > imag(b.(complex128)) {
			return 1
		} else if imag(a.(complex128)) < imag(b.(complex128)) {
			return -1
		}
	}
	return 0
}
func stringCmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if len(a.(string)) > len(b.(string)) {
		return 1
	} else if len(a.(string)) < len(b.(string)) {
		return -1
	} else {
		if a.(string) > b.(string) {
			return 1
		} else if a.(string) < b.(string) {
			return -1
		}
	}
	return 0
}
```

#### GetEqual

获取默认的相等器。

```go
func GetEqual() (equ Equaler) {
	return basicEqual
}
```

##### basicEqual

​		默认的相等器，即两个元素完全相等。

```go
func basicEqual(a, b interface{}) (B bool) {
   return a == b
}
```

#### Sort

​		排序，对于传入的数组，通过传入的比较函数（基本类型可通过GetCmp直接获取，不需要传入）进行排序。

​		为了简化实现，可以将待比较元素组限定为线性数组，而对于一些非线性结构，可以通过将其指针列为线性表，再通过传入的比较函数进行比较即可，但一般不建议非线性结构使用比较器，其结构最好自己去维护大小分布，这样可以更好的提高效率。

​		排序中分别实现了**二分排序**和**归并排序**，由于二分排序本身不稳定，所以更加适合数据量不太大的数组，而对于归并排序，其性能是十分稳定的，更加适合用于数据量较大的数组，所以对待排序数组根据长度进行了区分，以根据情况使用合适的排序方式。

​		传入的数组是其**指针**，传入指针可以减少复制的情况，以节约时间。

```go
func Sort(arr *[]interface{}, Cmp ...Comparator) {
	//如果传入一个空数组或nil,则直接结束
	if arr==nil || (*arr)==nil || len((*arr)) == 0 {
		return
	}
	var cmp Comparator
	cmp = nil
	if len(Cmp) > 0 {
		cmp = Cmp[0]
	} else {
		cmp = GetCmp((*arr)[0])
	}
	if cmp == nil {
		//未传入比较比较函数且并非默认类型导致未找到默认比较器则直接终止排序
		return
	}
	//根据数组长度进行分类选择排序函数
	if len((*arr)) < 2^26 {
		//当长度小于2^16时使用二分排序
		binary(arr,0,len((*arr))-1, cmp)
	} else {
		merge(arr,0,len((*arr))-1, cmp)
	}
}
```

##### Binary

​		对于二分排序，其原理主要是将一个无须数组，通过寻找一个中间量（一般是数组中间点的值）作为参考，通过比较和交换使得数组中间量左侧始终不大于中间量，右侧始终不小于中间量，即**相对有序的状态**，再对两侧的情况进行递归排序，从而使得每一部分都是相对有序的，依次保证整体的有序性。

​		当中间值取的较差甚至是极值的时候，二分排序将会退化为冒泡排序，该排序方案并不稳定，但由于不需要额外空间进行存储，所以可以在较小的数组中进行使用。

```go
func binary(arr *[]interface{},l,r int, cmp Comparator) {
	//对当前部分进行预排序,使得两侧都大于或小于中间值
	if l >= r {
		return
	}
	m := (*arr)[(r + l) / 2]
	i, j := l-1, r+1
	for i < j {
		i++
		for cmp((*arr)[i], m) < 0 {
			i++
		}
		j--
		for cmp((*arr)[j],m) > 0 {
			j--
		}
		if i < j {
			(*arr)[i],(*arr)[j]=(*arr)[j],(*arr)[i]
		}
	}
	//对分好的两侧进行迭代二分排序
	binary(arr,l,j, cmp)
	binary(arr,j+1,r, cmp)
}
```

##### Merge

​		对于归并排序，其原理主要是将一个无须数组，先通过对数组的左右拆分，随后从最小部分（一般是仅存在两个或一个元素）进行初步排序，然后依次向上合并，由于合并的两个小数组的内部都是有序的，所以只需要依次遍历比较其大小即可，与此同时，需要一个临时数组去存储比较结果，随后再将临时数组中的值放入待排序的数组中去，依次合并到整个数组来保证其有序性。

​		由于该方案是将数组逐步拆分为最小单元在进行合并，其情况必然是相对稳定的，虽然需要一定的额外空间进行存储，但相较于其稳定性来说也是值得的空间成本。

```go
func merge(arr *[]interface{},l,r int, cmp Comparator) {
	//对当前部分进行分组排序,将该部分近似平均的拆为两部分进行比较排序
	if l >= r {
		return
	}
	m := (r + l) / 2
	//对待排序内容进行二分
	merge(arr,l,m, cmp)
	merge(arr,m+1,r, cmp)
	//二分结束后依次比较进行归并
	i, j := l, m+1
	var tmp []interface{}=make([]interface{},0,r-l+1)
	for i <= m && j <= r {
		if cmp((*arr)[i], (*arr)[j]) <= 0 {
			tmp = append(tmp, (*arr)[i])
			i++
		} else {
			tmp = append(tmp, (*arr)[j])
			j++
		}
	}
	//当一方比较到头时将另一方剩余内容全部加入进去
	for ; i <= m; i++ {
		tmp = append(tmp, (*arr)[i])
	}
	for ; j <= r; j++ {
		tmp = append(tmp, (*arr)[j])
	}
	//将局部排序结果放入迭代器中
	for i, j = l, 0; i <= r; i, j = i+1, j+1 {
		(*arr)[i]=tmp[j]
	}
}
```

##### 示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/utils/comparator"
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
	comparator.Sort(&arr)
	for i:=0;i< len(arr);i++{
		println(arr[i].(int))
	}
}
```

#### Search

​		对于一个**有序线性表**来说，如果要从中查找某个元素，可以通过二分的方式进行查找，即通过比对该元素和当前区间的中值的大小，从而判断去掉左侧或右侧区间，然后继续进行比较直至剩下一个元素即可，此时只需要比较该元素和待查找元素是否相等，相等则返回该元素下标，不等返回-1表示未找到该元素。

​		先对传入是参数进行判断，以保证在查找时不会出现error，判断我完成后再调用search函数进行查找

```go
func Search(arr *[]interface{}, e interface{}, Cmp ...Comparator) (idx int) {
	if arr==nil || (*arr)==nil || len((*arr)) == 0 {
		return
	}
	//判断比较函数是否有效,若无效则寻找默认比较器
	var cmp Comparator
	cmp = nil
	if len(Cmp) == 0 {
		cmp = GetCmp(e)
	} else {
		cmp = Cmp[0]
	}
	if cmp == nil {
		//若并非默认类型且未传入比较比较函数则直接结束
		return -1
	}
	//查找开始
	return search(arr, e, cmp)
}
```

##### search

​		二分查找函数：

```go
func search(arr *[]interface{}, e interface{}, cmp Comparator) (idx int) {
	//通过二分查找的方式寻找该元素
	l, m, r := 0, (len((*arr))-1)/2, len((*arr))
	for l < r {
		m = (l + r) / 2
		if cmp((*arr)[m], e) < 0 {
			l = m + 1
		} else {
			r = m
		}
	}
	//查找结束
	if (*arr)[l] == e {
		//该元素存在,返回下标
		return l
	}
	//该元素不存在,返回-1
	return -1
}
```

##### 示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/utils/comparator"
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
	comparator.Sort(&arr)
	for i:=0;i< len(arr);i++{
		println(arr[i].(int))
	}
	fmt.Println("search:",comparator.Search(&arr,3))
}
```

#### NthElement

​		对于要找的第n个元素**（下标从0开始）**，可以使用和二分排序类似的方式，只不过，由于只需要将第n个元素放到第n位，所以只需要排序第n位所在的区间，即对利用中间值进行二分之后所获得的两个区间，**只需要对包含n的区间进行二分即可**，另一个区间可以直接不管。

​		该函数结果将会返回位于n的元素，该过程分两部分进行，第一部分验证其可执行情况，即指针是否为nil，数组是否为nil，n是否超出数组范围之类的情况，若出现类似情况，则直接返回nil即可。否则则对数组进行有限排序。

```go
func NthElement(arr *[]interface{}, n int, Cmp ...Comparator) (value interface{}){
	if arr==nil || (*arr)==nil || len((*arr)) == 0 {
		return nil
	}
	//判断比较函数是否有效
	var cmp Comparator
	cmp = nil
	if len(Cmp) > 0 {
		cmp = Cmp[0]
	} else {
		cmp = GetCmp((*arr)[0])
	}
	if cmp == nil {
		return nil
	}
	//判断待确认的第n位是否在该集合范围内
	if len((*arr)) < n || n<0 {
		return nil
	}
	//进行查找
	nthElement(arr,0,len((*arr))-1, n, cmp)
	return (*arr)[n]
}
```

##### nthElement

实现部分不返回任何值，仅仅是进行有限排序，即仅排序包含n的区间。

```go
func nthElement(arr *[]interface{},l,r int, n int, cmp Comparator){
	//二分该区域并对此进行预排序
	if l >= r {
		return
	}
	m := (*arr)[(r + l) / 2]
	i, j := l-1, r+1
	for i < j {
		i++
		for cmp((*arr)[i], m) < 0 {
			i++
		}
		j--
		for cmp((*arr)[j], m) > 0 {
			j--
		}
		if i < j {
			(*arr)[i],(*arr)[j]=(*arr)[j],(*arr)[i]
		}
	}
	//确认第n位的范围进行局部二分
	if n-1 >= i {
		nthElement(arr,j+1,r, n, cmp)
	} else {
		nthElement(arr,l,j, n, cmp)
	}
}
```

##### 示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/utils/comparator"
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
	for i:=0;i< len(arr);i++{
		fmt.Println("n:",comparator.NthElement(&arr,i))
	}
}
```

#### Bound

​		对于一组有序线性表，可以通过二分查找的方式，获得它的上下界，当待查找元素不存在于该线性表内时，则返回的上届是小于它的最大值，返回的下届是大于它的最小值。当元素存在于线性表内时，返回的上界是该元素最右侧的下标，下届是该元素最左侧的下标

​		查找方法是二分查找的变形，返回查找值的边界。

##### UpperBound

```go
func UpperBound(arr *[]interface{}, e interface{}, Cmp ...Comparator) (idx int) {
	if arr==nil || (*arr)==nil || len((*arr)) == 0 {
		return -1
	}
	//判断比较函数是否有效
	var cmp Comparator
	cmp = nil
	if len(Cmp) == 0 {
		cmp = GetCmp(e)
	} else {
		cmp = Cmp[0]
	}
	if cmp == nil {
		return -1
	}
	//寻找该元素的上界
	return upperBound(arr, e, cmp)
}
```

##### upperBound

```go
func upperBound(arr *[]interface{}, e interface{}, cmp Comparator) (idx int) {
	l, m, r := 0, len((*arr)) / 2, len((*arr))-1
	for l < r {
		m = (l + r + 1) / 2
		if cmp((*arr)[m], e) <= 0 {
			l = m
		} else {
			r = m - 1
		}
	}
	return l
}
```

##### LowerBound

```go
func LowerBound(arr *[]interface{}, e interface{}, Cmp ...Comparator) (idx int) {
	if arr==nil || (*arr)==nil || len((*arr)) == 0 {
		return -1
	}
	//判断比较函数是否有效
	var cmp Comparator
	cmp = nil
	if len(Cmp) == 0 {
		cmp = GetCmp(e)
	} else {
		cmp = Cmp[0]
	}
	if cmp == nil {
		return -1
	}
	//寻找该元素的下界
	return lowerBound(arr, e, cmp)
}
```

##### lowerBound

```go
func lowerBound(arr *[]interface{}, e interface{}, cmp Comparator) (idx int) {
	l, m, r := 0, len((*arr)) / 2, len((*arr))
	for l < r {
		m = (l + r) / 2
		if cmp((*arr)[m], e) >= 0 {
			r = m
		} else {
			l = m + 1
		}
	}
	return l
}
```

##### 示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/utils/comparator"
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
	comparator.Sort(&arr)
	for i:=0;i< len(arr);i++{
		fmt.Println(i,"=",arr[i])
	}
	fmt.Println("\n\n\n")
	for i:=0;i< len(arr);i++{
		fmt.Println(i)
		fmt.Println("upper:",comparator.UpperBound(&arr,i))
		fmt.Println("lower:",comparator.LowerBound(&arr,i))
		fmt.Println()
	}
}
```


github仓库存储地址：https://github.com/hlccd/goSTL

### 概述

​		树堆（Treap）是一个比较特殊的结构，它同之前实现的二叉搜索树和完全二叉搜有着类似的性质，它**即满足二叉搜索树的查找性质，又满足完全二叉树的极值处于堆顶的性质**。

​		树堆由一个根节点和根节点下属的多层次的子结点构成，任意一个结点最多只能拥有两个子结点，即左右子结点。它的这一特性**同二叉搜索树完全一致**，但除此之外，它还有另一个特质：给每一个结点赋予一个随机的优先级，使得其在插入的时候，通过**左右旋转**的方式，保证其满足二叉搜索树的性质，同时使得优先级更大的结点旋转至更靠近顶部的方式，即**满足完全二叉树的性质**。

​		除开以上两点之外，它还有一个比较特殊的优势，考虑到所有结点的优先级都是随机生成的，即其优先级和其结点值之前完全没有关系，故只要随机算法合理的情况下，整个树堆的形状就会更加扁平，从概率上来说，当插入结点数量够多的情况下，**增加和删除结点的时间复杂的仅为O（logn）**，即实现了**依概率平衡**。

### 原理

​		对于一个树堆来说，它需要满足的条件主要有两条：

1. 父结点的值必然大于左结点小于右结点
2. 父结点的优先级必然大于左右两结点

​		而为了实现这两个情况，我们可以考虑先实现分别实现第一点和第二点。对于第一点来说其实就是指的**二叉搜索树**的性质，第二点指的是完全二叉树的，考虑到之前的文章已经提及，便不在此处赘述。

​		于是，现在的核心问题就变成了，如何在满足二叉搜索树的性质的情况下，使得它同时也满足完全二叉树的性质。

​		我们考虑这样一种情况，对于一个二叉搜索树来说，父结点的左结点的值必然是小于父结点同时也小于右结点的，对应的，父结点的右结点的值必然也是大于父结点同时也大于左结点的，所以，对于前一个情况来说，将父结点设为右结点的左结点，同时将父结点的右结点设为右节点的左节点，其大小情况是保持不变的；对应的，将父结点设为左结点的右结点，同时将父结点的左结点设为左节点的右节点，其大小情况仍然是满足二叉搜索树的条件的。

​		以上两种情况分别是**右旋和左旋**，到此时可以知道，可以通过左旋或右旋的方式，保持该结构仍然满足二叉搜索树的特征的同时，转换结点的位置。

​		同时我们可以知道，对于一个完全二叉树来说，即对于一个堆来说，它的核心特征是父结点的优先级必然大于左右两个结点，所以只需要通过旋转的方式满足该特征即可，即当父结点的左结点的优先级大于父结点的时候，进行左旋，当父结点的右节点的优先级大于父结点的时候进行右旋，当左右结点的优先级都大于父结点时，选择更大的一侧旋转即可。

​		到此，即可根据左右旋转的方式同时满足二叉搜索树和完全二叉树的性质了，只不过不同点在于，满足二叉搜索树的值是结点承载的值，而满足完全二叉树的优先级并非承载值，而是随机生成的值，同时由于其是随机生成的，从概率上来说，当数量足够大的情况下，可以实现O(logn)级别的增删查。

#### 添加策略

​		对于树堆的添加元素的情况来说，需要的情况主要可以分为该元素是否已经存在于树堆中，如果已经存在，则只需要将结点的num+1即可，否则则需要新增结点，同时赋予一个随机的优先级，再根据情况左右旋转以满足性质。

​		步骤：

1. 先以根节点为父结点，递归找到等于该元素的结点位置，若插入元素小于当前父结点元素，则转入左子树，否则转入右子树，以此方式找到对应的插入位置
2. 找到后，如果该元素已经存在，则结点num+1即可，否则新建结点，同时建立父结点和新结点之间的关系，此时，整个结构是满足于二叉搜索树的性质的
3. 随后，根据新增结点的优先级和其父结点的优先级进行判断，如果大于父结点优先级则进行对应方向的旋转即可。

#### 删除策略

​		对于树堆的元素删除来说，首先要做的还是先找到要删除的结点，如果没找到则可以之间结束。找到后需要判断该节点是否存在子结点，如果不存在则可以直接删除，否则，如果右左子树则找左子树的前缀节点即左子树的最右结点，如果只有右子树则朝气后缀结点即右子树的最右结点，将寻找到的结点替换为当前父结点，然后删除被替换的结点即可。但对于左右两个结点都存在的情况来说，则找到优先级更小的一侧进行旋转，然后将结点转入优先级更小的一侧结点进行持续递归，直到找到一个只有左节点或只有右节点或无结点的情况，对于存在结点的情况来说，把其子结点替换为自己即可，对于无结点的情况来说，删除即可。

​		树堆删除结点的步骤：

1. 通过递归比对遍历寻找到要删除的结点
2. 如果该结点不存在则直接结束，如果有重复则num-1即可结束
3. 如果同时有左右子树，转换为优先级更小的一侧，然后进行递归直到不同时存在左右子树为止
4. 如果有左子树，将左子树替换为自己即可
5. 如果只有右子树，将右子树替换为自己即可
6. 如果无子树，直接删除即可

### 实现

​		treap树堆结构体，该实例存储树堆的根节点，同时保存该树堆中已经存储了多少个元素，二叉树中排序使用的比较器在创建时传入,若不传入则在插入首个节点时从默认比较器中寻找，该树堆实例中存储随机数生成器,用于后续新建节点时生成随机数，创建时传入是否允许该树堆出现重复值,如果不允许则进行覆盖,允许则对节点数目增加即可。

```go
type treap struct {
	root    *node                 //根节点指针
	size    int                   //存储元素数量
	cmp     comparator.Comparator //比较器
	rand    *rand.Rand            //随机数生成器
	isMulti bool                  //是否允许重复
	mutex   sync.Mutex            //并发控制锁
}
```

​		node树节点结构体，该节点是树堆的树节点，若该树堆允许重复则对节点num+1即可,否则对value进行覆盖，树堆节点将针对堆的性质通过左右旋转的方式做平衡。

```go
type node struct {
	value    interface{} //节点中存储的元素
	priority uint32      //该节点的优先级,随机生成
	num      int         //该节点中存储的数量
	left     *node       //左节点指针
	right    *node       //右节点指针
}
```

#### 接口

```go
type treaper interface {
	Iterator() (i *Iterator.Iterator) //返回包含该树堆的所有元素,重复则返回多个
	Size() (num int)                  //返回该树堆中保存的元素个数
	Clear()                           //清空该树堆
	Empty() (b bool)                  //判断该树堆是否为空
	Insert(e interface{})             //向树堆中插入元素e
	Erase(e interface{})              //从树堆中删除元素e
	Count(e interface{}) (num int)    //从树堆中寻找元素e并返回其个数
}
```

#### New

​		新建一个treap树堆容器并返回，初始根节点为nil，传入该树堆是否为可重复属性,如果为true则保存重复值,否则对原有相等元素进行覆盖，若有传入的比较器,则将传入的第一个比较器设为该树堆的比较器。

```go
func New(isMulti bool, Cmp ...comparator.Comparator) (t *treap) {
	//设置默认比较器
	var cmp comparator.Comparator
	if len(Cmp) > 0 {
		cmp = Cmp[0]
	}
	//创建随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &treap{
		root:    nil,
		size:    0,
		cmp:     cmp,
		rand:    r,
		isMulti: isMulti,
		mutex:   sync.Mutex{},
	}
}
```

​		新建一个树堆节点并返回，将传入的元素e作为该节点的承载元素，该节点的num默认为1,左右子节点设为nil，该节点优先级随机生成,范围在0~2^16内。

```go
func newNode(e interface{}, rand *rand.Rand) (n *node) {
	return &node{
		value:    e,
		priority: uint32(rand.Intn(4294967295)),
		num:      1,
		left:     nil,
		right:    nil,
	}
}
```

#### Iterator

​		以treap树堆做接收者，将该树堆中所有保存的元素将从根节点开始以中缀序列的形式放入迭代器中，若允许重复存储则对于重复元素进行多次放入。

```go
func (t *treap) Iterator() (i *Iterator.Iterator) {
	if t == nil {
		return nil
	}
	t.mutex.Lock()
	es := t.root.inOrder()
	i = Iterator.New(&es)
	t.mutex.Unlock()
	return i
}
```

​		以node树堆节点做接收者，以中缀序列返回节点集合，若允许重复存储则对于重复元素进行多次放入。

```go
func (n *node) inOrder() (es []interface{}) {
	if n == nil {
		return es
	}
	if n.left != nil {
		es = append(es, n.left.inOrder()...)
	}
	for i := 0; i < n.num; i++ {
		es = append(es, n.value)
	}
	if n.right != nil {
		es = append(es, n.right.inOrder()...)
	}
	return es
}
```

#### Size

​		以treap树堆做接收者，返回该容器当前含有元素的数量，如果容器为nil返回0。

```go
func (t *treap) Size() (num int) {
	if t == nil {
		return 0
	}
	return t.size
}
```

#### Clear

​		以treap树堆做接收者，将该容器中所承载的元素清空，将该容器的size置0。

```go
func (t *treap) Clear() {
	if t == nil {
		return
	}
	t.mutex.Lock()
	t.root = nil
	t.size = 0
	t.mutex.Unlock()
}
```

#### Empty

​		以treap树堆做接收者，判断该二叉搜索树是否含有元素，如果含有元素则不为空,返回false，如果不含有元素则说明为空,返回true，如果容器不存在,返回true。

```go
func (t *treap) Empty() (b bool) {
	if t == nil {
		return true
	}
	if t.size > 0 {
		return false
	}
	return true
}
```

#### rightRotate

​		以node树堆节点做接收者，新建一个节点作为n节点的右节点,同时将n节点的数值放入新建节点中作为右转后的n节点，右转后的n节点的左节点是原n节点左节点的右节点,右转后的右节点保持不变，原n节点改为原n节点的左节点,同时右节点指向新建的节点即右转后的n节点，该右转方式可以保证n节点的双亲节点不用更换节点指向。

```go
func (n *node) rightRotate() {
   if n == nil {
      return
   }
   if n.left == nil {
      return
   }
   //新建节点作为更换后的n节点
   tmp := &node{
      value:    n.value,
      priority: n.priority,
      num:      n.num,
      left:     n.left.right,
      right:    n.right,
   }
   //原n节点左节点上移到n节点位置
   n.right = tmp
   n.value = n.left.value
   n.priority = n.left.priority
   n.num = n.left.num
   n.left = n.left.left
}
```

#### leftRotate

​		以node树堆节点做接收者，新建一个节点作为n节点的左节点,同时将n节点的数值放入新建节点中作为左转后的n节点，左转后的n节点的右节点是原n节点右节点的左节点,左转后的左节点保持不变，原n节点改为原n节点的右节点,同时左节点指向新建的节点即左转后的n节点，该左转方式可以保证n节点的双亲节点不用更换节点指向。

```go
func (n *node) leftRotate() {
   if n == nil {
      return
   }
   if n.right == nil {
      return
   }
   //新建节点作为更换后的n节点
   tmp := &node{
      value:    n.value,
      priority: n.priority,
      num:      n.num,
      left:     n.left,
      right:    n.right.left,
   }
   //原n节点右节点上移到n节点位置
   n.left = tmp
   n.value = n.right.value
   n.priority = n.right.priority
   n.num = n.right.num
   n.right = n.right.right
}
```

#### Insert

​		以treap树堆做接收者，向二叉树插入元素e,若不允许重复则对相等元素进行覆盖，如果二叉树为空则之间用根节点承载元素e,否则以根节点开始进行查找，对于该树堆来说,通过赋予随机的优先级根据堆的性质来实现平衡。

```go
func (t *treap) Insert(e interface{}) {
	//判断容器是否存在
	if t == nil {
		return
	}
	t.mutex.Lock()
	if t.Empty() {
		//判断比较器是否存在
		if t.cmp == nil {
			t.cmp = comparator.GetCmp(e)
		}
		if t.cmp == nil {
			t.mutex.Unlock()
			return
		}
		//插入到根节点
		t.root = newNode(e, t.rand)
		t.size = 1
		t.mutex.Unlock()
		return
	}
	//从根节点向下插入
	if t.root.insert(newNode(e, t.rand), t.isMulti, t.cmp) {
		t.size++
	}
	t.mutex.Unlock()
}
```

​		以node二叉搜索树节点做接收者，从n节点中插入元素e，如果n节点中承载元素与e不同则根据大小从左右子树插入该元素，如果n节点与该元素相等,且允许重复值,则将num+1否则对value进行覆盖，插入成功返回true,插入失败或不允许重复插入返回false。

```go
func (n *node) insert(e *node, isMulti bool, cmp comparator.Comparator) (b bool) {
	if cmp(n.value, e.value) > 0 {
		if n.left == nil {
			//将左节点直接设为e
			n.left = e
			b = true
		} else {
			//对左节点进行递归插入
			b = n.left.insert(e, isMulti, cmp)
		}
		if n.priority > e.priority {
			//对n节点进行右转
			n.rightRotate()
		}
		return b
	} else if cmp(n.value, e.value) < 0 {
		if n.right == nil {
			//将右节点直接设为e
			n.right = e
			b = true
		} else {
			//对右节点进行递归插入
			b = n.right.insert(e, isMulti, cmp)
		}
		if n.priority > e.priority {
			//对n节点进行左转
			n.leftRotate()
		}
		return b
	}
	if isMulti {
		//允许重复
		n.num++
		return true
	}
	//不允许重复,对值进行覆盖
	n.value = e.value
	return false
}
```

#### Erase

​		以treap树堆做接收者，从树堆中删除元素e，若允许重复记录则对承载元素e的节点中数量记录减一即可，若不允许重复记录则删除该节点同时将前缀节点或后继节点更换过来以保证树堆的不发送断裂，交换后根据优先级进行左右旋转以保证符合堆的性质，如果该树堆仅持有一个元素且根节点等价于待删除元素,则将根节点置为nil。

```go
func (t *treap) Erase(e interface{}) {
	if t == nil {
		return
	}
	if t.Empty() {
		return
	}
	t.mutex.Lock()
	if t.size == 1 && t.cmp(t.root.value, e) == 0 {
		//该树堆仅持有一个元素且根节点等价于待删除元素,则将根节点置为nil
		t.root = nil
		t.size = 0
		t.mutex.Unlock()
		return
	}
	//从根节点开始删除元素
	if t.root.delete(e, t.isMulti, t.cmp) {
		//删除成功
		t.size--
	}
	t.mutex.Unlock()
}
```

​		以node二叉搜索树节点做接收者，从n节点中删除元素e，如果n节点中承载元素与e不同则根据大小从左右子树删除该元素，如果n节点与该元素相等,且允许重复值,则将num-1否则直接删除该元素，删除时先寻找该元素的前缀节点,若不存在则寻找其后继节点进行替换，替换后删除该节点。

```go
func (n *node) delete(e interface{}, isMulti bool, cmp comparator.Comparator) (b bool) {
	if n == nil {
		return false
	}
	//n中承载元素小于e,从右子树继续删除
	if cmp(n.value, e) < 0 {
		if n.right == nil {
			//右子树为nil,删除终止
			return false
		}
		if cmp(e, n.right.value) == 0 && (!isMulti || n.right.num == 1) {
			//待删除节点无子节点,直接删除即可
			if n.right.left == nil && n.right.right == nil {
				//右子树可直接删除
				n.right = nil
				return true
			}
		}
		//从右子树继续删除
		return n.right.delete(e, isMulti, cmp)
	}
	//n中承载元素大于e,从左子树继续删除
	if cmp(n.value, e) > 0 {
		if n.left == nil {
			//左子树为nil,删除终止
			return false
		}
		if cmp(e, n.left.value) == 0 && (!isMulti || n.left.num == 1) {
			//待删除节点无子节点,直接删除即可
			if n.left.left == nil && n.left.right == nil {
				//左子树可直接删除
				n.left = nil
				return true
			}
		}
		//从左子树继续删除
		return n.left.delete(e, isMulti, cmp)
	}
	if isMulti && n.num > 1 {
		//允许重复且数量超过1
		n.num--
		return true
	}
	//删除该节点
	tmp := n
	//左右子节点都存在则选择优先级较小一个进行旋转
	for tmp.left != nil && tmp.right != nil {
		if tmp.left.priority < tmp.right.priority {
			tmp.rightRotate()
			if tmp.right.left == nil && tmp.right.right == nil {
				tmp.right = nil
				return false
			}
			tmp = tmp.right
		} else {
			tmp.leftRotate()
			if tmp.left.left == nil && tmp.left.right == nil {
				tmp.left = nil
				return false
			}
			tmp = tmp.left
		}
	}
	if tmp.left == nil && tmp.right != nil {
		//到左子树为nil时直接换为右子树即可
		tmp.value = tmp.right.value
		tmp.num = tmp.right.num
		tmp.priority = tmp.right.priority
		tmp.left = tmp.right.left
		tmp.right = tmp.right.right
	} else if tmp.right == nil && tmp.left != nil {
		//到右子树为nil时直接换为左子树即可
		tmp.value = tmp.left.value
		tmp.num = tmp.left.num
		tmp.priority = tmp.left.priority
		tmp.right = tmp.left.right
		tmp.left = tmp.left.left
	}
	//当左右子树都为nil时直接结束
	return true
}
```

#### Count

​		以treap树堆做接收者，从树堆中查找元素e的个数，如果找到则返回该树堆中和元素e相同元素的个数，如果不允许重复则最多返回1，如果未找到则返回0。

```go
func (t *treap) Count(e interface{}) (num int) {
	if t == nil {
		//树堆不存在,直接返回0
		return 0
	}
	if t.Empty() {
		return
	}
	t.mutex.Lock()
	num = t.root.search(e, t.cmp)
	t.mutex.Unlock()
	//树堆存在,从根节点开始查找该元素
	return num
}
```

​		以node二叉搜索树节点做接收者，从n节点中查找元素e并返回存储的个数，如果n节点中承载元素与e不同则根据大小从左右子树查找该元素，如果n节点与该元素相等,则直接返回其个数。

```go
func (n *node) search(e interface{}, cmp comparator.Comparator) (num int) {
	if n == nil {
		return 0
	}
	if cmp(n.value, e) > 0 {
		return n.left.search(e, cmp)
	} else if cmp(n.value, e) < 0 {
		return n.right.search(e, cmp)
	}
	return n.num
}
```

### 使用示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/treap"
	"sync"
)

func main() {
	bs := treap.New(true)
	wg := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		bs.Insert(i)
		go func() {
			bs.Insert(i)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("遍历输出所有插入的元素")
	for i := bs.Iterator(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("删除一次树堆中存在的元素,存在重复的将会被剩下")
	for i := 0; i < 20; i++ {
		bs.Erase(i)
	}
	fmt.Println("输出剩余的重复元素")
	for i := bs.Iterator(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
}

```

注：由于过程中的增删过程是并发执行的，所以其结果和下方示例并不完全相同

> 遍历输出所有插入的元素
> 0
> 1
> 2
> 3
> 4
> 5
> 6
> 7
> 8
> 8
> 9
> 10
> 11
> 11
> 11
> 11
> 11
> 11
> 11
> 11
> 11
> 11
> 11
> 12
> 12
> 13
> 14
> 15
> 16
> 16
> 16
> 16
> 16
> 17
> 17
> 18
> 19
> 20
> 20
> 20
> 删除一次树堆中存在的元素,存在重复的将会被剩下
> 输出剩余的重复元素
> 8
> 11
> 11
> 11
> 11
> 11
> 11
> 11
> 11
> 11
> 11
> 12
> 16
> 16
> 16
> 16
> 17
> 20
> 20
> 20

#### 时间开销

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/treap"
	"time"
)

func main() {
	max:=10000000
	tv := time.Now()
	v := make([]interface{},max,max)
	num:=0
	for i := 0; i < max; i++ {
		if i%2==1{
			num=max-i
		}else{
			num=i
		}
		v[num]=num
	}
	fmt.Println("插入切片的消耗时间:",time.Since(tv))
	tt := time.Now()
	t := treap.New(false)
	for i := 0; i < max; i++ {
		if i%2==1{
			num=max-i
		}else{
			num=i
		}
		t.Insert(num)
	}
	fmt.Println("插入树堆的消耗时间:",time.Since(tt))
}

```

> 插入切片的消耗时间: 201.4934ms
> 插入树堆的消耗时间: 7.0841462s

github仓库存储地址：https://github.com/hlccd/goSTL

### 概述

​		单词查找树（Tire），又叫前缀树，字典树，是一种有序多叉树。不同于之前实现的二叉搜，它是一个具有多个分叉的树的结构，同时，树结点和其子结点之间并无大小关系，只存在前缀关系，即其**父结点是其子结点的前缀**，一般用于存储string类型，对于string类型的增删查效率极高，其**增删查时间等价于string的长度**。

### 原理

​		本次实现的单词查找树的每个结点共有64个分叉，即‘a'~'z','A'~'Z','0'~'9','+','/'一共64个字符，对应base64的64个字符，可用于存储base64。

​		对于一个前缀树来说，它需要满足的特征有三条：

- 父结点的前缀必然是子结点的前缀
- 根节点不包含字符，除根节点以外每个节点只包含一个字符
- 每个节点的所有子节点包含的字符串不相同。

​		它的核心策略是以空间换时间，即将一个string类型拆开，分层保存其byte或叫char，使得每次增删查都只需要其长度的时间，最坏的查找时间比hash表更好。同时也不会出现冲突，并且也必然是满足字典序，即按其中序遍历得到的结果必然是有序的。

​		但同时，如果出现了一个较长的string，就会让整个链条变得很长，造成较多的空间开销。

#### 添加策略

​		从根节点开始插入，将string按byte进行分段，每层插入一个，当插入到最后时该string指向的value存在时则说明之前已经插入过，故插入失败，否则插入成功。

​		插入不允许覆盖。

​		当中间结点在原Trie树中不存在时创建即可。

#### 删除策略

​		从根节点开始删除，将string按byte进行分段，逐层往下遍历寻找到最终点，如果此时有存储的元素则删除同时表示删除成功，随后逐层返回将对应结点的num-1即可，当num=0时表示无后续结点，将该结点删除即可。如果在逐层下推的过程中发现结点不存在，可视为删除失败。之间返回即可。

### 实现

​		trie单词查找树结构体，该实例存储单词查找树的根节点，同时保存该树已经存储了多少个元素，整个树不允许重复插入,若出现重复插入则直接失败。

```go
type trie struct {
	root  *node      //根节点指针
	size  int        //存放的元素数量
	mutex sync.Mutex //并发控制锁
}
```

​		node树节点结构体，该节点是trie的树节点，结点存储到此时的string的前缀数量，以son为分叉存储下属的string，该节点同时存储其元素。

```go
type node struct {
	num   int         //以当前结点为前缀的string的数量
	son   [64]*node   //分叉
	value interface{} //当前结点承载的元素
}
```

#### 接口

```go
type trieer interface {
	Iterator() (i *Iterator.Iterator)        //返回包含该trie的所有string
	Size() (num int)                         //返回该trie中保存的元素个数
	Clear()                                  //清空该trie
	Empty() (b bool)                         //判断该trie是否为空
	Insert(s string, e interface{}) (b bool) //向trie中插入string并携带元素e
	Erase(s string) (b bool)                 //从trie中删除以s为索引的元素e
	Delete(s string) (num int)               //从trie中删除以s为前缀的所有元素
	Count(s string) (num int)                //从trie中寻找以s为前缀的string单词数
	Find(s string) (e interface{})           //从trie中寻找以s为索引的元素e
}
```

#### New

​		新建一个trie单词查找树容器并返回，初始根节点为nil。

```go
func New() (t *trie) {
	return &trie{
		root:  newNode(nil),
		size:  0,
		mutex: sync.Mutex{},
	}
}
```

​		新建一个单词查找树节点并返回，将传入的元素e作为该节点的承载元素。

```go
func newNode(e interface{}) (n *node) {
	return &node{
		num:   0,
		value: e,
	}
}
```

#### Iterator

​		以trie单词查找树做接收者，将该trie中所有存放的string放入迭代器中并返回。

```go
func (t *trie) Iterator() (i *Iterator.Iterator) {
	if t == nil {
		return nil
	}
	t.mutex.Lock()
	//找到trie中存在的所有string
	es := t.root.inOrder("")
	i = Iterator.New(&es)
	t.mutex.Unlock()
	return i
}
```

​		以node单词查找树节点做接收者，遍历其分叉以找到其存储的所有string。

```go
func (n *node) inOrder(s string) (es []interface{}) {
	if n == nil {
		return es
	}
	if n.value != nil {
		es = append(es, s)
	}
	for i, p := 0, 0; i < 62 && p < n.num; i++ {
		if n.son[i] != nil {
			if i < 26 {
				es = append(es, n.son[i].inOrder(s+string(i+'a'))...)
			} else if i < 52 {
				es = append(es, n.son[i].inOrder(s+string(i-26+'A'))...)
			} else {
				es = append(es, n.son[i].inOrder(s+string(i-52+'0'))...)
			}
			p++
		}
	}
	return es
}
```

#### Size

​		以trie单词查找树做接收者，返回该容器当前含有元素的数量，如果容器为nil返回0。

```go
func (t *trie) Size() (num int) {
	if t == nil {
		return 0
	}
	return t.size
}
```

#### Clear

​		以trie单词查找树做接收者，将该容器中所承载的元素清空，将该容器的size置0。

```go
func (t *trie) Clear() {
	if t == nil {
		return
	}
	t.mutex.Lock()
	t.root = newNode(nil)
	t.size = 0
	t.mutex.Unlock()
}
```

#### Empty

​		以trie单词查找树做接收者，判断该trie是否含有元素，如果含有元素则不为空,返回false，如果不含有元素则说明为空,返回true，如果容器不存在,返回true。

```go
func (t *trie) Empty() (b bool) {
	if t == nil {
		return true
	}
	return t.size == 0
}
```

##### getIdx

​		传入一个byte并根据其值返回其映射到分叉的值，当不属于'a'~'z','A'~'Z','0'~'9','+','/'时返回-1。

```go
func getIdx(c byte) (idx int) {
   if c >= 'a' && c <= 'z' {
      idx = int(c - 'a')
   } else if c >= 'A' && c <= 'Z' {
      idx = int(c-'A') + 26
   } else if c >= '0' && c <= '9' {
      idx = int(c-'0') + 52
   } else if c == '+' {
      idx = 62
   } else if c == '/' {
      idx = 63
   } else {
      idx = -1
   }
   return idx
}
```

#### Insert

​		以trie单词查找树做接收者，向trie插入以string类型的s为索引的元素e，若存在重复的s则插入失败,不允许覆盖，否则插入成功。

```go
func (t *trie) Insert(s string, e interface{}) (b bool) {
	if t == nil {
		return
	}
	if len(s) == 0 {
		return false
	}
	t.mutex.Lock()
	if t.root == nil {
		//避免根节点为nil
		t.root = newNode(nil)
	}
	//从根节点开始插入
	b = t.root.insert(s, 0, e)
	if b {
		//插入成功,size+1
		t.size++
	}
	t.mutex.Unlock()
	return b
}
```

​		以node单词查找树节点做接收者，从n节点中继续插入以s为索引的元素e,且当前抵达的string位置为p，当到达s终点时进行插入,如果此时node承载了元素则插入失败,否则成功，当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可。

```go
func (n *node) insert(s string, p int, e interface{}) (b bool) {
	if p == len(s) {
		if n.value != nil {
			return false
		}
		n.value = e
		n.num++
		return true
	}
	idx := getIdx(s[p])
	if idx == -1 {
		return false
	}
	if n.son[idx] == nil {
		n.son[idx] = newNode(nil)
	}
	b = n.son[idx].insert(s, p+1, e)
	if b {
		n.num++
	}
	return b
}
```

#### Erase

​		以trie单词查找树做接收者，从trie树中删除元素以s为索引的元素e。

```go
func (t *trie) Erase(s string) (b bool) {
	if t == nil {
		return false
	}
	if t.Empty() {
		return false
	}
	if len(s) == 0 {
		//长度为0无法删除
		return false
	}
	if t.root == nil {
		//根节点为nil即无法删除
		return false
	}
	t.mutex.Lock()
	//从根节点开始删除
	b = t.root.erase(s, 0)
	if b {
		//删除成功,size-1
		t.size--
		if t.size == 0 {
			//所有string都被删除,根节点置为nil
			t.root = nil
		}
	}
	t.mutex.Unlock()
	return b
}
```

​		以node单词查找树节点做接收者，从n节点中继续删除以s为索引的元素e,且当前抵达的string位置为p，当到达s终点时进行删除,如果此时node未承载元素则删除失败,否则成功，当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可,若其分叉为nil则直接失败。

```go
func (n *node) erase(s string, p int) (b bool) {
	if p == len(s) {
		if n.value != nil {
			n.value = nil
			n.num--
			return true
		}
		return false
	}
	idx := getIdx(s[p])
	if idx == -1 {
		return false
	}
	if n.son[idx] == nil {
		return false
	}
	b = n.son[idx].erase(s, p+1)
	if b {
		n.num--
		if n.son[idx].num == 0 {
			n.son[idx] = nil
		}
	}
	return b
}
```

#### Delete

​		以trie单词查找树做接收者，从trie树中删除以s为前缀的所有元素。

```go
func (t *trie) Delete(s string) (num int) {
   if t == nil {
      return 0
   }
   if t.Empty() {
      return 0
   }
   if len(s) == 0 {
      //长度为0无法删除
      return 0
   }
   if t.root == nil {
      //根节点为nil即无法删除
      return 0
   }
   t.mutex.Lock()
   //从根节点开始删除
   num = t.root.delete(s, 0)
   if num > 0 {
      //删除成功
      t.size -= num
      if t.size <= 0 {
         //所有string都被删除,根节点置为nil
         t.root = nil
      }
   }
   t.mutex.Unlock()
   return num
}
```

​		以node单词查找树节点做接收者，从n节点中继续删除以s为索引的元素e,且当前抵达的string位置为p，当到达s终点时进行删除,删除所有后续元素,并返回其后续元素的数量，当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可,若其分叉为nil则直接返回0。

```go
func (n *node) delete(s string, p int) (num int) {
   if p == len(s) {
      return n.num
   }
   idx := getIdx(s[p])
   if idx == -1 {
      return 0
   }
   if n.son[idx] == nil {
      return 0
   }
   num = n.son[idx].delete(s, p+1)
   if num>0 {
      n.num-=num
      if n.son[idx].num <= 0 {
         n.son[idx] = nil
      }
   }
   return num
}
```

#### Count

​		以trie单词查找树做接收者，从trie中查找以s为前缀的所有string的个数，如果存在以s为前缀的则返回大于0的值即其数量，如果未找到则返回0。

```go
func (t *trie) Count(s string) (num int) {
	if t == nil {
		return 0
	}
	if t.Empty() {
		return 0
	}
	if t.root == nil {
		return 0
	}
	t.mutex.Lock()
	//统计所有以s为前缀的string的数量并返回
	num = int(t.root.count(s, 0))
	t.mutex.Unlock()
	return num
}
```

​		以node单词查找树节点做接收者，从n节点中继续查找以s为前缀索引的元素e,且当前抵达的string位置为p，当到达s终点时返回其值即可，当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可,当其分叉为nil则直接返回0。

```go
func (n *node) count(s string, p int) (num int) {
	if p == len(s) {
		return n.num
	}
	idx := getIdx(s[p])
	if idx == -1 {
		return 0
	}
	if n.son[idx] == nil {
		return 0
	}
	return n.son[idx].count(s, p+1)
}
```

#### Find

​		以trie单词查找树做接收者，从trie中查找以s为索引的元素e,找到则返回e，如果未找到则返回nil。

```go
func (t *trie) Find(s string) (e interface{}) {
	if t == nil {
		return nil
	}
	if t.Empty() {
		return nil
	}
	if t.root == nil {
		return nil
	}
	t.mutex.Lock()
	//从根节点开始查找以s为索引的元素e
	e = t.root.find(s, 0)
	t.mutex.Unlock()
	return e
}
```

​		以node单词查找树节点做接收者，从n节点中继续查找以s为前缀索引的元素e,且当前抵达的string位置为p，当到达s终点时返回其承载的元素即可，当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可,当其分叉为nil则直接返回nil。

```go
func (n *node) find(s string, p int) (e interface{}) {
	if p == len(s) {
		return n.value
	}
	idx := getIdx(s[p])
	if idx == -1 {
		return nil
	}
	if n.son[idx] == nil {
		return nil
	}
	return n.son[idx].find(s, p+1)
}
```



### 使用示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/trie"
)

func main() {
	t:=trie.New()
	t.Insert("hlccd","hlccd")
	t.Insert("ha","ha")
	t.Insert("hb","hb")
	t.Insert("hc","hc")
	t.Insert("hd","hd")
	t.Insert("he","he")
	t.Insert("hl","hl")
	t.Insert("hlccd1","hlccd1")
	t.Insert("hlccd2","hlccd2")
	t.Insert("hlccd3","hlccd3")
	t.Insert("hlccd+","hlccd")
	t.Insert("hlccd/","hlccd")
	fmt.Println("当前插入的所有string:")
	for i:=t.Iterator().Begin();i.HasNext();i.Next(){
		fmt.Println(i.Value())
	}
	t.Erase("h")
	t.Erase("ha")
	t.Erase("hb")
	t.Erase("hc")
	t.Erase("hd")
	t.Erase("he")
	fmt.Println("定向删除后剩余的string:")
	for i:=t.Iterator().Begin();i.HasNext();i.Next(){
		fmt.Println(i.Value())
	}
	t.Delete("h")
	fmt.Println("删除以'h'为前缀的所有元素后剩余的数量:")
	for i:=t.Iterator().Begin();i.HasNext();i.Next(){
		fmt.Println(i.Value())
	}
}
```

> 当前插入的所有string:
> ha
> hb
> hc
> hd
> he
> hl
> hlccd
> hlccd1
> hlccd2
> hlccd3
> 定向删除后剩余的string:
> hl
> hlccd
> hlccd1
> hlccd2
> hlccd3
> 删除以'h'为前缀的所有元素后剩余的数量:

#### 时间开销

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/trie"
	"math/rand"
	"time"
)

func main() {
	max := 3000000
	ss := ""
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ts := time.Now()
	s := make([]string, 0, 0)
	for i := 0; i < max; i++ {
		ss = fmt.Sprintf("%d", r.Intn(4294967295))
		s = append(s, ss)
	}
	fmt.Println("slice消耗时间:", time.Since(ts))
	tm := time.Now()
	m := make(map[string]bool)
	for i := 0; i < max; i++ {
		ss = fmt.Sprintf("%d", r.Intn(4294967295))
		m[ss] = true
	}
	fmt.Println("map消耗时间:", time.Since(tm))
	tt := time.Now()
	t := trie.New()
	for i := 0; i < max; i++ {
		ss = fmt.Sprintf("%d", r.Intn(4294967295))
		t.Insert(ss, true)
	}
	fmt.Println("trie消耗时间:", time.Since(tt))
	tt1 := time.Now()
	t.Iterator()
	fmt.Println("trie遍历消耗的时间:", time.Since(tt1))
}
```

> slice消耗时间: 586.3899ms
> map消耗时间: 1.192838s
> trie消耗时间: 6.4676663s
> trie遍历消耗的时间: 4.3793102s

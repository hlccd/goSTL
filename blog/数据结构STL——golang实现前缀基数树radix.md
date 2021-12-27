github仓库存储地址：https://github.com/hlccd/goSTL

### 概述

​		前缀基数树（Radix），又叫基数树，是前缀树的一种变种。

​		它和前缀树不同的地方在于，它前缀树是将一个string按char进行分段保存，而基数树是将**多个char设为一层**，然后将string进行**分层保存**，一般利用**‘/’**作为分层标识。

​		它可用于string的存储和索引，当加上**模糊匹配**时也可用于实现**动态路由**。

### 原理

​		本次实现的前缀基数树的每个结点分别存储一层的string，终结点存储整个的string以作为终结点标识。

​		string按以‘/’作为分层标识，在一组‘/’之间的string为一层，若该层string为空忽略不记。

​		同时，每一层的string如果首字符为‘:'则视为动态匹配层，即该层可以匹配任意的string，而当该层首字符为'*'时也是动态匹配，且可匹配后方的所有内容。

​		为了加快查找速度，每一个结点的子结点利用map对子结点对应层的string进行一个映射。

​		对于一个基数树来说，它需要满足的特征有三条：

- 父结点的前缀必然是子结点的前缀
- 根节点不包含字符，除根节点以外每个节点包含一串char
- 终结点保存整个string
- 每一层的首字符为‘:'时为动态匹配，即任意的不可分层string皆可
- 每一层的首字符为‘*'时为动态匹配，可匹配后续所有的string，不论是否可以分蹭

#### 添加策略

​		从根节点开始插入，将string按’/‘进行分段，每层插入一个，当插入到最后时结点存储的string存在时则说明之前已经插入过，故插入失败，否则插入成功。

​		当中间结点在原Radix树中不存在时创建即可。

​		若插入失败则需要将原Radix树种不存在的结点删除并从map中删除。

#### 删除策略

​		从根节点开始删除，将string按'/'进行分段，逐层往下遍历寻找到最终点，如果此时有存储的元素则删除同时表示删除成功，随后逐层返回将对应结点的num-1即可，当num=0时表示无后续结点，将该结点删除即可。如果在逐层下推的过程中发现结点不存在，可视为删除失败。之间返回即可。

#### 匹配策略

​		匹配策略主要用于从radix中进行动态匹配，即从中找到一个存储的可以用于进行模糊匹配的string，然后将其和待匹配的s进行匹配，同时将模糊匹配层的name作为key，将待匹配层的string作为value放入map中并返回。

​		于是要解决问题就变成了如何找到一个可以对当前待匹配的string进行模糊匹配的存储结点的string。

- 将待匹配的s进行分层处理
- 从radix树的根节点开始匹配
  - 匹配时遍历每一个子结点，如果不支持动态匹配则判断同层是否相等，相等则加入可匹配的子结点数组
  - 如果支持动态匹配也加入可匹配的子结点数组
  - 继续匹配所有支持匹配的子结点直到找到第一个满足匹配情况的结点并返回
- 从找到的结点获取它存储的string然后进行分层
- 将分层结果中的首字符为’:'或‘*'的层即支持模糊匹配的层与待匹配的string的对应层建立映射关系
- 返回映射关系和匹配是否成功

### 实现

​		radix前缀基数树结构体，该实例存储前缀基数树的根节点，同时保存该树已经存储了多少个元素。

```go
type radix struct {
	root  *node      //前缀基数树的根节点指针
	size  int        //当前已存放的元素数量
	mutex sync.Mutex //并发控制锁
}
```

​		node树节点结构体，该节点是radix的树节点，结点存储到此时的string的前缀数量，son存储其下属分叉的子结点指针，该节点同时存储其元素。

```go
type node struct {
	pattern string           //到终点时不为"",其他都为""
	part    string           //以当前结点的string内容
	num     int              //以当前结点为前缀的数量
	sons    map[string]*node //该结点下属结点的指针
	fuzzy   bool             //模糊匹配?该结点首字符为':'或'*'为模糊匹配
}
```

#### 接口

```go
type radixer interface {
	Iterator() (i *Iterator.Iterator)             //返回包含该radix的所有string
	Size() (num int)                              //返回该radix中保存的元素个数
	Clear()                                       //清空该radix
	Empty() (b bool)                              //判断该radix是否为空
	Insert(s string) (b bool)                     //向radix中插入string
	Erase(s string) (b bool)                      //从radix中删除string
	Delete(s string) (num int)                    //从radix中删除以s为前缀的所有string
	Count(s string) (num int)                     //从radix中寻找以s为前缀的string单词数
	Mate(s string) (m map[string]string, ok bool) //利用radix树中的string对s进行模糊匹配,':'可模糊匹配该层,'*'可模糊匹配后面所有
}
```

#### New

​		新建一个radix前缀基数树容器并返回，初始根节点为nil。

```go
func New() (r *radix) {
	return &radix{
		root:  newNode(""),
		size:  0,
		mutex: sync.Mutex{},
	}
}
```

​		新建一个前缀基数树节点并返回，将传入的元素e作为该节点的承载元素。

```go
func newNode(part string) (n *node) {
	fuzzy := false
	if len(part) > 0 {
		fuzzy = part[0] == ':' || part[0] == '*'
	}
	return &node{
		pattern: "",
		part:    part,
		num:     0,
		sons:    make(map[string]*node),
		fuzzy:   fuzzy,
	}
}
```

##### analysis

​		将string按'/'进行分段解析，为""部分直接舍弃,返回解析结果，同时按规则重组用以解析是string并返回。

```go
func analysis(s string) (ss []string, newS string) {
   vs := strings.Split(s, "/")
   ss = make([]string, 0)
   newS = "/"
   for _, item := range vs {
      if item != "" {
         ss = append(ss, item)
         newS = newS + "/" + item
         if item[0] == '*' {
            break
         }
      }
   }
   return ss, newS
}
```

#### Iterator

​		以radix前缀基数树做接收者，将该radix中所有存放的string放入迭代器中并返回。

```go
func (r *radix) Iterator() (i *Iterator.Iterator) {
	if r == nil {
		return nil
	}
	r.mutex.Lock()
	es := r.root.inOrder("")
	i = Iterator.New(&es)
	r.mutex.Unlock()
	return i
}
```

​		以node前缀基数树节点做接收者，遍历其分叉以找到其存储的所有string。

```go
func (n *node) inOrder(s string) (es []interface{}) {
	if n == nil {
		return es
	}
	if n.pattern != "" {
		es = append(es, s+n.part)
	}
	for _, son := range n.sons {
		es = append(es, son.inOrder(s+n.part+"/")...)
	}
	return es
}
```

#### Size

​		以radix前缀基数树做接收者，返回该容器当前含有元素的数量，如果容器为nil返回0。

```go
func (r *radix) Size() (num int) {
	if r == nil {
		return 0
	}
	if r.root == nil {
		return 0
	}
	return r.size
}
```

#### Clear

​		以radix前缀基数树做接收者，将该容器中所承载的元素清空，将该容器的size置0。

```go
func (r *radix) Clear() {
	if r == nil {
		return
	}
	r.mutex.Lock()
	r.root = newNode("")
	r.size = 0
	r.mutex.Unlock()
}
```

#### Empty

​		以radix前缀基数树做接收者，判断该radix是否含有元素，如果含有元素则不为空,返回false，如果不含有元素则说明为空,返回true，如果容器不存在,返回true。

```go
func (r *radix) Empty() (b bool) {
	if r == nil {
		return true
	}
	return r.size == 0
}
```

#### Insert

​		以radix前缀基数树做接收者，向radix插入string，将对string进行解析,按'/'进行分层,':'为首则为模糊匹配该层,'*'为首则为模糊匹配后面所有，已经存在则无法重复插入。

```go
func (r *radix) Insert(s string) (b bool) {
	if r == nil {
		return false
	}
	//解析s并按规则重构s
	ss, s := analysis(s)
	r.mutex.Lock()
	if r.root == nil {
		//避免根节点为nil
		r.root = newNode("")
	}
	//从根节点开始插入
	b = r.root.insert(s, ss, 0)
	if b {
		//插入成功,size+1
		r.size++
	}
	r.mutex.Unlock()
	return b
}
```

​		以node前缀基数树节点做接收者，从n节点中继续插入以s为索引的元素e,且当前抵达的string位置为p，当到达s终点时进行插入,如果此时node承载了string则插入失败,否则成功，当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可，当插入失败且对应子结点为新建节点时则需要删除该子结点。

```go
func (n *node) insert(pattern string, ss []string, p int) (b bool) {
	if p == len(ss) {
		if n.pattern != "" {
			//该节点承载了string
			return false
		}
		//成功插入
		n.pattern = pattern
		n.num++
		return true
	}
	//找到该层的string
	s := ss[p]
	//从其子结点的map中找到对应的方向
	son, ok := n.sons[s]
	if !ok {
		//不存在,新建并放入map中
		son = newNode(s)
		n.sons[s] = son
	}
	//从子结点对应方向继续插入
	b = son.insert(pattern, ss, p+1)
	if b {
		n.num++
	} else {
		if !ok {
			//插入失败且该子节点为新建结点则需要删除该子结点
			delete(n.sons, s)
		}
	}
	return b
}
```

#### Erase

​		以radix前缀基数树做接收者，从radix树中删除元素string。

```go
func (r *radix) Erase(s string) (b bool) {
	if r.Empty() {
		return false
	}
	if len(s) == 0 {
		return false
	}
	if r.root == nil {
		//根节点为nil即无法删除
		return false
	}
	//解析s并按规则重构s
	ss, _ := analysis(s)
	r.mutex.Lock()
	//从根节点开始删除
	b = r.root.erase(ss, 0)
	if b {
		//删除成功,size-1
		r.size--
		if r.size == 0 {
			//所有string都被删除,根节点置为nil
			r.root = nil
		}
	}
	r.mutex.Unlock()
	return b
}
```

​		以node前缀基数树节点做接收者，从n节点中继续删除以s为索引的元素e,且当前抵达的string位置为p，当到达s终点时进行删除,如果此时node未承载元素则删除失败,否则成功，当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可,若其分叉为nil则直接失败。

```go
func (n *node) erase(ss []string, p int) (b bool) {
	if p == len(ss) {
		if n.pattern != "" {
			//该结点承载是string是,删除成功
			n.pattern = ""
			n.num--
			return true
		}
		return false
	}
	//从map中找到对应下子结点位置并递归进行删除
	s := ss[p]
	son, ok := n.sons[s]
	if !ok || son == nil {
		//未找到或son不存在,删除失败
		return false
	}
	b = son.erase(ss, p+1)
	if b {
		n.num--
		if son.num <= 0 {
			//删除后子结点的num<=0即该节点无后续存储元素,可以销毁
			delete(n.sons, s)
		}
	}
	return b
}
```

#### Delete

​		以radix前缀基数树做接收者，从radix树中删除以s为前缀的所有string。

```go
func (r *radix) Delete(s string) (num int) {
	if r.Empty() {
		return 0
	}
	if len(s) == 0 {
		return 0
	}
	if r.root == nil {
		return 0
	}
	//解析s并按规则重构s
	ss, _ := analysis(s)
	r.mutex.Lock()
	//从根节点开始删除
	num = r.root.delete(ss, 0)
	if num > 0 {
		//删除成功
		r.size -= num
		if r.size <= 0 {
			//所有string都被删除,根节点置为nil
			r.root = nil
		}
	}
	r.mutex.Unlock()
	return num
}
```

​		以node前缀基数树节点做接收者，从n节点中继续删除以s为索引的元素e,且当前抵达的string位置为p，当到达s终点时进行删除,删除所有后续元素,并返回其后续元素的数量，当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可,若其分叉为nil则直接返回0。

```go
func (n *node) delete(ss []string, p int) (num int) {
	if p == len(ss) {
		return n.num
	}
	//从map中找到对应下子结点位置并递归进行删除
	s := ss[p]
	son, ok := n.sons[s]
	if !ok || son == nil {
		return 0
	}
	num = son.delete(ss, p+1)
	if num > 0 {
		son.num -= num
		if son.num <= 0 {
			//删除后子结点的num<=0即该节点无后续存储元素,可以销毁
			delete(n.sons, s)
		}
	}
	return num
}
```

#### Count

​		以radix前缀基数树做接收者，从radix中查找以s为前缀的所有string的个数，如果存在以s为前缀的则返回大于0的值即其数量，如果未找到则返回0。

```go
func (r *radix) Count(s string) (num int) {
	if r.Empty() {
		return 0
	}
	if r.root == nil {
		return 0
	}
	if len(s) == 0 {
		return 0
	}
	//解析s并按规则重构s
	ss, _ := analysis(s)
	r.mutex.Lock()
	num = r.root.count(ss, 0)
	r.mutex.Unlock()
	return num
}
```

​		以node前缀基数树节点做接收者，从n节点中继续查找以s为前缀索引的元素e,且当前抵达的string位置为p，当到达s终点时返回其值即可，当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可,当其分叉为nil则直接返回0。

```go
func (n *node) count(ss []string, p int) (num int) {
	if p == len(ss) {
		return n.num
	}
	//从map中找到对应下子结点位置并递归进行查找
	s := ss[p]
	son, ok := n.sons[s]
	if !ok || son == nil {
		return 0
	}
	return son.count(ss, p+1)
}
```

#### Mate

​		以radix前缀基数树做接收者，从radix中查找以s为信息的第一个可以模糊匹配到的key和value的映射表，key是radix树中的段名,value是s中的段名，如果未找到则返回nil和false，否则返回一个映射表和true。

```go
func (r *radix) Mate(s string) (m map[string]string, ok bool) {
	if r.Empty() {
		return nil, false
	}
	if len(s) == 0 {
		return nil, false
	}
	if r.root == nil {
		return nil, false
	}
	//将s按'/'进行分割,并去掉第一个即去掉"",随后一次按照分层结果进行查找
	m, ok = r.root.mate(s, 0)
	return m, ok
}
```

​		以node前缀基数树节点做接收者，先从radix树的根节点开始找到第一个可以满足该模糊匹配方案的string结点，随后将s和结点的pattern进行模糊映射,将模糊查找的值和匹配值进行映射并返回即可，若该结点未找到则直接返回nil和false即可。

```go
func (n *node) mate(s string, p int) (m map[string]string, ok bool) {
	//解析url
	searchParts, _ := analysis(s)
	//从该请求类型中寻找对应的路由结点
	q := n.find(searchParts, 0)
	if q != nil {
		//解析该结点的pattern
		parts, _ := analysis(q.pattern)
		//动态参数映射表
		params := make(map[string]string)
		for index, part := range parts {
			if part[0] == ':' {
				//动态匹配,将参数名和参数内容的映射放入映射表内
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				//通配符,将后续所有内容全部添加到映射表内同时结束遍历
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return params, true
	}
	return nil, false
}
```

​		以node前缀基数树节点做接收者，从radix树的根节点开始找到第一个可以满足该模糊匹配方案的string结点，若该结点未找到则直接返回nil。

```go
func (n *node) find(parts []string, height int) (q *node) {
   //根据长度和局部string的首字符进行判断
   if len(parts) == height || strings.HasPrefix(n.part, "*") {
      if n.pattern == "" {
         //匹配失败,该结点处无匹配的信息
         return nil
      }
      //匹配成功,返回该结点
      return n
   }
   //从该结点的所有子结点中查找可用于递归查找的结点
   //当局部string信息和当前层string相同时可用于递归查找
   //当该子结点是动态匹配时也可以用于递归查找
   part := parts[height]
   //从所有子结点中找到可用于递归查找的结点
   children := make([]*node, 0, 0)
   for _, child := range n.sons {
      if child.part == part || child.fuzzy {
         //局部string相同或动态匹配
         children = append(children, child)
      }
   }
   for _, child := range children {
      //递归查询,并根据结果进行判断
      result := child.find(parts, height+1)
      if result != nil {
         //存在一个满足时就可以返回
         return result
      }
   }
   return nil
}
```

### 使用示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/radix"
)

func main() {
	radix := radix.New()
	radix.Insert("/test/:name/:key")
	radix.Insert("/hlccd/:name/:key")
	radix.Insert("/hlccd/1")
	radix.Insert("/hlccd/a/*name")
	fmt.Println("分层匹配")
	m, _ := radix.Mate("/hlccd/test/abc")
	for k, v := range m {
		fmt.Println(k, v)
	}
	fmt.Println("匹配全部")
	m, _ = radix.Mate("/hlccd/a/abc")
	for k, v := range m {
		fmt.Println(k, v)
	}
	fmt.Println("利用迭代器遍历")
	for i := radix.Iterator(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	radix.Erase("/hlccd/a/*name")
	fmt.Println("利用迭代器遍历定向删除后的结果")
	for i := radix.Iterator(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	radix.Delete("/hlccd/")
	fmt.Println("利用迭代器遍历删除前缀的结果")
	for i := radix.Iterator(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
}
```

> 分层匹配
> name test
> key abc
> 匹配全部
> name a
> key abc
> 利用迭代器遍历
> /test/:name/:key
> /hlccd/:name/:key
> /hlccd/1
> /hlccd/a/*name
> 利用迭代器遍历定向删除后的结果
> /test/:name/:key
> /hlccd/:name/:key
> /hlccd/1
> 利用迭代器遍历删除前缀的结果
> /test/:name/:key

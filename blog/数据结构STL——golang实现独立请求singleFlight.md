github仓库存储地址：https://github.com/hlccd/goSTL

### 概述

​		singleFlight，独立请求，用于做并发控制，常用于防止**缓存击穿**。

​		对于缓存来说，它一般会有一个过期时间，过期后进行删除，当在删除后的短时间内，如果突然出现了一大批对该数据的并发请求，次数他们都没有从缓存中读取到数据，然后集体涌入数据库中进行IO，引起数据库过载造成故障。

​		为了解决这个问题，可以通过给一组相同的请求添加一个**可重入锁**，即对于拥有同一个关键词的请求来说，可以视为一组相同的请求，此时，只允许其中一个进行请求，对其他进行阻塞操作，当请求结束后获取的数据放入缓存中，其他请求再从缓存中读取即可（**为了防止缓存击穿，可以考虑对返回的nil也进行存储**），这样也就可以避免一大堆相同请求都涌入数据库进行操作导致的压力过载。

### 原理

​		实现原理也是相对比较简单的，主要是利用**map**和**可重入锁**（golang中可利用WaitGroup）。

​		当一类拥有相同key的请求发起时，先向map中添加该类请求，随后将可重入锁加一，其他同类请求直接进入阻塞。

​		当请求结束后返回数据，同时将数据放入到请求的结构体内以供同类请求进行使用。

​		本次实现主要有两种形式，第一种是返回数据和错误信息，该方案可用于时间较短的情况，即持续阻塞直到数据返回的情况，但可能会出现由于长期等待导致的阻塞。第二种是直接返回一个channel，当数据获取成功后再向channel中放入数据即可，该方案可用于进行超时控制。

### 实现

​		呼叫请求结构体。

```go
type call struct {
	wg  sync.WaitGroup //可重入锁
	val interface{}    //请求结果
	err error          //错误反馈
}
```

​		一组请求工作，每一类请求对应一个call,利用其内部的可重入锁避免一类请求在短时间内频繁执行，请求组工作由使用者自行分配空间来实现。

```go
type Group struct {
	m  map[string]*call //一类请求与同一类呼叫的映射表
	mu sync.Mutex       //并发控制锁,保证线程安全
}
```

#### Do

​		以请求组做接收者，传入一个请求类的key和请求调用函数fn，请求时候需要等待之前有的同类请求先完成在进行，防止击穿缓存,对同一个key进行请求时需要分别进行,利用可重入锁实现，请求完成后返回结果和错误信息即可。

```go
func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	//判断以key为关键词的该类请求是否存在
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		// 如果请求正在进行中，则等待
		c.wg.Wait()
		return c.val, c.err
	}
	//该类请求不存在,创建个请求
	c := new(call)
	// 发起请求前加锁,并将请求添加到请求组内以表示该类请求正在处理
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()
	//调用请求函数获取内容
	c.val, c.err = fn()
	//请求结束
	c.wg.Done()
	g.mu.Lock()
	//从请求组中删除该呼叫请求
	delete(g.m, key)
	g.mu.Unlock()
	return c.val, c.err
}
```

#### DoChan

​		以请求组做接收者，传入一个请求类的key和请求调用函数fn，请求时候需要等待之前有的同类请求先完成在进行，防止击穿缓存,对同一个key进行请求时需要分别进行,利用可重入锁实现，返回一个channel,利用fn函数获取到的数据将会传入其中，可利用channel做超时控制。

```go
func (g *Group) DoChan(key string, fn func() (interface{}, error)) (ch chan interface{}) {
	ch = make(chan interface{}, 1)
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if _, ok := g.m[key]; ok {
		g.mu.Unlock()
		return ch
	}
	c := new(call)
	c.wg.Add(1)  // 发起请求前加锁
	g.m[key] = c // 添加到 g.m，表明 key 已经有对应的请求在处理
	g.mu.Unlock()
	go func() {
		c.val, c.err = fn() // 调用 fn，发起请求
		c.wg.Done()         // 请求结束
		g.mu.Lock()
		delete(g.m, key) // 更新 g.m
		ch <- c.val
		g.mu.Unlock()
	}()
	return ch
}
```

#### ForgetUnshared

​		以请求组做接收者，传入一个请求类的key，如果该key存在于请求组内,则将其删除即可，从而实现遗忘该类请求的目的。

```go
func (g *Group) ForgetUnshared(key string) {
	g.mu.Lock()
	_, ok := g.m[key]
	if ok {
		delete(g.m, key)
	}
	g.mu.Unlock()
}
```

### 使用示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/algorithm/singleFlight"
	"sync"
	"time"
)
var mu sync.Mutex
var num = 0

func get() (interface{}, error) {
	mu.Lock()
	num++
	e:=num
	mu.Unlock()
	return e, nil
}
func main() {
	wg := sync.WaitGroup{}
	sf := singleFlight.Group{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			fmt.Println(sf.Do("hlccd", get))
			wg.Done()
		}()
		if i == 5 {
			sf.ForgetUnshared("hlccd")
		}
	}
	wg.Wait()
	ch1 := sf.DoChan("hlccd", func() (interface{}, error) {
		time.Sleep(1 * time.Second)
		return 1, nil
	})
	ch2 := sf.DoChan("hlccd", func() (interface{}, error) {
		time.Sleep(1 * time.Second)
		return 2, nil
	})
	ch3 := sf.DoChan("hlccd", func() (interface{}, error) {
		time.Sleep(1 * time.Second)
		return 3, nil
	})
	select {
	case p := <-ch1:
		fmt.Println(p)
	case p := <-ch2:
		fmt.Println(p)
	case p := <-ch3:
		fmt.Println(p)
	case <-time.After(3*time.Second):
		fmt.Println("超时")
	}
}
```

注：该过程是个并发的随即情况

> 1 <nil>
> 2 <nil>
> 7 <nil>
> 2 <nil>
> 8 <nil>
> 3 <nil>
> 5 <nil>
> 4 <nil>
> 9 <nil>
> 1

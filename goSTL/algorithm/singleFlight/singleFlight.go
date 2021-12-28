package singleFlight

//@Title		singleFlight
//@Description
//		单次请求-single flight
//		利用可重入锁避免对于一个同类的请求进行多次从而导致的缓存击穿的问题
//		缓存击穿：
//		缓存在某个时间点过期的时候
//		恰好在这个时间点对这个Key有大量的并发请求过来
//		这些请求发现缓存过期一般都会从后端DB加载数据并回设到缓存
//		这个时候大并发的请求可能会瞬间把后端DB压垮。

import "sync"

//呼叫请求结构体
type call struct {
	wg  sync.WaitGroup //可重入锁
	val interface{}    //请求结果
	err error          //错误反馈
}

//一组请求工作
//每一类请求对应一个call,利用其内部的可重入锁避免一类请求在短时间内频繁执行
//请求组工作由使用者自行分配空间来实现
type Group struct {
	m  map[string]*call //一类请求与同一类呼叫的映射表
	mu sync.Mutex       //并发控制锁,保证线程安全
}

//@title    Do
//@description
//		以请求组做接收者
//		传入一个请求类的key和请求调用函数fn
//		请求时候需要等待之前有的同类请求先完成在进行
//		防止击穿缓存,对同一个key进行请求时需要分别进行,利用可重入锁实现
//		请求完成后返回结果和错误信息即可
//@receiver		g			*Group							接受者请求组的指针
//@param    	key			string							请求的关键词key
//@param    	fn			func() (interface{}, error)		请求执行函数
//@return    	v			interface{}						请求执行得到的结果
//@return    	err			error							执行请求后的错误信息
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

//@title    DoChan
//@description
//		以请求组做接收者
//		传入一个请求类的key和请求调用函数fn
//		请求时候需要等待之前有的同类请求先完成在进行
//		防止击穿缓存,对同一个key进行请求时需要分别进行,利用可重入锁实现
//		返回一个channel,利用fn函数获取到的数据将会传入其中
//		可利用channel做超时控制
//@receiver		g			*Group							接受者请求组的指针
//@param    	key			string							请求的关键词key
//@param    	fn			func() (interface{}, error)		请求执行函数
//@return    	ch			chan interface{}				执行结果将会传入其中,可做超时控制
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

//@title    ForgetUnshared
//@description
//		以请求组做接收者
//		传入一个请求类的key
//		如果该key存在于请求组内,则将其删除即可
//		从而实现遗忘该类请求的目的
//@receiver		g			*Group							接受者请求组的指针
//@param    	key			string							请求的关键词key
//@return    	nil
func (g *Group) ForgetUnshared(key string) {
	g.mu.Lock()
	_, ok := g.m[key]
	if ok {
		delete(g.m, key)
	}
	g.mu.Unlock()
}

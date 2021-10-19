package Iterator

//@Title		Iterator
//@Description
//		迭代器
//		定义了一套迭代器接口和迭代器类型
//		本套接口定义了迭代器所要执行的基本函数
//		数据结构在使用迭代器时需要重写函数
//		其中主要包括:生成迭代器,移动迭代器,判断是否可移动

//Iterator迭代器
//包含泛型切片和该迭代器当前指向元素的下标
//可通过下标和泛型切片长度来判断是否可以前移或后移
//当index不小于0时迭代器可前移
//当index小于data的长度时可后移
type Iterator struct {
	data  *[]interface{} //该迭代器中存放的元素集合的指针
	index int            //该迭代器当前指向的元素下标，-1即不存在元素
}

//Iterator迭代器接口
//定义了一套迭代器接口函数
//函数含义详情见下列描述
type Iteratorer interface {
	Begin() (I *Iterator)      //将该迭代器设为位于首节点并返回新迭代器
	End() (I *Iterator)        //将该迭代器设为位于尾节点并返回新迭代器
	Get(idx int) (I *Iterator) //将该迭代器设为位于第idx节点并返回该迭代器
	Value() (e interface{})    //返回该迭代器下标所指元素
	HasNext() (b bool)         //判断该迭代器是否可以后移
	Next() (b bool)            //将该迭代器后移一位
	HasPre() (b bool)          //判罚该迭代器是否可以前移
	Pre() (b bool)             //将该迭代器前移一位
}

//@title    New
//@description
//		新建一个Iterator迭代器容器并返回
//		传入的切片指针设为迭代器所承载的元素集合
//		若传入下标，则将传入的第一个下标设为该迭代器当前指向下标
//		若该下标超过元素集合范围，则寻找最近的下标
//		若元素集合为空，则下标设为-1
//@receiver		nil
//@param    	data		*[]interface{}	迭代器所承载的元素集合的指针
//@param    	Idx			...int			预设的迭代器的下标
//@return    	i        	*Iterator		新建的Iterator迭代器指针
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

//@title    Begin
//@description
//		以Iterator迭代器指针做接收者
//		如果迭代器为空，直接结束
//		如果该迭代器元素集合为空，则将下标设为-1
//		如果该迭代器元素集合不为空，则将下标设为0
//		随后返回新迭代器指针
//@receiver		i			*Iterator		迭代器指针
//@param    	nil
//@return    	I        	*Iterator		修改后的新迭代器指针
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

//@title    End
//@description
//		以Iterator迭代器指针做接收者
//		如果迭代器为空，直接结束
//		如果该迭代器元素集合为空，则将下标设为-1
//		如果该迭代器元素集合不为空，则将下标设为元素集合的最后一个元素的下标
//		随后返回新迭代器指针
//@receiver		i			*Iterator		迭代器指针
//@param    	nil
//@return    	I        	*Iterator		修改后的新迭代器指针
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

//@title    Get
//@description
//		以Iterator迭代器指针做接收者
//		如果迭代器为空，直接结束
//		如果该迭代器元素集合为空，则将下标设为-1
//		如果该迭代器元素集合不为空，则将下标设为传入的预设下标
//		如果预设下标超过元素集合范围，则将下标设为最近元素的下标
//		随后返回该迭代器指针
//@receiver		i			*Iterator		迭代器指针
//@param    	idx			int				预设下标
//@return    	I        	*Iterator		修改后的该迭代器指针
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

//@title    Value
//@description
//		以Iterator迭代器指针做接收者
//		返回迭代器当前下标所指元素
//		若迭代器为nil或元素集合为空，返回nil
//		否则返回迭代器当前下标所指向的元素
//		如果该下标超过元素集合范围，则返回距离最近的元素
//@receiver		i			*Iterator		迭代器指针
//@param    	nil
//@return    	e			interface{}		迭代器下标所指元素
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

//@title    HasNext
//@description
//		以Iterator迭代器指针做接收者
//		判断该迭代器是否可以后移
//		当迭代器为nil时不能后移
//		当元素集合为空时不能后移
//		当下标到达元素集合范围上限时不能后移
//		否则可以后移
//@author     	hlccd		2021-07-1
//@receiver		i			*Iterator		迭代器指针
//@param    	nil
//@return    	b			bool			迭代器下标可以后移？
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

//@title    Next
//@description
//		以Iterator迭代器指针做接收者
//		将迭代器下标后移
//		当满足后移条件时进行后移同时返回true
//		当不满足后移条件时将下标设为尾元素下标同时返回false
//		当迭代器为nil时返回false
//		当元素集合为空时下标设为-1同时返回false
//@receiver		i			*Iterator		迭代器指针
//@param    	nil
//@return    	b			bool			迭代器下标后移成功？
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

//@title    HasPre
//@description
//		以Iterator迭代器指针做接收者
//		判断该迭代器是否可以前移
//		当迭代器为nil时不能前移
//		当元素集合为空时不能前移
//		当下标到达元素集合范围下限时不能前移
//		否则可以前移
//@receiver		i			*Iterator		迭代器指针
//@param    	nil
//@return    	b			bool			迭代器下标可以前移？
func (i *Iterator) HasPre() (b bool) {
	if i == nil {
		//迭代器为nil时不能前移
		return false
	}
	if len((*i.data)) == 0 {
		//元素集合为空时不能前移
		return false
	}
	//下标到达元素集合范围下限时不能前移,否则可以后移
	return i.index >= 0
}

//@title    Pre
//@description
//		以Iterator迭代器指针做接收者
//		将迭代器下标前移
//		当满足前移条件时进行前移同时返回true
//		当不满足前移条件时将下标设为首元素下标同时返回false
//		当迭代器为nil时返回false
//		当元素集合为空时下标设为-1同时返回false
//@receiver		i			*Iterator		迭代器指针
//@param    	nil
//@return    	b			bool			迭代器下标前移成功？
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

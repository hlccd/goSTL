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
	data  *[]interface{} 	//该迭代器中存放的元素集合的指针
	index int           	//该迭代器当前指向的元素下标，-1即不存在元素
}

//Iterator迭代器接口
//定义了一套迭代器接口函数
//函数含义详情见下列描述
type Iteratorer interface {
	New(data *[]interface{}) (I *Iteratorer)	//传输元素创建一个新迭代器
	Begin() (I *Iterator)       				//将该迭代器设为位于首节点并返回新迭代器
	End() (I *Iterator)         				//将该迭代器设为位于尾节点并返回新迭代器
	Get(idx int) (I *Iterator)  				//将该迭代器设为位于第idx节点并返回该迭代器
	Value() (e interface{})    					//返回该迭代器下标所指元素
	HasNext() (b bool)      				    //判断该迭代器是否可以后移
	Next() (b bool)         				    //将该迭代器后移一位
	HasPre() (b bool)        					//判罚该迭代器是否可以前移
	Pre() (b bool)           					//将该迭代器前移一位
}

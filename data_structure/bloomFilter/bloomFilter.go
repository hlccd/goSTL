package bloomFilter

//@Title		bloomFilter
//@Description
//		bloomFilter布隆过滤器容器包
//		内部使用uint64切片进行存储
//		将任意类型的值进行hash计算后放入布隆过滤器中
//		可用于查找某一值是否已经插入过,但查找存在误差,只能确定其不存在,不能保证其必然存在
//		不能用于删除某一特定值,但可清空整个布隆过滤器

import "fmt"

//bloomFilter布隆过滤器结构体
//包含其用于存储的uint64元素切片
//选用uint64是为了更多的利用bit位
type bloomFilter struct {
	bits []uint64
	hash Hash
}

//bloomFilter布隆过滤器接口
//存放了bloomFilter布隆过滤器可使用的函数
//对应函数介绍见下方
type bloomFilteror interface {
	Insert(v interface{})         //向布隆过滤器中插入v
	Check(v interface{}) (b bool) //检查该值是否存在于布隆过滤器中,该校验存在误差
	Clear()                       //清空该布隆过滤器
}

//允许自行传入hash函数
type Hash func(v interface{}) (h uint32)

//@title    hash
//@description
//		传入一个虚拟节点id和实际结点
//		计算出它的hash值
//		逐层访问并利用素数131计算其hash值随后返回
//@receiver		nil
//@param    	v			interface{}		待计算的值
//@return    	h			uint32			计算得到的hash值
func hash(v interface{}) (h uint32) {
	h = uint32(0)
	s := fmt.Sprintf("131-%v-%v", v,v)
	bs := []byte(s)
	for i := range bs {
		h += uint32(bs[i]) * 131
	}
	return h
}

//@title    New
//@description
//		新建一个bloomFilter布隆过滤器容器并返回
//		初始bloomFilter的切片数组为空
//@receiver		nil
//@param    	h			Hash					hash函数
//@return    	bf        	*bloomFilter			新建的bloomFilter指针
func New(h Hash) (bf *bloomFilter) {
	if h == nil {
		h = hash
	}
	return &bloomFilter{
		bits: make([]uint64, 0, 0),
		hash: h,
	}
}

//@title    Insert
//@description
//		以bloomFilter布隆过滤器容器做接收者
//		先将待插入的value计算得到其哈希值hash
//		再向布隆过滤器中第hash位插入一个元素(下标从0开始)
//		当hash大于当前所能存储的位范围时,需要进行扩增
//		若要插入的位比冗余的多不足2^16即1024*64时,则新增1024个uint64
//		否则则直接增加到可以容纳第hash位的位置,以此可以提高冗余量,避免多次增加
//@receiver		bf        	*bloomFilter			接收者bloomFilter指针
//@param    	v			interface{}				待插入的值
//@return    	nil
func (bf *bloomFilter) Insert(v interface{}) {
	//bm不存在时直接结束
	if bf == nil {
		return
	}
	//开始插入
	h := bf.hash(v)
	if h/64+1 > uint32(len(bf.bits)) {
		//当前冗余量小于num位,需要扩增
		var tmp []uint64
		//通过冗余扩增减少扩增次数
		if h/64+1 < uint32(len(bf.bits)+1024) {
			//入的位比冗余的多不足2^16即1024*64时,则新增1024个uint64
			tmp = make([]uint64, len(bf.bits)+1024)
		} else {
			//直接增加到可以容纳第num位的位置
			tmp = make([]uint64, h/64+1)
		}
		//将原有元素复制到新增的切片内,并将bm所指向的修改为扩增后的
		copy(tmp, bf.bits)
		bf.bits = tmp
	}
	//将第num位设为1即实现插入
	bf.bits[h/64] ^= 1 << (h % 64)
}

//@title    Check
//@description
//		以bloomFilter布隆过滤器容器做接收者
//		将待查找的值做哈希计算得到哈希值h
//		检验第h位在位图中是否存在
//		当h大于当前所能存储的位范围时,直接返回false
//		否则判断第h为是否为1,为1返回true,否则返回false
//		利用布隆过滤器做判断存在误差,即返回true可能也不存在,但返回false则必然不存在
//@receiver		bf        	*bloomFilter			接收者bloomFilter指针
//@param    	v			interface{}				待查找的值
//@return    	b			bool					待查找的值可能存在于布隆过滤器中吗?
func (bf *bloomFilter) Check(v interface{}) (b bool) {
	//bf不存在时直接返回false并结束
	if bf == nil {
		return false
	}
	h := bf.hash(v)
	//h超出范围,直接返回false并结束
	if h/64+1 > uint32(len(bf.bits)) {
		return false
	}
	//判断第num是否为1,为1返回true,否则为false
	if bf.bits[h/64]&(1<<(h%64)) > 0 {
		return true
	}
	return false
}

//@title    Clear
//@description
//		以bloomFilter布隆过滤器容器做接收者
//		清空整个布隆过滤器
//@receiver		bf        	*bloomFilter			接收者bloomFilter指针
//@param    	nil
//@return    	nums		[]uint					所有在位图中存在的元素的下标集合
func (bf *bloomFilter) Clear() {
	if bf == nil {
		return
	}
	bf.bits = make([]uint64, 0, 0)
}

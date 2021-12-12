package bitmap

//@Title		bitmap
//@Description
//		bitmap位图容器包
//		内部使用uint64切片进行存储
//		由于数字在计算机内部存储时采用多个bit组成一个字符
//		而一bit只有1和0两个情况,所以也可以使用一个bit表示任意一位存在
//		该数据结构主要可以进行过滤去重、标注是否存在、快速排序的功能

//bitmap位图结构体
//包含其用于存储的uint64元素切片
//选用uint64是为了更多的利用bit位

type Bitmap struct {
	bits []uint64
}

//bitmap位图容器接口
//存放了bitmap容器可使用的函数
//对应函数介绍见下方

type bitmaper interface {
	Insert(num uint)         //在num位插入元素
	Delete(num uint)         //删除第num位
	Check(num uint) (b bool) //检查第num位是否有元素
	All() (nums []uint)      //返回所有存储的元素的下标
	Clear()                  //清空
}

//@title    New
//@description
//		新建一个bitmap位图容器并返回
//		初始bitmap的切片数组为空
//@receiver		nil
//@param    	nil
//@return    	bm        	*Bitmap					新建的bitmap指针
func New() (bm *Bitmap) {
	return &Bitmap{
		bits: make([]uint64, 0, 0),
	}
}

//@title    Insert
//@description
//		以bitmap位图容器做接收者
//		向位图中第num位插入一个元素(下标从0开始)
//		当num大于当前所能存储的位范围时,需要进行扩增
//		若要插入的位比冗余的多不足2^16即1024*64时,则新增1024个uint64
//		否则则直接增加到可以容纳第num位的位置,以此可以提高冗余量,避免多次增加
//@receiver		bm			*Bitmap					接受者bitmap的指针
//@param    	num			int						待插入的位的下标
//@return    	nil
func (bm *Bitmap) Insert(num uint) {
	//bm不存在时直接结束
	if bm == nil {
		return
	}
	//开始插入
	if num/64+1 > uint(len(bm.bits)) {
		//当前冗余量小于num位,需要扩增
		var tmp []uint64
		//通过冗余扩增减少扩增次数
		if num/64+1 < uint(len(bm.bits)+1024) {
			//入的位比冗余的多不足2^16即1024*64时,则新增1024个uint64
			tmp = make([]uint64, len(bm.bits)+1024)
		} else {
			//直接增加到可以容纳第num位的位置
			tmp = make([]uint64, num/64+1)
		}
		//将原有元素复制到新增的切片内,并将bm所指向的修改为扩增后的
		copy(tmp, bm.bits)
		bm.bits = tmp
	}
	//将第num位设为1即实现插入
	bm.bits[num/64] ^= 1 << (num % 64)
}

//@title    Delete
//@description
//		以bitmap位图容器做接收者
//		向位图中第num位删除一个元素(下标从0开始)
//		当num大于当前所能存储的位范围时,直接结束即可
//		删除完成后对切片最后存储的uint64进行判断是否大于1,若大于1则不做缩容
//		若等于0则可以进行缩容
//		对于缩容而言,从后往前遍历判断最后有多少个连续的0,即可以删除的多少组
//		若可删除的组大于总组数的一半则进行删除,否则则当作冗余量即可
//		若可删除的组数超过1024个时,则先删除1024个
//@receiver		bm			*Bitmap					接受者bitmap的指针
//@param    	num			int						待删除的位的下标
//@return    	nil
func (bm *Bitmap) Delete(num uint) {
	//bm不存在时直接结束
	if bm == nil {
		return
	}
	//num超出范围,直接结束
	if num/64+1 > uint(len(bm.bits)) {
		return
	}
	//将第num位设为0
	bm.bits[num/64] &^= 1 << (num % 64)
	if bm.bits[len(bm.bits)-1] == 0 {
		//最后一组为0,可能进行缩容
		//从后往前遍历判断可缩容内容是否小于总组数
		i := len(bm.bits) - 1
		for ; i >= 0; i-- {
			if bm.bits[i] == 0  && i!=len(bm.bits)-256{
				continue
			} else {
				//不为0或到1024个时即可返回
				break
			}
		}
		if i <= len(bm.bits)/2 || i==len(bm.bits)-256 {
			//小于总组数一半或超过1023个,进行缩容
			bm.bits = bm.bits[:i+1]
		}
	} else {
		return
	}
}

//@title    Check
//@description
//		以bitmap位图容器做接收者
//		检验第num位在位图中是否存在
//		当num大于当前所能存储的位范围时,直接返回false
//		否则判断第num为是否为1,为1返回true,否则返回false
//@receiver		bm			*Bitmap					接受者bitmap的指针
//@param    	num			int						待检测位的下标
//@return    	b			bool					第num位存在于位图中
func (bm *Bitmap) Check(num uint) (b bool) {
	//bm不存在时直接返回false并结束
	if bm == nil {
		return false
	}
	//num超出范围,直接返回false并结束
	if num/64+1 > uint(len(bm.bits)) {
		return false
	}
	//判断第num是否为1,为1返回true,否则为false
	if bm.bits[num/64]&(1<<num%64) > 0 {
		return true
	}
	return false
}

//@title    All
//@description
//		以bitmap位图容器做接收者
//		返回所有在位图中存在的元素的下标
//		返回的下标是单调递增序列
//@receiver		bm			*Bitmap					接受者bitmap的指针
//@param    	nil
//@return    	nums		[]uint					所有在位图中存在的元素的下标集合
func (bm *Bitmap) All() (nums []uint) {
	//对要返回的集合进行初始化,以避免返回nil
	nums=make([]uint,0,0)
	//bm不存在时直接返回并结束
	if bm == nil {
		return nums
	}
	//分组遍历判断某下标的元素是否存在于位图中,即其值是否为1
	for j := 0; j < len(bm.bits); j++ {
		for i := 0; i < 64; i++ {
			if bm.bits[j]&(1<<i) > 0 {
				//该元素存在,添加入结果集合内
				nums = append(nums, uint(j*64+i))
			}
		}
	}
	return nums
}

//@title    Clear
//@description
//		以bitmap位图容器做接收者
//		清空位图
//@receiver		bm			*Bitmap					接受者bitmap的指针
//@param    	nil
//@return    	nums		[]uint					所有在位图中存在的元素的下标集合
func (bm *Bitmap) Clear() {
	if bm == nil {
		return
	}
	bm.bits = make([]uint64, 0, 0)
}

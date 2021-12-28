package comparator

//@Title		comparator
//@Description
//		该包内通过利用比较函数重排传入的数组,使得下标为n的元素必然是第n+1大的(考虑到下标从0开始)
//		对二分排序的变形,当只对该节点位置存在的某一局部进行查找即可

//@title    NthElement
//@description
//		若数组指针为nil或者数组为nil或数组长度为0则直接结束即可
//		通过利用比较函数重排传入的数组,使得下标为n的元素必然是第n+1大的(考虑到下标从0开始)
//		若n大于数组长度,直接结束,否则对第n+1大的元素排序后返回该元素
//@receiver		nil
//@param    	begin		*[]interface{}				待查找的元素数组指针
//@param    	n			int							待查找的是第n位,从0计数
//@param    	Cmp			...Comparator				比较函数
//@return    	value		interface{}					第n+1大的元素
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

//@title    nthElement
//@description
//		对传入的开启和结尾的两个比较器中的值进行查找
//		以传入的比较器进行比较
//		通过局部二分的方式进行查找并将第n小的元素放到第n位置(大小按比较器进行确认,默认未小)
//@receiver		nil
//@param    	begin		*[]interface{}				待查找的元素数组指针
//@param    	l			int							查找范围的左下标
//@param    	r			int							查找范围的右下标
//@param    	n			int							待查找的是第n位,从0计数
//@param    	Cmp			...Comparator				比较函数
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

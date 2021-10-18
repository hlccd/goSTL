package comparator

//@Title		comparator
//@Description
//		该部分为对带查找数组中元素进行二分查找
//		warning:仅对有序元素集合有效

//@title    Search
//@description
//		若数组指针为nil或者数组为nil或数组长度为0则直接结束即可
//		通过比较函数对传入的数组中的元素集合进行二分查找
//		若并未传入比较函数则寻找默认比较函数
//		找到后返回该元素的下标
//		若该元素不在该部分内存在,则返回-1
//@receiver		nil
//@param    	arr			*[]interface{}				待查找的有序数组
//@param    	e			interface{}					待查找元素
//@param    	Cmp			...Comparator				比较函数
//@return    	idx			int							待查找元素下标
func Search(arr *[]interface{}, e interface{}, Cmp ...Comparator) (idx int) {
	if arr==nil || (*arr)==nil || len((*arr)) == 0 {
		return
	}
	//判断比较函数是否有效,若无效则寻找默认比较函数
	var cmp Comparator
	cmp = nil
	if len(Cmp) == 0 {
		cmp = GetCmp(e)
	} else {
		cmp = Cmp[0]
	}
	if cmp == nil {
		//若并非默认类型且未传入比较器则直接结束
		return -1
	}
	//查找开始
	return search(arr, e, cmp)
}

//@title    search
//@description
//		通过比较函数对传入的数组中的元素集合进行二分查找
//		找到后返回该元素的下标
//		若该元素不在该部分内存在,则返回-1
//@receiver		nil
//@param    	arr			*[]interface{}				待查找的有序数组
//@param    	e			interface{}					待查找元素
//@param    	Cmp			...Comparator				比较函数
//@return    	idx			int							待查找元素下标
func search(arr *[]interface{}, e interface{}, cmp Comparator) (idx int) {
	//通过二分查找的方式寻找该元素
	l, m, r := 0, (len((*arr))-1)/2, len((*arr))
	for l < r {
		m = (l + r) / 2
		if cmp((*arr)[m], e) < 0 {
			l = m + 1
		} else {
			r = m
		}
	}
	//查找结束
	if (*arr)[l] == e {
		//该元素存在,返回下标
		return l
	}
	//该元素不存在,返回-1
	return -1
}

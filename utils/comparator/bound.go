package comparator

//@Title		comparator
//@Description
//		可查找有序序列中某一元素的上界和下届
//		当该元素存在时,上界和下届返回的下标位置指向该元素
//		当该元素不存在时,上界和下届指向下标错位且指向位置并非该元素

//@title    UpperBound
//@description
//		通过传入的比较函数对待查找数组进行查找以获取待查找元素的上界即不大于它的最大值的下标
//		以传入的比较函数进行比较
//		如果该元素存在,则上界指向元素为该元素
//		如果该元素不存在,上界指向元素为该元素的前一个元素
//@receiver		nil
//@param    	arr			*[]interface{}				待查找数组
//@param    	e			interface{}					待查找元素
//@param    	Cmp			...Comparator				比较函数
//@return    	idx 		int							待查找元素的上界
func UpperBound(arr *[]interface{}, e interface{}, Cmp ...Comparator) (idx int) {
	if arr==nil || (*arr)==nil || len((*arr)) == 0 {
		return -1
	}
	//判断比较函数是否有效
	var cmp Comparator
	cmp = nil
	if len(Cmp) == 0 {
		cmp = GetCmp(e)
	} else {
		cmp = Cmp[0]
	}
	if cmp == nil {
		return -1
	}
	//寻找该元素的上界
	return upperBound(arr, e, cmp)
}

//@title    upperBound
//@description
//		通过传入的比较函数对待查找数组进行查找以获取待查找元素的上界即不大于它的最大值的下标
//		以传入的比较函数进行比较
//		如果该元素存在,则上界指向元素为该元素,且为最右侧
//		如果该元素不存在,上界指向元素为该元素的前一个元素
//		以二分查找的方式寻找该元素的上界
//@receiver		nil
//@param    	arr			*[]interface{}				待查找数组
//@param    	e			interface{}					待查找元素
//@param    	Cmp			...Comparator				比较函数
//@return    	idx 		int							待查找元素的上界
func upperBound(arr *[]interface{}, e interface{}, cmp Comparator) (idx int) {
	l, m, r := 0, len((*arr)) / 2, len((*arr))-1
	for l < r {
		m = (l + r + 1) / 2
		if cmp((*arr)[m], e) <= 0 {
			l = m
		} else {
			r = m - 1
		}
	}
	return l
}

//@title    LowerBound
//@description
//		通过传入的比较函数对待查找数组进行查找以获取待查找元素的下界即不小于它的最小值的下标
//		以传入的比较函数进行比较
//		如果该元素存在,则上界指向元素为该元素
//		如果该元素不存在,上界指向元素为该元素的后一个元素
//@receiver		nil
//@param    	arr			*[]interface{}				待查找数组
//@param    	e			interface{}					待查找元素
//@param    	Cmp			...Comparator				比较函数
//@return    	idx 		int							待查找元素的下界
func LowerBound(arr *[]interface{}, e interface{}, Cmp ...Comparator) (idx int) {
	if arr==nil || (*arr)==nil || len((*arr)) == 0 {
		return -1
	}
	//判断比较函数是否有效
	var cmp Comparator
	cmp = nil
	if len(Cmp) == 0 {
		cmp = GetCmp(e)
	} else {
		cmp = Cmp[0]
	}
	if cmp == nil {
		return -1
	}
	//寻找该元素的下界
	return lowerBound(arr, e, cmp)
}

//@title    lowerBound
//@description
//		通过传入的比较函数对待查找数组进行查找以获取待查找元素的下界即不小于它的最小值的下标
//		以传入的比较函数进行比较
//		如果该元素存在,则上界指向元素为该元素,且为最右侧
//		如果该元素不存在,上界指向元素为该元素的后一个元素
//@receiver		nil
//@param    	arr			*[]interface{}				待查找数组
//@param    	e			interface{}					待查找元素
//@param    	Cmp			...Comparator				比较函数
//@return    	idx 		int							待查找元素的下界
func lowerBound(arr *[]interface{}, e interface{}, cmp Comparator) (idx int) {
	l, m, r := 0, len((*arr)) / 2, len((*arr))
	for l < r {
		m = (l + r) / 2
		if cmp((*arr)[m], e) >= 0 {
			r = m
		} else {
			l = m + 1
		}
	}
	return l
}

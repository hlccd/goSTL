package comparator

//@Title		comparator
//@Description
//		比较器
//		定义了一个比较器类型,该类型可传入两个泛型并返回一个整数判断大小
//		该比较器用于定义的数据结构中传入数据之间的比较
//		该包内定义了一些自带类型的比较类型
//		当使用自定义的数据结构时若不进行比较器的传入则使用默认比较器
//		若传入类型非系统自带类型,则返回空比较器同时对数据的传入失败

// 比较器将会返回数字num
// num > 0 ,if a > b
// num = 0 ,if a = b
// num < 0 ,if a < b

type Comparator func(a, b interface{}) int

//比较器的特种——相等器
//判断传入的两个元素是否相等
type Equaler func(a, b interface{}) (B bool)

//@title    GetCmp
//@description
//		传入一个数据并根据该数据类型返回一个对应的比较器
//		若该类型并非系统自带类型,则返回个空比较器
//		若传入元素为nil则之间返回nil
//@receiver		nil
//@param    	e			interface{}
//@return    	cmp        	Comparator		该类型对应的默认比较器
func GetCmp(e interface{}) (cmp Comparator) {
	if e == nil {
		return nil
	}
	switch e.(type) {
	case bool:
		return boolCmp
	case int:
		return intCmp
	case int8:
		return int8Cmp
	case uint8:
		return uint8Cmp
	case int16:
		return int16Cmp
	case uint16:
		return uint16Cmp
	case int32:
		return int32Cmp
	case uint32:
		return uint32Cmp
	case int64:
		return int64Cmp
	case uint64:
		return uint64Cmp
	case float32:
		return float32Cmp
	case float64:
		return float64Cmp
	case complex64:
		return complex64Cmp
	case complex128:
		return complex128Cmp
	case string:
		return stringCmp
	}
	return nil
}

//@title    GetEqual
//@description
//		传入一个数据并根据该数据类型返回一个对应的比较器
//		若该类型并非系统自带类型,则返回个空比较器
//		若传入元素为nil则之间返回nil
//@receiver		nil
//@param    	e			interface{}
//@return    	cmp        	Comparator		该类型对应的默认比较器
func GetEqual() (equ Equaler) {
	return basicEqual
}

//@title    basicEqual
//@description
//		返回基本比较器
//		即有且仅有判断量元素是否完全相等
//@receiver		a			interface{}		待判断相等的第一个元素
//@receiver		b			interface{}		待判断相等的第二个元素
//@param    	nil
//@return    	B			bool			这两个元素是否相等？
func basicEqual(a, b interface{}) (B bool) {
	return a == b
}

//以下为系统自带类型的默认比较器

func boolCmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(bool) {
		return 1
	} else if b.(bool) {
		return -1
	}
	return 0
}
func intCmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(int) > b.(int) {
		return 1
	} else if a.(int) < b.(int) {
		return -1
	}
	return 0
}
func int8Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(int8) > b.(int8) {
		return 1
	} else if a.(int8) < b.(int8) {
		return -1
	}
	return 0
}
func uint8Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(uint8) > b.(uint8) {
		return 1
	} else if a.(uint8) < b.(uint8) {
		return -1
	}
	return 0
}
func int16Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(int16) > b.(int16) {
		return 1
	} else if a.(int16) < b.(int16) {
		return -1
	}
	return 0
}
func uint16Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(uint16) > b.(uint16) {
		return 1
	} else if a.(uint16) < b.(uint16) {
		return -1
	}
	return 0
}
func int32Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(int32) > b.(int32) {
		return 1
	} else if a.(int32) < b.(int32) {
		return -1
	}
	return 0
}
func uint32Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(uint32) > b.(uint32) {
		return 1
	} else if a.(uint32) < b.(uint32) {
		return -1
	}
	return 0
}
func int64Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(int64) > b.(int64) {
		return 1
	} else if a.(int64) < b.(int64) {
		return -1
	}
	return 0
}
func uint64Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(uint64) > b.(uint64) {
		return 1
	} else if a.(uint64) < b.(uint64) {
		return -1
	}
	return 0
}
func float32Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(float32) > b.(float32) {
		return 1
	} else if a.(float32) < b.(float32) {
		return -1
	}
	return 0
}
func float64Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(float64) > b.(float64) {
		return 1
	} else if a.(float64) < b.(float64) {
		return -1
	}
	return 0
}
func complex64Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if real(a.(complex64)) > real(b.(complex64)) {
		return 1
	} else if real(a.(complex64)) < real(b.(complex64)) {
		return -1
	} else {
		if imag(a.(complex64)) > imag(b.(complex64)) {
			return 1
		} else if imag(a.(complex64)) < imag(b.(complex64)) {
			return -1
		}
	}
	return 0
}
func complex128Cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if real(a.(complex128)) > real(b.(complex128)) {
		return 1
	} else if real(a.(complex128)) < real(b.(complex128)) {
		return -1
	} else {
		if imag(a.(complex128)) > imag(b.(complex128)) {
			return 1
		} else if imag(a.(complex128)) < imag(b.(complex128)) {
			return -1
		}
	}
	return 0
}
func stringCmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if len(a.(string)) > len(b.(string)) {
		return 1
	} else if len(a.(string)) < len(b.(string)) {
		return -1
	} else {
		if a.(string) > b.(string) {
			return 1
		} else if a.(string) < b.(string) {
			return -1
		}
	}
	return 0
}

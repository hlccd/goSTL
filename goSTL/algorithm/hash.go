package algorithm

type Hasher func(key interface{}) uint64

func GetHash(e interface{}) (hash Hasher) {
	if e == nil {
		return nil
	}
	switch e.(type) {
	case bool:
		return boolHash
	case int:
		return intHash
	case int8:
		return int8Hash
	case uint8:
		return uint8Hash
	case int16:
		return int16Hash
	case uint16:
		return uint16Hash
	case int32:
		return int32Hash
	case uint32:
		return uint32Hash
	case int64:
		return int64Hash
	case uint64:
		return uint64Hash
	case float32:
		return float32Hash
	case float64:
		return float64Hash
	case complex64:
		return complex64Hash
	case complex128:
		return complex128Hash
	case string:
		return stringHash
	}
	return nil
}
func boolHash(key interface{}) uint64 {
	if key.(bool) {
		return 1
	}
	return 0
}
func intHash(key interface{}) uint64 {
	return uint64(key.(int) * key.(int)/2)
}
func int8Hash(key interface{}) uint64 {
	return uint64(key.(int8) * key.(int8)/2)
}
func uint8Hash(key interface{}) uint64 {
	return uint64(key.(uint8) * key.(uint8)/2)
}
func int16Hash(key interface{}) uint64 {
	return uint64(key.(int16) * key.(int16)/2)
}
func uint16Hash(key interface{}) uint64 {
	return uint64(key.(uint16) * key.(uint16)/2)
}
func int32Hash(key interface{}) uint64 {
	return uint64(key.(int32) * key.(int32)/2)
}
func uint32Hash(key interface{}) uint64 {
	return uint64(key.(uint32) * key.(uint32)/2)
}
func int64Hash(key interface{}) uint64 {
	return uint64(key.(int64) * key.(int64)/2)
}
func uint64Hash(key interface{}) uint64 {
	return uint64(key.(uint64) * key.(uint64)/2)
}
func float32Hash(key interface{}) uint64 {
	return uint64(key.(float32) * key.(float32)/2)
}
func float64Hash(key interface{}) uint64 {
	return uint64(key.(float64) * key.(float64)/2)
}
func complex64Hash(key interface{}) uint64 {
	r := uint64(real(key.(complex64)))
	i := uint64(imag(key.(complex64)))
	return uint64Hash(r) + uint64Hash(i)
}
func complex128Hash(key interface{}) uint64 {
	r := uint64(real(key.(complex64)))
	i := uint64(imag(key.(complex64)))
	return uint64Hash(r) + uint64Hash(i)
}
func stringHash(key interface{}) uint64 {
	bs := []byte(key.(string))
	ans := uint64(0)
	for i := range bs {
		ans += uint64(bs[i] * 251)
	}
	return ans
}

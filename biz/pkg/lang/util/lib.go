package util

// Of 获取任意类型的指针
func Of[T any](v T) *T {
	if v == nil {
		return nil
	}
	return &v
}

// UnPtr 获取任意类型的指针的值
func UnPtr[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}

package utils

func IntPtrToInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func Int32PtrToInt32(p *int32) int32 {
	if p == nil {
		return 0
	}
	return *p
}

func StringPtrToString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func Float64PtrToFloat32(p *float64) float32 {
	if p == nil {
		return 0
	}
	return float32(*p)
}

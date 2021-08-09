package utils

func CopyStringUint32Map(m map[string]uint32) map[string]uint32 {
	cp := make(map[string]uint32)
	for k, v := range m {
		cp[k] = v
	}

	return cp
}

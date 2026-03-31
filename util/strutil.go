package util

// MergeMaps merge map
func MergeMaps(dst, src map[string]interface{}) {
	for k, v := range src {
		if existing, ok := dst[k]; ok {
			evMap, evOk := existing.(map[string]interface{})
			vMap, vOk := v.(map[string]interface{})
			if evOk && vOk {
				MergeMaps(evMap, vMap)
				continue
			}
		}
		dst[k] = v
	}
}

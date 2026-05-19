package utils

import maps0 "maps"

func MergeMaps(maps ...map[string]any) map[string]any {
	result := make(map[string]any)
	for _, m := range maps {
		maps0.Copy(result, m)
	}
	return result
}

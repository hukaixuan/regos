package utils

// Shard split keys by simple hash key
func Shard(key string, partNum int) int {
	if len(key) == 0 {
		return 0
	}
	sum := 0
	if len(key) < 16 {
		for i := 0; i < len(key); i++ {
			sum += int(key[i])
		}
	} else {
		for i := 0; i < 8; i++ {
			sum += int(key[i])
		}
		for i := len(key) - 8; i < len(key); i++ {
			sum += int(key[i])
		}
	}
	return sum % partNum
}

// func Shard(key string, partNum int) int {
// 	return fnv32(key) % partNum
// }

// const prime32 = 16777619

// func fnv32(key string) int {
// 	hash := 2166136261
// 	for i := 0; i < len(key); i++ {
// 		hash *= prime32
// 		hash ^= int(key[i])
// 	}
// 	if hash < 0 {
// 		hash = -hash
// 	}
// 	return hash
// }

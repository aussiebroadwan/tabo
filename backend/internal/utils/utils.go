package utils

func Contains(s []uint8, e uint8) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

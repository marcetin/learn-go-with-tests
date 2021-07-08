package iteration

const repeatCount = 5

// Repeat враћа знак поновљен 5 пута.
func Repeat(character string) string {
	var repeated string
	for i := 0; i < repeatCount; i++ {
		repeated += character
	}
	return repeated
}

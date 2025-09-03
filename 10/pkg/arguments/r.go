package arguments

// Reverse разворачивает входной слайс.
func Reverse(in []string) {
	left := 0
	right := len(in) - 1

	for left < right {
		in[left], in[right] = in[right], in[left]
		left++
		right--
	}
}

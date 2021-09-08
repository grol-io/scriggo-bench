package main

const size = 400

func main() {
	s := make([]int, size)
	for i := 0; i < size; i++ {
		s[i] = i
	}
	for _, x := range s {
		for j := range s {
			s[j] += x
		}
	}
}
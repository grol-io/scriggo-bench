package main

const calls = 90000

func main() {
	var b int
	for i := 0; i < calls; i++ {
		func(x int) {
			b += x
		}(i)
	}
}

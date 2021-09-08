package main

func main() {

	var sum int
	var f func()

	func() {
		n := 1
		f = func() {
			sum += n
		}
	}()

	for i := 0; i < 10_000; i++ {
		f()
	}

	// fmt.Println(sum)

}

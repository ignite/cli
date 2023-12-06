package testdata

import "fmt"

func fooTest() {
	n := "test new method"
	bla := fmt.Sprintf("test new - %s", n)
	fmt.
		Println(bla)
}

func BazTest() {
	foo := 100
	bar := fmt.Sprintf("test - %d", foo)
	fmt.Println(bar)
}

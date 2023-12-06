package testdata

import "fmt"

func fooTest() {
	foo := 100
	bar := fmt.Sprintf("test - %d", foo)
	fmt.Println(bar)
}

func BazTest() {
	foo := 100
	bar := fmt.Sprintf("test - %d", foo)
	fmt.Println(bar)
}

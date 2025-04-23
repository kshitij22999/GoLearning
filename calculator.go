package main

import "fmt"

func main() {
	var inp int
	var a, b int
	runner(inp, a, b)
}

func runner(inp int, a int, b int) {
	fmt.Println("Calculator menu")
	fmt.Println("1. Add \n2. Subtract\n3. Multiply\n4. Divide")
	fmt.Scan(&inp)
	fmt.Println("Please enter two numbers")
	fmt.Scan(&a, &b)
	switch inp {
	case 1:
		fmt.Printf("Answer is %d\n", add(a, b))
	case 2:
		fmt.Printf("Answer is %d\n", subtract(a, b))
	case 3:
		fmt.Printf("Answer is %d\n", multiply(a, b))
	case 4:
		fmt.Printf("Answer is %d\n", divide(a, b))
	default:
		fmt.Printf("Please select between 1 to 4")
	}
}

func add(a int, b int) int {
	return a + b
}

func subtract(a int, b int) int {
	return a - b
}

func multiply(a int, b int) int {
	return a * b
}

func divide(a int, b int) int {
	return a / b
}

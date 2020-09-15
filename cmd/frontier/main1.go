package main

import "fmt"

func append1(data *[]string){
	var data1 *[]string
	data1 = data
	*data1 = make([]string,4)
	fmt.Println(data1)
}

func main() {
	n := 2
	var n1 []string
	append1(&n1)
	n1 = n1[:n]
	for i := 0; i < n; i++ {
		fmt.Println(n1[i])
	}
}

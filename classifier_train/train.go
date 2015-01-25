package main

import (
	"fmt"
	"os"
)

var input, positive, negative *os.File

func init() {
  var err error
	input, err = os.OpenFile("training.txt", os.O_RDONLY,  0777)
	if (err != nil) {
		panic(err)
	}
	positive, err = os.OpenFile("training-1.txt", os.O_APPEND, 0777)
	if (err != nil) {
		panic(err)
	}
	negative, err = os.OpenFile("training-0.txt", os.O_APPEND, 0777)
	if (err != nil) {
		panic(err)
	}
}

func main() {
	fmt.Println("moving data to positive and negative data files based on user input")
	defer input.Close()
	defer positive.Close()
	defer negative.Close()
	
	
	
}
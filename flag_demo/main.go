package main

import (
	"fmt"
	"os"
)

func main() {
	//os.Args[0] 是当前执行的程序
	fmt.Println(os.Args)
	if len(os.Args) > 0 {
		for i, arg := range os.Args {
			fmt.Printf("args[%d]=%v\n", i, arg)
		}
	}
}

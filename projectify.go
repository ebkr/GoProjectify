package main

import (
	"fmt"
	"strconv"
	"./Modules/projectify"
)

func main() {
	fmt.Println("Selection an option: ")
	fmt.Println("1. New Project")
	fmt.Println("2. Load Project")
	fmt.Println("3. Delete Project")
	fmt.Println("4. Exit")
	var name string = ""
	fmt.Scanln(&name)
	var num, err = strconv.Atoi(name)
	if err == nil {
		appController(num)
	} else {
		main()
	}
}

func appController(area int) {
	switch area {
	case 1:
		fmt.Println("New Project")
		break
	case 2:
		fmt.Println("Load Project")
		break
	case 3:
		fmt.Println("Delete Project")
		break
	case 4:
		fmt.Println("Exit")
		break
	}
}

type 
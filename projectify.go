package main

import (
	"fmt"
	"strconv"

	"./Libraries/projectify"
)

func main() {
	// Display selections
	fmt.Println("Selection an option: ")
	fmt.Println("1. New Project")
	fmt.Println("2. Load Project")
	fmt.Println("3. Delete Project")
	fmt.Println("4. Exit")
	// Allow option selection
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
	fmt.Println("----------")
	switch area {
	case 1:
		// Create a new projectify project under /Projects
		fmt.Println("Enter a project name (No Spaces): ")
		scan := ""
		fmt.Scanln(&scan)
		create := projectify.StructCreate{}.New(scan + ".projectify")
		create.OverwriteFile("# An empty GoProjectify project\n<<TEMPLATEs>>\n<<BINDS>>\n<<POSITIONS>>")
		break
	case 2:
		// Load project #TODO
		fmt.Println("Load Project")
		break
	case 3:
		// Delete project #TODO
		fmt.Println("Delete Project")
		break
	case 4:
		// Exits the program
		fmt.Println("Exit")
		break
	}
}

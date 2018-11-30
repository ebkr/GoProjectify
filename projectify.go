package main

import (
	"fmt"
	"strconv"

	"./Libraries/projectify"
)

var loadedProject projectify.StructProject

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
		if appController(num) == true {
			main()
		}
	} else {
		main()
	}
}

func appController(area int) bool {
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
		fmt.Println("Enter project name: ")
		scan := ""
		fmt.Scanln(&scan)
		loadCase(scan)
		break
	case 3:
		// Delete project #TODO
		fmt.Println("Delete Project")
		break
	case 4:
		// Exits the program
		fmt.Println("Exit")
		fmt.Println("----------")
		fmt.Println("APPLICATION EXIT")
		fmt.Println("----------")
		return false
	}
	fmt.Println("----------")
	return true
}

func loadCase(load string) {
	loadedProject = projectify.StructProject{}
	fileProject := projectify.StructCreate{}.New(load + ".projectify")
	nodes := fileProject.GenerateNodeTree()
	for i := 0; i < len(nodes); i++ {
		fmt.Println("::")
		nodes[i].Print(1)
	}
}

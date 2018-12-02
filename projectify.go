package main

import (
	"fmt"
	"strconv"

	"./Libraries/projectify"
)

var loadedProject projectify.StructProject

func readInput() string {
	var inp string
	fmt.Scanln(&inp)
	return inp
}

func main() {
	// Display selections
	fmt.Println("Select an option: ")
	fmt.Println("1. New Project")
	fmt.Println("2. Load Project")
	fmt.Println("3. Delete Project")
	fmt.Println("4. Exit")
	// Allow option selection
	var num, err = strconv.Atoi(readInput())
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
		create := projectify.StructCreate{}.New(readInput() + ".projectify")
		create.OverwriteFile("# An empty GoProjectify project\n<<TEMPLATES>>\n<<BINDS>>\n<<POSITIONS>>")
		break
	case 2:
		// Load project #TODO
		fmt.Println("Enter project name: ")
		loadCase(readInput())
		break
	case 3:
		// Delete project #TODO
		fmt.Println("Delete Project")
		break
	case 4:
		// Exits the program
		fmt.Println("Exit")
		fmt.Println("----------")
		fmt.Println("----------")
		fmt.Println("APPLICATION EXIT")
		fmt.Println("----------")
		return false
	}
	fmt.Println("----------")
	return true
}

func generateProjectTree(fileProject *projectify.StructCreate, proj *projectify.StructProject) {
	proj.Init()
	nodes := fileProject.GenerateNodeTree()
	myMap := map[*projectify.StructNode]string{}
	for i := 0; i < len(nodes); i++ {
		myMap[nodes[i]] = nodes[i].GetValue()
	}
	proj.SetTree(myMap)
}

func loadCase(load string) {
	loadedProject = projectify.StructProject{}
	fileProject := projectify.StructCreate{}.New(load + ".projectify")
	proj := projectify.StructProject{}
	generateProjectTree(&fileProject, &proj)
	var exit bool
	for !exit {
		fmt.Println("--------->")
		fmt.Println("\\ Select an option: ")
		fmt.Println("1. Add new node")
		fmt.Println("2. Remove a node")
		var num, err = strconv.Atoi(readInput())
		if err == nil {
			switch num {
			case 1:
				fmt.Println("-------->>")
				fmt.Println("Name new node: ")
				name := readInput()
				generateProjectTree(&fileProject, &proj)
				fileProject.AppendFile("<<TEMPLATES>>", strconv.Itoa(proj.GetAvailableId())+":"+name)
				break
			case 2:
				generateProjectTree(&fileProject, &proj)
				fmt.Println("Name of node: ")
				node := proj.GetNodeByName(readInput())
				if node != nil {
					fileProject.RemoveLine(strconv.Itoa(node.GetId()) + ":" + node.GetValue())
				}
			}
		}
	}
}

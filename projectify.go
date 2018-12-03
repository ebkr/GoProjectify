package main

import (
	"fmt"
	"strconv"

	"./Libraries/projectify"
)

var loadedProject projectify.StructProject

// Handle user input
func readInput() string {
	fmt.Print("> ")
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
		// Load project
		fmt.Println("Enter project name: ")
		loadCase(readInput())
		break
	case 3:
		// Delete project
		fmt.Println("Delete Project")
		deleteCase(readInput())
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

// Generate Project Tree
// Refreshes StructProject, and provide updated node list.
func generateProjectTree(fileProject *projectify.StructCreate, proj *projectify.StructProject) {
	proj.Init()
	nodes := fileProject.GenerateNodeTree()
	myMap := map[*projectify.StructNode]string{}
	for i := 0; i < len(nodes); i++ {
		myMap[nodes[i]] = nodes[i].GetValue()
	}
	proj.SetTree(myMap)
}

// Run when user is loading a project.
// Provides options for loaded project
func loadCase(load string) {
	fileProject := projectify.StructCreate{}.New(load + ".projectify")
	if !fileProject.CheckExistence() {
		fmt.Println("File Not Found")
		return
	}
	loadedProject = projectify.StructProject{}
	proj := projectify.StructProject{}
	generateProjectTree(&fileProject, &proj)
	var exit bool
	for !exit {
		// Opts contains menu options to display
		opts := []string{"Add new node", "Remove a node", "Link nodes", "Print", "Return"}
		// Funcs contains functions, with an index corresponding to their menu counterparts
		funcs := []func(){
			// New node
			func() {
				fmt.Println("-------->>")
				fmt.Println("Name new node: ")
				name := readInput()
				generateProjectTree(&fileProject, &proj)
				fileProject.AppendFile("<<TEMPLATES>>", strconv.Itoa(proj.GetAvailableID())+":"+name)
			},
			// Remove node
			func() {
				fmt.Println("-------->>")
				generateProjectTree(&fileProject, &proj)
				fmt.Println("Name of node: ")
				node := proj.GetNodeByName(readInput())
				if node != nil {
					fileProject.RemoveLine(strconv.Itoa(node.GetId()) + ":" + node.GetValue())
				}
			},
			// Link nodes
			func() {
				fmt.Println("-------->>")
				fmt.Println("Start Node: ")
				nodeA := proj.GetNodeByName(readInput())
				fmt.Println("End Node: ")
				nodeB := proj.GetNodeByName(readInput())
				if nodeA != nil && nodeB != nil {
					fmt.Println("...........")
					if nodeA.AddConnection(nodeB) {
						fmt.Println("Connected " + nodeA.GetValue() + " to " + nodeB.GetValue())
						fileProject.RemoveLine(strconv.Itoa(nodeA.GetId()) + ":" + strconv.Itoa(nodeB.GetId()))
						fileProject.AppendFile("<<BINDS>>", strconv.Itoa(nodeA.GetId())+":"+strconv.Itoa(nodeB.GetId()))
					} else {
						fmt.Println("Two nodes are already connected")
						fmt.Println("Recursion is not allowed")
					}
					fmt.Println("...........")
				}
				generateProjectTree(&fileProject, &proj)
			},
			// Display nodes
			func() {
				fmt.Println("-------->>")
				fmt.Println("Print Node: ")
				nodeA := proj.GetNodeByName(readInput())
				if nodeA != nil {
					nodeA.Print(1)
				}
			},
			// Back to main menu
			func() {
				exit = true
			},
		}
		// Display options, and handle selection
		fmt.Println("--------->")
		for i := 0; i < len(opts); i++ {
			fmt.Println(strconv.Itoa(i+1) + ". " + opts[i])
		}
		var num, err = strconv.Atoi(readInput())
		if err == nil {
			if num > 0 && num <= len(opts) {
				funcs[num-1]()
			}
		}
	}
}

func deleteCase(load string) {
	fileProject := projectify.StructCreate{}.New(load + ".projectify")
	if fileProject.CheckExistence() {
		fileProject.Delete()
		fmt.Println("Deleted Project: " + load)
	} else {
		fmt.Println("Invalid Name")
	}
}

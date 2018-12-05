package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"./ApplicationData/Libraries/projectify"
)

var loadedProject projectify.StructProject

// Handle user input
func readInput() string {
	fmt.Print("> ")
	var inp string
	fmt.Scanln(&inp)
	return inp
}

func buildString(array []string, index int) string {
	var str = "/"
	for i := index; i < len(array); i++ {
		str += array[i] + "/"
	}
	return str[:len(str)-1]
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./ApplicationData/Views")))
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		split := strings.Split(r.URL.Path, "/")
		if split[2] == "GetProjects" {
			files, err := filepath.Glob("./Projects/*.projectify")
			if err == nil {
				for i := 0; i < len(files); i++ {
					str := strings.TrimLeft(files[i], "Projects\\")
					str = str[:len(str)-11]
					w.Write([]byte(str))
					if i < len(files)-1 {
						w.Write([]byte("\n"))
					}
				}
			}
		} else if split[2] == "NewProject" {
			str := buildString(split, 3)
			create := projectify.StructCreate{}.New(str + ".projectify")
			create.OverwriteFile("# An empty GoProjectify project\n<<TEMPLATES>>\n<<BINDS>>\n<<POSITIONS>>")
		} else if split[2] == "LoadProject" {
			loadCase(buildString(split, 6), split[3], split, &w, r)
		}
	})
	http.ListenAndServe(":8080", nil)
	// Display selections
	/*
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
	*/
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
		//loadCase(readInput())
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
func loadCase(load string, use string, split []string, w *http.ResponseWriter, r *http.Request) {
	writer := *w
	fileProject := projectify.StructCreate{}.New(load + ".projectify")
	if !fileProject.CheckExistence() {
		fmt.Println(fileProject.Dir + fileProject.Name)
		fmt.Println("File Not Found")
		return
	}
	loadedProject = projectify.StructProject{}
	proj := projectify.StructProject{}
	generateProjectTree(&fileProject, &proj)
	// Funcs contains functions, with an index corresponding to their menu counterparts
	writer.Write([]byte("<<OUTPUT>>\n"))

	funcs := map[string]func(){
		"None": func() {},
		"NewNode": func() {
			fmt.Println("-------->>")
			fmt.Println("Name new node: ")
			name := split[4]
			generateProjectTree(&fileProject, &proj)
			// Prevent duplicate names
			for k := range proj.GetTree() {
				if k.GetValue() == name {
					return
				}
			}
			fileProject.AppendFile("<<TEMPLATES>>", strconv.Itoa(proj.GetAvailableID())+":"+name)
		},
		// Remove node
		"RemoveNode": func() {
			fmt.Println("-------->>")
			generateProjectTree(&fileProject, &proj)
			fmt.Println("Name of node: ")
			node := proj.GetNodeByName(split[4])
			if node != nil {
				fileProject.RemoveLine(strconv.Itoa(node.GetId()) + ":" + node.GetValue())
			}
		},
		// Link nodes
		"Link": func() {
			fmt.Println("-------->>")
			fmt.Println("Start Node: ")
			nodeA := proj.GetNodeByName(split[4])
			fmt.Println("End Node: ")
			nodeB := proj.GetNodeByName(split[5])
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
		// Print
		"Print": func() {
			nodeA := proj.GetNodeByName(split[4])
			if nodeA != nil {
				// nodeA.Print(1)
			}
		},
		// Get all nodes
		"Get": func() {
			fmt.Println("-------->>")
			fmt.Println("Nodes: ")
			for k := range proj.GetTree() {
				fmt.Println(strconv.Itoa(k.GetId()) + ":" + k.GetValue())
			}
		},
	}
	// Display options, and handle selection
	funcs[use]()
	generateProjectTree(&fileProject, &proj)
	writer.Write([]byte("<<GENERATE>>\n"))
	for k := range proj.GetTree() {
		writer.Write([]byte("Node:" + strconv.Itoa(k.GetId()) + ":" + k.GetValue() + "\n"))
		for i := 0; i < len(k.Connections); i++ {
			writer.Write([]byte("Connection:" + strconv.Itoa(k.Connections[i].GetId()) + "\n"))
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

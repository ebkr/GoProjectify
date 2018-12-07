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
			result := fileProject.NewNode(proj.GetAvailableID(), name)
			if !result {
				writer.Write([]byte("WARN:Illegal Character"))
			}
		},
		// Remove node
		"RemoveNode": func() {
			fmt.Println("-------->>")
			generateProjectTree(&fileProject, &proj)
			fmt.Println("Name of node: ")
			id, err := strconv.Atoi(split[4])
			if err == nil {
				node := proj.GetNodeByID(id)
				if node != nil {
					fileProject.RemoveNode(node.GetId())
				}
			}
		},
		"RemoveLink": func() {
			fmt.Println("-------->>")
			fmt.Println("Start Node: ")
			id1, _ := strconv.Atoi(split[4])
			id2, _ := strconv.Atoi(split[5])
			nodeA := proj.GetNodeByID(id1)
			fmt.Println("End Node: ")
			nodeB := proj.GetNodeByID(id2)
			if nodeA != nil && nodeB != nil {
				fileProject.RemoveLink(id1, id2)
			}
			generateProjectTree(&fileProject, &proj)
		},
		// Link nodes
		"Link": func() {
			fmt.Println("-------->>")
			fmt.Println("Start Node: ")
			id1, _ := strconv.Atoi(split[4])
			id2, _ := strconv.Atoi(split[5])
			nodeA := proj.GetNodeByID(id1)
			fmt.Println("End Node: ")
			nodeB := proj.GetNodeByID(id2)
			if nodeA != nil && nodeB != nil {
				fmt.Println("...........")
				if nodeA.AddConnection(nodeB) {
					fmt.Println("Connected " + nodeA.GetValue() + " to " + nodeB.GetValue())
					fileProject.AppendFile("<<BINDS>>", split[4]+":"+split[5])
				} else {
					writer.Write([]byte("Action not allowed. Nodes are already connected" + "\n"))
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
		// Set node position
		"Reposition": func() {
			id, err1 := strconv.Atoi(split[4])
			positions := strings.Split(split[5], ":")
			x, err2 := strconv.Atoi(positions[0])
			y, err3 := strconv.Atoi(positions[1])
			if err1 == err2 && err1 == err3 && err1 == nil {
				node := proj.GetNodeByID(id)
				if node != nil {
					fileProject.SetPosition(id, x, y)
				}
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
		x := int(k.GetPosition()[0])
		y := int(k.GetPosition()[1])
		writer.Write([]byte("Position:" + strconv.Itoa(x) + ":" + strconv.Itoa(y) + "\n"))
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

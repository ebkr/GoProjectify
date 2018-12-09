package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"./ApplicationData/Libraries/projectify"
)

// Globals
var loadedProject projectify.StructProject
var configData projectify.StructConf

// buildString : Build a string from an array, starting at index
func buildString(array []string, index int) string {
	var str = "/"
	for i := index; i < len(array); i++ {
		str += array[i] + "/"
	}
	return str[:len(str)-1]
}

// buildApplicationReguirements : Create Project directory, and default configurations
func buildApplicationRequirements(workingDir string) {
	configs := map[string][]string{
		"Server.conf": []string{"URL:localhost", "Port:8080", "ProjectDirectory:" + workingDir + "\\Projects"},
	}
	// Create Config directory
	configFolder := projectify.StructCreate{}.New(workingDir+"/Config", "")
	if !configFolder.CheckExistence() {
		os.Mkdir(workingDir+"/Config", os.FileMode(0755))
		// Create config files defined in "configs" variable.
		for file, data := range configs {
			f := projectify.StructCreate{}.New(workingDir+"/Config/", file)
			if !f.CheckExistence() {
				var builder string
				for _, line := range data {
					builder += line + "\n"
				}
				f.OverwriteFile(builder)
			}
		}
	}
	// Assign global to new config
	configData = projectify.StructConf{}.New(workingDir + "/Config/Server.conf")

	// Create Project directory
	projectFolder := projectify.StructCreate{}.New(configData.GetKey("ProjectDirectory"), "")
	if !projectFolder.CheckExistence() {
		os.Mkdir(configData.GetKey("ProjectDirectory"), os.FileMode(0755))
	}
}

// main : Main Method, loads web server
func main() {
	workingDir, dirErr := os.Getwd()
	if dirErr != nil {
		panic("Can't get the working directory")
	}
	buildApplicationRequirements(workingDir)
	http.Handle("/", http.FileServer(http.Dir(workingDir+"\\ApplicationData\\Views")))
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		split := strings.Split(r.URL.Path, "/")
		if split[2] == "GetProjects" {
			files, err := filepath.Glob(configData.GetKey("ProjectDirectory") + "\\*.projectify")
			if err == nil {
				for i := 0; i < len(files); i++ {
					str := files[i][(len(configData.GetKey("ProjectDirectory")) + 1):]
					str = str[:len(str)-11]
					w.Write([]byte(str))
					if i < len(files)-1 {
						w.Write([]byte("\n"))
					}
				}
			}
		} else if split[2] == "NewProject" {
			str := buildString(split, 3)
			create := projectify.StructCreate{}.New(configData.GetKey("ProjectDirectory")+"\\", str+".projectify")
			create.OverwriteFile("# An empty GoProjectify project\n<<TEMPLATES>>\n<<BINDS>>\n<<POSITIONS>>")
		} else if split[2] == "LoadProject" {
			loadCase(buildString(split, 6), split[3], split, &w, r)
		}
	})
	http.ListenAndServe(configData.GetKey("URL")+":"+configData.GetKey("Port"), nil)
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
	fileProject := projectify.StructCreate{}.New(configData.GetKey("ProjectDirectory")+"\\", load+".projectify")
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

// deleteCase : Delete a project
func deleteCase(load string) {
	fileProject := projectify.StructCreate{}.New(configData.GetKey("ProjectDirectory")+"\\", load+".projectify")
	if fileProject.CheckExistence() {
		fileProject.Delete()
		fmt.Println("Deleted Project: " + load)
	} else {
		fmt.Println("Invalid Name")
	}
}

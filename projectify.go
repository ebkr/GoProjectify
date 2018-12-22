package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ebkr/GoProjectify/ApplicationData/Libraries/projectify"
)

// Globals
var configData projectify.StructConf

// buildApplicationReguirements : Create Project directory, and default configurations
func buildApplicationRequirements(workingDir string) {
	configs := map[string][]string{
		"Server.ini": {"URL:localhost", "Port:8080", "ProjectDirectory:" + workingDir + "\\Projects"},
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
	configData = projectify.StructConf{}.New(workingDir + "/Config/Server.ini")
	// Create Project directory
	projectFolder := projectify.StructCreate{}.New(configData.GetKey("ProjectDirectory"), "")
	if !projectFolder.CheckExistence() {
		os.Mkdir(configData.GetKey("ProjectDirectory"), os.FileMode(0755))
	}
}

// main : Main Method, loads web server
func main() {
	fmt.Println("> Preparing")
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
			str := split[3]
			create := projectify.StructCreate{}.New(configData.GetKey("ProjectDirectory")+"\\", str+".projectify")
			if !create.CheckExistence() {
				create.OverwriteFile("# An empty GoProjectify project\n<<TEMPLATES>>\n<<BINDS>>\n<<POSITIONS>>")
			}
		} else if split[2] == "LoadProject" {
			log.Println(split[3] + " on project: " + split[6])
			loadCase(split[6], split[3], split, &w, r)
		}
	})
	fmt.Println("> Starting Server on Port: " + configData.GetKey("Port"))
	http.ListenAndServe(configData.GetKey("URL")+":"+configData.GetKey("Port"), nil)
}

// Run when user is loading a project.
// Provides options for loaded project
func loadCase(load string, use string, split []string, w *http.ResponseWriter, r *http.Request) {
	writer := *w
	fileProject := projectify.StructCreate{}.New(configData.GetKey("ProjectDirectory")+"\\", load+".projectify")
	if !fileProject.CheckExistence() {
		return
	}
	proj := projectify.StructProject{}
	loader := LoadCase{
		FileProject: &fileProject,
		Proj:        &proj,
	}
	err := loader.Use(use, split)
	loader.GenerateProjectTree()
	writer.Write([]byte("<<OUTPUT>>\n"))
	if err != nil {
		writer.Write([]byte(err.Error() + "\n"))
		log.Println(err.Error())
	}
	writer.Write([]byte("<<GENERATE>>\n"))
	for k := range loader.Proj.GetTree() {
		writer.Write([]byte("Node:" + strconv.Itoa(k.GetID()) + ":" + k.GetValue() + "\n"))
		for i := 0; i < len(k.Connections); i++ {
			writer.Write([]byte("Connection:" + strconv.Itoa(k.Connections[i].GetID()) + "\n"))
		}
		x := int(k.GetPosition()[0])
		y := int(k.GetPosition()[1])
		writer.Write([]byte("Position:" + strconv.Itoa(x) + ":" + strconv.Itoa(y) + "\n"))
	}
}

// LoadCase : Handles load cases
type LoadCase struct {
	FileProject *projectify.StructCreate
	Proj        *projectify.StructProject
	cases       map[string]func([]string) error
}

// initialiseCases : Set up {cases}
func (load *LoadCase) initialiseCases() {
	load.cases = map[string]func([]string) error{
		"NewNode":       load.newNode,
		"RemoveNode":    load.removeNode,
		"RemoveLink":    load.removeLink,
		"Link":          load.link,
		"Reposition":    load.reposition,
		"DeleteProject": load.deleteProject,
	}
}

// Use : Perform a {cases} action
func (load *LoadCase) Use(name string, split []string) error {
	if load.cases == nil {
		load.initialiseCases()
	}
	load.GenerateProjectTree()
	if load.cases[name] != nil {
		return load.cases[name](split)
	}
	return nil
}

// newNode : Creates a new node.
func (load *LoadCase) newNode(split []string) error {
	name := split[4]
	pos := split[5]
	var x, y int

	nums := strings.Split(pos, ":")
	if len(nums) > 1 {
		var err, err2 error
		x, err = strconv.Atoi(nums[0])
		y, err2 = strconv.Atoi(nums[1])
		if err != nil || err2 != nil {
			x = 0
			y = 0
		}
	}

	id := load.Proj.GetAvailableID()
	load.GenerateProjectTree()
	result := load.FileProject.NewNode(id, name)
	load.GenerateProjectTree()
	if !result {
		return errors.New("error: Name Contains Illegal Character")
	}
	node := load.Proj.GetNodeByID(id)
	if node != nil {
		load.FileProject.SetPosition(id, x, y)
	} else {
		return errors.New("error: Could not find node. Ensure write permissions are enabled")
	}
	return nil
}

// removeNode : Removes a node.
func (load *LoadCase) removeNode(split []string) error {
	id, err := strconv.Atoi(split[4])
	if err == nil {
		node := load.Proj.GetNodeByID(id)
		fmt.Println("Removing node: " + strconv.Itoa(node.GetID()))
		if node != nil {
			load.FileProject.RemoveNode(node.GetID())
		} else {
			return errors.New("error: Invalid Node ID")
		}
	} else {
		return err
	}
	return nil
}

// removeLink : Removes the connection between nodes
func (load *LoadCase) removeLink(split []string) error {
	id1, _ := strconv.Atoi(split[4])
	id2, _ := strconv.Atoi(split[5])
	nodeA := load.Proj.GetNodeByID(id1)
	nodeB := load.Proj.GetNodeByID(id2)
	if nodeA != nil && nodeB != nil {
		load.FileProject.RemoveLink(id1, id2)
	} else {
		return errors.New("error: Link doesn't exist")
	}
	return nil
}

// link : Link two nodes together
func (load *LoadCase) link(split []string) error {
	id1, _ := strconv.Atoi(split[4])
	id2, _ := strconv.Atoi(split[5])
	nodeA := load.Proj.GetNodeByID(id1)
	nodeB := load.Proj.GetNodeByID(id2)
	if nodeA != nil && nodeB != nil {
		if nodeA.AddConnection(nodeB) {
			load.FileProject.AppendFile("<<BINDS>>", split[4]+":"+split[5])
		} else {
			return errors.New("error: Action not allowed. Nodes are already connected")
		}
	}
	return nil
}

// reposition : Reposition a node
func (load *LoadCase) reposition(split []string) error {
	id, err1 := strconv.Atoi(split[4])
	positions := strings.Split(split[5], ":")
	x, err2 := strconv.Atoi(positions[0])
	y, err3 := strconv.Atoi(positions[1])
	if err1 == err2 && err1 == err3 && err1 == nil {
		node := load.Proj.GetNodeByID(id)
		if node != nil {
			load.FileProject.SetPosition(id, x, y)
		} else {
			return errors.New("error: Invalid Node ID")
		}
	} else {
		return errors.New("error: String->Int conversion error")
	}
	return nil
}

// deleteProject : Delete a project
func (load *LoadCase) deleteProject(split []string) error {
	err := load.FileProject.Delete()
	load.FileProject = &projectify.StructCreate{}
	return err
}

// GenerateProjectTree : Generates a node tree for use with the project
func (load *LoadCase) GenerateProjectTree() {
	load.Proj.Init()
	nodes := load.FileProject.GenerateNodeTree()
	myMap := map[*projectify.StructNode]string{}
	for i := 0; i < len(nodes); i++ {
		myMap[nodes[i]] = nodes[i].GetValue()
	}
	load.Proj.SetTree(myMap)
}

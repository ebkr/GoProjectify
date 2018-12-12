package projectify

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// StructCreate : Struct containing FileName, FileDirectory, and Data (string of file text)
type StructCreate struct {
	Name string
	Dir  string
	data string
}

// New : Used to generate a working Struct
func (ref StructCreate) New(Directory, Name string) StructCreate {
	c := StructCreate{Name, Directory, ""}
	return c
}

// OverwriteFile : Used to override file contents with specified string.
func (ref StructCreate) OverwriteFile(data string) bool {
	file, err := os.OpenFile(ref.Dir+ref.Name, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return false
	}
	file.Truncate(0)
	file.WriteString(data)
	file.Close()
	return true
}

// AppendFile : Used to append a string to the file.
func (ref StructCreate) AppendFile(after, newLine string) bool {
	file, err := os.OpenFile(ref.Dir+ref.Name, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return false
	}
	scanner := bufio.NewScanner(file)
	text := ""
	for scanner.Scan() {
		text += scanner.Text() + "\n"
		if scanner.Text() == after {
			text += newLine + "\n"
		}
	}
	file.Close()
	return ref.OverwriteFile(text)
}

// RemoveLine : Used to remove a line from the text document
func (ref *StructCreate) removeLine(search string) bool {
	ref.updateReadData()
	split := strings.Split(ref.data, "\n")
	ref.data = ""
	var removed bool
	for i := 0; i < len(split); i++ {
		if split[i] != search && split[i] != "" {
			ref.data += split[i] + "\n"
		} else {
			removed = true
		}
	}
	if removed {
		ref.OverwriteFile(ref.data)
	}
	return removed
}

// getTextInSection : Scan through update data, and get each line within a section (<<SECTION>>)
func (ref *StructCreate) getTextInSection(region string) []string {
	ref.updateReadData()
	var stringBuilder string
	split := strings.Split(ref.data, "\n")
	var section string
	for i := 0; i < len(split); i++ {
		if split[i] == region {
			section = split[i]
		} else if section == region {
			stringBuilder += split[i] + "\n"
		}
	}
	if stringBuilder == "" {
		return []string{}
	}
	group := strings.Split(stringBuilder, "\n")
	return group[:len(group)-1]
}

// NewNode : Creates a new node, assuming no "illegal" characters exist. Returns true if successful.
func (ref *StructCreate) NewNode(id int, nodeName string) bool {
	split := strings.Split(nodeName, "<")
	if len(split) > 1 {
		return false
	}
	split = strings.Split(split[0], ":")
	if len(split) > 1 {
		return false
	}
	ref.AppendFile("<<TEMPLATES>>", strconv.Itoa(id)+":"+split[0])
	return true
}

// RemoveNode : Remove node from projectify file with ID of nodeId
func (ref *StructCreate) RemoveNode(nodeID int) {
	// Remove Templates
	data := ref.getTextInSection("<<TEMPLATES>>")
	for i := 0; i < len(data); i++ {
		split := strings.Split(data[i], ":")
		if split[0] == strconv.Itoa(nodeID) {
			ref.removeLine(data[i])
		}
	}
	// Remove Links
	data = ref.getTextInSection("<<BINDS>>")
	for i := 0; i < len(data); i++ {
		split := strings.Split(data[i], ":")
		for j := 0; j < len(split); j++ {
			if split[j] == strconv.Itoa(nodeID) {
				ref.removeLine(data[i])
			}
		}
	}
}

// RemoveLink : Remove connection between NodeIdA and NodeIdB in projectify file
func (ref *StructCreate) RemoveLink(nodeIDA, nodeIDB int) {
	// Remove Links
	data := ref.getTextInSection("<<BINDS>>")
	for i := 0; i < len(data); i++ {
		split := strings.Split(data[i], ":")
		if split[0] == strconv.Itoa(nodeIDA) || split[0] == strconv.Itoa(nodeIDB) {
			ref.removeLine(data[i])
		} else if split[1] == strconv.Itoa(nodeIDA) || split[1] == strconv.Itoa(nodeIDB) {
			ref.removeLine(data[i])
		}
	}
}

// SetPosition : Update the position of nodeId in the projectify file.
func (ref *StructCreate) SetPosition(nodeID, x, y int) {
	// Remove Links
	data := ref.getTextInSection("<<POSITIONS>>")
	for i := 0; i < len(data); i++ {
		split := strings.Split(data[i], ":")
		if split[0] == "["+strconv.Itoa(nodeID)+"]" {
			ref.removeLine(data[i])
		}
	}
	ref.AppendFile("<<POSITIONS>>", "["+strconv.Itoa(nodeID)+"]:"+strconv.Itoa(x)+":"+strconv.Itoa(y))
}

// GenerateNodeTree : Creates an array of StructNodes, and links them together using the StructCreate Data
func (ref *StructCreate) GenerateNodeTree() []*StructNode {
	ref.updateReadData()
	split := strings.Split(ref.data, "\n")

	var action = "#"
	templateNodes := []*StructNode{}

	for i := 0; i < len(split); i++ {
		if strings.Trim(split[i], "#") == "#" {
			// Ignore
		} else if strings.Contains(split[i], "<<") {
			action = split[i]
		} else {
			if action == "<<TEMPLATES>>" {
				splitTwice := strings.Split(split[i], ":")
				id, err := strconv.Atoi(splitTwice[0])
				if err == nil {
					value := splitTwice[1]
					tempNode := StructNode{}.New(id, value, 0, 0)
					templateNodes = append(templateNodes, &tempNode)
				}
			} else if action == "<<BINDS>>" {
				nodeSectionBinds(split, i, &templateNodes)
			} else if action == "<<POSITIONS>>" {
				nodeSectionPositions(split, i, &templateNodes)
			}
		}
	}
	return templateNodes
}

// nodeSectionBinds : Generate connections for nodes, called in GenerateNodeTree()
func nodeSectionBinds(split []string, i int, nodes *[]*StructNode) {
	templateNodes := *nodes
	splitTwice := strings.Split(split[i], ":")
	id, err := strconv.Atoi(splitTwice[0])
	id2, err2 := strconv.Atoi(splitTwice[1])
	if err == nil && err2 == nil {
		var nodeA *StructNode
		var nodeB *StructNode
		for search := 0; search < len(templateNodes); search++ {
			if templateNodes[search].GetID() == id {
				nodeA = templateNodes[search]
			} else if templateNodes[search].GetID() == id2 {
				nodeB = templateNodes[search]
			}
		}
		if nodeA != nil && nodeB != nil {
			nodeA.AddConnection(nodeB)
		}
	}
}

// nodeSectionPositions : Generate positions for nodes, called in GenerateNodeTree()
func nodeSectionPositions(split []string, i int, nodes *[]*StructNode) {
	templateNodes := *nodes
	splitThrice := strings.Split(split[i], ":")
	if len(splitThrice) == 3 {
		idString := regexp.MustCompile("[0-9]+").FindString(splitThrice[0])
		id, err := strconv.Atoi(idString)
		if err == nil {
			x, errX := strconv.Atoi(splitThrice[1])
			y, errY := strconv.Atoi(splitThrice[2])
			if errX == errY && errX == nil {
				for search := 0; search < len(templateNodes); search++ {
					if templateNodes[search].GetID() == id {
						templateNodes[search].SetPosition(float64(x), float64(y))
					}
				}
			}
		}
	}
}

// updateReadData : Read the corresponding file, and place in to data field.
func (ref *StructCreate) updateReadData() {
	file, err := os.OpenFile(ref.Dir+ref.Name, os.O_RDONLY|os.O_CREATE, 0666)
	if err == nil {
		ref.data = ""
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			ref.data += scanner.Text() + "\n"
		}
		file.Close()
	}
}

// Delete : Delete the file/directory
func (ref *StructCreate) Delete() error {
	err := os.Remove(ref.Dir + ref.Name)
	return err
}

// CheckExistence : Used to check if file/directory exists.
func (ref *StructCreate) CheckExistence() bool {
	_, err := os.Stat(ref.Dir + ref.Name)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

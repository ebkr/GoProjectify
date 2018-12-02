package projectify

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Base data
type StructCreate struct {
	Name string
	Dir  string
	data string
}

// Used to generate a working Struct
func (ref StructCreate) New(Name string) StructCreate {
	c := StructCreate{Name, "./Projects/", ""}
	return c
}

// Used to override file contents with specified string.
func (ref StructCreate) OverwriteFile(data string) bool {
	file, err := os.OpenFile(ref.Dir+ref.Name, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return false
	} else {
		file.Truncate(0)
		file.WriteString(data)
		file.Close()
	}
	return true
}

// Used to append a string to the file.
func (ref StructCreate) AppendFile(after, newLine string) bool {
	file, err := os.OpenFile(ref.Dir+ref.Name, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return false
	} else {
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
}

// Used to remove a line from the text document
func (ref *StructCreate) RemoveLine(search string) bool {
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

func (ref StructCreate) GenerateNodeTree() []*StructNode {
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
					tempNode := StructNode{}.New(id, value)
					templateNodes = append(templateNodes, &tempNode)
				}
			} else if action == "<<BINDS>>" {
				splitTwice := strings.Split(split[i], ":")
				id, err := strconv.Atoi(splitTwice[0])
				id2, err2 := strconv.Atoi(splitTwice[1])
				if err == nil && err2 == nil {
					templateNodes[id].AddConnection(templateNodes[id2])
				}
			}
		}
	}
	return templateNodes
}

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

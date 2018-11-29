package projectify

import(
	"os"
	"bufio"
)

// Base data
type StructCreate struct {
	Name string
	Dir  string
}

// Used to generate a working Struct
func (ref StructCreate) New(Name string) StructCreate {
	c := StructCreate{Name, "./Projects/"}
	return c
}

// Used to override file contents with specified string.
func (ref StructCreate) OverwriteFile(data string) bool {
	file,err := os.OpenFile(ref.Dir + ref.Name, os.O_RDONLY|os.O_CREATE, 0666)
	if (err != nil) {
		return false;
	} else {
		file.WriteString(data)
		file.Close()
	}
	return true
}

// Used to append a string to the file.
func (ref StructCreate) AppendFile(data string) bool {
	file,err := os.OpenFile(ref.Name, os.O_RDONLY|os.O_CREATE, 0666)
	if (err != nil) {
		return false;
	} else {
		scanner := bufio.NewScanner(file)
		text := ""
    	for scanner.Scan() {
			text += scanner.Text() + "\n"
		}
		text += data
		file.Close();
		return ref.OverwriteFile(text)
	}
}
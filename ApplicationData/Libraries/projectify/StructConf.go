package projectify

import (
	"bufio"
	"os"
	"strings"
)

// StructConf : Struct to load and retrieve Config data
type StructConf struct {
	path    string
	configs map[string]string
}

// New : Generates defaults
func (ref StructConf) New(path string) StructConf {
	return StructConf{path, map[string]string{}}
}

// update : Updates configuration data
func (ref *StructConf) update() {
	file, err := os.OpenFile(ref.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			split := strings.Split(scanner.Text(), ":")
			var build string
			for i := 1; i < len(split); i++ {
				if i > 1 {
					build += ":"
				}
				build += split[i]
			}
			ref.configs[split[0]] = build
		}
		file.Close()
	}
}

// GetKey : Returns the string associated with the config key
func (ref *StructConf) GetKey(key string) string {
	ref.update()
	return ref.configs[key]
}

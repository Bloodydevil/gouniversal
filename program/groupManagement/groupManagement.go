package groupManagement

import (
	"encoding/json"
	"gouniversal/program/global"
	"gouniversal/program/programTypes"
	"gouniversal/program/ui/uifunc"
	"gouniversal/shared/config"
	"gouniversal/shared/io/file"
	"log"
	"os"
)

const GroupFile = "data/config/group"

func SaveGroup(gc programTypes.GroupConfigFile) error {

	gc.Header = config.BuildHeader("group", "groups", 1.0, "group config file")

	if _, err := os.Stat(GroupFile); os.IsNotExist(err) {
		// if not found, create default file

		newgroup := make([]programTypes.Group, 1)

		newgroup[0].UUID = "admin"
		newgroup[0].Name = "admin"
		newgroup[0].State = 1 // active

		pages := []string{"Program:Settings:User", "Program:Settings:User:List", "Program:Settings:User:Edit", "Program:Settings:Group", "Program:Settings:Group:List", "Program:Settings:Group:Edit"}
		newgroup[0].AllowedPages = pages

		gc.Group = newgroup
	}

	b, err := json.Marshal(gc)
	if err != nil {
		log.Fatal(err)
	}

	f := new(file.File)
	err = f.WriteFile(GroupFile, b)

	return err
}

func LoadGroup() programTypes.GroupConfigFile {

	var gc programTypes.GroupConfigFile

	if _, err := os.Stat(GroupFile); os.IsNotExist(err) {
		// if not found, create default file
		SaveGroup(gc)
	}

	f := new(file.File)
	b, err := f.ReadFile(GroupFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &gc)
	if err != nil {
		log.Fatal(err)
	}

	if config.CheckHeader(gc.Header, "groups") == false {
		log.Fatal("wrong config")
	}

	return gc
}

func IsPageAllowed(path string, gid string, checkState bool) bool {

	global.GroupConfig.Mut.Lock()
	defer global.GroupConfig.Mut.Unlock()

	for g := 0; g < len(global.GroupConfig.File.Group); g++ {

		if gid == global.GroupConfig.File.Group[g].UUID {

			if checkState {
				// if group is not active
				if global.GroupConfig.File.Group[g].State != 1 {
					return false
				}
			}

			for p := 0; p < len(global.GroupConfig.File.Group[g].AllowedPages); p++ {

				allowed := uifunc.RemovePFromPath(global.GroupConfig.File.Group[g].AllowedPages[p])
				requested := uifunc.RemovePFromPath(path)

				if requested == allowed {

					return true
				}
			}

			return false
		}
	}

	return false
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Env struct {
	Variables map[string]string `json:"Variables"`
}
type Configuration struct {
	Environment Env `json:"Environment"`
}

const exportVariablePath string = "/Users/architagarwal/.zprofile"

func main() {
	fmt.Println("setting env variables started")
	var allEnv map[string]string = make(map[string]string)
	var env Configuration
	f, err := os.OpenFile(exportVariablePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	if _, err = f.WriteString(fmt.Sprintf("\n")); err != nil {
		panic(err)
	}
	defer f.Close()

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.Contains(info.Name(), ".json") && strings.Contains(info.Name(), "configuration-") {
			fmt.Println("started processing ", info.Name())
			jsonFile, err := os.Open(info.Name())
			if err != nil {
				fmt.Printf("Error opening file %s:%+v", info.Name(), err)
			}
			byteValue, _ := ioutil.ReadAll(jsonFile)
			err1 := json.Unmarshal(byteValue, &env)
			if err1 != nil {
				fmt.Printf("Error unmarshalling file %s:%+v", info.Name(), err)
			}
			for key, value := range env.Environment.Variables {
				_, added := allEnv[key]
				_, exist := os.LookupEnv(key)
				if !exist && !added {
					allEnv[key] = value
					if _, err = f.WriteString(fmt.Sprintf("export %s=\"%s\"\n", key, value)); err != nil {
						panic(err)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

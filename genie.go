package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v2"
	"github.com/imdario/mergo"
)

type command struct {
	Executable string `yaml:"executable"`
	Arguments string `yaml:"arguments"`
}

type configFile struct {
	Commands map[string][]command `yaml:"commands"`
}

func findCommandFiles() []string {
	var discoveredFiles []string
	path, err := os.Getwd()
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Base(path) == "commands.yaml" {
			discoveredFiles = append(discoveredFiles, path)
		}

		return nil
	})
	if err != nil {
		log.Println(err)
	}

	return discoveredFiles
}

func (c *configFile) getConf() *configFile {
	files := findCommandFiles()
	fmt.Println(files)
	for k := range files {
		yamlFile, err := ioutil.ReadFile(files[k])
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		tempStruct := configFile{}
		err = yaml.Unmarshal(yamlFile, &tempStruct)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}

		mergo.Merge(c, tempStruct)
	}

	return c
}

func (c *configFile) getAvailableCommands() {
	fmt.Println("Available commands:")
	for k := range c.Commands {
		fmt.Println("\t" + k)
	}
}

func main() {
	var dryRun bool
	flag.BoolVar(&dryRun, "dry", false, "Run in dry mode to print commands")
	flag.Parse()

	var c configFile
	c.getConf()
	fmt.Println(c.Commands)
	args := os.Args[1:]
	if 0 == len(args) {
		c.getAvailableCommands()
		os.Exit(0)
	}

	commandName := flag.Args()[0]
	for key := range c.Commands[commandName] {
		if dryRun == true {
			fmt.Printf(">\t%s %s\n", c.Commands[commandName][key].Executable, c.Commands[commandName][key].Arguments)
		} else {
			executable := c.Commands[commandName][key].Executable
			arguments := c.Commands[commandName][key].Arguments
			cmd := exec.Command(executable, arguments)
			out, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatalf("Command failed with %s\n", err)
			}
			fmt.Printf("%s", string(out))
		}
	}
}

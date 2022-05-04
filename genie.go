package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/imdario/mergo"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type environment struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type command struct {
	Command     string        `yaml:"command"`
	Environment []environment `yaml:"environment"`
}

type configFile struct {
	Commands map[string][]command `yaml:"commands"`
	Shell    string               `yaml:"shell"`
}

func findCommandFiles() []string {
	var discoveredFiles []string
	homeDirectory, _ := homedir.Dir()
	if _, err := os.Stat(homeDirectory + "/.genie-commands.yaml"); !errors.Is(err, os.ErrNotExist) {
		discoveredFiles = append(discoveredFiles, homeDirectory+"/.genie-commands.yaml")
	}

	path, err := os.Getwd()
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Base(path) == "genie.yaml" {
			discoveredFiles = append(discoveredFiles, path)
		}

		return nil
	})
	if err != nil {
		log.Println(err)
	}

	return discoveredFiles
}

func initCommandsFile() (string, error) {
	path, _ := os.Getwd()
	targetFile := path + "/genie.yaml"

	_, err := os.Stat(targetFile)
	if err == nil {
		return "", errors.New("genie.yaml file already exists at " + path)
	}

	commandMap := make(map[string][]command)
	commandMap["example"] = []command{
		{Command: "echo this is an example command"},
	}
	t := configFile{
		Shell:    "/bin/bash",
		Commands: commandMap,
	}
	content, _ := yaml.Marshal(&t)
	err = ioutil.WriteFile(targetFile, []byte(content), 0644)
	if err != nil {
		log.Fatalf("Unable to create file " + path)
	}

	return path, nil
}

func (c *configFile) getConf() *configFile {
	files := findCommandFiles()
	for _, path := range files {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			continue
		}
		yamlFile, err := ioutil.ReadFile(path)
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
	fmt.Println("\tinit - create a commands.yaml file in the current directory")
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
	args := os.Args[1:]
	if len(args) == 0 {
		c.getAvailableCommands()
		os.Exit(0)
	}

	commandName := flag.Args()[0]

	if commandName == "init" {
		path, err := initCommandsFile()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("Created commands file at " + path)
	}

	for key := range c.Commands[commandName] {
		if dryRun {
			fmt.Printf(">\t%s\n", c.Commands[commandName][key].Command)
		} else {
			command := c.Commands[commandName][key].Command
			env := c.Commands[commandName][key].Environment
			cmd := exec.Command("bash", "-c", command)
			cmd.Env = append(os.Environ(), envBuilder(env)...)
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("%s", string(out))
				log.Fatalf("Command failed with %s\n", err)
			}
			fmt.Printf("%s", string(out))
		}
	}
}

func envBuilder(env []environment) []string {
	var result []string
	for _, row := range env {
		result = append(result, fmt.Sprintf("%s=%s\n", row.Name, row.Value))
	}

	return result
}

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type command struct {
	Executable string `yaml:"executable"`
	Arguments string `yaml:"arguments"`
}

type configFile struct {
	Commands map[string][]command `yaml:"commands"`
}

func (c *configFile) getConf() *configFile {

	yamlFile, err := ioutil.ReadFile("commands.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
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
	args := os.Args[1:]
	if 0 == len(args) {
		c.getAvailableCommands()
		os.Exit(0)
	}

	commandName := flag.Args()[0]
	// fmt.Printf("%+v\n", c.Commands)
	// fmt.Printf("%+v\n", c.Commands[commandName][0])
	for key := range c.Commands[commandName] {
		if dryRun == true {
			fmt.Printf(">\t%s %s\n", c.Commands[commandName][key].Executable, c.Commands[commandName][key].Arguments)
		} else {
			executable := c.Commands[commandName][key].Executable
			arguments := c.Commands[commandName][key].Arguments
			cmd := exec.Command(executable, arguments)
			out, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatalf("cmd.Run() failed with %s\n", err)
			}
			fmt.Printf("%s", string(out))
		}
	}
}

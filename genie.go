package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	// "os/exec"

	"gopkg.in/yaml.v2"
)

type conf struct {
	Commands map[string][]string `yaml:"commands"`
}

func (c *conf) getConf() *conf {

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

func (c *conf) getAvailableCommands() {
	fmt.Println("Available commands:")
	for k := range c.Commands {
		fmt.Println("\t" + k)
	}
}

func main() {
	var dryRun bool
	flag.BoolVar(&dryRun, "dry", false, "Run in dry mode to print commands")
	flag.Parse()

	var c conf
	c.getConf()
	args := os.Args[1:]
	if 0 == len(args) {
		c.getAvailableCommands()
		os.Exit(0)
	}

	fmt.Println(c)
	fmt.Println(c.Commands[args[0]])
	fmt.Println(dryRun)
	fmt.Println(flag.Args())
	// for key := range c.Commands[args[0]] {
	// 	if *dryMode == true {
	// 		fmt.Println(c.Commands[args[0]][key])
	// 	} else {
	// 		cmd := exec.Command(c.Commands[args[0]][key])
	// 		out, err := cmd.CombinedOutput()
	// 		if err != nil {
	// 			log.Fatalf("cmd.Run() failed with %s\n", err)
	// 		}
	// 		fmt.Printf("combined out:\n%s\n", string(out))
	// 	}
	// }
}

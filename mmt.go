package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)

type DbConfig struct {
	host   string
	port   string
	user   string
	pass   string
	schema string
}

type Table struct {
	name string
}

type Config struct {
	dbconfig DbConfig
	tables   []Table
}

func main() {
	app := cli.NewApp()
	app.Name = "mmt"
	app.Usage = "metadata management tool"
	app.Commands = []cli.Command{
		{
			Name:      "dump",
			ShortName: "d",
			Usage:     "Dump tables",
			Action: func(c *cli.Context) {
				do_dump()
			},
		},
		{
			Name:      "restore",
			ShortName: "r",
			Usage:     "Restore tables",
			Action: func(c *cli.Context) {
				do_restore()
			},
		},
		{
			Name:      "info",
			ShortName: "i",
			Usage:     "Show current configuration information",
			Action: func(c *cli.Context) {
				do_info()
			},
		},
	}
	app.Action = func(c *cli.Context) {
		println("TODO")
	}
	app.Run(os.Args)
}

func do_info() {
	config := ReadConfig()
	println("Host: " + config.dbconfig.host)
}

func ReadConfig() Config {
	file, readErr := ioutil.ReadFile("mmt.config")
	if readErr != nil {
		panic(readErr)
	}
	var config Config
	parseErr := json.Unmarshal([]byte(file), &config)
	if parseErr != nil {
		panic(parseErr)
	}
	fmt.Printf("%#v\n", config)
	return config
}

func get_binary() string {
	//Check that the mysql binary exists on this machine
	binary, lookErr := exec.LookPath("mysql")
	if lookErr != nil {
		panic(lookErr)
	}
	return binary
}

func do_dump() {
	args := []string{"mysql", "-u", "root", "-p"}
	env := os.Environ()
	execErr := syscall.Exec(get_binary(), args, env)
	if execErr != nil {
		panic(execErr)
	}
}

func do_restore() {
	println("TODO")
}

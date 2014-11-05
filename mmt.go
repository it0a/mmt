package main

import (
	"encoding/json"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)

type DbConfig struct {
	Host   string `json:"host"`
	Port   string `json:"port"`
	User   string `json:"user"`
	Schema string `json:"schema"`
}

type Table struct {
	Name string `json:"name"`
}

type TableProfile struct {
	Name   string  `json:"name"`
	Tables []Table `json:"tables"`
}

type DbProfile struct {
	Name     string   `json:"name"`
	DbConfig DbConfig `json:"dbConfig"`
}

type Config struct {
	DbProfiles    []DbProfile    `json:"dbProfiles"`
	TableProfiles []TableProfile `json:"tableProfiles"`
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
			Usage:     "Prints the current configuration information",
			Action: func(c *cli.Context) {
				print_info()
			},
		},
	}
	app.Action = func(c *cli.Context) {
		println("TODO")
	}
	app.Run(os.Args)
}

func print_info() {
	Config := read_config()
	for _, dbProfile := range Config.DbProfiles {
		println("Database Profile Name: " + dbProfile.Name)
		println("Host: " + dbProfile.DbConfig.Host)
		println("Port: " + dbProfile.DbConfig.Port)
		println("User: " + dbProfile.DbConfig.User)
		println("Schema: " + dbProfile.DbConfig.Schema)
		println("")
	}
	for _, tableProfile := range Config.TableProfiles {
		println("Table Profile Name: " + tableProfile.Name)
		println("Tables: ")
		for _, table := range tableProfile.Tables {
			println(table.Name)
		}
		println("")
	}
}

func read_config() Config {
	file, readErr := ioutil.ReadFile("./mmt.json")
	if readErr != nil {
		panic(readErr)
		os.Exit(1)
	}
	var config Config
	parseErr := json.Unmarshal(file, &config)
	if parseErr != nil {
		panic(parseErr)
		os.Exit(1)
	}
	return config
}

func get_binary() string {
	//Check that the mysql binary exists on this machine
	binary, lookErr := exec.LookPath("mysql")
	if lookErr != nil {
		panic(lookErr)
		os.Exit(1)
	}
	return binary
}

func validate_connection(dbconfig DbConfig) bool {
	retVal := false
	args := []string{}
	args = append(args, "-u", "root")
	args = append(args, "-p")
	command := exec.Command("mysql", args...)
	err := command.Run()
	if err == nil {
		retVal = true
	}
	return retVal
}

func do_dump() {
	args := []string{}
	args = append(args, "-u", "root")
	args = append(args, "-p")
	command := exec.Command("mysql", args...)
	err := command.Run()
	if err != nil {
		panic(err)
	}
	println("OK!")
}

func do_restore() {
	config := read_config()
	if validate_connection(config.DbProfiles[0].DbConfig) == true {
		println("Validated!")
	} else {
		println("Failed to validate")
	}
}

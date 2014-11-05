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
	Config := ReadConfig()
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

func ReadConfig() Config {
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
	}
	return binary
}

func do_dump() {
	args := []string{"mysql", "-u", "root", "-p"}
	env := os.Environ()
	execErr := syscall.Exec(get_binary(), args, env)
	if execErr != nil {
		panic(execErr)
		os.Exit(1)
	}
}

func do_restore() {
	println("TODO")
}

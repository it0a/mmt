package main

import (
	"encoding/json"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
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
				config := read_config()
				if validate_connection(config.DbProfiles[0]) {
					do_dump(config.DbProfiles[0], config.TableProfiles[0])
				}
			},
		},
		{
			Name:      "restore",
			ShortName: "r",
			Usage:     "Restore tables",
			Action: func(c *cli.Context) {
				config := read_config()
				if validate_connection(config.DbProfiles[0]) {
					do_restore(config)
				}
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

func build_args(dbProfile DbProfile) []string {
	args := []string{}
	args = append(args, "-h", dbProfile.DbConfig.Host)
	args = append(args, "-P", dbProfile.DbConfig.Port)
	args = append(args, "-u", dbProfile.DbConfig.User)
	args = append(args, "-p")
	return args
}

func exec_mysql(dbProfile DbProfile) error {
	args := build_args(dbProfile)
	args = append(args, dbProfile.DbConfig.Schema)
	return exec.Command("mysql", args...).Run()
}

func dump_table(dbProfile DbProfile, table Table) error {
	out, err := os.OpenFile(path.Join("./", "dump.sql"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	args := build_args(dbProfile)
	args = append(args, dbProfile.DbConfig.Schema, table.Name)
	println(table.Name)
	command := exec.Command("mysqldump", args...)
	command.Stdout = out
	return command.Run()
}

func validate_connection(dbProfile DbProfile) bool {
	retVal := false
	err := exec_mysql(dbProfile)
	if err == nil {
		retVal = true
	} else {
		println("Failed to validate connection for database profile " + dbProfile.Name)
	}
	return retVal
}

func do_dump(dbProfile DbProfile, tableProfile TableProfile) {
	println("Doing the dump...")
	dump_table(dbProfile, tableProfile.Tables[0])
}

func do_restore(config Config) {
	println("Doing the restore...")
}

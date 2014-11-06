package main

import (
	"code.google.com/p/gopass"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
)

type DbConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Schema   string `json:"schema"`
	Password string
}

type Table struct {
	Name string `json:"name"`
}

type TableProfile struct {
	Name    string  `json:"name"`
	DumpDir string  `json:"dumpDir"`
	Tables  []Table `json:"tables"`
}

type DbProfile struct {
	Name     string   `json:"name"`
	DbConfig DbConfig `json:"dbConfig"`
	Valid    bool
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
				dbProfile, tableProfile := initialize()
				if dbProfile.Valid {
					do_dump(dbProfile, tableProfile)
				}
			},
		},
		{
			Name:      "restore",
			ShortName: "r",
			Usage:     "Restore tables",
			Action: func(c *cli.Context) {
				dbProfile, tableProfile := initialize()
				if dbProfile.Valid {
					do_restore(dbProfile, tableProfile)
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

func read_config() Config {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	file, readErr := ioutil.ReadFile(usr.HomeDir + "/.mmt.json")
	if readErr != nil {
		panic(readErr)
	}
	var config Config
	parseErr := json.Unmarshal(file, &config)
	if parseErr != nil {
		fmt.Println("Error parsing .mmt.json")
		fmt.Println(parseErr)
		os.Exit(1)
	}
	return config
}

func build_args(dbProfile DbProfile) []string {
	args := []string{}
	args = append(args, "-h", dbProfile.DbConfig.Host)
	args = append(args, "-P", dbProfile.DbConfig.Port)
	args = append(args, "-u", dbProfile.DbConfig.User)
	args = append(args, "-p"+dbProfile.DbConfig.Password)
	return args
}

func exec_mysql(dbProfile DbProfile) error {
	args := build_args(dbProfile)
	args = append(args, dbProfile.DbConfig.Schema)
	return exec.Command("mysql", args...).Run()
}

func detect_diff(dumpDir string, table Table) {
	args := []string{}
	args = append(args, "--git-dir", dumpDir+"/.git")
	args = append(args, "--work-tree", dumpDir)
	args = append(args, "diff", "--quiet", dumpDir+"/"+table.Name+".sql")
	command := exec.Command("git", args...)
	err := command.Run()
	if err != nil {
		fmt.Println("Diff detected!")
	}
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

func initialize() (DbProfile, TableProfile) {
	config := read_config()
	dbProfile := config.DbProfiles[0]
	tableProfile := config.TableProfiles[0]
	password, err := gopass.GetPass("Enter password: ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbProfile.DbConfig.Password = password
	dbProfile.Valid = validate_connection(dbProfile)
	return dbProfile, tableProfile
}

func print_dbProfile(dbProfile DbProfile) {
	println("Database Profile Name: " + dbProfile.Name)
	println("Host: " + dbProfile.DbConfig.Host)
	println("Port: " + dbProfile.DbConfig.Port)
	println("User: " + dbProfile.DbConfig.User)
	println("Schema: " + dbProfile.DbConfig.Schema)
	println("")
}

func print_tableProfile(tableProfile TableProfile) {
	println("Table Profile Name: " + tableProfile.Name)
	println("Tables: ")
	for _, table := range tableProfile.Tables {
		print_table(table)
	}
	println("")
}

func print_table(table Table) {
	println(table.Name)
}

func print_info() {
	Config := read_config()
	for _, dbProfile := range Config.DbProfiles {
		print_dbProfile(dbProfile)
	}
	for _, tableProfile := range Config.TableProfiles {
		print_tableProfile(tableProfile)
	}
}

func dump_table(dbProfile DbProfile, dumpDir string, table Table) {
	println("Dumping " + table.Name + "...")
	out, err := os.OpenFile(path.Join(dumpDir, table.Name+".sql"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	args := build_args(dbProfile)
	args = append(args, "--extended-insert=FALSE")
	args = append(args, "--skip-comments")
	args = append(args, dbProfile.DbConfig.Schema, table.Name)
	command := exec.Command("mysqldump", args...)
	command.Stdout = out
	commandErr := command.Run()
	if commandErr != nil {
		fmt.Println(commandErr)
		os.Exit(1)
	}
	detect_diff(dumpDir, table)
}

func do_dump(dbProfile DbProfile, tableProfile TableProfile) {
	for _, table := range tableProfile.Tables {
		dump_table(dbProfile, tableProfile.DumpDir, table)
	}
}

func restore_table(dbProfile DbProfile, dumpDir string, table Table) {
	println("Restoring " + table.Name + "...")
	in, err := os.Open(dumpDir + "/" + table.Name + ".sql")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	args := build_args(dbProfile)
	args = append(args, dbProfile.DbConfig.Schema)
	command := exec.Command("mysql", args...)
	command.Stdin = in
	execErr := command.Run()
	if execErr != nil {
		fmt.Println(execErr)
		os.Exit(1)
	}
}

func do_restore(dbProfile DbProfile, tableProfile TableProfile) {
	for _, table := range tableProfile.Tables {
		restore_table(dbProfile, tableProfile.DumpDir, table)
	}
}

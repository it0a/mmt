package main

import (
	"github.com/codegangsta/cli"
	"os"
	"os/exec"
	"syscall"
)

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
	}
	app.Action = func(c *cli.Context) {
		println("TODO")
	}
	app.Run(os.Args)
}

func get_mysql() string {
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
	execErr := syscall.Exec(get_mysql(), args, env)
	if execErr != nil {
		panic(execErr)
	}
}

func do_restore() {
	println("TODO")
}

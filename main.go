package main

import (
   "fmt"
	"os"
	cmd"wgsl/cmds"
)
var available_commands = []string{"init","help","get","train","test","result"}
func contains(cmmd string) bool {
	for _,val := range available_commands {
		if cmmd == val {
			return true
		}
	}
	return false
}
func wgslSucess() bool {
	file := ".wgsl"
	_,err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
func main() {
	fmt.Println("Thank you for using the wgsl command line tool the one stop for an prediction of leukemia detection")
	fmt.Println("Disclaimer:\n1)This cli tool is intended for medical professinals to aid them to make a concrete decision\n2)The results are also shared to the patient only after an hour of the results sent to the doctor")
	fmt.Println("use wgsl help command to see all the supported commands")
	if len(os.Args) < 2 {
		panic("Not enough command to perform the operations. Please follow the structure wgsl <commandname>\n<commandname> is replaced by the commands\n1)init\n2)help\n3)get\n4)train\n5)test\n6)result")
	}
	cmmd := os.Args[1]
	if !contains(cmmd) {
		panic("Unrecognised command,please use help command to see available commands")
	}
	switch cmmd {
	case "init":
		if wgslSucess() {
			fmt.Println("already initialized in this folder")
			return
		}
		cmd.Init()
	case "help":
		 cmd.Help()
	}
}

package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Required Minimum 2 Args")
		return
	}
	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "install":
	case "uninstall":
	case "get":
	case "push":
	case "init" :
	}
		


}
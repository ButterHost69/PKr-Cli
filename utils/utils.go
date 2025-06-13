package utils

import (
	"math/rand"
	"os"
	"os/exec"
	"runtime"

	fake "github.com/brianvoe/gofakeit/v7"
)

func ClearScreen() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func CreateSlug() string {
	var gamerTag []string
	for range 1024 {
		gamerTag = append(gamerTag, fake.Gamertag())
	}
	g := rand.Intn(1024)
	return gamerTag[g]
}

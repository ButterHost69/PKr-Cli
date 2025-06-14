package utils

import (
	"fmt"
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

func PrintProgressBar(progress int, total int, barLength int) {
	percent := float64(progress) / float64(total)
	hashes := int(percent * float64(barLength))
	spaces := barLength - hashes

	fmt.Printf("\r[%s%s] %.2f%%",
		repeat("#", hashes),
		repeat(" ", spaces),
		percent*100)
}

func repeat(char string, count int) string {
	result := ""
	for range count {
		result += char
	}
	return result
}

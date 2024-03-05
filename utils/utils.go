package utils

import (
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

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

func CreateOTP(size int) int32 {
	mask := strings.Repeat("9", size)
	intmask, _ := strconv.Atoi(mask)
	return rand.Int31n(int32(intmask))
}

func CreateSlug() string {
	var gamerTag []string
	for i := 0; i < 1024; i++ {
		gamerTag = append(gamerTag, fake.Gamertag())
	}
	// fmt.Println(gamerTag)
	g := rand.Intn(1024)
	// fmt.Println(g)
	return gamerTag[g]
}

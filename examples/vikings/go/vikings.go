package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	// Example with a config file:
	//cmd := exec.Command("dosbox", "-conf", "my_dosbox.conf")

	// Example with commands (using -c):

	// Capture output (optional)
	cmd := exec.Command("dosbox", "-c", "mount c c:\\Games\\DosGames\\lost-vikings", "-c", "c:", "-c", "VIKINGS.EXE")
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DOSBox started with config/commands!")
}

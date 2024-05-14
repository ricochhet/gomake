package process

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func Exec(commands []string) error {
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "bash"
		flag = "-c"
	}

	for _, cmd := range commands {
		fmt.Println("gomake: executing command:", cmd)
		cmd := exec.Command(shell, flag, cmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func Errr(err error) {
	fmt.Println("gomake: " + err.Error())
}

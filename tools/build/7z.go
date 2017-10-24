package build

import (
	"os"
	"os/exec"
)

func SevenZipBuild(folder string, targetFile string, password string) error {
	args := []string{
		"a", "-p" + password, "-mhe=on", targetFile, folder,
	}

	cmd := exec.Command("7z", args...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		os.Stdout.Write(output)
	}

	return err
}

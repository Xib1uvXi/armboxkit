package shutil

import (
	"bytes"
	"os/exec"
	"strings"
)

type CmdExecutor struct {
}

func NewCmdExecutor() *CmdExecutor {
	return &CmdExecutor{}
}

func (e *CmdExecutor) ExecuteCmd(cmdStr string) (outStr string, errStr string, err error) {
	var (
		stdout, stderr bytes.Buffer
	)

	fields := strings.Fields(cmdStr)
	name := fields[0]

	var cmd *exec.Cmd

	if len(fields) == 1 {
		cmd = exec.Command(name)
	} else {
		args := fields[1:]
		cmd = exec.Command(name, args...)
	}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	outStr, errStr = stdout.String(), stderr.String()
	return outStr, errStr, err
}

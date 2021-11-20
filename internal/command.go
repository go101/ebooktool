package internal

import (
	"context"
	"log"
	"os"
	"os/exec"
	"time"

	"go101.org/ebooktool/internal/nstd"
)

func ExecCommand(timeout time.Duration, wd string, envs []string, cmdAndArgs ...string) ([]byte, error) {
	if len(cmdAndArgs) == 0 {
		panic("command is not specified")
	}

	if wd == "" {
		var err error
		wd, err = os.Getwd()
		if err != nil {
			log.Println(`Can't get current path. Set it as "."`)
			//wd = "."
			wd = ""
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	command := exec.CommandContext(ctx, cmdAndArgs[0], cmdAndArgs[1:]...)
	command.Dir = wd
	command.Env = envs
	return command.CombinedOutput() // ToDo: maybe it is better not to combine
}

func CommandOutputAsString(timeout time.Duration, wd string, envs []string, cmdAndArgs ...string) string {
	output, err := ExecCommand(timeout, wd, envs, cmdAndArgs...)
	if err != nil {
		return ""
	}
	return nstd.String(output).TrimSpace().String()
}

func GetVersionAndDateFromGit(dir string) (version string, date string) {
	tag := CommandOutputAsString(time.Second*32, dir, nil, "git", "describe", "--tags", "--abbrev=0")
	var rev string
	if tag != "" {
		rev = CommandOutputAsString(time.Second*32, dir, nil, "git", "rev-parse", tag)
		if len(rev) > 7 {
			rev = rev[:7]
		}
	} else {
		rev = CommandOutputAsString(time.Second*32, dir, nil, "git", "rev-parse", "HEAD")
	}
	date = CommandOutputAsString(time.Second*32, dir, nil, "git", "log", "-1", "--pretty=%ad", "--date=format:%Y/%m/%d", rev)
	if len(rev) > 7 {
		rev = rev[:7]
	}

	dirty := "" != CommandOutputAsString(time.Second*32, dir, nil, "git", "status", "-s")
	if dirty {
		defer func() {
			version += "-!"
		}()
	}

	if rev == "" {
		return tag, date
	}
	if tag == "" {
		if rev == "" {
			return "", date
		}
		return "rev-" + rev, date
	}
	return tag + "-rev-" + rev, date
}

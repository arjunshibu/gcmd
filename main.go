package main

import (
    "os"
    "os/user"
    "os/exec"
    "fmt"
    "flag"
    "errors"
    "strings"
    "encoding/json"
    "path/filepath"
)

type Cmd struct {
    Name    string  `json:"name,omitempty"`
    Args    string  `json:"args,omitempty"`
}

type Cmds struct {
    Cmds    []Cmd   `json:"cmds,omitempty"`
    Stdin   bool    `json:"stdin,omitempty"`
}

func main() {
    var saveMode, removeMode, echoMode, listMode, stdinMode bool

    flag.BoolVar(&saveMode, "save", false, "save a command")
    flag.BoolVar(&removeMode, "rm", false, "remove a command")
    flag.BoolVar(&echoMode, "echo", false, "prints the command rather than executing it")
    flag.BoolVar(&listMode, "ls", false, "list available commands")
    flag.BoolVar(&stdinMode, "i", false, "take input from stdin (optional, for -save only)")

    flag.Parse()

    if listMode {
	cmds, err := getCmds()

        if err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            return
        }

        fmt.Println(strings.Join(cmds, "\n"))
        return
    }

    dir, err := getCmdsDir()

    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        return
    }

    cmdName := flag.Arg(0)
    path := filepath.Join(dir, cmdName + ".json")

    if removeMode {
        err := os.Remove(path)

        if err != nil {
            fmt.Fprintln(os.Stderr, "No such command")
            return
        }

        fmt.Println("Command removed")
        return
    }

    var cmds Cmds

    if saveMode {
        cmd, err := parseCmd(flag.Arg(1))

        if err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            return
        }

	cmds = packCmds(cmd)

        if stdinMode {
            cmds.Stdin = true
        }

        file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)

	if err != nil {
	    fmt.Fprintf(os.Stderr, "Failed to create command file %s\n", path)
            return
	}

	defer file.Close()

        encoder := json.NewEncoder(file)
        encoder.SetIndent("", "   ")
        err = encoder.Encode(&cmds)

        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed to write command file %s\n", path)
            return
        }

	fmt.Println("Command saved")
        return
    }

    file, err := os.Open(path)

    if err != nil {
        if len(flag.Args()) == 0 {
            fmt.Fprintln(os.Stderr, "Provide a command")
            return
        }

        fmt.Fprintln(os.Stderr, "No such command")
        return
    }

    defer file.Close()

    decoder := json.NewDecoder(file)
    err = decoder.Decode(&cmds)

    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to parse command file %s\n", path)
        return
    }

    var cmdStr string

    for _, cmd := range cmds.Cmds {
        cmdStr += cmd.Name + " " + cmd.Args + " | "
    }

    cmdStr = cmdStr[:len(cmdStr) - 3]

    if echoMode {
        fmt.Println(cmdStr)
        return
    }

    fi, err := os.Stdin.Stat()

    if err != nil {
        panic(err)
    }

    execCmd := exec.Command("bash", "-c", cmdStr)
    execCmd.Stdout = os.Stdout
    execCmd.Stderr = os.Stderr

    if cmds.Stdin {
        if fi.Mode() & os.ModeNamedPipe == 0 {
            fmt.Fprintln(os.Stderr, "Command needs stdin")
            return
        } else {
            execCmd.Stdin = os.Stdin
        }
    }

    if err := execCmd.Run(); err != nil {
        fmt.Fprintln(os.Stderr, string(err.Error()))
        return
    }
}

func getCmdsDir() (string, error) {
    user, err := user.Current()

    if err != nil {
	return "", err
    }

    path := filepath.Join(user.HomeDir, ".config/gcmd")

    if _, err := os.Stat(path); !os.IsNotExist(err) {
	return path, nil
    }

    return filepath.Join(user.HomeDir, ".gcmd"), nil
}

func getCmds() ([]string, error) {
    dir, err := getCmdsDir()
    out := []string{}

    if err != nil {
        return out, err
    }

    files, err := filepath.Glob(dir + "/*.json")

    if err != nil {
        return out, err
    }

    for _, file := range files {
        file = file[len(dir) + 1 : len(file) - 5]
        out = append(out, file)
    }

    return out, nil
}

func parseCmd(cmd string) ([]string, error) {
    cmds := strings.Split(cmd, " | ")

    for i, cmd := range cmds {
        cmds[i] = strings.TrimSpace(cmd)
    }

    if len(flag.Args()) == 0 {
        return []string{}, errors.New("Name cannot be empty")
    }

    if len(cmds) == 1 && cmds[0] == "" {
        return []string{}, errors.New("Command cannot be empty")
    }

    return cmds, nil
}

func packCmds(cmd []string) (Cmds) {
    var cmds Cmds
    cmds.Cmds = make([]Cmd, len(cmd))

    for i, c := range cmd {
	cmdArr := strings.SplitN(c, " ", 2)
        cmds.Cmds[i].Name = cmdArr[0]
        if len(cmdArr) == 1 {
	    cmds.Cmds[i].Args = ""
        } else {
            cmds.Cmds[i].Args = cmdArr[1]
        }
    }

    return cmds
}

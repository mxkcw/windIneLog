package windIne_cmd

import (
	"bytes"
	"github.com/mxkcw/windIne/windIne_encoding"
	"github.com/mxkcw/windIne/windIne_string"
	"os/exec"
	"runtime"
	"sync"
)

var rc_wg sync.WaitGroup
var mutex sync.Mutex // 添加一个互斥锁来保护普通map

type WindIneCmd struct {
	results map[string]string
}

func RunWith(CommandMap map[string]string) map[string]string {
	gcmd := &WindIneCmd{
		results: make(map[string]string),
	}

	for key, command := range CommandMap {
		rc_wg.Add(1)
		go gcmd.execute(key, command)
	}
	rc_wg.Wait()
	return gcmd.results
}

func (gcmd *WindIneCmd) execute(key string, command string) {
	defer rc_wg.Done()
	var cmd *exec.Cmd
	switch os := runtime.GOOS; os {
	case "windows":
		git_bash_path := GetWindowsGitBashPath()
		if git_bash_path == "" {
		} else {
			cmd = exec.Command(git_bash_path, "-c", command)
		}
	case "darwin":
		cmd = exec.Command("/bin/zsh", "-c", command)
	default:
		cmd = exec.Command("/bin/bash", "-c", command)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := out.String()
	if err != nil {
		result = stderr.String()
	}

	if runtime.GOOS == "windows" {
		result, _ = windIne_encoding.ConvertToUTF8UsedLocalENV(result)
	}

	// 使用mutex保护对results map的写入
	mutex.Lock()
	windIne_string.DelStringEndNewlines(&result)
	gcmd.results[key] = result
	mutex.Unlock()
}

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	flag.Parse()
	args := flag.Args()
	prefix := ""
	if 0 < len(args) {
		prefix = args[0]
	}

	relPath := "../"

	root, err := filepath.Abs(relPath)
	if err != nil {
		panic(err)
	}

	files, err := ioutil.ReadDir(root)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}

	for _, file := range files {
		if file.IsDir() {
			// プレフィクス文字列が指定されている場合は一致したディレクトリのみ対象とする
			if prefix != "" && !strings.HasPrefix(file.Name(), prefix) {
				continue
			}
			wg.Add(1)
			fmt.Println(file.Name())
			go func(fileName string) {
				out, err := exec.Command("git", "-C", relPath+"./"+fileName, "pull", "--prune").Output()
				if err != nil {
					fmt.Printf("\u001B[31m[FAIL]\u001B[0m %s\n%v\n", fileName, err)
				} else {
					fmt.Printf("\u001B[34m[SUCCESS]\u001B[0m %s\n%s", fileName, string(out))
				}
				defer wg.Done()
			}(file.Name())

		}
	}
	wg.Wait()
}

package main

import (
	"os"
	"path/filepath"
	"strings"
)

func main() {
	/*app := InitializeApp()
	err := app.Run(":9080")
	if err != nil {
		println("运行异常!", err.Error())
		return
	}*/
	filePaths := scan("F:\\maven", 10)
	if filePaths != nil {
		for i := range filePaths {
			println(filePaths[i])
		}
	}
	println(len(filePaths))
}

func scan(path string, maxDepth int) []string {
	var res []string
	rootDepth := strings.Count(path, string(os.PathSeparator))
	if path == "" || path == "." || path == "."+string(os.PathSeparator) {
		rootDepth = 0 // 当前目录的特殊处理
	}
	err := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {

		if err != nil {
			println(err)
		}
		depth := strings.Count(path, string(os.PathSeparator)) - rootDepth
		if depth > maxDepth {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.IsDir() {
			res = append(res, path)
		}
		return nil
	})

	println(len(res))
	if err != nil {
		return nil
	}
	return res
}

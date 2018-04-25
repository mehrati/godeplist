package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"gopkg.in/fatih/set.v0"
)

func main() {
	dirs := GetAllDirs()
	files := GetAllFiles(dirs)
	err := ExecGoImports(files)
	if err != nil {
		fmt.Println(err)
	}
	deps := GetDeps(files)
	ShowDeps(deps)
}
// TODO
func GetNonStdDeps(deps []string) []string{

	return nil
}

func GetDeps(files []string) []string {
	RES := set.New()
	for _, file := range files {
		fo, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer fo.Close()
		scanner := bufio.NewScanner(fo)
	br:
		for scanner.Scan() {
			t := scanner.Text()
			if strings.HasPrefix(t, "import(") || strings.HasPrefix(t, "import (") {
				RES.Add(TrimImport(t))
				for scanner.Scan() {
					t = scanner.Text()
					if strings.Contains(t, ")") {
						break br
					}
					RES.Add(TrimImport(t))
				}
			} else if strings.HasPrefix(t, "import") {
				RES.Add(TrimImport(t))
				break br
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	listi := RES.List()
	list := []string{}
	for _, l := range listi {
		s := l.(string)
		if s != "" && s != "import (" && s != "import(" {
			list = append(list, TrimImport(s))
		}
	}
	RES.Clear()
	sort.Strings(list)
	return list
}

func GetAllDirs() []string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	listdirs := []string{wd}
	i := 0
	for {
		dirs, err := GetListDir(listdirs[i])
		if err != nil {
			break
		}
		for _, dir := range dirs {
			listdirs = append(listdirs, dir)
		}
		i++
		if i == len(listdirs) {
			break
		}
	}
	return listdirs
}

func GetListDir(dir string) ([]string, error) {
	listdir := []string{}
	list, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, l := range list {
		if l.IsDir() {
			listdir = append(listdir, dir+"/"+l.Name())
		}
	}
	return listdir, nil
}

func GetListFile(dir string) ([]string, error) {
	listfile := []string{}
	list, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, l := range list {
		if !l.IsDir() && strings.HasSuffix(l.Name(), ".go") {
			listfile = append(listfile, dir+"/"+l.Name())
		}
	}
	return listfile, nil
}

func GetAllFiles(dirs []string) []string {
	files := []string{}
	for _, d := range dirs {
		file, err := GetListFile(d)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range file {
			files = append(files, f)
		}
	}
	return files
}

func ShowDeps(deps []string) {
	for _, d := range deps {
		fmt.Println(d)
	}
}

func ExecGoImports(files []string) error {
	path, err := exec.LookPath("goimports")
	if err != nil {
		return err
	}
	for _, f := range files {
		err = exec.Command(path, "-w", f).Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func TrimImport(str string) string {
	if strings.Contains(str, "\"") {
		s := ""
		s = strings.TrimLeft(str, str[:strings.Index(str, "\"")+1])
		s = strings.TrimRight(s, str[strings.LastIndex(str, "\""):])
		return s
	} else {
		return str
	}
}
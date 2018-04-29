package main

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"gopkg.in/fatih/set.v0"
)

func main() {
	s := set.New()
	dirs := GetAllDirs()
	for _, d := range dirs {
		imports, _ := GetImport(d)
		for _,i := range imports {
			s.Add(i)
		}
	}
	list := s.List()
	ShowImports(list)

}

func GetImport(path string) ([]string, error) {
	pkg, err := build.ImportDir(path, 0)
	if err != nil {
		return nil, err
	}
	return pkg.Imports, nil
}

// TODO
func GetNonStdDeps(deps []string) []string {

	return nil
}

func GetAllDirs() []string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 1 {
		wd = os.Args[1]
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

func ShowImports(imp []interface{}) {
	for _, d := range imp {
		fmt.Println(d)
	}
}

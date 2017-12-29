package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"log"
	"path/filepath"
)

var files = make(map[string]FileDef)
var doubleFiles = make(map[string]FileDef)

var rootPath = flag.String("root-path", ".", "Root path for scan")
var showDoubles = flag.Bool("show-doubles", true, "Show list of double files")
var maxFileSize = flag.Int64("max-filesize", 10, "Max size of file to be scanned in MB, default 10 MB")
var deleteDoubles = flag.Bool("delete-doubles", false, "Delete found doubles from file system")

type FileDef struct {
	Path string
	Hash []byte
}

func walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Fatal("Read error ", err)
	}
	if !info.IsDir() {
		hash := sha256.New()
		if info.Size() > (*maxFileSize) * 1024 * 1024 {
			fmt.Printf("\nSkip file '%s' due to size\n", path)
			return nil
		}
		file, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal("Read error ", err)
		}
		hash.Write(file)
		shaHash := hash.Sum(nil)
		if (*showDoubles) {
			if doubleFilePath := files[fmt.Sprintf("%x", shaHash)].Path; doubleFilePath != "" {
				fmt.Printf("\n%s \t -> \t %s\n", doubleFilePath, path)
				doubleFiles[fmt.Sprintf("%x", shaHash)] = FileDef{path, shaHash}
				if (*deleteDoubles) {
					delete(path)
				}
			}
		} else {
			fmt.Printf("\n Indexing %s ", path)
		}
		files[fmt.Sprintf("%x", shaHash)] = FileDef{path, shaHash}
		fmt.Printf("\r %d scanned ", len(files))
	}
	return nil
}

func delete(path string) {
	fmt.Println("Deleting",path)
	os.Remove(path)
}

func init() {
	flag.Parse()
}
func main() {
	err := filepath.Walk(*rootPath, walk)
	if err != nil {
		log.Fatal("FilePath error ", err)
	}

	fmt.Printf("\r Scan completed. Indexed %d file(s). %d duplicate(s) found\n", len(files), len(doubleFiles))



	os.Exit(0)
}

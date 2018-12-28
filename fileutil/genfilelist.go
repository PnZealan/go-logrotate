package fileutil

import (
	"io/ioutil"
	"log"
	"os"
	"time"
)

// ReadDir read directory return os.FileInfo
func ReadDir(filename string) (files []os.FileInfo) {
	files, err := ioutil.ReadDir(filename)
	if err != nil {
		log.Println("open directory err")
		return nil
	}
	return
}

//GenFileList function, generate filelist through determine
//input <==  dir
func GenFileList(files []os.FileInfo, excludeList []string, lastModf float64, size int64) (fileList []string) {
	var flag bool
	size = size * 1024 * 1024
	t := time.Now()
	for _, file := range files {
		flag = true
		for _, s := range excludeList {
			if file.Name() == s || file.IsDir() {
				flag = false
				break
			} else {
				// calculate the difference between the current time and the last modified time of the file
				s := t.Sub(file.ModTime()).Hours()
				if s < lastModf || file.Size() < size {
					flag = false
					break
				}
			}
		}
		if flag {
			fileList = append(fileList, file.Name())
		}
	}
	return
}

// Filefilter examine file
// input <== file
func Filefilter(file string, excludeList []string, lastModf float64, size int64) (flag bool) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		log.Println("file not found")
		return
	}
	t := time.Now()
	size = size * 1024 * 1024
	flag = true
	for _, s := range excludeList {
		if fileInfo.Name() == s || fileInfo.IsDir() {
			flag = false
			return
		}
		// calculate the difference between the current time and the last modified time of the file
		s := t.Sub(fileInfo.ModTime()).Hours()
		if s < lastModf || fileInfo.Size() < size {
			flag = false
			return
		}

	}
	return
}

//FileStat determine the type of file
func FileStat(filename string) (bool, error) {
	f, err := os.Stat(filename)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return f.IsDir(), nil
}

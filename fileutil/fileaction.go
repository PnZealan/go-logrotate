package fileutil

import (
	"bufio"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// fileName is abs path !
func fileCompress(fileName string) error {
	cmpName := fileName + ".gz"
	// create compressed file
	outputFile, err := os.Create(cmpName)
	if err != nil {
		return err
	}

	//gzip writer
	gzipWriter := gzip.NewWriter(outputFile)
	defer gzipWriter.Close()

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// buffer io reader
	bfRD := bufio.NewReader(file)

	for {
		data, _, err := bfRD.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		// write bytes slice
		_, errWrite := gzipWriter.Write(data)
		if errWrite != nil {
			return errWrite
		}
	}

	log.Println("Compressed file :", cmpName)
	transfer(fileName)
	return nil
}

// FileRotate ,fileName, dest, pidPath are abs path !, when dest is null, truncate file
func FileRotate(fileName string, dest string, pidPath string, compress bool) error {
	t := time.Now().Format("2006-01-02")

	if dest != "" {
		newName := path.Join(dest, path.Base(fileName)) + "-" + string(t)
		err := os.Rename(fileName, newName)
		if err != nil {
			log.Println("file rename err:-----------<", err)
			return err
		}
		if err := pid2Reload(pidPath); err != nil {
			log.Println("pid reload err:-----------<", err)
			// recover the file name
			os.Rename(newName, fileName)
			return err
		}
		// newName abs path
		if compress == true {
			if err := fileCompress(newName); err != nil {
				log.Println(err)
				return err
			}
		}
	} else {
		if err := fileTruncate(fileName); err != nil {
			log.Println(err)
			return err
		}
	}

	// } else {
	// 	newName := fileName + "-" + string(t)
	// 	err := os.Rename(fileName, newName)
	// 	if err == nil {
	// 		pid2Reload(pidPath)
	// 		deleteFile(newName)
	// 	}
	log.Printf("fileLogRotate-->%v, dest:%v, compress:%t \n", fileName, dest, compress)
	return nil
}

// FileCopyTruncate is abs path !
func FileCopyTruncate(fileName string, dest string, compress bool) error {
	t := time.Now().Format("2006-01-02")

	if dest != "" {
		newName := path.Join(dest, path.Base(fileName)) + "-" + string(t)
		if err := fileCopy(fileName, newName); err != nil {
			log.Println(err)
			return err
		}
		if compress == true {
			if err := fileCompress(newName); err != nil {
				log.Println(err)
				return err
			}
		}
	}
	if err := fileTruncate(fileName); err != nil {
		log.Println(err)
		return err
	}
	log.Printf("fileCopyTruncate-->%v, dest:%v, compress:%t \n", fileName, dest, compress)

	return nil
}

func pid2Reload(pidPath string) error {
	pfile, err := os.Open(pidPath)
	defer pfile.Close()
	if err != nil {
		return err
	}
	pidData, _ := ioutil.ReadAll(pfile)
	pid := string(pidData)
	pid = strings.Replace(pid, "\n", "", -1)

	cmd := exec.Command("kill", "-USR1", pid)
	_, errCmd := cmd.Output()
	if errCmd != nil {
		log.Println("reload cmd exec failed：" + errCmd.Error())
		return errCmd
	}
	return nil
}

func deleteFile(filename string) error {
	if err := os.Remove(filename); err != nil {
		return err
	}
	return nil
}

func fileTruncate(filename string) error {
	// file, err := os.Open(filename)
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()
	// if err := file.Truncate(0); err != nil {
	// 	return err
	// }
	// return nil
	err := os.Truncate(filename, 0)
	if err != nil {
		log.Println(filename, "truncate err")
	}
	return err
}

func fileCopy(srcName string, destName string) error {
	src, err := os.Open(srcName)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, errD := os.Create(destName)
	if errD != nil {
		return errD
	}
	defer dest.Close()

	_, errC := io.Copy(dest, src)
	if errC != nil {
		return errC
	}
	return nil
}

//parmars: ENDPOINT,AK,AKSECRET,BKNAME,OBNAME
func transfer(fileName string) error {
	endpoint := os.Getenv("ENDPOINT")
	ak := os.Getenv("AK")
	aksecret := os.Getenv("AKSECRET")
	bkname := os.Getenv("BKNAME")
	obname := os.Getenv("OBNAME")

	//get oss client instance
	client, err := oss.New(endpoint, ak, aksecret)
	if err != nil {
		log.Println("get oss client instances error ")
		return err
	}

	//get oss bk 
	bucket, err := client.Bucket(bkname)
	if err != nil {
		log.Println("get oss bk error")
		return err
	}

	//set partSize 1024* 1024, 3 goroutines for upload, enable check back 
	err = bucket.UploadFile(obname, fileName, 100*1024, oss.Routines(3), oss.Checkpoint(true, ""))
	
	return nil
}
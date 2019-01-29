package main

import (
	"log"
	"path"
	"github.com/robfig/cron"

	"./fileutil"
	"./jsonutil"
)

type jsConfig struct {
	FilePathList map[string]fileConfig `json:"filePathList"`
	ExcludeFile  []string              `json:"excludeFile"`
}

type fileConfig struct {
	Size         int64   `json:"size"`
	Copytruncate bool    `json:"copytruncate"`
	PidPath      string  `json:"pidPath"`
	Compress     bool    `json:"compress"`
	Dest         string  `json:"dest"`
	LastModf     float64 `json:"lastModf"`
	SpecTime     string  `json:"specTime"`
}

type startJob struct {
	jc *jsConfig
	k  string
	v  fileConfig
}

func (j startJob) Run() {
	start(j.jc, j.k, j.v)
}

func main() {
	js := jsonutil.JSONStruct{}
	jc := &jsConfig{}
	js.Load("config.json", jc)
	c := cron.New()
	// addjob to each filepath
	for k, v := range jc.FilePathList {
		sj := startJob{jc, k, v}
		c.AddJob(v.SpecTime, sj)
	}
	c.Start()

	defer c.Stop()
	//block main thread
	select {}
}

func start(jc *jsConfig, k string, v fileConfig) {
	excludeList := jc.ExcludeFile

	b, err := fileutil.FileStat(k)
	// if k is dir
	if b == true {
		// get all the authorized files in the one directory (k)
		fileList := fileutil.GenFileList(fileutil.ReadDir(k), excludeList, v.LastModf, v.Size)
		if v.Copytruncate == true {
			dirStartCopyTruncate(k, fileList, v.Dest, v.Compress)
		} else {
			dirStartRotate(k, fileList, v.Dest, v.PidPath, v.Compress)
		}
	} else {
		if err != nil {
			log.Println("attach file stat err", err)
		}
		// if k is file
		if flag := fileutil.Filefilter(k, excludeList, v.LastModf, v.Size); flag == true {
			if v.Copytruncate == true {
				fileStartCopyTruncate(k, v.Dest, v.Compress)
			} else {
				fileStartRotate(k, v.Dest, v.PidPath, v.Compress)
			}
		}

	}

}

func dirStartCopyTruncate(base string, fileList []string, dest string, compress bool) {
	for _, file := range fileList {
		if err := fileutil.FileCopyTruncate(path.Join(base, file), dest, compress); err != nil {
		}

	}

}

func dirStartRotate(base string, fileList []string, dest string, pidPath string, compress bool) {
	for _, file := range fileList {
		if err := fileutil.FileRotate(path.Join(base, file), dest, pidPath, compress); err != nil {
		}
	}

}

func fileStartCopyTruncate(filename string, dest string, compress bool) {
	if err := fileutil.FileCopyTruncate(filename, dest, compress); err != nil {
		log.Println(err)
	}
}

func fileStartRotate(filename string, dest string, pidPath string, compress bool) {
	if err := fileutil.FileRotate(filename, dest, pidPath, compress); err != nil {

	}
}

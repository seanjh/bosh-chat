package message

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/seanjh/bosh-chat/util"
)

const storagePath = "/tmp/messages"
const storageExt = "msg"

type fileStorage struct{}

type writeJob struct {
	content []byte
	index   chan<- int
	err     chan<- error
}

// TODO provide 1 of these per "channel"
var writeQueue = make(chan writeJob)

func filenameIndex(filename string) int {
	_, file := filepath.Split(filename)
	ext := strings.Index(file, "."+storageExt)

	if ext == -1 {
		return errIndex
	}

	index, err := strconv.ParseInt(file[0:ext], 10, 32)
	if err != nil {
		return errIndex
	}

	return int(index)
}

func messageFilename(index int) string {
	return string(index) + storageExt
}

func (s *fileStorage) write(content []byte) (chan int, chan error) {
	ci := make(chan int, 1)
	ce := make(chan error, 1)

	writeQueue <- writeJob{content, ci, ce}

	return ci, ce
}

func largestIndex() int {
	util.EnsurePath(storagePath)
	files, err := ioutil.ReadDir(storagePath)
	if err != nil {
		log.Println("Failed to read directory", err)
	}

	if len(files) == 0 {
		log.Println("No existing files at", storagePath)
		return errIndex
	}

	filename := files[len(files)-1].Name()
	return filenameIndex(filename)
}

func queueWriter(c <-chan writeJob) {
	util.EnsurePath(storagePath)
	lastIndex := largestIndex()
	go func(index int) {
		for job := range c {
			log.Println("Received write job for index", index)
			filename := messageFilename(index)
			log.Println("Writing", filename)
			err := ioutil.WriteFile(filename, job.content, 0400)
			if err != nil {
				log.Println("Error writing output", job.content)
				job.err <- err
			} else {
				log.Println("Finished writing", filename)
				job.index <- index
				index++
			}
			close(job.index)
			close(job.err)
		}
	}(lastIndex + 1)
}

func (s *fileStorage) wait(last, limit int) <-chan *message {
	c := make(chan *message, limit)

	go func(timeout time.Duration) {
		defer close(c)
		<-time.Tick(timeout)
		log.Println("timeout waiting for messages")
	}(10 * time.Second)

	return c
}

package message

import (
	"fmt"
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
	content string
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
	name := fmt.Sprintf("%d.%s", index, storageExt)
	return filepath.Join(storagePath, name)
}

func (s *fileStorage) write(content string) (chan int, chan error) {
	log.Println("writing contents to file", content)
	ci := make(chan int, 1)
	ce := make(chan error, 1)

	log.Printf("sending contents to write queue: '%s'\n", content)
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

// StartWriter TODO
func StartWriter() {
	c := writeQueue
	util.EnsurePath(storagePath)
	lastIndex := largestIndex()
	log.Printf("Starting queue writer at index %d\n", lastIndex)
	go func(index int) {
		for job := range c {
			log.Println("Received write job for index", index)
			filename := messageFilename(index)
			log.Printf("Writing file %s: '%s'\n", filename, job.content)
			err := ioutil.WriteFile(filename, []byte(job.content), 0400)
			if err != nil {
				log.Println("Error writing file", filename)
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
	log.Printf("Waiting from index %d for up to %d messages\n", last, limit)
	c := make(chan *message, limit)

	go func(timeout time.Duration) {
		defer close(c)
		log.Println("Timed out waiting for messages")
		<-time.Tick(timeout)
	}(60 * time.Second)

	return c
}

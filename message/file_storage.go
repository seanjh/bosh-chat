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
const timeoutSecs = 10

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
	ci := make(chan int, 1)
	ce := make(chan error, 1)

	log.Printf("sending contents to write queue")
	writeQueue <- writeJob{content, ci, ce}
	return ci, ce
}

func largestIndex() int {
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

func (s *fileStorage) read(index int) (*message, error) {
	filename := messageFilename(index)
	log.Printf("Reading file at index %d: '%s'\n", index, filename)

	m := message{Index: index, Body: ""}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return &m, err
	}

	m.Body = string(content)
	return &m, nil
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
			err := ioutil.WriteFile(filename, []byte(job.content), 0600)
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

// yield up to limit new indices greater than last
func newIndices(last, limit int, c, done chan<- int) {
	files, err := ioutil.ReadDir(storagePath)
	if err != nil {
		log.Println("Failed to read directory", err)
	}

	for _, filename := range files {
		i := filenameIndex(filename.Name())
		if i > last {
			log.Printf("Found new index '%d' > '%d'\n", i, last)
			c <- i
			limit--
		}
		if limit == 0 {
			log.Println("Exhausted message limit")
			break
		}
	}

	done <- 1
	close(c)
}

// read the directory every Xms waiting for files
func scan(last, limit int, ci chan<- int, timeout <-chan time.Time) {
	c, done := make(chan int), make(chan int, 1)

loop:
	for {
		log.Println("Scanning from index", last)
		go newIndices(last, limit, c, done)

		select {
		case i, ok := <-c:
			log.Println("Scan received result")
			if ok {
				log.Println("Scan received new index", i)
				ci <- i
				for i = range c {
					log.Println("Scan received new index", i)
					ci <- i
				}
			}
			break loop

		case <-done:
			log.Println("No scan results. Time to snooze...")
			<-time.Tick(1 * time.Second)
			c = make(chan int)

		case <-timeout:
			log.Println("Timed out waiting for messages.")
			break loop

		}

	}

	log.Println("Scan complete.")
	close(ci)
}

func (s *fileStorage) wait(last, limit int) <-chan *message {
	log.Printf("Waiting from index %d for up to %d messages\n", last, limit)

	cm := make(chan *message)
	go func() {
		c := make(chan int)
		timeout := time.Tick(timeoutSecs * time.Second)
		go scan(last, limit, c, timeout)
		for i := range c {
			m, err := s.read(i)
			if err != nil {
				log.Printf("Error reading file: '%s'\n", messageFilename(i))
				continue
			} else {
				log.Printf("Wait delivering message %s: '%s'\n", m, m.Body)
				cm <- m
			}
		}

		log.Println("Wait finished delivering messages.")
		close(cm)
	}()

	return cm
}

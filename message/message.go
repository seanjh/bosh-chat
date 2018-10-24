package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/seanjh/bosh-chat/session"
)

const errIndex = -1
const maxMessages = 10

type message struct {
	index int    `json:"index"`
	body  string `json:"body"`
}

type reader interface {
	list(last int) ([]int, error)
	read(index int) (*message, error)
	wait(last, limit int) <-chan *message
}

type writer interface {
	write(content []byte) (<-chan int, <-chan error)
}

type readWriter interface {
	reader
	writer
}

func (m *message) String() string {
	return "message@" + string(m.index)
}

type messageHandler = func(index int, body io.ReadCloser) (string, int, error)

func messageSession(handler messageHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := session.LoadCookie(r)

		//log.Println("Request Method:", r.Method)
		//log.Println("Request Header:", r.Header)
		//log.Println("Request Host:", r.Host)
		//log.Println("Request ContentLength:", r.ContentLength)
		//log.Println("Request Form:", r.Form)
		//log.Println("Request PostForm:", r.PostForm)
		//log.Println("Request MultipartForm:", r.MultipartForm)
		//log.Println("Request RemoteAddr:", r.RemoteAddr)
		//log.Println("Request TransferEncoding:", r.TransferEncoding)

		content, index, err := handler(sess.Pos, r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to process message request")
			return
		}

		sess.Pos = index
		err = sess.Save()
		if err != nil {
			log.Println("Failed to save session", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to save session")
			return
		}
		http.SetCookie(w, sess.Cookie())

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, content)
	}
}

// HandleMessages waits for and returns at between 1 and 10 total
// messages for GET requests, and writes new messages for POST requests.
func HandleMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		log.Println("Handling message GET")
		messageSession(messageWait)(w, r)
	} else if r.Method == http.MethodPost {
		log.Println("Handling message POST")
		messageSession(messageAppend)(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func messageAppend(_ int, body io.ReadCloser) (string, int, error) {
	log.Println("Appending new message")
	s := fileStorage{}

	decoder := json.NewDecoder(body)
	m := &message{}
	err := decoder.Decode(&m)
	if err != nil {
		return "", errIndex, err
	}

	ci, ce := s.write(m.body)
	select {
	case index := <-ci:
		log.Println("Finished writing")
		return "", index, nil
	case err = <-ce:
		log.Println("Failed to write message")
		return "", errIndex, err
	}
}

func messageWait(index int, _ io.ReadCloser) (string, int, error) {
	log.Println("Waiting for new messages")
	s := fileStorage{}

	tmp, msgs := make([]message, maxMessages), make([]message, 0)
	i, last := 0, index
	for m := range s.wait(index, maxMessages) {
		log.Println("Received new message from wait", m)
		last = m.index
		tmp[i] = *m
		msgs = tmp[:i]
		i++
	}

	log.Printf("Received %d messages. New last index: %d\n", len(msgs), last)
	b, err := json.Marshal(msgs)
	if err != nil {
		return "", errIndex, errors.New("Failed to marshal messages")
	}

	return string(b), last, nil
}

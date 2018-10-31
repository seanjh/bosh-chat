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
	Index int    `json:"index"`
	Body  string `json:"body"`
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
	return fmt.Sprintf("message@%d", m.Index)
}

type messageHandler = func(index int, body io.ReadCloser) (string, int, error)

func messageSession(handler messageHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := session.LoadCookie(r)

		content, index, err := handler(sess.Pos, r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to process message request")
			return
		}

		if sess.Pos != index {
			log.Printf("Saving updated session index %d->%d\n", sess.Pos, index)
			sess.Pos = index
			err = sess.Save()
		}
		if err != nil {
			log.Println("Failed to save session", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to save session")
			return
		}
		http.SetCookie(w, sess.Cookie())

		_, err = fmt.Fprintf(w, content)
		if err != nil {
			log.Println("Error writing response", err)
		}
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
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func messageAppend(_ int, body io.ReadCloser) (string, int, error) {
	log.Println("Appending new message")
	s := fileStorage{}

	decoder := json.NewDecoder(body)
	m := &message{}
	err := decoder.Decode(m)
	if err != nil {
		return "", errIndex, err
	}

	ci, ce := s.write(m.Body)
	select {
	case index := <-ci:
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
		last = m.Index
		log.Println("New last index:", last)
		tmp[i] = *m
		msgs = tmp[:i+1]
		log.Println("Total messages:", len(msgs))
		i++
	}

	log.Printf("Received %d messages. New last index: %d\n", len(msgs), last)
	b, err := json.Marshal(msgs)
	if err != nil {
		return "", errIndex, errors.New("Failed to marshal messages")
	}

	return string(b), last, nil
}

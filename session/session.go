package session

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/seanjh/bosh-chat/util"
)

// Tail represents the end position in a collection of messages
const Tail = -1

// Cookie TODO
const cookie = "sessionid"

const sessionsPath = "/tmp/sessions"

// Session tracks the position in messages for a client
type Session struct {
	ID  string `json:"id"`
	Pos int    `json:"pos"`
}

func (s *Session) String() string {
	return fmt.Sprintf("session-%s@%d", s.ID, s.Pos)
}

// newID creates a 128 bit random hex-encoded ID number
func newID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate new ID: %s\n", err)
		return ""
	}
	return hex.EncodeToString(b)
}

// NewSession returns a Session with new random ID
func NewSession() *Session {
	log.Println("Creating new Session")
	return &Session{
		ID:  newID(),
		Pos: Tail,
	}
}

func getFilename(id string) string {
	return filepath.Join(sessionsPath, id)
}

// read sets the index to the last known position stored on disk
func read(filename string) (*Session, error) {
	util.EnsurePath(sessionsPath)

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Failed to read session:", err)
		return &Session{}, err
	}

	var s Session
	err = json.Unmarshal(b, &s)
	if err != nil {
		log.Println("Failed to unmarshal session:", err)
		return &Session{}, err
	}

	return &s, nil
}

// Load TODO
func LoadCookie(r *http.Request) *Session {
	cookie, err := r.Cookie(cookie)
	if err != nil {
		return NewSession()
	}

	s, err := read(getFilename(cookie.Value))
	if err != nil {
		s = NewSession()
		err = s.Save()
		log.Println("Creating new session")
		if err != nil {
			log.Println("Failed to create new session")
		}
	}

	log.Println("Session loaded:", s)
	return s
}

// save TODO
func (m *Session) Save() error {
	util.EnsurePath(sessionsPath)

	b, err := json.Marshal(m)
	if err != nil {
		log.Println("Failed to marshal session:", err)
		return err
	}
	err = ioutil.WriteFile(getFilename(m.ID), b, 0600)
	if err != nil {
		log.Println("Failed to save session:", err)
		return err
	}
	return nil
}

// Cookie TODO
func (m *Session) Cookie() *http.Cookie {
	return &http.Cookie{Name: cookie, Value: m.ID}
}

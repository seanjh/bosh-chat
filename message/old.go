package message

// Messages are a collection of Message
//type messages []message

//func (m messages) Len() int           { return len(m) }
//func (m messages) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
//func (m messages) Less(i, j int) bool { return m[i].index < m[j].index }

//func (m messages) unread(lastIndex int) messages {
//if len(m)-1 > lastIndex {
//sort.Sort(m)
//return m[lastIndex+1:]
//}
//return make(Messages, 0)
//}

//func (s *fileStorage) list() ([]Message, error) {
////func listFiles() ([]Message, error) {
//util.EnsurePath(messagePath)
//files, err := ioutil.ReadDir(messagePath)
//if err != nil {
//return []Message{}, err
//}

//messages := make([]Message, len(files))
//for i, f := range files {
//index := index(f.Name())
//if index == errIndex {
//log.Println("Error parsing message filename:", err)
//continue
//}
//messages[i] = Message{Index: int(index)}
//log.Println("Found message", messages[i])
//}

//log.Printf("Finished scanning %d total files from %s\n", len(messages), messagePath)
//return messages, nil
//}

//func (s *fileStorage) read(m *Message) error {
//filename := s.filename(m.Index)
//fmt.Println("Reading file:", filename)
//raw, err := ioutil.ReadFile(filename)
//if err != nil {
//return err
//}

//m.Body = string(raw)
//return nil
//}

//func (s *fileStorage) save(queue chan<- writeJob) error {
//c := make(chan int)
//queue <- writeJob{s.Body, c}
//index := <-c
//if index == errIndex {
////TODO handle this
//}
//s.Index = index
//return nil
//}

//func (s *fileStorage) write() error {
////TODO handle uninitialized Index
//content := []byte(s.Body)
//return ioutil.WriteFile(s.filename(), content, 0400)
//}

//func readDir(c chan<- string) {
//util.EnsurePath(messagePath)
//defer close(c)
//files, _ := ioutil.ReadDir(messagePath)
//for _, f := range files {
//filename := filepath.Join(messagePath, f.Name())
//content, err := read(filename)
//if err != nil {
//break
//}
//c <- content
//}
//}

// ReadMessages TODO
//func ReadMessages() <-chan string {
//contents := make(chan string)
//go readDir(contents)
//return contents
//}

//func newMessages(messages []Message, last, limit int) {
//log.Printf("Checking message %v against last %d\n", m, last)
//if m.Index > last {
//c <- m
//limit--
//}
//if limit == 0 {
//break
//}
//}

//func readMessages(c chan<- Message, last, limit int) {
//defer close(c)

//m, err := listFiles()
//if err != nil {
//log.Println("Failed to read messages", err)
//return
//}

//if len(m) == 0 {
////hold up
//} else {
////immediately deliver the next messages
//}

//for _, m := range m {
//log.Printf("Checking message %v against last %d\n", m, last)
//if m.Index > last {
//c <- m
//limit--
//}
//if limit == 0 {
//break
//}
//}
//log.Println("Finished reading messages")

//}

// Wait TODO
//func Wait(last int) ([]Message, int) {
//log.Println("Waiting for new messages after index: ", last)
//messages := make([]Message, maxMessages)

////wait for any message - return up to 10

//c := make(chan Message, maxMessages)
//go readMessages(c, last, maxMessages)

//i := -1
//next := last
//for m := range c {
//log.Println("Wait received new message: ", m)
//i++
//messages[i] = m
//next = m.Index
//}

//messages = messages[:i]
//log.Printf("Returning messages[:%d]: %v\n", i, messages)
//return messages[:i], next
//}

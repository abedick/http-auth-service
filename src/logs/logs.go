package logs

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

// Logs -- simple universal logger
type Logs struct {
	Directory  string                   `toml:"directory"`
	OutWriter  *bufio.Writer            `toml:"-"`
	ErrWriter  *bufio.Writer            `toml:"-"`
	AccWriter  *bufio.Writer            `toml:"-"`
	Writers    map[string]*bufio.Writer `toml:"-"`
	WriteMutex sync.Mutex               `toml:"-"`
}

type message struct {
	Time   string `json:"time"`
	Status string `json:"status,omitempty"`
	Type   string `json:"type"`
	URL    string `json:"url,omitempty"`
	IP     string `json:"ip,omitempty"`
	Msg    string `json:"msg"`
}

// Init - Initialize log files
func (l *Logs) Init() error {

	fileExtensions := []string{"out", "err", "acc"}
	l.Writers = make(map[string]*bufio.Writer)
	time := time.Now().Format("2006-01-02")

	for _, ext := range fileExtensions {

		err := createDir(l.Directory + "/" + time)
		if err != nil {
			return errors.New("could not create one or more directories for log files: " + err.Error())
		}

		path := l.Directory + "/" + time + "/log." + ext
		file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			file, err = os.Create(path)
			if err != nil {
				return errors.New("Could not create one or more log files: " + err.Error())
			}
		}
		l.Writers[ext] = bufio.NewWriter(file)
	}

	return nil
}

// UniversalLogger - use as middleware logger
func (l *Logs) UniversalLogger(i, u string) {
	LogMsg := message{
		Type: "middleware",
		IP:   i,
		URL:  u,
	}
	l.out(LogMsg)
}

// Message - puts a message in the log
func (l *Logs) Message(status, msg string) {
	LogMsg := message{
		Type:   "info",
		Status: status,
		Msg:    msg,
	}
	l.out(LogMsg)
}

func (l *Logs) out(Message message) {

	Message.Time = time.Now().Format("2006-01-02 15:04:05")
	MessageBytes, _ := json.Marshal(Message)

	l.WriteMutex.Lock()

	if Message.Type == "middleware" {
		l.Writers["acc"].Write(MessageBytes)
		l.Writers["acc"].WriteString("\n")
		l.Writers["acc"].Flush()
	} else if Message.Type == "info" {
		if Message.Status == "error" {
			l.Writers["err"].Write(MessageBytes)
			l.Writers["err"].WriteString("\n")
			l.Writers["err"].Flush()
		}
		l.Writers["out"].Write(MessageBytes)
		l.Writers["out"].WriteString("\n")
		l.Writers["out"].Flush()
	}

	l.WriteMutex.Unlock()
}

func createDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

package logger

import (
	"io/ioutil"
	l "log"
	"net/http"
	"os"
)

var log *l.Logger

type fileLog string

func (f1 fileLog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(string(f1), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(data)
}

func Run(dest string) {
	log = l.New(fileLog(dest), "", l.LstdFlags)
}

func RegisterHandler() {
	http.HandleFunc("/log", func(rw http.ResponseWriter, r *http.Request) {
		msg, err := ioutil.ReadAll(r.Body)
		if err != nil || len(msg) == 0 {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		write(string(msg))
	})
}

func write(msg string) {
	log.Printf("%v/n", msg)
}

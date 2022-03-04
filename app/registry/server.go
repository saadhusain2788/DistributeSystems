package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const (
	ServerPort = ":3000"
	ServicURL  = "http://localhost" + ServerPort + "/services"
)

type registry struct {
	registration []Registration
	mutex        *sync.Mutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registration = append(r.registration, reg)
	r.mutex.Unlock()
	return nil
}

func (r *registry) remove(serviceURL string) error {
	for i := range r.registration {
		if r.registration[i].ServiceURL == serviceURL {
			r.mutex.Lock()
			r.registration = append(r.registration[:i], r.registration[i+1:]...)
			r.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("Service at url %v not found", serviceURL)
}

var reg = registry{
	registration: make([]Registration, 0),
	mutex:        new(sync.Mutex),
}

type RegistryService struct{}

func (service RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")
	switch r.Method {
	case http.MethodPost:
		var registration Registration
		err := json.NewDecoder(r.Body).Decode(&registration)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service %v with url %v", registration.ServiceName, registration.ServiceURL)
		err = reg.add(registration)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(payload)
		log.Printf("Removing service %v with url %v", url)
		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}

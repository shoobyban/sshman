package backend

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type jsonFile struct {
	Key   string           `json:"key"`
	Hosts map[string]*Host `json:"hosts"`
	Users map[string]*User `json:"users"`
}

type JsonStorage struct {
	Path string
	data *Data
}

func (m *JsonStorage) Load(data *Data) error {
	b, err := os.ReadFile(m.Path)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}
	var cf jsonFile
	err = json.Unmarshal(b, &cf)
	if err != nil {
		log.Printf("Error: unable to decode into struct, please correct or remove broken %s %v\n", m.Path, err)
		return err
	}
	data.key = cf.Key
	data.hosts = cf.Hosts
	data.users = cf.Users
	for alias, host := range data.hosts {
		host.Alias, host.Config = alias, data
		data.hosts[alias] = host
	}
	data.updateGroups()
	return nil
}

func (m *JsonStorage) Save(data *Data) error {
	cf := jsonFile{Key: data.key, Hosts: data.hosts, Users: data.users}
	b, _ := json.MarshalIndent(cf, "", "  ")
	os.WriteFile(m.Path, b, 0600)
	return nil
}

func (m *JsonStorage) Watch(notify func()) {
	initWG := sync.WaitGroup{}
	initWG.Add(1)
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()
		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						eventsWG.Done()
						return
					}
					const writeOrCreateMask = fsnotify.Write | fsnotify.Create
					if filepath.Clean(event.Name) == m.Path &&
						event.Op&writeOrCreateMask != 0 {
						notify()
					} else if filepath.Clean(event.Name) == m.Path &&
						event.Op&fsnotify.Remove != 0 {
						eventsWG.Done()
						return
					}

				case err, ok := <-watcher.Errors:
					if ok {
						log.Printf("watcher error: %v\n", err)
					}
					eventsWG.Done()
					return
				}
			}
		}()
		configDir := filepath.Dir(m.Path)
		watcher.Add(configDir)
		initWG.Done()
		eventsWG.Wait()
	}()
	initWG.Wait()
}

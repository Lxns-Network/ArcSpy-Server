package main

import (
	_ "embed"
	"github.com/fsnotify/fsnotify"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	log "github.com/sirupsen/logrus"
	"os"
)

func initConfig() {
	config.AddDriver(yamlv3.Driver)
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		log.Panic("Config file does not exist")
	}
	if err := config.LoadFiles("config.yaml"); err != nil {
		log.Panic("Failed to load config file: ", err)
	}
}

func watchConfig() (err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					log.Info("Config file changed, reloading...")
					err := config.LoadFiles("config.yaml")
					if err != nil {
						log.Panic("Failed to load config file: ", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error("Failed to listen config file: ", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add("config.yaml")
	if err != nil {
		return err
	}

	// Block main goroutine forever.
	<-make(chan struct{})

	return nil
}

package gsos

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-fsnotify/fsnotify"
)

// Event .
type Event struct {
	Name string // Relative path to the file or directory.
	Op   Op     // File operation that triggered the event.
}

// Op .
type Op uint32

// Ops
const (
	Create Op = 1 << iota
	Write
	Remove
	Rename
	Chmod
)

// FSWatcher filesytem watcher
type FSWatcher struct {
	sync.Mutex                   // mixin mutex
	watcher    *fsnotify.Watcher // watcher implement
	watchDirs  map[string]bool   //watch directories
	Events     chan Event        // Event channel
	Errors     chan error        // Error event channel
}

// NewWatcher establishes a new watcher with the underlying OS and begins waiting for events.
func NewWatcher() (*FSWatcher, error) {

	fsnotifyWatcher, err := fsnotify.NewWatcher()

	if err != nil {
		return nil, err
	}

	watcher := &FSWatcher{
		watcher:   fsnotifyWatcher,
		watchDirs: make(map[string]bool),
		Events:    make(chan Event, 10),
		Errors:    make(chan error, 10),
	}

	go watcher.dispatch()

	return watcher, nil
}

// Add .
func (watcher *FSWatcher) Add(path string, recursively bool) error {

	absPath, err := filepath.Abs(path)

	if err != nil {
		return nil
	}

	if !IsDir(path) {
		watcher.watcher.Add(absPath)
		return nil
	}

	watcher.Lock()
	defer watcher.Unlock()

	return watcher.addDir(absPath, recursively)
}

func (watcher *FSWatcher) addDir(path string, recursively bool) error {

	fmt.Printf("add watch dir :%s\n", path)

	watcher.watchDirs[path] = recursively

	err := watcher.watcher.Add(path)

	if err != nil || !recursively {
		return err
	}

	entries, err := ioutil.ReadDir(path)

	for _, entry := range entries {
		if entry.IsDir() {
			watcher.addDir(filepath.Join(path, entry.Name()), recursively)
		}
	}

	return err
}

// Remove .
func (watcher *FSWatcher) Remove(path string) error {
	absPath, err := filepath.Abs(path)

	if err != nil {
		return nil
	}

	if !IsDir(path) {
		watcher.watcher.Remove(absPath)
		return nil
	}

	watcher.Lock()
	defer watcher.Unlock()

	return watcher.remove(path)
}

func (watcher *FSWatcher) remove(path string) error {

	if flag, ok := watcher.watchDirs[path]; ok {
		delete(watcher.watchDirs, path)

		err := watcher.watcher.Remove(path)

		if err != nil {
			return err
		}

		if flag {
			entries, _ := ioutil.ReadDir(path)
			for _, entry := range entries {
				if entry.IsDir() {
					watcher.remove(filepath.Join(path, entry.Name()))
				}
			}
		}
	}

	return nil
}

func (watcher *FSWatcher) dispatch() {
	for {
		select {
		case event := <-watcher.watcher.Events:

			if (event.Op & fsnotify.Create) == fsnotify.Create {
				if IsDir(event.Name) {
					watcher.onCreateDir(event)
				}
			}
			if (event.Op&fsnotify.Remove) == fsnotify.Remove ||
				(event.Op&fsnotify.Rename) == fsnotify.Rename {
				watcher.onDelDir(event)
			}

			watcher.Events <- Event{event.Name, Op(event.Op)}

		case err := <-watcher.watcher.Errors:
			watcher.Errors <- err
		}
	}
}

func (watcher *FSWatcher) onCreateDir(event fsnotify.Event) {

	watcher.Lock()
	defer watcher.Unlock()

	for name, flag := range watcher.watchDirs {

		if strings.HasPrefix(event.Name, name) && flag {
			watcher.watcher.Add(event.Name)
			watcher.watchDirs[event.Name] = true
		}
	}
}

func (watcher *FSWatcher) onDelDir(event fsnotify.Event) {

	watcher.Lock()
	defer watcher.Unlock()

	watcher.watcher.Remove(event.Name)
	delete(watcher.watchDirs, event.Name)

	for name, flag := range watcher.watchDirs {
		if strings.HasPrefix(name, event.Name) && flag {
			watcher.watcher.Remove(name)
			delete(watcher.watchDirs, name)
		}
	}
}

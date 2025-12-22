// Package memory is a memory source
package memory

import (
	"crypto/rand"
	"fmt"
	"github.com/bitxx/load-config/source"
	"sync"
	"time"
)

type memory struct {
	sync.RWMutex
	ChangeSet *source.ChangeSet
	Watchers  map[string]*watcher
}

func (m *memory) Read() (*source.ChangeSet, error) {
	m.RLock()
	cs := &source.ChangeSet{
		Format:    m.ChangeSet.Format,
		Timestamp: m.ChangeSet.Timestamp,
		Data:      m.ChangeSet.Data,
		Checksum:  m.ChangeSet.Checksum,
		Source:    m.ChangeSet.Source,
	}
	m.RUnlock()
	return cs, nil
}

func (m *memory) generateWatcherID() string {
	// 使用方案2（推荐）或方案1
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func (m *memory) Watch() (source.Watcher, error) {
	w := &watcher{
		Id:      m.generateWatcherID(),
		Updates: make(chan *source.ChangeSet, 100),
		Source:  m,
	}

	m.Lock()
	m.Watchers[w.Id] = w
	m.Unlock()
	return w, nil
}

func (m *memory) Write(cs *source.ChangeSet) error {
	m.Update(cs)
	return nil
}

// Update allows manual updates of the config data.
func (m *memory) Update(c *source.ChangeSet) {
	// don't process nil
	if c == nil {
		return
	}

	// hash the file
	m.Lock()
	// update changeset
	m.ChangeSet = &source.ChangeSet{
		Data:      c.Data,
		Format:    c.Format,
		Source:    "memory",
		Timestamp: time.Now(),
	}
	m.ChangeSet.Checksum = m.ChangeSet.Sum()

	// update watchers
	for _, w := range m.Watchers {
		select {
		case w.Updates <- m.ChangeSet:
		default:
		}
	}
	m.Unlock()
}

func (m *memory) String() string {
	return "memory"
}

func NewSource(opts ...source.Option) source.Source {
	var options source.Options
	for _, o := range opts {
		o(&options)
	}

	s := &memory{
		Watchers: make(map[string]*watcher),
	}

	if options.Context != nil {
		c, ok := options.Context.Value(changeSetKey{}).(*source.ChangeSet)
		if ok {
			s.Update(c)
		}
	}

	return s
}

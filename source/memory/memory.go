// Package memory is a memory source
package memory

import (
	"github.com/bitxx/load-config/source"
	"sync"
	"time"

	"github.com/google/uuid"
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

func (m *memory) Watch() (source.Watcher, error) {
	w := &watcher{
		Id:      uuid.New().String(),
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

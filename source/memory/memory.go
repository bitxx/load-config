// Package memory is a memory source
package memory

import (
	"github.com/bitxx/load-config/source"
	"sync"
	"time"
)

type memory struct {
	sync.RWMutex
	ChangeSet *source.ChangeSet
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

	s := &memory{}

	if options.Context != nil {
		c, ok := options.Context.Value(changeSetKey{}).(*source.ChangeSet)
		if ok {
			s.Update(c)
		}
	}

	return s
}

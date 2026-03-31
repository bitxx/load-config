package json

import (
	"encoding/json"
	"fmt"
	"github.com/bitxx/load-config/reader"
	"github.com/bitxx/load-config/source"
	"strconv"
	"strings"
	"time"
)

type jsonValues struct {
	ch   *source.ChangeSet
	data interface{}
}

type jsonValue struct {
	val interface{}
}

// ---------- init ----------

func newValues(ch *source.ChangeSet) (reader.Values, error) {
	var data interface{}

	b, _ := reader.ReplaceEnvVars(ch.Data)
	if err := json.Unmarshal(b, &data); err != nil {
		data = string(ch.Data)
	}

	return &jsonValues{ch: ch, data: data}, nil
}

// ---------- path helpers ----------

func get(data interface{}, path ...string) interface{} {
	cur := data

	for _, p := range path {
		switch v := cur.(type) {
		case map[string]interface{}:
			cur = v[p]
		case []interface{}:
			idx, err := strconv.Atoi(p)
			if err != nil || idx < 0 || idx >= len(v) {
				return nil
			}
			cur = v[idx]
		default:
			return nil
		}
	}
	return cur
}

func set(data interface{}, val interface{}, path ...string) interface{} {
	if len(path) == 0 {
		return val
	}

	m, ok := data.(map[string]interface{})
	if !ok {
		m = map[string]interface{}{}
	}

	cur := m
	for i := 0; i < len(path)-1; i++ {
		p := path[i]
		if next, ok := cur[p].(map[string]interface{}); ok {
			cur = next
		} else {
			newMap := map[string]interface{}{}
			cur[p] = newMap
			cur = newMap
		}
	}

	cur[path[len(path)-1]] = val
	return m
}

func del(data interface{}, path ...string) {
	if len(path) == 0 {
		return
	}

	m, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	if len(path) == 1 {
		delete(m, path[0])
		return
	}

	next, ok := m[path[0]]
	if !ok {
		return
	}

	del(next, path[1:]...)
}

// ---------- Values ----------

func (j *jsonValues) Get(path ...string) reader.Value {
	return &jsonValue{val: get(j.data, path...)}
}

func (j *jsonValues) Del(path ...string) {
	if len(path) == 0 {
		j.data = map[string]interface{}{}
		return
	}
	del(j.data, path...)
}

func (j *jsonValues) Set(val interface{}, path ...string) {
	j.data = set(j.data, val, path...)
}

func (j *jsonValues) Bytes() []byte {
	b, _ := json.Marshal(j.data)
	return b
}

func (j *jsonValues) Map() map[string]interface{} {
	if m, ok := j.data.(map[string]interface{}); ok {
		return m
	}
	return nil
}

func (j *jsonValues) Scan(v interface{}) error {
	b, err := json.Marshal(j.data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func (j *jsonValues) String() string {
	return "json"
}

// ---------- Value ----------

func (j *jsonValue) Bool(def bool) bool {
	switch v := j.val.(type) {
	case bool:
		return v
	case string:
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return def
}

func (j *jsonValue) Int(def int) int {
	switch v := j.val.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		i, err := strconv.Atoi(v)
		if err == nil {
			return i
		}
	}
	return def
}

func (j *jsonValue) String(def string) string {
	if s, ok := j.val.(string); ok {
		return s
	}
	return def
}

func (j *jsonValue) Float64(def float64) float64 {
	switch v := j.val.(type) {
	case float64:
		return v
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return f
		}
	}
	return def
}

func (j *jsonValue) Duration(def time.Duration) time.Duration {
	if s, ok := j.val.(string); ok {
		d, err := time.ParseDuration(s)
		if err == nil {
			return d
		}
	}
	return def
}

func (j *jsonValue) StringSlice(def []string) []string {
	switch v := j.val.(type) {
	case string:
		sl := strings.Split(v, ",")
		if len(sl) > 1 {
			return sl
		}
	case []interface{}:
		res := make([]string, 0, len(v))
		for _, item := range v {
			res = append(res, fmt.Sprintf("%v", item))
		}
		return res
	}
	return def
}

func (j *jsonValue) StringMap(def map[string]string) map[string]string {
	m, ok := j.val.(map[string]interface{})
	if !ok {
		return def
	}

	res := make(map[string]string, len(m))
	for k, v := range m {
		res[k] = fmt.Sprintf("%v", v)
	}
	return res
}

func (j *jsonValue) Scan(v interface{}) error {
	b, err := json.Marshal(j.val)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func (j *jsonValue) Bytes() []byte {
	b, err := json.Marshal(j.val)
	if err != nil {
		return []byte{}
	}
	return b
}

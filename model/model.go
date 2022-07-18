package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

var ModelMap sync.Map

type Result struct {
	Took     int64  `json:"took"`
	TimedOut bool   `json:"timed_out"`
	Shards   Shards `json:"_shards"`
	Hits     Hits   `json:"hits"`
}

type Hits struct {
	Total Total          `json:"total"`
	Hits  []*MarshalData `json:"hits"`
}

type Hit struct {
	Index  string      `json:"_index"`
	ID     string      `json:"_id"`
	Source interface{} `json:"_source"`
}
type help struct {
	Index string `json:"_index"`
}

type MarshalData Hit

func (h *MarshalData) UnmarshalJSON(data []byte) (ret error) {
	defer func() {
		err2 := recover()
		if err2 != nil {
			ret = fmt.Errorf("%v", err2)
		}
	}()
	tmp := &Hit{}
	helpVal := &help{}
	err := json.Unmarshal(data, helpVal)
	if err != nil {
		return err
	}
	rtype, ok := ModelMap.Load(helpVal.Index)
	retype := reflect.TypeOf(rtype)
	if retype.Kind() == reflect.Pointer {
		retype = retype.Elem()
	}
	newTmp := reflect.New(retype)
	if !ok && !newTmp.CanInterface() {
		err = json.Unmarshal(data, tmp)
		if err != nil {
			return err
		}
		h.ID = tmp.ID
		h.Index = tmp.Index
		h.Source = tmp.Source
		return nil
	}

	val := newTmp.Interface()
	tmp.Source = val
	err = json.Unmarshal(data, tmp)
	if err != nil {
		return err
	}
	h.ID = tmp.ID
	h.Index = tmp.Index
	h.Source = tmp.Source
	return nil
}

type Total struct {
	Value    int64  `json:"value"`
	Relation string `json:"relation"`
}

type Shards struct {
	Total      int64 `json:"total"`
	Successful int64 `json:"successful"`
	Skipped    int64 `json:"skipped"`
	Failed     int64 `json:"failed"`
}

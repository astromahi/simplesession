//Copyright 2014 Mahendra Kathirvel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simplesession

import (
	"encoding/gob"
	"fmt"
	"testing"
)

type TestType struct {
	Id   int
	Name string
	Auth bool
}

var testData = []map[string]interface{}{
	{"id": 1001},
	{"name": "John Doe"},
	{"auth": true},
	{"type": TestType{1001, "John Doe", true}},
}

func TestGenerateId(t *testing.T) {
	list := make([]string, 100000)
	for i := 0; i < 100000; i++ {
		id, err := generateId()
		if err != nil {
			t.Error(err)
		}
		for _, value := range list {
			if id == value {
				t.Errorf("Collision: %s; Rounds: %d", id, i)
			}
		}
		list[i] = id
	}
}

func TestSerialization(t *testing.T) {
	gob.Register(TestType{})
	var err error
	var serialized []byte
	var unserialized map[string]interface{}
	for _, value := range testData {
		if serialized, err = serialize(value); err != nil {
			t.Error(err)
		}
		unserialized = make(map[string]interface{})
		if err := unserialize(serialized, unserialized); err != nil {
			t.Error(err)
		}
		if fmt.Sprintf("%+v", value) != fmt.Sprintf("%+v", unserialized) {
			t.Errorf("Expected: %+v, Got: %+v", value, unserialized)
		}
	}
}

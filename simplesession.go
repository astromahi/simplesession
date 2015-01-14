//Copyright 2014 Mahendra Kathirvel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simplesession

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const (
	sidLength = 32
	sfileDir  = "/tmp"
)

func New(res http.ResponseWriter, option *Option) (*SimpleSession, error) {
	id, err := generateId()
	if err != nil {
		return nil, err
	}

	fpath := filepath.Join(sfileDir, "gosession_"+id)

	ss := &SimpleSession{
		name:   "GSESSIONID",
		id:     id,
		fpath:  fpath,
		option: option,
		data:   make(map[string]interface{}),
	}

	cookie := &http.Cookie{
		Name:   ss.name,
		Value:  ss.id,
		Path:   ss.option.Path,
		Domain: ss.option.Domain,
		//MaxAge:   ss.option.MaxAge,
		Secure:   ss.option.Secure,
		HttpOnly: ss.option.HttpOnly,
	}
	http.SetCookie(res, cookie)

	return ss, nil
}

type SimpleSession struct {
	name   string
	id     string
	fpath  string
	option *Option
	data   map[string]interface{}
}

func (ss *SimpleSession) Name() string {
	return ss.name
}

func (ss *SimpleSession) Id() string {
	return ss.id
}

func (ss *SimpleSession) FilePath() string {
	return ss.fpath
}

func (ss *SimpleSession) Set(key string, val interface{}) {
	ss.data[key] = val
}

func (ss *SimpleSession) Get(key string) interface{} {
	if val, ok := ss.data[key]; ok {
		return val
	}
	return nil
}

func (ss *SimpleSession) Del(key string) {
	if _, ok := ss.data[key]; ok {
		delete(ss.data, key)
	}
}

// Read reads stored session
func Read(req *http.Request) (*SimpleSession, error) {
	cke, err := req.Cookie("GSESSIONID")
	if err != nil {
		return nil, errors.New("simplesession: session not set")
	}

	option := &Option{
		Path:   cke.Path,
		Domain: cke.Domain,
		//MaxAge:   cke.MaxAge,
		Secure:   cke.Secure,
		HttpOnly: cke.HttpOnly,
	}

	fpath := filepath.Join(sfileDir, "gosession_"+cke.Value)

	fp, err := os.OpenFile(fpath, os.O_RDONLY, 0400)
	if err != nil {
		return nil, errors.New("simplesession: could not open session file")
	}
	defer fp.Close()

	var serialized []byte
	buf := make([]byte, 128)
	for {
		var n int
		n, err = fp.Read(buf)
		serialized = append(serialized, buf[0:n]...)
		if err != nil || err == io.EOF {
			break
		}
	}

	data := make(map[string]interface{})
	if err = unserialize(serialized, data); err != nil {
		return nil, err
	}

	ss := &SimpleSession{
		name:   cke.Name,
		id:     cke.Value,
		fpath:  fpath,
		option: option,
		data:   data,
	}

	return ss, nil
}

func (ss *SimpleSession) Write() error {
	serialized, err := serialize(ss.data)
	if err != nil {
		return err
	}

	var fmutex sync.RWMutex

	fmutex.Lock()
	fp, err := os.OpenFile(ss.fpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return errors.New("simplesession: could not open session file")
	}

	if err = os.Chown(ss.fpath, os.Getuid(), os.Getgid()); err != nil {
		return errors.New("simplesession: could not change file owner")
	}
	defer fp.Close()

	if _, err = fp.Write(serialized); err != nil {
		return errors.New("simplesession: could not write session file")
	}
	defer fmutex.Unlock()

	return nil
}

func (ss *SimpleSession) Destroy(res http.ResponseWriter) error {

	cookie := &http.Cookie{
		Name:   ss.name,
		Value:  "",
		Path:   ss.option.Path,
		Domain: ss.option.Domain,
		//Expires: time.Unix(1, 0),
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(res, cookie)

	if err := os.Remove(ss.fpath); err != nil {
		return errors.New("simplesession: no session file found")
	}

	return nil
}

// generateId generates a session id
func generateId() (string, error) {

	hash := sha1.New()
	io.WriteString(hash, strconv.FormatInt(time.Now().Unix(), 10))

	uran := make([]byte, 1024)
	if _, err := io.ReadFull(rand.Reader, uran); err != nil {
		return "", errors.New("simplesession: could not generate random key")
	}

	hash.Write(uran)
	id := hex.EncodeToString(hash.Sum(nil))

	fpath := filepath.Join(sfileDir, "gosession_"+id)
	if _, err := os.Stat(fpath); err == nil {
		return generateId()
	} else {
		return id[0:sidLength], nil
	}
}

// serialize encodes a value using binary.
func serialize(src map[string]interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(src); err != nil {
		return nil, errors.New("simplesession: error occured while serializing session")
	}
	return buf.Bytes(), nil
}

// unserialize decodes a value using binary.
func unserialize(src []byte, dst map[string]interface{}) error {
	buf := bytes.NewReader(src)
	if err := gob.NewDecoder(buf).Decode(&dst); err != nil {
		return errors.New("simplesession: error occured while unserializing session")
	}
	return nil
}

type Option struct {
	Path   string
	Domain string
	//MaxAge   int
	Secure   bool
	HttpOnly bool
}

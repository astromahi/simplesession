package session

import (
	"github.com/astromahi/simplesession"
	"net/http"
)

var session *simplesession.SimpleSession

func Start(res http.ResponseWriter) error {

	option := &simplesession.Option{
		Path:     "/",
		Domain:   "example.com",
		MaxAge:   24 * 60 * 60,
		Secure:   false,
		HttpOnly: true,
	}

	tmpSess, err := simplesession.New(res, option)
	if err != nil {
		return err
	}
	session = tmpSess
	return nil
}

func Session() *simplesession.SimpleSession {
	return session
}

func Set(key string, val interface{}) {
	session.Set(key, val)
}

func Get(key string) interface{} {
	return session.Get(key)
}

func Del(key string) {
	session.Del(key)
}

func Write(res http.ResponseWriter, req *http.Request) error {
	if err := session.Write(res, req); err != nil {
		return err
	}
	return nil
}

func Read(req *http.Request) error {
	ss, err := simplesession.Read(req)
	if err != nil {
		return err
	}
	session = ss
	return nil
}

func Destroy(res http.ResponseWriter) error {
	if err := session.Destroy(res); err != nil {
		return err
	}
	return nil
}

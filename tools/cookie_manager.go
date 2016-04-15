package tools

import (
	"errors"
	"github.com/gorilla/securecookie"
	"net/http"
	"strings"
	"time"
)

var cookieManger = &CookieManger{securecookie.New(hashKey, blockKey)}

type CookieManger struct {
	man *securecookie.SecureCookie
}

func GetCookieManager() *CookieManger {
	return cookieManger
}

func (cm *CookieManger) Set(rw http.ResponseWriter, name, value string, options map[string]interface{}) error {
	encoded, err := cm.man.Encode(name, value)
	if nil != err {
		return err
	}

	c := &http.Cookie{
		Name:  name,
		Value: encoded,
	}

	for k, v := range options {
		lowerAttr := strings.ToLower(k)
		switch lowerAttr {
		case "secure":
			c.Secure = v.(bool)
			continue
		case "httponly":
			c.HttpOnly = v.(bool)
			continue
		case "domain":
			c.Domain = v.(string)
			continue
		case "path":
			c.Path = v.(string)
			continue
		case "max-age":
		case "maxage":
			c.MaxAge = v.(int)
			continue
		case "expires":
			switch t := v.(type) {
			case time.Time:
				c.Expires = t
				break
			default:
				return errors.New("Expires should be a time")
			}
		}
	}

	http.SetCookie(rw, c)

	return nil
}

func (cm *CookieManger) Get(r *http.Request, name string) (string, error) {
	c, err := r.Cookie(name)
	if nil != err {
		return "", err
	}

	var value string
	err = cm.man.Decode(name, c.Value, &value)
	if nil != err {
		return "", err
	}

	return value, nil
}

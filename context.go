package middleware

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	//"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
)

const (
	DefaultKey    = "github.com/forease/gin-middleware"
	sessManageKey = "github.com/forease/gin-middleware-session-manager-key"
	sessStoreKey  = "github.com/forease/gin-middleware-session-store-key"
	errorFormat   = "[sessions] ERROR! %s\n"
)

type G map[string]interface{}

type Context struct {
	*gin.Context
	data     G
	dataLock sync.RWMutex
}

func NewContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := &Context{Context: c, data: make(G), dataLock: sync.RWMutex{}}
		c.Set(DefaultKey, s)
	}
}

func NewSession(opt ...session.Option) gin.HandlerFunc {
	manage := session.NewManager(opt...)
	return func(c *gin.Context) {

		c.Set(sessManageKey, manage)
		store, err := manage.Start(context.Background(), c.Writer, c.Request)
		if err != nil {

			return
		}
		c.Set(sessStoreKey, store)
		c.Next()
	}
}

// data
func (self *Context) Get(key string) interface{} {
	self.dataLock.RLock()
	defer self.dataLock.RUnlock()

	return self.data[key]
}

func (self *Context) Set(key string, val interface{}) {
	self.dataLock.Lock()
	if self.data == nil {
		self.data = make(G)
	}
	self.data[key] = val
	self.dataLock.Unlock()
}

func (self *Context) Data() G {
	return self.data
}

func (self *Context) Delete(key string) {
	self.dataLock.Lock()
	delete(self.data, key)
	self.dataLock.Unlock()
}

func (self *Context) Clear() {
	self.data = make(G)
}

// sessions
func (self *Context) Session() session.Store {
	return self.MustGet(sessStoreKey).(session.Store) //sessions.Default(self.Context)
}

func (self *Context) SessSet(key string, value interface{}) {
	self.Session().Set(key, value)
	//sessions.Default(self.Context).Set(key, value)
}

func (self *Context) SessGetValue(key string) (interface{}, bool) {
	return self.Session().Get(key)
	//return sessions.Default(self.Context).Get(key)
}

func (self *Context) SessGet(key string, value interface{}) (err error) {
	//tmp := sessions.Default(self.Context).Get(key)
	tmp, has := self.Session().Get(key)

	if tmp == nil || !has {
		return fmt.Errorf("Can't found session value for %s", key)
	}

	sType := getTypeOf(tmp)
	vType := getTypeOf(value)
	if sType != vType {
		return fmt.Errorf("Can't match value type (%s != %s)", sType, vType)
	}

	//
	refValue := reflect.Indirect(reflect.ValueOf(value))
	refValue.Set(reflect.Indirect(reflect.ValueOf(tmp)))

	return
}

func (self *Context) SessDelete(key string) interface{} {
	return self.Session().Delete(key)
}

func (self *Context) SessSave() {
	self.Session().Save()
	//sessions.Default(self.Context).Save()
}

func (self *Context) SessClear() error {
	//return self.Session().Flush()
	//sessions.Default(self.Context).Clear()
	return self.MustGet(sessManageKey).(*session.Manager).Destroy(context.Background(), self.Writer, self.Request)
}

// message
// Render JSON message
func (self *Context) JMessage(code int, url, message string, v ...interface{}) {
	if len(v) > 0 {
		message = fmt.Sprintf(message, v...)
	}
	msg := Msg{Code: code, Message: message, Url: url}
	self.JSON(200, msg)
}

// Render HTML message
func (self *Context) HMessage(code int, url, message string, v ...interface{}) {
	if len(v) > 0 {
		message = fmt.Sprintf(message, v...)
	}
	self.Set("msg", Msg{Code: code, Message: message, Url: url})
	self.HTML(200, "layout/message", self.Data())
}

// get Datatable
func (self *Context) DataTable(draw, total, datas interface{}) G {

	return G{"draw": draw, "recordsTotal": total,
		"recordsFiltered": total,
		"data":            datas,
	}
}

// shortcut to get context
func Default(c *gin.Context) *Context {
	return c.MustGet(DefaultKey).(*Context)
}

func getTypeOf(val interface{}) (typeName string) {

	tp := reflect.TypeOf(val)
	typeName = tp.Kind().String()

	if typeName == "ptr" {
		typeName = tp.Elem().Kind().String()
	}

	return
}

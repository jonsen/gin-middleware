package middleware

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
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

// G -
type G map[string]interface{}

// Context -
type Context struct {
	*gin.Context
	data     G
	dataLock sync.RWMutex
}

// NewContext -
func NewContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := &Context{Context: c, data: make(G), dataLock: sync.RWMutex{}}
		c.Set(DefaultKey, s)
	}
}

// NewSession -
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

// Get data
func (c *Context) Get(key string) interface{} {
	c.dataLock.RLock()
	defer c.dataLock.RUnlock()

	return c.data[key]
}

// Set -
func (c *Context) Set(key string, val interface{}) {
	c.dataLock.Lock()
	if c.data == nil {
		c.data = make(G)
	}
	c.data[key] = val
	c.dataLock.Unlock()
}

// Data -
func (c *Context) Data() G {
	return c.data
}

// Delete -
func (c *Context) Delete(key string) {
	c.dataLock.Lock()
	delete(c.data, key)
	c.dataLock.Unlock()
}

// Clear -
func (c *Context) Clear() {
	c.data = make(G)
}

// Session sessions
func (c *Context) Session() session.Store {
	return c.MustGet(sessStoreKey).(session.Store) //sessions.Default(c.Context)
}

// SessSet -
func (c *Context) SessSet(key string, value interface{}) {
	c.Session().Set(key, value)
	//sessions.Default(c.Context).Set(key, value)
}

// SessGetValue -
func (c *Context) SessGetValue(key string) (interface{}, bool) {
	return c.Session().Get(key)
	//return sessions.Default(c.Context).Get(key)
}

// SessGet -
func (c *Context) SessGet(key string, value interface{}) (err error) {
	//tmp := sessions.Default(c.Context).Get(key)
	tmp, has := c.Session().Get(key)

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

// SessDelete -
func (c *Context) SessDelete(key string) interface{} {
	return c.Session().Delete(key)
}

// SessSave -
func (c *Context) SessSave() {
	c.Session().Save()
	//sessions.Default(c.Context).Save()
}

// SessClear -
func (c *Context) SessClear() error {
	//return c.Session().Flush()
	//sessions.Default(c.Context).Clear()
	return c.MustGet(sessManageKey).(*session.Manager).Destroy(context.Background(), c.Writer, c.Request)
}

// JMessage message
// Render JSON message
func (c *Context) JMessage(code int, url, message string, v ...interface{}) {
	if len(v) > 0 {
		message = fmt.Sprintf(message, v...)
	}
	msg := Msg{Code: code, Message: message, Url: url}
	c.JSON(200, msg)
}

// HMessage Render HTML message
func (c *Context) HMessage(code int, url, message string, v ...interface{}) {
	if len(v) > 0 {
		message = fmt.Sprintf(message, v...)
	}
	c.Set("msg", Msg{Code: code, Message: message, Url: url})
	c.HTML(200, "layout/message", c.Data())
}

// DataTable get Datatable
func (c *Context) DataTable(draw, total, datas interface{}) G {

	return G{"draw": draw, "recordsTotal": total,
		"recordsFiltered": total,
		"data":            datas,
	}
}

// Default shortcut to get context
func Default(c *gin.Context) *Context {
	ctx, ok := c.Get(DefaultKey)
	if !ok {
		s := &Context{Context: c, data: make(G), dataLock: sync.RWMutex{}}
		c.Set(DefaultKey, s)
		return s
	}

	return ctx.(*Context)
}

func getTypeOf(val interface{}) (typeName string) {

	tp := reflect.TypeOf(val)
	typeName = tp.Kind().String()

	if typeName == "ptr" {
		typeName = tp.Elem().Kind().String()
	}

	return
}

//
// query value
//
func (c *Context) QueryIntDefault(key string, defaultValue int) int {
	return int(c.QueryInt64Default(key, int64(defaultValue)))
}

func (c *Context) QueryInt64Default(key string, defaultValue int64) int64 {
	value, _ := c.GetQuery(key)

	if value != "" {
		v, _ := strconv.ParseInt(value, 10, 64)
		return v
	}

	return defaultValue
}

//
// Param
//
func (c *Context) ParamInt64Default(key string, defaultValue int64) int64 {
	value := c.Params.ByName(key)
	v, _ := strconv.ParseInt(value, 10, 64)

	return v
}

func (c *Context) ParamDefault(key, defaultValue string) string {
	value := c.Params.ByName(key)
	if value != "" {
		return value
	}

	return defaultValue
}

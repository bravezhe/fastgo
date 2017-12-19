package fastgo

import (
	"net/http"
	"encoding/json"
	"strconv"
)

type Controller struct {
	rw http.ResponseWriter
	r *http.Request
}


func (c *Controller) Prepare(rw http.ResponseWriter, r *http.Request) {
	c.rw = rw
	c.r = r
}

func (c *Controller) GetString(key string, defaultVal string) string {
	value := c.r.FormValue(key)
	if value == "" {
		value = defaultVal
	}
	return value
}

func (c *Controller) GetStrings(key string, defaultVal ...[]string) []string {
	if c.r.Form == nil {
		return []string{}
	}
	value := c.r.Form[key]
	return value
}

func (c *Controller) GetInt(key string, defaultVal int64) int64 {
	value, err := strconv.ParseInt(c.r.FormValue(key), 10, 64)
	if err != nil {
		value = defaultVal
	}
	return value
}

func (c *Controller) GetFloat(key string, defaultVal float64) float64 {
	value, err := strconv.ParseFloat(c.r.FormValue(key), 64)
	if err != nil {
		value = defaultVal
	}
	return value
}

func (c *Controller) GetBool(key string, defaultVal bool) bool {
	value, err := strconv.ParseBool(c.r.FormValue(key))
	if err != nil {
		value = defaultVal
	}
	return value
}

func (c *Controller) RenderJson(data interface{}) {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic("json error")
	}
	c.rw.Header().Set("Content-Type", "application/json;charset=UTF-8")
	c.Write(content)
}

func (c *Controller) Write(bt []byte) {
	c.rw.Write(bt)
}
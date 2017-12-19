package fastgo

import (
	"reflect"
	"net/http"
	"strings"

	"runtime"
	"fmt"
	"time"
)


type HttpServer struct {
	Addr		string
	Port		int
	Timeout		int
	Handler 	*HttpHandler
}

func InitServer(addr string, port int, timeout int) *HttpServer {
	ret := &HttpServer {
		Addr:		addr,
		Port:		port,
		Timeout:	timeout,
		Handler:	&HttpHandler{routerMap:make(map[string]map[string]reflect.Type)},
	}
	return ret
}

func (s *HttpServer) Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	addr := fmt.Sprintf("%s:%d", s.Addr, s.Port)
	server := &http.Server{
		Addr:			addr,
		Handler:		s.Handler,
		ReadTimeout:	time.Duration(s.Timeout) * time.Millisecond,
		WriteTimeout:	time.Duration(s.Timeout) * time.Millisecond,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (s *HttpServer) AddController (c interface{}) {
	s.Handler.addController(c)
}

type HttpHandler struct {
	routerMap	map[string]map[string]reflect.Type
}

func (h *HttpHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(rw, fmt.Sprintln(err), http.StatusInternalServerError)
		}
	}()
	currentUrl := r.URL.Path
	splitUrl := strings.Split(currentUrl, "/")
	eleNum := len(splitUrl)
	var controller string
	var method string
	if eleNum == 0 || eleNum == 1 {
		controller = "Default"
		method = "Index"
	} else if eleNum == 2 {
		controller = strings.Title(splitUrl[1])
		method = "Index"
	} else if eleNum == 3 {
		controller = strings.Title(splitUrl[1])
		method = strings.Title(splitUrl[2])
	}
	var methodType reflect.Type
	var authFlag = false
	if controllerMap, ok := h.routerMap[controller]; ok {
		if methodType, ok = controllerMap[method]; ok {
			authFlag = true
		}
	}
	if !authFlag {
		http.NotFound(rw, r)
		return
	}

	var in []reflect.Value
	var methodAction reflect.Value

	methodValue := reflect.New(methodType)
	in = make([]reflect.Value, 2)
	in[0] = reflect.ValueOf(rw)
	in[1] = reflect.ValueOf(r)
	methodAction = methodValue.MethodByName("Prepare")
	methodAction.Call(in)

	in = make([]reflect.Value, 0)
	methodAction = methodValue.MethodByName(method)
	methodAction.Call(in)
}

func (h *HttpHandler) addController(c interface{}) {
	reflectVal := reflect.ValueOf(c)
	reflectType := reflectVal.Type()
	reflectRealType := reflect.Indirect(reflectVal).Type()
	controller := strings.TrimSuffix(reflectRealType.Name(), "Controller")
	if _, ok := h.routerMap[controller]; ok {
		return
	} else {
		h.routerMap[controller] = make(map[string]reflect.Type)
	}
	for i:=0; i < reflectType.NumMethod(); i++ {
		method := reflectType.Method(i).Name
		h.routerMap[controller][method] = reflectRealType
	}
}
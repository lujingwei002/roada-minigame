package gat

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/roada-go/gat/log"
)

type (
	handlerMethod struct {
		method   reflect.Method
		typ      reflect.Type
		isRawArg bool
	}

	Service struct {
		gate     *Gate
		typ      reflect.Type
		handlers map[string]*handlerMethod
	}
)

var (
	typeOfError   = reflect.TypeOf((*error)(nil)).Elem()
	typeOfBytes   = reflect.TypeOf(([]byte)(nil))
	typeOfRequest = reflect.TypeOf(&Request{})
)

func isExported(name string) bool {
	w, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(w)
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return isExported(t.Name()) || t.PkgPath() == ""
}

func isHandlerMethod(method reflect.Method) bool {
	mt := method.Type
	if method.PkgPath != "" {
		return false
	}
	if mt.NumIn() != 3 {
		return false
	}
	if mt.NumOut() != 1 {
		return false
	}
	if t1 := mt.In(1); t1.Kind() != reflect.Ptr || t1 != typeOfRequest {
		return false
	}
	if (mt.In(2).Kind() != reflect.Ptr && mt.In(2) != typeOfBytes) || mt.Out(0) != typeOfError {
		return false
	}
	return true
}

func newService(gate *Gate, svr interface{}) *Service {
	s := &Service{
		gate: gate,
		typ:  reflect.TypeOf(svr),
	}
	return s
}

func (self *Service) suitableHandlerMethods(typ reflect.Type) map[string]*handlerMethod {
	methods := make(map[string]*handlerMethod)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mt := method.Type
		mn := method.Name
		if isHandlerMethod(method) {
			raw := false
			if mt.In(2) == typeOfBytes {
				raw = true
			}
			mn = strings.ToLower(mn)
			methods[mn] = &handlerMethod{method: method, typ: mt.In(2), isRawArg: raw}
		}
	}
	return methods
}

func (self *Service) extractHandler() error {
	self.handlers = self.suitableHandlerMethods(self.typ)
	return nil
}

func (self *Service) Unpack(receiver interface{}, r *Request) {
	route := r.Route
	index := strings.LastIndex(route, ".")
	if index < 0 {
		log.Println(fmt.Sprintf("service invalid route, route:%s", route))
		return
	}
	handlerName := route[index+1:]
	handler, found := self.handlers[handlerName]
	if !found {
		log.Println(fmt.Sprintf("service handler not found, route:%s", route))
		return
	}
	var data interface{}
	if handler.isRawArg {
		data = r.Payload
	} else {
		data = reflect.New(handler.typ.Elem()).Interface()
		err := self.gate.serializer.Unmarshal(r.Payload, data)
		if err != nil {
			log.Println(fmt.Sprintf("Deserialize to %T failed: %+v (%v)", data, err, r.Payload))
			return
		}
	}
	args := []reflect.Value{reflect.ValueOf(receiver), reflect.ValueOf(r), reflect.ValueOf(data)}
	result := handler.method.Func.Call(args)
	if len(result) > 0 {
		if err := result[0].Interface(); err != nil {
			log.Println(fmt.Sprintf("Service %s error: %+v", route, err))
		}
	}
}

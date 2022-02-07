package cli

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/roada-go/gat/log"
)

type (
	handlerMethod struct {
		receiver reflect.Value
		method   reflect.Method
		typ      reflect.Type
		isRawArg bool
	}

	Service struct {
		client   *Client
		typ      reflect.Type
		receiver reflect.Value
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

func newService(client *Client, svr interface{}) *Service {
	s := &Service{
		client:   client,
		typ:      reflect.TypeOf(svr),
		receiver: reflect.ValueOf(svr),
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
	typeName := reflect.Indirect(self.receiver).Type().Name()
	if typeName == "" {
		return errors.New("no service name for type " + self.typ.String())
	}
	if !isExported(typeName) {
		return errors.New("type " + typeName + " is not exported")
	}
	self.handlers = self.suitableHandlerMethods(self.typ)
	for i := range self.handlers {
		self.handlers[i].receiver = self.receiver
	}
	return nil
}

func (self *Service) Unpack(receiver interface{}, r *Request) {
	payload := r.Payload
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
		data = payload
	} else {
		data = reflect.New(handler.typ.Elem()).Interface()
		err := self.client.serializer.Unmarshal(payload, data)
		if err != nil {
			log.Println(fmt.Sprintf("Deserialize to %T failed: %+v (%v)", data, err, payload))
			return
		}
	}
	args := []reflect.Value{reflect.ValueOf(receiver), reflect.ValueOf(r), reflect.ValueOf(data)}
	//task := func() {
	result := handler.method.Func.Call(args)
	if len(result) > 0 {
		if err := result[0].Interface(); err != nil {
			log.Println(fmt.Sprintf("Service %s error: %+v", route, err))
		}
	}
	//}
	//self.gate.scheduler.PushTask(task)
}

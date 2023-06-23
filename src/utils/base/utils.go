package base

import "reflect"

func getStructName(s interface{}) string {
	if t := reflect.TypeOf(s); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func GetProviderKey(s interface{}) string {
	return reflect.TypeOf(s).String()
}

func Bind(interfaceValue interface{}, structValue interface{}) interface{} {
	interfaceType := reflect.TypeOf(interfaceValue)
	if interfaceType.Kind() == reflect.Ptr {
		interfaceType = interfaceType.Elem()
	}
	if interfaceType.Kind() != reflect.Interface {
		panic("the type of parameter 'interfaceValue' should be interface")
	}

	in := []reflect.Type{reflect.TypeOf(structValue)}
	out := []reflect.Type{interfaceType}
	funcType := reflect.FuncOf(in, out, false)

	return reflect.MakeFunc(funcType, func(args []reflect.Value) (results []reflect.Value) {
		return args
	}).Interface()
}

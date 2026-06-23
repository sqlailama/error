package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime/debug"
)

type Values []any

type Result struct {
	Values
	err error
}

func Some(values ...any) Result {
	sz := len(values)
	if sz == 0 {
		panic("On: nothing to do")
	}

	var err error
	if e, ok := values[sz-1].(error); ok {
		err = e
		values = values[:sz-1]
	}

	return Result{
		Values: values,
		err:    err,
	}
}

func If[T any](pred bool, vt, vf T) T {
	if pred {
		return vt
	}
	return vf
}

func (r Result) Log() Result {
	log.Println(If(r.err != nil, r.err, fmt.Errorf("operations successfully completed")))
	return r
}

func (r Result) Must() Values {
	if r.err != nil {
		panic(r.err)
	}
	return r.Values
}

func (v Values) Unpack(vars ...any) {
	for i, itm := range vars {
		if i >= len(v) {
			return
		}

		rv := reflect.ValueOf(itm)
		if rv.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("Unpack: arg %d is not pointer", i))
		}

		val := reflect.ValueOf(v[i])
		if !val.Type().AssignableTo(rv.Elem().Type()) {
			panic(fmt.Sprintf(
				"Unpack: cannot assign %s to %s", val.Type(), rv.Elem().Type(),
			))
		}

		rv.Elem().Set(val)
	}
}

func RewindAndExit() {
	if err := recover(); err != nil {
		debug.PrintStack()
		os.Exit(1)
	}
}

func demoFunc(v any) (int, string, error) {
	if v == nil {
		return 0, "", fmt.Errorf("unknown result")
	}
	return 47, "success", nil
}

func main() {
	defer RewindAndExit()

	var (
		i int
		s string
	)

	Some(demoFunc("ok")).Must().Unpack(&i, &s)
	fmt.Println(i, s)

	var f *os.File
	Some(os.Open("LICENSE")).Log().Must().Unpack(&f)
	defer func() {
		log.Println("file is successfully closed")
		_ = f.Close()
	}()
}

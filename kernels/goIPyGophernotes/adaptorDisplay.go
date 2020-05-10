package main

// (proto)Adaptor Display

import (
//	"bytes"
//	"errors"
	"fmt"
	"image"
//	"io"
//	"io/ioutil"
//	"net/http"
	"reflect"
//	"strings"

	basereflect "github.com/cosmos72/gomacro/base/reflect"
	"github.com/cosmos72/gomacro/xreflect"
)

// fill kernel.renderer map used to convert interpreted types
// to known rendering interfaces
func (kernel *GoInterpreter) initRenderers() {
	kernel.render = make(map[string]xreflect.Type)
	for name, typ := range kernel.display.Types {
		if typ.Kind() == reflect.Interface {
			kernel.render[name] = typ
		}
	}
}

// if vals[] contain a single non-nil value which is auto-renderable,
// convert it to Data and return it.
// otherwise return MakeData("text/plain", fmt.Sprint(vals...))
func (kernel *GoInterpreter) autoRenderResults(vals []interface{}, types []xreflect.Type) Data {
	var nilcount int
	var obj interface{}
	var typ xreflect.Type
	for i, val := range vals {
		if kernel.canAutoRender(val, types[i]) {
			obj = val
			typ = types[i]
		} else if val == nil {
			nilcount++
		}
	}
	if obj != nil && nilcount == len(vals)-1 {
		return kernel.autoRender("", obj, typ)
	}
	if nilcount == len(vals) {
		// if all values are nil, return empty Data
		return Data{}
	}
	return MakeData(MIMETypeText, fmt.Sprint(vals...))
}

// return true if data type should be auto-rendered graphically
func (kernel *GoInterpreter) canAutoRender(data interface{}, typ xreflect.Type) bool {
	switch data.(type) {
	case Data, Renderer, SimpleRenderer, HTMLer, JavaScripter, JPEGer, JSONer,
		Latexer, Markdowner, PNGer, PDFer, SVGer, image.Image:
		return true
	}
	if kernel == nil || typ == nil {
		return false
	}
	// in gomacro, methods of interpreted types are emulated,
	// thus type-asserting them to interface types as done above cannot succeed.
	// Manually check if emulated type "pretends" to implement
	// at least one of the interfaces above
	for _, xtyp := range kernel.render {
		if typ.Implements(xtyp) {
			return true
		}
	}
	return false
}

// detect and render data types that should be auto-rendered graphically
func (kernel *GoInterpreter) autoRender(mimeType string, arg interface{}, typ xreflect.Type) Data {
	var data Data
	// try Data
	if x, ok := arg.(Data); ok {
		data = x
	}

	if kernel == nil || typ == nil {
		// try all autoRenderers
		for _, fun := range autoRenderers {
			data = fun(data, arg)
		}
	} else {
		// in gomacro, methods of interpreted types are emulated.
		// Thus type-asserting them to interface types as done by autoRenderer functions above cannot succeed.
		// Manually check if emulated type "pretends" to implement one or more of the above interfaces
		// and, in case, tell the interpreter to convert to them
		for name, xtyp := range kernel.render {
			fun := autoRenderers[name]
			if fun == nil || !typ.Implements(xtyp) {
				continue
			}
			conv := kernel.ir.Comp.Converter(typ, xtyp)
			x := arg
			if conv != nil {
				x = basereflect.Interface(conv(reflect.ValueOf(x)))
				if x == nil {
					continue
				}
			}
			data = fun(data, x)
		}
	}
	return fillDefaults(data, arg, "", nil, "", nil)
}

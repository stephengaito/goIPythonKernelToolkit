package goIPyGoMacroAdaptor

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

  tk "github.com/stephengaito/goIPythonKernelToolkit/goIPyKernel"
  
	basereflect "github.com/cosmos72/gomacro/base/reflect"
	"github.com/cosmos72/gomacro/xreflect"
)

/**
 * general interface, allows libraries to fully specify
 * how their data is displayed by Jupyter.
 * Supports multiple MIME formats.
 *
 * Note that Data defined above is an alias:
 * libraries can implement Renderer without importing goIPyGophernotes
 */
type Renderer = interface {
	Render() tk.Data
}

/**
 * simplified interface, allows libraries to specify
 * how their data is displayed by Jupyter.
 * Supports multiple MIME formats.
 *
 * Note that MIMEMap defined above is an alias:
 * libraries can implement SimpleRenderer without importing goIPyGophernotes
 */
type SimpleRenderer = interface {
	SimpleRender() tk.MIMEMap
}

/**
 * specialized interfaces, each is dedicated to a specific MIME type.
 *
 * They are type aliases to emphasize that method signatures
 * are the only important thing, not the interface names.
 * Thus libraries can implement them without importing goIPyGophernotes
 */
type HTMLer = interface {
	HTML() string
}
type JavaScripter = interface {
	JavaScript() string
}
type JPEGer = interface {
	JPEG() []byte
}
type JSONer = interface {
	JSON() map[string]interface{}
}
type Latexer = interface {
	Latex() string
}
type Markdowner = interface {
	Markdown() string
}
type PNGer = interface {
	PNG() []byte
}
type PDFer = interface {
	PDF() []byte
}
type SVGer = interface {
	SVG() string
}

// injected as placeholder in the interpreter, it's then replaced at runtime
// by a closure that knows how to talk with Jupyter
func stubDisplay(tk.Data) error {
	return errors.New("cannot display: connection with Jupyter not available")
}

// fill kernel.renderer map used to convert interpreted types
// to known rendering interfaces
func (adaptor *GoAdaptor) initRenderers() {
	adaptor.render = make(map[string]xreflect.Type)
	for name, typ := range adaptor.display.Types {
		if typ.Kind() == reflect.Interface {
			adaptor.render[name] = typ
		}
	}
}

// if vals[] contain a single non-nil value which is auto-renderable,
// convert it to Data and return it.
// otherwise return MakeData("text/plain", fmt.Sprint(vals...))
//
func (adaptor *GoAdaptor) autoRenderResults(
  vals []interface{},
  types []xreflect.Type,
) tk.Data {
	var nilcount int
	var obj interface{}
	var typ xreflect.Type
	for i, val := range vals {
		if adaptor.canAutoRender(val, types[i]) {
			obj = val
			typ = types[i]
		} else if val == nil {
			nilcount++
		}
	}
	if obj != nil && nilcount == len(vals)-1 {
		return adaptor.autoRender("", obj, typ)
	}
	if nilcount == len(vals) {
		// if all values are nil, return empty Data
		return tk.Data{}
	}
	return MakeData(tk.MIMETypeText, fmt.Sprint(vals...))
}

// return true if data type should be auto-rendered graphically
func (adaptor *GoAdaptor) canAutoRender(data interface{}, typ xreflect.Type) bool {
	switch data.(type) {
	case
    tk.Data,
    Renderer,
    SimpleRenderer,
    HTMLer,
    JavaScripter,
    JPEGer,
    JSONer,
		Latexer,
    Markdowner,
    PNGer,
    PDFer,
    SVGer,
    image.Image :
		return true
	}
	if adaptor == nil || typ == nil {
		return false
	}
	// in gomacro, methods of interpreted types are emulated,
	// thus type-asserting them to interface types as done above cannot succeed.
	// Manually check if emulated type "pretends" to implement
	// at least one of the interfaces above
	for _, xtyp := range adaptor.render {
		if typ.Implements(xtyp) {
			return true
		}
	}
	return false
}

var autoRenderers = map[string]func(tk.Data, interface{}) tk.Data{
	"Renderer": func(d tk.Data, i interface{}) tk.Data {
		if r, ok := i.(Renderer); ok {
			x := r.Render()
			d.Data = tk.MergeMIMEMap(d.Data, x.Data)
			d.Metadata = tk.MergeMIMEMap(d.Metadata, x.Metadata)
			d.Transient = tk.MergeMIMEMap(d.Transient, x.Transient)
		}
		return d
	},
	"SimpleRenderer": func(d tk.Data, i interface{}) tk.Data {
		if r, ok := i.(SimpleRenderer); ok {
			x := r.SimpleRender()
			d.Data = tk.MergeMIMEMap(d.Data, x)
		}
		return d
	},
	"HTMLer": func(d tk.Data, i interface{}) tk.Data {
		if r, ok := i.(HTMLer); ok {
			d.Data = tk.EnsureMIMEMap(d.Data)
			d.Data[tk.MIMETypeHTML] = r.HTML()
		}
		return d
	},
	"JavaScripter": func(d tk.Data, i interface{}) tk.Data {
		if r, ok := i.(JavaScripter); ok {
			d.Data = tk.EnsureMIMEMap(d.Data)
			d.Data[tk.MIMETypeJavaScript] = r.JavaScript()
		}
		return d
	},
	"JPEGer": func(d tk.Data, i interface{}) tk.Data {
		if r, ok := i.(JPEGer); ok {
			d.Data = tk.EnsureMIMEMap(d.Data)
			d.Data[tk.MIMETypeJPEG] = r.JPEG()
		}
		return d
	},
	"JSONer": func(d tk.Data, i interface{}) tk.Data {
		if r, ok := i.(JSONer); ok {
			d.Data = tk.EnsureMIMEMap(d.Data)
			d.Data[tk.MIMETypeJSON] = r.JSON()
		}
		return d
	},
	"Latexer": func(d tk.Data, i interface{}) tk.Data {
		if r, ok := i.(Latexer); ok {
			d.Data = tk.EnsureMIMEMap(d.Data)
			d.Data[tk.MIMETypeLatex] = r.Latex()
		}
		return d
	},
	"Markdowner": func(d tk.Data, i interface{}) tk.Data {
		if r, ok := i.(Markdowner); ok {
			d.Data = tk.EnsureMIMEMap(d.Data)
			d.Data[tk.MIMETypeMarkdown] = r.Markdown()
		}
		return d
	},
	"PNGer": func(d tk.Data, i interface{}) tk.Data {
		if r, ok := i.(PNGer); ok {
			d.Data = tk.EnsureMIMEMap(d.Data)
			d.Data[tk.MIMETypePNG] = r.PNG()
		}
		return d
	},
	"PDFer": func(d tk.Data, i interface{}) tk.Data {
		if r, ok := i.(PDFer); ok {
			d.Data = tk.EnsureMIMEMap(d.Data)
			d.Data[tk.MIMETypePDF] = r.PDF()
		}
		return d
	},
	"SVGer": func(d tk.Data, i interface{}) tk.Data {
		if r, ok := i.(SVGer); ok {
			d.Data = tk.EnsureMIMEMap(d.Data)
			d.Data[tk.MIMETypeSVG] = r.SVG()
		}
		return d
	},
	"Image": func(d tk.Data, i interface{}) tk.Data {
		if r, ok := i.(image.Image); ok {
			b, mimeType, err := encodePng(r)
			if err != nil {
				d = makeDataErr(err)
			} else {
				d.Data = tk.EnsureMIMEMap(d.Data)
				d.Data[mimeType] = b
				d.Metadata = tk.MergeMIMEMap(d.Metadata, imageMetadata(r))
			}
		}
		return d
	},
}

// detect and render data types that should be auto-rendered graphically
func (adaptor *GoAdaptor) autoRender(
  mimeType string,
  arg interface{},
  typ xreflect.Type,
) tk.Data {
	var data tk.Data
	// try Data
	if x, ok := arg.(tk.Data); ok {
		data = x
	}

	if adaptor == nil || typ == nil {
		// try all autoRenderers
		for _, fun := range autoRenderers {
			data = fun(data, arg)
		}
	} else {
  
    // in gomacro, methods of interpreted types are emulated. Thus 
    // type-asserting them to interface types as done by autoRenderer 
    // functions above cannot succeed. Manually check if emulated type 
    // "pretends" to implement one or more of the above interfaces and, in 
    // case, tell the interpreter to convert to them 
    //
    for name, xtyp := range adaptor.render {
			fun := autoRenderers[name]
			if fun == nil || !typ.Implements(xtyp) {
				continue
			}
			conv := adaptor.ir.Comp.Converter(typ, xtyp)
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

func fillDefaults(
  data tk.Data,
  arg interface{},
  s string,
  b []byte,
  mimeType string,
  err error,
) tk.Data {
	if err != nil {
		return makeDataErr(err)
	}
	if data.Data == nil {
		data.Data = make(tk.MIMEMap)
	}
	// cannot autodetect the mime type of a string
	if len(s) != 0 && len(mimeType) != 0 {
		data.Data[mimeType] = s
	}
	// ensure plain text is set
	if data.Data[tk.MIMETypeText] == "" {
		if len(s) == 0 {
			s = fmt.Sprint(arg)
		}
		data.Data[tk.MIMETypeText] = s
	}
	// if []byte is available, use it
	if len(b) != 0 {
		if len(mimeType) == 0 {
			mimeType = http.DetectContentType(b)
		}
		if len(mimeType) != 0 && mimeType != tk.MIMETypeText {
			data.Data[mimeType] = b
		}
	}
	return data
}

// do our best to render data graphically
func render(mimeType string, data interface{}) tk.Data {
	var adaptor *GoAdaptor // intentionally nil
	if adaptor.canAutoRender(data, nil) {
		return adaptor.autoRender(mimeType, data, nil)
	}
	var s string
	var b []byte
	var err error
	switch data := data.(type) {
	case string:
		s = data
	case []byte:
		b = data
	case io.Reader:
		b, err = ioutil.ReadAll(data)
	case io.WriterTo:
		var buf bytes.Buffer
		data.WriteTo(&buf)
		b = buf.Bytes()
	default:
		panic(fmt.Errorf("unsupported type, cannot render: %T", data))
	}
	return fillDefaults(tk.Data{}, data, s, b, mimeType, err)
}

func makeDataErr(err error) tk.Data {
	return tk.Data{
		Data: tk.MIMEMap{
			"ename":     "ERROR",
			"evalue":    err.Error(),
			"traceback": nil,
			"status":    "error",
		},
	}
}

func Any(mimeType string, data interface{}) tk.Data {
	return render(mimeType, data)
}

// same as Any("", data), autodetects MIME type
func Auto(data interface{}) tk.Data {
	return render("", data)
}

func MakeData(mimeType string, data interface{}) tk.Data {
	d := tk.Data{
		Data: tk.MIMEMap{
			mimeType: data,
		},
	}
	if mimeType != tk.MIMETypeText {
		d.Data[tk.MIMETypeText] = fmt.Sprint(data)
	}
	return d
}

func MakeData3(mimeType string, plaintext string, data interface{}) tk.Data {
	return tk.Data{
		Data: tk.MIMEMap{
			tk.MIMETypeText: plaintext,
			mimeType:     data,
		},
	}
}

func File(mimeType string, path string) tk.Data {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return Any(mimeType, bytes)
}

func HTML(html string) tk.Data {
	return MakeData(tk.MIMETypeHTML, html)
}

func JavaScript(javascript string) tk.Data {
	return MakeData(tk.MIMETypeJavaScript, javascript)
}

func JPEG(jpeg []byte) tk.Data {
	return MakeData(tk.MIMETypeJPEG, jpeg)
}

func JSON(json map[string]interface{}) tk.Data {
	return MakeData(tk.MIMETypeJSON, json)
}

func Latex(latex string) tk.Data {
	return MakeData3(tk.MIMETypeLatex, latex, "$"+strings.Trim(latex, "$")+"$")
}

func Markdown(markdown string) tk.Data {
	return MakeData(tk.MIMETypeMarkdown, markdown)
}

func Math(latex string) tk.Data {
	return MakeData3(tk.MIMETypeLatex, latex, "$$"+strings.Trim(latex, "$")+"$$")
}

func PDF(pdf []byte) tk.Data {
	return MakeData(tk.MIMETypePDF, pdf)
}

func PNG(png []byte) tk.Data {
	return MakeData(tk.MIMETypePNG, png)
}

func SVG(svg string) tk.Data {
	return MakeData(tk.MIMETypeSVG, svg)
}

// MIME encapsulates the data and metadata into a Data.
// The 'data' map is expected to contain at least one {key,value} pair,
// with value being a string, []byte or some other JSON serializable representation,
// and key equal to the MIME type of such value.
// The exact structure of value is determined by what the frontend expects.
// Some easier-to-use functions for common formats supported by the Jupyter frontend
// are provided by the various functions above.
func MIME(data, metadata tk.MIMEMap) tk.Data {
	return tk.Data{data, metadata, nil}
}

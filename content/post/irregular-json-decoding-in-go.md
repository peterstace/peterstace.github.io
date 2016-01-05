+++
Categories = ["Development", "GoLang"]
Description = ""
Tags = ["go", "golang", "json", "heterogeneous", "homogeneous", "marshal", "unmarshal", "regular", "irregular"]
date = "2015-06-18T22:03:45+11:00"
menu = "main"
title = "Irregular JSON decoding in Go"
aliases = ["/2015/06/18/irregular-json-decoding-in-go.html"]

+++

The Go standard library has an awesome [JSON encoding and decoding
package](http://golang.org/pkg/encoding/json/), which makes handling JSON a
breeze. If you're not familiar with the package, there are plenty of
[web](http://blog.golang.org/json-and-go) [pages](https://gobyexample.com/json)
[around](https://eager.io/blog/go-and-json/) that explain its basic usage.

Basically, if you know the structure of the JSON value you're encoding, you
create the zero value of the corresponding Go type, and pass a pointer to it
into the `json.Unmarshal` (along with the JSON value). Magic occurs, and your
Go object is now populated. If you don't know the structure of your JSON value
upfront, you can instead pass in a `map[string]interface{}`, and that will be
populated instead. Type assertions can then be used on the empty interfaces to
determine what actually got decoded.

But what if you know the precise structure of the JSON, but it's not _regular_?
For example, the following JSON value represents a ray tracer scene. The array
`"Objects"` contains a known structure, but each element isn't of the same
type. Objects of type `"box"` will always have `"Corner1"` and `"Corner2"`
fields, and objects of type `"sphere"` will always have `"Centre"` and
`"Radius"` fields.

    {
    	"Colour": "white",
    	"Material": "lambertian",
    	"Objects": [
    		{
    			"Type": "box",
    			"Corner1": {"X":0,"Y":1,"Z":2},
    			"Corner2": {"X":5,"Y":6,"Z":7}
    		},
    		{
    			"Type": "sphere",
    			"Centre": {"X":2,"Y":2.5,"Z":-1},
    			"Radius": 1.0
    		}
    	]
    }

I want to decode the JSON value into the follow Go data structure:

    type World struct {
    	Colour   string
    	Material string
    	Objects  []Object
    }
    
    type Object interface {
    	Contains(Vect) bool
    }
    
    type Box struct {
    	Corner1, Corner2 Vect
    }
    func (b Box) Contains(Vect) bool { ... }
    
    type Sphere struct {
    	Centre Vect
    	Radius float64
    }
    func (s Sphere) Contains(Vect) bool { ... }
    
    type Vect struct {
    	X, Y, Z float64
    }

Go can't unmarshal the JSON value into `World` directly. The `json.Unmarshal`
function will return an error "json: cannot unmarshal object into Go value of
type main.Object". This makes sense, since the JSON value and the `World` Go
type both have fields named `Objects`, but `Object` is a Go interface, so
cannot be unmarshalled into.

We need to perform custom unmarshalling into the `World` type by implementing
the [Unmarshaler](http://golang.org/pkg/encoding/json/#Unmarshaler) interface.

    func (w *World) UnmarshalJSON(p []byte) error {
    
    	// Unmarshal the regular parts of the JSON value.
    	record := struct {
    		Colour   string
    		Material string
    		Objects  []json.RawMessage
    	}{}
    	if err := json.Unmarshal(p, &record); err != nil {
    		return err
    	}
    	w.Colour = record.Colour
    	w.Material = record.Material
    
    	// Go back to the irregular parts, find each type and unmarshal accordingly.
    	for _, rawObject := range record.Objects {
    
    		// Find the type.
    		obj := struct{ Type string }{}
    		if err := json.Unmarshal(rawObject, &obj); err != nil {
    			return err
    		}
    
    		// Perform type specific unmarshalling.
    		switch obj.Type {
    		case "box":
    			var b Box
    			if err := json.Unmarshal(rawObject, &b); err != nil {
    				return err
    			}
    			w.Objects = append(w.Objects, b)
    		case "sphere":
    			var s Sphere
    			if err := json.Unmarshal(rawObject, &s); err != nil {
    				return err
    			}
    			w.Objects = append(w.Objects, s)
    		}
    	}
    	return nil
    }

So what's happening here? We are essentially doing the following:

1. Create a variable called `record` that allows us to decode the regular parts
   (`Colour` and `Material`). It also decodes the _irregular_ parts into
`json.RawMessage` objects.
2. Iterate over each `json.RawMessage`, and extract enough information to work
   out which type we should unmarshal into. In this case, it's easy, we just
look for the  `"Type"` field, and switch based on that.
3. Decode into the appropriate type as determined in step 2.

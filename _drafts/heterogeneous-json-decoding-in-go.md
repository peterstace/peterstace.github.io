---
layout: post
title: Heterogeneous JSON decoding in Go
tags: XXX
---

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

But what if you know the precise structure of the JSON, but it's not 'regular'?
For example, the following JSON value represents a ray tracer scene. The array
`"Objects"` contains a precise structure, but each element isn't of the same
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

I want to decode the JSON value into the follow data structure:

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
    
    type Sphere struct {
    	Centre Vect
    	Radius float64
    }
    
    type Vect struct {
    	X, Y, Z float64
    }

Go can't unmarshal the JSON value into `World` out of the box. The `Unmarshal`
error will return an error "json: cannot unmarshal object into Go value of type
main.Object".

The solution is to have `World` implement the
[Unmarshaler](http://golang.org/pkg/encoding/json/#Unmarshaler) interface. The
implementation should directly unmarshal and populate the `Colour` and
`Material` fields of `World`. Then for each element of the `"Objects"` array,
the appropriate unmarshalling can take place depending on if we have a box or a
sphere. This is essentially a two-pass unmarshal. First decode the type, then
decode the object itself now that we know the type. This can be achieved with
the [json.RawMessage](http://golang.org/pkg/encoding/json/#RawMessage) type,
which allows unmarshalling of a sub piece of the JSON value to be unmarshalled
later.

The implementation of the `Unmarshaler` interface is show below.

    func (w *World) UnmarshalJSON(p []byte) error {
    
    	// Unmarshal the homogeneous parts of the JSON value.
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
    
    	// Go back to the heterogeneous parts, find each type and unmarshal accordingly.
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

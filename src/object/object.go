package object

import "fmt"


type ObjectType string

type Object interface{
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ = "NULL"
	STRING_OBJ="STRING"
	ARRAY_OBJ="ARRAY"
)

type Integer struct{
	Value int64
}

func (i *Integer) Type() ObjectType{ return INTEGER_OBJ }
func (i *Integer) Inspect() string{ return fmt.Sprintf("%d", i.Value)}

type Boolean struct{
	Value bool
}

func (i *Boolean) Type() ObjectType{ return BOOLEAN_OBJ }
func (i *Boolean) Inspect() string{ return fmt.Sprintf("%d", i.Value)}

type Null struct{}
func (n *Null) Type() ObjectType{ return NULL_OBJ }
func (n *Null) Inspect() string{ return "null" }

type String struct{
	Value string
}

func (s *String) Type() ObjectType{ return STRING_OBJ}
func (s *String) Inspect() string {return s.Value}

type Array struct{
	Elements []Object
}
func (a *Array) Type() ObjectType{ return ARRAY_OBJ }
func (a *Array) Inspect() string{
	var out string
	for i, element := range a.Elements{
		out += element.Inspect()
		if i != len(a.Elements)-1{
			out += ", "
		}
	}
	return "[" + out + "]"
}
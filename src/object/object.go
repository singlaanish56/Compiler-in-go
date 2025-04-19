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
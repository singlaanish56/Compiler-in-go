package object

import "fmt"


type ObjectType string

type Object interface{
	Type() ObjectType
	Inspect() string
}

type Integer struct{
	Value int64
}

func (i *Integer) Type() ObjectType{ return "INTEGER" }
func (i *Integer) Inspect() string{ return fmt.Sprintf("%d", i.Value)}
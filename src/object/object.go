package object

import (
	"fmt"
	"hash/fnv"

	"github.com/singlaanish56/Compiler-in-go/code"
)


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
	HASHPAIR_OBJ="HASHPAIR"
	COMPILE_FUNCTION_OBJ="COMPILE_FUNCTION"
)

type HashKey struct{
	Type ObjectType
	Value uint64
}

type Hashable interface{
	HashKey() HashKey
}

type HashPair struct{
	Key Object
	Value Object
}

type Integer struct{
	Value int64
}

func (i *Integer) Type() ObjectType{ return INTEGER_OBJ }
func (i *Integer) Inspect() string{ return fmt.Sprintf("%d", i.Value)}
func (i *Integer) HashKey() HashKey{
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type Boolean struct{
	Value bool
}

func (b *Boolean) Type() ObjectType{ return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string{ return fmt.Sprintf("%d", b.Value)}
func (b *Boolean) HashKey() HashKey{
	if b.Value{
		return HashKey{Type: b.Type(), Value: uint64(1)}
	}

	return HashKey{Type: b.Type(), Value: uint64(0)}
}	

type Null struct{}
func (n *Null) Type() ObjectType{ return NULL_OBJ }
func (n *Null) Inspect() string{ return "null" }

type String struct{
	Value string
}

func (s *String) Type() ObjectType{ return STRING_OBJ}
func (s *String) Inspect() string {return s.Value}
func (s *String) HashKey() HashKey{
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

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

type Hash struct{
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType{ return HASHPAIR_OBJ }
func (h *Hash) Inspect() string{
	var out string
	for _, pair := range h.Pairs{
		out += pair.Key.Inspect() + ": " + pair.Value.Inspect()
	}
	return "{" + out + "}"
}

type CompiledFunction struct{
	Instructions code.Instructions
	NumberOfLocals int
}

func (cf *CompiledFunction) Type() ObjectType{ return COMPILE_FUNCTION_OBJ }
func (cf *CompiledFunction) Inspect() string{ return fmt.Sprintf("CompiledFunction[%p]",cf) }
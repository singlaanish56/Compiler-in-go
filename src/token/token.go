package token

type TokenType string

type Token struct{
	Type TokenType
	Identifier string
	StartPosition int
	EndPosition int
}

var KeywordMap = map[string]TokenType{
	"fn":FUNCTION,
	"let":LET,
	"true":TRUE,
	"false":FALSE,
	"if":IF,
	"else":ELSE,
	"return":RETURN,
	"null":NULL,
	"var":VARIABLE,
}


const(
	FUNCTION ="fn"
	LET = "let"
	IF="if"
	ELSE="else"
	RETURN="return"

	VARIABLE="var"
	STRING="str"
	NUMBER="int"
	TRUE="t"
	FALSE="f"
	NULL="nullptr"

	OPENROUND="("
	CLOSEROUND=")"
	OPENBRACE="{"
	CLOSEBRACE="}"
	OPENBRACKET="["
	CLOSEBRACKET="]"
	OPENANGLE="<"
	CLOSEANGLE=">"

	SEMICOLON=";"
	COLON=":"
	COMMA=","
	PLUS="+"
	MINUS="-"
	DIVIDE="/"
	MULTIPLY="*"
	EQUALTO="="
	UNDERSCORE="_"
	DOUBLEEQUALTO="=="
	EXCLAMATION="!"
	EXCLAMATIONEQUALTO="!="

	INVALID="inv"
	EOF="eof"
)
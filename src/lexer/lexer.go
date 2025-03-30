package lexer

import (

	"github.com/singlaanish56/Compiler-in-go/token"
)


type Lexer struct {
	input []rune
	char rune
	currentPosition int
	nextReadPosition int
}

func New(input string) *Lexer{
	lexer := &Lexer{input: []rune(input), nextReadPosition: 0}
	lexer.nextChar()
	return lexer
}

func (l *Lexer) NextToken() token.Token{
 
	for isEscapeSequence(l.char){
		l.nextChar()
	}

	// is it a number
	if l.char>= '0' && l.char<='9' {
		return l.retrieveTheNumber()
	}

	// is it string
	if l.char == '"'{
		return l.retrieveTheString()
	}

	// any variables
	if((l.char>='a' && l.char<='z') || (l.char>='A' && l.char<='Z')){
		return l.retrieveTheVariable()
	}

	tk := l.retrieveTheSign()
	l.nextChar()

	return  tk
}

func (l *Lexer) nextChar(){
	if l.nextReadPosition>= len(l.input){
		l.char = 0
	}else{
		l.char =l.input[l.nextReadPosition]
	}

	l.currentPosition = l.nextReadPosition
	l.nextReadPosition++
}

func (l *Lexer) peekChar() rune{
	if l.currentPosition >= len(l.input){
		return 0
	}

	return l.input[l.nextReadPosition]
}

func (l *Lexer) retrieveTheNumber() token.Token{
	start := l.currentPosition

	for l.char >= '0' && l.char <= '9'{
		l.nextChar()
	}

	return token.Token{Type: token.NUMBER, Identifier : string(l.input[start:l.currentPosition]), StartPosition: start, EndPosition: l.currentPosition}
}

func (l *Lexer) retrieveTheString() token.Token{
	start := l.currentPosition+1
	var strBuilder []rune
	for {
		l.nextChar()

		if l.char=='"' || l.char==0{
			break
		}

		if l.char=='\n' || l.char=='\t' || l.char=='\r'{
			continue
		}

		strBuilder = append(strBuilder, l.char)
	}

	str := string(strBuilder)
	endIndex := l.currentPosition
	l.nextChar()
	return token.Token{Type: token.STRING, Identifier: str, StartPosition: start, EndPosition: endIndex}
}

func (l *Lexer) retrieveTheVariable() token.Token{
	start := l.currentPosition

	for (l.char>='a' && l.char<='z') || (l.char>='A' && l.char<='Z'){
		l.nextChar()
	}

	str := string(l.input[start:l.currentPosition])
	tt:= token.KeywordMap["var"]
	if tokenType, exists := token.KeywordMap[str]; exists{
		tt = tokenType
	}

	return token.Token{Type: tt, Identifier: str, StartPosition: start, EndPosition: l.currentPosition}	
}

func (l *Lexer) retrieveTheSign() token.Token{
	var tk token.Token
	switch l.char{
	case '=':	
		str:= string(l.char)
		if l.peekChar() == '='{
			l.nextChar()
			str = str + string(l.char)
			tk = token.Token{Type: token.DOUBLEEQUALTO, Identifier: str, StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
		}else{
			tk = token.Token{Type: token.EQUALTO, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
		}
	case '_':
		tk = token.Token{Type: token.UNDERSCORE, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
	case ';':
		tk = token.Token{Type: token.SEMICOLON, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
	case ':':
		tk = token.Token{Type: token.COLON, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
	case '(':
		tk = token.Token{Type: token.OPENROUND, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
	case ')':
		tk = token.Token{Type: token.CLOSEROUND, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
	case '{':
		tk = token.Token{Type: token.OPENBRACE, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
	case '}':
		tk = token.Token{Type: token.CLOSEBRACE, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
	case '[':
		tk = token.Token{Type: token.OPENBRACKET, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
	case ']':
		tk = token.Token{Type: token.CLOSEBRACKET, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
	case '<':
		tk = token.Token{Type: token.OPENANGLE, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
	case '>':
		tk = token.Token{Type: token.CLOSEANGLE, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
	case ',':
		tk = token.Token{Type: token.COMMA, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
	case '+':
		tk = token.Token{Type: token.PLUS, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
	case '-':
		tk = token.Token{Type: token.MINUS, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
	case '/':
		tk = token.Token{Type: token.DIVIDE, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
	case '*':
		tk = token.Token{Type: token.MULTIPLY, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
	case '!':
		str := string(l.char)
		if l.peekChar() == '='{
			l.nextChar()
			str = str + string(l.char)
			tk = token.Token{Type: token.EXCLAMATIONEQUALTO, Identifier: str, StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
		}else{
			tk = token.Token{Type: token.EXCLAMATION, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
		}
	default:
		if l.char == 0{
			tk = token.Token{Type: token.EOF, Identifier: "", StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
		}else{
			tk = token.Token{Type: token.INVALID, Identifier: string(l.char), StartPosition: l.currentPosition, EndPosition: l.nextReadPosition}
		}		
	}

	return tk
}

func isEscapeSequence(c rune) bool{
	return  c==' ' || c=='\n' || c=='\t' || c=='\r'
}
package lexer

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var reNewLine = regexp.MustCompile("\r\n|\n\r|\n|\r")
var reIdentifier = regexp.MustCompile(`^[_\d\w]+`)
var reNumber = regexp.MustCompile(`^0[xX][0-9a-fA-F]*(\.[0-9a-fA-F]*)?([pP][+\-]?[0-9]+)?|^[0-9]*(\.[0-9]*)?([eE][+\-]?[0-9]+)?`)
var reShortStr = regexp.MustCompile(`(?s)(^'(\\\\|\\'|\\\n|\\z\s*|[^'\n])*')|(^"(\\\\|\\"|\\\n|\\z\s*|[^"\n])*")`)
var reOpeningLongBracket = regexp.MustCompile(`^\[=*\[`)

var reDecEscapeSeq = regexp.MustCompile(`^\\[0-9]{1,3}`)
var reHexEscapeSeq = regexp.MustCompile(`^\\x[0-9a-fA-F]{2}`)
var reUnicodeEscapeSeq = regexp.MustCompile(`^\\u\{[0-9a-fA-F]+\}`)

type Lexer struct {
	chunk         string //源代码
	chunkName     string //源文件名
	line          int    //当前行号
	nextToken     string //下一个token
	nextTokenKind int    //下一个token的编译结果
	nextTokenLine int    //下一个token的行号
}

//创建行号为1的Lexer
func NewLexer(chunk, chunkName string) *Lexer {
	return &Lexer{chunk, chunkName, 1, "", 0, 0}
}

func (self *Lexer) Line() int {
	return self.line
}

//判断字符串的开头
func (self *Lexer) test(s string) bool {
	return strings.HasPrefix(self.chunk, s)
}

//跳过n个字符
func (self *Lexer) next(n int) {
	self.chunk = self.chunk[n:]
}

//判断是否是空白字符
func isWhiteSpace(c byte) bool {
	switch c {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		return true
	}
	return false
}

//是回车或者换行
func isNewLine(c byte) bool {
	return c == '\r' || c == '\n'
}

func (self *Lexer) error(f string, a ...interface{}) {
	err := fmt.Sprintf(f, a...)
	err = fmt.Sprintf("%s:%d: %s", self.chunkName, self.line, err)
	panic(err)
}

//检测是数字
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

//判断是字母
func isLetter(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

//跳过一些换行和注释
func (self *Lexer) skipWhiteSpaces() {
	for len(self.chunk) > 0 {
		if self.test("--") {
			self.skipComment()
		} else if self.test("\r\n") || self.test("\n\r") {
			self.next(2)
			self.line += 1
		} else if isNewLine(self.chunk[0]) {
			self.next(1)
			self.line += 1
		} else if isWhiteSpace(self.chunk[0]) {
			self.next(1)
		} else {
			break
		}

	}
}

//去除长注释
func (self *Lexer) scanLongString() string {
	//找出[=?[长注释的头
	openingLongBracket := reOpeningLongBracket.FindString(self.chunk)
	if openingLongBracket == "" {
		self.error("invalid long string delimiter near '%s'", self.chunk[0:2])
	}

	//头尾数量应该对齐 所以替换
	closingLongBracket := strings.Replace(openingLongBracket, "[", "]", -1)
	//找出尾
	closingLongBracketIdx := strings.Index(self.chunk, closingLongBracket)
	if closingLongBracketIdx < 0 {
		self.error("unfinished long string or comment")
	}

	str := self.chunk[len(openingLongBracket):closingLongBracketIdx]
	//跳过尾开始index + 尾长度
	self.next(closingLongBracketIdx + len(closingLongBracket))

	str = reNewLine.ReplaceAllString(str, "\n")
	self.line += strings.Count(str, "\n")
	//跳过第一个换行符
	if len(str) > 0 && str[0] == '\n' {
		str = str[1:]
	}
	return str
}

//扫描短字符串
func (self *Lexer) scanShortString() string {
	if str := reShortStr.FindString(self.chunk); str != "" {
		self.next(len(str))
		str = str[1 : len(str)-1]
		if strings.Index(str, `\`) >= 0 {
			self.line += len(reNewLine.FindAllString(str, -1))
			str = self.escape(str)
		}
		return str
	}
	self.error("unfinished string")
	return ""
}

//转义
func (self *Lexer) escape(str string) string {
	var buf bytes.Buffer

	for len(str) > 0 {
		if str[0] != '\\' {
			buf.WriteByte(str[0])
			str = str[1:]
			continue
		}
		if len(str) == 1 {
			self.error("unfinished string")
		}
		switch str[1] {
		case 'a':
			buf.WriteByte('\a')
			str = str[2:]
			continue
		case 'b':
			buf.WriteByte('\b')
			str = str[2:]
			continue
		case 'f':
			buf.WriteByte('\f')
			str = str[2:]
			continue
		case 'n', '\n':
			buf.WriteByte('\n')
			str = str[2:]
			continue
		case 'r':
			buf.WriteByte('\r')
			str = str[2:]
			continue
		case 't':
			buf.WriteByte('\t')
			str = str[2:]
			continue
		case 'v':
			buf.WriteByte('\v')
			str = str[2:]
			continue
		case '"':
			buf.WriteByte('"')
			str = str[2:]
			continue
		case '\'':
			buf.WriteByte('\'')
			str = str[2:]
			continue
		case '\\':
			buf.WriteByte('\\')
			str = str[2:]
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // \ddd
			//超过0xFF则报错
			if found := reDecEscapeSeq.FindString(str); found != "" {
				d, _ := strconv.ParseInt(found[1:], 10, 32)
				if d <= 0xFF {
					buf.WriteByte(byte(d))
					str = str[len(found):]
					continue
				}
				self.error("decimal escape too large near '%s'", found)
			}
		case 'x': // \xXX
			//16进制转换 \x0-9a-f
			if found := reHexEscapeSeq.FindString(str); found != "" {
				d, _ := strconv.ParseInt(found[2:], 16, 32)
				buf.WriteByte(byte(d))
				str = str[len(found):]
				continue
			}
		case 'u': // \u{XXX}
			//unicode代码 <=十六进制
			if found := reUnicodeEscapeSeq.FindString(str); found != "" {
				d, err := strconv.ParseInt(found[3:len(found)-1], 16, 32)
				if err == nil && d <= 0x10FFFF {
					buf.WriteRune(rune(d))
					str = str[len(found):]
					continue
				}
				self.error("UTF-8 value too large near '%s'", found)
			}
		case 'z':
			//先跳过这个转义 和空白字符串  正则表达式用
			str = str[2:]
			for len(str) > 0 && isWhiteSpace(str[0]) { // todo
				str = str[1:]
			}
			continue
		}
		self.error("invalid escape sequence near '\\%c'", str[1])
	}

	return buf.String()
}

//跳过注释
func (self *Lexer) skipComment() {
	self.next(2)        // skip --
	if self.test("[") { //long comment ?
		if reOpeningLongBracket.FindString(self.chunk) != "" {
			self.scanLongString()
			return
		}
	}
	for len(self.chunk) > 0 && !isNewLine(self.chunk[0]) {
		self.next(1)
	}
}

func (self *Lexer) scan(re *regexp.Regexp) string {
	if token := re.FindString(self.chunk); token != "" {
		self.next(len(token))
		return token
	}
	panic("unreachable!")
}

//判断是数字
func (self *Lexer) scanNumber() string {
	return self.scan(reNumber)
}

//判断是关键字or变量
func (self *Lexer) scanIdentifier() string {
	return self.scan(reIdentifier)
}

//预测下一个token
func (self *Lexer) LookAhead() int {
	if self.nextTokenLine > 0 {
		return self.nextTokenKind
	}

	currentLine := self.line
	line, kind, token := self.NextToken()
	self.line = currentLine
	self.nextTokenLine = line
	self.nextTokenKind = kind
	self.nextToken = token
	return kind
}

//提取指定的kind
func (self *Lexer) NextTokenOfKind(kind int) (line int, token string) {
	line, _kind, token := self.NextToken()
	if kind != _kind {
		self.error("syntax error near '%s'", token)
	}
	return line, token
}

//提取 自己定义的变量
func (self *Lexer) NextIdentifier() (line int, token string) {
	return self.NextTokenOfKind(TOKEN_IDENTIFIER)
}

func (self *Lexer) NextToken() (line, kind int, token string) {
	//已经预测好了
	if self.nextTokenLine > 0 {
		line = self.nextTokenLine
		kind = self.nextTokenKind
		token = self.nextToken
		self.line = self.nextTokenLine
		self.nextTokenLine = 0
		return
	}

	self.skipWhiteSpaces()
	if len(self.chunk) == 0 {
		return self.line, TOKEN_EOF, "EOF"
	}
	switch self.chunk[0] {
	case ';':
		self.next(1)
		return self.line, TOKEN_SEP_SEMI, ";"
	case ',':
		self.next(1)
		return self.line, TOKEN_SEP_COMMA, ","
	case '(':
		self.next(1)
		return self.line, TOKEN_SEP_LPAREN, "("
	case ')':
		self.next(1)
		return self.line, TOKEN_SEP_RPAREN, ")"
	case ']':
		self.next(1)
		return self.line, TOKEN_SEP_RBRACK, "]"
	case '{':
		self.next(1)
		return self.line, TOKEN_SEP_LCURLY, "{"
	case '}':
		self.next(1)
		return self.line, TOKEN_SEP_RCURLY, "}"
	case '+':
		self.next(1)
		return self.line, TOKEN_OP_ADD, "+"
	case '-':
		self.next(1)
		return self.line, TOKEN_OP_MINUS, "-"
	case '*':
		self.next(1)
		return self.line, TOKEN_OP_MUL, "*"
	case '^':
		self.next(1)
		return self.line, TOKEN_OP_POW, "^"
	case '%':
		self.next(1)
		return self.line, TOKEN_OP_MOD, "%"
	case '&':
		self.next(1)
		return self.line, TOKEN_OP_BAND, "&"
	case '|':
		self.next(1)
		return self.line, TOKEN_OP_BOR, "|"
	case '#':
		self.next(1)
		return self.line, TOKEN_OP_LEN, "#"
	case ':':
		if self.test("::") {
			self.next(2)
			return self.line, TOKEN_SEP_LABEL, "::"
		} else {
			self.next(1)
			return self.line, TOKEN_SEP_COLON, ":"
		}
	case '/':
		if self.test("//") {
			self.next(2)
			return self.line, TOKEN_OP_IDIV, "//"
		} else {
			self.next(1)
			return self.line, TOKEN_OP_DIV, "/"
		}
	case '~':
		if self.test("~=") {
			self.next(2)
			return self.line, TOKEN_OP_NE, "~="
		} else {
			self.next(1)
			return self.line, TOKEN_OP_WAVE, "~"
		}
	case '=':
		if self.test("==") {
			self.next(2)
			return self.line, TOKEN_OP_EQ, "=="
		} else {
			self.next(1)
			return self.line, TOKEN_OP_ASSIGN, "="
		}
	case '<':
		if self.test("<<") {
			self.next(2)
			return self.line, TOKEN_OP_SHL, "<<"
		} else if self.test("<=") {
			self.next(2)
			return self.line, TOKEN_OP_LE, "<="
		} else {
			self.next(1)
			return self.line, TOKEN_OP_LT, "<"
		}
	case '>':
		if self.test(">>") {
			self.next(2)
			return self.line, TOKEN_OP_SHR, ">>"
		} else if self.test(">=") {
			self.next(2)
			return self.line, TOKEN_OP_GE, ">="
		} else {
			self.next(1)
			return self.line, TOKEN_OP_GT, ">"
		}
	case '.':
		if self.test("...") {
			self.next(3)
			return self.line, TOKEN_VARARG, "..."
		} else if self.test("..") {
			self.next(2)
			return self.line, TOKEN_OP_CONCAT, ".."
		} else if len(self.chunk) == 1 || !isDigit(self.chunk[1]) {
			self.next(1)
			return self.line, TOKEN_SEP_DOT, "."
		}
	case '[':
		if self.test("[[") || self.test("[=") {
			return self.line, TOKEN_STRING, self.scanLongString()
		} else {
			self.next(1)
			return self.line, TOKEN_SEP_LBRACK, "["
		}
	case '\'', '"':
		return self.line, TOKEN_STRING, self.scanShortString()
	}

	//数字扫描
	c := self.chunk[0]
	if c == '.' || isDigit(c) {
		token := self.scanNumber()
		return self.line, TOKEN_NUMBER, token
	}

	//关键字和变量扫描
	if c == '_' || isLetter(c) {
		token := self.scanIdentifier()
		if kind, found := keywords[token]; found {
			return self.line, kind, token //keyword
		} else {
			return self.line, TOKEN_IDENTIFIER, token
		}
	}

	//return default or null
	self.error("unexpected symbol near %q", c)
	return
}

package comment

type StateCode int

const (
	text StateCode = iota
	leftSlash
	leftStar
	comment
	rightSlash
	rightStar
	slashSlash
	commentLine
)

type Pos struct {
	FP   int
	Line int
	Word int
}

type Context struct {
	LastState StateCode
	Src       *[]byte
	Len       int
	TokenPos  Pos
	CurPos    Pos
}

type Token struct {
	Content  string
	Position Pos
}

func GetTokenFromBytes(b []byte) []Token {
	initialPos := Pos{
		FP:   0,
		Line: 1,
		Word: 1,
	}

	c := Context{
		LastState: text,
		Src:       &b,
		Len:       len(b),
		TokenPos:  initialPos,
		CurPos:    initialPos,
	}

	var res []Token
	for {
		t, stop := c.scan()
		if t.Content != "" {
			res = append(res, t)
		}
		if stop {
			break
		}
	}
	return res
}

func (c *Context) scan() (content Token, end bool) {

	if !c.hasCur() { //end of file
		if c.LastState == rightSlash {
			return Token{
				Content:  string((*c.Src)[c.TokenPos.FP : c.CurPos.FP-2]),
				Position: c.TokenPos,
			}, true
		} else if c.LastState == commentLine {
			return Token{
				Content:  string((*c.Src)[c.TokenPos.FP:c.CurPos.FP]),
				Position: c.TokenPos,
			}, true
		} else {
			return Token{}, true
		}
	}

	switch c.LastState {

	case text:
		c.preReadText()
	case leftSlash:
		if c.curWord() == '*' {
			c.next()
			c.LastState = leftStar
		} else if c.curWord() == '/' {
			c.next()
			c.LastState = slashSlash
		} else {
			c.next()
			c.LastState = text
		}
	case leftStar:
		if c.curWord() == '*' {
			c.next()
			c.LastState = rightStar
		} else {
			c.TokenPos = c.CurPos
			c.next()
			c.LastState = comment
		}
	case comment:
		c.preReadComment()
	case rightStar:
		if c.curWord() == '/' {
			c.next()
			c.LastState = rightSlash
		} else if c.curWord() == '*' {
			c.next()
		} else {
			c.next()
			c.LastState = comment
		}
	case rightSlash:
		if c.curWord() == '/' {
			c.next()
			c.LastState = leftSlash
		} else {
			c.next()
			c.LastState = text
		}
		return Token{
			Content:  string((*c.Src)[c.TokenPos.FP : c.CurPos.FP-3]),
			Position: c.TokenPos,
		}, false
	case slashSlash:
		if c.curWord() == '\n' {
			c.next()
			c.LastState = text
		} else {
			c.TokenPos = c.CurPos
			c.next()
			c.LastState = commentLine
		}
	case commentLine:
		token, got := c.preReadCommentLine()
		if got {
			return token, false
		}
	}
	return Token{
		Content:  "",
		Position: c.TokenPos,
	}, false
}

func (c *Context) preReadText() {
	for c.hasCur() && c.curWord() != '/' {
		c.next()
	}
	if c.hasCur() && c.curWord() == '/' {
		c.next()
		c.LastState = leftSlash
	}
}

func (c *Context) preReadComment() {
	for c.hasCur() && c.curWord() != '*' {
		c.next()
	}
	if c.hasCur() && c.curWord() == '*' {
		c.next()
		c.LastState = rightStar
	} //else wait for next scan
}

func (c *Context) preReadCommentLine() (Token, bool) {
	for c.hasCur() && c.curWord() != '\n' {
		c.next()
	}
	if c.hasCur() && c.curWord() == '\n' {
		c.next()
		c.LastState = text
		return Token{
			Content:  string((*c.Src)[c.TokenPos.FP : c.CurPos.FP-1]),
			Position: c.TokenPos,
		}, true
	} //else wait for next scan
	return Token{}, false
}

func (c *Context) hasCur() bool {
	if c.CurPos.FP >= c.Len {
		return false
	}
	return true
}

func (c *Context) curWord() byte {
	return (*c.Src)[c.CurPos.FP]
}

func (c *Context) nextWord() byte {
	return (*c.Src)[c.CurPos.FP+1]
}

func (c *Context) next() {
	if (*c.Src)[c.CurPos.FP] == '\n' {
		c.CurPos.Line++
		c.CurPos.Word = 1
	} else {
		c.CurPos.Word++
	}
	c.CurPos.FP++
}

package comment

type stateCode int

const (
	text stateCode = iota
	slashStar
	comment
	slashSlash
	commentLine
)

type Token struct {
	Content  string
	Position pos
}

func GetTokenFromBytes(b []byte) []Token {
	initialPos := pos{
		FP:   0,
		Line: 1,
		Word: 1,
	}

	c := context{
		LastState: text,
		Src:       &b,
		Len:       len(b),
		TokenPos:  initialPos,
		CurPos:    initialPos,
	}

	for {
		if stop := c.scan(); stop {
			break
		}
	}
	return c.Results
}

type pos struct {
	FP   int
	Line int
	Word int
}

type context struct {
	LastState stateCode
	Src       *[]byte
	Len       int
	TokenPos  pos
	CurPos    pos
	Results   []Token
}

func (c *context) scan() (end bool) {

	//c.hasCur() => true
	switch c.LastState {

	case text:
		c.readText()
	case slashStar:
		// /*_?
		c.TokenPos = c.CurPos
		if c.curWord() == '*' {
			if !c.preReadStarSlash() {
				c.LastState = comment
			}
		} else {
			c.LastState = comment
		}
	case comment:
		c.readComment()
	case slashSlash:
		// //_?
		if c.curWord() == '\n' {
			if c.hasNext() {
				c.next()
			}
			c.LastState = text
		} else {
			c.TokenPos = c.CurPos
			c.LastState = commentLine
		}
	case commentLine:
		// //_?
		got := c.readCommentLine()
		if got { //c.curWord() => \n
			c.LastState = text
			c.Results = append(c.Results, Token{
				Content:  string((*c.Src)[c.TokenPos.FP:c.CurPos.FP]),
				Position: c.TokenPos,
			})
		} //else eof
	} // switch

	if !c.hasNext() { //end of file
		if c.LastState == commentLine {
			c.Results = append(c.Results, Token{
				Content:  string((*c.Src)[c.TokenPos.FP : c.CurPos.FP+1]),
				Position: c.TokenPos,
			})
			return true
		} else {
			return true
		}
	}
	return false
}

// // or /*
func (c *context) preReadSlash() stateCode {
	//c.curWord == '/'
	if !c.hasNext() {
		return -1
	}
	if c.nextWord() == '*' {
		return slashStar
	} else if c.nextWord() == '/' {
		return slashSlash
	}
	return -1
}

// */
func (c *context) preReadStarSlash() bool {
	//c.curWord == '*'
	if !c.hasNext() {
		return false
	}
	if c.nextWord() == '/' {
		c.Results = append(c.Results, Token{
			Content:  string((*c.Src)[c.TokenPos.FP:c.CurPos.FP]),
			Position: c.TokenPos,
		})
		c.LastState = text
		c.next()
		if c.hasNext() {
			c.next()
		}
		return true
	}
	return false
}

func (c *context) readText() {
	//c.LastState = text
	for {
		if !c.hasNext() {
			break
		}
		if c.curWord() == '/' {
			if code := c.preReadSlash(); code != -1 {
				c.LastState = code
				c.next()
				if c.hasNext() {
					c.next()
				}
				break
			}
		}
		c.next()
	}
}

func (c *context) readComment() {
	//c.LastState = comment
	for c.hasNext() {
		if c.curWord() == '*' && c.preReadStarSlash() {
			break
		}
		c.next()
	}
}

func (c *context) readCommentLine() (got bool) {
	for {
		if c.curWord() == '\n' {
			return true
		}
		if c.hasNext() {
			c.next()
		} else {
			break
		}
	}
	return false
}

func (c *context) preWord() byte {
	return (*c.Src)[c.CurPos.FP-1]
}

func (c *context) curWord() byte {
	return (*c.Src)[c.CurPos.FP]
}

func (c *context) nextWord() byte {
	return (*c.Src)[c.CurPos.FP+1]
}

func (c *context) hasNext() bool {
	if c.CurPos.FP+1 >= c.Len {
		return false
	}
	return true
}

func (c *context) next() {
	if (*c.Src)[c.CurPos.FP] == '\n' {
		c.CurPos.Line++
		c.CurPos.Word = 1
	} else {
		c.CurPos.Word++
	}
	c.CurPos.FP++
}

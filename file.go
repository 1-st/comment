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

type Context struct {
	LastState       StateCode
	Src        *[]byte
	Cur        int
	Len        int
	TokenStart int
}

func GetCommentStrings(b []byte) []string {
	s := Context{
		LastState:       text,
		Src:        &b,
		Cur:        0,
		Len:        len(b),
		TokenStart: 0,
	}
	var res []string
	for {
		s, stop := s.scan()

		if s != "" {
			res = append(res, s)
		}

		if stop {
			break
		}
	}
	return res
}

func (s *Context) scan() (content string, end bool) {

	if !s.hasCur() { //end of file
		if s.LastState == rightSlash {
			return string((*s.Src)[s.TokenStart : s.Cur-2]), true
		} else if s.LastState == commentLine {
			return string((*s.Src)[s.TokenStart:s.Cur]), true
		} else {
			return "", true
		}
	}


	switch s.LastState {

	case text:
		if s.curWord() == '/' {
			s.Cur++
			s.LastState = leftSlash
		} else {
			s.Cur++
		}
	case leftSlash:
		if s.curWord() == '*' {
			s.Cur++
			s.LastState = leftStar
		} else if s.curWord() == '/' {
			s.Cur++
			s.LastState = slashSlash
		} else {
			s.Cur++
			s.LastState = text
		}
	case leftStar:
		if s.curWord() == '*' {
			s.Cur++
			s.LastState = rightStar
		} else {
			s.TokenStart = s.Cur
			s.Cur++
			s.LastState = comment
		}
	case comment:
		if s.curWord() == '*' {
			s.Cur++
			s.LastState = rightStar
		} else {
			s.Cur++
		}
	case rightStar:
		if s.curWord() == '/' {
			s.Cur++
			s.LastState = rightSlash
		} else if s.curWord() == '*' {
			s.Cur++
		} else {
			s.Cur++
			s.LastState = comment
		}
	case rightSlash:
		if s.curWord() == '/' {
			s.Cur++
			s.LastState = leftSlash
		} else {
			s.Cur++
			s.LastState = text
		}
		return string((*s.Src)[s.TokenStart : s.Cur-3]), false
	case slashSlash:
		if s.curWord() == '\n' {
			s.Cur++
			s.LastState = text
		} else {
			s.TokenStart = s.Cur
			s.Cur++
			s.LastState = commentLine
		}
	case commentLine:
		if s.curWord() == '\n' {
			s.Cur++
			s.LastState = text
			return string((*s.Src)[s.TokenStart:s.Cur-1]), false
		} else {
			s.Cur++
		}
	}
	return "", false
}

func (s *Context) hasCur() bool {
	if s.Cur >= s.Len {
		return false
	}
	return true
}

func (s *Context) curWord() byte {
	return (*s.Src)[s.Cur]
}

func (s *Context) nextWord() byte {
	return (*s.Src)[s.Cur+1]
}

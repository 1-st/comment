package comment

type StateCode int

const (
	text StateCode = iota
	leftSlash
	leftStar
	comment
	rightSlash
	rightStar
	slashSlash //
	commentLine
)

type Context struct {
	Code       StateCode
	Src        *[]byte
	Cur        int
	Len        int
	TokenStart int
}

func GetCommentStrings(b []byte) []string {
	s := Context{
		Code:       text,
		Src:        &b,
		Cur:        0,
		Len:        len(b),
		TokenStart: 0,
	}
	var res []string
	for {
		s, stop := s.next()

		if s != "" {
			res = append(res, s)
		}

		if stop {
			break
		}
	}
	return res
}

func (s *Context) next() (content string, end bool) {
	if !s.hasCur() {
		if s.Code == rightSlash {
			return string((*s.Src)[s.TokenStart : s.Cur-2]), true
		} else if s.Code == commentLine {
			return string((*s.Src)[s.TokenStart:s.Cur]), true
		} else {
			return "", true
		}
	}
	switch s.Code {
	case text:
		if s.curWord() == '/' {
			s.Cur++
			s.Code = leftSlash
		} else {
			s.Cur++
		}
	case leftSlash:
		if s.curWord() == '*' {
			s.Cur++
			s.Code = leftStar
		} else if s.curWord() == '/' {
			s.Cur++
			s.Code = slashSlash
		} else {
			s.Cur++
			s.Code = text
		}
	case leftStar:
		if s.curWord() == '*' {
			s.Cur++
			s.Code = rightStar
		} else {
			s.TokenStart = s.Cur
			s.Cur++
			s.Code = comment
		}
	case comment:
		if s.curWord() == '*' {
			s.Cur++
			s.Code = rightStar
		} else {
			s.Cur++
		}
	case rightStar:
		if s.curWord() == '/' {
			s.Cur++
			s.Code = rightSlash
		} else if s.curWord() == '*' {
			s.Cur++
		} else {
			s.Cur++
			s.Code = comment
		}
	case rightSlash:
		if s.curWord() == '/' {
			s.Cur++
			s.Code = leftSlash
		} else {
			s.Cur++
			s.Code = text
		}
		return string((*s.Src)[s.TokenStart : s.Cur-3]), false
	case slashSlash:
		if s.curWord() == '\n' {
			s.Cur++
			s.Code = text
		} else {
			s.TokenStart = s.Cur
			s.Cur++
			s.Code = commentLine
		}
	case commentLine:
		if s.curWord() == '\n' {
			s.Cur++
			s.Code = text
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
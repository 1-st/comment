package comment

type StateCode int

const (
	start StateCode = iota
	beforeComment
	leftSlash
	leftStar
	comment
	rightStar
	rightSlash
	afterComment
	stop
)

type State struct {
	Code StateCode
	Src  *[]byte
	Cur  int
	Len  int
}

func keyWord(b byte) bool {
	if b == '/' || b == '*' {
		return true
	}
	return false
}

func (s *State) next() (content *string, end bool) {
	switch s.Code {
	case start:
		if !s.hasNext() {
			return nil, true
		}
		b := s.nextByte()
		if keyWord(b) {
			switch b {
			case '/':
				s.Code = leftSlash
				return nil,false
			case '*':

			}
		}

	case beforeComment:

	case leftSlash:

	case leftStar:

	case comment:

	case rightStar:

	case rightSlash:

	case afterComment:

	case stop:
		return nil, true
	}
	return nil, true
}

func (s *State) hasNext() bool {
	if s.Cur+1 == s.Len {
		return false
	}
	return true
}

func (s *State) nextByte() byte {
	s.Cur++
	return (*s.Src)[s.Cur]
}

func GetCommentStrings(b []byte) []string {
	s := State{
		Code: start,
		Src:  &b,
		Cur:  0,
		Len:  len(b),
	}
	var res []string
	for {
		s, stop := s.next()
		if stop {
			break
		}
		res = append(res, s)
	}
	return res
}

func getState(b byte) int {

}

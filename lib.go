package comment

import (
	"bytes"
	"strconv"
)

func PurifyToken(t *Token,filePath string)string{
	var res bytes.Buffer
	res.WriteString("{")
	res.WriteString(t.Content)
	res.WriteString("}\n")
	if filePath!=""{
		res.WriteString("  "+filePath+":"+strconv.Itoa(t.Position.Line)+":"+strconv.Itoa(t.Position.Word)+"  \n")
	}
	return res.String()
}

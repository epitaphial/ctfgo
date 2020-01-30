package tools

import(
	"fmt"
	"io"
    "crypto/md5"
)

func Md5Encode(str string) string {
    w := md5.New()
    io.WriteString(w, str)
    md5str := fmt.Sprintf("%x", w.Sum(nil))
    return md5str
}
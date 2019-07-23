package main

import (
	"crypto/md5"
	_ "crypto/md5"
	_ "crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log/log"
)

func genInputValue(value string, vposition string) string {
	tmplen := 64 - len(value)
	tmpstr := make([]rune, tmplen)
	for i := range tmpstr {
		tmpstr[i] = '0'
	}
	if vposition == "front" {
		return value + string(tmpstr)
	} else if vposition == "back" {
		return string(tmpstr) + value
	}
	return ""
}

func main() {
	files, err := ioutil.ReadDir("./files")
	if err != nil {
		log.Info(err)
	}
	for _, f := range files {

		Md5Inst := md5.New()
		Md5Inst.Write([]byte(f.Name()))
		result := Md5Inst.Sum([]byte(""))
		filename := hex.EncodeToString([]byte(f.Name()))
		filename = genInputValue(filename, "front")
		fmt.Printf("%s\n", filename)
		fmt.Printf(f.Name()+"\n%x\n\n", result)
	}
}

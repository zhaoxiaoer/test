package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	// RFC4648标准化了两种字符集，默认字符集用于MIME(RFC2025)和PEM(RFC1421)编码
	// 另一种用于URL和文件名编码，用'-'和'_'替换了'+'和'/'
	fmt.Printf("使用标准base64字符集:\n")
	hello := "hello, world"
	enc := base64.StdEncoding.EncodeToString([]byte(hello))
	fmt.Printf("enc: %s\n", enc)
	dec, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("dec: %s\n", string(dec))

	fmt.Printf("\n使用自定义base64字符集:\n")
	// 可以自定义base64字符集，用自定义的base64字符集编解码称为base64加解密
	base64Table := "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
	coder := base64.NewEncoding(base64Table)
	// 设置padding字符
	//	coder := base64.NewEncoding(base64Table).WithPadding(base64.NoPadding)
	hello = "hello, world"
	enc = coder.EncodeToString([]byte(hello))
	fmt.Printf("enc: %s\n", enc)
	dec, err = coder.DecodeString(enc)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("dec: %s\n", string(dec))
}

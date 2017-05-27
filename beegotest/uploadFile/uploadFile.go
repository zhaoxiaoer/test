package uploadFile

import (
	"io"
	"os"

	"github.com/astaxie/beego"
)

type UploadFile struct {
	beego.Controller
}

func (uf *UploadFile) Get() {
	uf.TplNames = "./uploadFile/uploadFile.tpl"
}

func (uf *UploadFile) Post() {
	uf.TplNames = "./uploadFile/uploadFile.tpl"

	file, fileHeader, err := uf.GetFile("fileUpload")
	if err != nil {
		uf.Data["result"] = "上传失败"
		return
	}
	uf.Data["result"] = "上传成功"
	defer file.Close()

	//创建文件
	f, err := os.OpenFile("./uploadFile/"+fileHeader.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		uf.Data["result"] = "文件创建失败"
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		uf.Data["result"] = "文件保存失败"
		return
	}
}

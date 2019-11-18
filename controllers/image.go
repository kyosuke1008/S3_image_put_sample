package controllers

import (
	"S3_image_put_sample/models"
	"github.com/astaxie/beego"
)

// Operations about Images
type S3ImageController struct {
	beego.Controller
}

// @Title 画像登録
// @Description create Images
// @Param	id		path 	string	true		"imageID"
// @Success 200 {string} image path
// @Failure 403 body is empty
// @router /:id [put]
func (p *S3ImageController) PutImage() {
	id := p.Ctx.Input.Param(":id")
	err := p.Ctx.Input.ParseFormOrMulitForm(32 << 20)
	var picture, _, _ = p.GetFile("picture")
	print(picture)
	uid, err := models.Upload(id, picture)
	if err != nil {
		p.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		p.Data["json"] = map[string]string{"image_url": uid}
	}
	p.ServeJSON()
}
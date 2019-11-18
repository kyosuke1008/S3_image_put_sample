// @APIVersion 0.9.0
// @Title APPO API
// @Description APPOのswagger
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl https://www.tokyotower.co.jp/
package routers

import (
	"S3_image_put_sample/controllers"
	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/image",
			beego.NSInclude(
				&controllers.S3ImageController{},
			),
		),
	)
	beego.AddNamespace(ns)
}

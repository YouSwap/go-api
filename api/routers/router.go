// @APIVersion 1.0.0
// @Title Youswap API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://api.youswap.info/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"v1-go-api/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		//beego.NSNamespace("/bundle",
		//	beego.NSInclude(
		//		&controllers.BundleController{},
		//	),
		//),
		//beego.NSNamespace("/buyback",
		//	beego.NSInclude(
		//		&controllers.BuyBackController{},
		//	),
		//),
		//beego.NSNamespace("/pair",
		//	beego.NSInclude(
		//		&controllers.PairController{},
		//	),
		//),
		//beego.NSNamespace("/pool",
		//	beego.NSInclude(
		//		&controllers.PoolController{},
		//	),
		//),
		//beego.NSNamespace("/token",
		//	beego.NSInclude(
		//		&controllers.TokenController{},
		//	),
		//),
		beego.NSNamespace("/tvl",
			beego.NSInclude(
				&controllers.TvlController{},
			),
		),
		beego.NSNamespace("/airdrop",
			beego.NSInclude(
				&controllers.AirdropController{},
			),
		),
		beego.NSNamespace("/pool",
			beego.NSInclude(
				&controllers.RewardPoolController{},
			),
		),
	)
	beego.AddNamespace(ns)
}

package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"v1-go-api/models"
)

// Operations about Tokens
type TvlController struct {
	beego.Controller
}

// @Title GetAll
// @Description get all TvlInfo
// @Success 200 {object} map
// @router / [get]
func (u *TvlController) GetAll() {
	tvl := models.GetAllTvls()

	tvlMap := make(map[string]float64)
	for k, v := range tvl {
		tvlMap[k] = v.TotalAmount
		tvlMap["total"] += v.TotalAmount

	}
	u.Data["json"] = tvlMap
	u.ServeJSON()
}

// @Title Get
// @Description get tvl by id
// @Param	id		path 	string	true		"The key for id"
// @Success 200 {object} models.Tvl
// @Failure 403 :id is empty
// @router /:id [get]
func (u *TvlController) Get() {
	id := u.GetString(":id")
	if id != "" {
		token, err := models.GetTvl(id)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = token
		}
	}
	u.ServeJSON()
}

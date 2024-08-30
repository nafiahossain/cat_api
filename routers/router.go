package routers

import (
	"cat_api/controllers"
	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/", &controllers.CatController{})
	web.Router("/api/cat", &controllers.CatController{}, "get:GetCatData")
	web.Router("/api/breeds", &controllers.CatController{}, "get:GetBreeds")
	web.Router("/api/breed/:id", &controllers.CatController{}, "get:GetBreedInfo")
	web.Router("/api/favorites", &controllers.CatController{}, "post:AddFavorite")
    web.Router("/api/favorites", &controllers.CatController{}, "get:GetFavorites")
	web.Router("/api/vote", &controllers.CatController{}, "post:SubmitVote")
}

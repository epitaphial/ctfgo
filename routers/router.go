package routers

import (
	"ctfgo/controllers"
	"github.com/astaxie/beego"
	"ctfgo/models"
	"github.com/astaxie/beego/context"
)

var FilterSetup = func(ctx *context.Context) {
    game,_ := models.GetGameSetting()
    if !game.IfSetup && ctx.Request.RequestURI != "/setup"{
        ctx.Redirect(302, "/setup")
	}else if game.IfSetup && ctx.Request.RequestURI == "/setup"{
		ctx.Abort(404,"404")
	}
}

func init() {

	//routers
	beego.InsertFilter("/*",beego.BeforeRouter,FilterSetup)
	//安装页面
	beego.Router("/setup",&controllers.SetupController{})

	//主页
	beego.Router("/", &controllers.IndexController{})

	//个人设置
	beego.Router("/usersetting", &controllers.UserSettingController{})
	beego.Router("/changepwd", &controllers.ChangePwdController{})

	//登录注册相关
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/logout", &controllers.LogoutController{})
	beego.Router("/signup", &controllers.RegisterController{})
	beego.Router("/active/:user/:activestring", &controllers.ActiveUserController{})

	//管理页面
	beego.Router("/admin", &controllers.AdminController{})
	beego.Router("/admin/gamesetting", &controllers.GameManageController{})
	beego.Router("/admin/subjects", &controllers.AdminSubjectsController{})
	beego.Router("/admin/subjects/add", &controllers.SubjectsAddController{})
	beego.Router("/admin/subjects/delete/:id", &controllers.SubjectsDeleteController{})
	beego.Router("/admin/subjects/edit/:id", &controllers.SubjectsEditController{})

	//比赛页面
	beego.Router("/game", &controllers.GameController{})
	beego.Router("/rank", &controllers.RankController{})

	//题目文件上传、下载与删除
	beego.Router("/admin/subjects/file/upload/:id", &controllers.SubjectsFileUploadController{})
	beego.Router("/game/file/download/:id", &controllers.SubjectFileDownloadController{})
	beego.Router("/admin/subjects/file/delete/:id", &controllers.SubjectsFileDeleteController{})	
}

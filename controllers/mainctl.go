package controllers

import (
	"ctfgo/models"
	"ctfgo/tools"
	_ "fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/utils/captcha"
	"html/template"
	"net/url"
	"strings"
	"time"
)

//验证码生成器
var cpt *captcha.Captcha
var gamecommon models.Game

func init() {
	store := cache.NewMemoryCache()
	cpt = captcha.NewWithFilter("/captcha/", store)
	gamecommon, _ = models.GetGameSetting()
}

//安装页面的控制器
type SetupController struct {
	beego.Controller
}

func (c *SetupController) Prepare() {
	c.EnableXSRF = true
}

//Get method to the controller.
func (c *SetupController) Get() {
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		c.Data["Notice"] = true
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["Error"] = true
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "setup.html"
}

//Post method to the controller.
func (c *SetupController) Post() {
	flash := beego.NewFlash()
	gamename := strings.TrimSpace(c.GetString("gamename"))
	username := strings.TrimSpace(c.GetString("adminname"))
	password := strings.TrimSpace(c.GetString("password"))
	vrpassword := strings.TrimSpace(c.GetString("veripassword"))
	if password != vrpassword {
		flash.Error("两次输入密码不一致！")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/signup")
		return
	}
	email := strings.TrimSpace(c.GetString("email"))
	activestring := tools.Md5Encode(time.Now().String())
	if status := models.RegisterUser(username, password, email, activestring,1,true); status == models.WellOp {
		var game models.Game
		game.IfSetup = true
		game.GameName = gamename
		game.Id = 1
		game.IfUseEmail = false
		if status := models.GameSetting(game); status == models.WellOp{
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			gamecommon, _ = models.GetGameSetting()
			return
		}else {
			flash.Error("数据库错误")
			//TODO:数据库错误后清空user表格
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/setup")
			return
		}
	} else {
		flash.Error("数据库错误")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/setup")
		return
	}
}

//主页的控制器
type IndexController struct {
	beego.Controller
}

func (c *IndexController) Prepare() {
	userSess := c.GetSession("user")
	if userSess != nil {
		c.Data["IsLogin"] = true
		c.Data["UserName"] = userSess
		adminSess := c.GetSession("admin")
		if adminSess != nil {
			c.Data["IsAdmin"] = true
		}
	}
	c.Data["GameName"] = gamecommon.GameName
}

//Get method to the controller.
func (c *IndexController) Get() {
	c.TplName = "index.html"
}

//登录页面的控制器
type LoginController struct {
	beego.Controller
}

func (c *LoginController) Prepare() {
	c.EnableXSRF = true
	userSess := c.GetSession("user")
	if userSess != nil {
		c.Data["IsLogin"] = true
		c.Data["UserName"] = userSess
		adminSess := c.GetSession("admin")
		if adminSess != nil {
			c.Data["IsAdmin"] = true
		}
		c.Ctx.Redirect(302, "/")
	}
	c.Data["GameName"] = gamecommon.GameName
}

//Get method to the controller.
func (c *LoginController) Get() {
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["Error"] = true
	}

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "login.html"
}

//Post method to the controller
func (c *LoginController) Post() {
	flash := beego.NewFlash()
	if cpt.VerifyReq(c.Ctx.Request) {

	} else {
		flash.Error("验证码错误!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/login")
		return
	}
	username := strings.TrimSpace(c.GetString("username"))
	password := strings.TrimSpace(c.GetString("password"))
	var state models.State
	state = models.LoginUser(username, password)
	if state == models.PassWrong {
		flash.Error("密码错误!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/login")
		return
	} else if state == models.NoExistUser {
		flash.Error("用户不存在!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/login")
		return
	} else if state == models.WellOp {
		c.SetSession("user", username)

		state, isadmin := models.IfAdmin(username)
		if state != models.WellOp {
			c.Ctx.Redirect(302, "/")
			return
		}
		if isadmin {
			c.SetSession("admin", username)
		}
		c.Ctx.Redirect(302, "/")
	} else if state == models.NoActive {
		flash.Error("用户未激活，请先前往邮箱激活!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/login")
		return
	} else {
		flash.Error("数据库错误!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/login")
		return
	}
}

//注册页面的控制器
type RegisterController struct {
	beego.Controller
}

func (c *RegisterController) Prepare() {
	c.EnableXSRF = true
	userSess := c.GetSession("user")
	if userSess != nil {
		c.Data["IsLogin"] = true
		c.Data["UserName"] = userSess
		adminSess := c.GetSession("admin")
		if adminSess != nil {
			c.Data["IsAdmin"] = true
		}
	}
	c.Data["GameName"] = gamecommon.GameName
}

//Get method to the controller.
func (c *RegisterController) Get() {
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		c.Data["Notice"] = true
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["Error"] = true
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "signup.html"
}

//Post method to the controller.
func (c *RegisterController) Post() {
	flash := beego.NewFlash()
	username := strings.TrimSpace(c.GetString("username"))
	password := strings.TrimSpace(c.GetString("password"))
	vrpassword := strings.TrimSpace(c.GetString("veripassword"))
	if password != vrpassword {
		flash.Error("两次输入密码不一致！")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/signup")
		return
	}
	email := strings.TrimSpace(c.GetString("email"))
	activestring := tools.Md5Encode(time.Now().String())
	gameconfig, _ := models.GetGameSetting()
	if status := models.RegisterUser(username, password, email, activestring,0,!gameconfig.IfUseEmail); status == models.WellOp {
		if gameconfig.IfUseEmail{
			emailhost, _ := url.Parse(gameconfig.EmailHost)
			gameurl, _ := url.Parse(gameconfig.GameUrl)
			tools.SendEmailActive(email, username, activestring, emailhost.Hostname(), gameurl.Hostname(), gameurl.Port(), gameconfig.EmailAcount, gameconfig.EmailPass, gameconfig.EmailPort)
			flash.Notice("注册成功！请前往邮箱激活。")
		}else{
			flash.Notice("注册成功！")
		}

		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/signup")
		return
	} else {
		if status == models.EmailRepeat {
			flash.Error("邮箱已被注册！")
		} else if status == models.UserRepeat {
			flash.Error("用户名已存在！")
		} else {
			flash.Error("数据库错误")
		}
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/signup")
		return
	}
}

//登出的控制器
type LogoutController struct {
	beego.Controller
}

func (c *LogoutController) Prepare() {
	userSess := c.GetSession("user")
	if userSess == nil {
		c.Ctx.Redirect(302, "/")
	}
}

func (c *LogoutController) Get() {
	c.DestroySession()
	c.Ctx.Redirect(302, "/")
}

func (c *LogoutController) Post() {
	c.DestroySession()
	c.Ctx.Redirect(302, "/")
}

//激活页面的控制器
type ActiveUserController struct {
	beego.Controller
}

func (c *ActiveUserController) Prepare() {
	userSess := c.GetSession("user")
	if userSess != nil {
		c.Data["IsLogin"] = true
		c.Data["UserName"] = userSess
		adminSess := c.GetSession("admin")
		if adminSess != nil {
			c.Data["IsAdmin"] = true
		}
	}
	c.Data["GameName"] = gamecommon.GameName
}

func (c *ActiveUserController) Get() {
	username := c.Ctx.Input.Param(":user")
	activestring := c.Ctx.Input.Param(":activestring")
	if status := models.ActiveUser(username, activestring); status == models.WellOp {
		c.Data["Info"] = "您的账户已经激活成功！"
	} else if status == models.NoExistUser {
		c.Ctx.Redirect(302, "/")
		return
	} else if status == models.FailActive {
		c.Data["Info"] = "您的账户激活失败！"
	} else if status == models.DatabaseErr {
		c.Data["Info"] = "数据库错误，请联系管理员！"
	} else if status == models.ActiveRepeat {
		c.Data["Info"] = "您已激活,请勿重复激活！"
	}
	c.TplName = "active.html"
}

//个人设置页面的控制器
type UserSettingController struct {
	beego.Controller
}

//Prepare method to the controller.
func (c *UserSettingController) Prepare() {
	c.EnableXSRF = true
	userSess := c.GetSession("user")
	if userSess == nil {
		c.Ctx.Redirect(302, "/login")
	} else {
		c.Data["IsLogin"] = true
		c.Data["UserName"] = userSess
		adminSess := c.GetSession("admin")
		if adminSess != nil {
			c.Data["IsAdmin"] = true
		}
	}
	c.Data["GameName"] = gamecommon.GameName
}

func (c *UserSettingController) Get() {
	username := c.GetSession("user").(string)
	_, user := models.GetUserInfo(username)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["username"] = user.Username
	c.Data["name"] = user.Name
	c.Data["userid"] = user.Stuid
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		c.Data["Notice"] = true
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["Error"] = true
	}
	c.TplName = "usersetting.html"
}

func (c *UserSettingController) Post() {
	flash := beego.NewFlash()
	name := strings.TrimSpace(c.GetString("name"))
	stuid := strings.TrimSpace(c.GetString("stuid"))
	username := c.GetSession("user").(string)
	var user models.User
	user.Name = name
	user.Stuid = stuid
	user.Username = username
	if models.UpdateUserInfo(user) == models.WellOp {
		flash.Notice("修改成功！")
		flash.Store(&c.Controller)
	} else {
		flash.Error("修改失败！")
		flash.Store(&c.Controller)
	}
	c.Ctx.Redirect(302, "/usersetting")
}

//修改密码页面的控制器
type ChangePwdController struct {
	beego.Controller
}

//Prepare method to the controller.
func (c *ChangePwdController) Prepare() {
	c.EnableXSRF = true
	userSess := c.GetSession("user")
	if userSess == nil {
		c.Ctx.Redirect(302, "/login")
	} else {
		c.Data["IsLogin"] = true
		c.Data["UserName"] = userSess
		adminSess := c.GetSession("admin")
		if adminSess != nil {
			c.Data["IsAdmin"] = true
		}
	}
	c.Data["GameName"] = gamecommon.GameName
}

func (c *ChangePwdController) Get() {
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		c.Data["Notice"] = true
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["Error"] = true
	}
	c.TplName = "changepass.html"
}

//Post method to the controller.
func (c *ChangePwdController) Post() {
	username := c.GetSession("user").(string)
	flash := beego.NewFlash()
	oldpass := strings.TrimSpace(c.GetString("oldpass"))
	password := strings.TrimSpace(c.GetString("password"))
	vrpassword := strings.TrimSpace(c.GetString("veripassword"))
	if password != vrpassword {
		flash.Error("两次输入密码不一致！")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/changepwd")
		return
	}
	if password == oldpass {
		flash.Error("新旧密码不能一致！")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/changepwd")
		return
	}
	if status := models.UpdatePassword(username, oldpass, password); status == models.WellOp {
		flash.Notice("修改成功！")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/changepwd")
		return
	} else {
		if status == models.NewAndOldDiff {
			flash.Error("新旧密码不一致！")
		} else {
			flash.Error("数据库错误")
		}
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/changepwd")
		return
	}

}

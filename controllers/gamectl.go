package controllers

import (
	"ctfgo/models"
	_ "fmt"
	"github.com/astaxie/beego"
	"html/template"
	"strings"
	"strconv"
)

//比赛题目页面
type GameController struct {
	beego.Controller
}

//Prepare method to the controller.
func (c *GameController) Prepare() {
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

//Get method to the controller.
func (c *GameController) Get() {
	userName := c.GetSession("user").(string)
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		c.Data["FlagRight"] = true
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["FlagWrong"] = true
	}
	status, subjects := models.GetUnhiddenSubjects()
	var subjectType map[string]bool
	subjectType = make(map[string]bool)
	for index, subject := range subjects {
		subjectType[subject.SubType] = true
		if models.IfSolved(subject.Id, userName) == models.HasRightSubmit {
			subjects[index].IfDone = true
		} else {
			subjects[index].IfDone = false
		}
	}
	if status == models.WellOp {
		c.Data["Subjects"] = subjects
	}
	c.Data["SubjectType"] = subjectType
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	allfiles,state :=  models.GetAllFiles()
	if state != models.WellOp{
		c.Ctx.Redirect(302, "/game")
	}
	c.Data["allfiles"] = allfiles
	c.TplName = "game.html"
}

//提交flag，判断正确性
func (c *GameController) Post() {
	subjectId := strings.TrimSpace(c.GetString("subjectid"))
	userFlag := strings.TrimSpace(c.GetString("userflag"))
	userName := c.GetSession("user").(string)
	flash := beego.NewFlash()
	if models.UserCommitFlag(subjectId, userFlag, userName) == models.FlagWrong {
		flash.Error("Flag错误!")
		flash.Store(&c.Controller)
		c.Redirect("/game", 302)
		return
	} else {
		flash.Notice("Flag正确!")
		flash.Store(&c.Controller)
		c.Redirect("/game", 302)
		return
	}
}

//排行榜页面
type RankController struct {
	beego.Controller
}

func (c *RankController) Prepare() {
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

//Get method to the controller.
func (c *RankController) Get() {
	_, users := models.GetUnhiddenUsers()
	c.Data["UsersRank"] = users
	c.TplName = "rank.html"
}

//题目附件下载
type SubjectFileDownloadController struct {
	beego.Controller
}

//Prepare method to the controller.
func (c *SubjectFileDownloadController) Prepare() {
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

func (c *SubjectFileDownloadController) Get() {
	fileid, errors := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if errors != nil {
		c.Ctx.Redirect(302, "/game")
		return
	}
	subfile,state := models.GetFileById(fileid)
	if state != models.WellOp{
		c.Redirect("/game", 302)		
	}
	c.Ctx.Output.Download("upload/"+subfile.Md5FileName,subfile.FileName)
}
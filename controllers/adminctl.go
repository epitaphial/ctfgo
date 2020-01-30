package controllers

import (
	"ctfgo/models"
	_"fmt"
	"github.com/astaxie/beego"
	"html/template"
	"strconv"
	"strings"
	"io"
	"os"
	"time"
	"ctfgo/tools"
)

//管理页面主页
type AdminController struct {
	beego.Controller
}

func (c *AdminController) Prepare() {
	adminSess := c.GetSession("admin")
	if adminSess == nil {
		c.Ctx.Abort(404,"404")
	} else {
		adminSess := c.GetSession("user")
		c.Data["IsLogin"] = true
		c.Data["IsAdmin"] = true
		c.Data["UserName"] = adminSess
		c.Data["IsAdminPage"] = true
	}
}

//Get method to the controller.
func (c *AdminController) Get() {
	c.TplName = "admin/index.html"
}

//题目管理页面
type AdminSubjectsController struct {
	beego.Controller
}

func (c *AdminSubjectsController) Prepare() {
	adminSess := c.GetSession("admin")
	if adminSess == nil {
		c.Ctx.Abort(404,"404")
	} else {
		adminSess := c.GetSession("admin")
		c.Data["IsLogin"] = true
		c.Data["IsAdmin"] = true
		c.Data["UserName"] = adminSess
		c.Data["IsAdminPage"] = true
	}
}

//Get method to the controller.
func (c *AdminSubjectsController) Get() {
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		c.Data["welldone"] = true //在题目页面flash成功消息，修改，增加，删除成功消息
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["baddone"] = true //失败消息
	}
	var subjects []models.Subject
	_, subjects = models.GetSubjects()
	c.Data["Subjects"] = subjects
	c.TplName = "admin/subjects.html"
	return
}

//题目编辑页面，动态路由
type SubjectsEditController struct {
	beego.Controller
}

func (c *SubjectsEditController) Prepare() {
	c.EnableXSRF = true
	adminSess := c.GetSession("admin")
	if adminSess == nil {
		c.Ctx.Abort(404,"404")
	} else {
		adminSess := c.GetSession("admin")
		c.Data["IsLogin"] = true
		c.Data["IsAdmin"] = true
		c.Data["UserName"] = adminSess
		c.Data["IsAdminPage"] = true
	}
}

//Get method to the controller.
func (c *SubjectsEditController) Get() {
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		c.Data["UploadOk"] = true
	} else if _, ok = flash.Data["error"];ok  {
		c.Data["EditError"] = true
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	EditId, error := strconv.Atoi(c.Ctx.Input.Param(":id"))
	c.Data["EditId"] = EditId
	if error != nil {
		c.Ctx.Redirect(302, "/admin/subjects")
		return
	}
	state, subject := models.GetSubject(EditId)
	if state != models.WellOp {
		c.Ctx.Redirect(302, "/admin/subjects")
		return
	}
	subjectid, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	subfiles,_ := models.GetSubjectFile(subjectid)
	c.Data["SubFiles"] = subfiles
	c.Data["Subject"] = subject
	c.TplName = "admin/subjectedit.html"
	return
}

//Post method to the controller.
func (c *SubjectsEditController) Post() {
	var subject models.Subject
	var errors error
	flash := beego.ReadFromRequest(&c.Controller)
	subject.Id, errors = strconv.Atoi(c.Ctx.Input.Param(":id"))
	if errors != nil {
		c.Ctx.Redirect(302, "/admin/subjects/")
		return
	}
	subject.SubName = strings.TrimSpace(c.GetString("subname"))
	subject.SubType = strings.TrimSpace(c.GetString("subtype"))
	subject.SubFlag = strings.TrimSpace(c.GetString("subflag"))
	subject.SubDescribe = strings.TrimSpace(c.GetString("subdescribe"))
	subject.SubMark, errors = c.GetInt("submark")
	if errors != nil {
		flash.Error("分值必须为整数!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/edit/"+c.Ctx.Input.Param(":id"))
		return
	}
	if strings.TrimSpace(c.GetString("ifhidden")) == "on" {
		subject.IfHidden = true
	} else {
		subject.IfHidden = false
	}
	models.EditSubject(subject)
	flash.Notice("修改成功!")
	flash.Store(&c.Controller)
	c.Ctx.Redirect(302, "/admin/subjects/")
}

//添加题目
type SubjectsAddController struct {
	beego.Controller
}

func (c *SubjectsAddController) Prepare() {
	c.EnableXSRF = true
	adminSess := c.GetSession("admin")
	if adminSess != nil {
		c.Data["IsLogin"] = true
		c.Data["IsAdmin"] = true
		c.Data["UserName"] = adminSess
		c.Data["IsAdminPage"] = true
	} else {
		c.Data["IsLogin"] = false
		c.Ctx.Abort(404,"404")
	}
}

//Get method to the controller.
func (c *SubjectsAddController) Get() {
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["adderror"] = true
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "admin/subjectadd.html"
	return
}

//Post method to the controller.
func (c *SubjectsAddController) Post() {
	var subject models.Subject
	var errors error
	flash := beego.NewFlash()
	subject.SubName = strings.TrimSpace(c.GetString("subname"))
	subject.SubType = strings.TrimSpace(c.GetString("subtype"))
	subject.SubFlag = strings.TrimSpace(c.GetString("subflag"))
	subject.SubDescribe = strings.TrimSpace(c.GetString("subdescribe"))
	subject.SubMark, errors = c.GetInt("submark")
	if errors != nil {
		flash.Error("分值必须为整数!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/add")
		return
	}
	if strings.TrimSpace(c.GetString("ifhidden")) == "on" {
		subject.IfHidden = true
	} else {
		subject.IfHidden = false
	}
	if models.AddSubject(subject) != models.WellOp {
		flash.Error("添加题目失败!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/add")
		return
	} else {
		flash.Notice("添加题目成功!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/")
		return
	}

}

//删除题目
type SubjectsDeleteController struct {
	beego.Controller
}

func (c *SubjectsDeleteController) Prepare() {
	c.EnableXSRF = true
	adminSess := c.GetSession("admin")
	if adminSess != nil {
		c.Data["IsLogin"] = true
		c.Data["UserName"] = adminSess
		c.Data["IsAdmin"] = true
		c.Data["IsAdminPage"] = true
	} else {
		c.Data["IsLogin"] = false
		c.Ctx.Abort(404,"404")
	}
}

//Get method to the controller.
func (c *SubjectsDeleteController) Get() {
	flash := beego.NewFlash()
	subjectId, errors := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if errors != nil {
		flash.Error("删除失败!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/")
		return
	}
	if models.DeleteSubject(subjectId) == models.WellOp {
		flash.Notice("删除成功!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/")
	} else {
		flash.Error("删除失败!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/")
		return
	}
	return
}

//管理比赛
type GameManageController struct {
	beego.Controller
}

func (c *GameManageController) Prepare() {
	c.EnableXSRF = true
	adminSess := c.GetSession("admin")
	if adminSess != nil {
		c.Data["IsLogin"] = true
		c.Data["IsAdmin"] = true
		c.Data["UserName"] = adminSess
		c.Data["IsAdminPage"] = true
	} else {
		c.Data["IsLogin"] = false
		c.Ctx.Abort(404,"404")
	}
}

func (c *GameManageController) Get() {
	game, _ := models.GetGameSetting()
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		c.Data["Notice"] = true
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["Error"] = true
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["game"] = game
	c.TplName = "admin/gamesettings.html"
	return
}

func (c *GameManageController) Post() {
	gamename := strings.TrimSpace(c.GetString("gamename"))
	domainname := strings.TrimSpace(c.GetString("domainname"))
	emailserver := strings.TrimSpace(c.GetString("emailserver"))
	emailport, _ := c.GetInt("emailport")
	emailaccount := strings.TrimSpace(c.GetString("emailaccount"))
	emailpass := strings.TrimSpace(c.GetString("emailpass"))
	var game models.Game
	game.GameName = gamename
	game.GameUrl = domainname
	game.EmailHost = emailserver
	game.EmailPort = int(emailport)
	game.EmailAcount = emailaccount
	game.EmailPass = emailpass
	game.Id = 1
	game.IfSetup = true
	if strings.TrimSpace(c.GetString("ifuseemail")) == "on" {
		game.IfUseEmail = true
	} else {
		game.IfUseEmail = false
	}
	flash := beego.NewFlash()
	if models.GameSetting(game) == models.WellOp {
		flash.Notice("修改成功！")
		flash.Store(&c.Controller)
		gamecommon, _ = models.GetGameSetting()
	} else {
		flash.Error("修改失败！")
		flash.Store(&c.Controller)
	}
	c.Ctx.Redirect(302, "/admin/gamesetting/")
	return
}


//题目文件上传
type SubjectsFileUploadController struct {
	beego.Controller
}

func (c *SubjectsFileUploadController) Prepare() {
	adminSess := c.GetSession("admin")
	if adminSess == nil {
		c.Ctx.Abort(404,"404")
	} else {
		adminSess := c.GetSession("admin")
		c.Data["IsLogin"] = true
		c.Data["IsAdmin"] = true
		c.Data["UserName"] = adminSess
		c.Data["IsAdminPage"] = true
	}
}

//Post method to the controller.
func (c *SubjectsFileUploadController) Post() {
	flash := beego.ReadFromRequest(&c.Controller)
	subjectid, errors := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if errors != nil {
		flash.Error("路径错误!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/")
		return
	}
	files, err := c.GetFiles("files")
	if err != nil {
		c.Ctx.WriteString("Invalid file")
		return
	}
	for i,_ := range files{
		file,err := files[i].Open()
		defer file.Close()
		if err != nil {
			flash.Error("上传失败!")
			flash.Store(&c.Controller)
			break
		}
		//判断文件夹是否存在，不存在就创建文件夹
		_, err = os.Stat("upload")
		if os.IsNotExist(err) {
			err := os.Mkdir("upload", os.ModePerm)
			if err != nil {
				flash.Error("上传失败!")
				flash.Store(&c.Controller)
				break
				}
			}
		md5filename := tools.Md5Encode(files[i].Filename+time.Now().String())
		dst, err := os.Create("upload/" + md5filename)
		defer dst.Close()
		if err != nil {
			flash.Error("上传失败!")
			flash.Store(&c.Controller)
			break
		}
		if _, err := io.Copy(dst, file); err != nil {
			flash.Error("上传失败!")
			flash.Store(&c.Controller)
			break
		}
		if models.UploadSubjectFile(files[i].Filename,md5filename,subjectid) != models.WellOp{
			flash.Error("数据库错误!")
			flash.Store(&c.Controller)
			break			
		}
	}
	if err == nil{
		flash.Notice("上传成功!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/edit/"+strconv.Itoa(subjectid)+"#theupload")
	}else{
		c.Ctx.Redirect(302, "/admin/subjects")
	}

}

//题目文件删除
type SubjectsFileDeleteController struct {
	beego.Controller
}

func (c *SubjectsFileDeleteController) Prepare() {
	adminSess := c.GetSession("admin")
	if adminSess == nil {
		c.Ctx.Abort(404,"404")
	} else {
		adminSess := c.GetSession("admin")
		c.Data["IsLogin"] = true
		c.Data["IsAdmin"] = true
		c.Data["UserName"] = adminSess
		c.Data["IsAdminPage"] = true
	}
}

func (c *SubjectsFileDeleteController) Get(){
	flash := beego.ReadFromRequest(&c.Controller)
	fileid, errors := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if errors != nil {
		flash.Error("路径错误!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/")
		return
	}
	subfile,state := models.GetFileById(fileid)
	md5filename := subfile.Md5FileName
	if state != models.WellOp{
		flash.Error("数据库错误!")
		flash.Store(&c.Controller)			
	}else{
		if models.DeleteSubjectFile(fileid) != models.WellOp{
			flash.Error("附件删除失败!")
			flash.Store(&c.Controller)	
		}else{
			err := os.Remove("upload/" + md5filename)
			if err != nil {
				flash.Error("附件删除失败!")
				flash.Store(&c.Controller)
			}else {
				flash.Notice("附件删除成功！")
				flash.Store(&c.Controller)
			}
		}
	}
	c.Ctx.Redirect(302, "/admin/subjects/")
	return
}
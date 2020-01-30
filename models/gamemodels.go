package models

import (
	//_ "github.com/go-sql-driver/mysql"
	_"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_"fmt"
	"strconv"
	"time"
	"os"
)


//题目的数据库模型
type Subject struct {
	Id       int
	IfHidden bool    //题目是否隐藏
	SubName     string `orm:"size(100)"` //题目名称
	SubMark		int //题目分数
	SubFlag		string `orm:"size(200)"` //flag
	SubDescribe string `orm:"size(1000)"` //题目描述
	SubType		string `orm:"size(50)"` //题目类型
	CreatedTime	time.Time	`orm:"auto_now_add;type(datetime)"`
	UpdatedTime	time.Time	`orm:"auto_now;type(datetime)"`
	IfDone	bool	//临时确定当前用户是否正确回答该题
}

//比赛的数据库模型
type Game struct {
	Id	int
	IfSetup	bool//是否配置完成
	GameName	string//比赛名
	GameUrl	string//比赛域名
	IfUseEmail	bool
	EmailHost	string//邮件服务器
	EmailPort	int//服务器端口
	EmailAcount	string//邮件账户
	EmailPass	string//邮件密码
	CreatedTime	time.Time	`orm:"auto_now_add;type(datetime)"`
	UpdatedTime	time.Time	`orm:"auto_now;type(datetime)"`			
}

//题目附件的数据库模型
type SubjectFile struct{
	Id	int
	SubId	int//对应题目的ID
	FileName	string//下载下来的文件名
	Md5FileName	string//Md5后存储的文件名
	CreatedTime	time.Time	`orm:"auto_now_add;type(datetime)"`
	UpdatedTime	time.Time	`orm:"auto_now;type(datetime)"`		
}


//取得所有题目的数组，用来进一步操作或者渲染页面
func GetSubjects() (state State,subjects []Subject){
	o := orm.NewOrm()
	//o.QueryTable("subject").Filter("Status", 1).All(&subjects, "Id", "Title")
	o.QueryTable("subject").All(&subjects, "Id", "SubName","SubType")
	state = WellOp
	return
}

//获取某一道题目信息，用来编辑
func GetSubject(id int) (state State,subject Subject){
	o := orm.NewOrm()
	subject.Id = id
	err := o.Read(&subject)
	if err == orm.ErrNoRows {
		state = NoSuchSubject
	} else if err == orm.ErrMissPK {
		state = NoSuchKey
	} else {
		state = WellOp
	}
	return
}

//添加题目
func AddSubject(subject Subject) (state State){
	o := orm.NewOrm()
	_, err := o.Insert(&subject)
	if err != nil{
		state = DatabaseErr
		return
	}else{
		state = WellOp
		return
	}
}

//修改题目
func EditSubject(subject Subject) (state State){

	o := orm.NewOrm()
	oldsubject := Subject{Id:subject.Id}
	if o.Read(&oldsubject) == nil {
		oldsubject.IfHidden = subject.IfHidden
		oldsubject.SubName = subject.SubName
		oldsubject.SubMark = subject.SubMark
		oldsubject.SubFlag	 = subject.SubFlag
		oldsubject.SubDescribe = subject.SubDescribe
		oldsubject.SubType = subject.SubType
    if _, err := o.Update(&oldsubject); err == nil {
		state = WellOp
		return
	}
	state = NoSuchSubject
	return
}else{
	state = DatabaseErr
	return
}
}

//删除题目,并删除相应的文件
func DeleteSubject(id int) (state State){
	o := orm.NewOrm()
	var subfiles []SubjectFile
	subfiles,state = GetSubjectFile(id)
	if state != WellOp{
		return
	}
	_, err := o.QueryTable("subject_file").Filter("sub_id", id).Delete()
	if err !=nil{
		state = DatabaseErr
		return
	}
	for i,_ := range subfiles{
		err := os.Remove("upload/" + subfiles[i].Md5FileName)
		if err != nil{
			state = FileDeleteError
			return
		}
	}
	if id, err := o.Delete(&Subject{Id: id}); err == nil {
	state = WellOp
	if id == 0{
		state = NoSuchSubject
	}
	}else{
		state = DatabaseErr
	}
	return
}

//获取所有未隐藏题目
func GetUnhiddenSubjects() (state State,subjects []Subject){
	o := orm.NewOrm()
	//o.QueryTable("subject").Filter("Status", 1).All(&subjects, "Id", "Title")
	o.QueryTable("subject").Filter("IfHidden", false).All(&subjects, "Id", "SubName","SubType","SubMark","SubDescribe")
	state = WellOp
	return
}

//提交flag，并记录提交历史
func UserCommitFlag(subjectId,userFlag,userName string) (state State){
	o := orm.NewOrm()
	subject := new(Subject)
	var errors error
	subject.Id,errors = strconv.Atoi(subjectId)
	if errors != nil{
		state = NoSuchId
		return
	}
	errors = o.Read(subject,"Id")
	if errors != nil{
		state = NoSuchSubject
		return
	}
	if subject.SubFlag != userFlag{
		state = FlagWrong
		WrongSubmit(subject.Id,userName,userFlag)
		return
	}else{
		if IfSolved(subject.Id,userName) == NoRightSubmit{
			state = EditUserMark(userName,subject.SubMark)
			RightSubmit(subject.Id,userName)
			return
		}else{
			return
		}
	}

}

//加减分操作
func EditUserMark(userName string,userMark int) (state State){
	o := orm.NewOrm()
	user := new(User)
	user.Username = userName
	if o.Read(user,"Username") == nil {
		user.Mark += userMark 
		if _, err := o.Update(user,"Mark"); err != nil {
			state = MarkEditWrong
			return
		}
		state = WellOp
		return
	}else{
		state = MarkEditWrong
		return
	}
}

//比赛全局设置
func GameSetting(game Game)(state State){
	o := orm.NewOrm()
	oldgame := game
	if created, _, err := o.ReadOrCreate(&oldgame, "Id"); err == nil {
		if created {
			state = WellOp
		} else {
			oldgame = game
			if _, err := o.Update(&oldgame); err == nil {
				state = WellOp
			}else{
				state = DatabaseErr
			}
		}
	}else{
		state = DatabaseErr
	}
	return
}

//获取比赛全局设置
func GetGameSetting()(game Game,state State){
	o := orm.NewOrm()
	game.Id = 1
	err := o.Read(&game)
	if err == orm.ErrNoRows {
		state = NoSuchId
	} else if err == orm.ErrMissPK {
		state = NoSuchKey
	} else {
		state = WellOp
	}
	return
}

//记录题目对应的文件
func UploadSubjectFile(filename,md5filename string,subjectid int)(state State){
	o := orm.NewOrm()
	var subfile SubjectFile
	subfile.Md5FileName = md5filename
	subfile.FileName = filename
	subfile.SubId = subjectid
	_, err := o.Insert(&subfile)
	if err == nil {
		state = WellOp
	}else{
		state = DatabaseErr
	}
	return
}

//删除题目对应的文件
func DeleteSubjectFile(fileid int)(state State){
	o := orm.NewOrm()
	if _, err := o.Delete(&SubjectFile{Id: fileid}); err == nil {
		state = WellOp
	}else{
		state = DatabaseErr
	}
	return
}

//获取题目对应的文件
func GetSubjectFile(subjectid int)(subfile []SubjectFile,state State){
	o := orm.NewOrm()
	_, err := o.QueryTable("subject_file").Filter("sub_id", subjectid).All(&subfile)
	if err != nil{
		state = DatabaseErr
	}else{
		state = WellOp
	}
	return
}

//获取所有文件
func GetAllFiles()(subfile []SubjectFile,state State){
	o := orm.NewOrm()
	_, err := o.QueryTable("subject_file").All(&subfile)
	if err != nil{
		state = DatabaseErr
	}else{
		state = WellOp
	}
	return
}

//获取文件id对应的md5文件名
func GetFileById(fileid int)(subfile SubjectFile,state State){
	o := orm.NewOrm()
	subfile.Id = fileid
	err := o.Read(&subfile,"Id")
	if err == orm.ErrNoRows {
		state = NoSuchId
	} else if err == orm.ErrMissPK{
		state = NoSuchKey
	}else{
		state = WellOp
		}
	return
}
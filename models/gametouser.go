package models

import (
	//_ "github.com/go-sql-driver/mysql"
	_"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_"fmt"
	"time"
)

//错误提交记录的表
type WrongSubmitTable struct{
	Id	int
	SubmitFlag	string	//提交的内容
	UserName	string	//用户名
	SubjectId	int	//题目id
	CreatedTime	time.Time	`orm:"auto_now_add;type(datetime)"`
}

//正确提交记录的表
type RightSubmitTable struct{
	Id	int
	UserName	string
	SubjectId	int
	CreatedTime	time.Time	`orm:"auto_now_add;type(datetime)"`
}

//判断用户是否已经成功提交过正确的flag
func IfSolved(subjectId int,userName string)(state State){
	o := orm.NewOrm()
	rstable := RightSubmitTable{UserName: userName,SubjectId:subjectId}
	err := o.Read(&rstable,"UserName","SubjectId")
	if err == orm.ErrNoRows {
		state = NoRightSubmit
		return
	} else if err == orm.ErrMissPK {
		state = NoSuchKey
		return
	} else {
		state = HasRightSubmit
		return
	}
}

//记录正确提交的记录
func RightSubmit(subjectId int,userName string)(state State){
	o := orm.NewOrm()
	var rstable RightSubmitTable
	rstable.SubjectId = subjectId
	rstable.UserName = userName
	_, err := o.Insert(&rstable)
	if err == nil {
		state = WellOp
		return
	}else{
		state = DatabaseErr
		return
	}
}

//记录错误提交的记录
func WrongSubmit(subjectId int,userName string,submitFlag string)(state State){
	o := orm.NewOrm()
	var wstable WrongSubmitTable
	wstable.SubjectId = subjectId
	wstable.UserName = userName
	wstable.SubmitFlag = submitFlag
	_, err := o.Insert(&wstable)
	if err == nil {
		state = WellOp
		return
	}else{
		state = DatabaseErr
		return
	}
}
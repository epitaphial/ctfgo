package models

import (
	_"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_"github.com/mattn/go-sqlite3"
	_"fmt"
)

func init() {
	//mysql
	//orm.RegisterDriver("mysql", orm.DRMySQL)
	//databaseconfig := fmt.Sprintf("%s:%s@/%s?charset=UTF8MB4",beego.AppConfig.String("mysqluser"),beego.AppConfig.String("mysqlpass"),beego.AppConfig.String("databasename"))
	//orm.RegisterDataBase("default", "mysql", databaseconfig)
	//sqlite3
	orm.RegisterDriver("sqlite", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", "./datas/ctfgo.db")

	orm.RegisterModel(new(User),new(Subject),new(WrongSubmitTable),new(RightSubmitTable),new(Game),new(SubjectFile))
	orm.RunSyncdb("default", false, true)
	orm.Debug = true
}

type State int

const (
	WellOp  State = iota	//Everything is ok
	DatabaseErr              // 数据库内部错误
	NoSuchKey

	//用户状态
    PassWrong // 密码错误
	UserRepeat            // 已经存在用户（注册时）
	EmailRepeat            // 已经存在Email（注册时）
	NoExistUser   //用户不存在
	MarkEditWrong //修改分数失败
	NoActive //未激活
	FailActive //激活失败
	ActiveRepeat//重复激活
	NewAndOldDiff//新旧密码不一致

	//题目状态
	NoSuchSubject
	NoSuchId
	
	//提交flag状态
	FlagWrong //flag错误
	NoRightSubmit //没有成功提交记录
	HasRightSubmit //有成功提交记录

	//题目文件状态
	FileDeleteError//题目文件删除失败
)
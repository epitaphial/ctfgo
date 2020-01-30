package models

import (
	//_ "github.com/go-sql-driver/mysql"
	_"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ctfgo/tools"
	"time"
	"fmt"
)

//用户的数据模型.
type User struct {
	Id	int
	Mark	int //用户分数
	Name	string `orm:"size(100)"` //用户姓名
	Email	string
	Stuid	string //学号
	Username	string //用户名
	Hashpass	string //哈希后的密码
	Identity	int //标识是管理员(1)还是普通用户(0)
	IfActive	bool//是否激活
	IfHidden	bool//是否隐藏
	ActiveString string//激活链接
	CreatedTime	time.Time	`orm:"auto_now_add;type(datetime)"`
	UpdatedTime	time.Time	`orm:"auto_now;type(datetime)"`
}


//注册用户的数据库操作
func RegisterUser(username, password, email ,activestring string,ifadmin int,ifactive bool)(state State){
	o := orm.NewOrm()
	user := new(User)
	user.Username = username
	user.Email = email
	user.Hashpass = tools.Md5Encode(password)
	user.Identity = ifadmin
	if ifadmin == 1{
		user.IfActive = true
		user.IfHidden = true
	}else{
		user.IfActive = ifactive
		user.IfHidden = false
	}
	user.Mark = 0
	user.ActiveString = activestring
	var err error

	//先查是否用户名重复
	err = o.Read(user,"Username")
	if err == orm.ErrNoRows {
	} else if err == orm.ErrMissPK {
		
	} else {
		state = UserRepeat
		return
	}

	//再查是否邮箱重复
	err = o.Read(user,"Email")
	if err == orm.ErrNoRows {
		
	} else if err == orm.ErrMissPK {
		
	} else {
		state = EmailRepeat
		return
	}

	_, err = o.Insert(user)
	if err != nil {
		state = DatabaseErr
		return
	}else{
		state = WellOp
		return
	}
}

//登录的数据库操作
func LoginUser(username,password string) (state State) {
	o := orm.NewOrm()
	user := new(User)
	user.Username = username
	err := o.Read(user,"Username")
	if err != nil{
		state = NoExistUser
		return
	}
	if user.Hashpass != tools.Md5Encode(password){
		state = PassWrong
		return
	}else if user.IfActive == false{
		state = NoActive
		return
	}else{
		state = WellOp
		return
	}
}

//激活用户的数据库操作
func ActiveUser(username,activestring string) (state State) {
	o := orm.NewOrm()
	user := new(User)
	user.Username = username
	err := o.Read(user,"Username")
	if err != nil{
		state = NoExistUser
		return
	}
	if user.IfActive{
		state = ActiveRepeat
		return
	}

	if user.ActiveString != activestring{
		state = FailActive
		return
	}else{
		user.IfActive = true
		_, err := o.Update(user,"IfActive")
		if err != nil {
			state = DatabaseErr
		}else{
			state = WellOp
		}
		return
	}
}


//判断用户角色的数据库操作

func IfAdmin(username string) (state State,isadmin bool) {
	o := orm.NewOrm()
	user := new(User)
	user.Username = username
	err := o.Read(user,"Username")
	if err != nil{
		state = NoExistUser
		return
	}
	if user.Identity == 1{
		state = WellOp
		isadmin = true
		return
	}else{
		state = WellOp
		isadmin = false
		return
	}

}

//获取所有未隐藏用户的操作（排行榜）
func GetUnhiddenUsers() (state State,users []User){
	o := orm.NewOrm()
	o.QueryTable("user").OrderBy("-Mark").Filter("IfHidden", false).All(&users, "Username","Mark")
	state = WellOp
	return
}

//获取所有用户的操作（用户管理）
func GetUsers() (state State,users []User){
	o := orm.NewOrm()
	o.QueryTable("user").All(&users, "Username","Mark")
	state = WellOp
	return
}

//获取某一用户信息（进行进一步操作）
func GetUserInfo(username string)(state State,user User){
	o := orm.NewOrm()
	user = User{Username:username}
	err := o.Read(&user,"Username")

	if err == orm.ErrNoRows {
		state = NoExistUser
	} else if err == orm.ErrMissPK {
    	state = NoSuchKey
	} else {
    	state = WellOp
	}
	return
}

//更新用户信息
func UpdateUserInfo(user User)(state State){
	o := orm.NewOrm()
	fmt.Println(user.Username,user.Stuid,user.Name)
	var olduser User
	olduser.Username = user.Username
	if o.Read(&olduser,"Username") == nil {
		olduser.Name = user.Name
		olduser.Stuid = user.Stuid
		if _, err := o.Update(&olduser,"Name","Stuid"); err == nil {
			state = WellOp
		}
	}else{
		state = DatabaseErr
	}
	return
}

//修改密码
func UpdatePassword(username,oldpassword,password string) (state State){
	o := orm.NewOrm()
	user := new(User)
	user.Username = username
	user.Hashpass = tools.Md5Encode(oldpassword)
	err := o.Read(user,"Username","Hashpass")

	if err == orm.ErrNoRows {
    	state = NewAndOldDiff
	} else if err == orm.ErrMissPK {
    	state = NoSuchKey
	} else {
			user.Hashpass = tools.Md5Encode(password)
			if _, err := o.Update(user,"Hashpass"); err == nil {
				state = WellOp
			}else{
				state = DatabaseErr
			}
	}
	return
}
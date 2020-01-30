package tools

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
)

//发送激活邮件
func SendEmailActive(reciever, username, activestring, emailhost, gameurl, gameport, emailacount, emailpass string, emailport int) {
	config := fmt.Sprintf(`{"username":"%s","password":"%s","host":"%s","port":%d}`, emailacount, emailpass, emailhost, emailport)
	// 通过存放配置信息的字符串，创建Email对象
	temail := utils.NewEMail(config)
	// 指定邮件的基本信息
	temail.To = []string{reciever} //指定收件人邮箱地址
	temail.From = "ctf@stega.cn"   //指定发件人的邮箱地址
	temail.Subject = "激活你的CTFGO账号" //指定邮件的标题
	temail.HTML = fmt.Sprintf(`<html>
		<head>
		</head>
			 <body>
			   <h1>你好啊，%s，欢迎来参加这次的CTF比赛!</h1>
			   <br>
			   <h2>点击超链接即可完成激活 <a href="http://%s:%s/active/%s/%s" target="_brank">点我</a></h2>
			   <br>
			   <h2>Hack and have fun!</h2>
	     	</body>
	 	</html>`, username, gameurl, gameport, username, activestring) //指定邮件内容
	// 发送邮件
	err := temail.Send()
	if err != nil {
		beego.Error("邮件发送失败：", err)
		return
	}
}

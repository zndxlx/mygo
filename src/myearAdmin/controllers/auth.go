package controllers

import (
    //"github.com/astaxie/beego"
    //"github.com/astaxie/beego/validation"
    // "github.com/lisijie/goblog/models"
    // "github.com/lisijie/goblog/util"
    // "strconv"
    "encoding/json"
    "myearAdmin/models"
    //"myearAdmin/util"
    //"strings"
)

type AuthController struct {
    baseController
}

//LoginReq 登陆请求
type LoginReq struct {
    Name string `json:"managerName" valid:"Required"`
    Pwd  string `json:"managerPwd" valid:"Required"`
}

type LoginRsp struct {
    Token string
    Role  int
}

// Login 登录
func (c *AuthController) Login() {
    var loginRsp LoginRsp

    loginReq := LoginReq{}
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &loginReq); err != nil {
        c.RspErr(ReqParaErr, err.Error())
        return
    }

    if c.ValidRequest(&loginReq) == false {
        return
    }

    manager, err := models.GetMangerByName(loginReq.Name)
    if err != nil {
        c.RspErr(1, "用户名不存在")
        return
    }

    if manager.Pass != loginReq.Pwd {
        c.RspErr(2, "用户名或密码错误")
        return
    }

    // util.ULog.Info("managers=%+v, err=%v", managers, err)

    loginRsp.Token, err = generateToken(manager.Id, manager.Name, manager.Role)
    if err != nil {
        c.RspErr(3, "创建token失败")
        return
    }
    loginRsp.Role = manager.Role
    c.RspData(loginRsp)
    return
    // if c.ValidRequest(&loginReq) == false {
    //     return
    // }

    // valid := validation.Validation{}
    // b, err := valid.Valid(&loginReq)
    // if err != nil {
    //     c.RspErr(10000, err.Error())
    //     return
    // }

    // if !b {
    //     errMsgs := []string{}
    //     for _, err := range valid.Errors {
    //         errMsgs = append(errMsgs, err.Key+":"+err.Message)
    //     }
    //     c.RspErr(10000, strings.Join(errMsgs, "-"))
    //     return
    // }
    //校验

    // managerName := strings.TrimSpace(c.GetString("managerName"))
    // managerPwd := strings.TrimSpace(c.GetString("managerPwd"))
    // if managerName == "" || managerPwd == "" {
    //     c.RspErr(10000, "中文测试abc")
    // }

    // beego.Informational("this is informational")
    // beego.Debug("this is debug")
}

//退出登录
func (c *AuthController) Logout() {
    c.RspData(struct{}{})
    //beego.Informational("this is informational")
    //beego.Debug("this is debug")
}

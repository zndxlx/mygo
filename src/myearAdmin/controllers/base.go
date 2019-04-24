package controllers

import (
    "github.com/astaxie/beego"
    //"github.com/astaxie/beego/context"
    "github.com/astaxie/beego/validation"
    "strings"
    // "time"
    "errors"
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "myearAdmin/util"
    "strconv"
    "time"
)

const (
    //SecretKey jwt密钥
    SecretKey = "laden"
)

const (
    TokenAuthErr = 10000 //TokenAuthErr token错误
    PermisionErr = 10001 //PermisionErr 权限错误
    ReqParaErr   = 10002 //ReqParaErr  请求参数错误
)

type baseController struct {
    beego.Controller
    managerID   int64
    role        int
    managerName string
}

// Prepare 在路由前处理
func (c *baseController) Prepare() {
    c.auth()
    //controller.checkPermission()
}

// auth 登录状态验证
func (c *baseController) auth() {
    controllerName, actionName := c.GetControllerAndAction()
    util.ULog.Info("util controllerName=%s, actionName=%s", controllerName, actionName)

    if actionName != "Login" {
        token := c.Ctx.Input.Header("Authorization")
        var err error = nil
        util.ULog.Info("Request token=%v", token)
        c.managerID, c.managerName, c.role, err = parseToken(token)
        if err != nil {
            c.RspErr(TokenAuthErr, err.Error())
            c.StopRun()
            return
        }
        if c.role != 1 && controllerName == "ManagerController" {
            id := c.Ctx.Input.Param(":id")
            managerIDStr := strconv.FormatInt(c.managerID, 10)
            if managerIDStr != id {
                c.RspErr(PermisionErr, "只允许操作自己")
                fmt.Printf("managerIDStr=%s, id=%s\n", managerIDStr, id)
                c.StopRun()
                return
            }
        }
    }
}

//是否post提交
func (c *baseController) isPost() bool {
    return c.Ctx.Request.Method == "POST"
}

//获取用户IP地址
func (c *baseController) getClientIp() string {
    s := strings.Split(c.Ctx.Request.RemoteAddr, ":")
    return s[0]
}

// func (this *baseController) getTime() time.Time {
//     timezone, _ := strconv.ParseFloat(option.Get("timezone"), 64)
//     add := timezone * float64(time.Hour)
//     return time.Now().UTC().Add(time.Duration(add))
// }

//ValidRequest 校验请求参数
func (c *baseController) ValidRequest(obj interface{}) (b bool) {
    valid := validation.Validation{}
    b, err := valid.Valid(obj)
    if err != nil {
        c.RspErr(ReqParaErr, err.Error())
        return false
    }

    if !b {
        errMsgs := []string{}
        for _, err := range valid.Errors {
            errMsgs = append(errMsgs, err.Key+":"+err.Message)
        }
        c.RspErr(ReqParaErr, strings.Join(errMsgs, "-"))
        return false
    }

    return true
}

func (c *baseController) RspJSON(code int, data interface{}, Msg string) {
    c.Data["json"] = map[string]interface{}{"code": code, "msg": Msg, "data": data}
    c.ServeJSON()
    util.ULog.Info("reqUrl=%s, reqBody=%s, rspBody=%+v", c.Ctx.Input.URI(),
        c.Ctx.Input.RequestBody, c.Data["json"])
}

func (c *baseController) RspData(data interface{}) {
    c.RspJSON(0, data, "success")
}

func (c *baseController) RspErr(code int, msg string) {
    c.RspJSON(code, struct{}{}, msg)
}

func parseToken(tokenString string) (managerID int64, managerName string, role int, err error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(SecretKey), nil
    })
    if err != nil {
        return
    }
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        err = errors.New("token is invalid")
        return
    }

    role = int(claims["role"].(float64))
    managerID = int64(claims["mID"].(float64))
    managerName = claims["mName"].(string)

    return
}

func generateToken(managerID int64, managerName string, role int) (token string, err error) {
    exp := time.Now().UTC().Add(6000 * time.Second).Unix() //60秒过期
    tokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "exp":   exp,
        "mID":   managerID,
        "mName": managerName,
        "role":  role,
    })
    token, err = tokenJwt.SignedString([]byte(SecretKey))
    return
}

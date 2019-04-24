package control

type rsp map[string]interface{}

func NewRsp(status int) rsp {
    r := rsp{}
    r["status"] = status
    return r
}

func (self rsp) GetStatus() (status int) {
    return self["status"].(int)
}

func (self rsp) SetStatus(status int) {
    self["status"] = status
}

func (self rsp) SetMsg(msg string) {
    self["msg"] = msg
}

func (self rsp) SetData(data interface{}) {
    self["data"] = data
}

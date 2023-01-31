package controller

const (
	CodeSuccess int32 = iota
	CodeInvalidParams
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodeInvalidToken
	CodeInvalidAuthFormat
	CodeNotLogin
	CodeInsertError
	CodeUpdateError
	CodeQueryError
	CodeDeleteError
	CodeGenerateTokenError
)

var MsgFlags = map[int32]string{
	CodeSuccess:            "success",
	CodeInvalidParams:      "请求参数错误",
	CodeUserExist:          "用户名重复",
	CodeUserNotExist:       "用户不存在",
	CodeInvalidPassword:    "用户名或密码错误",
	CodeServerBusy:         "服务繁忙",
	CodeInvalidToken:       "无效的Token",
	CodeInvalidAuthFormat:  "认证格式有误",
	CodeNotLogin:           "未登录",
	CodeInsertError:        "插入数据错误",
	CodeUpdateError:        "更新数据错误",
	CodeQueryError:         "查询数据错误",
	CodeDeleteError:        "删除数据错误",
	CodeGenerateTokenError: "生成token错误",
}

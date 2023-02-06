package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"tiktok-demo/dao/mysql"
	"tiktok-demo/logger"
	"tiktok-demo/middleware/jwt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//WS消息传送:客1 -> 服 -> 客2
//

// ws客户端发送给服务端的消息
// 客1 -> 服
type C2SMsgData struct {
	UserId     int64  `json:"user_id"`
	ToUserId   int64  `json:"to_user_id"`
	MsgContent string `json:"msg_content"`
}

// ws服务推送数据格式
// 服 -> 客2
type S2CMsgData struct {
	FromUserId int64  `json:"from_user_id"`
	MsgContent string `json:"msg_content"`
}

// 客-服 socket连接类
type SocketConn struct {
	UserId   int64 //消息发送方
	ToUserId int64 //消息接收方
	Conn     *websocket.Conn
	MsgChan  chan []byte //消息通道
}

// 信使类
// 服务端保存的,用于尝试转发给客2的,以及记入消息记录数据库的 中间类型
type Messeger struct {
	Conn    *SocketConn // 客1->服
	Message []byte
	Type    int
}

// socket连接的索引
type ConnKey struct {
	SourceId int64
	TargetId int64
}

// socket连接管理中心
// 管理时,使用select监听这三个管道.如果有来者,就做相应的动作.
type ConnManager struct {
	WsConns       map[ConnKey]*SocketConn //连接注册表
	MessegerChan  chan *Messeger          //信使等待处理的队列
	SetConnection chan *SocketConn        //建立连接,注册到管理中心
	DisConnection chan *SocketConn        //断开连接,从管理中心删除
}

var Manager = ConnManager{
	WsConns:       make(map[ConnKey]*SocketConn),
	MessegerChan:  make(chan *Messeger),
	SetConnection: make(chan *SocketConn),
	DisConnection: make(chan *SocketConn),
}

//gin 中间件无法在websocket中使用,只能另写专用中间件.
//gin response无法传回websocket客户端,只能修改响应头.
func WsChatHandler(c *gin.Context) {
	token := c.Query("token")
	jwtUid, err := jwt.WsAuthInHeader(token)
	if err != nil {
		c.Header("Error-code", fmt.Sprint(WsInvalidToken))
		c.Header("Error-type", WsMsgFlags[WsInvalidToken])
		return
	}
	uidRaw := c.Query("uid")
	toUidRaw := c.Query("toUid")
	uid, err := strconv.Atoi(uidRaw)
	if err != nil {
		c.Header("Error-code", fmt.Sprint(WsInvalidUid))
		c.Header("Error-type", WsMsgFlags[WsInvalidUid])
		return
	}
	if jwtUid != uid {
		c.Header("Error-code", fmt.Sprint(WsUidMismatch))
		c.Header("Error-type", WsMsgFlags[WsUidMismatch])
		return
	}
	toUid, err := strconv.Atoi(toUidRaw)
	if err != nil {
		c.Header("Error-code", fmt.Sprint(WsInvalidToUid))
		c.Header("Error-type", WsMsgFlags[WsInvalidToUid])
		return
	}
	_,err=mysql.CheckUserExist(toUid)
	if err != nil {
		c.Header("Error-code", fmt.Sprint(WsToUidNotExist))
		c.Header("Error-type", WsMsgFlags[WsToUidNotExist])
		return
	}
	//升级为websocket连接
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Header("Error-code", fmt.Sprint(WsUpgradeFail))
		c.Header("Error-type", WsMsgFlags[WsUpgradeFail])
		http.NotFound(c.Writer, c.Request)
		return
	}
	//创建一个连接实例,用于建立连接
	wsConn := &SocketConn{
		UserId:   int64(uid),
		ToUserId: int64(toUid),
		Conn:     conn,
		MsgChan:  make(chan []byte),
	}
	// 新连接注册到管理中心
	Manager.SetConnection <- wsConn
	go wsConn.TCP2Messeger()
	go wsConn.Messege2TCP()
}

// 如果服务端接收到了ws的tcp报文,就把它做成一个信使,丢到信使队列
func (c *SocketConn) TCP2Messeger() {
	defer func() {
		Manager.DisConnection <- c
		_ = c.Conn.Close()
	}()
	for {
		//读取tcp报文的json数据体
		c.Conn.PongHandler()
		msgData := new(C2SMsgData)
		err := c.Conn.ReadJSON(&msgData)
		if err != nil {
			logger.Log.Sugar().Errorf("tcp报文json数据格式不正确, %v", err)
			Manager.DisConnection <- c
			_ = c.Conn.Close()
			break
		}
		//把tcp报文做成信使
		Manager.MessegerChan <- &Messeger{
			Conn:    c,
			Message: []byte(msgData.MsgContent),
		}
	}
}

// 服务端 如果c.MsgChan有消息,就把消息发给客户端
func (c *SocketConn) Messege2TCP() {
	defer func() {
		_ = c.Conn.Close()
	}()
	for {
		select {
		case msgContent, ok := <-c.MsgChan:
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			msgData := S2CMsgData{
				FromUserId: c.UserId,
				MsgContent: string(msgContent),
			}
			msg, _ := json.Marshal(msgData)
			_ = c.Conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

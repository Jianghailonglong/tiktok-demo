package controller

import (
	"encoding/json"
	"fmt"
	"tiktok-demo/logger"
	"tiktok-demo/service"

	"github.com/gorilla/websocket"
)

func (manager *ConnManager) Start() {
	for {
		logger.Log.Info("<-----监听管道通信----->")
		select {
		//建立连接,需要自动给返回未读消息+10条已读消息
		case wsConn := <-manager.SetConnection:
			connKey := ConnKey{
				SourceId: wsConn.UserId,
				TargetId: wsConn.ToUserId,
			}
			logger.Log.Info(fmt.Sprintf(
				"有新的连接建立:[%v->%v]",
				connKey.SourceId, connKey.TargetId))
			manager.WsConns[connKey] = wsConn //注册到manager上
			msgData := &S2CMsgData{
				FromUserId: 0,
				MsgContent: "已连接至服务器",
			}
			msg, _ := json.Marshal(msgData)
			_ = wsConn.Conn.WriteMessage(websocket.TextMessage, msg)
			//返回未读消息+20条历史消息
			msgRecords, err := service.UnreadAndHistory(int(wsConn.UserId), int(wsConn.ToUserId))
			if err != nil {
				logger.Log.Sugar().Errorf("UnreadAndHistory失败:%v", err)
			}
			for _, mr := range msgRecords {
				msgData := S2CMsgData{
					FromUserId: int64(mr.SourceId),
					MsgContent: mr.Content,
				}
				msg, _ := json.Marshal(msgData)
				_ = wsConn.Conn.WriteMessage(websocket.TextMessage, msg)
			}
		//注销连接
		case wsConn := <-manager.DisConnection:
			connKey := ConnKey{
				SourceId: wsConn.UserId,
				TargetId: wsConn.ToUserId,
			}
			logger.Log.Info(fmt.Sprintf(
				"本连接即将注销:[%v->%v]",
				connKey.SourceId, connKey.TargetId))
			if _, ok := manager.WsConns[connKey]; ok {
				msgData := &S2CMsgData{
					FromUserId: 0,
					MsgContent: "连接已断开",
				}
				msg, _ := json.Marshal(msgData)
				_ = wsConn.Conn.WriteMessage(websocket.TextMessage, msg)
				close(wsConn.MsgChan)
				delete(manager.WsConns, connKey)
			}
		//发送消息
		case messenger := <-manager.MessegerChan:
			msgContent := messenger.Message
			//查找对方的连接是否存在
			connKey := ConnKey{
				SourceId: messenger.Conn.ToUserId,
				TargetId: messenger.Conn.UserId,
			}
			flag := false //指明对方是否在线,默认不在线
			for ck, wsConn := range manager.WsConns {
				if ck != connKey {
					continue
				}
				//没按他的来
				wsConn.MsgChan <- msgContent
				flag = true
			}
			//插入ws消息记录表
			if flag {
				msgData := &S2CMsgData{
					FromUserId: 0,
					MsgContent: "对方在线",
				}
				msg, _ := json.Marshal(msgData)
				_ = messenger.Conn.Conn.WriteMessage(websocket.TextMessage, msg)
				err := service.InsertWsMsg(
					int(messenger.Conn.UserId),
					int(messenger.Conn.ToUserId),
					string(msgContent), 1)
				if err != nil {
					logger.Log.Sugar().Errorf("InsertWsMsg 失败: %v", err)
				}
			} else {
				msgData := &S2CMsgData{
					FromUserId: 0,
					MsgContent: "对方离线",
				}
				msg, _ := json.Marshal(msgData)
				_ = messenger.Conn.Conn.WriteMessage(websocket.TextMessage, msg)
				err := service.InsertWsMsg(
					int(messenger.Conn.UserId),
					int(messenger.Conn.ToUserId),
					string(msgContent), 0)
				if err != nil {
					logger.Log.Sugar().Errorf("InsertWsMsg 失败: %v", err)
				}
			}

		}
	}
}

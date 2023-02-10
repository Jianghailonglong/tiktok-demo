package controller

const (
	WsInvalidToken int32 = iota
	WsUidMismatch
	WsInvalidUid
	WsInvalidToUid
	WsUpgradeFail
	WsToUidNotExist
)

var WsMsgFlags = map[int32]string{
	WsInvalidToken:  "Token is wronge.",
	WsUidMismatch:   "Token mismatch with userId.",
	WsInvalidUid:    "UserId should be a integer.",
	WsInvalidToUid:  "ToUserId should be a integer.",
	WsUpgradeFail:   "Http failed to upgrade to websocket.",
	WsToUidNotExist: "Don't have a User whose uid is ToUserId. ",
}

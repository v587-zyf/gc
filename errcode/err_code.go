package errcode

import (
	"github.com/v587-zyf/gc/enums"
)

var (
	ERR_SUCCEED      = CreateErrCode(0, NewCodeLang("成功", enums.LANG_CN), NewCodeLang("succeed", enums.LANG_EN))
	ERR_STANDARD_ERR = CreateErrCode(1, NewCodeLang("失败", enums.LANG_CN), NewCodeLang("failed", enums.LANG_EN))
	ERR_SIGN         = CreateErrCode(2, NewCodeLang("验证未通过", enums.LANG_CN), NewCodeLang("Verification failed", enums.LANG_EN))
	ERR_PARAM        = CreateErrCode(3, NewCodeLang("参数错误", enums.LANG_CN), NewCodeLang("The parameter is incorrect", enums.LANG_EN))
	ERR_CONFIG_NIL   = CreateErrCode(4, NewCodeLang("配置为空", enums.LANG_CN), NewCodeLang("The Config Is Nil", enums.LANG_EN))

	ERR_NET_SEND_TIMEOUT   = CreateErrCode(11, NewCodeLang("发送数据超时", enums.LANG_CN), NewCodeLang("The sending data timed out", enums.LANG_EN))
	ERR_NET_PKG_LEN_LIMIT  = CreateErrCode(12, NewCodeLang("数据包长度限制", enums.LANG_CN), NewCodeLang("Packet length limit", enums.LANG_EN))
	ERR_SERVER_INTERNAL    = CreateErrCode(13, NewCodeLang("服务器内部错误", enums.LANG_CN), NewCodeLang("Server internal error", enums.LANG_EN))
	ERR_WP_TOO_MANY_WORKER = CreateErrCode(14, NewCodeLang("工作池任务太多", enums.LANG_CN), NewCodeLang("There are too many work pool tasks", enums.LANG_EN))
	ERR_JSON_MARSHAL_ERR   = CreateErrCode(15, NewCodeLang("json打包错误", enums.LANG_CN), NewCodeLang("JSON packaging error", enums.LANG_EN))
	ERR_JSON_UNMARSHAL_ERR = CreateErrCode(16, NewCodeLang("json解包错误", enums.LANG_CN), NewCodeLang("JSON unpacking error", enums.LANG_EN))

	ERR_EVENT_PARAM_INVALID     = CreateErrCode(31, NewCodeLang("事件参数错误", enums.LANG_CN), NewCodeLang("Event parameter error", enums.LANG_EN))
	ERR_EVENT_LISTENER_LIMIT    = CreateErrCode(32, NewCodeLang("事件监听器数量限制", enums.LANG_CN), NewCodeLang("Event listener limit", enums.LANG_EN))
	ERR_EVENT_LISTENER_EMPTY    = CreateErrCode(33, NewCodeLang("事件监听器为空", enums.LANG_CN), NewCodeLang("Event listener is empty", enums.LANG_EN))
	ERR_EVENT_LISTENER_NOT_FIND = CreateErrCode(34, NewCodeLang("事件监听器未找到", enums.LANG_CN), NewCodeLang("Event listener not found", enums.LANG_EN))

	ERR_MQ_CONNECT_FAIL = CreateErrCode(41, NewCodeLang("mq连接失败", enums.LANG_CN), NewCodeLang("MQ connection failed", enums.LANG_EN))

	ERR_USER_DATA_NOT_FOUND  = CreateErrCode(101, NewCodeLang("用户信息未找到", enums.LANG_CN), NewCodeLang("User information not found", enums.LANG_EN))
	ERR_USER_DATA_INVALID    = CreateErrCode(102, NewCodeLang("用户信息错误", enums.LANG_CN), NewCodeLang("The user information is incorrect", enums.LANG_EN))
	ERR_REDIS_UPDATE_USER    = CreateErrCode(103, NewCodeLang("redis更新玩家数据错误", enums.LANG_CN), NewCodeLang("Redis update player data error", enums.LANG_EN))
	ERR_REDIS_LOGIN_DATA_NIL = CreateErrCode(104, NewCodeLang("redis登陆数据数据为空", enums.LANG_CN), NewCodeLang("The Redis login data is empty", enums.LANG_EN))
	ERR_MONGO_UPSERT         = CreateErrCode(105, NewCodeLang("upsert错误", enums.LANG_CN), NewCodeLang("upsert error", enums.LANG_EN))
	ERR_MONGO_FIND           = CreateErrCode(106, NewCodeLang("未找到数据", enums.LANG_CN), NewCodeLang("Data not found", enums.LANG_EN))
)

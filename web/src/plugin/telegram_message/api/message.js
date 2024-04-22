import service from '@/utils/request'

// SendMessage
// @Tags Telegram_bot
// @Summary 指定chat_id发送消息
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param    token chat_id content  "发送消息必须的参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /telegram_bot/sendMessage [post]
export const sendMessage = (data) => {
  return service({
    url: '/telegram_bot/sendMessage',
    method: 'post',
    data
  })
}

// IsMember
// @Tags Telegram_bot
// @Summary 查询用户是否为频道用户
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param    token user_id channel_id  "查询用户是否为频道用户必须的参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /telegram_bot/isMember [post]
export const isMember = (data) => {
  return service({
    url: '/telegram_bot/isMember',
    method: 'post',
    data
  })
}


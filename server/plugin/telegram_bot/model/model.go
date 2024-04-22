package model

type TelegramChat struct {
	Token       string `json:"token"`        // Telegram的Bot的Token 可单一可多个，多个token以英文逗号分隔
	ChatId      string `json:"chat_id"`      // 发送的目标ID
	Content     string `json:"content"`      // 发送的内容
	MessageType string `json:"message_type"` // 发送内容的文本类型 html 或 markdown 为空则默认为 markdown
}

type MemberCheck struct {
	Token     string `json:"token"`      // Telegram的Bot的Token 可单一可多个，多个token以英文逗号分隔
	UserId    string `json:"user_id"`    // 需要查询的用户id
	ChannelID string `json:"channel_id"` // 被查询的频道id，该频道需要已加入bot并给予了管理员权限
}

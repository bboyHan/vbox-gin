## Telegram_Bot插件

开发者：spiritlhl

版本：v0.0.4

### 使用步骤

#### 1. 注册插件

查看 ```server/initialize/plugin.go``` 文件中是否已注册插件，如若未注册，在```import```中增加

```
"github.com/flipped-aurora/gin-vue-admin/server/plugin/telegram_bot"
```

然后在函数

```
func InstallPlugin
```

的尾部插入

```
    PluginInit(PrivateGroup, telegram_bot.CreateTelegram_botPlug())
```

### 2. 消息发送

#### 入参结构说明

```
type TelegramChat struct {
	Token       string `json:"token"`        // Telegram的Bot的Token 可单一可多个，多个token以英文逗号分隔
	ChatId      string `json:"chat_id"`      // 发送的目标ID
	Content     string `json:"content"`      // 发送的内容
	MessageType string `json:"message_type"` // 发送内容的文本类型 html 或 markdown 为空则默认为 纯文本
}
```

#### 可直接调用的接口

发送消息的接口： /telegram_bot/sendMessage [post]

### 3. 是否为频道用户

#### 入参结构说明

```
type MemberCheck struct {
	Token     string `json:"token"`      // Telegram的Bot的Token 可单一可多个，多个token以英文逗号分隔
	UserId    string `json:"user_id"`    // 需要查询的用户id
	ChannelID string `json:"channel_id"` // 被查询的频道id，该频道需要已加入bot并给予了管理员权限
}
```

#### 可直接调用的接口

发送消息的接口： /telegram_bot/isMember [post]

### 4. 后台设置

![image](https://github.com/spiritysdx/tgm/assets/97792170/5dd643a3-4b42-4d92-9dd1-1cefd72e5ebc)

路径：```plugin/telegram_message/view/index.vue```

![image](https://github.com/spiritysdx/tgm/assets/97792170/7b40a5ec-78a5-47b6-9cfa-d4a219578369)

路径：```/telegram_bot/sendMessage```

![image](https://github.com/spiritysdx/tgm/assets/97792170/be2479fe-c51a-40c6-9a3c-569258a3a8b2)

路径：```/telegram_bot/isMember```

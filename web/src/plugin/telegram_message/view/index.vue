<template>
  <div>
    <warning-bar title="需要提前在Telegram申请Bot的Token，务必注意Token不要泄露给任何人，token可使用一个或多个" />
    <div class="gva-form-box">
      <el-form
        ref="telegramMessageForm"
        label-position="right"
        label-width="80px"
        :model="form"
      > 
        <el-form-item label="token">
          <el-input v-model="form.token" />
        </el-form-item>
        <el-form-item label="chat_id">
          <el-input v-model="form.chat_id" />
        </el-form-item>
        <el-form-item label="消息内容">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="6"
          />
        </el-form-item>
        <el-form-item label="消息类型">
          <el-select v-model="form.message_type" placeholder="消息类型">
            <el-option label="Markdown" value="markdown"></el-option>
            <el-option label="HTML" value="html"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button @click="sendTelegramMessage">发送消息</el-button>
        </el-form-item>
      </el-form>
      <el-form
        ref="isMemberForm"
        label-position="right"
        label-width="80px"
        :model="memberForm"
      >
        <el-form-item label="token">
          <el-input v-model="memberForm.token" />
        </el-form-item>
        <el-form-item label="user_id">
          <el-input v-model="memberForm.user_id" />
        </el-form-item>
        <el-form-item label="channel_id">
          <el-input v-model="memberForm.channel_id" />
        </el-form-item>
        <el-form-item>
          <el-button @click="checkMembership">查询是否为频道用户</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import WarningBar from '@/components/warningBar/warningBar.vue'
import { sendMessage, isMember } from '@/plugin/telegram_message/api/message.js'
import { ElMessage } from 'element-plus'
import { reactive, ref } from 'vue'

defineOptions({
  name: 'TelegramMessage',
})

const telegramMessageForm = ref(null)
const form = reactive({
  token: '可单一可多个，多个token时以英文逗号分隔',
  chat_id: '',
  content: 'HTML示例：<code>1234</code>\nMarkdown示例：*1234*',
  message_type: 'markdown',
})
const sendTelegramMessage = async() => {
  const res = await sendMessage(form)
  if (res.code === 0) {
    ElMessage.success('发送成功,请查收')
  }
}
const memberForm = reactive({
  token: '可单一可多个，多个token时以英文逗号分隔',
  user_id: '',
  channel_id: ''
})
const checkMembership = async () => {
  const res = await isMember(memberForm);
  if (res.code === 0) {
    ElMessage.success('查询成功');
    ElMessage.success(`是否为频道用户：${res.data}`);
  } else {
    ElMessage.error('查询失败');
  }
}
</script>

<style scoped>
/* 样式可以根据您的需求进行调整 */
.gva-form-box {
  max-width: 600px;
  margin: 0 auto;
}
</style>

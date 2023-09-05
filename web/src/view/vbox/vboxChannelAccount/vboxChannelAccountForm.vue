<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="归属用户ID:" prop="uid">
          <el-input v-model.number="formData.uid" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="通道账户:" prop="acAccount">
          <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="账户密码:" prop="acPwd">
          <el-input v-model="formData.acPwd" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="账户备注:" prop="acRemark">
          <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="账户id:" prop="acId">
          <el-input v-model.number="formData.acId" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="通道id:" prop="cid">
          <el-input v-model.number="formData.cid" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="token:" prop="token">
          <el-input v-model="formData.token" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="日限额:" prop="dailyLimit">
          <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="总限额:" prop="totalLimit">
          <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="状态开关:" prop="status">
          <el-input v-model.number="formData.status" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="系统开关:" prop="sysStatus">
          <el-input v-model.number="formData.sysStatus" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ChannelAccount'
}
</script>

<script setup>
import {
  createChannelAccount,
  updateChannelAccount,
  findChannelAccount
} from '@/api/vboxChannelAccount'

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'
// 富文本组件
// import RichEdit from '@/components/richtext/rich-edit.vue'

const route = useRoute()
const router = useRouter()

const type = ref('')
const formData = ref({
            uid: 0,
            acAccount: '',
            acPwd: '',
            acRemark: '',
            token: '',
            acId: 0,
            cid: 0,
            dailyLimit: 0,
            totalLimit: 0,
            status: 0,
            sysStatus: 0,
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findChannelAccount({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.revca
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
}

init()
// 保存按钮
const save = async() => {
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return
            let res
           switch (type.value) {
             case 'create':
               res = await createChannelAccount(formData.value)
               break
             case 'update':
               res = await updateChannelAccount(formData.value)
               break
             default:
               res = await createChannelAccount(formData.value)
               break
           }
           if (res.code === 0) {
             ElMessage({
               type: 'success',
               message: '创建/更改成功'
             })
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>

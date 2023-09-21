<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="团队ID:" prop="teamId">
          <el-input v-model="formData.teamId" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="团队名:" prop="teamName">
          <el-input v-model="formData.teamName" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="用户ID:" prop="uid">
          <el-input v-model="formData.uid" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="团队领导角色名ID:" prop="leaderId">
          <el-input v-model="formData.leaderId" :clearable="true" placeholder="请输入" />
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
  name: 'VboxTeamsUser'
}
</script>

<script setup>
import {
  createVboxTeamsUser,
  updateVboxTeamsUser,
  findVboxTeamsUser
} from '@/api/vboxTeamsUser'

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'

const route = useRoute()
const router = useRouter()

const type = ref('')
const formData = ref({
            teamId: '',
            teamName: '',
            uid: '',
            leaderId: '',
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findVboxTeamsUser({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.retUsers
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
               res = await createVboxTeamsUser(formData.value)
               break
             case 'update':
               res = await updateVboxTeamsUser(formData.value)
               break
             default:
               res = await createVboxTeamsUser(formData.value)
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

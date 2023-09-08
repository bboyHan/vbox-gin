<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="uid:" prop="uid">
          <el-input v-model.number="formData.uid" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="用户名:" prop="username">
          <el-input v-model="formData.username" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="充值数:" prop="recharge">
          <el-input v-model.number="formData.recharge" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="费率:" prop="tariff">
          <el-input-number v-model="formData.tariff" :precision="2" :clearable="true"></el-input-number>
       </el-form-item>
        <el-form-item label="划转/分配:" prop="remark">
          <el-input v-model="formData.remark" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="创建时间:" prop="createTime">
          <el-date-picker v-model="formData.createTime" type="date" placeholder="选择日期" :clearable="true"></el-date-picker>
       </el-form-item>
        <el-form-item label="会员等级:" prop="vipLevel">
          <el-input v-model.number="formData.vipLevel" :clearable="true" placeholder="请输入" />
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
  name: 'VboxUserWallet'
}
</script>

<script setup>
import {
  createVboxUserWallet,
  updateVboxUserWallet,
  findVboxUserWallet
} from '@/api/vboxUserWallet'

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'

const route = useRoute()
const router = useRouter()

const type = ref('')
const formData = ref({
            uid: 0,
            username: '',
            recharge: 0,
            tariff: 0,
            remark: '',
            createTime: new Date(),
            vipLevel: 0,
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findVboxUserWallet({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.revuw
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
               res = await createVboxUserWallet(formData.value)
               break
             case 'update':
               res = await updateVboxUserWallet(formData.value)
               break
             default:
               res = await createVboxUserWallet(formData.value)
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

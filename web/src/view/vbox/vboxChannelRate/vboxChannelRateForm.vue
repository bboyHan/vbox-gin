<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="uid:" prop="uid">
          <el-input v-model.number="formData.uid" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="通道编码ID:" prop="channelCode">
          <el-input v-model="formData.channelCode" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="产品名称:" prop="productName">
          <el-input v-model="formData.productName" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="产品ID:" prop="productId">
          <el-input v-model="formData.productId" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="费率:" prop="rate">
          <el-input-number v-model="formData.rate" :precision="2" :clearable="true"></el-input-number>
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
  name: 'VboxChannelRate'
}
</script>

<script setup>
import {
  createVboxChannelRate,
  updateVboxChannelRate,
  findVboxChannelRate
} from '@/api/vboxChannelRate'

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
            channelCode: '',
            productName: '',
            productId: '',
            rate: 0,
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findVboxChannelRate({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.rechRate
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
               res = await createVboxChannelRate(formData.value)
               break
             case 'update':
               res = await updateVboxChannelRate(formData.value)
               break
             default:
               res = await createVboxChannelRate(formData.value)
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

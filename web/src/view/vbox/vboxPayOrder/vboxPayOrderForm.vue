<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="订单id:" prop="order_id">
          <el-input v-model="formData.order_id" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="付方uuid:" prop="p_account">
          <el-input v-model="formData.p_account" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="金额:" prop="cost">
          <el-input v-model.number="formData.cost" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="uid:" prop="uid">
          <el-input v-model.number="formData.uid" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="帐号id:" prop="ac_id">
          <el-input v-model="formData.ac_id" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="所属通道:" prop="c_channel_id">
          <el-input v-model="formData.c_channel_id" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="平台id:" prop="platform_oid">
          <el-input v-model="formData.platform_oid" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="客户ip:" prop="pay_ip">
          <el-input v-model="formData.pay_ip" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="区域:" prop="pay_region">
          <el-input v-model="formData.pay_region" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="支付链接:" prop="resource_url">
          <RichEdit v-model="formData.resource_url"/>
       </el-form-item>
        <el-form-item label="回调地址:" prop="notify_url">
          <el-input v-model="formData.notify_url" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="订单状态:" prop="order_status">
          <el-input v-model.number="formData.order_status" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="回调状态:" prop="callback_status">
          <el-input v-model.number="formData.callback_status" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="取码状态:" prop="code_use_status">
          <el-input v-model.number="formData.code_use_status" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="创建时间:" prop="create_time">
          <el-date-picker v-model="formData.create_time" type="date" placeholder="选择日期" :clearable="true"></el-date-picker>
       </el-form-item>
        <el-form-item label="异步执行时间:" prop="async_time">
          <el-date-picker v-model="formData.async_time" type="date" placeholder="选择日期" :clearable="true"></el-date-picker>
       </el-form-item>
        <el-form-item label="回调时间:" prop="call_time">
          <el-date-picker v-model="formData.call_time" type="date" placeholder="选择日期" :clearable="true"></el-date-picker>
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
  name: 'VboxPayOrder'
}
</script>

<script setup>
import {
  createVboxPayOrder,
  updateVboxPayOrder,
  findVboxPayOrder
} from '@/api/vboxPayOrder'

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'
// 富文本组件
import RichEdit from '@/components/richtext/rich-edit.vue'

const route = useRoute()
const router = useRouter()

const type = ref('')
const formData = ref({
            order_id: '',
            p_account: '',
            cost: 0,
            uid: 0,
            ac_id: '',
            c_channel_id: '',
            platform_oid: '',
            pay_ip: '',
            pay_region: '',
            notify_url: '',
            order_status: 0,
            callback_status: 0,
            code_use_status: 0,
            create_time: new Date(),
            async_time: new Date(),
            call_time: new Date(),
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findVboxPayOrder({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.revpo
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
               res = await createVboxPayOrder(formData.value)
               break
             case 'update':
               res = await updateVboxPayOrder(formData.value)
               break
             default:
               res = await createVboxPayOrder(formData.value)
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

<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="通道id:" prop="c_channel_id">
          <el-input v-model="formData.c_channel_id" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="文件base 64编码:" prop="img_base_str">
          <el-input v-model="formData.img_base_str" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="图片顺序:" prop="img_num">
          <el-input v-model.number="formData.img_num" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="文件名:" prop="file_name">
          <el-input v-model="formData.file_name" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="图片地址:" prop="url">
          <el-input v-model="formData.url" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="文件标签:" prop="tag">
          <el-input v-model="formData.tag" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="编号:" prop="key">
          <el-input v-model="formData.key" :clearable="true" placeholder="请输入" />
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
  name: 'Channel_guideimg'
}
</script>

<script setup>
import {
  createChannel_guideimg,
  updateChannel_guideimg,
  findChannel_guideimg
} from '@/api/channelGuideImg'

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'

const route = useRoute()
const router = useRouter()

const type = ref('')
const formData = ref({
            c_channel_id: '',
            img_base_str: '',
            img_num: 0,
            file_name: '',
            url: '',
            tag: '',
            key: '',
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findChannel_guideimg({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.rechGuideImg
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
               res = await createChannel_guideimg(formData.value)
               break
             case 'update':
               res = await updateChannel_guideimg(formData.value)
               break
             default:
               res = await createChannel_guideimg(formData.value)
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

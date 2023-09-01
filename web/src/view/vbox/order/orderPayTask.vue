<template>
    <div>
        <h4>订单支付页面</h4>
    <warning-bar
        title="订单支付页面测试中"
      />

    </div>
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog"  destroy-on-close class="el-dialog__wrapper" fullscreen>
        <h2 class="dialog-title">👇👇👇操作流程提示👇👇👇</h2>
        <div >
            <el-image :src="imgData.img_base_str" fit="contain" class="thumbnail-image"/>
        </div>
        <!-- <template #footer> -->
        <div class="dialog-footer">
          <el-button @click="changImgPrev">上一步</el-button>
          <el-button @click="changImgNext">下一步</el-button>
          <el-button type="primary" @click="enterDialog">我知道了</el-button>
        </div>
      <!-- </template> -->
    </el-dialog>
    
</template>

<script>
export default {
  name: 'OrderPayTask'
}
</script>

<script setup>
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive,watch } from 'vue'

import {
  getChannel_guideimgList
} from '@/api/channelGuideImg'


const imgData = ref({
    c_channel_id: '',
    img_base_str: '',
    img_num: 0
    })
const page = ref(1)
const total = ref(0)
const imgNum = ref(1)
const pageSize = ref(10)
const searchInfo = ref({})
const tableData = ref([])

const getTableData = async() => {
    const table = await getChannel_guideimgList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    if (table.code === 0) {
        tableData.value = table.data.list
        console.log('imgs=' + JSON.stringify(tableData.value))
        total.value = table.data.total
        imgData.value = tableData.value[imgNum.value - 1]
    }
}

getTableData()

const elFormRef = ref()
// 弹窗控制标记
const dialogFormVisible = ref(false)

const changImgPrev = () => {
    if (imgNum.value > 1){
        imgNum.value --
    }else {
        imgNum.value = 1
    }
    
    getTableData()
}
const changImgNext = () => {
    
    if (imgNum.value >= total.value){
        imgNum.value = total.value
    }else {
        imgNum.value ++
    }
    getTableData()
}

// 打开弹窗
const openDialog = () => {
    
    dialogFormVisible.value = true
}
openDialog()
// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    // formData.value = {
    //     c_channel_id: '',
    //     img_base_str: '',
    //     img_num: 0,
    //     file_name: '',
    //     url: '',
    //     tag: '',
    //     key: '',
    //     }
}
// 弹窗确定
const enterDialog = async () => {
    //  elFormRef.value?.validate( async (valid) => {
    //          if (!valid) return
            //   let res
            //   switch (type.value) {
            //     case 'create':
            //       res = await createChannel_guideimg(formData.value)
            //       break
            //     case 'update':
            //       res = await updateChannel_guideimg(formData.value)
            //       break
            //     default:
            //       res = await createChannel_guideimg(formData.value)
            //       break
            //   }
            //   if (res.code === 0) {
            //     ElMessage({
            //       type: 'success',
            //       message: '创建/更改成功'
            //     })
            //     closeDialog()
            //     getTableData()
            //   }
    //   })
      closeDialog()
}

</script>

<style>
.el-dialog__wrapper {
  background-color: transparent !important;
  display: flex;
  align-items: center;
  justify-content: center;
}
.thumbnail-image {
  /* position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%); */
  margin-bottom: 20px;
}
.dialog-footer {
  display: flex;
  width: 100%;
  justify-content: flex-end;
}
.dialog-title {
  color: red;
  text-align: center;
}
</style>
<template>
  <div class="container">
    <div class="content">
      <h1>欢迎来到网页</h1>
      <count-down
          v-model:fire="fire"
          :tiping="tiping"
          :tipend="tipend"
          time="60"
          @statusChange="onStatusChange"
          @end="onEnd"
      >
      </count-down>

      <div class="buttons">
        <el-row :gutter="12">
          <el-col :span="24">
            <el-button type="primary" class="button" round @click="dialogFormVisible = true">操作指南</el-button>
          </el-col>
          <el-col :span="24">
            <el-button type="success" class="button" round>跳转支付</el-button>
          </el-col>
        </el-row>
      </div>

      <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" destroy-on-close class="el-dialog__wrapper" fullscreen>
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
      </el-dialog>
    </div>
  </div>
</template>

<script setup>
import { defineProps } from 'vue';
import { ElButton } from 'element-plus';
import { onMounted, ref } from 'vue';
import CountDown from 'vue-canvas-countdown';
import {
  getChannelGuideImgTaskList
} from '@/api/channelGuideImg'

// 弹窗控制标记
const dialogFormVisible = ref(false)
const changImgPrev = () => {
  if (imgNum.value > 1){
    imgNum.value --
  }else {
    imgNum.value = 1
  }
}
const changImgNext = () => {
  if (imgNum.value >= total.value){
    imgNum.value = total.value
  }else {
    imgNum.value ++
  }
}
const imgData = ref({
  c_channel_id: '',
  img_base_str: '',
  img_num: 0
})
// const page = ref(1)
const total = ref(0)
const imgNum = ref(1)
// const pageSize = ref(10)
// const searchInfo = ref({})
const tableData = ref([])

const chId = ref("tx_jd")

const getTableData = async() => {
  const table = await getChannelGuideImgTaskList({ channelId: chId.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    console.log('imgs=' + JSON.stringify(tableData.value))
    total.value = table.data.total
    imgData.value = tableData.value[imgNum.value - 1]
  }
}

getTableData()

// 打开弹窗
const openDialog = () => {
  dialogFormVisible.value = true
}
openDialog()

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
}

// 弹窗确定
const enterDialog = async () => {
  closeDialog()
}

// ---------- 倒计时 ----------------
const fire = ref(0);
const tiping = {
  text: '倒计时进行中',
  color: '#fff'
};
const tipend = {
  text: '倒计时结束',
  color: '#fff'
};

const fireCD = () => {
  // 配置参数（更多配置如下表）
  tiping.text = '请支付';
  tiping.color = '#fff';
  tipend.text = '停止支付';
  tipend.color = '#fff';

  // 启动倒计时(效果如上图所示)
  fire.value++;
};

const onStatusChange = (payload) => {
  console.log('倒计时状态改变：', payload);
};

const onEnd = () => {
  console.log('倒计时结束的回调函数');
};

onMounted(() => {
  // 启动倒计时
  fireCD();
});

const props = defineProps({
  title: {
    type: String,
    default: '欢迎来到网页',
  },
});
</script>

<style scoped>
.container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.content {
  text-align: center;
  color: #333;
}

h1 {
  font-size: 36px;
  margin-bottom: 24px;
}

.buttons {
  display: flex;
  flex-direction: column;
  gap: 16px;
  align-items: center; /* 垂直居中对齐按钮 */
  margin-top: 20px;
}

.button {
  padding: 12px 24px;
  font-size: 18px;
  margin-top: 10px;
  width: 80%;
}

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

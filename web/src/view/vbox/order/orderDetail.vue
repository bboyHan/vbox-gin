<template>
  <div class="container">
    <div class="content">
      <div v-show="dialogCountVisible">
        <count-down
          :fire="fire"
          :tiping="tiping"
          :tipend="tipend"
          time="12"
          @statusChange="onStatusChange"
          @end="onEnd"
          :statusChange="[2000,500]"
          width="180"
          height="180"
        >
        </count-down>
        <div class="buttons">
          <el-row :gutter="6">
            <el-col :span="24">
              <el-button type="success" class="button" round>正在通过安全验证，请等待...</el-button>
            </el-col>
            <el-col :span="24">
              <el-button type="primary" class="button" round>订单正在匹配中，预计5-20秒</el-button>
            </el-col>
          </el-row>
        </div>
      </div>

      <div v-show="payVisible">
        <!-- 显示新的 div 的代码... -->
        <h1>付款页面</h1>
      </div>
      <div v-show="notFoundVisible">
        <!-- 显示新的 div 的代码... -->
        <h1>订单不存在</h1>
      </div>
      <div v-show="finishedVisible">
        <!-- 显示新的 div 的代码... -->
        <h1>已付款成功</h1>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ElButton } from 'element-plus';
import { onMounted, ref, onUnmounted } from 'vue';
import CountDown from 'vue-canvas-countdown';
import { queryOrderSimple } from '@/api/vboxPayOrder';
import { useRoute } from 'vue-router';

// 弹窗控制标记
const dialogCountVisible = ref(true)
const payVisible = ref(false)
const finishedVisible = ref(false)
const notFoundVisible = ref(false)
const route = useRoute()

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

const fireCD = async () => {
  // 配置参数（更多配置如下表）
  tiping.text = '匹配中';
  tiping.color = '#fff';
  tipend.text = '停止匹配';
  tipend.color = '#fff';

  // 启动倒计时(效果如上图所示)
  fire.value++;
};

const onStatusChange = async (payload) => {
  console.log('倒计时状态改变：', payload);
};

const onEnd = async () => {
  console.log('倒计时结束的回调函数');
};

// 添加一个空变量作为定时器的 ID
let timerId = null;
let count = 1;

onMounted(() => {
  // 启动倒计时
  fireCD();
  // 启动定时器，每秒钟请求一次 HTTP 接口
  timerId = setInterval(queryOrder, 2000);
});

onUnmounted(() => {
  // 组件销毁时清除定时器
  clearInterval(timerId);
});

const queryOrder = async () => {
  try {
    const orderId = route.query.orderId;
    console.log(orderId)
    const result = await queryOrderSimple({order_id: orderId}); // 发送 HTTP 请求
    console.log(result)
    console.log(result.code)
    if (result) {
      clearInterval(timerId); // 如果状态发生变化，则停止定时器
    }
    if (result.code === 404) {
      dialogCountVisible.value = false;
      notFoundVisible.value = true;
    } else if (result.code === 0) {
      // clearInterval(timerId); // 如果状态发生变化，则停止定时器
      dialogCountVisible.value = false;
      if (result.data.status === 1) {
        finishedVisible.value = true;
      }
      if (result.data.status === 2) {
        payVisible.value = true;
      }
    }
  } catch (error) {
    console.log(error);
  }
}


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
  margin-top: 30px;
}

.button {
  padding: 12px 24px;
  font-size: 18px;
  margin-top: 6px;
  width: 80%;
  height: 42px;
}
</style>

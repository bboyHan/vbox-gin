<template>
  <div>
    <div v-show="dialogCountVisible">
      <div class="c_container">
        <div class="c_content">
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
                <el-button type="success" class="c_button" round>正在通过安全验证，请等待...</el-button>
              </el-col>
              <el-col :span="24">
                <el-button type="primary" class="c_button" round>订单正在匹配中，预计5-20秒</el-button>
              </el-col>
            </el-row>
          </div>
        </div>
      </div>
    </div>

    <div v-show="payVisible">
      <!-- 显示新的 div 的代码... -->
      <div class="p_container">
        <!--        <div class="p_blue-section" v-for="(color, index) in blueColors" :key="index" :style="{ backgroundColor: color }"></div>-->
        <div class="p_blue-section" v-for="index in 10" :key="index" :style="{ backgroundColor: generateColor(index) }"></div>
        <div class="p_content">
          <el-row :gutter="12">
            <el-col>
              <img src="@/assets/logo.png" alt="" style="width: 80px; height: 80px">
            </el-col>
            <el-col>
              <div style="color: #6B7687; margin-top: 10px; font-size: 16px">无法充值或提示错误，请联系客服！</div>
            </el-col>
            <el-col>
              <div style="color: #6B7687; margin-top: 20px; font-size: 60px">￥{{ payData.money }}.00</div>
            </el-col>
            <el-col>
              <div style="color: #e81239; margin-top: 10px; font-size: 16px">
                <el-icon style="margin-right: 5px"><WarningFilled /></el-icon>请在规定时间内付款！
              </div>
            </el-col>
            <el-col>
              <div style="height: 100px;">
              </div>
            </el-col>
            <el-col :span="24">
              <!--                <el-button class="p_button" type="primary">复制充值账号 立即付款</el-button>-->
            </el-col>
          </el-row>
        </div>
      </div>
      <div class="p_content_button">
        <el-row :gutter="12">
          <el-col>
            <button class="p_button" @click="openPay(payData)">立即付款</button>
          </el-col>
        </el-row>
      </div>
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

</template>
<script>
export default {
  name: 'Pay',
}
</script>
<script setup>
import { ElButton } from 'element-plus';
import { onMounted, ref, onUnmounted } from 'vue';
import CountDown from 'vue-canvas-countdown';
import { queryOrderSimple } from '@/api/payOrder';
import { useRoute } from 'vue-router';
import { WarningFilled } from '@element-plus/icons-vue';

// 弹窗控制标记
const dialogCountVisible = ref(true)
const payVisible = ref(false)
const finishedVisible = ref(false)
const notFoundVisible = ref(false)
const route = useRoute()

// ---------- 付款页 ----------------

const generateColor = (index) => {
  const hue = 220; // 蓝色的色调，范围为0-360
  const saturation = 90; // 蓝色的饱和度，范围为0-100
  const lightness = 100 - ((index - 0) * 8); // 蓝色的亮度，范围为0-100
  return `hsl(${hue}, ${saturation}%, ${lightness}%)`;
};

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

const openPay = async (payData) => {
  location.href = payData.resource_url;
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

const payData = ref({
  money: 0,
  resource_url: '',
  status: 0,
})

const queryOrder = async () => {
  try {
    const orderId = route.query.orderId;
    console.log(orderId)
    const result = await queryOrderSimple({order_id: orderId}); // 发送 HTTP 请求
    payData.value = result.data
    if (result) {
      clearInterval(timerId); // 如果状态发生变化，则停止定时器
    }
    if (result.code === 7) {
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
.c_container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.c_content {
  text-align: center;
  color: #333;
}

h1 {
  font-size: 36px;
  margin-bottom: 24px;
}


.c_button {
  padding: 12px 24px;
  font-size: 18px;
  margin-top: 6px;
  width: 80%;
  height: 42px;
}

.p_container {
  display: flex;
  justify-content: space-between;
  height: 45vh;
}

.p_blue-section {
  flex-grow: 1;
}

.p_content {
  text-align: center;
  color: #333;
  position: absolute;
  top: 5%;
  left: 5%;
  right: 5%;
  height: 45%;
  background-color: #f2f2f2;
  border-radius: 20px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.p_content_button {
  text-align: center;
  color: #333;
  position: absolute;
  top: 52%;
  left: 5%;
  right: 5%;
}

.p_button {
  border: none;
  padding: 12px 24px;
  font-size: 22px;
  color: #e7dfdf;
//background-color: #1660f3;
  background: linear-gradient(to right, #064954, #125280, #1247c9);
  margin-top: 6px;
  width: 80%;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  height: 50px;
}
</style>

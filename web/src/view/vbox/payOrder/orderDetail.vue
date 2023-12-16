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
      <!--   1000 引导   -->
      <div v-if="payTypeVisible === 1">
        <div class="p_container">
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
                  <div>
                    <span v-if="countdowns[0] > 0">{{ formatTime(countdowns[0]) }} </span>
                    <span v-else>-1 （已过期）</span>
                  </div>
                </div>
              </el-col>
              <el-col :span="24">

              </el-col>
            </el-row>
          </div>
        </div>
        <div class="p_content_inner" style="margin-top: 20px;">
          <el-row>
            <el-col>
              <div style="height: 5px; margin-top: 20px">
              </div>
            </el-col>
            <el-col>
              <div style="height: 100px; margin-top: 20px">
                <div class="medicine-bag">
                  <span>{{ payData.account }}</span>
                </div>
                <div class="copy-tip">
                  <span>长按框内</span>
                  <span class="jtone"></span>
                  <span>全选</span>
                  <span class="jttwo"></span>
                  <span>复制</span>
                  <span class="jtthree"></span>
                  <span>打开跳转</span>
                </div>
<!--                <button class="copy_button" @click="">一键复制</button>-->
              </div>
            </el-col>
          </el-row>
        </div>
        <div class="p_content_button">
          <el-row :gutter="12">
            <el-col>
              <button class="p_button" @click="openPay(payData)">点击付款</button>
            </el-col>
          </el-row>
        </div>
      </div>
      <!--   3000 直付   -->
      <div v-if="payTypeVisible === 2">
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
                  <div>
                    <span v-if="countdowns[0] > 0">{{ formatTime(countdowns[0]) }} </span>
                    <span v-else>-1 （已过期）</span>
                  </div>
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
        <div class="p_content_inner" style="margin-top: 20px;">
          <el-row>
            <el-col>
              <div style="height: 5px; margin-top: 5px">
              </div>
            </el-col>
            <el-col>
              <div>
                <div class="medicine-bag-qr">
                  <div v-if="qrcodeUrl">
                    <img :src="qrcodeUrl" alt="QR Code" style="height: 174px"/>
                  </div>
                  <div v-else>
                    <span>暂无二维码</span>
                  </div>
                </div>
                <div class="copy-tip">
                  <span>核对信息</span>
                  <span class="jtone"></span>
                  <span>打开微信</span>
                  <span class="jttwo"></span>
                  <span>扫一扫</span>
                  <span class="jtthree"></span>
                  <span>确认付款</span>
                </div>
                <!--                <button class="copy_button" @click="">一键复制</button>-->
              </div>
            </el-col>
          </el-row>
        </div>
        <div class="p_content_button_qr">
          <el-row :gutter="12">
            <el-col>
              <button class="p_button" @click="">扫码直付</button>
            </el-col>
          </el-row>
        </div>
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
    <div v-show="timeoutVisible">
      <!-- 显示新的 div 的代码... -->
      <h1>订单已超时</h1>
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
import {formatTime} from "@/utils/format";
import QRCode from "qrcode";

// 弹窗控制标记
const dialogCountVisible = ref(true)
const payVisible = ref(false)
const payTypeVisible = ref(0)
const finishedVisible = ref(false)
const timeoutVisible = ref(false)
const notFoundVisible = ref(false)
const route = useRoute()

// ---------- 付款页 ----------------

const generateColor = (index) => {
  const hue = 220; // 蓝色的色调，范围为0-360
  const saturation = 90; // 蓝色的饱和度，范围为0-100
  const lightness = 100 - ((index - 0) * 8); // 蓝色的亮度，范围为0-100
  return `hsl(${hue}, ${saturation}%, ${lightness}%)`;
};

// 直付扫码
const qrcodeUrl = ref('');

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
  timerId = setInterval(queryOrder, 1500);
  calculateCountdown();
});

onUnmounted(() => {
  // 组件销毁时清除定时器
  clearInterval(timerId);
});

const payData = ref({
  money: 0,
  expTime: 0,
  account: '',
  order_id: '',
  resource_url: '',
  status: 0,
})

const queryOrder = async () => {
  try {
    const orderId = route.query.orderId;
    console.log(orderId)
    const result = await queryOrderSimple({order_id: orderId}); // 发送 HTTP 请求
    payData.value = result.data
    const content = result.data.resource_url;
    if(content){
      clearInterval(timerId); // 如果状态发生变化，则停止定时器
      dialogCountVisible.value = false;

      if (result.code === 7) {
        dialogCountVisible.value = false;
        notFoundVisible.value = true;
      } else if (result.code === 0) {

        let cid = Number(result.data.channel_code);
        if (cid >= 1000 && cid < 1099) {
          payTypeVisible.value = 1;
        } else if (cid >= 2000 && cid < 2099){
          payTypeVisible.value = 1;
        } else if (cid >= 3000 && cid < 3099){
          payTypeVisible.value = 2;
          if (content){
            QRCode.toDataURL(content)
                .then((dataUrl) => {
                  qrcodeUrl.value = dataUrl
                })
                .catch((error) => {
                  console.error('Failed to generate QR code:', error);
                });
          } else {
            // 付款码异常
          }
        }

        if (result.data.status === 1) {
          finishedVisible.value = true;
        }
        if (result.data.status === 2) {
          payVisible.value = true;
        }
        if (result.data.status === 3) {
          timeoutVisible.value = true;
        }
      }
    }else {
    }
    // if (result) {
    //   clearInterval(timerId); // 如果状态发生变化，则停止定时器
    // }

  } catch (error) {
    console.log(error);
  }
}

// 倒计时数组
const countdowns = ref([]);

// 计算倒计时
const calculateCountdown = () => {
  setInterval(() => {
    const currentTime = new Date();
    const timeLimit = new Date(payData.value.expTime);
    const timeDiffInSeconds = (timeLimit - currentTime) / 1000;
    countdowns.value[0] = timeDiffInSeconds > 0 ? Math.floor(timeDiffInSeconds) : -1;
  }, 1000);
};


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
  top: 50px;
  left: 5%;
  right: 5%;
  height: 280px;
  background-color: #f2f2f2;
  border-radius: 20px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.p_content_inner {
  text-align: center;
  color: #333;
  position: absolute;
  top: 320px;
  left: 5%;
  right: 5%;
  height: 200px;
  background-color: #f2f2f2;
  border-radius: 20px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.p_content_button {
  text-align: center;
  color: #333;
  position: absolute;
  top: 550px;
  left: 5%;
  right: 5%;
}

.p_content_button_qr {
  text-align: center;
  color: #333;
  position: absolute;
  top: 580px;
  left: 5%;
  right: 5%;
}

.p_button {
  border: none;
  padding: 12px 24px;
  font-size: 22px;
  color: #e7dfdf;
  background: linear-gradient(to right, #064954, #125280, #1247c9);
  margin-top: 6px;
  width: 80%;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  height: 50px;
}

.copy_button {
  border: none;
  font-size: 16px;
  color: #e7dfdf;
  background: linear-gradient(to right, #d71010, #ab4f34);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  height: 30px;
  margin-top: 6px;
  margin-left: 5%;
  margin-right: 5%;
  width: 90%;
}

.copy-tip {
  display: flex;
  justify-content: space-around;
  align-items: center;
  //background: linear-gradient(to right, #f9c492, #fb6a65);
  background: linear-gradient(to right, #d71010, #064954);
  color: white;
  font-size: 14px;
  height: 24px;
  margin-top: 6px;
  margin-left: 5%;
  margin-right: 5%;
  width: 90%;
}

.jtone {
  position: relative;
  width: 0;
  height: 0;
  border-top: 12px solid transparent;
  border-left: 10px solid white;
  border-bottom: 12px solid transparent;
}

.jtone::before {
  content: "";
  position: absolute;
  top: -12px;
  left: -13px;
  width: 0;
  height: 0;
  border-top: 12px solid transparent;
  border-left: 10px solid #9f330c;
  border-bottom: 12px solid transparent;
}

.jttwo {
  position: relative;
  width: 0;
  height: 0;
  border-top: 12px solid transparent;
  border-left: 10px solid white;
  border-bottom: 12px solid transparent;
}

.jttwo::before {
  content: "";
  position: absolute;
  top: -12px;
  left: -13px;
  width: 0;
  height: 0;
  border-top: 12px solid transparent;
  border-left: 10px solid rgba(126, 60, 46, 0.99);
  border-bottom: 12px solid transparent;
}

.jtthree {
  position: relative;
  width: 0;
  height: 0;
  border-top: 12px solid transparent;
  border-left: 10px solid white;
  border-bottom: 12px solid transparent;
}

.jtthree::before {
  content: "";
  position: absolute;
  top: -12px;
  left: -13px;
  width: 0;
  height: 0;
  border-top: 12px solid transparent;
  border-left: 10px solid rgb(77, 67, 61);
  border-bottom: 12px solid transparent;
}

.medicine-bag {
  background: rgba(215, 197, 197, 0.1);
  border: 1px dashed rgba(59, 28, 23, 0.99);
  margin-top: 6px;
  border-radius: 5px;
  padding-top: 10px;
  margin-left: 5%;
  margin-right: 5%;
  width: 90%;
  font-size: 18px;
  height: 30px;
}

.medicine-bag-qr {
  background: rgba(215, 197, 197, 0.1);
  border: 1px dashed rgba(59, 28, 23, 0.99);
  border-radius: 5px;
  padding-top: 5px;
  margin-left: 10%;
  margin-right: 10%;
  width: 80%;
  font-size: 18px;
  height: 180px;
}

</style>

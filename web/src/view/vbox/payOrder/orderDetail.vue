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
      <div v-if="payTypeVisible >= 1000 && payTypeVisible < 1099">
        <div class="p_container">
          <div class="p_blue-section" v-for="index in 10" :key="index"
               :style="{ backgroundColor: generateColor(index) }"></div>
          <div class="p_content" :style="backgroundImageStyle">
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
                  <el-icon style="margin-right: 5px">
                    <WarningFilled/>
                  </el-icon>
                  请在规定时间内付款！
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
        <div class="p_content_inner" :style="backgroundImageStyle" style="margin-top: 20px;">
          <el-row>
            <el-col>
              <div style="height: 100px; margin-top: 20px">
                <div class="medicine-money-bag">
                  <span><span style="color: red">牢记</span>充值金额：<span style="color: blue">￥{{
                      payData.money
                    }}.00</span></span>
                </div>
                <div class="medicine-bag">
                  <span>{{ payData.account }}</span>
                </div>
                <div class="copy-tip">
                  <span>长按框内</span>
                  <span class="jtone"></span>
                  <span>复制</span>
                  <span class="jttwo"></span>
                  <span>记金额</span>
                  <span class="jtthree"></span>
                  <span>打开跳转</span>
                </div>
                <div v-if="!copyInfoVisible">
                  <button class="btn-copy copy_button" @click="copyInfo">① 一键复制</button>
                </div>
                <div v-else>
                  <button class="btn-copy copy_success_button" @click="copyInfo">复制成功</button>
                </div>
              </div>
            </el-col>
          </el-row>
        </div>
        <div class="p_content_button">
          <el-row :gutter="12">
            <el-col :span="24">
              <button class="btn-copy p_button" @click="openYdVisible">② 点击付款</button>
            </el-col>
          </el-row>
        </div>
      </div>

      <!--   3000 直付   -->
      <div v-if="payTypeVisible >= 3000 && payTypeVisible < 3099">
        <div class="p_container">
          <!--        <div class="p_blue-section" v-for="(color, index) in blueColors" :key="index" :style="{ backgroundColor: color }"></div>-->
          <div class="p_blue-section" v-for="index in 10" :key="index"
               :style="{ backgroundColor: generateColor(index) }"></div>
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
                  <el-icon style="margin-right: 5px">
                    <WarningFilled/>
                  </el-icon>
                  请在规定时间内付款！
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
                  <span>切换手机</span>
                  <span class="jtone"></span>
                  <span>打开微信</span>
                  <span class="jttwo"></span>
                  <span>扫一扫</span>
                  <span class="jtthree"></span>
                  <span>扫码付款</span>
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

      <!--   步骤指导   -->
      <el-dialog width="360px" v-model="dialogYd1000Visible" :draggable="true" :before-close="closeYdDialog" :style="backgroundYdImageStyle"
                 top="5vh" destroy-on-close>
        <div style="padding: 0; margin: -20px 0 0;">
          <div>
            <img alt style="width: 100%; height: 100%;border-radius: 20px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);" src="@/assets/yd_qb_jym.png">
          </div>
          <div class="p_content_yd_inner" :style="backgroundImageStyle" style="margin-top: 10px;">
            <el-row>
              <el-col>
                <div style="height: 100px; margin-top: 20px">
                  <div class="medicine-money-bag">
                  <span><span style="color: red">牢记</span>充值金额：<span style="color: blue">￥{{
                      payData.money
                    }}.00</span></span>
                  </div>
                  <div class="medicine-bag">
                    <span>{{ payData.account }}</span>
                  </div>
                  <div class="copy-tip">
                    <span>长按框内</span>
                    <span class="jtone"></span>
                    <span>复制</span>
                    <span class="jttwo"></span>
                    <span>记金额</span>
                    <span class="jtthree"></span>
                    <span>打开跳转</span>
                  </div>
                  <div v-if="!copyInfoVisible">
                    <button class="btn-copy copy_button" @click="copyInfo">一键复制</button>
                  </div>
                  <div v-else>
                    <button class="btn-copy copy_success_button" @click="copyInfo">复制成功</button>
                  </div>
                </div>
              </el-col>
            </el-row>
          </div>
        </div>
        <template #footer>
          <div class="dialog-footer">
            <div class="yd_p_content_button_qr">
              <el-row :gutter="12">
                <el-col>
                  <div v-if="readInfoVisible">
                    <button class="yd_read_p_button" @click="">我已阅读并知晓({{ countdownTime }}s)</button>
                  </div>
                  <div v-else>
                    <button class="yd_p_button" @click="openPay">点此支付</button>
                  </div>
                </el-col>
              </el-row>
            </div>
          </div>
        </template>
      </el-dialog>
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
  <div v-show="exVisible">
    <!-- 显示新的 div 的代码... -->
    <h1>订单异常，请重新下单</h1>
  </div>
</template>
<script>
export default {
  name: 'Pay',
}
</script>
<script setup>
import {ElButton, ElMessage} from 'element-plus';
import {onMounted, ref, onUnmounted, onBeforeUnmount, watch} from 'vue';
import CountDown from 'vue-canvas-countdown';
import {queryOrderSimple} from '@/api/payOrder';
import {useRoute} from 'vue-router';
import {WarningFilled} from '@element-plus/icons-vue';
import {formatTime} from "@/utils/format";
import QRCode from "qrcode";
import ClipboardJS from "clipboard";
import bgImage from '@/assets/od_info_bg.png'; // 背景图片

const backgroundImageStyle = `background-image: url(${bgImage});background-size: 100% 100%;`;
const backgroundYdImageStyle = `background-image: url(${bgImage});background-size: 100% 100%;border-radius: 10px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);top: -20px`;

// 弹窗控制标记
const dialogCountVisible = ref(true)
const payVisible = ref(false)
const payTypeVisible = ref(0)
const copyInfoVisible = ref(false)
const readInfoVisible = ref(false)
const finishedVisible = ref(false)
const timeoutVisible = ref(false)
const exVisible = ref(false)
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

// 复制
const copyInfo = async () => {
  let copyInfo = `${payData.value.account}`
  console.log("copyInfo", copyInfo)
  const clipboard = new ClipboardJS('.btn-copy', {
    text: () => copyInfo
  });

  clipboard.on('success', () => {
    ElMessage({
      type: 'success',
      message: '复制成功'
    })
    copyInfoVisible.value = true
    setTimeout(() => {
      copyInfoVisible.value = false
    }, 2000)
    clipboard.destroy(); // 销毁 ClipboardJS 实例
  });

  clipboard.on('error', () => {
    ElMessage({
      type: 'error',
      message: '复制异常'
    })
    clipboard.destroy(); // 销毁 ClipboardJS 实例
  });

};

const dialogYd1000Visible = ref(false)
const dialogYd2000Visible = ref(false)
const dialogYd3000Visible = ref(false)

const closeYdDialog = async () => {
  dialogYd1000Visible.value = false
  dialogYd2000Visible.value = false
  dialogYd3000Visible.value = false
}

const openYdVisible = async () => {
  let cid = payData.value.channel_code;
  if (cid >= 3000 && cid < 3099) {
    dialogYd3000Visible.value = true
  } else if (cid >= 2000 && cid < 2099) {
    startCountdown()
    readInfoVisible.value = true
    setTimeout(() => {
      readInfoVisible.value = false
    }, 3000)
    dialogYd2000Visible.value = true
  } else if (cid >= 1000 && cid < 1099) {
    startCountdown()
    readInfoVisible.value = true
    setTimeout(() => {
      readInfoVisible.value = false
    }, 3000)
    dialogYd1000Visible.value = true
  } else {

  }
}

const openPay = async () => {
  const clipboard = new ClipboardJS('.btn-copy', {
    text: () => payData.value.account
  });

  clipboard.on('success', () => {
    ElMessage({
      type: 'success',
      message: '复制成功'
    })
    clipboard.destroy(); // 销毁 ClipboardJS 实例
  });

  clipboard.on('error', () => {
    ElMessage({
      type: 'error',
      message: '复制异常'
    })
    clipboard.destroy(); // 销毁 ClipboardJS 实例
  });

  location.href = payData.value.resource_url;
};

// 添加一个空变量作为定时器的 ID
let timerId = null;
let timerExp = null;
let timerYD = null;

onMounted(() => {
  // 启动倒计时
  fireCD();
  // 启动定时器，每秒钟请求一次 HTTP 接口
  // timerId = setInterval(queryOrder, 1000);
  startCountdownQryOrder();

  // setCacheControl(180) // 设置缓存时间为180秒
});

onUnmounted(() => {
  // 组件销毁时清除定时器
  clearInterval(timerId);
});

const payData = ref({
  money: 0,
  exp_time: 0,
  channel_code: 0,
  account: '',
  order_id: '',
  resource_url: '',
  status: 0,
})

const reqCnt = ref(0)

const queryOrder = async (timerId) => {
  try {
    const orderId = route.query.orderId;
    console.log(orderId)
    const result = await queryOrderSimple({order_id: orderId}); // 发送 HTTP 请求
    let nowTime = new Date().getTime();
    let resExp = new Date(result.data?.exp_time).getTime();
    console.log(nowTime)
    console.log(resExp)
    payData.value = result.data
    const content = result.data.resource_url;
    const account = result.data.account;
    if (content && account) {
      console.log("qry time id ", timerId)
      clearInterval(timerId);
      startCountdownExp();
      // 如果状态发生变化，则停止定时器
      dialogCountVisible.value = false;

      if (result.code === 7) {
        dialogCountVisible.value = false;
        notFoundVisible.value = true;
      } else if (result.code === 0) {

        payTypeVisible.value = Number(result.data.channel_code);

        if (content) {
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
    } else if (result.data?.status === 0) {
      dialogCountVisible.value = false;
      clearInterval(timerId);
      exVisible.value = true;
    } else if (result.data?.status === 2 && resExp < nowTime) {
      dialogCountVisible.value = false;
      clearInterval(timerId);
      exVisible.value = true;
      console.log('超时')
    } else if (reqCnt.value < 10) {
      reqCnt.value++
    } else {
      dialogCountVisible.value = false;
      clearInterval(timerId);
      exVisible.value = true;
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
// const calculateCountdown = () => {
//   setInterval(() => {
//     const currentTime = new Date();
//     const timeLimit = new Date(payData.value.exp_time);
//     const timeDiffInSeconds = (timeLimit - currentTime) / 1000;
//     countdowns.value[0] = timeDiffInSeconds > 0 ? Math.floor(timeDiffInSeconds) : -1;
//   }, 1000);
// };

// 引导倒计时
const countdownTime = ref(3);

const countdown = (countdownTimeRef, timerRef) => {
  countdownTimeRef.value--;

  // 监听倒计时时间变化，可以在这里执行倒计时结束后的操作
  watch(countdownTimeRef, (newVal) => {
    if (newVal === 0) {
      // 在此处执行倒计时结束后的操作
      console.log(`倒计时结束 - ${countdownTimeRef.value}`);
      countdownTime.value = 3;
      // 清除定时器
      clearInterval(timerRef);
    }
  });
};

const startCountdown = () => {
  timerYD = setInterval(() => countdown(countdownTime, timerYD), 1000);
};

const countdownTimeQryOrder = ref(20);
const countdownQryOrder = (countdownTimeRef, timerRef) => {
  countdownTimeRef.value--;

  // 监听倒计时时间变化，可以在这里执行倒计时结束后的操作
  watch(countdownTimeRef, (newVal) => {
    console.log(newVal)
    queryOrder(timerRef)
    if (newVal === 0) {
      // 在此处执行倒计时结束后的操作
      console.log(`倒计时结束 - ${countdownTimeRef.value}`);
      // 清除定时器
      clearInterval(timerRef);
    }
  });
};

const startCountdownQryOrder = () => {
  timerId = setInterval(() => countdownQryOrder(countdownTimeQryOrder, timerId), 1000);
};

const startCountdownExp = () => {
  timerExp = setInterval(() => {
    const currentTime = new Date();
    const timeLimit = new Date(payData.value.exp_time);
    const timeDiffInSeconds = (timeLimit - currentTime) / 1000;
    countdowns.value[0] = timeDiffInSeconds > 0 ? Math.floor(timeDiffInSeconds) : -1;
    // console.log('timeLimit', timeLimit);
    // console.log('timeDiffInSeconds', timeDiffInSeconds);
    // console.log(timerExp)
    if (timeDiffInSeconds < 0) {
      clearInterval(timerExp);
    }
  }, 1000);
};

// 在组件销毁前清除定时器
onBeforeUnmount(() => {
  clearInterval(timerYD);
});

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
  border-radius: 10px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.p_content_yd_inner {
  text-align: center;
  color: #333;
  top: 320px;
  left: 5%;
  right: 5%;
  height: 200px;
  background-color: #f2f2f2;
  border-radius: 10px;
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

.yd_p_content_button_qr {
  text-align: center;
  color: #333;
  position: absolute;
  left: 5%;
  right: 5%;
}

.yd_read_p_button {
  border: none;
  padding: 12px 24px;
  font-size: 22px;
  color: #e7dfdf;
  background: linear-gradient(to right, #a5abb4, rgba(122, 129, 140, 0.99), #a5abb4);
  margin-top: 6px;
  width: 80%;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  height: 50px;
}

.yd_p_button {
  border: none;
  padding: 12px 24px;
  font-size: 22px;
  color: #e7dfdf;
  background: linear-gradient(90deg,#5498ff 1%,#00d9d0 100%);
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

.copy_success_button {
  border: none;
  font-size: 16px;
  color: #e7dfdf;
  background: linear-gradient(to right, #0d650a, #71ab34);
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

.medicine-money-bag {
  background: rgba(215, 197, 197, 0.1);
  border-radius: 5px;
  padding-top: 10px;
  margin-left: 5%;
  margin-right: 5%;
  width: 90%;
  height: 30px;
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

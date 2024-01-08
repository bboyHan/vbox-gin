<template>
  <div>
    <div v-show="dialogCountVisible">
      <div class="background-image">
        <img
            src="@/assets/login_bg.jpg"
            alt="banner">
      </div>
      <div class="c_container">
        <div class="c_content">
          <el-row :gutter="12">
            <el-col :span="24">
              <el-cascader
                  :options="channelCodeOptions"
                  :props="channelCodeProps"
                  @change="handleChange"
                  style="width: 80%"
              />
            </el-col>
            <el-col :span="24" style="padding-top: 10px">
              <el-input v-model="money" style="width: 80%"
                        placeholder="输入测试金额"
                        :formatter="(value) => `￥ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')"
                        :parser="(value) => value.replace(/￥\s?|(,*)/g, '')"></el-input>
            </el-col>
            <el-col :span="24" style="padding-top: 10px">
              <el-input v-model="authCaptcha" style="width: 80%" placeholder="防爆验证码"></el-input>
            </el-col>
            <el-col :span="24" style="padding-top: 10px">
              <el-button type="primary" @click="orderTest(rowId, money, authCaptcha)" style="width: 80%">订单测试</el-button>
            </el-col>
          </el-row>
        </div>
      </div>
    </div>

    <div v-show="payVisible">
      <!-- 显示新的 div 的代码... -->
      <div class="p_container">
        <!--        <div class="p_blue-section" v-for="(color, index) in blueColors" :key="index" :style="{ backgroundColor: color }"></div>-->
        <div class="p_content">
          <el-row :gutter="12">
            <el-col>
              <img src="@/assets/logo.png" alt="" style="width: 80px; height: 80px">
            </el-col>
          </el-row>
        </div>
      </div>
      <div class="p_content_button">
        <el-row :gutter="12">
          <el-col>
            <button class="p_button" @click="">立即付款</button>
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
  name: 'PayTest',
}
</script>
<script setup>
import {ElButton, ElMessageBox} from 'element-plus';
import { onMounted, ref, onUnmounted } from 'vue';
import { useRouter } from 'vue-router'
import { WarningFilled } from '@element-plus/icons-vue';
import { getChannelAccountList } from '@/api/channelAccount';
import { getChannelProductSelf } from '@/api/channelProduct';
import { createOrderTest } from '@/api/payOrder';

// 弹窗控制标记
const dialogCountVisible = ref(true)
const payVisible = ref(false)
const finishedVisible = ref(false)
const notFoundVisible = ref(false)
const router = useRouter()

// ---------- 测试页 ----------------
const money = ref('')
const authCaptcha = ref('')
const channelCodeOptions = ref([])

const channelCodeProps = {
  expandTrigger: 'hover',
  checkStrictly: false,
  emitPath: false,
}

const handleChange = (value) => {
  // console.log(value)
  rowId.value = value
}

const vcpTableData = ref([])
const rowId = ref(1)
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const orderData = ref({})

const orderTest = async (chanId, money, authCaptcha) => {
  ElMessageBox.confirm('确定要测试吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    await orderTestCreate(chanId, money, authCaptcha)
  })
}

const orderTestCreate = async (chanId, money, authCaptcha) => {
  console.log(chanId)
  console.log(money)
  console.log(authCaptcha)
  const order = await createOrderTest({
    money: Number(money),
    channel_code: chanId,
    auth_captcha: authCaptcha,
  })
  console.log("------------ o ---------")
  console.log(order)
  console.log(order.data)
  orderData.value = order.data.order_id
  if (orderData.value){
    await router.push({name: 'Pay', query: {orderId: orderData.value}})
  }
}

const getTableData = async () => {
  const table = await getChannelAccountList({page: page.value, pageSize: pageSize.value, ...searchInfo.value})
  const vcpTable = await getChannelProductSelf({page: 1, pageSize: 999, ...searchInfo.value})

  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
    vcpTableData.value = vcpTable.data.list
    setOptions()
  }
}

getTableData()

const setOptions = () => {
  channelCodeOptions.value = []
  setChannelCodeOptions(vcpTableData.value, channelCodeOptions.value, false)
}

const setChannelCodeOptions = (ChannelCodeData, optionsData, disabled) => {
  ChannelCodeData &&
  ChannelCodeData.forEach(item => {
    if (item.children && item.children.length) {
      const option = {
        value: item.channelCode + '',
        label: item.productName,
        children: []
      }
      setChannelCodeOptions(
          item.children,
          option.children,
      )
      optionsData.push(option)
    } else {
      const option = {
        value: item.channelCode + '',
        label: item.productName,
      }
      optionsData.push(option)
    }
  })
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
  background: linear-gradient(to right, #064954, #125280, #1247c9);
  margin-top: 6px;
  width: 80%;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  height: 50px;
}

.background-image {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}
.background-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>

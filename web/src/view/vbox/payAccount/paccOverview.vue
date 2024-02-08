<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
        <el-form-item label="商户ID" prop="cid">
          <el-input v-model="searchInfo.pAccount" placeholder="商户ID"/>
        </el-form-item>
        <el-form-item label="通道ID" prop="cid">
          <el-cascader
              v-model="searchInfo.cid"
              :options="channelCodeOptions"
              :props="channelCodeProps"
              @change="handleChange"
              style="width: 100%"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div>
      <!-- 商户信息 -->
      <el-row :gutter="24">
        <el-col :span="24" :xs="24">
          <div class="flex justify-between items-center flex-wrap" style="margin-left: 10px"><h2>付方账号</h2></div>
        </el-col>
        <el-col :span="6" :xs="24">
          <CenterCard title="当前查询商户" :custom-style="accCustomStyle">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999" ></span>
            </template>
            <template #body>
              <!--              <ReclaimMileage :channel-code="searchInfo.cid" :acc-on="accOn" :acc-off="accOff" :acc-total="accTotal"/>-->
              <div class="acc-container">
                <div class="indicator">
                  <span>
                    <div class="label">商户ID</div>
                    <div class="value">{{ searchInfo.pAccount }}</div>
                  </span>
                  <span>
                    <div class="label">通道编码</div>
                    <div class="value">{{ searchInfo.cid }}</div>
                  </span>
                </div>
              </div>
            </template>
          </CenterCard>
        </el-col>
      </el-row>

      <!-- 成单统计 -->
      <el-row :gutter="24">
        <el-col :span="24" :xs="24">
          <div class="flex justify-between items-center flex-wrap" style="margin-left: 10px"><h2>成单数据概览</h2></div>
        </el-col>
        <el-col :span="6" :xs="24">
          <CenterCard title="近1小时成单" :custom-style="order1CustomStyle">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999" />
            </template>
            <template #body>
              <!--              <Order :channel-code="searchInfo.cid"/>-->
              <div class="acc-container">
                <div class="indicator">
                    <span>
                      <div class="label">总笔数</div>
                      <div class="value">{{ nearOneHourRate.x2 }}</div>
                    </span>
                  <span>
                      <div class="label">成单数</div>
                      <div class="value">{{ nearOneHourRate.x1 }}</div>
                    </span>
                  <span>
                      <div class="label">成率</div>
                      <div class="value">{{ calculatePercentage(nearOneHourRate.x1, nearOneHourRate.x2) }}% </div>
                    </span>
                </div>
              </div>
            </template>
          </CenterCard>
        </el-col>

        <el-col :span="6" :xs="24">
          <CenterCard title="昨日成单" :custom-style="order2CustomStyle">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999" />
            </template>
            <template #body>
              <!--              <Order :channel-code="searchInfo.cid"/>-->
              <div class="acc-container">
                <div class="indicator">
                  <span>
                    <div class="label">总笔数</div>
                    <div class="value">{{ nearYesterdayRate.x2 }}</div>
                  </span>
                  <span>
                    <div class="label">成单数</div>
                    <div class="value">{{ nearYesterdayRate.x1 }}</div>
                  </span>
                  <span>
                    <div class="label">成率</div>
                    <div class="value">{{ calculatePercentage(nearYesterdayRate.x1, nearYesterdayRate.x2) }}% </div>
                  </span>
                </div>
              </div>
            </template>
          </CenterCard>
        </el-col>

        <el-col :span="6" :xs="24">
          <CenterCard title="今日成单" :custom-style="order3CustomStyle">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999" />
            </template>
            <template #body>
              <!--              <Order :channel-code="searchInfo.cid"/>-->
              <div class="acc-container">
                <div class="indicator">
                  <span>
                    <div class="label">总笔数</div>
                    <div class="value">{{ nearTodayRate.x2 }}</div>
                  </span>
                  <span>
                    <div class="label">成单数</div>
                    <div class="value">{{ nearTodayRate.x1 }}</div>
                  </span>
                  <span>
                    <div class="label">成率</div>
                    <div class="value">{{ calculatePercentage(nearTodayRate.x1, nearTodayRate.x2) }}% </div>
                  </span>
                </div>
              </div>
            </template>
          </CenterCard>
        </el-col>
        <el-col :span="6" :xs="24"></el-col>

        <el-col :span="6" :xs="24">
          <CenterCard title="近1小时成单金额" :custom-style="order1CustomStyle">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999" />
            </template>
            <template #body>
              <!--              <Order :channel-code="searchInfo.cid"/>-->
              <div class="acc-container">
                <div class="indicator">
                  <span>
                    <div class="label">金额</div>
                    <div class="value">{{ formatMoney(sumData.x4) }}</div>
                  </span>
                </div>
              </div>
            </template>
          </CenterCard>
        </el-col>

        <el-col :span="6" :xs="24">
          <CenterCard title="昨日金额" :custom-style="order2CustomStyle">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999" />
            </template>
            <template #body>
              <!--              <Order :channel-code="searchInfo.cid"/>-->
              <div class="acc-container">
                <div class="indicator">
                  <span>
                    <div class="label">金额</div>
                    <div class="value">{{ formatMoney(sumData.x2) }}</div>
                  </span>
                </div>
              </div>
            </template>
          </CenterCard>
        </el-col>

        <el-col :span="6" :xs="24">
          <CenterCard title="今日金额" :custom-style="order3CustomStyle">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999" />
            </template>
            <template #body>
              <!--              <Order :channel-code="searchInfo.cid"/>-->
              <div class="acc-container">
                <div class="indicator">
                  <span>
                    <div class="label">金额</div>
                    <div class="value">{{ formatMoney(sumData.x3) }}</div>
                  </span>
                </div>
              </div>
            </template>
          </CenterCard>
        </el-col>

        <el-col :span="6" :xs="24">
          <CenterCard title="今日待付金额（含失效单）" :custom-style="order4CustomStyle">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999" />
            </template>
            <template #body>
              <!--              <Order :channel-code="searchInfo.cid"/>-->
              <div class="acc-container">
                <div class="indicator">
                  <span>
                    <div class="label">待付金额</div>
                    <div class="value">{{ formatMoney(nearTodayRate.x4-nearTodayRate.x3) }}</div>
                  </span>
                </div>
              </div>
            </template>
          </CenterCard>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import CenterCard from '@/view/vbox/dashboard/dataCenterComponents/CenterCard.vue'
import {reactive, ref, nextTick, defineEmits, onMounted, watch, toRefs} from "vue";
import {getChannelProductSelf} from "@/api/channelProduct";
import {getChannelAccountList} from "@/api/channelAccount";
import {getOrderDataOverview, getPayOrderOverview, getPayOrderRate} from "@/api/payOrder";
import {calculatePercentage, formatMoney} from "@/utils/format";
import {getUserWalletSelf} from "@/api/userWallet";


const searchInfo = ref({})
// 获取当前时间
let localTime = new Date();
// 获取当前时区相对于UTC的偏移量（分钟）
let offset = localTime.getTimezoneOffset();
// 计算东八区相对于UTC的偏移量（分钟）
let easternOffset = 8 * 60;
// let easternOffset = 0;
// 计算东八区当前时间
// let endTime = new Date(localTime.getTime() + (easternOffset + offset) * 60 * 1000);

// 获取当前时间的分钟数
let currentMinute = localTime.getMinutes();
// 计算所在的分钟数分组
let adjustedMinute = Math.ceil(currentMinute / 5) * 5;

// 设置分钟数为分组的结束值，秒和毫秒都设置为0
localTime.setMinutes(adjustedMinute, 0, 0);
// 前一个小时
const startTimeOneHour = new Date(localTime.getTime() + (easternOffset + offset) * 60 * 1000 - 60 * 60 * 1000);
const endTimeOneHour = new Date(localTime.getTime() + (easternOffset + offset) * 60 * 1000);

// 设置时间为0点
let zeroTime = new Date();
zeroTime.setHours(0, 0, 0, 0);
let zero24Time = new Date();
zero24Time.setHours(24, 0, 0, 0);
const startTimeToday = new Date(zeroTime.getTime() + (easternOffset + offset) * 60 * 1000);
const endTimeToday = new Date(zero24Time.getTime() + (easternOffset + offset) * 60 * 1000);

//昨天
const startTimeYesterday = new Date(zeroTime.getTime() + (easternOffset + offset) * 60 * 1000 - 24 * 60 * 60 * 1000);
const endTimeYesterday = new Date(zero24Time.getTime() + (easternOffset + offset) * 60 * 1000 - 24 * 60 * 60 * 1000);

const elSearchFormRef = ref()

const setOptions = async () =>{
  channelCodeOptions.value = []
  setChannelCodeOptions(vcpTableData.value, channelCodeOptions.value, false)
}

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

const isUpd = ref(false)
const accCustomStyle = ref({
  background: 'linear-gradient(to right, #3498db, #2ecc71)',
  color: '#FFF',
  height: '150px',
})
const order1CustomStyle = ref({
  background: 'linear-gradient(to right, #be2eba, #5b2ecc)',
  color: '#FFF',
  height: '150px',
})
const order2CustomStyle = ref({
  background: 'linear-gradient(to right, #be2eba, #5b2ecc)',
  color: '#FFF',
  height: '150px',
})
const order3CustomStyle = ref({
  background: 'linear-gradient(to right, #be2eba, #5b2ecc)',
  color: '#FFF',
  height: '150px',
})
const order4CustomStyle = ref({
  background: 'linear-gradient(to right, #22111a, #606266)',
  color: '#FFF',
  height: '150px',
})
const walletCustomStyle = ref({
  background: 'linear-gradient(to right, #22111a, #606266)',
  color: '#FFF',
  height: '140px',
})
// 余额
const userBalance = ref(0)


// 搜索
const onSubmit = () => {
  console.log("searchInfo.value",searchInfo.value)
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    console.log("elSearchFormRef.value",elSearchFormRef.value)
    getTableData()
  })
}

const searchRule = reactive({})

//通道产品
const channelCodeOptions = ref([])
const vcpTableData = ref([])
//账号明细
const accTableData = ref([])

const channelCodeProps = {
  expandTrigger: 'hover',
  checkStrictly: false,
  emitPath: false,
}

const handleChange = (value) => {

  console.log('cid:',value)
  searchInfo.value.cid = value
  // getTableData()
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

const accOn = ref()
const accOff = ref()
const accTotal = ref()

const setAccSwitchView = async () => {
  const re = accTableData.value
  let countOn = 0
  let countOff = 0
  for (let i = 0; i < re.length; i++) {
    if (re[i].status === 1){
      countOn++;
    }else {
      countOff++;
    }
  }

  accOn.value = countOn
  accOff.value = countOff
  accTotal.value = re.length
  console.log(re.length)

}

const nearOneHourRate = ref({})
const nearYesterdayRate = ref({})
const nearTodayRate = ref({})
const nearOneHourCnt = ref()
const nearOneHourSum = ref()
const nearTodayCnt = ref()
const nearTodaySum = ref()
const sumData = ref({})

const getTableData = async() => {
  await nextTick()
  const vcpTable = await getChannelProductSelf({ page: 1, pageSize: 9999, ...searchInfo.value })
  if (vcpTable.code === 0) {
    vcpTableData.value = vcpTable.data.list
    setOptions()
  }
  let nearOneHourRateResult = await getPayOrderRate({ page: 1, pageSize: 9999, pAccount: searchInfo.value.pAccount, channelCode: searchInfo.value.cid, startTime: Math.floor(startTimeOneHour.getTime() / 1000), endTime: Math.floor(endTimeOneHour.getTime() / 1000), keyword: 'cas'})
  let nearYesterdayRateResult = await getPayOrderRate({ page: 1, pageSize: 9999, pAccount: searchInfo.value.pAccount, channelCode: searchInfo.value.cid, startTime: Math.floor(startTimeYesterday.getTime() / 1000), endTime: Math.floor(endTimeYesterday.getTime() / 1000), keyword: 'cas'})
  let nearTodayRateResult = await getPayOrderRate({ page: 1, pageSize: 9999, pAccount: searchInfo.value.pAccount, channelCode: searchInfo.value.cid, startTime: Math.floor(startTimeToday.getTime() / 1000), endTime: Math.floor(endTimeToday.getTime() / 1000), keyword: 'cas'})
  let nearOneHourCntResult = await getPayOrderOverview({ page: 1, pageSize: 9999, orderStatus:1, pAccount: searchInfo.value.pAccount, channelCode: searchInfo.value.cid, startTime: Math.floor(startTimeOneHour.getTime() / 1000), endTime: Math.floor(endTimeOneHour.getTime() / 1000), interval:  '5m', keyword: 'cnt', format: 'HH:mm'})
  let nearOneHourSumResult = await getPayOrderOverview({ page: 1, pageSize: 9999, orderStatus:1, pAccount: searchInfo.value.pAccount, channelCode: searchInfo.value.cid, startTime: Math.floor(startTimeOneHour.getTime() / 1000), endTime: Math.floor(endTimeOneHour.getTime() / 1000), interval:  '5m', keyword:'sum', format: 'HH:mm'})
  let nearTodayCntResult = await getPayOrderOverview({ page: 1, pageSize: 9999, orderStatus:1, pAccount: searchInfo.value.pAccount, channelCode: searchInfo.value.cid, startTime: Math.floor(startTimeToday.getTime() / 1000), endTime: Math.floor(endTimeToday.getTime() / 1000), interval:  '30m', keyword: 'cnt', format: 'HH:mm'})
  let nearTodaySumResult = await getPayOrderOverview({ page: 1, pageSize: 9999, orderStatus:1, pAccount: searchInfo.value.pAccount, channelCode: searchInfo.value.cid, startTime: Math.floor(startTimeToday.getTime() / 1000), endTime: Math.floor(endTimeToday.getTime() / 1000), interval:  '30m', keyword:'sum', format: 'HH:mm'})

  let sumDataOverview = await getOrderDataOverview({
    channelCode: searchInfo.value.cid,
    pAccount: searchInfo.value.pAccount,
  })

  sumData.value = sumDataOverview.data.list[0]

  nearOneHourCnt.value = nearOneHourCntResult
  nearOneHourSum.value = nearOneHourSumResult
  nearTodayCnt.value = nearTodayCntResult
  nearTodaySum.value = nearTodaySumResult
  nearOneHourRate.value = (nearOneHourRateResult.data)[0]
  nearTodayRate.value = (nearTodayRateResult.data)[0];
  nearYesterdayRate.value = (nearYesterdayRateResult.data)[0];
  console.log(nearOneHourRate.value)
  console.log(nearTodayRate.value)

  let balanceVal = await getUserWalletSelf({ ...searchInfo.value })
  userBalance.value = balanceVal.data.balance

}

getTableData()


// 在 mounted 钩子中，初始化 watch
onMounted(() => {
  // 监听 searchInfo.cid 的变化
  // watch(() => searchInfo.value.cid, (newCid, oldCid) => {
  //   // 当 searchInfo.cid 变化时，手动触发子组件刷新
  //   console.log("发生cid变化了", newCid, oldCid)
  //   // 你需要确保 lineCharts 组件中有 refresh 方法
  // });
});

watch(() => searchInfo.value.cid, (newCid, oldCid) => {
  // 当 searchInfo.cid 变化时，手动触发子组件刷新
  console.log("发生cid变化了", newCid, oldCid)
  // 你需要确保 lineCharts 组件中有 refresh 方法
});

</script>

<script>
export default {
  name: 'DataCenter',
}
</script>

<style lang="scss" scoped>

.data-center-box{
  width: 100%;
  display: grid;
  grid-template-columns: 2fr 4fr;
  column-gap: 10px;
}

.acc-container{
  color: #FFFFFF;
}
.indicator {
  display: flex;
  justify-content: space-around; // 使子元素水平居中展开
  padding: 15px;
  border-radius: 8px; // 添加圆角
}

.indicator span {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 10px; // 调整间距

  &:not(:last-child) {
    border-right: 2px solid #fff; // 白色边框
    margin-right: 15px; // 调整间距
  }
}

.label {
  color: #F5F5F5;
  font-size: 14px;
}

.value {
  color: #FFFFFF;
  font-size: 30px;
  font-weight: bold;
  margin-top: 5px; // 调整间距

}

</style>

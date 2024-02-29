<template>
  <div>
    <div>
      <el-row :gutter="12">
        <el-col :span="6" :xs="24">
          <el-col :span="24" :xs="24">
            <div class="flex justify-between items-center flex-wrap" style="margin-left: 10px"><h2>我的积分</h2></div>
          </el-col>
              <CenterCard title="我的积分" :custom-style="walletCustomStyle">
                <template #action>
                  <span class="gvaIcon-prompt" style="color: #999"/>
                </template>
                <template #body>
                  <!--              <Order :channel-code="searchInfo.cid"/>-->
                  <div class="acc-container">
                    <div class="indicator">
                    <span>
                      <div class="label"></div>
                      <div class="value">{{ userBalance }}</div>
                    </span>
                    </div>
                  </div>
                </template>
              </CenterCard>
        </el-col>

<!--        <el-col :span="6" :xs="24">
          &lt;!&ndash; 通道账号 &ndash;&gt;
            <el-col :span="24" :xs="24">
              <div class="flex justify-between items-center flex-wrap" style="margin-left: 10px"><h2>通道账号</h2></div>
            </el-col>
              <CenterCard title="当前通道账号" :custom-style="accCustomStyle">
                <template #action>
                  <span class="gvaIcon-prompt" style="color: #999"></span>
                </template>
                <template #body>
                  &lt;!&ndash;              <ReclaimMileage :channel-code="searchInfo.cid" :acc-on="accOn" :acc-off="accOff" :acc-total="accTotal"/>&ndash;&gt;
                  <div class="acc-container">
                    <div class="indicator">
                  <span>
                    <div class="label">账户总数</div>
                    <div class="value">{{ accTotal }}</div>
                  </span>
                      <span>
                    <div class="label">开启数</div>
                    <div class="value">{{ accOn }}</div>
                  </span>
                      <span>
                    <div class="label">关闭数</div>
                    <div class="value">{{ accOff }}</div>
                  </span>
                    </div>
                  </div>
                </template>
              </CenterCard>
        </el-col>-->
      </el-row>
    </div>

    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule"
               @keyup.enter="onSubmit">
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

      <!-- 成单统计 -->
      <el-row :gutter="24">
        <el-col :span="24" :xs="24">
          <div class="flex justify-between items-center flex-wrap" style="margin-left: 10px"><h2>成单数据概览</h2></div>
        </el-col>
        <el-col :span="6" :xs="24">
          <CenterCard title="近1小时成单" :custom-style="order1CustomStyle">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999"/>
            </template>
            <template #body>
              <!--              <Order :channel-code="searchInfo.cid"/>-->
              <div class="acc-container">
                <div class="indicator">
                  <span>
                    <div class="label">成单数 / 总笔数</div>
                    <div class="value">{{ nearOneHourRate.x1 }} / {{ nearOneHourRate.x2 }}</div>
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
              <span class="gvaIcon-prompt" style="color: #999"/>
            </template>
            <template #body>
              <!--              <Order :channel-code="searchInfo.cid"/>-->
              <div class="acc-container">
                <div class="indicator"><span>
                    <div class="label">成单数 / 总笔数</div>
                    <div class="value">{{ nearYesterdayRate.x1 }} / {{ nearYesterdayRate.x2 }}</div>
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
              <span class="gvaIcon-prompt" style="color: #999"/>
            </template>
            <template #body>
              <!--              <Order :channel-code="searchInfo.cid"/>-->
              <div class="acc-container">
                <div class="indicator">
                  <span>
                    <div class="label">成单数 / 总笔数</div>
                    <div class="value">{{ nearTodayRate.x1 }} / {{ nearTodayRate.x2 }}</div>
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
              <span class="gvaIcon-prompt" style="color: #999"/>
            </template>
            <template #body>
              <!--              <Order :channel-code="searchInfo.cid"/>-->
              <div class="acc-container">
                <div class="indicator">
                  <span>
                    <div class="label">金额</div>
                    <div class="value">{{ formatMoney(sumData.x4) }}</div>
                  </span>
                  <!--                  <span>
                                      <div class="label">待付金额</div>
                                      <div class="value">{{ formatMoney(nearOneHourRate.x4) }}</div>
                                    </span>-->
                </div>
              </div>
            </template>
          </CenterCard>
        </el-col>

        <el-col :span="6" :xs="24">
          <CenterCard title="昨日金额" :custom-style="order2CustomStyle">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999"/>
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
              <span class="gvaIcon-prompt" style="color: #999"/>
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
              <span class="gvaIcon-prompt" style="color: #999"/>
            </template>
            <template #body>
              <!--              <Order :channel-code="searchInfo.cid"/>-->
              <div class="acc-container">
                <div class="indicator">
                  <span>
                    <div class="label">待付金额</div>
                    <div class="value">{{ formatMoney(nearTodayRate.x4 - nearTodayRate.x3) }}</div>
                  </span>
                </div>
              </div>
            </template>
          </CenterCard>
        </el-col>
      </el-row>

      <!--   趋势图   -->
      <el-row :gutter="24">
        <el-col :span="24" :xs="24">
          <div class="flex justify-between items-center flex-wrap" style="margin-left: 10px"><h2>近1小时实时成单</h2>
          </div>
        </el-col>
        <el-col :span="12" :xs="24">
          <CenterCard title="近1小时实时成单(金额)">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999"/>
            </template>
            <template #body>
              <!--              <lineCharts :channel-code="searchInfo.cid" :start-time="startTimeOneHour" :end-time="endTimeOneHour"-->
              <!--                          interval="5m" keyword="sum" format="HH:mm" unit="元"/>-->
              <lineCharts :chart-data="nearOneHourSum" format="HH:mm" unit="元"/>
            </template>
          </CenterCard>
          <!--          <CenterCard title="近1小时实时成单(数量)" style="grid-column-start: span 2;">-->
        </el-col>
        <el-col :span="12" :xs="24">
          <CenterCard title="近1小时实时成单(数量)">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999"/>
            </template>
            <template #body>
              <!--              <lineCharts :channel-code="searchInfo.cid" :start-time="startTimeOneHour" :end-time="endTimeOneHour" interval="5m" keyword="cnt" format="HH:mm" unit="笔"/>-->
              <lineCharts :chart-data="nearOneHourCnt" format="HH:mm" unit="笔"/>
            </template>
          </CenterCard>
        </el-col>
      </el-row>
      <el-row :gutter="24">
        <el-col :span="24" :xs="24">
          <div class="flex justify-between items-center flex-wrap" style="margin-left: 10px"><h2>今日实时成单</h2></div>
        </el-col>
        <el-col :span="12" :xs="24">
          <CenterCard title="今日成单（金额）">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999"/>
            </template>
            <template #body>
              <!--              <lineCharts :channel-code="searchInfo.cid" :start-time="startTimeToday" :end-time="endTimeToday" interval="30m" keyword="sum" format="HH:mm" unit="元"/>-->
              <lineCharts :chart-data="nearTodaySum" format="HH:mm" unit="元"/>
            </template>
          </CenterCard>
        </el-col>
        <el-col :span="12" :xs="24">
          <CenterCard title="今日成单（数量）">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999"/>
            </template>
            <template #body>
              <!--              <lineCharts :channel-code="searchInfo.cid" :start-time="startTimeToday" :end-time="endTimeToday" interval="30m" keyword="cnt" format="HH:mm" unit="笔"/>-->
              <lineCharts :chart-data="nearTodayCnt" format="HH:mm" unit="笔"/>
            </template>
          </CenterCard>
        </el-col>
      </el-row>


    </div>
  </div>
</template>

<script setup>
import CenterCard from './dataCenterComponents/centerCard.vue'
import ChainRatio from './dataCenterComponents/chainRatio.vue'
import ReclaimMileage from './dataCenterComponents/ReclaimMileage.vue'
import RecoveryRate from './dataCenterComponents/RecoveryRate.vue'
import Order from './dataCenterComponents/order.vue'
import lineCharts from './dataCenterComponents/lineCharts.vue'
import part from './dataCenterComponents/part.vue'
import {reactive, ref, nextTick, defineEmits, onMounted, watch, toRefs} from "vue";
import {getChannelProductSelf} from "@/api/channelProduct";
import {getChannelAccountList} from "@/api/channelAccount";
import {getOrderDataOverview, getPayOrderOverview, getPayOrderRate} from "@/api/payOrder";
import {calculatePercentage, formatMoney} from "@/utils/format";
import {getUserWalletSelf} from "@/api/userWallet";
import bgImage from '@/assets/bg.jpg'; // 背景图片

const backgroundImageStyle = `background-image: url(${bgImage});background-size: 100% 100%;`;

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


const setOptions = async () => {
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
  background: 'linear-gradient(to right, #be2eba, #5b2ecc)',
  color: '#FFF',
  height: '120px',
})
const order1CustomStyle = ref({
  background: 'linear-gradient(to right, #2ecc71, #3498db)',
  color: '#FFF',
  height: '120px',
})
const order2CustomStyle = ref({
  background: 'linear-gradient(to right, #2ecc71, #3498db)',
  color: '#FFF',
  height: '120px',
})
const order3CustomStyle = ref({
  background: 'linear-gradient(to right, #2ecc71, #3498db)',
  color: '#FFF',
  height: '120px',
})
const order4CustomStyle = ref({
  background: 'linear-gradient(to right, #22111a, #606266)',
  color: '#FFF',
  height: '120px',
})
const walletCustomStyle = ref({
  background: 'linear-gradient(to right, #22111a, #606266)',
  color: '#FFF',
  height: '120px',
})
// 余额
const userBalance = ref(0)


// 搜索
const onSubmit = () => {
  console.log("searchInfo.value", searchInfo.value)
  elSearchFormRef.value?.validate(async (valid) => {
    if (!valid) return
    console.log("elSearchFormRef.value", elSearchFormRef.value)
    getTableData()
  })
}

const searchRule = reactive({})

//通道产品
const channelCodeOptions = ref([])
const vcpTableData = ref([])

const channelCodeProps = {
  expandTrigger: 'hover',
  checkStrictly: false,
  emitPath: false,
}

const handleChange = (value) => {

  console.log('cid:', value)
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

const nearOneHourRate = ref({})
const nearYesterdayRate = ref({})
const nearTodayRate = ref({})
const nearOneHourCnt = ref()
const nearOneHourSum = ref()
const nearTodayCnt = ref()
const nearTodaySum = ref()
const sumData = ref({})

const getTableData = async () => {
  await nextTick()
  const vcpTablePromise = getChannelProductSelf({page: 1, pageSize: 9999, ...searchInfo.value})

  let nearOneHourRateResultPromise = getPayOrderRate({
    page: 1,
    pageSize: 9999,
    channelCode: searchInfo.value.cid,
    startTime: Math.floor(startTimeOneHour.getTime() / 1000),
    endTime: Math.floor(endTimeOneHour.getTime() / 1000),
    keyword: 'cas'
  })
  let nearYesterdayRateResultPromise = getPayOrderRate({
    page: 1,
    pageSize: 9999,
    channelCode: searchInfo.value.cid,
    startTime: Math.floor(startTimeYesterday.getTime() / 1000),
    endTime: Math.floor(endTimeYesterday.getTime() / 1000),
    keyword: 'cas'
  })
  let nearTodayRateResultPromise = getPayOrderRate({
    page: 1,
    pageSize: 9999,
    channelCode: searchInfo.value.cid,
    startTime: Math.floor(startTimeToday.getTime() / 1000),
    endTime: Math.floor(endTimeToday.getTime() / 1000),
    keyword: 'cas'
  })
  let nearOneHourCntResultPromise = getPayOrderOverview({
    page: 1,
    pageSize: 9999,
    orderStatus: 1,
    channelCode: searchInfo.value.cid,
    startTime: Math.floor(startTimeOneHour.getTime() / 1000),
    endTime: Math.floor(endTimeOneHour.getTime() / 1000),
    interval: '5m',
    keyword: 'cnt',
    format: 'HH:mm'
  })
  let nearOneHourSumResultPromise = getPayOrderOverview({
    page: 1,
    pageSize: 9999,
    orderStatus: 1,
    channelCode: searchInfo.value.cid,
    startTime: Math.floor(startTimeOneHour.getTime() / 1000),
    endTime: Math.floor(endTimeOneHour.getTime() / 1000),
    interval: '5m',
    keyword: 'sum',
    format: 'HH:mm'
  })
  let nearTodayCntResultPromise = getPayOrderOverview({
    page: 1,
    pageSize: 9999,
    orderStatus: 1,
    channelCode: searchInfo.value.cid,
    startTime: Math.floor(startTimeToday.getTime() / 1000),
    endTime: Math.floor(endTimeToday.getTime() / 1000),
    interval: '30m',
    keyword: 'cnt',
    format: 'HH:mm'
  })
  let nearTodaySumResultPromise = getPayOrderOverview({
    page: 1,
    pageSize: 9999,
    orderStatus: 1,
    channelCode: searchInfo.value.cid,
    startTime: Math.floor(startTimeToday.getTime() / 1000),
    endTime: Math.floor(endTimeToday.getTime() / 1000),
    interval: '30m',
    keyword: 'sum',
    format: 'HH:mm'
  })

  let sumDataOverviewPromise = getOrderDataOverview({
    channelCode: searchInfo.value.cid,
    pAccount: searchInfo.value.pAccount,
  })

  let balanceValPromise = getUserWalletSelf({...searchInfo.value})

  const [vcpTable, nearOneHourRateResult, nearYesterdayRateResult, nearTodayRateResult, nearOneHourCntResult,
    nearOneHourSumResult, nearTodayCntResult, nearTodaySumResult, sumDataOverview, balanceVal
  ] = await Promise.all([vcpTablePromise, nearOneHourRateResultPromise, nearYesterdayRateResultPromise,
    nearTodayRateResultPromise, nearOneHourCntResultPromise, nearOneHourSumResultPromise, nearTodayCntResultPromise,
    nearTodaySumResultPromise, sumDataOverviewPromise, balanceValPromise])

  if (vcpTable.code === 0) {
    vcpTableData.value = vcpTable.data.list
    setOptions()
  }
  if (sumDataOverview.code === 0) {
    sumData.value = sumDataOverview.data.list[0];
  }
  if (nearOneHourCntResult.code === 0) {
    nearOneHourCnt.value = nearOneHourCntResult;
  }
  if (nearOneHourSumResult.code === 0) {
    nearOneHourSum.value = nearOneHourSumResult;
  }
  if (nearTodayCntResult.code === 0) {
    nearTodayCnt.value = nearTodayCntResult;
  }
  if (nearTodaySumResult.code === 0) {
    nearTodaySum.value = nearTodaySumResult;
  }
  if (nearOneHourRateResult.code === 0) {
    nearOneHourRate.value = (nearOneHourRateResult.data)[0]
  }
  if (nearTodayRateResult.code === 0) {
    nearTodayRate.value = (nearTodayRateResult.data)[0];
  }
  if (nearYesterdayRateResult.code === 0) {
    nearYesterdayRate.value = (nearYesterdayRateResult.data)[0];
  }
  if (balanceVal.code === 0) {
    userBalance.value = balanceVal.data.balance
  }

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

.data-center-box {
  width: 100%;
  display: grid;
  grid-template-columns: 2fr 4fr;
  column-gap: 10px;
}

.acc-container {
  color: #FFFFFF;
}

.indicator {
  display: flex;
  justify-content: space-around; // 使子元素水平居中展开
  padding: 10px;
  border-radius: 8px; // 添加圆角
}

.indicator span {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px; // 调整间距

  &:not(:last-child) {
    border-right: 2px solid #fff; // 白色边框
    margin-right: 5px; // 调整间距
  }
}

.label {
  color: #F5F5F5;
  font-size: 14px;
}

.value {
  color: #FFFFFF;
  font-size: 22px;
  font-weight: bold;
  margin-top: 5px; // 调整间距
}

.value-small {
  color: #FFFFFF;
  font-size: 20px;
  font-weight: bold;
  margin-top: 5px; // 调整间距
}

</style>

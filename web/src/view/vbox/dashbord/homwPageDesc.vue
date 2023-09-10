<template>
  <div class="page">


    <div class="gva-card-box">
    <div class="gva-card quick-entrance">
      <h3>近一小时交易实时概况</h3>
      <el-row>
        <el-col :span="6">
          <el-statistic  :value="usersOrderHData.tInCome" >
            <template #title>
            <div style="font-size: 14px; align-items: center">最近一小时收入</div>
            </template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>元(RMB) </span>
            <span :class="hColor">
              环比 {{ hPercentage}}%
              
              <el-icon v-if="hColor > 0">
                <CaretTop />
              </el-icon>
              <el-icon v-else>
                <CaretBottom />
              </el-icon>
              
            </span>
            <el-tooltip
                effect="dark"
                content="环比是指与前一周期相比较的增长或降低的比例或百分比"
                placement="top"
              >
              <el-icon style="margin-left: 4px" :size="12">
                  <Warning />
                </el-icon>
            </el-tooltip>
          </div>
        </div>
        </el-col>
        <el-col :span="6">
          <el-statistic :value="usersOrderHData.newOpenChId">
            <template #title>
              <div style="font-size: 14px; display: inline-flex; align-items: center">
                通道概览
                <!-- <el-icon style="margin-left: 4px" :size="12">
                  <Male />
                </el-icon> -->
              </div>
            </template>
            <template #suffix>/{{ usersOrderHData.openChId }}/{{ usersOrderHData.chIdCnt}}</template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>新开启个数/开启总数/通道总数</span>
            <!-- <span class="green">
              
              <el-icon>
                <CaretTop />
              </el-icon>
            </span> -->
          </div>
        </div>
        </el-col>
        <el-col :span="6">
          <el-statistic  :value="usersOrderHData.tOrderQuantify" >
            <template #title>
            <div style="font-size: 14px; align-items: center">最近一小时订单量</div>
            </template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>单</span>
            <!-- <span class="green">
              {{ usersOrderData.tOkRate }} %)
            </span> -->
          </div>
        </div>
        </el-col>
        <el-col :span="6">
          <el-statistic  :value="usersOrderHData.tOkOrderQuantify" >
            <template #title>
            <div style="font-size: 14px; align-items: center">最近一小时成单量</div>
            </template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>单</span>
            <span class="green">
              (成功率 {{ usersOrderHData.tOkRate }} %)
            </span>
          </div>
        </div>
        </el-col>
      </el-row>
    </div>
        
  </div>
  

 
  <div class="gva-card-box">
    <div class="gva-card quick-entrance">
      <h3>当日交易实时概况</h3>
      <el-row>
        <el-col :span="6">
          <el-statistic  :value="usersOrderData.tInCome" >
            <template #title>
            <div style="font-size: 14px; align-items: center">今日收入</div>
            </template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>元(RMB)</span>
            <!-- <span class="green">
              元
              <el-icon>
                <CaretTop />
              </el-icon>
            </span> -->
          </div>
        </div>
        </el-col>
        <el-col :span="6">
          <el-statistic :value="usersOrderData.openChId">
            <template #title>
              <div style="font-size: 14px; display: inline-flex; align-items: center">
                通道概览
                <!-- <el-icon style="margin-left: 4px" :size="12">
                  <Male />
                </el-icon> -->
              </div>
            </template>
            <template #suffix>/{{ usersOrderData.chIdCnt }}</template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>开启个数/通道总数</span>
            <!-- <span class="green">
              
              <el-icon>
                <CaretTop />
              </el-icon>
            </span> -->
          </div>
        </div>
        </el-col>
        <el-col :span="6">
          <el-statistic  :value="usersOrderData.tOrderQuantify" >
            <template #title>
            <div style="font-size: 14px; align-items: center">今日订单量</div>
            </template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>单</span>
            <!-- <span class="green">
              {{ usersOrderData.tOkRate }} %)
            </span> -->
          </div>
        </div>
        </el-col>
        <el-col :span="6">
          <el-statistic  :value="usersOrderData.tOkOrderQuantify" >
            <template #title>
            <div style="font-size: 14px; align-items: center">今日成单量</div>
            </template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>单</span>
            <span class="green">
              (成功率 {{ usersOrderData.tOkRate }} %)
            </span>
          </div>
        </div>
        </el-col>
      </el-row>
    </div>
        
  </div>
  
  <div class="gva-card-box">
    <div class="gva-card quick-entrance">
      <h3>昨日交易概况</h3>
      <el-row>
        <el-col :span="6">
          <el-statistic  :value="usersOrderData.yInCome" >
            <template #title>
            <div style="font-size: 14px; align-items: center">昨日收入</div>
            </template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>元(RMB)</span>
            <!-- <span class="green">
              元
              <el-icon>
                <CaretTop />
              </el-icon>
            </span> -->
          </div>
        </div>
        </el-col>
        <el-col :span="6">
          <el-statistic :value="usersOrderData.openChId">
            <template #title>
              <div style="font-size: 14px; display: inline-flex; align-items: center">
                通道概览
                <!-- <el-icon style="margin-left: 4px" :size="12">
                  <Male />
                </el-icon> -->
              </div>
            </template>
            <template #suffix>/{{ usersOrderData.chIdCnt }}</template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>开启个数/通道总数</span>
            <!-- <span class="green">
              
              <el-icon>
                <CaretTop />
              </el-icon>
            </span> -->
          </div>
        </div>
        </el-col>
        <el-col :span="6">
          <el-statistic  :value="usersOrderData.yOrderQuantify" >
            <template #title>
            <div style="font-size: 14px; align-items: center">昨日订单量</div>
            </template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>单</span>
            <!-- <span class="green">
              {{ usersOrderData.tOkRate }} %)
            </span> -->
          </div>
        </div>
        </el-col>
        <el-col :span="6">
          <el-statistic  :value="usersOrderData.yOkOrderQuantify" >
            <template #title>
            <div style="font-size: 14px; align-items: center">昨日成单量</div>
            </template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>单</span>
            <span class="green">
              (成功率 {{ usersOrderData.yOkRate }} %)
            </span>
          </div>
        </div>
        </el-col>
      </el-row>
    </div>
        
  </div>
  <div class="gva-card-box">

    <div class="gva-card quick-entrance">
      <div
            id="incomeEchart"
            class="dashboard-line"
          />
    </div>
  </div>
</div>


  </template>

<script>
export default {
  name: 'homePageDesc'
}
</script>

<script setup>
import { ChatLineRound, Male } from '@element-plus/icons-vue'
import {
  getSelectUserPayOrderAnalysisH,
  getHomePagePayOrderAnalysis,
  getVboxUserPayOrderAnalysisIncomeCharts,
} from '@/api/vboxPayOrder'
import { ref,shallowRef,onMounted} from 'vue'
import * as echarts from 'echarts'

const usersOrderData = ref({
          username: '',
          balance:0,
          chIdCnt:0,
          openChId:0,
          yOrderQuantify:0,
          yOkOrderQuantify:0,
          yOkRate:0,
          yInCome:0,
          tOrderQuantify:0,
          tOkOrderQuantify:0,
          tOkRate:0,
          tInCome:0,
        })

const getUsersOrderData = async() => {
  const table = await getHomePagePayOrderAnalysis()
  if (table.code === 0) {
    usersOrderData.value = table.data.resultData
    console.log("usersOrderData= " + JSON.stringify(usersOrderData.value))
  }
}
getUsersOrderData()


const usersOrderHData = ref({
          username: '',
          balance:0,
          chIdCnt:0,
          openChId:0,
          newOpenChId:0,
          yOrderQuantify:0,
          yOkOrderQuantify:0,
          yOkRate:0,
          yInCome:0,
          tOrderQuantify:0,
          tOkOrderQuantify:0,
          tOkRate:0,
          tInCome:0,
        })

const hColor = ref("")
const hPercentage = ref(0)
const getUsersOrderHData = async() => {
  const table = await getSelectUserPayOrderAnalysisH()
  if (table.code === 0) {
    usersOrderHData.value = table.data.resultData
    console.log("usersOrderHData= " + JSON.stringify(usersOrderHData.value))
    let per = 0
    if (usersOrderHData.value.yInCome != 0){
      per = (1.0 * usersOrderHData.value.tInCome - usersOrderHData.value.yInCome / usersOrderHData.value.yInCome * 100 )
    }
    hPercentage.value = per
    if (per > 0){
      hColor.value = "green"
    }else{
      hColor.value = "red"
    }
  }
}
getUsersOrderHData()


const incomeLineChartSeries = ref([])
const incomeLineChartXData = ref([])
const incomeLineChartLegend = ref([])
const chart = shallowRef(null)
const initChart = async() => {
  const table = await getVboxUserPayOrderAnalysisIncomeCharts()
  if (table.code === 0) {
    // console.log(JSON.stringify(table.data))
    // const res = JSON.parse(table.data)
    incomeLineChartSeries.value = table.data.lists
    incomeLineChartXData.value = table.data.xData
    incomeLineChartLegend.value = table.data.legendData
  }

  chart.value = echarts.init(document.getElementById("incomeEchart") /* 'macarons' */)
  setOptions()
}
// onMounted(async() => {
//   // await nextTick()
  
//   initChart()
// })
const setOptions = () => {
  chart.value.setOption({
  title: {
    text: '销售收入'
  },
  tooltip: {
    trigger: 'axis'
  },
  legend: {
    data: incomeLineChartLegend.value
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: '3%',
    containLabel: true
  },
  toolbox: {
    feature: {
      magicType: { show: true, type: ['stack', 'tiled'] },
      saveAsImage: {}
    }
  },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: incomeLineChartXData.value
  },
  yAxis: {
    type: 'value'
  },
  series: incomeLineChartSeries.value
})
}
initChart()
</script>

<style lang="scss" scoped>
.el-col {
  text-align: center;
}

.page {
    @apply p-0;
    .gva-card-box{
      @apply p-4;
      &+.gva-card-box{
        @apply pt-0;
      }
    }
    .gva-card {
      @apply box-border bg-white rounded h-auto px-6 py-8 overflow-hidden shadow-sm;
      .gva-card-title{
        @apply pb-5 border-t-0 border-l-0 border-r-0 border-b border-solid border-gray-100;
      }
    }
    .gva-top-card {
        @apply h-72 flex items-center justify-between text-gray-500;
        &-left {
          @apply h-full flex flex-col w-auto;
            &-title {
              @apply text-3xl text-gray-600;
            }
            &-dot {
              @apply mt-4 text-gray-600 text-lg;
            }
            &-item{
              +.gva-top-card-left-item{
                margin-top: 24px;
              }
              margin-top: 14px;
            }
        }
        &-right {
            height: 600px;
            width: 600px;
            margin-top: 28px;
        }
    }
     ::v-deep(.el-card__header){
          @apply p-0  border-gray-200;
        }
        .card-header{
          @apply pb-5 border-b border-solid border-gray-200 border-t-0 border-l-0 border-r-0;
        }
    .quick-entrance-items {
      @apply flex items-center justify-center text-center text-gray-800;
        .quick-entrance-item {
          @apply px-8 py-6 flex items-center flex-col transition-all duration-100 ease-in-out rounded-lg cursor-pointer;
          &:hover{
            @apply shadow-lg;
          }
            &-icon {
              @apply flex items-center h-16 w-16 rounded-lg justify-center mx-0 my-auto text-2xl;
            }
            p {
                @apply mt-2.5;
            }
        }
    }
}

.statistic-footer {
  // display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  font-size: 12px;
  color: var(--el-text-color-regular);
  margin-top: 16px;
}

.statistic-footer .footer-item {
  // display: flex;
  justify-content: space-between;
  align-items: center;
}

.statistic-footer .footer-item span:last-child {
  display: inline-flex;
  align-items: center;
  margin-left: 4px;
}

.custom-title {
  font-size: 12px; /* Adjust the font size as needed */
}

.green {
  color: var(--el-color-success);
}
.red {
  color: var(--el-color-error);
}
.dashboard-line {
  background-color: #fff;
  height: 500px;
  width: 1000px;
}
</style>
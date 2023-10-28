<template>
  <div class="ReclaimMileage-box">
    <div
        ref="echart"
        class="ReclaimMileage-box-echarts"
        :style="`width : ${chart?.clientWidth}px`"
    />

    <div class="ReclaimMileage-box-data">
      <div class="ReclaimMileage-box-data-left">
        <div>今日成单</div>
        <div class="number"><span style="color: #D3B379">￥</span> {{ orderFinishedTotalToday }} </div>
        <div>单数:{{ orderFinishedTotalToday }}</div>
      </div>
      <div class="ReclaimMileage-box-data-center">
        <div>今日总单</div>
        <div class="number"><span>-62% ↓</span> 3,067 </div>
        <div>单数:{{ orderTotalToday }}</div>
      </div>
      <div class="ReclaimMileage-box-data-right">
        <div>今日成率</div>
        <div class="number">1,298 <span>%</span></div>
        <div>{{ orderFinishedTotalToday }} / {{ orderTotalToday }}</div>
      </div>
    </div>

  </div>
</template>
<script setup>
import * as echarts from 'echarts'
import {nextTick, onMounted, onUnmounted, ref, watch} from 'vue'
import {getPayOrderList} from "@/api/payOrder";

const props = defineProps({
  channelCode: String
})

const orderTotalToday = ref(0)
const orderFinishedTotalToday = ref(0)

const orderView = async () => {
  let param = props.channelCode
  console.log(param)
  let res = await getPayOrderList({channelCode: param})
  console.log("----==")
  console.log(res)
  let total = 0
  let finished = 0
  if (res.code === 0) {
    let l = res.data.list
    for (let i = 0; i < l.length; i++) {
      if (l[i].orderStatus === 1) {
        finished++
      }
    }
    total = l.length
  }
  orderFinishedTotalToday.value = finished
  orderTotalToday.value = total
}
orderView()

const chart = ref(null)
const echart = ref(null)
const initChart = () => {
  chart.value = echarts.init(echart.value)
  setOptions()
  document.addEventListener('resize', () => {
    chart.value?.resize()
  })
}

const setOptions = () => {
  chart.value.setOption({
    backgroundColor: 'transparent',
    grid: {
      left: '0',
      right: '0',
      top: '60%',
      bottom: '0',
      containLabel: false
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'line',
        lineStyle: {
          color: 'rgba(50, 216, 205, 1)'
        },
      }
    },
    xAxis: [{
      type: 'category',
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        show: false
      }
    }],
    yAxis: [
      {
        splitLine: {
          show: false,
        },
        axisLine: {
          show: false
        },
        axisLabel: {
          show: false,
        },
        axisTick: {
          show: false
        }
      },

    ],
    series: [
      {
        name: '今日',
        type: 'line',
        smooth: true,
        stack: '总量',
        symbolSize: 5,
        showSymbol: false,
        itemStyle: {
          normal: {
            color: '#7723d0',
            lineStyle: {
              color: '#5423d0',
              width: 1
            },
          }
        },
        areaStyle: {
          normal: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              {
                offset: 0,
                color: 'rgb(141,50,216)'
              },
              {
                offset: 1,
                color: 'rgba(255, 255, 255, 0.2)'
              }
            ], false),
          }
        },
        data: [220, 182, 191, 234, 290, 330, 310, 201, 154, 190, 330, 410]
      },
    ]
  })
}

onMounted(() => {
  nextTick(() => {
    initChart()
  })
})

onUnmounted(() => {
  if (!chart.value) {
    return
  }
  chart.value.dispose()
  chart.value = null
})
</script>
<style lang="scss" scoped>

.ReclaimMileage-box {
  height: 120px;
  overflow: hidden;
  position: relative;

  &-echarts {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 2;
    width: 150%;
    height: 100%;
  }

  &-data {
    position: absolute;
    top: 0;
    left: 0;
    bottom: 10px;
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 10;
    user-select: none;
    pointer-events: none;

    &-left {
      height: 60%;
      display: flex;
      align-items: flex-start;
      flex-direction: column;
      justify-content: center;

      div {
        color: #999;
        font-size: 12px;
        margin-bottom: 8px;
      }
    }

    &-center {
      height: 60%;
      border-right: 1px solid #eee;
      border-left: 1px solid #eee;
      display: flex;
      align-items: flex-start;
      flex-direction: column;
      justify-content: center;
      padding: 0 10px;

      div {
        color: #999;
        font-size: 12px;
        margin-bottom: 8px;
      }
    }

    &-right {
      height: 60%;
      padding-left: 10px;
      display: flex;
      align-items: flex-start;
      flex-direction: column;
      justify-content: center;

      div {
        color: #999;
        font-size: 16px;
        margin-bottom: 8px;
      }
    }
  }
}

.in-line {
  --color: #5BC2A4;
}

.out-line {
  --color: #DF534E;
}

.number {
  color: #1d1d1f !important;
  font-size: 18px !important;
  font-weight: 500;

  span {
    font-size: 12px;
    color: #DF534E;
    font-weight: 400;
    transform: scale(0.8333);
  }
}
</style>

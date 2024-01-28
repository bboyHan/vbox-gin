<template>
  <div class="lineCharts-box">
    <div
      ref="echart"
      class="lineCharts-box-echarts"
      :style="`width : ${chart?.clientWidth}px`"
    />
  </div>
</template>
<script setup>
import * as echarts from 'echarts'
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import {formatTimeStr2Date, get5MinNearlyOneHour} from "@/utils/date";



const props = defineProps({
  uid: Number,
  channelCode: String,
  startTime: Date,
  endTime: Date,
  interval: String,
  keyword: String,
  format: String,
  unit: String,
  chartData: Object,
})


let legendData = [];
let seriesData = [];
let xAxisData = [];

const orderView = async () => {
  let res =  props.chartData
  // console.log('getAccIncomeSumData lineCharts res:', JSON.stringify(res))
  // debugger
  if (res.code === 0) {
    let resdata = res.data.chartData
    // console.log('getAccIncomeSumData ---->:', JSON.stringify(resdata))
 
    legendData = resdata.legendData
    // console.log('getAccIncomeSumData legendData:', legendData)
    xAxisData = resdata.xAxisData
    // console.log('getAccIncomeSumData xAxisData:', xAxisData)
    seriesData = resdata.seriesData
    // console.log('getAccIncomeSumData seriesData:', seriesData)
    initChart()
  }
  // // console.log('xData:', xData, 'yData', yData)
  
}



let chart = null
let echart = null
const initChart = () => {
  chart = echarts.init(echart)
  setOptions()
  document.addEventListener('resize', () => {
    chart?.resize()
  })
}


const generateEChartsOptions = (legendData, xAxisData, seriesData) => {
  return {
    tooltip: {
      show: true,
      trigger: 'axis',
      axisPointer: {
      type: 'line'
      } 
    },
    legend: {
      data: legendData,
      type: 'scroll',
    orient: 'vertical',
      right: 10,
    },
    grid: {
      left: '3%',
      right: '18%',
      bottom: '7%',
      containLabel: true
    },
    toolbox: {
      feature: {
        saveAsImage: {}
      }
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: xAxisData
    },
    yAxis: {
      name: '单位：' + props.unit,
      type: 'value'
    },
    series: seriesData
  };
};



// const options = generateEChartsOptions(legendData, xAxisData, seriesData);

const setOptions = () => {
  chart.setOption(generateEChartsOptions(legendData, xAxisData, seriesData))
}
orderView()

onMounted(() => {
  nextTick(() => {
    setTimeout(() => {
      orderView()
      initChart()
    }, 1000)
  })

  watch(props.chartData, (newVal, oldVal) => {
    console.log("发生了。。", newVal, oldVal)
    if (newVal) {
      orderView()
      chart.resize()
      console.log("resize了")
    }
  })
})

onUnmounted(() => {
  if (!chart) {
    return
  }
  chart.dispose()
  chart = null
})
</script>
<style lang="scss" scoped>
.lineCharts-box{
  height: 360px;
  overflow: hidden;
  position: relative;
  &-echarts{
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 2;
    width: 100%;
    height: 100%;
  }
}

.in-line{
  --color : #5BC2A4;
}
.out-line{
  --color: #DF534E;
}
</style>
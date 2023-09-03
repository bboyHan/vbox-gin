<template>
<el-tabs type="border-card">
    <el-tab-pane label="卡片模式">
      <el-descriptions
    class="margin-top"
    title="核销信息"
    :column="4"
    :size="size"
    border
  >
  <div v-for="(user, index) in users" :key="index">
    <!-- <template #extra> -->
      <!-- <el-button type="primary">Operation</el-button> -->
    <!-- </template> -->
    <!-- <el-divider content-position="left">{{ user.username }}</el-divider> -->
    <el-descriptions-item :span="4">
      <!-- <template #label>
        <div class="cell-item">
          <el-icon :style="iconStyle">
            <user />
          </el-icon>
          账户名
        </div>
      </template> -->
      <!-- {{ user.username }} -->
      <el-divider content-position="center">
        <el-icon><star-filled /></el-icon>
        {{ user.username }}
        <el-icon><star-filled /></el-icon>
      </el-divider>
    </el-descriptions-item>
    
    <el-descriptions-item :span="2">
      <template #label>
        <div class="cell-item">
          <el-icon :style="iconStyle">
            <user />
          </el-icon>
          账户名
        </div>
      </template>
      {{ user.username }}
    </el-descriptions-item>
    <el-descriptions-item :span="2">
      <template #label>
        <div class="cell-item">
          <el-icon :style="iconStyle">
            <Money />
          </el-icon>
          余额
        </div>
      </template>
      <el-tag type="danger" size="small">{{ user.balance }} 元</el-tag>
    </el-descriptions-item>
    <el-descriptions-item :span="2">
      <template #label>
        <div class="cell-item">
          <el-icon :style="iconStyle">
            <Shop />
          </el-icon>
          通道账号
        </div>
      </template>
      <el-tag size="small">{{ user.chIdCnt }}</el-tag>
    </el-descriptions-item>
    <el-descriptions-item :span="2">
      <template #label>
        <div class="cell-item">
          <el-icon :style="iconStyle">
            <Shop />
          </el-icon>
          开启账号
        </div>
      </template>
      <el-tag size="small">{{ user.openChId }}</el-tag>
    </el-descriptions-item>
    <el-descriptions-item>
      <template #label>
        <div class="cell-item">
          <el-icon :style="iconStyle">
            <DataAnalysis />
          </el-icon>
          昨日订单量
        </div>
      </template>
      <el-tag size="small">{{ user.yOrderQuantify }}</el-tag>
    </el-descriptions-item>
    <el-descriptions-item>
      <template #label>
        <div class="cell-item">
          <el-icon :style="iconStyle">
            <DataAnalysis />
          </el-icon>
          昨日成单量
        </div>
      </template>
      <el-tag size="small">{{ user.yOkOrderQuantify }}</el-tag>
    </el-descriptions-item>
    <el-descriptions-item>
      <template #label>
        <div class="cell-item">
          <el-icon :style="iconStyle">
            <PieChart />
          </el-icon>
          昨日成功率
        </div>
      </template>
      <el-tag type="success" size="small">{{ user.yOkRate }} %</el-tag>
    </el-descriptions-item>
    <el-descriptions-item>
      <template #label>
        <div class="cell-item">
          <el-icon :style="iconStyle">
            <Money />
          </el-icon>
          昨日收入
        </div>
      </template>
      <el-tag type="danger" size="small">{{ user.yInCome }} 元</el-tag>
    </el-descriptions-item>
    <el-descriptions-item>
      <template #label>
        <div class="cell-item">
          <el-icon :style="iconStyle">
            <DataAnalysis />
          </el-icon>
          今日订单量
        </div>
      </template>
      <el-tag size="small">{{ user.tOrderQuantify }}</el-tag>
    </el-descriptions-item>
    <el-descriptions-item>
      <template #label>
        <div class="cell-item">
          <el-icon :style="iconStyle">
            <DataAnalysis />
          </el-icon>
          今日成单量
        </div>
      </template>
      <el-tag size="small">{{ user.tOkOrderQuantify }}</el-tag>
    </el-descriptions-item>
    <el-descriptions-item>
      <template #label>
        <div class="cell-item">
          <el-icon :style="iconStyle">
            <PieChart />
          </el-icon>
          今日成功率
        </div>
      </template>
      <el-tag type="success" size="small">{{ user.tOkRate }} %</el-tag>
    </el-descriptions-item>
    <el-descriptions-item>
      <template #label>
        <div class="cell-item">
          <el-icon :style="iconStyle">
            <Money />
          </el-icon>
          今日收入
        </div>
      </template>
      <el-tag type="danger" size="small">{{ user.tInCome }} 元</el-tag>
    </el-descriptions-item>
  </div>
  </el-descriptions>

    </el-tab-pane>
    <el-tab-pane label="图形模式">
        <div >
          <!-- <div class="dashboard-line-title">
            收入趋势
          </div> -->
          <div
            id="incomeEchart"
            class="dashboard-line"
          />
        </div>
    </el-tab-pane>
  </el-tabs>
 

</template>

<script>
export default {
  name: 'OrderPayDesc'
}
</script>

<script setup>
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive,watch ,computed,shallowRef,nextTick,onMounted} from 'vue'
import * as echarts from 'echarts'

import {
  Iphone,
  Location,
  Money,
  OfficeBuilding,
  PieChart,
  Tickets,
  User,
} from '@element-plus/icons-vue'

import {
  getVboxUserPayOrderAnalysis
} from '@/api/vboxPayOrder'


const total = ref(0)
const users = ref([])
const size = ref('')
const iconStyle = computed(() => {
  const marginMap = {
    large: '8px',
    default: '6px',
    small: '4px',
  }
  return {
    marginRight: marginMap[size.value] || marginMap.default,
  }
})
const blockMargin = computed(() => {
  const marginMap = {
    large: '32px',
    default: '28px',
    small: '24px',
  }
  return {
    marginTop: marginMap[size.value] || marginMap.default,
  }
})


const getTableData = async() => {
    const table = await getVboxUserPayOrderAnalysis()
    if (table.code === 0) {
        users.value = table.data.list
        console.log('payOrderAnalysis=' + JSON.stringify(users.value))
        total.value = table.data.total
    }
}

getTableData()



const chart = shallowRef(null)
// const incomeEchart = ref(null)
const initChart = () => {
  chart.value = echarts.init(document.getElementById("incomeEchart") /* 'macarons' */)
  setOptions()
}
onMounted(async() => {
  // await nextTick()
  initChart()
})
const setOptions = () => {
  chart.value.setOption({
  title: {
    text: '销售收入'
  },
  tooltip: {
    trigger: 'axis'
  },
  legend: {
    data: ['xiaoming', 'xiaobai', 'xiaohei', 'xiaohei1', 'xiaowang']
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: '3%',
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
    data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
  },
  yAxis: {
    type: 'value'
  },
  series: [
    {
      name: 'xiaoming',
      type: 'line',
      stack: 'Total',
      data: [120, 132, 101, 134, 90, 230, 210]
    },
    {
      name: 'xiaobai',
      type: 'line',
      stack: 'Total',
      data: [220, 182, 191, 234, 290, 330, 310]
    },
    {
      name: 'xiaohei',
      type: 'line',
      stack: 'Total',
      data: [150, 232, 201, 154, 190, 330, 410]
    },
    {
      name: 'xiaohei1',
      type: 'line',
      stack: 'Total',
      data: [320, 332, 301, 334, 390, 330, 320]
    },
    {
      name: 'xiaowang',
      type: 'line',
      stack: 'Total',
      data: [820, 932, 901, 934, 1290, 1330, 1320]
    }
  ]
})
}



</script>

<style>

.el-descriptions {
  margin-top: 20px;
}
.cell-item {
  display: flex;
  align-items: center;
}
.margin-top {
  margin-top: 20px;
}

.dashboard-line {
  background-color: #fff;
  height: 500px;
  width: 1000px;
}
.dashboard-line-title {
  font-weight: 600;
  margin-bottom: 14px;
}

</style>
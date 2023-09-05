<template>
<el-tabs type="border-card">
  


    <el-tab-pane label="卡片模式" v-model="activeTab">
      <div class="app-container">
        <el-row :gutter="20">
          <el-col :span="4" :xs="24">
        <div class="head-container">
          <el-autocomplete 
          v-model="usersFormData.name" 
          :fetch-suggestions="querySearchAsync" 
          placeholder="请输入" 
          :clearable="true"
          @select="handleSelect">
            <template #default="{ item }">
              {{ item.value }}
            </template>
          </el-autocomplete>
        </div>
        <div class="head-container">
          <div v-for="name in usersItem" :key="name.id" class="name-item" @click="handleItemClick(name)">
          {{ name.value }}
        </div>
          <!-- <el-tree
            :data="deptOptions"
            :props="defaultProps"
            :expand-on-click-node="false"
            :filter-node-method="filterNode"
            ref="tree"
            node-key="id"
            default-expand-all
            highlight-current
            @node-click="handleNodeClick"
          /> -->
        </div>
      </el-col>
      <el-col :span="20" :xs="24">
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
    </el-col>
    </el-row>
    </div>
    <div>
      <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" width="80%" destroy-on-close>
      
      <el-descriptions
        class="margin-top"
        title="个人核销信息"
        :column="4"
        :size="size"
        border
      >
      <div v-for="(user, index) in selectUsers" :key="index">
        <!-- <template #extra> -->
          <!-- <el-button type="primary">Operation</el-button> -->
        <!-- </template> -->
        <!-- <el-divider content-position="left">{{ user.username }}</el-divider> -->
        <el-descriptions-item :span="4">
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
      <div
          id="quantifyEchart"
          class="dashboard-line"
        />
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">关 闭</el-button>
          <!-- <el-button type="primary" @click="enterDialog">确 定</el-button> -->
        </div>
      </template>
    </el-dialog>
    </div>
    </el-tab-pane>


    <el-tab-pane label="图形模式" v-model="activeTab">
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
import {
  getOwnerUserListForSelect
} from '@/api/user'

import {
  getVboxUserPayOrderAnalysis,
  getSelectUserPayOrderAnalysis,
  getVboxUserPayOrderAnalysisIncomeCharts,
  getSelectPayOrderAnalysisQuantifyCharts
} from '@/api/vboxPayOrder'

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

import { stringify } from 'qs'


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
        // console.log('payOrderAnalysis=' + JSON.stringify(users.value))
        total.value = table.data.total
    }
}

getTableData()


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


const activeTab = ref('图形模式');

    // 监听标签页的变化
watch(activeTab, (newTab, oldTab) => {
  if (newTab === '图形模式') {
    // 执行数据初始化的操作
    initChart();
  }
  if (newTab === '卡片模式') {
    // 执行数据初始化的操作
    getTableData();
  }
  
});


const usersFormData = ref({
        name: ''
        })
const usersItem = ref([])
// 可搜索店铺
const loadAll = async ()  => {

  const res = await getOwnerUserListForSelect({ page: 1, pageSize: 100 })
  // console.log('== res ==>' + JSON.stringify(res))
  if (res.code === 0) {
    // console.log('== res.data.marks==>' + JSON.stringify(res.data.list))
    usersItem.value = res.data.list
  }
  return usersItem.value
}
loadAll()
let timeout
const querySearchAsync = (queryString, cb) => {
  loadAll()
  const results = queryString
    ? usersItem.value.filter(createFilter(queryString))
    : usersItem.value

  clearTimeout(timeout)
  timeout = setTimeout(() => {
    cb(results)
  }, 2000 * Math.random())
}


const createFilter = (queryString) => {
  return (restaurant) => {
    return (
      restaurant.value.toLowerCase().indexOf(queryString.toLowerCase()) === 0
    )
  }
}

const handleSelect = (item) => {
  console.log(item)
  usersFormData.value.name = item.value
}

onMounted(() => {
  console.log('onMounted')
  usersItem.value = loadAll()
})


const selectUserName = ref('')
// 点击名字时的操作
const handleItemClick = (item) => {
      // 在这里处理点击名字后的操作
      console.log(item.value); // 举例：打印选中的名字
      selectUserName.value = item.value
      dialogFormVisible.value = true
      getDialogTable()
      initQuantifyChart()
};


const selectTotal = ref(0)
const selectUsers = ref([])


const getDialogTable = async() => {
  const table = await getSelectUserPayOrderAnalysis({ Username: selectUserName.value })
  if (table.code === 0) {
    selectUsers.value = table.data.list
    selectTotal.value = table.data.total
  }
}


// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
    dialogFormVisible.value = true
    getDialogTable()
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    selectUsers.value = []
    selectTotal.value = 0
}





const quantifyLineChartSeries = ref([])
const quantifyLineChartXData = ref([])
const quantifyLineChartLegend = ref([])
const quantifyChart = shallowRef(null)
const initQuantifyChart = async() => {
  const table = await getSelectPayOrderAnalysisQuantifyCharts({ Username: selectUserName.value })
  if (table.code === 0) {
    // console.log(JSON.stringify(table.data))
    // const res = JSON.parse(table.data)
    quantifyLineChartSeries.value = table.data.lists
    quantifyLineChartXData.value = table.data.xData
    quantifyLineChartLegend.value = table.data.legendData
  }

  quantifyChart.value = echarts.init(document.getElementById("quantifyEchart") /* 'macarons' */)
  setQuantifyOptions()
}

const setQuantifyOptions = () => {
  quantifyChart.value.setOption({
  title: {
    text: '各个通道成单量'
  },
  tooltip: {
    trigger: 'axis'
  },
  legend: {
    data: quantifyLineChartLegend.value
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
    data: quantifyLineChartXData.value
  },
  yAxis: {
    type: 'value'
  },
  series: quantifyLineChartSeries.value
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
.name-item {
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  padding: 10px;
  /* background-color: #f5f5f5; */
  /* border-radius: 20px; */
  margin-right: 1px;
  margin-bottom: 1px;
  transition: background-color 0.2s ease-in-out;
}
#quantifyEchart{
  margin-top: 40px;
}
</style>
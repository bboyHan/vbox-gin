<template>

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

</template>

<script>
export default {
  name: 'OrderPayDesc'
}
</script>

<script setup>
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive,watch ,computed} from 'vue'
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
</style>
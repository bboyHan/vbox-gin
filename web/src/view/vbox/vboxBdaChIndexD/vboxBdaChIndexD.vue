<template>
  <div class="gva-search-box">
    <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
      <el-form-item label="创建日期" prop="dt">
      <template #label>
        <span>
           选择日期
          <!-- <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
            <el-icon><QuestionFilled /></el-icon>
          </el-tooltip> -->
        </span>
      </template>
      <el-date-picker  v-model="searchInfo.startCreatedAt" type="date" format="YYYY-MM-DD" placeholder="日期" :disabled-date="time=> searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"></el-date-picker>
       <!-- — -->
      <!-- <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束日期" :disabled-date="time=> searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"></el-date-picker> -->
      </el-form-item>

    </el-form>
  </div>
     <div class="gva-table-box">
        
        <el-table
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        border
        >
        <el-table-column type="selection" width="55" />
        <!-- <el-table-column align="left" label="日期" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column> -->

        <el-table-column align="left" label="日期" prop="dt" width="120" />
        <!-- <el-table-column align="left" label="uid" prop="uid" width="120" /> -->
        <el-table-column align="left" label="用户名" prop="username" width="120" />
        <!-- <el-table-column align="left" label="付方账户名" prop="p_account" width="120" /> -->
        <el-table-column align="left" label="通道code" prop="channelCode" width="120" />
        <!-- <el-table-column align="left" label="通道id" prop="productId" width="120" /> -->
        <el-table-column align="left" label="通道名称" prop="productName" width="120" />
        <el-table-column align="left" label="订单数量" prop="orderQuantify" width="120" />
        <el-table-column align="left" label="订单成交数量" prop="okOrderQuantify" width="120" />
        <el-table-column align="left" label="成交率" prop="ratio" width="120" >
          <template #default="scope">{{ scope.row.ratio }} %</template>
        </el-table-column>
        <el-table-column align="left" label="成交金额" prop="income" width="120" />
     
        </el-table>
        <div class="gva-pagination">
            <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
            />
        </div>
    </div>

</template>

<script setup>
import {
  findVboxBdaChaccIndexD,
  getVboxBdaChaccIndexDList
} from '@/api/vboxBdaChaccIndexD'
import {
  findVboxBdaChIndexD,
  getVboxBdaChIndexDList
} from '@/api/vboxBdaChIndexD'
import { reactive, computed,ref } from 'vue';
import { getDictFunc, formatDate, formatBoolean, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const searchInfo = ref({})
const tableData = ref([])

const totalData = ref({
        uid: 0,
        username: '',
        p_account: '',
        // channelCode: 0,
        productId: 0,
        productName: '',
        orderQuantify: 0,
        okOrderQuantify: 0,
        ratio: 0,
        income: 0,
        dt: '',
        })
// const tableData = reactive([
//   { name: 'John', age: 25, score: 80 },
//   { name: 'Alice', age: 28, score: 90 },
//   { name: 'Tom', age: 22, score: 95 },
//   { name: 'Tom2', age: 22, score: 95 },
//   { name: 'Tom3', age: 13, score: 95 },
//   // 更多数据行...
// ]);
const columns = reactive([
  { label: '日期', prop: 'dt' },
  { label: '用户', prop: 'username' },
  { label: '通道名', prop: 'productName' },
  { label: '订单量', prop: 'orderQuantify' },
  { label: '成交单量', prop: 'okOrderQuantify' },
  { label: '成交率', prop: 'ratio' },
  { label: '成交金额', prop: 'income' },
]);

// 计算总和行的数据
// const totalRow = computed(() => {
//   const total = { name: 'Total' };
//   for (let i = 0; i < columns.length; i++) {
//     const column = columns[i];
//     let sum = 0;
//     for (let j = 0; j < tableData.length; j++) {
//       sum += tableData[j][column.prop];
//     }
//     total[column.prop] = sum;
//   }
//   return [total];
// });


// 验证规则
const rule = reactive({
               dt : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              }
              ],
})

const searchRule = reactive({
  createdAt: [
    { validator: (rule, value, callback) => {
      if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
        callback(new Error('请填写结束日期'))
      } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
        callback(new Error('请填写开始日期'))
      } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
        callback(new Error('开始日期应当早于结束日期'))
      } else {
        callback()
      }
    }, trigger: 'change' }
  ],
})

const elFormRef = ref()
const elSearchFormRef = ref()

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}
const getTableData = async() => {
  const table = await getVboxBdaChIndexDList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    console.log(JSON.stringify(table.data.list))
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
    let sumOrderQuantify = 0
    let sumOkOrderQuantify = 0
    let sumIncome = 0
    let ratio = 0
    let cnt = 0
    tableData.value.forEach(item => {
      cnt ++
      sumOkOrderQuantify += item.okOrderQuantify
      sumOrderQuantify += item.orderQuantify
      sumIncome += item.income
      ratio +=item.ratio
    })
    totalData.value.orderQuantify = sumOrderQuantify
    totalData.value.okOrderQuantify = sumOkOrderQuantify
    totalData.value.income = sumIncome
    totalData.value.ratio = ratio / cnt
    totalData.value.productName = "合计"

    tableData.value.push(totalData.value)


  }
}
getTableData()

const getTableTotalData = async() => {
     let sumOrderQuantify = 0
     let sumOkOrderQuantify = 0
     let sumIncome = 0
      tableData.value.forEach(item => {
        sumOkOrderQuantify += item.okOrderQuantify
        sumOrderQuantify += item.orderQuantify
        sumIncome += item.income
      })
      totalData.value.orderQuantify = sumOrderQuantify
      totalData.value.okOrderQuantify = sumOkOrderQuantify
      totalData.value.income = sumIncome
      totalData.value.productName = "合计"

}
</script>

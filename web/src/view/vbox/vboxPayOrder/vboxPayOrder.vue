<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
      <el-form-item label="创建日期" prop="createdAt">
      <template #label>
        <span>
          创建日期
          <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
            <el-icon><QuestionFilled /></el-icon>
          </el-tooltip>
        </span>
      </template>
      <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始日期" :disabled-date="time=> searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"></el-date-picker>
       —
      <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束日期" :disabled-date="time=> searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"></el-date-picker>
      </el-form-item >
        <el-form-item label="付方订单ID" prop="order_id">
         <el-input v-model="searchInfo.order_id" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="平台订单ID" prop="platform_oid">
          <el-input v-model="searchInfo.platform_oid" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="通道帐号" prop="acAccount">
         <el-input v-model="searchInfo.acAccount" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="帐号备注" prop="acRemark">
         <el-input v-model="searchInfo.acRemark" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="所属通道" prop="c_channel_id">
         <el-input v-model="searchInfo.c_channel_id" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="订单状态" prop="order_status">
          <el-input v-model.number="searchInfo.order_status" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="通知状态" prop="callback_status">
          <el-input v-model.number="searchInfo.callback_status" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        border
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="归属用户" prop="uid" width="240">
          <template #default="scope">
            <el-row>
              <el-col :span="24">
                <el-text tag="b">付方ID: </el-text>
                <el-text>{{ scope.row.p_account }}</el-text>
              </el-col>
              <el-col :span="24">
                <el-text tag="b">归属用户: </el-text>
                <el-text>{{ scope.row.uid }}</el-text>
              </el-col>
            </el-row>
          </template>
        </el-table-column>
        <el-table-column align="left" label="充值账号" prop="uid" width="240">
          <template #default="scope">
            <el-row>
              <el-col :span="24">
                <el-text tag="b">通道账号: </el-text>
                <el-tag>{{ scope.row.acAccount }}</el-tag>
              </el-col>
              <el-col :span="24">
                <el-text tag="b">账号备注: </el-text>
                <el-tag>{{ scope.row.acRemark }}</el-tag>
              </el-col>
            </el-row>
          </template>
        </el-table-column>
        <el-table-column align="left" label="所属通道" prop="channel_code" width="120">
          <template #default="scope">
            <el-tag effect="light" round>{{ scope.row.channel_code }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="订单ID" prop="order_id" width="420">
          <template #default="scope">
            <el-row>
              <el-col :span="24">
                <el-text tag="b">上游订单ID: </el-text>
                <el-tag>{{ scope.row.order_id }}</el-tag>
              </el-col>
              <el-col :span="24">
                <el-text tag="b">平台订单ID: </el-text>
                <el-tag>{{ scope.row.platform_oid }}</el-tag>
              </el-col>
            </el-row>
          </template>
        </el-table-column>
        <el-table-column align="left" label="金额" prop="money" width="80" />
        <el-table-column align="left" label="订单状态" prop="order_status" width="180">
          <template #default="scope">
            <el-row>
              <el-col :span="24">
                <el-tag effect="dark" :type="formatOrderStatusColor(scope.row.order_status)">{{ formatOrderStatus(scope.row.order_status) }}</el-tag>
              </el-col>
              <el-col :span="24">
                <el-tag effect="dark" :type="formatCallbackStatusColor(scope.row.callback_status)">{{ formatCallbackOrderStatus(scope.row.callback_status) }}</el-tag>
              </el-col>
            </el-row>
          </template>
        </el-table-column>
         <el-table-column align="left" label="创建时间" width="180">
            <template #default="scope">{{ formatDate(scope.row.create_time) }}</template>
         </el-table-column>
         <el-table-column align="left" label="通知时间" width="180">
            <template #default="scope">{{ formatDate(scope.row.call_time) }}</template>
         </el-table-column>
        <el-table-column align="left" label="客户端信息" prop="pay_ip" width="240">
          <template #default="scope">
            <el-row>
              <el-col :span="24">
                <el-tag>设备: {{ scope.row.pay_device }}</el-tag>
              </el-col>
              <el-col :span="24">
                <el-tag>IP: {{ scope.row.pay_ip }}</el-tag>
              </el-col>
              <el-col :span="24">
                <el-tag>区域: {{ scope.row.pay_region }}</el-tag>
              </el-col>
            </el-row>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作">
            <template #default="scope">
              <el-button type="success" @click="queryHisRecords(scope.row)">查询</el-button>
              <el-button type="primary" @click="queryHisRecords(scope.row)">补单</el-button>
            </template>
        </el-table-column>
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
  </div>
</template>

<script>
export default {
  name: 'VboxPayOrder'
}
</script>

<script setup>
import {
  findVboxPayOrder,
  getVboxPayOrderList
} from '@/api/vboxPayOrder'
// 全量引入格式化工具 请按需保留
import {
  formatCallbackOrderStatus,
  formatOrderStatus,
  formatDate,
  formatBoolean,
  filterDict,
  formatOrderStatusColor,
  formatCallbackStatusColor
} from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { useUserStore } from '@/pinia/modules/user';
import { QuestionFilled } from "@element-plus/icons-vue";

const userStore = useUserStore()

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        order_id: '',
        p_account: '',
        cost: 0,
        uid: 0,
        ac_id: '',
        channel_code: '',
        platform_oid: '',
        pay_device: '',
        pay_ip: '',
        resource_url:'',
        pay_region: '',
        notify_url: '',
        order_status: 0,
        callback_status: 0,
        code_use_status: 0,
        create_time: new Date(),
        async_time: new Date(),
        call_time: new Date(),
        })

// 验证规则
const rule = reactive({
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

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

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

// 查询
const getTableData = async() => {
  await userStore.LoadAllUser();
  await userStore.GetChannelProductList();
  const table = await getVboxPayOrderList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  console.log(userStore.chanCodeMap)
  console.log(userStore.ownerUsersMap)

  if (table.code === 0) {
    // tableData.value = table.data.list
    tableData.value = table.data.list.map(item => {
      // 将每个元素的c_channel_id替换为对应的productName值
      const productName = userStore.chanCodeMap.get(item.c_channel_id + '');
      const userName = userStore.ownerUsersMap.get(item.user_id);
      return { ...item,
        uid: userName ? userName : '- | ' + item.uid,
        c_channel_id: productName ? productName : '- | ' + item.c_channel_id
      };
    })
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
}

// 获取需要的字典 可能为空 按需保留
setOptions()

// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
    multipleSelection.value = val
}

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        order_id: '',
        p_account: '',
        cost: 0,
        uid: 0,
        ac_id: '',
        channel_code: '',
        platform_oid: '',
        pay_device: '',
        pay_ip: '',
        pay_region: '',
        notify_url: '',
        order_status: 0,
        callback_status: 0,
        code_use_status: 0,
        create_time: new Date(),
        async_time: new Date(),
        call_time: new Date(),
      }
}
// 弹窗确定
const enterDialog = async () => {
   elFormRef.value?.validate( async (valid) => {
     if (!valid) return
      let res
      switch (type.value) {
        case 'create':
          break
        case 'update':
          break
        default:
          break
      }
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '创建/更改成功'
        })
        closeDialog()
        getTableData()
      }
  })
}

</script>

<style>
.limit-height {
  max-height: 20px;
}
</style>

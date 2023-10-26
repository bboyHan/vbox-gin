<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
        <el-form-item label="付方单号" prop="orderId">
          <el-input v-model="searchInfo.orderId" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button icon="refresh" @click="resetSimple(true)">简约版</el-button>
          <el-button icon="refresh" @click="resetSimple(false)">详情版</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <!--   简约版   -->
      <el-table
          v-if="isSimple"
          ref="multipleTable"
          style="width: 100%"
          tooltip-effect="dark"
          :data="tableData"
          row-key="ID"
          border
          @selection-change="handleSelectionChange"
      >
        <el-table-column align="center" label="账号ID" prop="acId" width="180" >
          <template #default="scope">
            <div v-if="isPendingAcc(scope.row)">
              {{ scope.row.acId }}
              <el-button type="primary" link class="table-button" @click="getAccDetails(scope.row)">
                <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>
              </el-button>
            </div>
            <div v-else>
              <el-button type="info" :loading-icon="Eleme" loading link class="table-button">匹配中</el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="center" label="订单ID" prop="orderId" width="220" />
        <el-table-column align="center" label="金额" prop="money" width="120" />
        <el-table-column align="center" label="订单状态" prop="orderStatus" width="120">
          <template #default="scope">
            <el-button style="width: 90px" :color="formatPayedColor(scope.row.orderStatus, scope.row.acId)">{{ formatPayed(scope.row.orderStatus, scope.row.acId) }}</el-button>
          </template>
        </el-table-column>
        <el-table-column align="center" label="回调状态" prop="callbackStatus" width="120">
          <template #default="scope">
            <el-button style="width: 90px" :color="formatNotifyColor(scope.row.callbackStatus)">{{ formatNotify(scope.row.callbackStatus) }}</el-button>
          </template>
        </el-table-column>
        <el-table-column align="center" label="创建时间" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
              <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>
              详情
            </el-button>
            <el-button type="primary" link class="table-button" @click="notifyPayOrder(scope.row)">
              <el-icon style="margin-right: 5px"><Position /></el-icon>
              补单
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!--   详情版   -->
      <el-table
        v-else
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        border
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="日期" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="付方ID" prop="pAccount" width="120">
          <template #default="scope">
            {{ scope.row.pAccount }}
            <el-button type="primary" link class="table-button" @click="getPADetails(scope.row.pAccount)">
              <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>
              详情
            </el-button>
          </template>
        </el-table-column>
        <el-table-column align="left" label="单价积分" prop="unitPrice" width="120" />
        <el-table-column align="left" label="用户ID" prop="uid" width="120" />
        <el-table-column align="left" label="通道编码" prop="channelCode" width="120" />
        <el-table-column align="left" label="平台id" prop="platformOid" width="220" />
        <el-table-column align="left" label="客户ip" prop="payIp" width="120" />
        <el-table-column align="left" label="区域" prop="payRegion" width="120" />
        <el-table-column align="left" label="客户端设备" prop="payDevice" width="120" />
        <el-table-column align="left" label="支付链接" prop="resourceUrl" width="200"/>
        <el-table-column align="left" label="回调地址" prop="notifyUrl" width="120" />
        <el-table-column align="left" label="账号ID" prop="acId" width="180" >
          <template #default="scope">
            <div v-if="isPendingAcc(scope.row)">
              {{ scope.row.acId }}
              <el-button type="primary" link class="table-button" @click="getAccDetails(scope.row)">
                <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>
              </el-button>
            </div>
            <div v-else>
              <el-button type="info" :loading-icon="Eleme" loading>匹配中</el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="订单ID" prop="orderId" width="220" />
        <el-table-column align="left" label="金额" prop="money" width="120" />
        <el-table-column align="left" label="订单状态" prop="orderStatus" width="120">
          <template #default="scope">
            <el-button style="width: 90px" :color="formatPayedColor(scope.row.orderStatus, scope.row.acId)">{{ formatPayed(scope.row.orderStatus, scope.row.acId) }}</el-button>
          </template>
        </el-table-column>
        <el-table-column align="left" label="回调状态" prop="callbackStatus" width="120">
          <template #default="scope">
            <el-button style="width: 90px" :color="formatNotifyColor(scope.row.callbackStatus)">{{ formatNotify(scope.row.callbackStatus) }}</el-button>
          </template>
        </el-table-column>
                <el-table-column align="left" label="取码状态" prop="codeUseStatus" width="120">
                  <template #default="scope">{{ formatBoolean(scope.row.codeUseStatus) }}</template>
                </el-table-column>
                <el-table-column align="left" label="异步执行时间" width="180">
                  <template #default="scope">{{ formatDate(scope.row.asyncTime) }}</template>
                </el-table-column>
        <el-table-column align="left" label="回调时间" width="180">
          <template #default="scope">{{ formatDate(scope.row.callTime) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
              <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>
              详情
            </el-button>
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

    <!-- 订单查看详情 -->
    <el-dialog v-model="detailShow" style="width: 800px" lock-scroll :before-close="closeDetailShow" title="查看详情" destroy-on-close>
      <el-scrollbar height="550px">
        <el-descriptions column="1" border>
          <el-descriptions-item label="订单ID">{{ formData.orderId }}</el-descriptions-item>
          <el-descriptions-item label="付方ID">{{ formData.pAccount }}</el-descriptions-item>
          <el-descriptions-item label="金额">{{ formData.money }}</el-descriptions-item>
          <el-descriptions-item label="单价积分">{{ formData.unitPrice }}</el-descriptions-item>
          <el-descriptions-item label="用户ID">{{ formData.uid }}</el-descriptions-item>
          <el-descriptions-item label="账号ID">{{ formData.acId }}</el-descriptions-item>
          <el-descriptions-item label="通道编码">{{ formData.channelCode }}</el-descriptions-item>
          <el-descriptions-item label="平台id">{{ formData.platformOid }}</el-descriptions-item>
          <el-descriptions-item label="客户ip">{{ formData.payIp }}</el-descriptions-item>
          <el-descriptions-item label="区域">{{ formData.payRegion }}</el-descriptions-item>
          <el-descriptions-item label="客户端设备">{{ formData.payDevice }}</el-descriptions-item>
          <el-descriptions-item label="支付链接">{{ formData.resourceUrl }}</el-descriptions-item>
          <el-descriptions-item label="回调地址">{{ formData.notifyUrl }}</el-descriptions-item>
          <el-descriptions-item label="订单状态">{{ formatBoolean(formData.orderStatus) }}</el-descriptions-item>
          <el-descriptions-item label="回调状态">{{ formatBoolean(formData.callbackStatus) }}</el-descriptions-item>
          <el-descriptions-item label="取码状态">{{ formatBoolean(formData.codeUseStatus) }}</el-descriptions-item>
          <el-descriptions-item label="异步执行时间">{{ formatDate(formData.asyncTime) }}</el-descriptions-item>
          <el-descriptions-item label="回调时间">{{ formatDate(formData.callTime) }}</el-descriptions-item>
        </el-descriptions>
      </el-scrollbar>
    </el-dialog>

    <!--  补单  -->
    <el-dialog
        v-model="dialogFormVisible"
        :before-close="closeDialog"
        :title="typeTitle"
        destroy-on-close
        style="width: 450px"
    >
      <el-scrollbar height="100px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="120px">
          <el-form-item label="订单ID" prop="authCaptcha">
            <el-input disabled v-model="formData.orderId" :clearable="true" placeholder="请输入" style="width: 80%"/>
          </el-form-item>
          <el-form-item label="防爆验证码" prop="authCaptcha">
            <el-input v-model="formData.authCaptcha" :clearable="true" placeholder="请输入防爆验证码" style="width: 80%"/>
          </el-form-item>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 账号查看详情 -->
    <el-dialog v-model="detailAccShow" style="width: 800px" lock-scroll :before-close="closeAccDetailShow" title="查看详情" destroy-on-close>
      <el-scrollbar height="550px">
        <el-descriptions column="6" border>
          <el-descriptions-item label="用户id" :span="6">{{ formAccData.uid }}</el-descriptions-item>
          <el-descriptions-item label="账户ID" :span="6">{{ formAccData.acId }}</el-descriptions-item>
          <el-descriptions-item label="账户备注" :span="6">{{ formAccData.acRemark }}</el-descriptions-item>
          <el-descriptions-item label="通道账户" :span="3">{{ formAccData.acAccount }}</el-descriptions-item>
          <el-descriptions-item label="账户密码" :span="3">{{ formAccData.acPwd }}</el-descriptions-item>
          <el-descriptions-item label="ck" :span="6">{{ formAccData.token }}</el-descriptions-item>
          <el-descriptions-item label="通道id" :span="6">{{ formAccData.cid }}</el-descriptions-item>
          <el-descriptions-item label="笔数限制" :span="2">{{ formAccData.countLimit }}</el-descriptions-item>
          <el-descriptions-item label="日限额" :span="2">{{ formAccData.dailyLimit }}</el-descriptions-item>
          <el-descriptions-item label="总限额" :span="2">{{ formAccData.totalLimit }}</el-descriptions-item>
          <el-descriptions-item label="状态开关" :span="3">{{ formAccData.status===0?'关闭':'开启' }}</el-descriptions-item>
          <el-descriptions-item label="系统开关" :span="3">{{ formAccData.sysStatus===0?'关闭':'开启' }}</el-descriptions-item>
        </el-descriptions>
      </el-scrollbar>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createPayOrder,
  findPayOrder,
  getPayOrderList
} from '@/api/payOrder'
import {
  findChannelAccount,
} from '@/api/channelAccount'
// 全量引入格式化工具 请按需保留
import {
  getDictFunc,
  formatDate,
  formatBoolean,
  filterDict,
  ReturnArrImg,
  onDownloadFile,
  formatNotify, formatPayed, formatPayedColor, formatNotifyColor
} from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import {Eleme, InfoFilled, Position} from "@element-plus/icons-vue";

defineOptions({
  name: 'PayOrder'
})

//页面简约切换
const isSimple = ref(true)

// 获取付方账号详情
const getPADetails = (paID) => {
  console.log("查当前详情付方账号：" +paID)
}

// 重置
const resetSimple = (status) => {
  isSimple.value = status
}

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  authCaptcha: '',
  orderId: '',
  pAccount: '',
  money: 0,
  unitPrice: 0,
  uid: 0,
  acId: '',
  channelCode: '',
  platformOid: '',
  payIp: '',
  payRegion: '',
  payDevice: '',
  resourceUrl: '',
  notifyUrl: '',
  orderStatus: false,
  callbackStatus: false,
  codeUseStatus: false,
  asyncTime: new Date(),
  callTime: new Date(),
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
    if (searchInfo.value.orderStatus === ""){
      searchInfo.value.orderStatus=null
    }
    if (searchInfo.value.callbackStatus === ""){
      searchInfo.value.callbackStatus=null
    }
    if (searchInfo.value.codeUseStatus === ""){
      searchInfo.value.codeUseStatus=null
    }
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
  const table = await getPayOrderList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
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

// 删除行
const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    deletePayOrderFunc(row)
  })
}


// 批量删除控制标记
const deleteVisible = ref(false)

// 多选删除
const onDelete = async() => {
  const ids = []
  if (multipleSelection.value.length === 0) {
    ElMessage({
      type: 'warning',
      message: '请选择要删除的数据'
    })
    return
  }
  multipleSelection.value &&
  multipleSelection.value.map(item => {
    ids.push(item.ID)
  })
  const res = await deletePayOrderByIds({ ids })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '删除成功'
    })
    if (tableData.value.length === ids.length && page.value > 1) {
      page.value--
    }
    deleteVisible.value = false
    getTableData()
  }
}

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')
const typeTitle = ref('')


// 删除行
const deletePayOrderFunc = async (row) => {
  const res = await deletePayOrder({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '删除成功'
    })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)


// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findPayOrder({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.repayOrder
    openDetailShow()
  }
}

// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
    orderId: '',
    pAccount: '',
    money: 0,
    unitPrice: 0,
    uid: 0,
    acId: '',
    channelCode: '',
    platformOid: '',
    payIp: '',
    payRegion: '',
    payDevice: '',
    notifyUrl: '',
    orderStatus: false,
    callbackStatus: false,
    codeUseStatus: false,
    asyncTime: new Date(),
    callTime: new Date(),
  }
}

// ---- 补单 ----
// 打开弹窗
const openDialog = () => {
  type.value = 'notify'
  dialogFormVisible.value = true
  typeTitle.value = '补单'
}

// 打开详情（补单使用）
const notifyPayOrder = async (row) => {
  // 打开弹窗
  const res = await findPayOrder({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.repayOrder
    openDialog()
  }
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    authCaptcha: '',
    orderId: '',
    pAccount: '',
    money: 0,
    unitPrice: 0,
    uid: 0,
    acId: '',
    channelCode: '',
    platformOid: '',
    payIp: '',
    payRegion: '',
    payDevice: '',
    notifyUrl: '',
    orderStatus: false,
    callbackStatus: false,
    codeUseStatus: false,
    asyncTime: new Date(),
    callTime: new Date(),
  }
}
// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate( async (valid) => {
    if (!valid) return
    switch (type.value) {
      case 'notify':
        console.log(formData.value)
        await console.log("准备补单")
    }
  })
}



// ------ 账户匹配 ------
const isPendingAcc = (row) => {
  console.log(row.adId)
  return row.acId !== "";
}
// ------ 账户匹配 ------

// ------ 账户详情 ------
const formAccData = ref({
  acId: '',
  acRemark: '',
  acAccount: '',
  acPwd: '',
  token: '',
  cid: '',
  countLimit: 0,
  dailyLimit: 0,
  totalLimit: 0,
  status: 0,
  sysStatus: 0,
  uid: 0,
})
const detailAccShow = ref(false)
// 打开详情弹窗
const openAccDetailShow = () => {
  detailAccShow.value = true
}
// 打开详情
const getAccDetails = async(row) => {
  // 打开弹窗
  const res = await findChannelAccount({ acId: row.acId })
  if (res.code === 0) {
    formAccData.value = res.data.revca
    openAccDetailShow()
  }
}
// 关闭详情弹窗
const closeAccDetailShow = () => {
  detailAccShow.value = false
  formAccData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    uid: 0,
  }
}
// ------ 账户详情 ------

</script>

<style>
.el-table__body, .el-table__header {
  border-collapse: collapse;
  border-bottom-width: 2px;
  border-bottom-style: solid;
}

.el-table__body tr, .el-table__header tr {
  border-bottom-width: 2px;
  border-bottom-style: solid;
}

.el-table__body td, .el-table__header th {
  border-right-width: 2px;
  border-right-style: solid;
}

// loading
.el-button .custom-loading .circular {
  margin-right: 6px;
  width: 18px;
  height: 18px;
  animation: loading-rotate 2s linear infinite;
}
.el-button .custom-loading .circular .path {
  animation: loading-dash 1.5s ease-in-out infinite;
  stroke-dasharray: 90, 150;
  stroke-dashoffset: 0;
  stroke-width: 2;
  stroke: var(--el-button-text-color);
  stroke-linecap: round;
}
</style>

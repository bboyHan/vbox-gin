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
      </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
            <el-popover v-model:visible="deleteVisible" :disabled="!multipleSelection.length" placement="top" width="160">
            <p>确定要删除吗？</p>
            <div style="text-align: right; margin-top: 8px;">
                <el-button type="primary" link @click="deleteVisible = false">取消</el-button>
                <el-button type="primary" @click="onDelete">确定</el-button>
            </div>
            <template #reference>
                <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="deleteVisible = true">删除</el-button>
            </template>
            </el-popover>
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
        >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="日期" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="订单ID" prop="orderId" width="120" />
        <el-table-column align="left" label="付方ID" prop="pAccount" width="120" />
        <el-table-column align="left" label="金额" prop="money" width="120" />
        <el-table-column align="left" label="单价积分" prop="unitPrice" width="120" />
        <el-table-column align="left" label="用户ID" prop="uid" width="120" />
        <el-table-column align="left" label="账号ID" prop="acId" width="120" />
        <el-table-column align="left" label="通道编码" prop="channelCode" width="120" />
        <el-table-column align="left" label="平台id" prop="platformOid" width="120" />
        <el-table-column align="left" label="客户ip" prop="payIp" width="120" />
        <el-table-column align="left" label="区域" prop="payRegion" width="120" />
        <el-table-column align="left" label="客户端设备" prop="payDevice" width="120" />
                      <el-table-column label="支付链接" width="200">
                         <template #default="scope">
                            [富文本内容]
                         </template>
                      </el-table-column>
        <el-table-column align="left" label="回调地址" prop="notifyUrl" width="120" />
        <el-table-column align="left" label="订单状态" prop="orderStatus" width="120">
            <template #default="scope">{{ formatBoolean(scope.row.orderStatus) }}</template>
        </el-table-column>
        <el-table-column align="left" label="回调状态" prop="callbackStatus" width="120">
            <template #default="scope">{{ formatBoolean(scope.row.callbackStatus) }}</template>
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
                查看详情
            </el-button>
            <el-button type="primary" link icon="edit" class="table-button" @click="updatePayOrderFunc(scope.row)">变更</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
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
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :title="type==='create'?'添加':'修改'" destroy-on-close>
      <el-scrollbar height="500px">
          <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
            <el-form-item label="订单ID:"  prop="orderId" >
              <el-input v-model="formData.orderId" :clearable="true"  placeholder="请输入订单ID" />
            </el-form-item>
            <el-form-item label="付方ID:"  prop="pAccount" >
              <el-input v-model="formData.pAccount" :clearable="true"  placeholder="请输入付方ID" />
            </el-form-item>
            <el-form-item label="金额:"  prop="money" >
              <el-input v-model.number="formData.money" :clearable="true" placeholder="请输入金额" />
            </el-form-item>
            <el-form-item label="单价积分:"  prop="unitPrice" >
              <el-input v-model.number="formData.unitPrice" :clearable="true" placeholder="请输入单价积分" />
            </el-form-item>
            <el-form-item label="用户ID:"  prop="uid" >
              <el-input v-model.number="formData.uid" :clearable="true" placeholder="请输入用户ID" />
            </el-form-item>
            <el-form-item label="账号ID:"  prop="acId" >
              <el-input v-model="formData.acId" :clearable="true"  placeholder="请输入账号ID" />
            </el-form-item>
            <el-form-item label="通道编码:"  prop="channelCode" >
              <el-input v-model="formData.channelCode" :clearable="true"  placeholder="请输入通道编码" />
            </el-form-item>
            <el-form-item label="平台id:"  prop="platformOid" >
              <el-input v-model="formData.platformOid" :clearable="true"  placeholder="请输入平台id" />
            </el-form-item>
            <el-form-item label="客户ip:"  prop="payIp" >
              <el-input v-model="formData.payIp" :clearable="true"  placeholder="请输入客户ip" />
            </el-form-item>
            <el-form-item label="区域:"  prop="payRegion" >
              <el-input v-model="formData.payRegion" :clearable="true"  placeholder="请输入区域" />
            </el-form-item>
            <el-form-item label="客户端设备:"  prop="payDevice" >
              <el-input v-model="formData.payDevice" :clearable="true"  placeholder="请输入客户端设备" />
            </el-form-item>
            <el-form-item label="支付链接:"  prop="resourceUrl" >
              <RichEdit v-model="formData.resourceUrl"/>
            </el-form-item>
            <el-form-item label="回调地址:"  prop="notifyUrl" >
              <el-input v-model="formData.notifyUrl" :clearable="true"  placeholder="请输入回调地址" />
            </el-form-item>
            <el-form-item label="订单状态:"  prop="orderStatus" >
              <el-switch v-model="formData.orderStatus" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
            </el-form-item>
            <el-form-item label="回调状态:"  prop="callbackStatus" >
              <el-switch v-model="formData.callbackStatus" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
            </el-form-item>
            <el-form-item label="取码状态:"  prop="codeUseStatus" >
              <el-switch v-model="formData.codeUseStatus" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
            </el-form-item>
            <el-form-item label="异步执行时间:"  prop="asyncTime" >
              <el-date-picker v-model="formData.asyncTime" type="date" style="width:100%" placeholder="选择日期" :clearable="true"  />
            </el-form-item>
            <el-form-item label="回调时间:"  prop="callTime" >
              <el-date-picker v-model="formData.callTime" type="date" style="width:100%" placeholder="选择日期" :clearable="true"  />
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

    <el-dialog v-model="detailShow" style="width: 800px" lock-scroll :before-close="closeDetailShow" title="查看详情" destroy-on-close>
      <el-scrollbar height="550px">
        <el-descriptions column="1" border>
                <el-descriptions-item label="订单ID">
                        {{ formData.orderId }}
                </el-descriptions-item>
                <el-descriptions-item label="付方ID">
                        {{ formData.pAccount }}
                </el-descriptions-item>
                <el-descriptions-item label="金额">
                        {{ formData.money }}
                </el-descriptions-item>
                <el-descriptions-item label="单价积分">
                        {{ formData.unitPrice }}
                </el-descriptions-item>
                <el-descriptions-item label="用户ID">
                        {{ formData.uid }}
                </el-descriptions-item>
                <el-descriptions-item label="账号ID">
                        {{ formData.acId }}
                </el-descriptions-item>
                <el-descriptions-item label="通道编码">
                        {{ formData.channelCode }}
                </el-descriptions-item>
                <el-descriptions-item label="平台id">
                        {{ formData.platformOid }}
                </el-descriptions-item>
                <el-descriptions-item label="客户ip">
                        {{ formData.payIp }}
                </el-descriptions-item>
                <el-descriptions-item label="区域">
                        {{ formData.payRegion }}
                </el-descriptions-item>
                <el-descriptions-item label="客户端设备">
                        {{ formData.payDevice }}
                </el-descriptions-item>
                <el-descriptions-item label="支付链接">
                        [富文本内容]
                </el-descriptions-item>
                <el-descriptions-item label="回调地址">
                        {{ formData.notifyUrl }}
                </el-descriptions-item>
                <el-descriptions-item label="订单状态">
                    {{ formatBoolean(formData.orderStatus) }}
                </el-descriptions-item>
                <el-descriptions-item label="回调状态">
                    {{ formatBoolean(formData.callbackStatus) }}
                </el-descriptions-item>
                <el-descriptions-item label="取码状态">
                    {{ formatBoolean(formData.codeUseStatus) }}
                </el-descriptions-item>
                <el-descriptions-item label="异步执行时间">
                      {{ formatDate(formData.asyncTime) }}
                </el-descriptions-item>
                <el-descriptions-item label="回调时间">
                      {{ formatDate(formData.callTime) }}
                </el-descriptions-item>
        </el-descriptions>
      </el-scrollbar>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createPayOrder,
  deletePayOrder,
  deletePayOrderByIds,
  updatePayOrder,
  findPayOrder,
  getPayOrderList
} from '@/api/payOrder'
// 富文本组件
import RichEdit from '@/components/richtext/rich-edit.vue'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({
    name: 'PayOrder'
})

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
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

// 更新行
const updatePayOrderFunc = async(row) => {
    const res = await findPayOrder({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data.repayOrder
        dialogFormVisible.value = true
    }
}


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


// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
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
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              let res
              switch (type.value) {
                case 'create':
                  res = await createPayOrder(formData.value)
                  break
                case 'update':
                  res = await updatePayOrder(formData.value)
                  break
                default:
                  res = await createPayOrder(formData.value)
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

</style>

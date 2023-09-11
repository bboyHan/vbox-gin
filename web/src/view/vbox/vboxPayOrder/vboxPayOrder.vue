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
        <el-form-item label="订单id" prop="order_id">
         <el-input v-model="searchInfo.order_id" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="付方uuid" prop="p_account">
         <el-input v-model="searchInfo.p_account" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="金额" prop="cost">
             <el-input v-model.number="searchInfo.cost" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="uid" prop="uid">
             <el-input v-model.number="searchInfo.uid" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="帐号id" prop="ac_id">
         <el-input v-model="searchInfo.ac_id" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="所属通道" prop="c_channel_id">
         <el-input v-model="searchInfo.c_channel_id" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="平台id" prop="platform_oid">
         <el-input v-model="searchInfo.platform_oid" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="客户ip" prop="pay_ip">
         <el-input v-model="searchInfo.pay_ip" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="区域" prop="pay_region">
         <el-input v-model="searchInfo.pay_region" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="支付链接" prop="resource_url">
         <el-input v-model="searchInfo.resource_url" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="回调地址" prop="notify_url">
         <el-input v-model="searchInfo.notify_url" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="订单状态" prop="order_status">
             <el-input v-model.number="searchInfo.order_status" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="回调状态" prop="callback_status">
             <el-input v-model.number="searchInfo.callback_status" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="取码状态" prop="code_use_status">
             <el-input v-model.number="searchInfo.code_use_status" placeholder="搜索条件" />
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
            <el-popover v-model:visible="deleteVisible" placement="top" width="160">
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
        <!-- <el-table-column align="left" label="日期" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column> -->
        <el-table-column align="left" label="订单id" prop="order_id" width="120"  show-overflow-tooltip="true"/>
        <el-table-column align="left" label="付方uuid" prop="p_account" width="120"  show-overflow-tooltip="true"/>
        <el-table-column align="left" label="金额" prop="cost" width="80" />
        <el-table-column align="left" label="uid" prop="uid" width="50" />
        <el-table-column align="left" label="帐号id" prop="ac_id" width="120"  show-overflow-tooltip="true"/>
        <el-table-column align="left" label="所属通道" prop="c_channel_id" width="120" />
        <el-table-column align="left" label="平台id" prop="platform_oid" width="120"  show-overflow-tooltip="true"/>
        <el-table-column align="left" label="客户ip" prop="pay_ip" width="150" />
        <el-table-column align="left" label="区域" prop="pay_region" width="150"  show-overflow-tooltip="true"/>
        <el-table-column label="支付链接" prop="resource_url" width="200" show-overflow-tooltip="true">
          <template #default="{row}">
            <div class="table-cell limit-height">{{row.resource_url}}</div>
          </template>
        </el-table-column>
                      
        <el-table-column align="left" label="回调地址" prop="notify_url" width="120"  show-overflow-tooltip="true"/>
        <el-table-column align="left" label="订单状态" prop="order_status" width="80" />
        <el-table-column align="left" label="回调状态" prop="callback_status" width="80" />
        <el-table-column align="left" label="取码状态" prop="code_use_status" width="80" />
         <el-table-column align="left" label="创建时间" width="180">
            <template #default="scope">{{ formatDate(scope.row.create_time) }}</template>
         </el-table-column>
         <el-table-column align="left" label="异步执行时间" width="180">
            <template #default="scope">{{ formatDate(scope.row.async_time) }}</template>
         </el-table-column>
         <el-table-column align="left" label="回调时间" width="180">
            <template #default="scope">{{ formatDate(scope.row.call_time) }}</template>
         </el-table-column>
        <!-- <el-table-column align="left" label="操作">
            <template #default="scope">
            <el-button type="primary" link icon="edit" class="table-button" @click="updateVboxPayOrderFunc(scope.row)">变更</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
            </template>
        </el-table-column> -->
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
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="订单id:"  prop="order_id" >
          <el-input v-model="formData.order_id" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="付方uuid:"  prop="p_account" >
          <el-input v-model="formData.p_account" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="金额:"  prop="cost" >
          <el-input v-model.number="formData.cost" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="uid:"  prop="uid" >
          <el-input v-model.number="formData.uid" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="帐号id:"  prop="ac_id" >
          <el-input v-model="formData.ac_id" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="所属通道:"  prop="c_channel_id" >
          <el-input v-model="formData.c_channel_id" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="平台id:"  prop="platform_oid" >
          <el-input v-model="formData.platform_oid" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="客户ip:"  prop="pay_ip" >
          <el-input v-model="formData.pay_ip" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="区域:"  prop="pay_region" >
          <el-input v-model="formData.pay_region" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="支付链接:"  prop="resource_url" >
          <RichEdit v-model="formData.resource_url"/>
        </el-form-item>
        <el-form-item label="回调地址:"  prop="notify_url" >
          <el-input v-model="formData.notify_url" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="订单状态:"  prop="order_status" >
          <el-input v-model.number="formData.order_status" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="回调状态:"  prop="callback_status" >
          <el-input v-model.number="formData.callback_status" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="取码状态:"  prop="code_use_status" >
          <el-input v-model.number="formData.code_use_status" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="创建时间:"  prop="create_time" >
          <el-date-picker v-model="formData.create_time" type="date" style="width:100%" placeholder="选择日期" :clearable="true"  />
        </el-form-item>
        <el-form-item label="异步执行时间:"  prop="async_time" >
          <el-date-picker v-model="formData.async_time" type="date" style="width:100%" placeholder="选择日期" :clearable="true"  />
        </el-form-item>
        <el-form-item label="回调时间:"  prop="call_time" >
          <el-date-picker v-model="formData.call_time" type="date" style="width:100%" placeholder="选择日期" :clearable="true"  />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'VboxPayOrder'
}
</script>

<script setup>
import {
  createVboxPayOrder,
  deleteVboxPayOrder,
  deleteVboxPayOrderByIds,
  updateVboxPayOrder,
  findVboxPayOrder,
  getVboxPayOrderList
} from '@/api/vboxPayOrder'
// 富文本组件
import RichEdit from '@/components/richtext/rich-edit.vue'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        order_id: '',
        p_account: '',
        cost: 0,
        uid: 0,
        ac_id: '',
        c_channel_id: '',
        platform_oid: '',
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
  const table = await getVboxPayOrderList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
            deleteVboxPayOrderFunc(row)
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
      const res = await deleteVboxPayOrderByIds({ ids })
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
const updateVboxPayOrderFunc = async(row) => {
    const res = await findVboxPayOrder({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data.revpo
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteVboxPayOrderFunc = async (row) => {
    const res = await deleteVboxPayOrder({ ID: row.ID })
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
        c_channel_id: '',
        platform_oid: '',
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
                  res = await createVboxPayOrder(formData.value)
                  break
                case 'update':
                  res = await updateVboxPayOrder(formData.value)
                  break
                default:
                  res = await createVboxPayOrder(formData.value)
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

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
      </el-form-item>
        <el-form-item label="通道账户" prop="acAccount">
         <el-input v-model="searchInfo.acAccount" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="账户备注" prop="acRemark">
         <el-input v-model="searchInfo.acRemark" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="账户id" prop="acId">
             <el-input v-model.number="searchInfo.acId" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="通道id" prop="cid">
             <el-input v-model.number="searchInfo.cid" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="状态开关" prop="status">
            <el-input v-model.number="searchInfo.startStatus" placeholder="搜索条件（起）" />
            —
            <el-input v-model.number="searchInfo.endStatus" placeholder="搜索条件（止）" />
        </el-form-item>
        <el-form-item label="系统开关" prop="sysStatus">
            <el-input v-model.number="searchInfo.startSysStatus" placeholder="搜索条件（起）" />
            —
            <el-input v-model.number="searchInfo.endSysStatus" placeholder="搜索条件（止）" />
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
        <el-table-column align="left" label="日期" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="归属用户ID" prop="uid" width="120" />
        <el-table-column align="left" label="通道账户" prop="acAccount" width="120" />
        <el-table-column align="left" label="账户密码" prop="acPwd" width="120" />
        <el-table-column align="left" label="账户备注" prop="acRemark" width="120" />
        <el-table-column align="left" label="账户id" prop="acId" width="120" />
        <el-table-column align="left" label="通道id" prop="cid" width="120" />
        <el-table-column align="left" label="token" prop="cid" width="200" />
        <el-table-column align="left" label="日限额" prop="dailyLimit" width="120" />
        <el-table-column align="left" label="总限额" prop="totalLimit" width="120" />
        <el-table-column align="left" label="状态开关" prop="status" width="120" />
        <el-table-column align="left" label="系统开关" prop="sysStatus" width="120" />
        <el-table-column align="left" label="操作">
          <template #default="scope">
          <el-button type="primary" link icon="edit" class="table-button" @click="updateChannelAccountFunc(scope.row)">变更</el-button>
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
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="账户备注"  prop="acRemark" >
          <el-input v-model="formData.acRemark" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-row>
          <el-col :span="12">
            <el-form-item label="通道账户"  prop="acAccount" >
              <el-input v-model="formData.acAccount" :clearable="true"  placeholder="请输入" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="账户密码"  prop="acPwd" >
              <el-input v-model="formData.acPwd" :clearable="true"  placeholder="请输入" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="通道id"  prop="cid" >
          <el-input v-model.number="formData.cid" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="token"  prop="token" >
          <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-row>
          <el-col :span="12">
            <el-form-item label="日限额"  prop="dailyLimit" >
              <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="总限额"  prop="totalLimit" >
              <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="状态开关"  prop="status" >
          <el-switch v-model="formData.status" active-value="1" inactive-value="0" active-text="开启" inactive-text="关闭"></el-switch>
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
  name: 'ChannelAccount'
}
</script>

<script setup>
import {
  createChannelAccount,
  deleteChannelAccount,
  deleteChannelAccountByIds,
  updateChannelAccount,
  findChannelAccount,
  getChannelAccountList
} from '@/api/vboxChannelAccount'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

function switchValue(status) {
  return status === 1;
}

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        uid: 0,
        acAccount: '',
        acPwd: '',
        acRemark: '',
        token: '',
        acId: 0,
        cid: 0,
        dailyLimit: 0,
        totalLimit: 0,
        status: 1,
        sysStatus: 0,
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
  const table = await getChannelAccountList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
            deleteChannelAccountFunc(row)
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
      const res = await deleteChannelAccountByIds({ ids })
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
const updateChannelAccountFunc = async(row) => {
    const res = await findChannelAccount({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data.revca
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteChannelAccountFunc = async (row) => {
    const res = await deleteChannelAccount({ ID: row.ID })
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
        uid: 0,
        acAccount: '',
        acPwd: '',
        acRemark: '',
        token: '',
        acId: 0,
        cid: 0,
        dailyLimit: 0,
        totalLimit: 0,
        status: 0,
        sysStatus: 0,
        }
}
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              let res
              switch (type.value) {
                case 'create':
                  res = await createChannelAccount(formData.value)
                  break
                case 'update':
                  res = await updateChannelAccount(formData.value)
                  break
                default:
                  res = await createChannelAccount(formData.value)
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

<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule"
               @keyup.enter="onSubmit">
        <el-form-item label="通道账户" prop="acAccount">
          <el-input v-model="searchInfo.acAccount" placeholder="搜索条件"/>
        </el-form-item>
        <el-form-item label="账户备注" prop="acRemark">
          <el-input v-model="searchInfo.acRemark" placeholder="搜索条件"/>
        </el-form-item>
        <el-form-item label="账户id" prop="acId">
          <el-input v-model.number="searchInfo.acId" placeholder="搜索条件"/>
        </el-form-item>
        <el-form-item label="通道id" prop="cid">
          <el-cascader
              v-model="searchInfo.cid"
              :options="channelCodeOptions"
              :props="channelCodeProps"
              @change="handleChange"
              style="width: 100%"
              placeholder="选择通道"
          />
        </el-form-item>
        <el-form-item label="开关状态" prop="status">
          <el-select v-model="searchInfo.status" placeholder="选择状态">
            <el-option label="已开启" value="1"/>
            <el-option label="已关闭" value="0"/>
          </el-select>
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
            <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length"
                       @click="deleteVisible = true">删除
            </el-button>
          </template>
        </el-popover>
        <el-popover v-model:visible="switchOnVisible" placement="top" width="160">
          <p>确定批量开启吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button type="primary" link @click="switchOnVisible = false">取消</el-button>
            <el-button type="primary" @click="onSwitchEnable">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length"
                       @click="switchOnVisible = true">批量开启
            </el-button>
          </template>
        </el-popover>
        <el-popover v-model:visible="switchOffVisible" placement="top" width="160">
          <p>确定批量关闭吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button type="primary" link @click="switchOffVisible = false">取消</el-button>
            <el-button type="primary" @click="onSwitchDisable">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length"
                       @click="switchOffVisible = true">批量关闭
            </el-button>
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
        <el-table-column type="selection" width="55"/>
        <el-table-column align="left" label="归属用户ID" prop="uid" width="120"/>
        <el-table-column align="left" label="通道账户" prop="acAccount" width="120">
          <template #default="scope">
            <el-popover trigger="hover" placement="right-end" width="auto">
              <template #default>
                <div>ID: {{ scope.row.acId }}</div>
                <div>备注: {{ scope.row.acRemark }}</div>
                <div>帐户名: {{ scope.row.acAccount }}</div>
                <div>密钥: {{ scope.row.acPwd }}</div>
                <div>创建时间: {{ formatDate(scope.row.CreatedAt) }}</div>
              </template>
              <template #reference>
                <el-tag>备注: {{ scope.row.acRemark }}</el-tag>
              </template>
            </el-popover>
            <el-tag>账户名: {{ scope.row.acAccount }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="通道id" prop="cid" width="120"/>
        <el-table-column align="left" label="CK" prop="token" width="160">
          <template #default="scope">
            <el-input
                v-model="scope.row.token"
                :rows="2"
                readonly="readonly"
            >
              <template #append>
                <el-button type="primary" link icon="edit" @click="updTokenInfo(scope.row)"></el-button>
              </template>
            </el-input>
          </template>
        </el-table-column>
        <el-table-column align="left" label="限额设置" prop="dailyLimit" width="120">
          <template #default="scope">
            <el-tag>日限额: {{ scope.row.dailyLimit }}</el-tag>
            <el-tag>总限额: {{ scope.row.totalLimit }}</el-tag>
            <el-tag>笔数限制: {{ scope.row.countLimit }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="总限额" prop="totalLimit" width="120"/>
        <el-table-column align="left" label="状态 / 系统开关" prop="status" width="140">
          <template #default="scope">
            <el-row :gutter="12">
              <el-col :span="11">
                <el-switch
                    v-model="scope.row.status"
                    inline-prompt
                    :active-value="1"
                    active-text="开启"
                    :inactive-value="0"
                    inactive-text="关闭"
                    size="large"
                    @change="()=>{switchEnable(scope.row)}"
                />
              </el-col>
              <el-col :span="11">
                <el-switch
                    v-model="scope.row.sysStatus"
                    inline-prompt
                    :active-value="1"
                    active-text="开启"
                    :inactive-value="0"
                    inactive-text="关闭"
                    size="large"
                    disabled
                />
              </el-col>
            </el-row>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">详情</el-button>
            <el-button type="primary" link icon="edit" class="table-button" @click="updateChannelAccountFunc(scope.row)">变更</el-button>
            <el-button type="warning" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
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

    <!--  创建  -->
    <el-dialog
      v-model="dialogFormVisible"
      :before-close="closeDialog"
      :title="typeTitle"
      destroy-on-close
    >
      <el-scrollbar height="500px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
          <el-form-item label="账户备注" prop="acRemark">
            <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
          </el-form-item>
          <el-row>
            <el-col :span="12">
              <el-form-item label="通道账户" prop="acAccount">
                <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="账户密码" prop="acPwd">
                <el-input v-model="formData.acPwd" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="通道id" prop="cid">
                <el-cascader
                    v-model="formData.cid"
                    :options="channelCodeOptions"
                    :props="channelCodeProps"
                    @change="handleChange"
                    style="width: 100%"
                />
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="token" prop="token">
            <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入"/>
          </el-form-item>
          <el-row>
            <el-col :span="8">
              <el-form-item label="日限额" prop="dailyLimit">
                <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限额" prop="totalLimit">
                <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="笔数限额" prop="countLimit">
                <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="状态开关" prop="status">
            <el-switch v-model="formData.status" active-value="1" inactive-value="0" active-text="开启"
                       inactive-text="关闭"></el-switch>
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

    <!--  修改  -->
    <el-dialog
      v-model="dialogUpdFormVisible"
      :before-close="closeDialog"
      :title="typeTitle"
      destroy-on-close
    >
      <el-scrollbar height="500px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
          <el-form-item label="账户备注" prop="acRemark">
            <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
          </el-form-item>
          <el-row>
            <el-col :span="12">
              <el-form-item label="通道账户" prop="acAccount">
                <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="账户密码" prop="acPwd">
                <el-input v-model="formData.acPwd" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="通道id" prop="cid">
                <el-input v-model="formData.cid" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="token" prop="token">
            <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入"/>
          </el-form-item>
          <el-row>
            <el-col :span="8">
              <el-form-item label="日限额" prop="dailyLimit">
                <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限额" prop="totalLimit">
                <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="笔数限额" prop="countLimit">
                <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="状态开关" prop="status">
            <el-switch v-model="formData.status" active-value="1" inactive-value="0" active-text="开启"
                       inactive-text="关闭"></el-switch>
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

    <!--  CK  -->
    <el-dialog v-model="dialogTokenFormVisible" :before-close="closeDialog" :title="'变更CK'"
               destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="token" prop="token">
          <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 查看详情 -->
    <el-dialog v-model="detailShow" style="width: 800px" lock-scroll :before-close="closeDetailShow" title="查看详情" destroy-on-close>
      <el-scrollbar height="550px">
        <el-descriptions column="6" border>
          <el-descriptions-item label="用户id" :span="6">{{ formData.uid }}</el-descriptions-item>
          <el-descriptions-item label="账户ID" :span="6">{{ formData.acId }}</el-descriptions-item>
          <el-descriptions-item label="账户备注" :span="6">{{ formData.acRemark }}</el-descriptions-item>
          <el-descriptions-item label="通道账户" :span="3">{{ formData.acAccount }}</el-descriptions-item>
          <el-descriptions-item label="账户密码" :span="3">{{ formData.acPwd }}</el-descriptions-item>
          <el-descriptions-item label="ck" :span="6">{{ formData.token }}</el-descriptions-item>
          <el-descriptions-item label="通道id" :span="6">{{ formData.cid }}</el-descriptions-item>
          <el-descriptions-item label="笔数限制" :span="2">{{ formData.countLimit }}</el-descriptions-item>
          <el-descriptions-item label="日限额" :span="2">{{ formData.dailyLimit }}</el-descriptions-item>
          <el-descriptions-item label="总限额" :span="2">{{ formData.totalLimit }}</el-descriptions-item>
          <el-descriptions-item label="状态开关" :span="3">{{ formData.status===0?'关闭':'开启' }}</el-descriptions-item>
          <el-descriptions-item label="系统开关" :span="3">{{ formData.sysStatus===0?'关闭':'开启' }}</el-descriptions-item>
        </el-descriptions>
      </el-scrollbar>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createChannelAccount,
  deleteChannelAccount,
  deleteChannelAccountByIds,
  updateChannelAccount,
  findChannelAccount,
  getChannelAccountList,
  switchEnableCA,
  switchEnableCAByIds,
} from '@/api/channelAccount'
import {
  getChannelProductSelf
} from '@/api/channelProduct'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, nextTick } from 'vue'

defineOptions({
  name: 'ChannelAccount'
})

//通道产品

const channelCodeOptions = ref([])
const vcpTableData = ref([])

const channelCodeProps = {
  expandTrigger: 'hover',
  checkStrictly: false,
  emitPath: false,
}

const handleChange = (value) => {
  console.log(value)
}

const setChannelCodeOptions = (ChannelCodeData, optionsData, disabled) => {
  ChannelCodeData &&
  ChannelCodeData.forEach(item => {
    if (item.children && item.children.length) {
      const option = {
        value: item.channelCode + '',
        label: item.productName,
        children: []
      }
      setChannelCodeOptions(
          item.children,
          option.children,
      )
      optionsData.push(option)
    } else {
      const option = {
        value: item.channelCode + '',
        label: item.productName,
      }
      optionsData.push(option)
    }
  })
}

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
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

// 验证规则
const rule = reactive({
  acAccount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  {
    whitespace: true,
    message: '不能只输入空格',
    trigger: ['input', 'blur'],
  }
  ],
  cid: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
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
  const vcpTable = await getChannelProductSelf({page: 1, pageSize: 999, ...searchInfo.value})

  if (table.code === 0) {
    tableData.value = table.data.list
    vcpTableData.value = vcpTable.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
    setOptions()
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async() => {
  channelCodeOptions.value = []
  setChannelCodeOptions(vcpTableData.value, channelCodeOptions.value, false)
}

// 获取需要的字典 可能为空 按需保留

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
const typeTitle = ref('')

// 更新行
// ca 更新
const dialogUpdFormVisible = ref(false)
const updateChannelAccountFunc = async(row) => {
  const res = await findChannelAccount({ ID: row.ID })
  type.value = 'update'
  typeTitle.value = '修改'
  if (res.code === 0) {
    formData.value = res.data.revca
    dialogUpdFormVisible.value = true
  }
}

// 删除行
const deleteChannelAccountFunc = async(row) => {
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

// 查看详情控制标记
const detailShow = ref(false)

// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}

// 打开详情
const getDetails = async(row) => {
  // 打开弹窗
  const res = await findChannelAccount({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.revca
    openDetailShow()
  }
}

// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
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

// 打开弹窗
const openDialog = () => {
  type.value = 'create'
  typeTitle.value = '创建'
  dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  dialogUpdFormVisible.value = false
  dialogTokenFormVisible.value = false
  formData.value = {
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
// 弹窗确定
const enterDialog = async() => {
  elFormRef.value?.validate(async(valid) => {
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

// 通道账号开关（批量）
const switchOnVisible = ref(false)
const switchOffVisible = ref(false)
// 通道账号开关
const caInfo = ref({
  status: 1,
  id: '',
})
// 批量ca data
const switchData = ref({
  ids: [],
  status: 0,
})

const switchEnable = async (row) => {
  console.log(row)
  caInfo.value = JSON.parse(JSON.stringify(row))
  await nextTick()
  const req = {
    ...caInfo.value
  }
  const res = await switchEnableCA(req)
  if (res.code === 0) {
    ElMessage({type: 'success', message: `${req.status === 0 ? '禁用' : '启用'}成功`})
    await getTableData()
  }
}

// 批量开启
const onSwitchEnable = async() => {
  const ids = []
  if (multipleSelection.value.length === 0) {
    ElMessage({
      type: 'warning',
      message: '请选择要开启的数据'
    })
    return
  }
  multipleSelection.value &&
  multipleSelection.value.map(item => {
    ids.push(item.ID)
  })
  switchData.value.ids = ids
  switchData.value.status = 1
  const req = {
    ...switchData.value
  }
  const res = await switchEnableCAByIds(req)
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '开启成功'
    })
    if (tableData.value.length === ids.length && page.value > 1) {
      page.value--
    }
    switchOnVisible.value = false
    getTableData()
  }
}

//批量关闭
const onSwitchDisable = async() => {
  const ids = []
  if (multipleSelection.value.length === 0) {
    ElMessage({
      type: 'warning',
      message: '请选择要关闭的数据'
    })
    return
  }
  multipleSelection.value &&
  multipleSelection.value.map(item => {
    ids.push(item.ID)
  })
  switchData.value.ids = ids
  switchData.value.status = 0
  const req = {
    ...switchData.value
  }
  const res = await switchEnableCAByIds(req)
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '关闭成功'
    })
    if (tableData.value.length === ids.length && page.value > 1) {
      page.value--
    }
    switchOffVisible.value = false
    getTableData()
  }
}

// 通道账号token更新
const dialogTokenFormVisible = ref(false)

const caTokenInfo = ref({
  token: '',
  id: '',
})
const updTokenInfo = async (row) => {
  console.log(row)
  caTokenInfo.value = JSON.parse(JSON.stringify(row))
  await nextTick()
  const req = {
    ...caTokenInfo.value
  }
  const res = await updateTokenInfoFunc(req)
  if (res.code === 0) {
    ElMessage({type: 'success', message: `更新CK${req.status === 0 ? '成功' : '失败'}`})
    await getTableData()
  }
}

const updateTokenInfoFunc = async (row) => {
  const res = await findChannelAccount({ID: row.ID})
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.revca
    dialogTokenFormVisible.value = true
  }
}

</script>

<style>

</style>

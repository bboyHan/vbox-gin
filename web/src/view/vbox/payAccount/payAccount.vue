<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
        <el-form-item label="付方备注" prop="pRemark">
         <el-input v-model="searchInfo.pRemark" placeholder="搜索条件" />
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
        <el-table-column align="left" label="商户ID" prop="pAccount" width="120" />
        <el-table-column align="left" label="商户Key" prop="pKey" width="360" />
        <el-table-column align="left" label="商户备注" prop="pRemark" width="120" />
        <el-table-column align="left" label="状态开关" prop="status" width="120">
          <template #default="scope">
            <el-switch
                v-model="scope.row.status"
                inline-prompt
                :active-value="1"
                :inactive-value="0"
                @change="()=>{switchEnable(scope.row)}"
            />
          </template>
        </el-table-column>
        <el-table-column align="left" label="创建时间" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作">
          <template #default="scope">
            <el-button type="primary" link icon="edit" class="table-button" @click="showVPAInfo(scope.row)">对接</el-button>
            <el-button type="primary" link icon="edit" class="table-button" @click="updatePayAccountFunc(scope.row)">变更</el-button>
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
    <el-dialog width="360px" v-model="dialogFormVisible" :before-close="closeDialog" :title="type==='create'?'添加':'修改'" destroy-on-close>
      <el-scrollbar height="50px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
          <el-form-item label="商户备注"  prop="pRemark" >
            <el-input v-model="formData.pRemark" :clearable="true"  placeholder="请输入" />
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
          <el-descriptions-item label="UID">
                  {{ formData.uid }}
          </el-descriptions-item>
          <el-descriptions-item label="付方账户">
                  {{ formData.pAccount }}
          </el-descriptions-item>
          <el-descriptions-item label="付方Key">
                  {{ formData.pKey }}
          </el-descriptions-item>
          <el-descriptions-item label="付方备注">
                  {{ formData.pRemark }}
          </el-descriptions-item>
          <el-descriptions-item label="状态开关">
                  {{ formData.status }}
          </el-descriptions-item>
        </el-descriptions>
      </el-scrollbar>
    </el-dialog>

    <!--  对接信息复制  -->
    <el-dialog v-model="dialogInfoVisible" :before-close="closeInfoDialog" :title="infoTitle" destroy-on-close>
      <el-form :model="formData" label-position="left" ref="elFormRef" label-width="80px">
        <el-form-item label="商户ID"  prop="pAccount" >
          <el-input v-model="formData.pAccount" readonly />
        </el-form-item>
        <el-form-item label="商户Key"  prop="pKey" >
          <el-input v-model="formData.pKey" readonly />
        </el-form-item>
        <el-form-item label="通道编码"  prop="cid" >
          <el-cascader v-model="formData.cid" :options="channelCodeOptions" :props="channelCodeProps" @change="" style="width: 100%"/>
        </el-form-item>
        <el-form-item label="商户备注"  prop="pRemark" >
          <el-input v-model="formData.pRemark" readonly />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeInfoDialog">取 消</el-button>
          <el-button class="btn-copy" type="primary" @click="previewInfo">预 览</el-button>
          <el-button class="btn-copy" type="primary" @click="copyVPAInfo">复 制</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  预览  -->
    <el-dialog v-model="previewFlag">
      <PreviewCodeDialog v-if="previewFlag" ref="previewNode" :preview-code="preViewCode"/>
      <template #footer>
        <div class="dialog-footer" style="padding-top:14px;padding-right:14px">
          <el-button type="primary" @click="copy">复 制</el-button>
          <el-button type="primary" @click="previewFlag = false">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createPayAccount,
  deletePayAccount,
  deletePayAccountByIds,
  updatePayAccount,
  findPayAccount,
  getPayAccountList,
  switchEnablePA,
} from '@/api/payAccount'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, nextTick } from 'vue'
import ClipboardJS from "clipboard";
import {getChannelProductSelf} from "@/api/channelProduct";
import PreviewCodeDialog from "@/view/systemTools/autoCode/component/previewCodeDialg.vue";

defineOptions({
    name: 'PayAccount'
})

const previewFlag = ref(false)
const preViewCode = ref({})
const previewNode = ref(null)
const copy = () => {
  previewNode.value.copy()
}

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        uid: 0,
        pAccount: '',
        pKey: '',
        cid: '',
        pRemark: '',
        status: 0,
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
  const table = await getPayAccountList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  const vcpTable = await getChannelProductSelf({page: 1, pageSize: 999, type: formData.value.type})
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
  if (vcpTable.code === 0) {
    vcpTableData.value = vcpTable.data.list
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
            deletePayAccountFunc(row)
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
      const res = await deletePayAccountByIds({ ids })
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
const updatePayAccountFunc = async(row) => {
    const res = await findPayAccount({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data.repacc
        dialogFormVisible.value = true
    }
}


// 删除行
const deletePayAccountFunc = async (row) => {
    const res = await deletePayAccount({ ID: row.ID })
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
  const res = await findPayAccount({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.repacc
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
          uid: 0,
          pAccount: '',
          pKey: '',
          pRemark: '',
          status: 0,
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
        uid: 0,
        pAccount: '',
        pKey: '',
        pRemark: '',
        status: 0,
        }
}
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              let res
              switch (type.value) {
                case 'create':
                  res = await createPayAccount(formData.value)
                  break
                case 'update':
                  res = await updatePayAccount(formData.value)
                  break
                default:
                  res = await createPayAccount(formData.value)
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

// 对接信息
const dialogInfoVisible = ref(false)
const infoType = ref('')
const infoTitle = ref('')

// 关闭弹窗
const closeInfoDialog = () => {
  dialogInfoVisible.value = false
  formData.value = {
    uid: 0,
    pAccount: '',
    pKey: '',
    pRemark: '',
    status: 0,
  }
}

// 查看对接信息
const showVPAInfo = async(row) => {
  const res = await findPayAccount({ ID: row.ID })
  infoType.value = 'show'
  infoTitle.value = '对接信息'
  if (res.code === 0) {
    formData.value = res.data.repacc
    dialogInfoVisible.value = true
  }
}


// 复制对接信息
const previewInfo = () => {
  let res = formData.value;
  if (!res.cid){
    ElMessage({
      showClose: true,
      message: "未指定通道产品编码,请核实",
      type: 'error'
    })
    return
  }
  console.log(res)
  let copyInfo = `
    商户备注: ${res.pRemark}
    商户ID: ${res.pAccount}
    商户Key: ${res.pKey}
    通道编码: ${res.cid}
    服务网关:
    `
  preViewCode.value = {
    '对接信息': "```shell" + copyInfo + "```"
  }
  previewFlag.value = true
}
// 复制对接信息
const copyVPAInfo = () => {
  let res = formData.value;
  console.log(res)
  let copyInfo = `
    商户备注: ${res.pRemark}
    商户ID: ${res.pAccount}
    商户Key: ${res.pKey}
    通道编码: ${res.cid}
  `
  const clipboard = new ClipboardJS('.btn-copy', {
    text: () => copyInfo
  });

  clipboard.on('success', () => {
    ElMessage({
      type: 'success',
      message: '复制成功'
    })
    clipboard.destroy(); // 销毁 ClipboardJS 实例
  });

  clipboard.on('error', () => {
    ElMessage({
      type: 'error',
      message: '复制异常'
    })
    clipboard.destroy(); // 销毁 ClipboardJS 实例
  });
  closeInfoDialog()
};

// 付方的启用与禁用
const paInfo = ref({
  status: 1,
  id: '',
})

const switchEnable = async(row) => {
  console.log(row)
  paInfo.value = JSON.parse(JSON.stringify(row))
  await nextTick()
  const req = {
    ...paInfo.value
  }
  const res = await switchEnablePA(req)
  if (res.code === 0) {
    ElMessage({ type: 'success', message: `${req.status === 0 ? '禁用' : '启用'}成功` })
    await getTableData()
  }
}
</script>

<style>

</style>

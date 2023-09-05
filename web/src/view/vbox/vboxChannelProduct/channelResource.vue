<template>
  <div>
    <div class="gva-table-box">
      <div>
        <el-row :gutter="12">
          <el-col v-for="(item, index) in tableData" :key="index" :span="6">
            <el-card shadow="hover">
              <template #header>
                <el-descriptions :title="item.productName" :column="2" border>
                  <template #extra>
                    <el-button link>通道编码：{{ item.channelCode }}</el-button>
                  </template>
                  <el-descriptions-item>
                    <template #label><div>形式</div></template>
                    {{ typeMap[item.type] }}
                  </el-descriptions-item>
                  <el-descriptions-item>
                    <template #label><div>支付</div></template>
                    {{ payTypeMap[item.payType] }}
                  </el-descriptions-item>
                </el-descriptions>
              </template>
              <el-row :gutter="12">
                <el-col :span="12">
                  <el-button>批量添加</el-button>
                </el-col>
                <el-col :span="12">
                  <el-button>添加账号</el-button>
                </el-col>
              </el-row>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </div>
    <el-dialog v-model="dialogFormVisible" :title="dialogTitle">
      <el-form ref="channelProductForm" :model="form" :rules="rules" label-width="80px">
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
  name: 'VboxChannelResource'
}
</script>

<script setup>
import {
  createVboxChannelProduct,
  deleteVboxChannelProduct,
  updateVboxChannelProduct,
  getVboxChannelProductList
} from '@/api/vboxChannelProduct'
// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { typeMap, payTypeMap } from "@/utils/channel";

const mustUint = (rule, value, callback) => {
  if (!/^[0-9]*[1-9][0-9]*$/.test(value)) {
    return callback(new Error('请输入正整数'))
  }
  return callback()
}

const ChannelCodeOption = ref([
  {
    channelCode: 0,
    productName: '根产品名'
  }
])
const dialogType = ref('add')
const dialogTitle = ref('新增通道产品')
const dialogFormVisible = ref(false)

// 自动化生成的字典（可能为空）以及字段
const form = ref({
        parentId: 0,
        channelCode: 0,
        productName: '',
        productId: '',
        ext: '',
        type: 0,
        payType: '',
        })

// 验证规则
const rules = ref({
  channelCode: [
    { required: true, message: '请输入通道编码', trigger: 'blur' },
    { validator: mustUint, trigger: 'blur', message: '必须为正整数' }
  ],
  productName: [
    { required: true, message: '请输入产品名', trigger: 'blur' }
  ],
  parentId: [
    { required: true, message: '请选择父编码', trigger: 'blur' },
  ]
})


// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(999)
const tableData = ref([])
const searchInfo = ref({})

// 查询
const getTableData = async() => {
  const table = await getVboxChannelProductList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    console.log('==>bb' + JSON.stringify(tableData.value))
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()


// ============== 表格控制部分结束 ===============

// 删除通道
const deleteChannel = (row) => {
  ElMessageBox.confirm('此操作将永久删除该产品, 是否继续?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
      .then(async() => {
        const res = await deleteVboxChannelProduct({ ID: row.ID , channelCode: row.channelCode})
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: '删除成功!'
          })
          if (tableData.value.length === 1 && page.value > 1) {
            page.value--
          }
          getTableData()
        }
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: '已取消删除'
        })
      })
}
// 初始化表单
const channelProductForm = ref(null)
const initForm = () => {
  if (channelProductForm.value) {
    channelProductForm.value.resetFields()
  }
  form.value = {
    channelCode: 0,
    productName: '',
    parentId: 0
  }
}
// 关闭窗口
const closeDialog = () => {
  initForm()
  dialogFormVisible.value = false
}
// 确定弹窗
const enterDialog = () => {
  channelProductForm.value.validate(async valid => {
    if (valid) {
      form.value.channelCode = Number(form.value.channelCode)
      switch (dialogType.value) {
        case 'add':
        {
          const res = await createVboxChannelProduct(form.value)
          if (res.code === 0) {
            ElMessage({
              type: 'success',
              message: '添加成功!'
            })
            getTableData()
            closeDialog()
          }
        }
          break
        case 'edit':
        {
          const res = await updateVboxChannelProduct(form.value)
          if (res.code === 0) {
            ElMessage({
              type: 'success',
              message: '添加成功!'
            })
            getTableData()
            closeDialog()
          }
        }
          break
      }

      initForm()
      dialogFormVisible.value = false
    }
  })
}

const setOptions = () => {
  ChannelCodeOption.value = [
    {
      channelCode: 0,
      productName: '根通道产品'
    }
  ]
  setChannelCodeOptions(tableData.value, ChannelCodeOption.value, false)
}

const setChannelCodeOptions = (ChannelCodeData, optionsData, disabled) => {
  form.value.channelCode = String(form.value.channelCode)
  ChannelCodeData &&
  ChannelCodeData.forEach(item => {
    if (item.children && item.children.length) {
      const option = {
        channelCode: item.channelCode,
        productName: item.productName,
        disabled: disabled || item.channelCode === form.value.channelCode,
        children: []
      }
      setChannelCodeOptions(
          item.children,
          option.children,
          disabled || item.channelCode === form.value.channelCode
      )
      optionsData.push(option)
    } else {
      const option = {
        channelCode: item.channelCode,
        productName: item.productName,
        disabled: disabled || item.channelCode === form.value.channelCode
      }
      optionsData.push(option)
    }
  })
}

// 增加通道产品
const addChannelCode = (parentId) => {
  initForm()
  dialogTitle.value = '新增子产品项'
  dialogType.value = 'add'
  form.value.parentId = parentId
  setOptions()
  dialogFormVisible.value = true
}

// 编辑产品
const editChannelCode = (row) => {
  setOptions()
  dialogTitle.value = '编辑产品'
  dialogType.value = 'edit'
  for (const key in form.value) {
    form.value[key] = row[key]
  }
  setOptions()
  dialogFormVisible.value = true
}

</script>

<style lang="scss">
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>

<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="addChannelCode(0)">新增通道产品</el-button>
      </div>
      <el-table
        :tree-props="{children: 'children', hasChildren: 'hasChildren'}"
        :data="tableData"
        row-key="channelCode"
        style="width: 100%"
        default-expand-all
      >
        <el-table-column align="left" label="通道编码" prop="channelCode" width="120" />
        <el-table-column align="left" label="产品名称" prop="productName" width="160" />
        <el-table-column align="left" label="产品ID" prop="productId" width="120" />
        <el-table-column align="left" label="附加参数" prop="ext" width="120" />
        <el-table-column align="left" label="产码方式" prop="type" width="120">
          <template #default="scope"><el-tag>{{ typeMap[scope.row.type] }}</el-tag></template>
        </el-table-column>
        <el-table-column align="left" label="支付方式" prop="payType" width="120">
          <template #default="scope"><el-tag>{{ payTypeMap[scope.row.payType] }}</el-tag></template>
        </el-table-column>
        <el-table-column align="left" label="创建时间" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作">
          <template #default="scope">
            <el-button icon="plus" type="primary" link @click="addChannelCode(scope.row.channelCode)">新增子产品项</el-button>
            <el-button type="primary" link icon="edit" class="table-button" @click="editChannelCode(scope.row)">编辑</el-button>
            <el-button type="primary" link icon="delete" @click="deleteChannel(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <el-dialog v-model="dialogFormVisible" :title="dialogTitle">
      <el-scrollbar height="500px">
        <el-form ref="channelProductForm" :model="form" :rules="rules" label-width="80px">
          <el-form-item label="父级编码" prop="parentId">
            <el-cascader
              v-model="form.parentId"
              style="width:100%"
              :disabled="dialogType==='add'"
              :options="ChannelCodeOption"
              :props="{ checkStrictly: true,label:'productName',value:'channelCode',disabled:'disabled',emitPath:false}"
              :show-all-levels="false"
              filterable
            />
          </el-form-item>
          <el-form-item label="通道编码"  prop="channelCode" >
            <el-input v-model="form.channelCode" :disabled="dialogType==='edit'" autocomplete="off" maxlength="15" />
          </el-form-item>
          <el-form-item label="产品名称"  prop="productName" >
            <el-input v-model="form.productName" :clearable="true"  placeholder="请输入" />
          </el-form-item>
          <el-form-item label="产品ID"  prop="productId" >
            <el-input v-model="form.productId" :clearable="true"  placeholder="请输入" />
          </el-form-item>
          <el-form-item label="附加参数"  prop="ext" >
            <el-input v-model="form.ext" :clearable="true"  placeholder="请输入" />
          </el-form-item>
          <el-form-item label="产码方式"  prop="type" >
            <el-radio-group v-model="form.type">
              <el-radio label="1"><template #default><span>引导</span></template></el-radio>
              <el-radio label="2"><template #default><span>预产</span></template></el-radio>
              <el-radio label="3"><template #default><span>原生</span></template></el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="支付方式"  prop="payType" >
            <el-radio-group v-model="form.payType">
              <el-radio label="wechat"><template #default><span>微信</span></template></el-radio>
              <el-radio label="alipay"><template #default><span>支付宝</span></template></el-radio>
              <el-radio label="app"><template #default><span>三方App</span></template></el-radio>
            </el-radio-group>
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
  </div>
</template>

<script setup>
import {
  createChannelProduct,
  deleteChannelProduct,
  deleteChannelProductByIds,
  updateChannelProduct,
  findChannelProduct,
  getChannelProductList
} from '@/api/channelProduct'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { typeMap, payTypeMap, mustUint } from "@/utils/channel";
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({
    name: 'ChannelProduct'
})

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
  parentId: '0',
  channelCode: '0',
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
  const table = await getChannelProductList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
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
        const res = await deleteChannelProduct({ ID: row.ID , channelCode: row.channelCode})
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
    channelCode: '0',
    productName: '',
    parentId: '0'
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
      form.value.channelCode = String(form.value.channelCode)
      form.value.parentId = String(form.value.parentId)
      form.value.type = Number(form.value.type)
      switch (dialogType.value) {
        case 'add':
        {
          const res = await createChannelProduct(form.value)
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
          const res = await updateChannelProduct(form.value)
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

<style>

</style>

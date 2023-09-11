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
                  <el-button type="primary" icon="plus" @click="openDialog">添加账号</el-button>
<!--                  <el-button>添加账号</el-button>-->
                </el-col>
              </el-row>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </div>
    <!--  账号添加/修改  -->
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :title="type==='create'?'添加':'修改'"
               destroy-on-close>
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
  createChannelAccount,
  updateChannelAccount,
  queryCAHisRecords,
} from '@/api/vboxChannelAccount'
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
  countLimit: 0,
  status: 1,
  sysStatus: 0,
})

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
// 打开弹窗
const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  dialogTokenFormVisible.value = false
  dialogQueryFormVisible.value = false
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
    countLimit: 0,
    status: 0,
    sysStatus: 0,
  }
}
// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    formData.value.cid = Number(formData.value.cid)
    formData.value.status = Number(formData.value.status)
    let res
    switch (type.value) {
      case 'create':
        res = await createChannelAccount(formData.value)
        break
      case 'update':
        res = await updateChannelAccount(formData.value)
        break
      case 'query':
        res = await queryCAHisRecords(formData.value)
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

</script>

<style lang="scss">
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>

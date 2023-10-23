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
      <el-form :model="accFormData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="账户备注" prop="acRemark">
          <el-input v-model="accFormData.acRemark" :clearable="true" placeholder="请输入"/>
        </el-form-item>
        <el-row>
          <el-col :span="12">
            <el-form-item label="通道账户" prop="acAccount">
              <el-input v-model="accFormData.acAccount" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="账户密码" prop="acPwd">
              <el-input v-model="accFormData.acPwd" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="24">
            <el-form-item label="通道id" prop="cid">
              <el-cascader
                  v-model="accFormData.cid"
                  :options="channelCodeOptions"
                  :props="channelCodeProps"
                  @change="handleChange"
                  style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="token" prop="token">
          <el-input v-model="accFormData.token" type="textarea" :clearable="true" placeholder="请输入"/>
        </el-form-item>
        <el-row>
          <el-col :span="8">
            <el-form-item label="日限额" prop="dailyLimit">
              <el-input v-model.number="accFormData.dailyLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限额" prop="totalLimit">
              <el-input v-model.number="accFormData.totalLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="笔数限额" prop="countLimit">
              <el-input v-model.number="accFormData.countLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="状态开关" prop="status">
          <el-switch v-model="accFormData.status" active-value="1" inactive-value="0" active-text="开启"
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

<script setup>
import {
  findChannelProduct,
  getChannelProductSelf
} from '@/api/channelProduct'
import {
  createChannelAccount
} from '@/api/channelAccount';

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { typeMap, payTypeMap } from "@/utils/channel";

defineOptions({
  name: 'OrgProduct'
})

const dialogFormVisible = ref(false)
const accFormData = ref({
  uid: 0,
  acAccount: '',
  acPwd: '',
  acRemark: '',
  token: '',
  acId: '',
  cid: 0,
  dailyLimit: 0,
  totalLimit: 0,
  countLimit: 0,
  status: 1,
  sysStatus: 0,
})

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  channelCode: '',
  ext: '',
  parentId: 0,
  payType: '',
  productId: '',
  productName: '',
  type: false,
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
    if (searchInfo.value.type === ""){
      searchInfo.value.type=null
    }
    getTableData()
  })
}

// 查询
const getTableData = async() => {
  const table = await getChannelProductSelf({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize

    vcpTableData.value = table.data.list
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

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 打开弹窗
const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    channelCode: '',
    ext: '',
    parentId: 0,
    payType: '',
    productId: '',
    productName: '',
    type: false,
  }
}

// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    accFormData.value.status = Number(accFormData.value.status)
    let res
    switch (type.value) {
      case 'create':
        res = await createChannelAccount(accFormData.value)
        break
      default:
        res = await createChannelAccount(accFormData.value)
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
</script>

<style>

</style>


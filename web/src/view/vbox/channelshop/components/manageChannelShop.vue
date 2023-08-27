<!-- <template>
    <div>
        <el-dialog v-model="dialogManageFormVisible" :before-close="closeManageDialog" :title="店铺详情" destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <!-- <el-form-item label="用户ID:"  prop="uid" >
          <el-input v-model.number="formData.uid" :clearable="true" placeholder="请输入" />
        </el-form-item> -->
        <!-- <el-form-item label="通道ID:"  prop="cid" >
          <el-input v-model="formData.cid" :clearable="true"  placeholder="请输入" />
        </el-form-item> -->
        <el-form-item label="通道:"  prop="channel" >
          <el-input v-model="formData.channel" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="店铺备注:"  prop="shop_remark" >
          <el-input v-model="formData.shop_remark" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="店地址:"  prop="address" >
          <el-input v-model="formData.address" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="金额:"  prop="money" >
          <el-input v-model.number="formData.money" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <!-- <el-form-item label="开关:"  prop="status" >
          <el-input v-model.number="formData.status" :clearable="true" placeholder="请输入" />
        </el-form-item> -->
        <el-table-column align="left" label="操作">
            <template #default="scope">
             <!-- <el-button type="primary" link icon="edit" class="table-button" @click="updateChannelShopFunc(scope.row)">管理</el-button> -->
            <el-button type="primary" link icon="edit" class="table-button" @click="updateChannelShopFunc(scope.row)">更新</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
            </template>
        </el-table-column>
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
  name: 'manageChannelShop'
}
</script>

<script setup>
import {
  createChannelShop,
  deleteChannelShop,
  deleteChannelShopByIds,
  updateChannelShop,
  findChannelShop,
  getChannelShopList
} from '@/api/channelshop'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        uid: 0,
        cid: '',
        channel: '',
        shop_remark: '',
        address: '',
        money: 0,
        status: 0,
        })

// 验证规则
const rule = reactive({
})

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
  const table = await getChannelShopList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    console.log('==>channel1 ' + JSON.stringify(tableData.value))
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
  console.log('==>channel2 ' + JSON.stringify(tableData.value))
}

getTableData()

function handleSwitchChange(row, value) {
  row.status = value ? 1 : 0;
}

function switchValue(status) {
  return status === 1 ? true : false;
}

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
            deleteChannelShopFunc(row)
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
      const res = await deleteChannelShopByIds({ ids })
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
const updateChannelShopFunc = async(row) => {
    const res = await findChannelShop({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data.rechShop
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteChannelShopFunc = async (row) => {
    const res = await deleteChannelShop({ ID: row.ID })
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
const dialogManageFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    dialogManageFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogManageFormVisible.value = false
    formData.value = {
        uid: 0,
        cid: '',
        channel: '',
        shop_remark: '',
        address: '',
        money: 0,
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
                  res = await createChannelShop(formData.value)
                  break
                case 'update':
                  res = await updateChannelShop(formData.value)
                  break
                default:
                  res = await createChannelShop(formData.value)
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
</script> -->
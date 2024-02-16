<template>
  <div>
    <div class="gva-search-box">
      <el-row :gutter="12">
        <el-col :xs="24" :span="6">
            <CenterCard title="我的积分" :custom-style="walletCustomStyle">
              <template #action>
                <span class="gvaIcon-prompt" style="color: #999" />
              </template>
              <template #body>
                <!--              <Order :channel-code="searchInfo.cid"/>-->
                <div class="acc-container">
                  <div class="indicator">
                  <span>
                    <div class="label"></div>
                    <div class="value">{{ userBalance }}</div>
                  </span>
                  </div>
                </div>
              </template>
            </CenterCard>
        </el-col>
      </el-row>
    </div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
        <el-form-item label="ID" prop="eventId">
         <el-input v-model="searchInfo.eventId" placeholder="搜索ID" />
        </el-form-item>
        <el-form-item label="用户名" prop="username">
         <el-input v-model="searchInfo.username" placeholder="搜索用户名" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="searchInfo.type" placeholder="选择类型">
            <el-option label="充值" value="1"/>
            <el-option label="划转" value="2"/>
            <el-option label="消费" value="3"/>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button icon="refresh" @click="onReset"></el-button>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
      >
        <el-table-column align="left" label="ID" prop="eventId" width="220" />
<!--        <el-table-column align="left" label="用户名" prop="username" width="120" />-->
        <el-table-column align="left" label="积分" prop="recharge" width="120" />
        <el-table-column align="left" label="事件类型" prop="type" width="120">
          <template #default="scope">
            <el-tag effect="dark">
              <div v-if="scope.row.type === 1">充值</div>
              <div v-else-if="scope.row.type === 2">划转</div>
              <div v-else>消费</div>
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="说明" prop="remark" width="480" />
        <el-table-column align="left" label="日期" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
              <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>
              详情
            </el-button>
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

    <el-dialog v-model="detailShow" style="width: 800px" lock-scroll :before-close="closeDetailShow" title="查看详情" :draggable="true" destroy-on-close>
      <el-scrollbar height="550px">
        <el-descriptions column="1" border>
                <el-descriptions-item label="用户ID">
                        {{ formData.uid }}
                </el-descriptions-item>
                <el-descriptions-item label="用户名">
                        {{ formData.username }}
                </el-descriptions-item>
                <el-descriptions-item label="积分">
                        {{ formData.recharge }}
                </el-descriptions-item>
                <el-descriptions-item label="事件ID">
                        {{ formData.eventId }}
                </el-descriptions-item>
                <el-descriptions-item label="事件类型">
                        {{ formData.type }}
                </el-descriptions-item>
                <el-descriptions-item label="说明">
                        {{ formData.remark }}
                </el-descriptions-item>
        </el-descriptions>
      </el-scrollbar>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  findUserWallet,
  getUserWalletSelf,
  getUserWalletList
} from '@/api/userWallet'

// 全量引入格式化工具 请按需保留
import {
  getDictFunc,
  formatDate,
  formatBoolean,
  filterDict,
  ReturnArrImg,
  onDownloadFile,
  formatMoney
} from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import CenterCard from "@/view/vbox/dashboard/dataCenterComponents/centerCard.vue";
import {InfoFilled} from "@element-plus/icons-vue";

defineOptions({
    name: 'UserWallet'
})
const walletCustomStyle = ref({
  background: 'linear-gradient(to right, #22111a, #606266)',
  color: '#FFF',
  height: '140px',
})

const userBalance = ref(0)

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        uid: 0,
        username: '',
        recharge: 0,
        eventId: '',
        type: 0,
        remark: '',
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
  const table = await getUserWalletList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  const balanceVal = await getUserWalletSelf()
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
  if (balanceVal.code === 0) {
    userBalance.value = balanceVal.data.balance
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
}

// 获取需要的字典 可能为空 按需保留
setOptions()

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

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
  const res = await findUserWallet({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.reuserWallet
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
          uid: 0,
          username: '',
          recharge: 0,
          eventId: '',
          type: 0,
          remark: '',
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
        username: '',
        recharge: 0,
        eventId: '',
        type: 0,
        remark: '',
        }
}
</script>

<style lang="scss" scoped>
.card {
  width: 100%;/* 设置卡片宽度 */
  height: 50px; /* 设置卡片高度 */
  background-color: #119b8b; /* 设置卡片背景颜色 */
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.19); /* 设置卡片阴影效果 */
  display: flex; /* 使用flex布局 */
  align-items: center; /* 使内容垂直居中 */
  justify-content: center; /* 使内容水平居中 */
  text-align: center; /* 设置文字水平居中 */
  color: #FFFFFF;
  font-size: 20px;
}

.acc-container{
  color: #FFFFFF;
}

.indicator {
  display: flex;
  justify-content: space-around; // 使子元素水平居中展开
  padding: 15px;
  border-radius: 8px; // 添加圆角
}

.indicator span {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 10px;
  &:not(:last-child) {
    border-right: 2px solid #fff; // 白色边框
    margin-right: 15px; // 调整间距
  }
}

.label {
  color: #F5F5F5;
  font-size: 14px;
}

.value {
  color: #FFFFFF;
  font-size: 30px;
  font-weight: bold;
  margin-top: 5px; // 调整间距
}
</style>

<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
 
        <el-form-item label="用户名" prop="username">
         <el-input v-model="searchInfo.username" placeholder="搜索条件" />

        </el-form-item>
        <!-- <el-form-item label="账号名" prop="acAccount">
         <el-input v-model="searchInfo.acAccount" placeholder="搜索条件" />

        </el-form-item> -->
        <!-- <el-form-item label="账户备注" prop="acRemark">
         <el-input v-model="searchInfo.acRemark" placeholder="搜索条件" />

        </el-form-item> -->
        <el-form-item label="店铺备注" prop="shopRemark">
           <el-input v-model="searchInfo.shopRemark" placeholder="搜索条件" />
          </el-form-item>

        <el-form-item label="通道code" prop="channelCode">
         <el-input v-model="searchInfo.channelCode" placeholder="搜索条件" />
        </el-form-item>
        <!-- <el-form-item label="产品ID" prop="productId">
         <el-input v-model="searchInfo.productId" placeholder="搜索条件" />

        </el-form-item> -->
        <!-- <el-form-item label="产品名称" prop="productName">
         <el-input v-model="searchInfo.productName" placeholder="搜索条件" />
        </el-form-item> -->
        <el-form-item label="天" prop="dt">
         <el-input v-model="searchInfo.dt" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
 
        </div>
        <!--  历史数据 -->
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
        >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="天" prop="dt" width="150" />
        <!-- <el-table-column align="left" label="日期" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column> -->
        <!-- <el-table-column align="left" label="用户id" prop="uid" width="120" /> -->
        <el-table-column align="left" label="用户名" prop="username" width="180" >
          <template #default="scope">
            <!-- {{ scope.row.username }} -->
            <el-button type="text"  @click="getAccDetails(scope.row)">
                
                <el-icon style="margin-right: 5px">
                  <InfoFilled/>
                </el-icon>
                {{ scope.row.username }}
              </el-button>
          </template>
        </el-table-column>
        <!-- <el-table-column align="left" label="账号ID" prop="acId" width="120" /> -->
        <!-- <el-table-column align="left" label="通道账户名" prop="acAccount" width="180" />
        <el-table-column align="left" label="账户备注" prop="acRemark" width="180" /> -->
        <el-table-column align="left" label="店铺备注" prop="shopRemark" width="150" />
        <el-table-column align="left" label="通道code" prop="channelCode" width="150" />
        <!-- <el-table-column align="left" label="产品ID" prop="productId" width="120" /> -->
        <el-table-column align="left" label="产品名称" prop="productName" width="180" />
        <el-table-column align="left" label="订单量" prop="orderQuantify" width="100" />
        <el-table-column align="left" label="成功订单量" prop="okOrderQuantify" width="120" />
        <el-table-column align="left" label="成交率" prop="ratio" width="120" >
          <template #default="scope">
            <span>{{scope.row.ratio}}%</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="成交金额" prop="income" width="120" />
        <!-- <el-table-column align="left" label="操作" min-width="120">
            <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
                <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>
                详情
            </el-button>
            <el-button type="primary" link icon="edit" class="table-button" @click="updateBdaChShopIndexDFunc(scope.row)">变更</el-button>
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

    <!-- <el-dialog v-model="accDetailShow" style="width: 800px" lock-scroll :before-close="closeAccDetailShow" title="查看详情" destroy-on-close> -->

    <!-- </el-dialog> -->
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :draggable="true" :title="type==='create'?'添加':'修改'" destroy-on-close>
        <el-scrollbar height="500px">
            <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
              <el-form-item label="用户id:"  prop="uid" >
                <el-input v-model.number="formData.uid" :clearable="true" placeholder="请输入用户id" />
              </el-form-item>
              <el-form-item label="用户名:"  prop="username" >
                <el-input v-model="formData.username" :clearable="true"  placeholder="请输入用户名" />
              </el-form-item>
              <el-form-item label="通道ID:"  prop="cid" >
                <el-input v-model="formData.cid" :clearable="true"  placeholder="请输入通道ID" />
              </el-form-item>
              <el-form-item label="店铺备注:"  prop="shopRemark" >
                <el-input v-model="formData.shopRemark" :clearable="true"  placeholder="请输入店铺备注" />
              </el-form-item>
              <el-form-item label="通道code:"  prop="channelCode" >
                <el-input v-model="formData.channelCode" :clearable="true"  placeholder="请输入通道code" />
              </el-form-item>
              <el-form-item label="产品ID:"  prop="productId" >
                <el-input v-model="formData.productId" :clearable="true"  placeholder="请输入产品ID" />
              </el-form-item>
              <el-form-item label="产品名称:"  prop="productName" >
                <el-input v-model="formData.productName" :clearable="true"  placeholder="请输入产品名称" />
              </el-form-item>
              <el-form-item label="订单量:"  prop="orderQuantify" >
                <el-input v-model.number="formData.orderQuantify" :clearable="true" placeholder="请输入订单量" />
              </el-form-item>
              <el-form-item label="成功订单量:"  prop="okOrderQuantify" >
                <el-input v-model.number="formData.okOrderQuantify" :clearable="true" placeholder="请输入成功订单量" />
              </el-form-item>
              <el-form-item label="成交率:"  prop="ratio" >
                <el-input-number v-model="formData.ratio"  style="width:100%" :precision="2" :clearable="true"  />
              </el-form-item>
              <el-form-item label="成交金额:"  prop="income" >
                <el-input v-model.number="formData.income" :clearable="true" placeholder="请输入成交金额" />
              </el-form-item>
              <el-form-item label="天:"  prop="dt" >
                <el-input v-model="formData.dt" :clearable="true"  placeholder="请输入天" />
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

    <el-dialog v-model="detailShow" style="width: 1400px" lock-scroll :before-close="closeDetailShow" :title="`店铺数据看板 - 用户: ${formData.username}`"  destroy-on-close>
 <!-- 成单统计 -->
    <el-row :gutter="24">
        <el-col :span="24" :xs="24">
          <div class="flex justify-between items-center flex-wrap" style="margin-left: 10px"><h2>近三日该用户店铺数据概览</h2></div>
        </el-col>

        <el-col :span="8" :xs="24">
          <CenterCard title="今日" :custom-style="order1CustomStyle">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999" />
            </template>
            <template #body>
              <!--              <Order :channel-code="searchInfo.cid"/>-->
              <div class="acc-container">
                <div class="indicator">
                  <span>
                    <div class="label">通道数</div>
                    <div class="value">{{ cardData1.channelCnt }}</div>
                  </span>
                  <span>
                    <div class="label">店铺数</div>
                    <div class="value">{{ cardData1.shopCnt }}</div>
                  </span>
                  <span>
                    <div class="label">成单数</div>
                    <div class="value">{{ cardData1.okOrderCnt }}</div>
                  </span>
                  <span>
                    <div class="label">成单金额</div>
                    <div class="value">{{ cardData1.okIncome }}</div>
                  </span>
<!--                  <span>
                    <div class="label">待付金额</div>
                    <div class="value">{{ formatMoney(nearOneHourRate.x4) }}</div>
                  </span>-->
                </div>
              </div>
            </template>
          </CenterCard>
        </el-col>

        <el-col :span="8" :xs="24">
          <CenterCard title="昨日" :custom-style="order2CustomStyle">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999" />
            </template>
            <template #body>
              <!--              <Order :channel-code="searchInfo.cid"/>-->
              <div class="acc-container">
                <div class="indicator">
                  <span>
                    <div class="label">通道数</div>
                    <div class="value">{{ cardData2.channelCnt }}</div>
                  </span>
                  <span>
                    <div class="label">店铺数</div>
                    <div class="value">{{ cardData2.shopCnt }}</div>
                  </span>
                  <span>
                    <div class="label">成单数</div>
                    <div class="value">{{ cardData2.okOrderCnt }}</div>
                  </span>
                  <span>
                    <div class="label">成单金额</div>
                    <div class="value">{{ cardData2.okIncome }}</div>
                  </span>
                </div>
              </div>
            </template>
          </CenterCard>
        </el-col>

        <el-col :span="8" :xs="24">
          <CenterCard title="两天前" :custom-style="order3CustomStyle">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999" />
            </template>
            <template #body>
              <!--              <Order :channel-code="searchInfo.cid"/>-->
              <div class="acc-container">
                <div class="indicator">
                  <span>
                    <div class="label">通道数</div>
                    <div class="value">{{ cardData3.channelCnt }}</div>
                  </span>
                  <span>
                    <div class="label">店铺数</div>
                    <div class="value">{{ cardData3.shopCnt }}</div>
                  </span>
                  <span>
                    <div class="label">成单数</div>
                    <div class="value">{{ cardData3.okOrderCnt }}</div>
                  </span>
                  <span>
                    <div class="label">成单金额</div>
                    <div class="value">{{ cardData3.okIncome }}</div>
                  </span>
                </div>
              </div>
            </template>
          </CenterCard>
        </el-col>


    </el-row>


<!--   趋势图   -->
<el-row :gutter="24">
        <el-col :span="24" :xs="24">
          <div class="flex justify-between items-center flex-wrap" style="margin-left: 10px"><h2>各个店铺今日成单</h2></div>
        </el-col>
        <el-col :span="12" :xs="24">
          <CenterCard title="近6小时实时成单(金额)">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999" />
            </template>
            <template #body>
              <!--              <lineCharts :channel-code="searchInfo.cid" :start-time="startTimeOneHour" :end-time="endTimeOneHour"-->
              <!--                          interval="5m" keyword="sum" format="HH:mm" unit="元"/>-->
              <StackedLineCharts :chartData="todayIncomeSum" :uid=dialogUid unit="元"/>
            </template>
          </CenterCard>
          <!--          <CenterCard title="近1小时实时成单(数量)" style="grid-column-start: span 2;">-->
        </el-col>
        <el-col :span="12" :xs="24">
          <CenterCard title="近6小时实时成单(数量)">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999" />
            </template>
            <template #body>
              <StackedLineCharts :chartData="todayOkPayCnt" :uid=dialogUid unit="笔"/>
            </template>
          </CenterCard>
        </el-col>
      </el-row>

      <!-- <el-row :gutter="24">
        <el-col :span="24" :xs="24">
          <div class="flex justify-between items-center flex-wrap" style="margin-left: 10px"><h2>各个通道各个账户今日成单</h2></div>
        </el-col>
        <el-col :span="12" :xs="24">
          <CenterCard title="当日实时成单(金额)">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999" />
            </template>
            <template #body>
              <lineCharts :chart-data="nearOneHourSum" format="HH:mm" unit="元"/>
            </template>
          </CenterCard>
        </el-col>
        <el-col :span="12" :xs="24">
          <CenterCard title="当日实时成单(数量)">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999" />
            </template>
            <template #body>
              <lineCharts :chart-data="nearOneHourSum" format="HH:mm" unit="笔"/>
            </template>
          </CenterCard>
        </el-col>
      </el-row> -->

    
<el-row :gutter="24">
        <el-col :span="24" :xs="24">
          <div class="flex justify-between items-center flex-wrap" style="margin-left: 10px"><h2>各个店铺近一周成单</h2></div>
        </el-col>
        <el-col :span="12" :xs="24">
          <CenterCard title="近一周成单(金额)">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999" />
            </template>
            <template #body>
              <StackedLineCharts :chartData="nearWeekIncomeSum" :uid=dialogUid unit="元"/>
            </template>
          </CenterCard>
        </el-col>
        <el-col :span="12" :xs="24">
          <CenterCard title="近一周成单(数量)">
            <template #action>
              <span class="gvaIcon-prompt" style="color: #999" />
            </template>
            <template #body>
              <StackedLineCharts :chartData="nearWeekOkPayCnt" :uid=dialogUid unit="笔"/>
            </template>
          </CenterCard>
        </el-col>
      </el-row>

      <el-row :gutter="24">
        <el-col :span="24" :xs="24">
          <div class="flex justify-between items-center flex-wrap" style="margin-left: 10px"><h2>各个店铺近一周(不含当天)成单详情</h2></div>
        </el-col>
        <el-col :span="24" :xs="24">
          <div class="gva-search-box">
            <el-form ref="elSearchFormRef" :inline="true" :model="formData" class="demo-form-inline"  @keyup.enter="onSubmitAcid">
              <!-- <el-form-item label="通道账户">
                    <el-select
                    v-model="formData.acId"
                    placeholder="请选择通道账号"
                    filterable
                    clearable
                    style="width: 100%"
                    @change="handleAccChange"
                >
                  <el-option
                      v-for="item in accList"
                      :key="item.acAccount"
                      :label="formatJoin(' -- 备注： ', item.acAccount, item.acRemark)"
                      :value="item.acId"
                  />
                </el-select>
              </el-form-item> -->
              <el-form-item label="店铺名" prop="shopRemark">
              <el-input v-model="viewSearchInfo.shopRemark" placeholder="搜索条件" />
              </el-form-item>

              <el-form-item label="通道code" prop="channelCode">
              <el-input v-model="viewSearchInfo.channelCode" placeholder="搜索条件" />
              </el-form-item>
          
              <el-form-item label="产品名称" prop="productName">
              <el-input v-model="viewSearchInfo.productName" placeholder="搜索条件" />
              </el-form-item>

              <el-form-item label="天" prop="dt">
              <el-input v-model="viewSearchInfo.dt" placeholder="搜索条件" />
              </el-form-item>
              
              <el-form-item>
                <el-button type="primary" icon="search" @click="viewOnSubmit">查询</el-button>
                <el-button icon="refresh" @click="viewOnReset">重置</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-col>
        <el-col :span="24" :xs="24">
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="viewTableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
        >
        <!-- <el-table-column type="selection" width="55" /> -->
        <el-table-column align="left" label="天" prop="dt" width="120" />
        <!-- <el-table-column align="left" label="日期" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column> -->
        <!-- <el-table-column align="left" label="用户id" prop="uid" width="120" /> -->
        <el-table-column align="left" label="用户名" prop="username" width="150" >
         
        </el-table-column>
        <!-- <el-table-column align="left" label="账号ID" prop="acId" width="120" /> -->
        <!-- <el-table-column align="left" label="通道账户名" prop="acAccount" width="150" /> -->
        <el-table-column align="left" label="店铺名" prop="shopRemark" width="150" />
        <el-table-column align="left" label="通道code" prop="channelCode" width="100" />
        <!-- <el-table-column align="left" label="产品ID" prop="productId" width="120" /> -->
        <el-table-column align="left" label="产品名称" prop="productName" width="150" />
        <el-table-column align="left" label="订单量" prop="orderQuantify" width="100" />
        <el-table-column align="left" label="成功订单量" prop="okOrderQuantify" width="120" />
        <el-table-column align="left" label="成交率" prop="ratio" width="120" >
          <template #default="scope">
            <span>{{scope.row.ratio}}%</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="成交金额" prop="income" width="120" />

        </el-table>
        <div class="gva-pagination">
            <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="viewPage"
            :page-size="viewPageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="viewTotal"
            @current-change="viewHandleCurrentChange"
            @size-change="viewHandleSizeChange"
            />
        </div>
      </el-col>
      </el-row>

  
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createBdaChShopIndexD,
  deleteBdaChShopIndexD,
  deleteBdaChShopIndexDByIds,
  updateBdaChShopIndexD,
  findBdaChShopIndexD,
  getBdaChShopIndexDList,
  getBdaChShopIndexDListWeek,
  getBdaChShopIndexDUesrOverview,
  getBdaChShopIndexToDayIncome,
  getBdaChShopIndexToDayInOkCnt,
  getBdaChShopIndexToWeekIncome,
  getBdaChShopIndexToWeekInOkCnt
} from '@/api/bdaChshopIndexD'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import {calculatePercentage, formatMoney} from '@/utils/format';
import {getPayOrderOverview, getPayOrderRate} from "@/api/payOrder";


import CenterCard from '../centerCard.vue'
import lineCharts from '../lineCharts.vue'
import StackedLineCharts from '../stackedLineCharts.vue'

defineOptions({
    name: 'BdaChShopIndexD'
})

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
          uid: 0,
          username: '',
          cid: '',
          shopRemark: '',
          channelCode: '',
          productId: '',
          productName: '',
          orderQuantify: 0,
          okOrderQuantify: 0,
          ratio: 0,
          income: 0,
          dt: '',
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
  const table = await getBdaChShopIndexDListWeek({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
            deleteBdaChShopIndexDFunc(row)
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
      const res = await deleteBdaChShopIndexDByIds({ ids })
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
const updateBdaChShopIndexDFunc = async(row) => {
    const res = await findBdaChShopIndexD({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data.rebdaChaccIndexD
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteBdaChShopIndexDFunc = async (row) => {
    const res = await deleteBdaChShopIndexD({ ID: row.ID })
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
const dialogUid = ref(0)

// 查看详情控制标记
const detailShow = ref(false)

// 打开详情弹窗
const openDetailShow = () => {
  getDetails(formData)
  detailShow.value = true
  console.log('row==>' + JSON.stringify(formData.value))
  dialogUid.value = formData.value.uid
  
  
}

const cardList = ref([])
const cardData1 = ref({
        uid: 0,
        shopCnt: 0,
        channelCnt: 0,
        okOrderCnt: 0,
        okIncome: 0,
        dt: '',
        })
const cardData2 = ref({
        uid: 0,
        shopCnt: 0,
        channelCnt: 0,
        okOrderCnt: 0,
        okIncome: 0,
        dt: '',
})
const cardData3 = ref({
        uid: 0,
        shopCnt: 0,
        channelCnt: 0,
        okOrderCnt: 0,
        okIncome: 0,
        dt: '',
})
const todayOkPayCnt = ref()

const todayIncomeSum = ref()
const nearWeekIncomeSum = ref()
const nearWeekOkPayCnt = ref()

const getAccIncomeSumData = async(row) => {
  console.log('getAccIncomeSumData=',row.value)
  let resData = await getBdaChShopIndexToDayIncome({uid : row.value.uid})
  console.log('getAccIncomeSumData:', JSON.stringify(resData))
  todayIncomeSum.value = resData
}


// 打开详情
const getDetails = async (row) => {

  let resData = await getBdaChShopIndexToDayIncome({uid : row.value.uid})
  // console.log('getBdaChShopIndexToDayIncome:', JSON.stringify(resData))
  if (resData.code === 0) {
    todayIncomeSum.value = resData
  }else{
    todayIncomeSum.value = defaultShow
  }

  let resCntData = await getBdaChShopIndexToDayInOkCnt({uid : row.value.uid})
  // console.log('getBdaChShopIndexToDayInOkCnt:', JSON.stringify(resCntData))
  if (resCntData.code === 0) {
    todayOkPayCnt.value = resCntData
  }else{
    todayOkPayCnt.value = defaultShow
  }

  let resWeekIncomeData = await getBdaChShopIndexToWeekIncome({uid : row.value.uid})
  // console.log('getBdaChShopIndexToWeekIncome:', JSON.stringify(resCntData))
  if (resWeekIncomeData.code === 0) {
    nearWeekIncomeSum.value = resWeekIncomeData
  }
  
  let resWeekOkCntData = await getBdaChShopIndexToWeekInOkCnt({uid : row.value.uid})
  // console.log('getBdaChShopIndexToWeekInOkCnt:', JSON.stringify(resCntData))
  if (resWeekOkCntData.code === 0) {
    nearWeekOkPayCnt.value = resWeekOkCntData
  }
  

  // 打开弹窗
  // const res = await findBdaChShopIndexD({ ID: row.ID })
  // if (res.code === 0) {
  //   formData.value = res.data.rebdaChaccIndexD
  //   openDetailShow()
  // }
  // console.log('getDetails ==>' + JSON.stringify(row.value))
  const res = await getBdaChShopIndexDUesrOverview(row.value)
  // console.log('getDetails res ==>' + JSON.stringify(res))
  if (res.code === 0) {
    cardList.value = res.data.list
    if (cardList.value.length >= 1) {
      cardData1.value = cardList.value[0];
    }

    if (cardList.value.length >= 2) {
      cardData2.value = cardList.value[1];
    }

    if (cardList.value.length >= 3) {
      cardData3.value = cardList.value[2];
    }

  }

  getViewTableData()
  // getAccIncomeSumData(row)
}

let defaultShow = ref({"legendData":["09966548",""],"xAxisData":["00:00","00:05","00:10","00:15","00:20","00:25","00:30","00:35","00:40","00:45","00:50","00:55","01:00","01:05","01:10","01:15","01:20","01:25","01:30","01:35","01:40","01:45","01:50","01:55","02:00","02:05","02:10","02:15","02:20","02:25","02:30","02:35","02:40","02:45","02:50","02:55","03:00","03:05","03:10","03:15","03:20","03:25","03:30","03:35","03:40","03:45","03:50","03:55","04:00","04:05","04:10","04:15","04:20","04:25","04:30","04:35","04:40","04:45","04:50","04:55","05:00","05:05","05:10","05:15","05:20","05:25","05:30","05:35","05:40","05:45","05:50","05:55","06:00","06:05","06:10","06:15","06:20","06:25","06:30","06:35","06:40","06:45","06:50","06:55","07:00","07:05","07:10","07:15","07:20","07:25","07:30","07:35","07:40","07:45","07:50","07:55","08:00","08:05","08:10","08:15","08:20","08:25","08:30","08:35","08:40","08:45","08:50","08:55","09:00","09:05","09:10","09:15","09:20","09:25","09:30","09:35","09:40","09:45","09:50","09:55","10:00","10:05","10:10","10:15","10:20","10:25","10:30","10:35","10:40","10:45","10:50","10:55","11:00","11:05","11:10","11:15","11:20","11:25","11:30","11:35","11:40","11:45","11:50","11:55","12:00","12:05","12:10","12:15","12:20","12:25","12:30","12:35","12:40","12:45","12:50","12:55","13:00","13:05","13:10","13:15","13:20","13:25","13:30","13:35","13:40","13:45","13:50","13:55","14:00","14:05","14:10","14:15","14:20","14:25","14:30","14:35","14:40","14:45","14:50","14:55","15:00","15:05","15:10","15:15","15:20","15:25","15:30","15:35","15:40","15:45","15:50","15:55","16:00","16:05","16:10","16:15","16:20","16:25","16:30","16:35","16:40","16:45","16:50","16:55","17:00","17:05","17:10","17:15","17:20","17:25","17:30","17:35","17:40","17:45","17:50","17:55","18:00","18:05","18:10","18:15","18:20","18:25","18:30","18:35","18:40","18:45","18:50","18:55","19:00","19:05","19:10","19:15","19:20","19:25","19:30","19:35","19:40","19:45","19:50","19:55","20:00","20:05","20:10","20:15","20:20","20:25","20:30","20:35","20:40","20:45","20:50","20:55","21:00","21:05","21:10","21:15","21:20","21:25","21:30","21:35","21:40","21:45","21:50","21:55","22:00","22:05","22:10","22:15","22:20","22:25","22:30","22:35","22:40","22:45","22:50","22:55","23:00","23:05","23:10","23:15","23:20","23:25","23:30","23:35","23:40","23:45","23:50","23:55"],"seriesData":[{"name":"09966548","type":"line","stack":"Total","data":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]}]})
// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
            uid: 0,
            username: '',
            cid: '',
            shopRemark: '',
            channelCode: '',
            productId: '',
            productName: '',
            orderQuantify: 0,
            okOrderQuantify: 0,
            ratio: 0,
            income: 0,
            dt: '',
            }
   dialogUid.value = 0
   viewSearchInfo.value = {}
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
            cid: '',
            shopRemark: '',
            channelCode: '',
            productId: '',
            productName: '',
            orderQuantify: 0,
            okOrderQuantify: 0,
            ratio: 0,
            income: 0,
            dt: '',
            }
}
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              let res
              switch (type.value) {
                case 'create':
                  res = await createBdaChShopIndexD(formData.value)
                  break
                case 'update':
                  res = await updateBdaChShopIndexD(formData.value)
                  break
                default:
                  res = await createBdaChShopIndexD(formData.value)
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



// ----
const getAccDetails = async (row) => {
  console.log(row)
  // 打开弹窗
  // const res = await findBdaChShopIndexD({ ID: row.ID })
  // if (res.code === 0) {
  //   formData.value = res.data.rebdaChaccIndexD
  //   // openDetailShow()
  // }
  formData.value.uid = row.uid
  formData.value.username = row.username
  openDetailShow()
  // openDetailShow()
  // 打开弹窗
  // const res = await findChannelAccount({acId: row.acId})
  // if (res.code === 0) {
  //   formAccData.value = res.data.revca
  //   openAccDetailShow()
  // }
  getAccShowTableData()
}



// --趋势图
// 获取当前时间
let localTime = new Date();
// 获取当前时区相对于UTC的偏移量（分钟）
let offset = localTime.getTimezoneOffset();
// 计算东八区相对于UTC的偏移量（分钟）
let easternOffset = 8 * 60;
// let easternOffset = 0;
// 计算东八区当前时间
// let endTime = new Date(localTime.getTime() + (easternOffset + offset) * 60 * 1000);

// 获取当前时间的分钟数
let currentMinute = localTime.getMinutes();
// 计算所在的分钟数分组
let adjustedMinute = Math.ceil(currentMinute / 5) * 5;

// 设置分钟数为分组的结束值，秒和毫秒都设置为0
localTime.setMinutes(adjustedMinute, 0, 0);
// 前一个小时
const startTimeOneHour = new Date(localTime.getTime() + (easternOffset + offset) * 60 * 1000 - 60 * 60 * 1000);
const endTimeOneHour = new Date(localTime.getTime() + (easternOffset + offset) * 60 * 1000);

const nearOneHourSum = ref()

const getAccShowTableData = async() => {
let nearOneHourSumResult = await getPayOrderOverview({ page: 1, pageSize: 9999, orderStatus:1, channelCode: 3000, startTime: Math.floor(startTimeOneHour.getTime() / 1000), endTime: Math.floor(endTimeOneHour.getTime() / 1000), interval:  '5m', keyword:'sum', format: 'HH:mm'})
nearOneHourSum.value = nearOneHourSumResult
}




const order1CustomStyle = ref({
  background: 'linear-gradient(to right, #be2eba, #5b2ecc)',
  color: '#FFF',
  height: '150px',
})
const order2CustomStyle = ref({
  background: 'linear-gradient(to right, #be2eba, #5b2ecc)',
  color: '#FFF',
  height: '150px',
})
const order3CustomStyle = ref({
  background: 'linear-gradient(to right, #be2eba, #5b2ecc)',
  color: '#FFF',
  height: '150px',
})
const order4CustomStyle = ref({
  background: 'linear-gradient(to right, #22111a, #606266)',
  color: '#FFF',
  height: '150px',
})

//------ 账户筛选


// 搜索
const onSubmitAcid = () => {
  console.log("searchInfo.value",searchInfo.value)
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    console.log("elSearchFormRef.value",elSearchFormRef.value)
    getTableData()
  })
}





// =========== view表格控制部分 ===========
const viewPage = ref(1)
const viewTotal = ref(0)
const viewPageSize = ref(10)
const viewTableData = ref([])
const viewSearchInfo = ref({})

// 重置
const viewOnReset = () => {
  viewSearchInfo.value = {}
  getViewTableData()
}

// 搜索
const viewOnSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    viewPage.value = 1
    viewPageSize.value = 10
    getViewTableData()
  })
}

// 分页
const viewHandleSizeChange = (val) => {
  viewPageSize.value = val
  getViewTableData()
}

// 修改页面容量
const viewHandleCurrentChange = (val) => {
  viewPage.value = val
  getViewTableData()
}

// 查询
const getViewTableData = async() => {
  viewSearchInfo.value.uid = dialogUid.value 
  const table = await getBdaChShopIndexDListWeek({ page: viewPage.value, pageSize: viewPageSize.value, ...viewSearchInfo.value })
  if (table.code === 0) {
    viewTableData.value = table.data.list
    viewTotal.value = table.data.total
    viewPage.value = table.data.page
    viewPageSize.value = table.data.pageSize
  }
}



// ============== 表格控制部分结束 ===============
</script>


<style lang="scss" scoped>

.data-center-box{
  width: 100%;
  display: grid;
  grid-template-columns: 2fr 4fr;
  column-gap: 10px;
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
  padding: 10px; // 调整间距

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


<template>
  <div>
    <div class="gva-search-box">
      <el-row :gutter="12">
        <el-col :span="24"></el-col>
        <el-col :span="24">
          <el-col :span="24">
            <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" :rules="searchRule" @keyup.enter="onSubmit"
                     label-width="auto" label-position="right">
              <el-form-item label="查单池账户" prop="acAccount">
                <el-input v-model="searchInfo.acAccount" placeholder="搜索通道账户"/>
              </el-form-item>
              <el-form-item label="账户备注" prop="acRemark">
                <el-input v-model="searchInfo.acRemark" placeholder="搜索备注"/>
              </el-form-item>
              <el-form-item label="账户ID" prop="acId">
                <el-input v-model.number="searchInfo.acId" placeholder="搜索账户ID"/>
              </el-form-item>
              <el-form-item label="通道ID" prop="cid">
                <el-input v-model.number="searchInfo.cid" placeholder="搜索通道ID"/>
              </el-form-item>
              <el-form-item label="开关状态" prop="status">
                <el-select v-model="searchInfo.status" placeholder="选择状态">
                  <el-option label="已开启" value="1"/>
                  <el-option label="已关闭" value="0"/>
                </el-select>
              </el-form-item>
              <el-form-item label="系统状态" prop="sysStatus">
                <el-select v-model="searchInfo.sysStatus" placeholder="选择系统状态">
                  <el-option label="已开启" value="1"/>
                  <el-option label="已关闭" value="0"/>
                </el-select>
              </el-form-item>
              <el-form-item label="归属用户" prop="username">
                <el-input v-model.number="searchInfo.username" placeholder="搜索归属用户"/>
              </el-form-item>
              <el-form-item>
                <el-button icon="refresh" @click="onReset"></el-button>
                <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
              </el-form-item>
            </el-form>
          </el-col>
        </el-col>
      </el-row>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
        <el-popover v-model:visible="deleteVisible" placement="top" width="240">
          <p><span style="color: red;">注意：</span>确定要删除吗？</p>
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
        <el-popover v-model:visible="switchOnVisible" placement="top" width="240">
          <p><span style="color: red;">注意：开启后，查单池账户将进入待使用状态</span>，确定批量开启吗？
          </p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button type="primary" link @click="switchOnVisible = false">取消</el-button>
            <el-button type="primary" @click="onSwitchEnable">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="turn-off" style="margin-left: 10px;" :disabled="!multipleSelection.length"
                       @click="switchOnVisible = true">批量开启
            </el-button>
          </template>
        </el-popover>
        <el-popover v-model:visible="switchOffVisible" placement="top" width="240">
          <p><span style="color: red;">注意：关闭后，查单池账户将无法使用</span>，确定批量关闭吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button type="primary" link @click="switchOffVisible = false">取消</el-button>
            <el-button type="primary" @click="onSwitchDisable">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="open" style="margin-left: 10px;" :disabled="!multipleSelection.length"
                       @click="switchOffVisible = true">批量关闭
            </el-button>
          </template>
        </el-popover>
      </div>

      <el-table ref="multipleTable" tooltip-effect="dark" :data="tableData" row-key="ID" border resizable="true"
                @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55"/>
        <el-table-column align="left" label="ID" prop="acId" width="140">
          <template #default="scope">
            {{ scope.row.acId }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="通道ID" prop="cid" width="80"/>
        <el-table-column align="left" label="账户备注" prop="acRemark" width="180"/>
        <el-table-column align="left" label="通道账户" prop="acAccount" width="180"/>
        <el-table-column align="left" label="CK" prop="token" width="260">
          <template #default="scope">
            <el-input v-model="scope.row.token" :rows="4" readonly="readonly">
              <template #append>
                <el-button type="primary" link icon="edit" @click="updateCaTokenFunc(scope.row)"></el-button>
              </template>
            </el-input>
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态 / 系统开关" prop="status" width="140">
          <template #default="scope">
            <el-row :gutter="12">
              <el-col :span="12">
                <el-popover trigger="hover" placement="top" width="240">
                  <p><span style="color: red;">注意</span>：操作后将通过系统审核，<span style="color: red;">审核通过后开启（或关闭）账号关联资源；</span><span
                      style="color: blue;">未通过系统审核请查看"操作日志"</span>核查原因，确定操作？</p>
                  <template #reference>
                    <el-switch v-model="scope.row.status" inline-prompt :active-value="1" active-text="开启"
                               :inactive-value="0" inactive-text="关闭" size="large"
                               @change="()=>{switchEnable(scope.row)}"/>
                  </template>
                </el-popover>
              </el-col>
              <el-col :span="12">
                <el-switch v-model="scope.row.sysStatus" inline-prompt :active-value="1" active-text="开启"
                           :inactive-value="0" inactive-text="关闭" size="large" disabled/>
              </el-col>
            </el-row>
          </template>
        </el-table-column>
        <el-table-column align="center" label="冷却状态" prop="cbStatus" width="120">
          <template #default="scope">
            <el-button style="width: 80px" :color="formatCDStatusColor(scope.row.cdStatus)">
              {{ formatCDStatus(scope.row.cdStatus) }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作" width="120">
          <template #default="scope">
            <el-row>
              <el-col :span="24">
                <el-button type="primary" link icon="info-filled" class="table-button"
                           @click="getDetails(scope.row)"></el-button>
                <el-button type="primary" link icon="edit" class="table-button"
                           @click="updateCardAccFunc(scope.row)"></el-button>
                <el-button type="warning" link icon="delete" @click="deleteRow(scope.row)"></el-button>
              </el-col>
            </el-row>
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

    <el-dialog v-model="dialogChanFormVisible" :before-close="closeChanDialog" :title="typeTitle" destroy-on-close
               style="width: 1000px">
      <el-scrollbar height="700px">
        <div>
          <div>
            <el-row :gutter="24">
              <el-col :span="24" :xs="24">
                <div v-for="(parent, index) in parentNodes" :key="index" class="card-container">
                  <el-col :span="24" :xs="24">
                    <div class="flex justify-between items-center flex-wrap"
                         style="margin-left: 10px;margin-bottom: -30px"><h2>
                      {{ parent.productName }}</h2></div>
                    <el-divider></el-divider>
                  </el-col>
                  <el-row :gutter="12">
                    <div v-if="parent.children && parent.children.length > 0"
                         style="flex-wrap: wrap;  justify-content: center;display: flex;">
                      <div v-for="(node, childIndex) in parent.children" :key="childIndex">
                        <el-col class="card" :span="24" :xs="24">
                          <CenterCard title="" :custom-style="accCustomStyle">
                            <template #action>
                              <span class="gvaIcon-prompt" style="color: #999"></span>
                            </template>
                            <template #body>
                              <div class="acc-container">
                                <div class="indicator">
                                  <span>
                                    <div class="label">编码</div>
                                    <div class="value">{{ node.channelCode }}</div>
                                  </span><span>
                                    <div class="label">名称</div>
                                    <div class="value">{{ node.productName }}</div>
                                  </span>
                                </div>
                              </div>
                              <el-row :gutter="12">
                                <el-col class="card" :span="12" :xs="24">
                                  <el-button @click="handleProdClick(node)">
                                    添加
                                  </el-button>
                                </el-col>
                                <el-col class="card" :span="12" :xs="24">
                                  <el-button @click="handleProdBatchClick(node)">
                                    批量添加
                                  </el-button>
                                </el-col>
                              </el-row>
                            </template>
                          </CenterCard>
                        </el-col>
                      </div>
                    </div>
                    <div v-else>
                      <el-col class="card" :span="24" :xs="24">
                        <CenterCard title="" :custom-style="accCustomStyle">
                          <template #action>
                            <span class="gvaIcon-prompt" style="color: #999"></span>
                          </template>
                          <template #body>
                            <div class="acc-container">
                              <div class="indicator">
                                  <span>
                                    <div class="label">编码</div>
                                    <div class="value">{{ parent.channelCode }}</div>
                                  </span>
                                <span>
                                    <div class="label">名称</div>
                                    <div class="value">{{ parent.productName }}</div>
                                  </span>
                              </div>
                              <el-row :gutter="12">
                                <el-col class="card" :span="12" :xs="24">
                                  <el-button @click="handleProdClick(parent)">
                                    添加
                                  </el-button>
                                </el-col>
                                <el-col class="card" :span="12" :xs="24">
                                  <el-button @click="handleProdBatchClick(parent)">
                                    批量添加
                                  </el-button>
                                </el-col>
                              </el-row>
                            </div>
                          </template>
                        </CenterCard>
                      </el-col>
                    </div>
                  </el-row>
                </div>
              </el-col>
            </el-row>
          </div>
        </div>
      </el-scrollbar>
    </el-dialog>

    <!--  创建 6000 -->
    <el-dialog v-model="dialog6000FormVisible" :before-close="close6000Dialog" :draggable="true" :title="typeTitle"
               destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="产码方式" prop="type" disabled>
              <el-button v-model="formData.type" readonly disabled>
                {{ formatProdType(Number(formData.type)) }}
              </el-button>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通道ID" prop="cid" disabled>
              <el-input v-model="formData.cid" readonly disabled></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="CK" prop="token">
          <el-input v-model="formData.token" type="textarea" :clearable="true" :rows="4"
                    placeholder="请输入CK" @input="handleTokenInput"/>
        </el-form-item>
        <el-row :gutter="24">
          <el-col :span="24">
            <el-form-item label="通道账户" prop="acAccount">
              <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入账户"/>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="备注" prop="acRemark">
              <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入备注"/>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" active-value="1" inactive-value="0" active-text="开启"
                     inactive-text="关闭"></el-switch>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="returnPreStep(6000)">上一步</el-button>
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  修改 6000 -->
    <el-dialog v-model="dialogUpd6000FormVisible" :before-close="closeUpd6000Dialog" :draggable="true"
               :title="typeTitle" destroy-on-close>
      <el-scrollbar height="300px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
          <el-row>
            <el-col :span="24">
              <el-form-item label="通道ID" prop="cid">
                <el-input v-model="formData.cid" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="账户备注" prop="acRemark">
                <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="通道账户" prop="acAccount">
                <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="CK" prop="token">
            <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入"/>
          </el-form-item>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeUpd6000Dialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  CK  -->
    <el-dialog v-model="dialogTokenFormVisible" :before-close="closeUpdTokenDialog" :draggable="true" title="变更CK"
               destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="token" prop="token">
          <el-input v-model="formData.token" type="textarea" :rows="4" :clearable="true" placeholder="请输入"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeUpdTokenDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 查看详情 -->
    <el-dialog v-model="detailShow" style="width: 800px" lock-scroll :draggable="true" :before-close="closeDetailShow"
               title="查看详情"
               destroy-on-close>
      <el-scrollbar height="550px">
        <el-descriptions :column="6" border>
          <el-descriptions-item label="用户归属" :span="6">{{ formData.username }}</el-descriptions-item>
          <el-descriptions-item label="账户ID" :span="6">{{ formData.acId }}</el-descriptions-item>
          <el-descriptions-item label="通道账户" :span="3">{{ formData.acAccount }}</el-descriptions-item>
          <el-descriptions-item label="账户备注" :span="3">{{ formData.acRemark }}</el-descriptions-item>
          <el-descriptions-item label="ck" :span="6">
            <el-input v-model="formData.token" type="textarea" readonly/>
          </el-descriptions-item>
          <el-descriptions-item label="通道id" :span="6">{{ formData.cid }}</el-descriptions-item>
          <el-descriptions-item label="状态开关" :span="3">{{ formData.status === 0 ? '关闭' : '开启' }}
          </el-descriptions-item>
          <el-descriptions-item label="系统开关" :span="3">{{ formData.sysStatus === 0 ? '关闭' : '开启' }}
          </el-descriptions-item>
        </el-descriptions>
      </el-scrollbar>
    </el-dialog>

    <!--  账号批量添加  -->
    <el-dialog v-model="dialog6000BatchFormVisible" :before-close="close6000BatchDialog" :draggable="true"
               title="批量添加" destroy-on-close>
      <el-form :model="accBatchFormData" label-position="right" ref="elBatchFormRef" :rules="rule" label-width="80px">
        <el-row>
          <el-col :span="24">
            <el-form-item label="通道ID" prop="cid">
              <el-input v-model="accBatchFormData.cid" :clearable="true" placeholder="请输入" disabled/>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <warning-bar
                title="格式说明：一行填写CK信息用回车（按键提示：Enter）,格式示例：__jdu=1222706.......，一次最多添加50个"/>
          </el-col>
          <el-col :span="24">
            <el-form-item label="" prop="acRemark" >
              <!--          <el-input type="textarea" v-model="accBatchFormData.token" :autosize="{ minRows: 5, maxRows: 30 }" placeholder="请输入账号信息，格式内容为：账号备注|CK，一次最多添加50个，回车换行（Enter），一行一条！"/>-->
              <el-row>
                <el-col :span="24">
                  <el-input
                      v-model="inputValue"
                      type="textarea"
                      :rows="6"
                      placeholder="请输入CK，格式为一行对应一个Cookie，回车确认"
                      @keyup.enter="addTag"
                  >
                  </el-input>
                </el-col>
                <el-col :span="24">
                  <el-tag v-for="(tag, index) in tags" :key="index" closable @close="removeTag(index)">
                    <el-popover trigger="hover" placement="top" width="240">
                      <p>{{ tag }}</p>
                      <template #reference>
                        {{ fmtSimpleBodyByWidth(tag, 35) }}
                      </template>
                    </el-popover>
                  </el-tag>
                </el-col>
                <el-col :span="24">
                  <span style="margin-top: 10px">你已经输入了 <span style="color: red">{{ tags.length }} </span>个</span>
                </el-col>
              </el-row>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="close6000BatchDialog">取 消</el-button>
          <el-button type="primary" @click="enterBatchDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createCardAcc,
  deleteCardAcc,
  deleteCardAccByIds,
  updateCardAcc,
  findCardAcc,
  getCardAccList,
  switchEnableCA,
  switchEnableCAByIds,
} from '@/api/cardCheckPool'
import {
  getChannelProductSelf
} from '@/api/channelProduct'
import {codeToText, regionData} from 'element-china-area-data';
import {useRouter} from 'vue-router'

// 全量引入格式化工具 请按需保留
import {
  formatCDStatusColor,
  formatCDStatus,
  formatProdType,
  fmtSimpleBodyByWidth,
  validateCookie, findValuesContainingString,
} from '@/utils/format'
import {ElMessage, ElMessageBox} from 'element-plus'
import {ref, reactive, nextTick} from 'vue'
import WarningBar from "@/components/warningBar/warningBar.vue";
import dayjs from "dayjs";
import utcPlugin from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';
import CenterCard from "@/view/vbox/dashboard/dataCenterComponents/centerCard.vue";

defineOptions({
  name: 'CardAcc'
})

// 注册插件
dayjs.extend(utcPlugin);
dayjs.extend(timezone);

//通道产品

const channelCodeOptions = ref([])
const vcpTableData = ref([])

const channelCodeProps = {
  expandTrigger: 'hover',
  checkStrictly: false,
  emitPath: false,
}

const handleChange = (value) => {
  getOptionData()
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
  type: 0,
  status: 0,
  sysStatus: 0,
  cbStatus: 0,
  username: '',
})

// 验证规则
const rule = reactive({
  token: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
    {
      whitespace: true,
      message: '不能只输入空格',
      trigger: ['input', 'blur'],
    }],
  cid: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
    {
      whitespace: true,
      message: '不能只输入空格',
      trigger: ['input', 'blur'],
    }],
})

const searchRule = reactive({
  createdAt: [
    {
      validator: (rule, value, callback) => {
        if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
          callback(new Error('请填写结束日期'))
        } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
          callback(new Error('请填写开始日期'))
        } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
          callback(new Error('开始日期应当早于结束日期'))
        } else {
          callback()
        }
      }, trigger: 'change'
    }
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

const handleTokenInput = () => {
  const token = formData.value.token;
  const acAccount = formData.value.acAccount
  if (!acAccount) {
    const pinValue = findValuesContainingString(token, "pin");
    if (pinValue) {
      formData.value.acAccount = pinValue;
    }
  }
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async (valid) => {
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
const getTableData = async () => {
  const table = await getCardAccList({page: page.value, pageSize: pageSize.value, ...searchInfo.value})
  const vcpTable = await getChannelProductSelf({page: 1, pageSize: 999, ...searchInfo.value, productId: "card", type: 1})

  if (table.code === 0) {
    tableData.value = table.data.list
    vcpTableData.value = vcpTable.data.list
    //card select
    parentNodes.value = getParentNodes(vcpTableData.value);

    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

// 根据不同的产品类型切换 Option
const payType = ref(0)

const getOptionData = async () => {
  const vcpTable = await getChannelProductSelf({page: 1, pageSize: 999, type: formData.value.type})

  if (vcpTable.code === 0) {
    vcpTableData.value = vcpTable.data.list
    setOptions()
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () => {
  channelCodeOptions.value = []
  setChannelCodeOptions(vcpTableData.value, channelCodeOptions.value, false)
}

const setRegionOptions = (ChannelCodeData, optionsData, disabled) => {
  ChannelCodeData &&
  ChannelCodeData.forEach(item => {
    if (item.children && item.children.length) {
      const option = {
        value: item.code + '',
        label: item.name,
        children: []
      }
      setRegionOptions(
          item.children,
          option.children,
      )
      optionsData.push(option)
    } else {
      const option = {
        value: item.code + '',
        label: item.name,
      }
      optionsData.push(option)
    }
  })
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
    deleteCardAccFunc(row)
  })
}

// 批量删除控制标记
const deleteVisible = ref(false)

// 多选删除
const onDelete = async () => {
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
  const res = await deleteCardAccByIds({ids})
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
const dialogUpd6000FormVisible = ref(false)

const updateCardAccFunc = async (row) => {
  const res = await findCardAcc({ID: row.ID})
  type.value = 'update'
  typeTitle.value = '修改'
  if (res.code === 0) {
    formData.value = res.data.ret
    let cid = Number(res.data.ret.cid)
    if (cid >= 6000 && cid <= 6099) {
      dialogUpd6000FormVisible.value = true
    }
  }
}

// 删除行
const deleteCardAccFunc = async (row) => {
  const res = await deleteCardAcc({ID: row.ID})
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
const dialog6000FormVisible = ref(false)
const dialog6000BatchFormVisible = ref(false)

// 上一步
const returnPreStep = (cid) => {
  dialogChanFormVisible.value = true
  if (cid >= 6000 && cid <= 6099) {
    dialog6000FormVisible.value = false
  }
}

// 查看详情控制标记
const detailShow = ref(false)

// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}

// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findCardAcc({ID: row.ID})
  if (res.code === 0) {
    formData.value = res.data.ret
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
  channelCodeOptions.value = []
  dialogChanFormVisible.value = true
}

// 关闭弹窗
const close6000Dialog = () => {
  dialog6000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    uid: 0,
  }
}


const closeUpd6000Dialog = () => {
  dialogUpd6000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    uid: 0,
  }
}

const closeUpdTokenDialog = () => {
  dialogTokenFormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    uid: 0,
  }
}
const closeDialog = () => {
  dialogTokenFormVisible.value = false
  dialog6000FormVisible.value = false
  dialogUpd6000FormVisible.value = false

  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    status: 0,
    sysStatus: 0,
    uid: 0,
  }
}
// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    formData.value.status = Number(formData.value.status)
    let res
    switch (type.value) {
      case 'create':
        res = await createCardAcc(formData.value)
        break
      case 'update':
        res = await updateCardAcc(formData.value)
        break
      default:
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
  });
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
const onSwitchEnable = async () => {
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
const onSwitchDisable = async () => {
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

const updateCaTokenFunc = async (row) => {
  const res = await findCardAcc({ID: row.ID})
  type.value = 'update'
  typeTitle.value = '修改'
  if (res.code === 0) {
    formData.value = res.data.ret
    dialogTokenFormVisible.value = true
  }
}

const updTokenInfo = async (row) => {
  console.log(row)
  caTokenInfo.value = JSON.parse(JSON.stringify(row))
  await nextTick()
  const req = {
    ...caTokenInfo.value
  }
  const res = await updateTokenInfoFunc(req)
  if (res.code === 0) {
    ElMessage({type: 'success', message: `更新CK${req.code === 0 ? '成功' : '失败'}`})
    await getTableData()
  }
}

const updateTokenInfoFunc = async (row) => {
  const res = await findCardAcc({ID: row.ID})
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.ret
    dialogTokenFormVisible.value = true
  }
}

// 充值记录查询
const channelCode = ref("")

const router = useRouter()

// 弹窗控制标记
const dialogChanFormVisible = ref(false)

// 关闭弹窗
const closeChanDialog = () => {
  dialogChanFormVisible.value = false
}

// ---------- 通道卡片 ---------
// 获取所有父节点
// const parentNodes = getParentNodes(list);
const parentNodes = ref([])

function getParentNodes(nodes) {
  const result = [];
  for (const node of nodes) {
    if (node.parentId === "0") {
      result.push(node);
    }
  }
  return result;
}

function handleProdClick(node) {
  console.log(node)
  formData.value.cid = node.channelCode

  dialogChanFormVisible.value = false
  let channelCode = Number(formData.value.cid);
  console.log(formData.value)

  if (channelCode >= 6000 && channelCode < 6099) {
    formData.value.type = 4
    dialog6000FormVisible.value = true
  }
}

// 批量添加e card ck info
const batchFormData = ref({
  acId: '',
  acRemark: '',
  acAccount: '',
  acPwd: '',
  token: '',
  cid: '',
  type: 0,
  status: 0,
  sysStatus: 0,
  cdStatus: 0,
  username: '',
})

function handleProdBatchClick(node) {
  console.log(node)
  formData.value.cid = node.channelCode
  accBatchFormData.value.cid = node.channelCode

  dialogChanFormVisible.value = false
  let channelCode = Number(formData.value.cid);
  console.log(formData.value)

  if (channelCode >= 6000 && channelCode < 6099) {
    formData.value.type = 4
    type.value = 'createBatch'
    tags.value = []
    dialog6000BatchFormVisible.value = true
  }
}

// 关闭弹窗
const close6000BatchDialog = () => {
  dialog6000BatchFormVisible.value = false
  batchFormData.value = {
    channelCode: '',
    ext: '',
    parentId: 0,
    payType: '',
    productId: '',
    productName: '',
    type: false,
  }
}
const elBatchFormRef = ref()
const accBatchFormData = ref({
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
// 弹窗确定
const enterBatchDialog = async () => {
  let res
  switch (type.value) {
    case 'createBatch':
      console.log("我要开始批量添加了", tags.value)
      for (let i = 0; i < tags.value.length; i++) {
        let tag = tags.value[i]
        let formData = {
          acId: '',
          acRemark: '',
          acAccount: findValuesContainingString(tag, 'pin'),
          acPwd: '',
          token: tag,
          cid: accBatchFormData.value.cid,
        }

        let res = await createCardAcc(formData)
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: '创建成功'
          })
        }
      }
      break
    default:
      break
  }
  close6000Dialog()
  getTableData()
}

const inputValue = ref('');
const tags = ref([]);
const customTrim = str => {
  return String(str).replace(/^\s+|\s+$/g, '');
};
const addTag = () => {
  const lines = customTrim(inputValue.value).split('\n');
  lines.forEach(cookie => {
    // let [remark, cookie] = line.split(',');
    // console.log('r',remark)
    console.log('c', cookie)
    // if (remark === '' && (typeof cookie === 'undefined' || cookie === '')) {
    //   console.log(cookie)
    //   return
    // }
    let remark = '';
    if (typeof cookie === 'undefined') {
      remark = findValuesContainingString(cookie, 'pin');
    }
    // if (!validateCookie(cookie)) return;
    const formattedLine = `${customTrim(remark)},${customTrim(cookie)}`;
    if (!tags.value.includes(customTrim(cookie)) && tags.value.length <= 50) {
      tags.value.push(customTrim(cookie));
    }else if (tags.value.length > 50) {
      ElMessage({
        type: 'error',
        message: '最多允许一次性添加50个'
      })
    }
  });
  inputValue.value = '';
};

const removeTag = index => {
  tags.value.splice(index, 1);
};

// 批量添加e card ck info
const accCustomStyle = ref({
  background: 'linear-gradient(to right, #2ecc71, #3498db)',
  color: '#FFF',
  height: '140px',
})

</script>

<style lang="scss" scoped>
.region-card-container {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  width: 100%;
}

.region-card {
  margin: 10px;
  width: 250px;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  text-align: center;
  color: white;
  position: relative;
  transition: transform 0.3s;
}

.region-card:hover {
  transform: translateY(-5px);
}

.region-tag {
  position: absolute;
  top: 10px;
  left: 10px;
  color: white;
  padding: 5px 10px;
  border-radius: 5px;
}

.region-name {
  margin: 0;
  font-size: 14px;
}

.region-code {
  margin: 0;
  font-size: 12px;
}


.region-title {
  margin-bottom: 10px;
}

.region-title h2 {
  font-size: 20px;
  color: white;
}

.region-title p {
  font-size: 14px;
}

.region-business-data {
  display: flex;
  justify-content: space-around;
  margin-top: 15px;
}

.region-data-item {
  flex: 1;
}

.region-label {
  font-size: 14px;
}

.region-value {
  padding-top: 5px;
  font-size: 18px;
  font-weight: bold;
}

.tab {
  margin-bottom: 20px;
}

.tab h2 {
  cursor: pointer;
  padding: 10px;
  background-color: #ccc;
  margin: 0;
}

.card-container {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
}

.card {
  width: 260px;
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
    margin-right: 10px; // 调整间距
  }
}

.acc-container {
  color: #FFFFFF;
}

.label {
  color: #F5F5F5;
  font-size: 16px;
}

.value {
  color: #FFFFFF;
  font-size: 16px;
  font-weight: bold;
  margin-top: 15px; // 调整间距
}

.scrolling-text {
  height: 30px; /* 设置显示区域的高度 */
  overflow: hidden; /* 隐藏超出显示区域的内容 */
  position: relative; /* 设置为相对定位，以便在其中添加绝对定位的子元素 */
}

.scrolling-text ul {
  list-style-type: none; /* 移除列表默认样式 */
  padding: 0;
  margin: 0;
  animation: scroll-text 4s linear infinite; /* 使用动画实现滚动效果，20s表示滚动完成需要的时间，可根据需要调整 */
}

@keyframes scroll-text {
  0% {
    transform: translateY(0); /* 初始位置在顶部 */
  }
  100% {
    transform: translateY(-100%); /* 最终位置在顶部的上方，根据行数和行高进行计算 */
  }
}

.input-container .el-input__prepend {
  display: flex;
  flex-wrap: wrap;
}

.input-container .el-input__prepend .el-tag {
  margin-right: 5px;
  margin-bottom: 5px;
}
</style>

<template>
  <div class="organization">
    <!-- <div class="gva-search-box org-top">
      组织人员管理
    </div> -->
    <div class="gva-search-box org-top" >
      <span class="demonstration">选择查看日期：</span>
      <el-date-picker
        v-model="dt"
        type="date"
        placeholder="选择日期"
        :size="size"
        value-format="YYYY-MM-DD"
        @change="selectDt"
        :disabled-date="disabledDate"
      />
    </div>
    <div class="gva-organization-box">
      <div class="gva-organization-box-left">
        <!-- <div class="toolbar">
          <el-button type="primary" @click="addOrg(0)">新增组织</el-button>
        </div> -->
        <div class="tree-body">
          <el-tree
            ref="treeRef"
            :data="treeData"
            node-key="ID"
            :props="defaultProps"
            lazy
            :load="loadDeptData"
            highlight-current
            @current-change="getNowOrg"
            default-expand-all
          >
            <template #default="{ node ,data }">
              <span class="tree-body-tree-node">
                <el-tooltip
                  class="box-item"
                  effect="dark"
                  :show-after="600"
                  :content="node.label"
                  placement="top-start"
                >
                  <span class="tree-body-tree-node-label">{{ node.label }}</span>
                </el-tooltip>
                <!-- <span>
                  <el-dropdown>
                    <span class="el-more-filled">
                      <el-icon><more-filled /></el-icon>
                    </span>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item @click="addOrg(data.ID)">
                          <el-icon><Plus /></el-icon>
                          新增子组织</el-dropdown-item>
                        <el-dropdown-item @click="editOrg(data)">
                          <el-icon><CirclePlus /></el-icon>
                          编辑组织</el-dropdown-item>
                        <el-dropdown-item @click="deleteOrg(data)">
                          <el-icon><Delete /></el-icon>
                          删除组织</el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </span> -->
              </span>
            </template>
          </el-tree>
        </div>
      </div>
      <div class="gva-organization-box-right">
        <div class="toolbar">
          <div class="toolbar-search">
            <el-form :inline="true">
             
                
                <!-- <el-input v-model="userSearch.username" placeholder="请输入要搜索的用户名" /> -->
                <!-- <el-form-item label="团队" >
                  <el-select
                    v-model="orgUserForm.sysUserIDS"
                    placeholder="请选择团队"
                    filterable
                    clearable
                  >
                    <el-option
                      v-for="item in userList"
                      :key="item.ID"
                      :disabled="disabledUserMap[item.ID]"
                      :label="item.nickname"
                      :value="item.ID"
                    />
                  </el-select>
                </el-form-item> -->
                <el-form-item label="产品" >
                  <!-- <el-select
                    v-model="orgSelectForm.channelCode"
                    placeholder="请选择产品"
                    filterable
                    clearable
                  >
                    <el-option
                      v-for="item in userList"
                      :key="item.ID"
                      :disabled="disabledUserMap[item.ID]"
                      :label="item.nickname"
                      :value="item.ID"
                    />
                  </el-select> -->
                  <el-cascader
                      v-model="orgSelectForm.cid"
                      :options="channelCodeOptions"
                      :props="channelCodeProps"
                      @change="handleChange"
                      style="width: 100%"
                      clearable
                  />
                </el-form-item>
                <el-form-item label="用户" >
                  <el-select
                    v-model="orgSelectForm.sysUserID"
                    placeholder="请选择成员"
                    filterable
                    clearable
                    @change="handleChangeUid"
                  >
                    <el-option
                      v-for="item in userList"
                      :key="item.ID"
                      :disabled="disabledUserMap[item.ID]"
                      :label="item.nickname"
                      :value="item.ID"
                    />
                  </el-select>
                </el-form-item>
                <el-form-item label="付款账户" >
                  <el-select
                    v-model="orgSelectForm.pAccount"
                    placeholder="请选择付款账户"
                    filterable
                    clearable
                    @change="handleChangePAccount"
                  >
                    <el-option
                      v-for="item in pAccountList"
                      :key="item.ID"
                      :disabled="disabledUserMap[item.ID]"
                      :label="item.pRemark"
                      :value="item.pAccount"
                    />
                  </el-select>
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
                </el-form-item>
              </el-form>
          </div>
          <div>
            <el-button type="primary" @click="orgShow()">团队</el-button>
            <el-button type="primary" @click="cidShow()">产品</el-button>
            <el-button type="primary" @click="uidShow()">用户</el-button>
            <el-button type="primary" @click="paccShow()">付款账户</el-button>
          </div>

        </div>
        <div>
            <!-- 成单统计 -->
            <el-row :gutter="24">
              <el-col :span="24" :xs="24">
                <div class="flex flex-wrap items-center justify-between" style="margin-left: 10px"><h2>该组今日成单数据概览</h2></div>
              </el-col>
              <el-col 
                v-for="(item, index) in cardsData"
                :key="index" 
                :span="8" 
                :xs="24"
              >
                <CenterCard :title="`${item.title}-当天成单数`" :custom-style="order1CustomStyle">
                  <template #action>
                    <span class="gvaIcon-prompt" style="color: #999"/>
                  </template>
                  <template #body>
                    <!--              <Order :channel-code="searchInfo.cid"/>-->
                    <div class="acc-container">
                      <div class="indicator">
                        <span>
                          <div class="label">成交总金额</div>
                          <div class="value">{{ formatMoney(item.income) }}</div>
                        </span>
                        <span>
                          <div class="label">成单数 / 总笔数</div>
                          <div class="value">{{ item.okOrderQuantify }} / {{ item.orderQuantify }}</div>
                        </span>
                        <span>
                          <div class="label">成率</div>
                          <div class="value">{{ calculatePercentage(item.okOrderQuantify,item.orderQuantify) }} % </div>
                        </span>
                      </div>
                    </div>
                  </template>
                </CenterCard>
              </el-col>
            </el-row>
        </div>
        <!-- <div class="table-body">
          <el-table :data="userTable" @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55" />
            <el-table-column prop="sysUser.nickname" label="姓名" width="120" />
            <el-table-column prop="sysUser.authority.authorityName" label="用户角色" width="120" />
            <el-table-column prop="isAdmin" label="是否管理员" width="120">
              <template #default="{row}">
                {{ row.isAdmin?"是":"否" }}
              </template>
            </el-table-column>
            <el-table-column label="操作列" min-width="220">
              <template #default="{row}">
                <el-button link type="primary" @click="openTransferOrgUser(row.sysUser.ID)"> 更换组织</el-button>
                <el-button link type="primary" @click="deleteUser([row.sysUser.ID])"> 踢出组织</el-button>
                <el-button v-if="!row.isAdmin" link type="primary" @click="setAdmin(row.sysUser.ID,true)"> 设置管理员</el-button>
                <el-button v-else link type="primary" @click="setAdmin(row.sysUser.ID,false)"> 取消管理员</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination
              :current-page="page"
              :page-size="pageSize"
              :page-sizes="[10, 30, 50, 100]"
              :total="total"
              layout="total, sizes, prev, pager, next, jumper"
              @current-change="handleCurrentChange"
              @size-change="handleSizeChange"
            />
          </div>
        </div> -->

      </div>
    </div>
   
  </div>
</template>

<script>
import { ref } from 'vue'
import { createOrganization,
  deleteOrganization,
  getOrganizationList,
  updateOrganization,
  findOrganization,
  findOrgUserAll,
  createOrgUser,
  findOrgUserList,
  setOrgUserAdmin,
  deleteOrgUserApi,
  transferOrgUserApi,
  
} from '@/plugin/organization/api/organization'
import {
  createPayAccount,
  deletePayAccount,
  deletePayAccountByIds,
  updatePayAccount,
  findPayAccount,
  getPayAccountList,
  switchEnablePA,
  getPAccGateway,
} from '@/api/payAccount'

import {
  getBdaChorgIndexRealList,
  getBdaChorgIndexRealListBySelect
} from '@/api/bdaChorgIndexD'

import {getChannelProductSelf} from "@/api/channelProduct";
import { getUserList } from '@/api/user.js'
import { ElMessageBox, ElMessage } from 'element-plus'

import CenterCard from '../centerCard.vue'
import lineCharts from '../lineCharts.vue'
import StackedLineCharts from '../stackedLineCharts.vue'
import {calculatePercentage, formatMoney} from "@/utils/format";


export default {
  name: 'OrganizationView',
}
</script>

<script setup>

//  --cards
const cardsData = ref([
  {
          title: '测试团队',
          // uid: 0,
          // username: '',
          // channelCode: '',
          // productId: '',
          // productName: '',
          orderQuantify: 1000,
          okOrderQuantify: 100,
          ratio: 0,
          income: 10000,
          dt: '',
  },
  {
          title: 'laoshang团队',
          orderQuantify: 3000,
          okOrderQuantify: 800,
          ratio: 0,
          income: 100000,
          dt: '',
  }
])
const cardsDataValue = ref({
        title: '',
        // uid: 0,
        // username: '',
        // channelCode: '',
        // productId: '',
        // productName: '',
        orderQuantify: 0,
        okOrderQuantify: 0,
        ratio: 0,
        income: 0,
        dt: '',
})



const dt = ref('')
// 切换选中组织
const selectDt = (e) => {
  dt.value = e
  console.log('dt',dt.value)
  // getParentId()
}
const disabledDate = (time) => {
  return time.getTime() > Date.now()
}

const orgSelectForm = ref({
  sysUserID: '',
  organizationID: 0,
  cid: '',
  pAccount: '',
  dt: dt,
})




// 所有付款列表
const pAccountList = ref([])
// 
const getPAccountList = async(e) => {
  const res = await getPayAccountList({page: page.value, pageSize: pageSize.value, ...orgSelectForm.value})
  if (res.code === 0) {
    pAccountList.value = res.data.list
    console.log('pAccountList:', JSON.stringify(pAccountList.value))
  }
}
const handleChangePAccount = (value) => {

console.log('pAccount:', value)
orgSelectForm.value.pAccount = value
// getTableData()
}

// 用户

const handleChangeUid = (value) => {

console.log('sysUserID:', value)
orgSelectForm.value.sysUserID = value
// getTableData()
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

console.log('cid:', value)
orgSelectForm.value.cid = value
// getTableData()
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

const setOptions = async () => {
  channelCodeOptions.value = []
  setChannelCodeOptions(vcpTableData.value, channelCodeOptions.value, false)
}


// 查询
// 搜索
const onSubmit = () => {
  console.log("searchInfo.value", JSON.stringify(orgSelectForm.value))
  // elSearchFormRef.value?.validate(async (valid) => {
  //   if (!valid) return
  //   console.log("elSearchFormRef.value", elSearchFormRef.value)
  //   getTableData()
  // })
  getYyCardsBySelect()
  // getYyCards()
}


const getYyCardsBySelect = async() => {
  const res = await getBdaChorgIndexRealListBySelect(orgSelectForm.value)
  if (res.code === 0) {
    total.value = res.data.total
    cardsData.value = res.data.list
    console.log('cardsData.value', JSON.stringify(cardsData.value))
  }
}



const getYyCards = async() => {
  console.log('orgSelectForm.value', JSON.stringify(orgSelectForm.value))

  const res = await getBdaChorgIndexRealList(orgSelectForm.value)
  if (res.code === 0) {
    total.value = res.data.total
    cardsData.value = res.data.list
    console.log('cardsData.value', JSON.stringify(cardsData.value))
  }
}

const parentID = ref(0)
const getParentId = async() => {
  console.log('orgSelectForm.value.organizationID', orgSelectForm.value.organizationID)
  const res = await findOrganization({ID: orgSelectForm.value.organizationID})
  // console.log('getParentId.value res', JSON.stringify(res))
  if (res.code === 0) {
    console.log('getParentId.value', JSON.stringify(res.data.reorg.parentID))
    parentID.value = res.data.reorg.parentID
  }
}


const orgShow = async() => {
  getYyCards()
  
}


const cidShow = async() => {
  const selectForm = ref({ 
                sysUserID: '',
                organizationID: orgSelectForm.value.organizationID,
                cid: '1001',
                pAccount: '',
                dt:dt
              } )

 const res = await getBdaChorgIndexRealList(selectForm.value)
  if (res.code === 0) {
    total.value = res.data.total
    cardsData.value = res.data.list
    console.log('cardsData.value', JSON.stringify(cardsData.value))
  }
}

const uidShow = async() => {
  const selectForm = ref({ 
                sysUserID: 10,
                organizationID: orgSelectForm.value.organizationID,
                cid: '',
                pAccount: '',
                dt:dt
              } )

  const res = await getBdaChorgIndexRealList(selectForm.value)
  if (res.code === 0) {
    total.value = res.data.total
    cardsData.value = res.data.list
    console.log('cardsData.value', JSON.stringify(cardsData.value))
  }
}

const paccShow = async() => {
  const selectForm = ref({  
                sysUserID: '',
                organizationID: orgSelectForm.value.organizationID,
                cid: '',
                pAccount: '10000',
                dt:dt
              }) 

  const res = await getBdaChorgIndexRealList(selectForm.value)
  if (res.code === 0) {
    total.value = res.data.total
    cardsData.value = res.data.list
    console.log('cardsData.value', JSON.stringify(cardsData.value))
  }
}



const defaultProps = {
  children: 'children',
  label: 'name',
}

const currentOrg = ref(0)

const data = ref([])

const orgDialog = ref(false)

const selectData = ref([])

const loadDeptData = async(node, resolve) => {
  
  if (node.level === 0) {
    const res = await getOrgTree()
    resolve(res.data.list)
    return
  }
  const res = await getOrganizationList({ parentID: node.data.ID })
  // console.log('loadDeptData',JSON.stringify(res))
  const data = res.data.list
  if (data) {
    resolve(data)
  } else {
    resolve([])
  }
}

// 获取组织树
const treeData = ref([])
const getOrgTree = async() => {
  const res = await getOrganizationList({ parentID: 0 })
  data.value = res.data.list
  const dataMap = {}
  data.value.forEach(item => {
    if (!dataMap[item.parentID]) {
      dataMap[item.parentID] = []
    }
    dataMap[item.parentID].push(item)
  })

  const treeDataOrg = []
  data.value.forEach(item => {
    item.children = dataMap[item.ID]
    if (!item.parentID) {
      treeDataOrg.push(item)
    }
    treeData.value = treeDataOrg
  })
}

// 左侧组织树操作 start

const orgForm = ref({
  ID: 0,
  parentID: 0,
  name: '',
})

// 组织弹窗操作标记
const orgType = ref('add')

// 新增组织或子组织
const addOrg = (ID) => {
  orgForm.value = {
    ID: 0,
    parentID: ID,
    name: '',
  }
  orgType.value = 'add'
  orgDialog.value = true
}

// 组织弹窗关闭
const orgClear = () => {
  orgForm.value = {
    ID: 0,
    parentID: 0,
    name: '',
  }
  orgDialog.value = false
}


const disabledUserMap = ref({})

// 人员操作弹窗
const orgUserDialog = ref(false)




// 切换选中组织
const getNowOrg = (e) => {
  currentOrg.value = e.ID
  // getUserTable()
  orgSelectForm.value.organizationID = currentOrg.value
  getYyCardsBySelect()
  console.log('getNowOrg',currentOrg.value)
  // getParentId()
}

// 组织树组件获取
const treeRef = ref(null)



// 当前组织用户列表
const userTable = ref([])

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)


// 获取当前组织用户列表
const getUserTable = async() => {
  const res = await findOrgUserList({ organizationID: currentOrg.value, page: page.value, pageSize: pageSize.value, ...userSearch.value })
  if (res.code === 0) {
    page.value = res.data.page
    pageSize.value = res.data.pageSize
    total.value = res.data.total
    userTable.value = res.data.list
  }
}

// 所有用户列表
const userList = ref([])

// 组织用户搜索
const userSearch = ref({
  username: '',
})


// 获取所有用户（用于弹窗内选择）
const getAllUser = async(e) => {
  const res = await getUserList({ page: 1, pageSize: 9999 })
  // console.log('getAllUser',JSON.stringify(res))
  userList.value = res.data.list
  total.value = res.data.total
}

// 初始化方法
const init = async() => {
  await getOrgTree()
  treeRef.value.setCurrentKey(data.value[0].ID)
  getAllUser()
  const vcpTable = await getChannelProductSelf({page: 1, pageSize: 9999, ...orgSelectForm.value})
  if (vcpTable.code === 0) {
    vcpTableData.value = vcpTable.data.list
    setOptions()
  }
  getPAccountList()
  currentOrg.value = data.value[0].ID
  orgSelectForm.value.organizationID = data.value[0].ID
  getUserTable()
  getYyCards()
}





init()


const order1CustomStyle = ref({
  background: 'linear-gradient(to right, #2ecc71, #3498db)',
  color: '#FFF',
  height: '120px',
})
</script>

<style scope lang="scss">
.org-top{
  padding-bottom: 20px;
}

.demonstration {
  // display: block;
  color: var(--el-text-color-secondary);
  font-size: 14px;
  margin-bottom: 20px;
}

.gva-organization-box{
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  height: calc(100vh - 260px) ;
  &-left{
    box-sizing: border-box;
    padding: 20px;
    width: 260px;
    background: #fff;
    &>.toolbar{
      margin-bottom: 20px;
    }
     .tree-body-tree-node{
        padding-right: 20px;
        width: 100%;
        display: flex;
        justify-content: space-between;
        flex:1;
        &-label{
          display: inline-block;
          max-width: 100px;
          overflow:hidden;
          text-overflow:ellipsis;
          white-space:nowrap;
            }
      }
    &>.tree-body{
      height: calc(100% - 52px);
      overflow: auto;

      &::-webkit-scrollbar{
        width: 2px;
        height: 2px;
      }
    }
  }
  &-right{
    box-sizing: border-box;
    padding: 20px;
    background: #fff;
     &>.toolbar{
       padding-bottom: 20px;
       display: flex;
       justify-content: space-between;
       .toolbar-search{
          display: flex;
          align-items: center;
          margin-bottom: 20px;
          .el-input{
            flex: 1;
          }
          .el-button{
            margin-left: 10px;
          }
       }
    }
    width: calc(100% - 270px);
  }
}



.data-center-box {
  width: 100%;
  display: grid;
  grid-template-columns: 2fr 4fr;
  column-gap: 10px;
}

.acc-container {
  color: #FFFFFF;
}

.indicator {
  display: flex;
  justify-content: space-around; // 使子元素水平居中展开
  padding: 10px;
  border-radius: 8px; // 添加圆角
}

.indicator span {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px; // 调整间距

  &:not(:last-child) {
    border-right: 2px solid #fff; // 白色边框
    margin-right: 5px; // 调整间距
  }
}

.label {
  color: #F5F5F5;
  font-size: 14px;
}

.value {
  color: #FFFFFF;
  font-size: 22px;
  font-weight: bold;
  margin-top: 5px; // 调整间距
}

.value-small {
  color: #FFFFFF;
  font-size: 20px;
  font-weight: bold;
  margin-top: 5px; // 调整间距
}

</style>

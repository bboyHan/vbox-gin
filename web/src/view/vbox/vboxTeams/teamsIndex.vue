<template>
  <div class="organization">
    <div class="gva-search-box org-top">
      组织管理
    </div>
    <div class="gva-organization-box">
      <div class="gva-organization-box-left">
        <div class="toolbar">
          <el-button type="primary" @click="addOrg(0)">新增团队</el-button>
        </div>
        <div class="tree-body">
          <el-tree
            ref="treeRef"
            :data="treeData"
            show-checkbox
            node-key="ID"
            :props="defaultProps"
            lazy
            :load="loadDeptData"
            highlight-current
            @current-change="getNowOrg"
          >
            <template #default="{ node ,data }">
              <span class="tree-body-tree-node">
                <el-tooltip
                  class="box-item"
                  effect="dark"
                  open-delay
                  
                  :content="node.label"
                  placement="top-start"
                >
                  <span class="tree-body-tree-node-label">{{ node.label }}</span>
                </el-tooltip>
                <span>
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
                </span>
              </span>
            </template>
          </el-tree>
        </div>
      </div>
      <div class="gva-organization-box-right">
        <div class="toolbar">
          <div class="toolbar-search">
            <el-input v-model="userSearch.userName" placeholder="请输入要搜索的用户名" />
            <el-button type="primary" @click="getUserTable">搜索</el-button>
          </div>
          <div>
            <el-button type="primary" @click="openTransferOrgUser()">更换团队</el-button>
            <el-button type="primary" @click="deleteUser(selectData.map(item=>item.uid))">移出团队</el-button>
            <el-button type="primary" @click="addUser">人员入队</el-button>
          </div>
          
        </div>
        <!-- <div>
            <span>
              {{ selectData.map(item=>item) }}
            </span>
          </div> -->
        <div class="table-body">
          <el-table :data="userTable" @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55" />
            <el-table-column prop="userName" label="姓名" width="120" />
            <el-table-column prop="authorityName" label="用户角色" width="120" />
            <el-table-column prop="isAdmin" label="是否队长" width="120">
              <template #default="{row}">
                {{ row.uid==row.leaderId?"是":"否" }}
              </template>
            </el-table-column>
            <el-table-column label="操作列" min-width="220">
              <template #default="{row}">
                <el-button link type="primary" @click="opdendrawer(row)"> 费率详情</el-button>
                <el-button link type="primary" @click="openTransferOrgUser(row.uid)"> 更换团队</el-button>
                <el-button link type="primary" @click="deleteUser([row.uid])"> 移出团队</el-button>
                <el-button v-if="!row.isAdmin" link type="primary" @click="setAdmin(row.uid,true)"> 设置队长</el-button>
                <el-button v-else link type="primary" @click="setAdmin(row.uid,false)"> 取消队长</el-button>
                
                <el-drawer
                title="费率详情"
                size="40%"
                v-if="drawer" 
                v-model="drawer"
                :direction="direction"
                :before-close="handleClose">
                <span><el-tag>{{ row.userName }} 的费率</el-tag></span>
                <div>
                  <el-table
                    :tree-props="{children: 'children', hasChildren: 'hasChildren'}"
                    :data="tableData"
                    row-key="channelCode"
                    style="width: 100%"
                  >
                    <el-table-column align="left" label="通道编码" prop="channelCode" width="120" />
                    <el-table-column align="left" label="产品名称" prop="productName" width="160" />
                    <el-table-column align="left" label="产品ID" prop="productId" width="120" />
                    <el-table-column align="left" label="费率" prop="rate" width="160">
                    <template #default="scope">
                      <el-input
                          v-model="scope.row.rate"
                          :rows="1"
                          readonly="readonly"
                      >
                        <template #append>
                         <!-- <i class="percent-symbol">%</i> -->
                          <el-button type="primary" link icon="edit" @click="upChannelRate(scope.row)"></el-button>
                        </template>
                      </el-input>
                    </template>
                  </el-table-column>
                  </el-table>

                </div>


              </el-drawer>
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
        </div>

      </div>
    </div>
    <el-dialog
      v-model="orgDialog"
      title="组织管理"
    >
      <el-form :model="orgForm" label-width="120px">
        <el-form-item label="父组织">
          <el-select
            v-model="orgForm.parentID"
            filterable
            remote
            reserve-keyword
            placeholder="请选择组织"
            :remote-method="remoteMethod"
            :loading="loading"
          >
            <el-option
              label="根组织"
              :value="0"
            />
            <el-option
              v-for="item in data"
              :key="item.ID"
              :label="item.name"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="组织名称">
          <el-input v-model="orgForm.name" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="orgClear">取消</el-button>
        <el-button type="primary" @click="orgEnter">确认</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="orgUserDialog" title="组织人员">
      <el-form v-model="orgUserForm" label-width="120px">
        <el-form-item label="组织">
          <el-select v-model="orgUserForm.teamID" disabled placeholder="请选择目标组织">
            <el-option
              v-for="item in data"
              :key="item.ID"
              :label="item.name"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="成员">
          <el-select
            v-model="orgUserForm.sysUserIDS"
            placeholder="请选择需要加入组织的成员"
            multiple
            filterable
          >
            <el-option
              v-for="item in userList"
              :key="item.ID"
              :disabled="disabledUserMap[item.ID]"
              :label="item.nickName"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="orgUserClear">取消</el-button>
        <el-button type="primary" @click="orgUserEnter">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="transferOrgUserFlag" title="更换组织">
      <el-form v-model="targetOrg" label-width="120px">
        <el-form-item label="父组织">
          <el-select
            v-model="targetOrg.ToTeamID"
            filterable
            remote
            reserve-keyword
            placeholder="请选择组织"
            :remote-method="remoteMethod"
            :loading="loading"
          >
            <el-option
              label="根组织"
              :value="0"
            />
            <el-option
              v-for="item in data"
              :key="item.ID"
              :label="item.name"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="cloesTransferOrgUser">取消</el-button>
        <el-button type="primary" @click="transferOrgUser">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="dialogRateFormVisible" :before-close="closeRateDialog" :title="'变更费率'"
               destroy-on-close>
      <el-form :model="formData" label-position="right"  label-width="80px">
        <el-form-item label="费率" prop="rate">
          <!-- <el-input v-model="formData.rate" type="textarea" :clearable="true" placeholder="请输入"/> -->
          <el-input-number v-model="formData.rate" :precision="2" :step="0.1" :max="100"></el-input-number>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeRateDialog">取 消</el-button>
          <el-button type="primary" @click="enterRateDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref,nextTick } from 'vue'
import { createVboxTeams,
  deleteVboxTeams,
  getVboxTeamsList,
  updateVboxTeams,
  findVboxTeams,
} from '@/api/vboxTeams.js'
import { 
  findTeamUserAll,
  createVboxTeamsUser,
  findVboxTeamsUser,
  getVboxTeamsUserList,
  // setOrgUserAdmin,
  deleteVboxTeamsUser,
  transferTeamUserApi
} from '@/api/vboxTeamsUser.js'
import{
  createLatestVboxChannelRate,
  createVboxChannelRate,
  getVboxTeamUserChannelRateList
} from '@/api/vboxChannelRate.js'


import {
  getVboxChannelProductList
} from '@/api/vboxChannelProduct'

import { getUserList } from '@/api/user.js'
import { ElMessageBox, ElMessage } from 'element-plus'
import icon from '@/view/superAdmin/menu/icon.vue'
export default {
  components: { icon },
  name: 'Organization',
}
</script>

<script setup>
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
    // console.log('loadDeptData' + JSON.stringify(res.data))
    // resolve(res.data.list)
    return
  }
  const res = await getVboxTeamsList({ parentId: node.data.ID })
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
  const res = await getVboxTeamsList({ parentId: 0 })
  data.value = res.data.list
  console.log('getOrgTree' + JSON.stringify(res.data))
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

// 编辑组织弹窗打开
const editOrg = async(row) => {
  const res = await findVboxTeams({ ID: row.ID })
  if (res.code === 0) {
    orgForm.value = {
      ID: res.data.reorg.ID,
      parentID: res.data.reorg.parentID,
      name: res.data.reorg.name,
    }
    orgType.value = 'edit'
    orgDialog.value = true
  }
}

// 删除组织
const deleteOrg = async(row) => {
  ElMessageBox.confirm('确定删除该组织吗？', '提示', {
    type: 'warning',
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    showClose: false,
  }).then(async() => {
    const res = await deleteVboxTeams({
      ID: row.ID,
    })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      getOrgTree()
    }
  }).catch(() => {
    ElMessage.info('取消删除')
  })
}

// 对选中组织进行编辑或添加操作
const orgEnter = async() => {
  switch (orgType.value) {
    case 'add':
    {
      const res = await createVboxTeams(orgForm.value)
      if (res.code === 0) {
        orgClear()
        getOrgTree()
      }
      break
    }
    case 'edit':
    {
      const res = await updateVboxTeams(orgForm.value)
      if (res.code === 0) {
        orgClear()
        getOrgTree()
      }
      break
    }
    default:
      break
  }
}

// 左侧组织树操作 end

// 右侧人员操作

// 多选人员
const handleSelectionChange = (val) => {
  selectData.value = val
}

const disabledUserMap = ref({})

// 人员操作弹窗
const orgUserDialog = ref(false)

// 人员操作弹窗数据
const orgUserForm = ref({
  sysUserIDS: [],
  teamID: '',
})


// 增加当前组织用户
const addUser = async() => {
  disabledUserMap.value = {}
  const res = await findTeamUserAll({ teamID: currentOrg.value })
  if (res.code === 0) {
    res.data && res.data.forEach(item => {
      disabledUserMap.value[item] = true
    })
  }
  orgUserForm.value.teamID = currentOrg.value
  orgUserDialog.value = true
}

// 切换选中组织
const getNowOrg = (e) => {
  currentOrg.value = e.ID
  getUserTable()
}

// 组织树组件获取
const treeRef = ref(null)

// 设置为管理员
const setAdmin = async(id, flag) => {
  // const res = await setOrgUserAdmin({
  //   sysUserID: id,
  //   isAdmin: flag
  // })
  console.log("设置队长=" + id)
  // if (res.code === 0) {
  //   ElMessage.success('设置成功')
  //   getUserTable()
  // }
}

// 当前组织用户列表
const userTable = ref([])

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const handleCurrentChange = (e) => {
  page.value = e
  getUserTable()
}

const handleSizeChange = (e) => {
  pageSize.value = e
  getUserTable()
}

// 获取当前组织用户列表
const getUserTable = async() => {
  console.log('获取当前组织用户列表currentOrg:' + JSON.stringify(currentOrg))
  const res = await getVboxTeamsUserList({ teamId: currentOrg.value, page: page.value, pageSize: pageSize.value, ...userSearch.value })
  console.log('获取当前组织用户列表:' + JSON.stringify(res))
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
  userName: '',
})

// 踢出组织
const deleteUser = async(sysUserIDS) => {
  ElMessageBox.confirm('确定删除选中用户吗？', '提示', {
    type: 'warning',
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    showClose: false,
  }).then(async() => {
    const res = await deleteVboxTeamsUser({ sysUserIDS, teamID: currentOrg.value })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      getUserTable()
      selectData.value = []
    }
  }).catch(() => {
    ElMessage.info('取消删除')
  })
}

// 获取所有用户（用于弹窗内选择）
const getAllUser = async(e) => {
  const res = await getUserList({ page: 1, pageSize: 9999 })
  userList.value = res.data.list
  total.value = res.data.total
}

// 成员入职功能关闭弹窗
const orgUserClear = () => {
  orgUserForm.value.sysUserIDS = []
  orgUserDialog.value = false
}

// 切换组织功能数据
const targetOrg = ref({
  sysUserIDS: [],
  TeamID: 0,
  ToTeamID: undefined,
})

// 切换组织弹窗标记
const transferOrgUserFlag = ref(false)

// 打开切换组织弹窗
const openTransferOrgUser = (sysUserID) => {
  targetOrg.value.TeamID = currentOrg.value
  if (sysUserID) {
    targetOrg.value.sysUserIDS = [sysUserID]
  } else {
    targetOrg.value.sysUserIDS = selectData.value.map(item => item.sysUser.ID)
  }
  transferOrgUserFlag.value = true
}

// 关闭切换组织弹窗
const cloesTransferOrgUser = () => {
  transferOrgUserFlag.value = false
  targetOrg.value.sysUserIDS = []
  targetOrg.value.ToTeamID = undefined
}

// 成员切换组织
const transferOrgUser = async() => {
  const res = await transferTeamUserApi(targetOrg.value)
  if (res.code === 0) {
    ElMessage.success('转移成功')
    getUserTable()
    selectData.value = []
    cloesTransferOrgUser()
  }
}




// 添加组织成员
const orgUserEnter = async() => {
  orgUserDialog.value = false
  console.log(JSON.stringify(orgUserForm.value))
  const res = await createVboxTeamsUser(orgUserForm.value)
  if (res.code === 0) {
    ElMessage.success('添加成功')
    orgUserClear()
    getUserTable()
  }
}

// 初始化方法
const init = async() => {
  await getOrgTree()
  treeRef.value.setCurrentKey(data.value[0].ID)
  getAllUser()
  currentOrg.value = data.value[0].ID
  getUserTable()
}

init()


const drawer = ref(false);
const direction = ref('rtl');
const userId = ref(0)

const handleClose = (done) => {
  window.confirm('确认关闭？') ? done() : {};
};
const opdendrawer = (row) => {
  drawer.value = true
  console.log('opdendrawer row ==> ' + JSON.stringify(row))
  userId.value = row.uid
  getDataRow.value = row
  getTableData(row)
  // activeRow.value = row

}

//费率展示
const showUesrChannelRate = async() => {
  const res = await transferTeamUserApi(targetOrg.value)
  if (res.code === 0) {
    ElMessage.success('转移成功')
    getUserTable()
    selectData.value = []
    cloesTransferOrgUser()
  }
}

// =========== 表格控制部分 ===========
const dPage = ref(1)
const dTotal = ref(0)
const dPageSize = ref(999)
const tableData = ref([])
const searchInfo = ref({})

// 查询
const getTableData = async(row) => {
  console.log('getTableData row ==> ' + row.uid)
  const table = await getVboxTeamUserChannelRateList({ page: page.value, pageSize: pageSize.value, teamId: row.teamId,uid :row.uid })
  if (table.code === 0) {
    tableData.value = table.data.list
    console.log('==>bb' + JSON.stringify(tableData.value))
    dTotal.value = table.data.total
    dPage.value = table.data.page
    dPageSize.value = table.data.pageSize
  }
}

const dialogRateFormVisible = ref(false)
const upChannelRate = async (row) => {
  console.log(row)
  console.log('==>row1' + JSON.stringify(row))

  // caTokenInfo.value = JSON.parse(JSON.stringify(row))
  await nextTick()
  dialogRateFormVisible.value = true
  formData.value.rate = row.rate
  formRateDataRow.value = row
  // const res = await updateTokenInfoFunc(req)
  // if (res.code === 0) {
  //   ElMessage({type: 'success', message: `更新CK${req.status === 0 ? '成功' : '失败'}`})
  //   await getTableData()
  // }
}

const getDataRow = ref({
})

// rate update form
const formData = ref({
  uid: 0,
  channelCode:0,
  rate: 0
})

const formRateData = ref({
        uid: 0,
        channelCode: '',
        productName: '',
        productId: '',
        rate: 0,
        })
const formRateDataRow = ref({
})

// 点击确定弹窗
const enterRateDialog = () => {
  console.log(JSON.stringify(formData))
  formRateDataRow.value.rate = formData.value.rate
  updateRateForChannel(formRateDataRow)
}

const updateRateForChannel = async (row) => {
  console.log('row2 ==> ' + JSON.stringify(formRateDataRow.value))

  console.log('row3 ==> ' + JSON.stringify(row.value))
  // const res = await findChannelAccount({ID: row.ID})
  formRateData.value.channelCode = row.value.channelCode
  formRateData.value.productName = row.value.productName
  formRateData.value.productId = row.value.productId
  formRateData.value.rate = row.value.rate
  formRateData.value.uid = userId.value


  console.log('row4 ==> ' + JSON.stringify(formRateData.value))
  const res = await createVboxChannelRate(formRateData.value)
  if (res.code === 0) {
    ElMessage.success('添加成功')
    getTableData(getDataRow.value)
    dialogRateFormVisible.value = false
  }
}

// 关闭弹窗
const closeRateDialog = () => {
  dialogRateFormVisible.value = false
  formData.value = {
    uid: 0,
    rate: 0
  }
}
</script>

<style scope lang="scss">
.org-top{
  padding-bottom: 20px;
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
.percent-symbol {
  margin-right: 30px;
}
</style>

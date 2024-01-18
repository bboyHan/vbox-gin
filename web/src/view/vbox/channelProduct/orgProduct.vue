<template>
  <div class="organization">
    <div class="gva-search-box org-top">
      组织产品管理
    </div>
    <div class="gva-organization-box">
      <div class="gva-organization-box-left">
        <div class="toolbar">
          <el-button type="primary" @click="addOrg(0)">新增组织</el-button>
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
            <el-input v-model="productSearch.productName" placeholder="请输入要搜索的通道产品" />
            <el-button type="primary" @click="getProductTable">搜索</el-button>
          </div>
          <div>
            <el-button type="primary" @click="deleteProduct(selectData.map(item=>item.product.ID))">踢出组织</el-button>
            <el-button type="primary" @click="addProduct">产品入队</el-button>
          </div>

        </div>
        <div class="table-body">
          <el-table :data="productTable" @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55" />
            <el-table-column prop="channelCode" label="通道编码" width="120" />
            <el-table-column prop="productName" label="产品" width="120" />
            <el-table-column prop="name" label="归属团队" width="120" />
            <el-table-column prop="oid" label="组织ID" width="120" />
            <el-table-column prop="cpId" label="产品ID" width="120" />
            <el-table-column label="操作列" min-width="220">
              <template #default="{row}">
                <el-button link type="primary" @click="deleteProduct([row.cpId])"> 踢出组织</el-button>
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
      v-model="orgDialog" :draggable="true"
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

    <el-dialog v-model="orgProductDialog" :draggable="true" title="组织产品">
      <el-form v-model="orgProductForm" label-width="120px">
        <el-form-item label="组织">
          <el-select v-model="orgProductForm.organizationID" disabled placeholder="请选择目标组织">
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
            v-model="orgProductForm.channelProductIDS"
            placeholder="请选择需要加入组织的产品"
            multiple
            filterable
          >
            <el-option
              v-for="item in productList"
              :key="item.ID"
              :disabled="disabledProductMap[item.ID]"
              :label="item.productName"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="orgProductClear">取消</el-button>
        <el-button type="primary" @click="orgProductEnter">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref } from 'vue'
import { createOrganization,
  deleteOrganization,
  getOrganizationList,
  updateOrganization,
  findOrganization,
  findOrgProductAll,
  createOrgProduct,
  findOrgProductList,
  deleteOrgProduct,
} from '@/plugin/organization/api/organization'
import { getChannelProductAll } from '@/api/channelProduct'
import { ElMessageBox, ElMessage } from 'element-plus'
export default {
  name: 'Organization',
}
</script>

<script setup>
import { CirclePlus, Delete, MoreFilled, Plus } from '@element-plus/icons-vue';

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

// 编辑组织弹窗打开
const editOrg = async(row) => {
  const res = await findOrganization({ ID: row.ID })
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
    const res = await deleteOrganization({
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
      const res = await createOrganization(orgForm.value)
      if (res.code === 0) {
        orgClear()
        getOrgTree()
      }
      break
    }
    case 'edit':
    {
      const res = await updateOrganization(orgForm.value)
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

const disabledProductMap = ref({})

// 产品操作弹窗
const orgProductDialog = ref(false)

// 产品操作弹窗数据
const orgProductForm = ref({
  channelProductIDS: [],
  organizationID: '',
})

// 增加当前产品至组织
const addProduct = async() => {
  disabledProductMap.value = {}
  const res = await findOrgProductAll({ organizationID: currentOrg.value })
  if (res.code === 0) {
    res.data && res.data.forEach(item => {
      disabledProductMap.value[item] = true
    })
  }
  orgProductForm.value.organizationID = currentOrg.value
  orgProductDialog.value = true
}

// 切换选中组织
const getNowOrg = (e) => {
  currentOrg.value = e.ID
  getProductTable()
}

// 组织树组件获取
const treeRef = ref(null)

// 当前组织用户列表
const productTable = ref([])

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const handleCurrentChange = (e) => {
  page.value = e
  getProductTable()
}

const handleSizeChange = (e) => {
  pageSize.value = e
  getProductTable()
}

// 获取当前组织用户列表
const getProductTable = async() => {
  const res = await findOrgProductList({ organizationID: currentOrg.value, page: page.value, pageSize: pageSize.value, ...productSearch.value })
  if (res.code === 0) {
    page.value = res.data.page
    pageSize.value = res.data.pageSize
    total.value = res.data.total
    productTable.value = res.data.list
  }
}

// 所有产品列表
const productList = ref([])

// 组织用户搜索
const productSearch = ref({
  productName: '',
})

//TODO 踢出组织
const deleteProduct = async(channelProductIDS) => {
  ElMessageBox.confirm('确定删除选中产品吗？', '提示', {
    type: 'warning',
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    showClose: false,
  }).then(async() => {
    const res = await deleteOrgProduct({ channelProductIDS, organizationID: currentOrg.value })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      getProductTable()
      selectData.value = []
    }
  }).catch(() => {
    ElMessage.info('取消删除')
  })
}

// 获取所有产品（用于弹窗内选择）
const getAllProduct = async(e) => {
  const res = await getChannelProductAll({ page: 1, pageSize: 9999 })
  productList.value = res.data.list
  total.value = res.data.total
}

// 产品入队功能关闭弹窗
const orgProductClear = () => {
  orgProductForm.value.channelProductIDS = []
  orgProductDialog.value = false
}

// 添加组织成员
const orgProductEnter = async() => {
  orgProductDialog.value = false
  const res = await createOrgProduct(orgProductForm.value)
  if (res.code === 0) {
    ElMessage.success('添加成功')
    orgProductClear()
    getProductTable()
  }
}

// 初始化方法
const init = async() => {
  await getOrgTree()
  treeRef.value.setCurrentKey(data.value[0].ID)
  getAllProduct()
  currentOrg.value = data.value[0].ID
  getProductTable()
}

init()

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

</style>

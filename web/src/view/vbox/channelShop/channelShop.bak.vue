<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
        <el-form-item label="通道ID" prop="cid">
          <el-cascader
              v-model="searchInfo.cid"
              :options="channelCodeOptions"
              :props="channelCodeProps"
              @change="handleChange"
              style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="店名" prop="shopRemark">
          <el-input v-model="searchInfo.shopRemark" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">新增商铺</el-button>
      </div>
      <div>
        <div>
          <el-tabs tab-position="left" :span="24">
            <el-tab-pane :label="item.productName" v-for="(item, index) in vcpTableData" :key="index">
              <div>
                <el-tabs type="card" tab-position="top" :span="24">
                  <el-tab-pane :label="itemInfo.productName" v-for="(itemInfo, index2) in item.children" :key="index2">
                    <el-row :gutter="12">
                      <el-col v-for="(itemTable, indexTable) in chanMap[itemInfo.channelCode]" :key="indexTable" :span="8">
                        <el-card shadow="hover">
                          <template #header>
                            <el-button type="info" style="font-size: 24px">{{ item.shopRemark }}</el-button>
                            <el-descriptions :title="item.shopRemark" :column="6" border>
                              <template #extra>
                                <el-button link>商铺ID：{{ itemTable.productId }}</el-button>
                              </template>
                              <el-descriptions-item :span="3">
                                <template #label><div>店名</div></template>
                                {{ itemTable.shopRemark }}
                              </el-descriptions-item>
                              <el-descriptions-item :span="3">
                                <template #label><div>通道编码</div></template>
                                {{ itemTable.cid }}
                              </el-descriptions-item>
                              <el-descriptions-item :span="6">
                                <template #label><div>商品数</div></template>
                                <el-tag type="" effect="dark"> {{ lengthFunc(itemTable.list) }} </el-tag>
                              </el-descriptions-item>
                              <el-descriptions-item :span="3">
                                <template #label><div>已开启</div></template>
                                <el-tag type="success" effect="dark"> {{ statusOnCountFunc(itemTable.list) }} </el-tag>
                              </el-descriptions-item>
                              <el-descriptions-item :span="3">
                                <template #label><div>未开启</div></template>
                                <el-tag type="danger" effect="dark"> {{ statusOffCountFunc(itemTable.list) }} </el-tag>
                              </el-descriptions-item>
                              <el-descriptions-item :span="6" align="center">
                                <template #label><div>启用占比</div></template>
                                <el-row>
                                  <el-col :span="12">
                                    <el-progress type="dashboard" :percentage="calPercentage(itemTable.list)">
                                      <template #default="{ percentage }">
                                        <span class="percentage-value">{{ percentage }}%</span>
                                        <span class="percentage-label">Running</span>
                                      </template>
                                    </el-progress>
                                  </el-col>
                                  <el-col :span="12">
                                    <el-row>
                                      <el-col :span="12">
                                        <el-popconfirm @confirm="switchEnableAll(itemTable, 1)" width="220" confirm-button-text="Yes" cancel-button-text="No, Thanks" :icon="InfoFilled" icon-color="#626AEF" title="确定要一键启用所有商品？">
                                          <template #reference>
                                            <el-button type="success" style="margin-top: 10px; margin-bottom: 5px; width: 110px">一键开启</el-button>
                                          </template>
                                        </el-popconfirm>
                                      </el-col>
                                      <el-col :span="12">
                                        <el-popconfirm @confirm="switchEnableAll(itemTable, 0)" width="220" confirm-button-text="Yes" cancel-button-text="No, Thanks" :icon="InfoFilled" icon-color="#626AEF" title="确定要一键启用所有商品？">
                                          <template #reference>
                                            <el-button type="danger" style="margin-top: 10px; margin-bottom: 5px; width: 110px">一键关闭</el-button>
                                          </template>
                                        </el-popconfirm>
                                      </el-col>
                                      <el-col :span="12">
                                        <el-button type="primary" icon="edit" @click="updShopNameDialog(itemTable)" round style="margin-top: 10px; margin-bottom: 5px; width: 110px">店名修改</el-button>
                                      </el-col>
                                      <el-col :span="12">
                                        <el-button round color="#626aef" icon="edit" @click="updDialog(itemTable)" style="margin-top: 10px; margin-bottom: 5px; width: 110px">地址管理</el-button>
                                      </el-col>
                                    </el-row>


                                  </el-col>
                                </el-row>
                              </el-descriptions-item>
                            </el-descriptions>
                          </template>
                          <el-row :gutter="12">
                            <el-button round color="#626aef" icon="search" @click="">统计概览</el-button>
                          </el-row>
                        </el-card>
                      </el-col>
                    </el-row>
                  </el-tab-pane>
                </el-tabs>
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>
    </div>

    <!--  创建商铺  -->
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :draggable="true" :title="typeTitle" destroy-on-close width="60%">
      <el-scrollbar height="500px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="120px">
          <el-row>
            <el-col :span="20">
              <el-form-item label="通道ID" prop="cid">
                <el-cascader
                  v-model="formData.cid"
                  :options="channelCodeOptions"
                  :props="channelCodeProps"
                  @change="handleChange"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
            <el-col :span="20">
              <el-form-item label="店名"  prop="shopRemark" >
                <el-input v-model="formData.shopRemark" :clearable="true"  placeholder="请输入店铺备注" />
              </el-form-item>
            </el-col>
          </el-row>

          <el-card style="{width: 100% !important}" shadow="hover">
            <template #header>
              <div class="card-header">
                <span>商品明细</span>
              </div>
            </template>
            <div>
              <el-table :data="formData.list" style="width: 100%">
                <el-table-column label="地址" prop="address" style="width: 100%">
                  <template #default="scope">
                    <el-input :rows="2" type="textarea" v-if="activeIndex === scope.$index" v-model="scope.row.address"></el-input>
                    <el-input :rows="2" type="textarea" disabled v-model="scope.row.address" readonly v-else></el-input>
                  </template>
                </el-table-column>
                <el-table-column label="金额（元）" prop="money" width="120px">
                  <template #default="scope">
                    <el-input type="number" v-if="activeIndex === scope.$index" v-model.number="scope.row.money" :step="10"></el-input>
                    <span v-else>{{ scope.row.money }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="开关" prop="status" width="100px">
                  <template #default="scope">
                    <el-switch v-if="activeIndex === scope.$index" v-model="scope.row.status" :active-value="1" :inactive-value="0" active-text="开启"
                               inactive-text="关闭" inline-prompt size="large" width="70px"></el-switch>
                    <el-switch v-else v-model="scope.row.status" :active-value="1" :inactive-value="0" active-text="开启"
                               inactive-text="关闭" inline-prompt size="large" width="70px"></el-switch>
                  </template>
                </el-table-column>
                <el-table-column align="right" width="200">
                  <template #header>
                    <el-button type="primary" @click="handleAdd"><Plus style="width:1em; height:1em;" /></el-button>
                  </template>
                  <template #default="scope">
                    <div v-if="activeIndex === scope.$index">
                      <el-button type="primary" @click="handleSave"><Select style="width:1em; height:1em;" /></el-button>
                    </div>
                    <div v-else>
                      <el-button type="success" @click="handleEdit(scope.$index)"><Edit style="width:1em; height:1em;" /></el-button>
                      <el-popconfirm @confirm="handleDelete(scope.$index)" width="220" confirm-button-text="Yes" cancel-button-text="No, Thanks" :icon="InfoFilled" icon-color="#626AEF" title="Are you sure to delete this?">
                        <template #reference>
                          <el-button type="danger"><Delete style="width:1em; height:1em;" /></el-button>
                        </template>
                      </el-popconfirm>
                    </div>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-card>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  修改店名备注  -->
    <el-dialog v-model="dialogUpdShopRemarkFormVisible" :draggable="true" :before-close="closeDialog" :title="typeTitle" destroy-on-close width="20%">
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-row>
          <el-col :span="24">
            <el-form-item label="通道ID" prop="cid" disabled="disabled">
              <el-cascader
                  v-model="formData.cid"
                  :options="channelCodeOptions"
                  :props="channelCodeProps"
                  @change="handleChange"
                  style="width: 100%"
                  disabled
              />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="店铺ID" prop="productId">
              <el-input disabled v-model="formData.productId" />
            </el-form-item>
            <el-form-item label="店名备注" prop="cid">
              <el-input v-model="formData.shopRemark" :clearable="true"  placeholder="请输入店名备注" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  修改商铺  -->
    <el-dialog v-model="dialogUpdFormVisible" :before-close="closeDialog" :draggable="true" :title="typeTitle" destroy-on-close width="60%">
      <el-scrollbar height="500px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
          <el-row>
            <el-col :span="24">
              <el-form-item label="通道ID" prop="cid" disabled="disabled">
                <el-cascader
                  v-model="formData.cid"
                  :options="channelCodeOptions"
                  :props="channelCodeProps"
                  @change="handleChange"
                  style="width: 100%"
                  disabled
                />
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="商铺ID"  prop="productId" >
                <el-input disabled v-model="formData.productId" :clearable="true" placeholder="请输入产品ID" />
              </el-form-item>
            </el-col>
            <el-col :span="16">
              <el-form-item label="店名"  prop="shopRemark">
                <el-input v-model="formData.shopRemark" disabled />
              </el-form-item>
            </el-col>
          </el-row>

          <el-card style="{width: 100% !important}" shadow="never">
            <template #header>
              <div class="card-header">
                <span>商品明细</span>
              </div>
            </template>
            <div>
              <el-table :data="formData.list" style="width: 100%">
                <el-table-column label="地址" prop="address" style="width: 100%">
                  <template #default="scope">
                    <el-input :rows="2" type="textarea" v-if="activeUpdIndex === scope.$index" v-model="scope.row.address"></el-input>
                    <el-input :rows="2" type="textarea" disabled v-model="scope.row.address" readonly v-else></el-input>
                  </template>
                </el-table-column>
                <el-table-column label="金额（元）" prop="money" width="120px">
                  <template #default="scope">
                    <el-input type="number" v-if="activeUpdIndex === scope.$index" v-model.number="scope.row.money" :step="10"></el-input>
                    <span v-else>{{ scope.row.money }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="开关" prop="status" width="100px">
                  <template #default="scope">
                    <el-switch v-if="activeUpdIndex === scope.$index" v-model="scope.row.status" :active-value="1" :inactive-value="0" active-text="开启"
                               inactive-text="关闭" inline-prompt size="large" width="70px" @change="()=>{switchEnable(scope.row)}"></el-switch>
                    <el-switch v-else v-model="scope.row.status" :active-value="1" :inactive-value="0" active-text="开启"
                               inactive-text="关闭" inline-prompt size="large" width="70px" @change="()=>{switchEnable(scope.row)}"></el-switch>
                  </template>
                </el-table-column>
                <el-table-column align="right" width="200">
                  <template #header>
                    <el-button type="primary" @click="handleAdd2Upd"><Plus style="width:1em; height:1em;" /></el-button>
                  </template>
                  <template #default="scope">
                    <div v-if="activeUpdIndex === scope.$index">
                      <el-button type="primary" @click="handleSave2Upd()"><Select style="width:1em; height:1em;" /></el-button>
                    </div>
                    <div v-else>
                      <el-button type="success" @click="handleEdit2Upd(scope.$index)"><Edit style="width:1em; height:1em;" /></el-button>
                      <el-popconfirm @confirm="handleDelete2Upd(scope.$index)" width="220" confirm-button-text="Yes" cancel-button-text="No, Thanks" :icon="InfoFilled" icon-color="#626AEF" title="确定要删除该商品吗？">
                        <template #reference>
                          <el-button type="danger"><Delete style="width:1em; height:1em;" /></el-button>
                        </template>
                      </el-popconfirm>
                    </div>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-card>
        </el-form>
      </el-scrollbar>
    </el-dialog>

    <el-dialog v-model="detailShow" style="width: 800px" lock-scroll :draggable="true" :before-close="closeDetailShow" title="查看详情" destroy-on-close>
      <el-scrollbar height="550px">
        <el-descriptions column="1" border>
          <el-descriptions-item label="用户ID">
            {{ formData.uid }}
          </el-descriptions-item>
          <el-descriptions-item label="通道ID">
            {{ formData.cid }}
          </el-descriptions-item>
          <el-descriptions-item label="产品ID">
            {{ formData.productId }}
          </el-descriptions-item>
          <el-descriptions-item label="店名">
            {{ formData.shopRemark }}
          </el-descriptions-item>
          <el-descriptions-item label="店地址">
            {{ formData.address }}
          </el-descriptions-item>
          <el-descriptions-item label="金额">
            {{ formData.money }}
          </el-descriptions-item>
          <el-descriptions-item label="开关">
            {{ formData.status }}
          </el-descriptions-item>
        </el-descriptions>
      </el-scrollbar>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createChannelShop,
  deleteChannelShop,
  updateChannelShop,
  findChannelShop,
  findChannelShopByProductID,
  getChannelShopList
} from '@/api/channelShop'
import {
  findChannelProduct,
  getChannelProductSelf
} from '@/api/channelProduct'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import {ref, reactive, nextTick} from 'vue'
import {CircleCheck, CircleClose, Delete, Edit, InfoFilled, Plus, Select} from '@element-plus/icons-vue';
import {setUserInfo} from "@/api/user";

defineOptions({
  name: 'ChannelShop'
})

// -------------- 子表编辑(创建) ------------------------
let activeIndex = ref(-1);
// 新增行
const handleAdd = function () {
  let item = {
    address: '',
    money: 10,
    status: 0
  };
  formData.value.list.push(item);
  activeIndex.value = formData.value.list.length - 1;
};
// 编辑行
const handleEdit = (index) => {
  activeIndex.value = index;
};
// 保存行
const handleSave = () => {
  console.log(formData.value.list[activeIndex.value])
  activeIndex.value = -1;
};
// 删除行
const handleDelete = function (index) {
  formData.value.list.splice(index, 1);
};
// -------------- 子表编辑(创建) ------------------------

// -------------- 子表编辑(修改) ------------------------
let activeUpdIndex = ref(-1);
// 新增行
const handleAdd2Upd = function () {
  let item = {
    address: '',
    money: 10,
    status: 0
  };
  formData.value.list.push(item);
  activeUpdIndex.value = formData.value.list.length - 1;
};
// 编辑行
const handleEdit2Upd = (index) => {
  activeUpdIndex.value = index;
};
// 保存行
const handleSave2Upd = () => {
  let create = { ...formData.value}
  let newList = []
  newList.push(formData.value.list[activeUpdIndex.value])
  create.list = newList
  createChannelShop(create)
  activeUpdIndex.value = -1;
};
// 删除行
const handleDelete2Upd = function (index) {
  let ele = formData.value.list[index];
  console.log(ele)
  let id = ele.id;
  if (id) {
    console.log("有id，要删库 -> id: " + id)
    deleteChannelShopFunc({ID: id})
  }else {
    console.log("没id的临时数据，随便删")
  }
  formData.value.list.splice(index, 1);
};
// -------------- 子表编辑(修改) ------------------------


// -------------- 同一通道产品的归集 ------------------------
const chanMap = ref({})
const lengthFunc = (list) => {
  return list.length;
}
const calPercentage = (list) => {
  let c = 0;
  for (let i = 0; i < list.length; i++) {
    if (list[i].status === 1) {
      c++;
    }
  }
  return Math.round(c * 100 / list.length);
}

const statusOnCountFunc = (list) => {
  let c = 0;
  for (let i = 0; i < list.length; i++) {
    if (list[i].status === 1) {
      c++;
    }
  }
  return c
}

const statusOffCountFunc = (list) => {
  let c = 0;
  for (let i = 0; i < list.length; i++) {
    if (list[i].status === 0) {
      c++;
    }
  }
  return c
}
const processChanMap = (list) => {
  if(!list) list= []
  for (let i = 0; i < list.length; i++) {
    const item = list[i];

    if (item.cid in chanMap.value) {
      chanMap.value[item.cid].push(item);
    } else {
      chanMap.value[item.cid] = [item];
    }
  }
}
// -------------- 同一通道产品的归集 ------------------------

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

// ------------ 开关商品 -----------------
// 开关单条
const switchEnable = async(row) => {
  console.log(row)
  // userInfo.value = JSON.parse(JSON.stringify(row))
  await nextTick()
  const req = {
    type: 2,
    ...row
  }
  console.log(req)
  if (req.id){
    const res = await updateChannelShop(req)
    if (res.code === 0) {
      ElMessage({ type: 'success', message: `${req.status === 0 ? '禁用' : '启用'}单条商品成功` })
      await getTableData()
    }
  } else {
    console.log("没id的数据，不操作")
  }
}

// 开关所有
const switchEnableAll = async(row, status) => {
  console.log(row)
  // userInfo.value = JSON.parse(JSON.stringify(row))
  await nextTick()
  const req = {
    type: 3,
    status: status,
    ...row
  }
  console.log(req)
  if (req.productId){
    const res = await updateChannelShop(req)
    if (res.code === 0) {
      ElMessage({ type: 'success', message: `${req.status === 0 ? '禁用' : '启用'}店铺成功` })
      await getTableData()
    }
  } else {
    console.log("没productId的数据，不操作")
  }
}
// ------------ 开关商品 -----------------

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  uid: 0,
  cid: '',
  productId: '',
  shopRemark: '',
  list: [
      {
        address: '',
        money: 0,
        status: 0
      }
  ]
})


// 验证规则
const rule = reactive({
  cid : [
    {
      required: true,
      message: '',
      trigger: ['input','blur'],
    },
    {
      whitespace: true,
      message: '不能只输入空格',
      trigger: ['input', 'blur'],
    }
  ],
})

const searchRule = reactive({})

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

const initForm = () => {
  tableData.value = []
  vcpTableData.value = []
  chanMap.value = {}
}

// 查询
const getTableData = async() => {
  const table = await getChannelShopList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  const vcpTable = await getChannelProductSelf({ page: 1, pageSize: 999, ...searchInfo.value })

  if (table.code === 0) {
    initForm()
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize

    vcpTableData.value = vcpTable.data.list
    setOptions()
    processChanMap(tableData.value)
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
  channelCodeOptions.value = []
  setChannelCodeOptions(vcpTableData.value, channelCodeOptions.value, false)
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

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')
const typeTitle = ref('')

// 更新行
const updateChannelShopFunc = async(row) => {
  const res = await findChannelShopByProductID({ productId: row.productId })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.rechannelShop
    dialogUpdFormVisible.value = true
  }
}


// 删除行
const deleteChannelShopFunc = async (row) => {
  console.log(row)
  const res = await deleteChannelShop({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '删除成功'
    })
    getTableData()
  }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)
const dialogUpdFormVisible = ref(false)
const dialogUpdShopRemarkFormVisible = ref(false)


// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findChannelShop({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.rechannelShop
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
    uid: 0,
    cid: '',
    productId: '',
    shopRemark: '',
    money: 0,
    status: 0,
  }
}

// 打开弹窗
const openDialog = () => {
  type.value = 'create'
  typeTitle.value = '创建商铺'
  dialogFormVisible.value = true
}

// 修改弹窗
const updDialog = (itemTable) => {
  type.value = 'update'
  typeTitle.value = '修改商铺'
  dialogUpdFormVisible.value = true
  formData.value = itemTable
}

// 店名修改弹窗
const updShopNameDialog = (itemTable) => {
  type.value = 'updateShopRemark'
  typeTitle.value = '修改店名备注'
  dialogUpdShopRemarkFormVisible.value = true
  formData.value = itemTable
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  dialogUpdFormVisible.value = false
  dialogUpdShopRemarkFormVisible.value = false
  formData.value = {
    uid: 0,
    cid: '',
    productId: '',
    list: [
        {
          shopRemark: '',
          money: 0,
          status: 0,
        }
    ]
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
      case 'updateShopRemark':
        res = await updateChannelShop({
          type: 1,
          cid: formData.value.cid,
          productId: formData.value.productId,
          shopRemark: formData.value.shopRemark,
        })
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
</script>

<style>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 20px;
  color: #6B7687;
}
.percentage-value {
  display: block;
  margin-top: 10px;
  font-size: 28px;
}
.percentage-label {
  display: block;
  margin-top: 10px;
  font-size: 12px;
}
</style>

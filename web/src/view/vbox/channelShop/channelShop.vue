<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule"
               @keyup.enter="onSubmit">
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
          <el-input v-model="searchInfo.shopRemark" placeholder="搜索条件"/>
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
          <div>
            <el-row :gutter="12">
              <div class="card_list_wrap">
                <div v-for="item in tableData">
                  <div class="card_wrap sm_card hot">
                    <div class="card_header" :style="backgroundImageStyle">
                      <div class="card_header_tag new_customers">
                        <span>通道ID: {{ item.cid }}</span>
                      </div>
                      <div class="card_header_title"><span>{{ item.shopRemark }}</span></div>
                      <div class="card_header_description">
                        <div>商铺数量： <span style="color: red">{{ lengthFunc(item.list) }}</span></div>
                      </div>
                    </div>
                    <div class="card_content">
                      <div class="label">
                        <span class="label_item">商铺ID：{{ item.productId }}</span>
                      </div>
                    </div>
                    <div class="card_footer">
                      <div class="label">
                        <span class="label_item">已开启： {{ statusOnCountFunc(item.list) }}</span>
                        <span class="label_item_err">未开启： {{ statusOffCountFunc(item.list) }}</span>
                      </div>
                      <div class="list_wrap">
                        <div class="percent">
                          <div class="font-num en3">66.6%</div>
                          <div class="percent_des">
                            <span>成率</span>
                            <span class="org_percent">66 / 100</span>
                          </div>
                        </div>
                        <div class="buy btn_wrap">
                          <el-popconfirm @confirm="switchEnableAll(item, 0)" width="220"
                                         confirm-button-text="Yes" cancel-button-text="No, Thanks"
                                         :icon="InfoFilled" icon-color="#626AEF"
                                         title="确定要一键关闭所有商品？">
                            <template #reference>
                              <button class="card button-base red" type="button">一键关闭</button>
                            </template>
                          </el-popconfirm>
                          <el-popconfirm @confirm="switchEnableAll(item, 1)" width="220"
                                         confirm-button-text="Yes" cancel-button-text="No, Thanks"
                                         :icon="InfoFilled" icon-color="#626AEF"
                                         title="确定要一键启用所有商品？">
                            <template #reference>
                              <button class="card button-base yellow" type="button">一键开启</button>
                            </template>
                          </el-popconfirm>
                        </div>
                        <div class="buy btn_wrap">
                          <button class="card button-base yellow" type="button" @click="updShopNameDialog(item)">
                            店名修改
                          </button>
                          <button class="card button-base yellow" type="button" @click="updDialog(item)">地址管理</button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </el-row>
          </div>

        </div>
      </div>
    </div>

    <!--  创建商铺  -->
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :draggable="true" :title="typeTitle"
               destroy-on-close width="80%" draggable overflow>
      <el-scrollbar height="500px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
          <el-row>
            <el-col :span="8">
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
            <el-col :span="16">
              <el-form-item label="店名" prop="shopRemark">
                <el-input v-model="formData.shopRemark" :clearable="true" placeholder="请输入店铺备注"/>
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
                    <el-input :rows="2" type="textarea" v-if="activeIndex === scope.$index"
                              v-model="scope.row.address"></el-input>
                    <el-input :rows="2" type="textarea" disabled v-model="scope.row.address" readonly v-else></el-input>
                  </template>
                </el-table-column>
                <el-table-column label="金额（元）" prop="money" width="120px">
                  <template #default="scope">
                    <el-input type="number" v-if="activeIndex === scope.$index" v-model.number="scope.row.money"
                              :step="10"></el-input>
                    <span v-else>{{ scope.row.money }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="设备" prop="device" width="120px">
                  <template #default="scope">
                    <el-select v-if="activeIndex === scope.$index" v-model="scope.row.device">
                      <el-option label="默认" value="default"/>
                      <el-option label="安卓" value="Android"/>
                      <el-option label="苹果" value="iOS"/>
                    </el-select>
                    <span v-else>{{ scope.row.device }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="标识" prop="markId" width="120px">
                  <template #default="scope">
                    <el-input v-if="activeIndex === scope.$index" v-model="scope.row.markId"></el-input>
                    <span v-else>{{ scope.row.markId }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="开关" prop="status" width="100px">
                  <template #default="scope">
                    <el-switch v-if="activeIndex === scope.$index" v-model="scope.row.status" :active-value="1"
                               :inactive-value="0" active-text="开启"
                               inactive-text="关闭" inline-prompt size="large" width="70px"></el-switch>
                    <el-switch v-else v-model="scope.row.status" :active-value="1" :inactive-value="0"
                               active-text="开启"
                               inactive-text="关闭" inline-prompt size="large" width="70px"></el-switch>
                  </template>
                </el-table-column>
                <el-table-column align="right" width="200">
                  <template #header>
                    <el-button type="primary" @click="handleAdd">
                      <Plus style="width:1em; height:1em;"/>
                    </el-button>
                  </template>
                  <template #default="scope">
                    <div v-if="activeIndex === scope.$index">
                      <el-button type="primary" @click="handleSave"><Select style="width:1em; height:1em;"/></el-button>
                      <el-button type="primary" @click="handleDelete(scope.$index)">
                        <Delete style="width:1em; height:1em;"/>
                      </el-button>
                    </div>
                    <div v-else>
                      <el-button type="success" @click="handleEdit(scope.$index)">
                        <Edit style="width:1em; height:1em;"/>
                      </el-button>
                      <el-popconfirm @confirm="handleDelete(scope.$index)" width="220" confirm-button-text="Yes"
                                     cancel-button-text="No, Thanks" :icon="InfoFilled" icon-color="#626AEF"
                                     title="Are you sure to delete this?">
                        <template #reference>
                          <el-button type="danger">
                            <Delete style="width:1em; height:1em;"/>
                          </el-button>
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
    <el-dialog v-model="dialogUpdShopRemarkFormVisible" :before-close="closeDialog" :draggable="true" :title="typeTitle"
               destroy-on-close
               width="20%">
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
              <el-input disabled v-model="formData.productId"/>
            </el-form-item>
            <el-form-item label="店名备注" prop="cid">
              <el-input v-model="formData.shopRemark" :clearable="true" placeholder="请输入店名备注"/>
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
    <el-dialog v-model="dialogUpdFormVisible" :before-close="closeDialog" :draggable="true" :title="typeTitle"
               destroy-on-close draggable overflow
               width="80%">
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
              <el-form-item label="商铺ID" prop="productId">
                <el-input disabled v-model="formData.productId" :clearable="true" placeholder="请输入产品ID"/>
              </el-form-item>
            </el-col>
            <el-col :span="16">
              <el-form-item label="店名" prop="shopRemark">
                <el-input v-model="formData.shopRemark" disabled/>
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
                    <el-input :rows="2" type="textarea" v-if="editUpdFieldVisible(scope.row, scope.$index)"
                              v-model="scope.row.address"></el-input>
                    <el-input :rows="2" type="textarea" v-model="scope.row.address" readonly v-else></el-input>
                  </template>
                </el-table-column>
                <el-table-column label="金额（元）" prop="money" width="120px">
                  <template #default="scope">
                    <el-input type="number" v-if="editUpdFieldVisible(scope.row, scope.$index)"
                              v-model.number="scope.row.money"
                              :step="10"></el-input>
                    <el-input type="number" v-else v-model.number="scope.row.money" readonly></el-input>
                    <!--                    <span v-else>{{ scope.row.money }}</span>-->
                  </template>
                </el-table-column>
                <el-table-column label="设备" prop="device" width="120px">
                  <template #default="scope">
                    <el-select v-if="editUpdFieldVisible(scope.row, scope.$index)" v-model="scope.row.device">
                      <el-option label="默认" value="default"/>
                      <el-option label="安卓" value="Android"/>
                      <el-option label="苹果" value="iOS"/>
                    </el-select>
                    <el-input v-model="scope.row.device" readonly v-else></el-input>
                  </template>
                </el-table-column>
                <el-table-column label="标识" prop="markId" width="120px">
                  <template #default="scope">
                    <el-input v-if="editUpdFieldVisible(scope.row, scope.$index)" v-model="scope.row.markId"></el-input>
                    <el-input v-else v-model="scope.row.markId" readonly></el-input>
                  </template>
                </el-table-column>
                <el-table-column label="开关" prop="status" width="100px">
                  <template #default="scope">
                    <el-switch v-if="activeUpdIndex === scope.$index" v-model="scope.row.status" :active-value="1"
                               :inactive-value="0" active-text="开启"
                               inactive-text="关闭" inline-prompt size="large" width="70px"
                               @change="()=>{switchEnable(scope.row)}"></el-switch>
                    <el-switch v-else v-model="scope.row.status" :active-value="1" :inactive-value="0"
                               active-text="开启"
                               inactive-text="关闭" inline-prompt size="large" width="70px"
                               @change="()=>{switchEnable(scope.row)}"></el-switch>
                  </template>
                </el-table-column>
                <el-table-column align="right" width="100">
                  <template #header>
                    <el-popconfirm @confirm="handleAdd2Upd" width="320" confirm-button-text="Yes"
                                   cancel-button-text="No, Thanks" :icon="InfoFilled" icon-color="#626AEF"
                                   title="确定要添加商品吗？">
                      <template #reference>
                        <el-button type="danger">
                          <Plus style="width:1em; height:1em;"/>
                        </el-button>
                      </template>
                    </el-popconfirm>
                  </template>
                  <template #default="scope">
                    <div v-if="activeUpdIndex === scope.$index">
                      <el-button type="primary" @click="handleSave2Upd()"><Select style="width:1em; height:1em;"/>
                      </el-button>
                    </div>
                    <div v-else>
<!--                      <el-button v-if="editUpdFieldVisible(scope.row, scope.$index)" type="success" @click="handleEdit2Upd(scope.$index)">
                        <Edit style="width:1em; height:1em;"/>
                      </el-button>-->
                      <el-popconfirm @confirm="handleDelete2Upd(scope.$index)" width="320" confirm-button-text="Yes"
                                     cancel-button-text="No, Thanks" :icon="InfoFilled" icon-color="#626AEF"
                                     title="注意：如果只剩一条记录，删除商品时将连同店铺一起删除。确定要删除该商品吗？">
                        <template #reference>
                          <el-button type="danger">
                            <Delete style="width:1em; height:1em;"/>
                          </el-button>
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

    <el-dialog v-model="detailShow" style="width: 800px" lock-scroll :draggable="true" :before-close="closeDetailShow"
               title="查看详情"
               destroy-on-close>
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

    <!--  dialog -->
<!--    <div class="form-dialog" v-show="dialogPFormVisible">
      <div class="form-wrapper">
        <div class="form-item">
          <label for="uid">uid:</label>
          <input type="text" v-model="formData.uid" id="uid">
        </div>
        <div class="form-item">
          <label for="cid">cid:</label>
          <input type="text" v-model="formData.cid" id="cid">
        </div>
        <div class="form-item">
          <label for="productId">productId:</label>
          <input type="text" v-model="formData.productId" id="productId">
        </div>
        <div class="form-item">
          <label for="shopRemark">shopRemark:</label>
          <input type="text" v-model="formData.shopRemark" id="shopRemark">
        </div>
        <div class="form-item">
          <label>list:</label>
          <button @click="addPFormRow">新增一行</button>
          <button @click="deletePFormSelected">删除选中行</button>
          <div class="table-wrapper">
            <div class="table-scroll">
              <div class="table">
                <div v-for="(item, index) in formData.list" :key="index" class="table-row" :class="{ 'selected': selectedPFormRows.includes(index) }">
                  <input type="checkbox" v-model="selectedPFormRows" :value="index">
                  <input type="text" v-model="item.address">
                  <input type="number" v-model="item.money">
                  <input type="number" v-model="item.status">
                  <input type="checkbox" v-model="item.enable">
                  <button @click="deletePFormRow(index)">删除</button>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="form-actions">
          <button @click="dialogPFormVisible = false">取消</button>
          <button @click="handlePFormSubmit">确定</button>
        </div>
      </div>
    </div>-->

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
import {getDictFunc, formatDate, formatBoolean, filterDict, ReturnArrImg, onDownloadFile} from '@/utils/format'
import {ElMessage, ElMessageBox} from 'element-plus'
import {ref, reactive, nextTick} from 'vue'
import {CircleCheck, CircleClose, Delete, Edit, InfoFilled, Plus, Select} from '@element-plus/icons-vue';
import {setUserInfo} from "@/api/user";
import bgImage from "@/assets/od_info_bg.png";
import shopBgImage from '@/assets/shop_bg.png'

defineOptions({
  name: 'ChannelShop'
})
const backgroundImageStyle = `background-image: url(${shopBgImage});background-size: 100% 100%;background-repeat: no-repeat;`;

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
  for (let ele of formData.value.list) {
    if (ele.address === '') {
      ElMessage({
        type: 'error',
        message: '请先正确填写上一条记录中的地址'
      })
      return
    }
    if (ele.money === 0) {
      ElMessage({
        type: 'error',
        message: '请先正确填写上一条记录中的金额'
      })
      return
    }
  }
  let item = {
    address: '',
    money: 0,
    status: 0
  };
  formData.value.list.push(item);
  activeUpdIndex.value = formData.value.list.length - 1;
};


const editUpdFieldVisible = (row, index) => {
  console.log(row)
  if (row.id) {
    // formData.value.list[index].enable = false
    return false
  } else {
    // formData.value.list[index].enable = true
    return true
  }
  // let flag = activeUpdIndex.value === index
  // console.log(flag)
}

// 编辑行
const handleEdit2Upd = (index) => {
  activeUpdIndex.value = index;
  let ele = formData.value.list[index];
  console.log("handleEdit2Upd ele", ele)
  // let money = Number(ele.money);
  let id = ele.id;
  if (id) {
    formData.value.list[index].enable = false
    // editUpdMoneyVisible.value = false
  } else {
    formData.value.list[index].enable = true
  }
};
// 保存行
const handleSave2Upd = () => {
  let create = {...formData.value}
  let newList = []
  newList.push(formData.value.list[activeUpdIndex.value])
  create.list = newList
  console.log(create.list)
  if (create.list.length === 0) {
    ElMessage({
      type: 'error',
      message: '请至少添加一个商铺地址'
    })
    return
  } else {
    for (let i = 0; i < create.list.length; i++) {
      let addr = create.list[i].address
      let money = create.list[i].money
      console.log(addr)
      console.log(money)
      if (addr === '') {
        ElMessage({
          type: 'error',
          message: '请填写正确地址'
        })
        return
      } else if (money === 0) {
        ElMessage({
          type: 'error',
          message: '请填写正确金额'
        })
        return
      }
    }
  }
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
  } else {
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
  if (!list) list = []
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
const switchEnable = async (row) => {
  console.log(row)
  // userInfo.value = JSON.parse(JSON.stringify(row))
  await nextTick()
  const req = {
    type: 2,
    ...row
  }
  console.log(req)
  if (req.id) {
    const res = await updateChannelShop(req)
    if (res.code === 0) {
      ElMessage({type: 'success', message: `${req.status === 0 ? '禁用' : '启用'}单条商品成功`})
      await getTableData()
    }
  } else {
    console.log("没id的数据，不操作")
  }
}

// 开关所有
const switchEnableAll = async (row, status) => {
  console.log(row)
  // userInfo.value = JSON.parse(JSON.stringify(row))
  await nextTick()
  const req = {
    type: 3,
    status: status,
    ...row
  }
  console.log(req)
  if (req.productId) {
    const res = await updateChannelShop(req)
    if (res.code === 0) {
      ElMessage({type: 'success', message: `${req.status === 0 ? '禁用' : '启用'}店铺成功`})
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
      device: '',
      markId: '',
      money: 0,
      status: 0,
      enable: true,
    }
  ]
})


// 验证规则
const rule = reactive({
  cid: [
    {
      required: true,
      message: '',
      trigger: ['input', 'blur'],
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
  elSearchFormRef.value?.validate(async (valid) => {
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
const getTableData = async () => {
  const table = await getChannelShopList({page: page.value, pageSize: pageSize.value, ...searchInfo.value})
  const vcpTable = await getChannelProductSelf({page: 1, pageSize: 999, type: 1})

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
const setOptions = async () => {
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
const updateChannelShopFunc = async (row) => {
  const res = await findChannelShopByProductID({productId: row.productId})
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.rechannelShop
    dialogUpdFormVisible.value = true
  }
}


// 删除行
const deleteChannelShopFunc = async (row) => {
  console.log(row)
  const res = await deleteChannelShop({ID: row.ID})
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
  const res = await findChannelShop({ID: row.ID})
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
  // dialogPFormVisible.value = true
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
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    switch (type.value) {
      case 'create':
        if (formData.value.list.length === 0) {
          ElMessage({
            type: 'error',
            message: '请至少添加一个商铺地址'
          })
          return
        } else {
          for (let i = 0; i < formData.value.list.length; i++) {
            let addr = formData.value.list[i].address
            let money = formData.value.list[i].money
            console.log(addr)
            console.log(money)
            if (addr === '') {
              ElMessage({
                type: 'error',
                message: '请填写地址'
              })
              return
            } else if (money === 0) {
              ElMessage({
                type: 'error',
                message: '请填写正确金额'
              })
              return
            }
          }
        }
        res = await createChannelShop(formData.value)
        break
      case 'update':
        if (formData.value.list.length === 0) {
          ElMessage({
            type: 'error',
            message: '请至少添加一个商铺地址'
          })
          return
        } else {
          for (let i = 0; i < formData.value.list.length; i++) {
            let addr = formData.value.list[i].address
            let money = formData.value.list[i].money
            console.log(addr)
            console.log(money)
            if (addr === '') {
              ElMessage({
                type: 'error',
                message: '请填写地址'
              })
              return
            } else if (money === 0) {
              ElMessage({
                type: 'error',
                message: '请填写正确金额'
              })
              return
            }
          }
        }
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

//
const addPFormRow = () => {
  formData.value.list.push({
    address: '',
    money: 0,
    status: 0,
    enable: true,
  });
};
const dialogPFormVisible = ref(false);
const handlePFormSubmit = () => {
  // 提交表单逻辑，可以在这里处理表单数据
  console.log(formData.value);
  dialogPFormVisible.value = false;
};
const selectedPFormRows = ref([]);

const deletePFormRow = (index) => {
  formData.value.list.splice(index, 1);
};

const deletePFormSelected = () => {
  const sortedSelection = selectedPFormRows.value.sort((a, b) => b - a);
  sortedSelection.forEach(index => {
    formData.value.list.splice(index, 1);
  });
};
</script>

<style>
.card_list_wrap {
  display: flex;
  gap: 32px 24px;
  flex-flow: wrap;
}

.card_wrap.sm_card {
  min-width: 252px;
  width: 252px;
  position: relative;
}

.card_warp.hot::before {
  background: #FFD4C8;
}

.card_warp::before {
  background: #fe5f47;
}

.card_warp::before {
  content: '';
  position: absolute;
  left: 24px;
  width: 60px;
  height: 4px;
  background: #3860f4;
  box-shadow: 0 16px 16px 0 rgba(55, 69, 103, 2%), 0 8px 8px 0 rgba(235, 240, 252, 2%);
  border-radius: 0 0 2px 2px;
  z-index: 1;
}

.card_wrap {
  display: flex;
  flex: 1;
  flex-direction: column;
  position: relative;
  min-width: 344px;
  box-sizing: border-box;
  background-color: #fff;
  box-shadow: 0px 16px 16px 0px rgb(55 69 103 / 2%), 0px 8px 8px 0px rgb(235 240 252 / 2%);
  border: 1px solid #e1e6f0;
}

.card_warp.hot .card_header .card_header_title {
  color: #fff;
}

.card_warp.hot .card_header .card_header_description {
  color: #1d1585;
}

.card_header {
  flex: 1 1 auto;
  min-height: 0;
  padding-top: 32px;
  padding-bottom: 18px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  position: relative;
}

.card_header_tag.new_customers {
  color: #ffffff;
  background-image: linear-gradient(130deg, rgb(131, 128, 208) 40%, rgb(94, 50, 255) 100%);
}

.card_header_tag.new_customers::before {
  background: #ffffff;
}

.card_header_tag.new_customers::after {
  background: rgb(131, 128, 208);
}

.card_header_tag::after {
  background: rgba(255, 235, 214, 1);
}

.card_header_tag::before {
  content: '';
  position: absolute;
  left: -4px;
  top: 12px;
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background-color: #3860f4;
  z-index: 10;
}

.card_header_tag::after {
  content: " ";
  height: 28px;
  width: 28px;
  background: linear-gradient(90deg, #f9faff 0%, #ebf0fc 100%);
  position: absolute;
  top: 0px;
  left: -9px;
  border-radius: 4px;
  transform: rotate(-41deg) skew(16deg, 9deg) translateZ(-1px);
  z-index: 0;
}

.card_header_tag {
  background: rgba(255, 235, 214, 1);
  color: #724040;
}

.card_header_tag {
  position: absolute;
  top: 12px;
  right: 0;
  padding: 2px 10px;
  border-radius: 4px 0 0 4px;
  font-size: 12px;
  line-height: 24px;
  font-weight: bold;
  background: linear-gradient(90deg, #f9faff 0%, #ebf0fc 100%);
  color: #3860f4;
  transform-style: preserve-3d;
}

.card_header_title {
  margin-top: 8px;
  padding-left: 24px;
  padding-right: 24px;
  font-size: 18px;
  font-weight: 700;
  line-height: 36px;
  color: rgb(131, 128, 208);
}

.card_header_description {
  padding-left: 24px;
  padding-right: 24px;
  margin-bottom: -20px;
  font-size: 14px;
  line-height: 28px;
  color: rgb(121, 118, 206);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 20px;
  color: #6B7687;
}

.card_content {
  padding: 16px 24px;
}

.label .label_item{
  color: #0a0e31;
  margin-bottom: -20px
}

.card_footer::before {
  content: '';
  position: absolute;
  top: 0;
  left: 50%;
  transform: translate(-50%, 0);
  width: calc(100% - 48px);
  height: 1px;
  background-color: #ebf0fc;
}

.card_footer {
  position: relative;
  padding: 16px 0 0 0;
}

.card_footer .label {
  display: flex;
  padding: 0 24px;
}

.card_footer .label .label_item {
  background-color: rgb(41, 21, 107);
  color: rgba(254, 95, 71, 1);
}

.card_footer .label .label_item {
  margin-right: 8px;
  min-width: max-content;
  padding: 0 8px;
  height: 24px;
  line-height: 24px;
  font-size: 12px;
  font-weight: bold;
  color: #3860f4;
  background: #ebf0fc;
  border-radius: 2px;
}

.card_footer .label .label_item_blue {
  margin-right: 8px;
  min-width: max-content;
  padding: 0 8px;
  height: 24px;
  line-height: 24px;
  font-size: 12px;
  font-weight: bold;
  color: #3848f4;
  background: #ebf0fc;
  border-radius: 2px;
}

.card_footer .label .label_item_err {
  margin-right: 8px;
  min-width: max-content;
  padding: 0 8px;
  height: 24px;
  line-height: 24px;
  font-size: 12px;
  font-weight: bold;
  color: #f43838;
  background: #ebf0fc;
  border-radius: 2px;
}

.sm_card .card_footer .list_wrap {
  flex-direction: column;
  padding: 0;
}

.card_footer .list_wrap {
  margin-top: 24px;
  padding: 0 24px 36px 24px;
  display: flex;
  justify-content: space-between;
}

.sm_card .card_footer .list_wrap .percent {
  padding: 0 24px;
}

.card_footer .list_wrap .percent {
  display: flex;
  gap: 6px;
  align-items: center;
  min-width: 100px;
}

.sm_card .card_footer .list_wrap .percent .font-num {
  font-size: 48px;
}

.card_footer .list_wrap .percent .font-num {
  font-size: 50px;
  padding-right: 8px;
  color: #8294ee;
}

.card_footer .list_wrap .percent .percent_des {
  display: flex;
  flex-direction: column;
  font-size: 12px;
  line-height: 16px;
  color: #374567;
}

.card_footer .list_wrap .percent .percent_des > span:last-child {
  color: #7a8ba6;
}

.card_footer .org_percent {
  margin-top: 4px;
  line-height: 12px;
}

.sm_card .card_footer .list_wrap .btn_wrap {
  margin-top: 12px;
}

.card_footer .list_wrap .btn_wrap {
  display: flex;
  gap: 8px;
  flex: 1 1 auto;
  justify-content: flex-end;
}

.card_footer .buy .button-base.yellow {
  color: #ffffff;
  border: none;
  height: 34px;
  line-height: 34px;
}

.card_footer .buy .button-base.red {
  color: #ffffff;
  border: none;
  height: 34px;
  line-height: 34px;
}

.card_footer .buy .button-base.blue {
  color: #ffffff;
  border: none;
  height: 34px;
  line-height: 34px;
}

.sm_card .card_footer .list_wrap .button-base {
  width: 100%;
}

.card.button-base.yellow {
  background: linear-gradient(to right, #8f87ff, #769dff);
  color: #4e342e;
}

.card.button-base.red {
  background: linear-gradient(to right, #2a54e1, #49a4ef);
  color: #4e342e;
}

.card.button-base.blue {
  background: linear-gradient(to right, #0c87ec, #0daef3);
  color: #4e342e;
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

.form-dialog {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background-color: #f5f5f5;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  z-index: 1000;
}

.table-wrapper {
  max-height: 300px;
  overflow-y: auto;
}

.table-scroll {
  width: calc(100% + 17px); /* 17px 是滚动条的宽度 */
  overflow-x: hidden;
}

.form-item {
  margin-bottom: 15px;
}

.table {
  margin-top: 10px;
  border: 1px solid #ccc;
  border-collapse: collapse;
  width: 100%;
}
.table-row.selected {
  background-color: #f0f0f0;
}
.table-row {
  display: flex;
  align-items: center;
  padding: 8px;
  border-bottom: 1px solid #eee;
}

.table-row:last-child {
  border-bottom: none;
}

.table-row.selected {
  background-color: #e1f5fe;
}

.table-row input,
.table-row button {
  margin-right: 10px;
  flex: 1;
  height: 30px;
  padding: 5px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.form-actions {
  margin-top: 20px;
  text-align: right;
}

.form-actions button {
  padding: 8px 16px;
  margin-left: 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  background-color: #4CAF50;
  color: white;
}

.form-actions button:hover {
  background-color: #45a049;
}

</style>

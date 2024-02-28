<template>
  <div class="organization">
    <div class="gva-search-box org-top">
      成员管理
    </div>
    <div>
      <el-row :gutter="12">
<!--        <el-col :span="8">-->

<!--        </el-col>-->
<!--        <el-col :span="16">-->
        <el-col :span="24">
          <div class="">
            <div class="toolbar">
              <div class="gva-search-box" style=" margin-left: 10px;margin-right: 10px">
                <el-form :inline="true">
                  <el-form-item label="商户ID" prop="cid">
                    <el-input v-model="userSearch.username" placeholder="请输入要搜索的用户名"/>
                  </el-form-item>
                  <el-form-item>
                    <el-button type="primary" @click="getUserTable">搜索</el-button>
                  </el-form-item>
                  <el-form-item>
                    <el-button type="primary" @click="addUser">开户</el-button>
                  </el-form-item>
                </el-form>
              </div>
              <div>
              </div>
            </div>
            <div>
              <div style=" margin-left: 10px;">
                <el-table :data="treeData" style="margin-right: 10px" border>
                  <el-table-column type="expand">
                    <template #default="props">
                      <!--                    <div v-if="props.row.children"></div>
                                          <div v-else>-->
                      <div>
                        <el-row :gutter="2">
                          <el-col :span="6">
                          </el-col>
                          <el-col :span="6">
                            <el-tag effect="dark">团队归属: {{ props.row.organization.name }}</el-tag>
                          </el-col>
                          <el-col :span="6">
                            <el-tag effect="dark">团队长: {{ props.row.sysUser.username }}</el-tag>
                          </el-col>
                          <el-col :span="6">
                          </el-col>
                          <el-col :span="24">
                            <h3 style="margin-left: 40%; margin-bottom: -40px">团队成员</h3>
                          </el-col>
                          <el-col :span="24">
                            <div>
                              <el-table :data="props.row.children" style="width: 100%;margin: 60px">
                                <el-table-column label="团队成员" prop="sysUser.username" width="190"/>
                                <el-table-column align="center" prop="x9" label="积分" width="90"/>
                                <el-table-column label="操作列" fixed="right">
                                  <template #default="{row}">
                                    <el-button v-auth="btnAuth.rechargeBtn" link type="primary" icon="wallet"
                                               @click="showOperateRecharge(row)"> 充值
                                    </el-button>
                                    <el-button type="primary" link icon="magic-stick" @click="resetPasswordFunc(row)">
                                      重置密码
                                    </el-button>
                                    <el-button type="primary" link icon="lock" @click="getAuthCaptcha(row)"> 动态验证
                                    </el-button>
                                    <el-popover v-model="row.visible" placement="top" width="160">
                                      <p>确定要删除此用户吗</p>
                                      <div style="text-align: right; margin-top: 8px;">
                                        <el-button type="primary" link @click="row.visible = false">取消</el-button>
                                        <el-button type="primary" @click="deleteUserFunc(row)">确定</el-button>
                                      </div>
                                      <template #reference>
                                        <el-button type="danger" link icon="delete">删除</el-button>
                                      </template>
                                    </el-popover>
                                  </template>
                                </el-table-column>
                              </el-table>
                            </div>
                          </el-col>
                        </el-row>

                      </div>
                    </template>
                  </el-table-column>
                  <el-table-column prop="organization.name" label="团队" width="160"></el-table-column>
                  <el-table-column prop="sysUser.username" label="负责人" width="160"></el-table-column>
                  <el-table-column align="center" prop="x9" label="积分" width="90"/>
                  <el-table-column label="操作列">
                    <template #default="{row}">
                      <el-button v-auth="btnAuth.rechargeBtn" link type="primary" icon="wallet"
                                 @click="showOperateRecharge(row)"> 充值
                      </el-button>
                      <el-button type="primary" link icon="magic-stick" @click="resetPasswordFunc(row)"> 重置密码
                      </el-button>
                      <el-button type="primary" link icon="lock" @click="getAuthCaptcha(row)"> 动态验证</el-button>
                      <el-popover v-model="row.visible" placement="top" width="160">
                        <p>确定要删除此用户吗</p>
                        <div style="text-align: right; margin-top: 8px;">
                          <el-button type="primary" link @click="row.visible = false">取消</el-button>
                          <el-button type="primary" @click="deleteUserFunc(row)">确定</el-button>
                        </div>
                        <template #reference>
                          <el-button type="danger" link icon="delete">删除</el-button>
                        </template>
                      </el-popover>
                    </template>
                  </el-table-column>

                </el-table>
              </div>
            </div>
          </div>
        </el-col>
        <el-col :span="12">
        </el-col>
      </el-row>
    </div>

    <!--  开户  -->
    <el-dialog v-model="addUserDialog" :draggable="true" title="开户" width="600px">
      <el-form ref="elFormRef" :model="userForm" label-width="80px" :rules="rules">
        <el-row>
          <el-col :span="12">
            <el-form-item label="用户名" label-width="80px">
              <el-input v-model="userForm.username" prop="username"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="团队名" label-width="80px">
              <el-input v-model="userForm.orgName" prop="orgName"/>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item :minlength="6" label="新密码" prop="newPassword" label-width="80px">
              <el-input v-model="userForm.newPassword" show-password/>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item :minlength="6" label="确认密码" prop="confirmPassword" label-width="80px">
              <el-input v-model="userForm.confirmPassword" show-password/>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="充值积分" label-width="80px">
              <el-input v-model="userForm.recharge" prop="recharge"/>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="产品分配" style="width: 100%">
              <el-cascader
                  :options="channelCodeOptions"
                  :props="channelCodeProps"
                  @change="handleChange"
                  style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="addUserClear">取 消</el-button>
        <el-button type="primary" @click="addUserEnter">确 定</el-button>
      </template>
    </el-dialog>

    <!-- 防爆验证码 -->
    <el-dialog v-model="showAuthCaptcha" title="重置谷歌动态验证" :draggable="true" width="360px"
               @close="clearAuthCaptcha">
      <el-form ref="modifyCapForm" :model="capModify" label-width="80px">
        <el-form-item label="用户ID" prop="toUid">
          <el-input v-model="capModify.ID" disabled/>
        </el-form-item>
        <el-form-item label="登录密码" prop="password">
          <el-input v-model="capModify.password" show-password/>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showAuthCaptcha = false">取 消</el-button>
          <el-button type="primary" @click="resetAuthCaptcha">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 获取订单acc统计数据 -->
    <el-dialog v-model="showCostOrderAccVisible" :title="showCostOrderAccTitle" :draggable="true" width="1000px"
               @close="closeCostOrderAcc">
      <!--      <div class="gva-search-box">-->
      <!--        <el-form :inline="true" :model="searchAccInfo" class="demo-form-inline" @keyup.enter="onAccSubmit">-->
      <!--          <el-form-item label="通道账户名" prop="acId">-->
      <!--            <el-input v-model="searchAccInfo.acAccount" placeholder="搜索通道账户"/>-->
      <!--          </el-form-item>-->
      <!--          <el-form-item label="通道账户ID" prop="acAccount">-->
      <!--            <el-input v-model="searchAccInfo.acId" placeholder="搜索通道账户ID"/>-->
      <!--          </el-form-item>-->
      <!--          <el-form-item label="通道ID" prop="cid">-->
      <!--            <el-input v-model.number="searchAccInfo.channelCode" placeholder="搜索通道ID"/>-->
      <!--          </el-form-item>-->
      <!--          <el-form-item>-->
      <!--            <el-button type="primary" icon="search" @click="onAccSubmit">查询</el-button>-->
      <!--            <el-button icon="refresh" @click="onAccReset">重置</el-button>-->
      <!--          </el-form-item>-->
      <!--        </el-form>-->
      <!--      </div>-->
      <!--      <div class="gva-table-box">
              <el-scrollbar>
                <el-table ref="multipleTable" tooltip-effect="dark" :data="costOrderAccTable" border resizable="true" height="400"
                          show-summary>
                  <el-table-column align="center" label="通道ID" width="80">
                    <template #default="{row}">
                      {{ String(row.channelCode) }}
                    </template>
                  </el-table-column>
                  <el-table-column align="center" label="账号ID" width="90">
                    <template #default="{row}">
                      {{ String(row.acId) }}
                    </template>
                  </el-table-column>
                  <el-table-column align="center" label="通道账号" width="160">
                    <template #default="{row}">
                      {{ String(row.acAccount) }}
                    </template>
                  </el-table-column>
                  <el-table-column align="center" sortable label="3日前" prop="x1" width="120"/>
                  <el-table-column align="center" sortable label="2日前" prop="x2" width="120"/>
                  <el-table-column align="center" sortable label="昨日" prop="x3" width="120"/>
                  <el-table-column align="center" sortable label="今日" prop="x4" width="120"/>
                  <el-table-column align="center" sortable label="总充值" prop="x0" width="120"/>
                </el-table>
              </el-scrollbar>
            </div>-->
    </el-dialog>

    <!-- 积分结算 -->
    <el-dialog v-model="showCostRechargeVisible" title="积分结算" :draggable="true" width="560px"
               @close="closeCostRecharge">
      <el-form :model="costRechargeForm" label-width="80px">
        <!--        <el-row :gutter="8">-->
        <!--          <el-col :span="12">-->
        <!--          </el-col>-->
        <!--          <el-col :span="12">-->
        <!--          </el-col>-->

        <!--        </el-row>-->
        <el-row :gutter="8">
          <el-col :span="12">
            <el-form-item label="用户ID" prop="x0">
              <el-input v-model="costRechargeForm.x0" readonly/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="costRechargeForm.username" readonly/>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="总余额" prop="x9">
              <el-input v-model="costRechargeForm.x9" readonly/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="8">
          <el-col :span="8">
            <el-form-item label="3日前消费" prop="x1">
              <el-input v-model="costRechargeForm.x1" readonly/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="2日前消费" prop="x2">
              <el-input v-model="costRechargeForm.x2" readonly/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="昨日消费" prop="x3">
              <el-input v-model="costRechargeForm.x3" readonly/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="8">
          <el-col :span="24">
            <el-form-item label="今日消费" prop="x4">
              <el-input v-model="costRechargeForm.x4" readonly/>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="8">
          <el-col :span="12">
            <el-form-item label="（昨日）" prop="x5">
              转账: {{ costRechargeForm.x6 }} / 充值: {{ costRechargeForm.x5 }}
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="（今日）" prop="x6">
              转账: {{ costRechargeForm.x8 }} / 充值: {{ costRechargeForm.x7 }}
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-dialog>

    <!-- 查看 -->
    <el-dialog v-model="showQRCode" title="谷歌动态验证" :draggable="true" width="300px" @close="closeAuthCaptcha">
      <div class="qrcode-generator">
        <div v-if="isNotSetting" style="margin-bottom: 20px">
          暂未设置谷歌动态验证，请尽快设置！
        </div>
        <div v-else>
          <img :src="qrcodeUrl" alt="QR Code" style="height: 200px"/>
        </div>
        <el-button link type="primary" icon="lock" @click="resetShowAuthCaptcha"> 设置(或重置)</el-button>
      </div>
    </el-dialog>

    <!-- 积分划转 -->
    <el-dialog v-model="showRechargeVisible" title="积分划转" :draggable="true" width="360px" @close="clearRecharge">
      <el-form :model="rechargeForm" label-width="80px">
        <el-form-item label="用户ID" prop="toUid">
          <el-input v-model="rechargeForm.toUid" disabled/>
        </el-form-item>
        <el-form-item label="用户名" prop="toUsername">
          <el-input v-model="rechargeForm.toUsername" disabled/>
        </el-form-item>
        <el-form-item label="划转积分" prop="recharge">
          <el-input v-model.number="rechargeForm.recharge" :min="0"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showRechargeVisible = false">取 消</el-button>
          <el-button type="primary" @click="transferRecharge(2)">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 直充 -->
    <el-dialog v-model="operateRechargeVisible" title="积分充值" :draggable="true" width="360px"
               @close="clearOperateRecharge">
      <el-form :model="rechargeForm" label-width="80px">
        <el-form-item label="用户ID" prop="toUid">
          <el-input v-model="rechargeForm.toUid" disabled/>
        </el-form-item>
        <el-form-item label="用户名" prop="toUsername">
          <el-input v-model="rechargeForm.toUsername" disabled/>
        </el-form-item>
        <el-form-item label="充值积分" prop="recharge">
          <el-input v-model.number="rechargeForm.recharge"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="operateRechargeVisible = false">取 消</el-button>
          <el-button type="primary" @click="transferRecharge(1)">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import QRCode from 'qrcode';
import {
  createOrganization,
  deleteOrganization,
  getOrganizationList,
  updateOrganization,
  findOrganization,
  findOrgUserAll,
  createOrgUser,
  findOrgUserList,
  setOrgUserAdmin,
  deleteOrgUserApi,
  transferOrgUserApi
} from '@/plugin/organization/api/organization'
import {getUserList} from '@/api/user.js'
import {ref} from 'vue'
import {ElMessageBox, ElMessage} from 'element-plus'

export default {
  name: 'Organization',
}
</script>

<script setup>
import {findOrgUserListSelf} from '@/plugin/organization/api/organization';
import {deleteUser, selfRegister, resetCaptcha, resetPassword, openAccRegister} from '@/api/user';
import {ElMessage, ElMessageBox} from "element-plus";
import {getUserWalletCostOV, getUserWalletOverview, getUserWalletSelf, transferUserWallet} from "@/api/userWallet";
import {useBtnAuth} from '@/utils/btnAuth'
import {reactive, ref} from "vue";
import CenterCard from "@/view/vbox/dashboard/dataCenterComponents/centerCard.vue";
import {getOrderAccOverview} from "@/api/payOrder";
import {getChannelProductSelf} from "@/api/channelProduct";

const walletCustomStyle = ref({
  background: 'linear-gradient(to right, #22111a, #606266)',
  color: '#FFF',
  height: '140px',
})

const userBalance = ref(0)

const btnAuth = useBtnAuth()
const data = ref([])

const selectData = ref([])

// 右侧人员操作

// 多选人员
const handleSelectionChange = (val) => {
  selectData.value = val
}

// 人员新增
const addUserDialog = ref(false)

// 人员操作弹窗数据
const userForm = ref({
  username: '',
  confirmPassword: '',
  newPassword: '',
  orgName: '',
  codeProdIDS: [],
  recharge: 0,
  enable: 1,
  enableAuth: 1,
})

const rules = reactive({
  username: [
    {min: 6, message: '最少6个字符', trigger: ['input', 'blur']},
  ],
  newPassword: [
    {min: 6, message: '最少6个字符', trigger: ['input', 'blur']},
  ],
  confirmPassword: [
    {min: 6, message: '最少6个字符', trigger: ['input', 'blur']},
    {
      validator: (rule, value, callback) => {
        // console.log('value', value)
        if (value !== userForm.value.newPassword) {
          callback(new Error('两次密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
})

// 增加新用户
const addUser = async () => {
  addUserDialog.value = true
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
const getUserTable = async () => {
  const vcpTable = await getChannelProductSelf({page: 1, pageSize: 999})

  const res = await findOrgUserListSelf({page: 1, pageSize: 9999, ...userSearch.value})
  // let d = generateTreeData(res.data.list);
  // console.log(...d)
  const walletRes = await getUserWalletOverview({...userSearch.value})
  // console.log("walletRes",walletRes.data)
  let resultMap;
  if (walletRes.code === 0) {
    resultMap = reactive(new Map(walletRes.data.map(item => [item.x0, item])));
  } else {
    resultMap = reactive(new Map());
  }
  if (res.code === 0) {
    page.value = res.data.page
    pageSize.value = res.data.pageSize
    total.value = res.data.total
    userTable.value = res.data.list
    vcpTableData.value = vcpTable.data.list
    setOptions()

    for (let i = 0; i < userTable.value.length; i++) {
      let userID = userTable.value[i].sysUserID;
      const resultForX0 = resultMap.has(userID) ? resultMap.get(userID) : null;
      if (resultForX0) {
        userTable.value[i].x0 = resultForX0.x0
        userTable.value[i].x1 = resultForX0.x1
        userTable.value[i].x2 = resultForX0.x2
        userTable.value[i].x3 = resultForX0.x3
        userTable.value[i].x4 = resultForX0.x4
        userTable.value[i].x5 = resultForX0.x5
        userTable.value[i].x6 = resultForX0.x6
        userTable.value[i].x7 = resultForX0.x7
        userTable.value[i].x8 = resultForX0.x8
        userTable.value[i].x9 = resultForX0.x9
      } else {
        userTable.value[i].x0 = userID
        userTable.value[i].x1 = 0
        userTable.value[i].x2 = 0
        userTable.value[i].x3 = 0
        userTable.value[i].x4 = 0
        userTable.value[i].x5 = 0
        userTable.value[i].x6 = 0
        userTable.value[i].x7 = 0
        userTable.value[i].x8 = 0
        userTable.value[i].x9 = 0
      }
    }

    treeData.value = generateTreeData(userTable.value);
    // console.log(treeData.value)
    // console.log(userTable.value)
  }
}

// 组织用户搜索
const userSearch = ref({
  username: '',
})

const deleteUserFunc = async (row) => {
  ElMessageBox.confirm('删除用户为高危操作，请核实是否删除？', '提示', {
    type: 'warning',
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    showClose: false,
  }).then(async () => {
    const res = await deleteUser({id: row.sysUser.ID})
    if (res.code === 0) {
      ElMessage.success('删除成功')
      await getUserTable()
      selectData.value = []
    }
  }).catch(() => {
    ElMessage.info('取消删除')
  })
}

// 成员入职功能关闭弹窗
const addUserClear = () => {
  addUserDialog.value = false
}

const elFormRef = ref(null)
// 添加组织成员
const addUserEnter = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (valid) {
      userForm.value.enableAuth = Number(userForm.value.enableAuth)
      userForm.value.enable = Number(userForm.value.enable)
      userForm.value.recharge = Number(userForm.value.recharge)
      let res = await openAccRegister(userForm.value)
      if (res.code === 0) {
        ElMessage.success('添加成功')
        await getUserTable()
      }
      userForm.value = {
        username: '',
        confirmPassword: '',
        newPassword: '',
      }
      addUserDialog.value = false
    } else {
      return false
    }
  })
}

// 初始化方法
const init = async () => {
  await getUserTable()
}

init()

// 重置密码
const resetPasswordFunc = (row) => {
  ElMessageBox.confirm(
      '是否将此用户密码重置为123456?',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  ).then(async () => {
    const res = await resetPassword({
      ID: row.sysUser.ID,
    })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: res.msg,
      })
    } else {
      ElMessage({
        type: 'error',
        message: res.msg,
      })
    }
  })
}


// ---------- 重置防爆码 ----------
const modifyCapForm = ref(null)
const showAuthCaptcha = ref(false)
const capModify = ref({})
const resetShowAuthCaptcha = async () => {
  showAuthCaptcha.value = true
}
const resetAuthCaptcha = () => {
  modifyCapForm.value.validate((valid) => {
    if (valid) {
      resetCaptcha({
        password: capModify.value.password,
        toUid: capModify.value.ID,
        type: 1,
      }).then((res) => {
        if (res.code === 0) {
          ElMessage.success('重置谷歌动态验证成功！')
          // 刷新页面
          window.location.reload()
        }
        showQRCode.value = false
        showAuthCaptcha.value = false
      })
    } else {
      return false
    }
  })
}
const clearAuthCaptcha = () => {
  capModify.value = {
    password: '',
  }
  modifyCapForm.value.clearValidate()
}

const closeAuthCaptcha = async () => {
  showQRCode.value = false
}

const url = ref('');
const qrcodeUrl = ref('');
const showQRCode = ref(false);
const isNotSetting = ref(false);

// 查看防爆码
const getAuthCaptcha = (row) => {
  let authCaptcha = row.sysUser.authCaptcha;
  capModify.value = JSON.parse(JSON.stringify(row.sysUser))
  if (authCaptcha !== "") {
    QRCode.toDataURL(authCaptcha)
        .then((dataUrl) => {
          console.log(dataUrl)
          qrcodeUrl.value = dataUrl;
          isNotSetting.value = false;
          showQRCode.value = true;
        })
        .catch((error) => {
          console.error('Failed to generate QR code:', error);
        });
  } else {
    isNotSetting.value = true;
    showQRCode.value = true;
  }
};
// ---------- 重置防爆码 end ----------

// ---------- 消费历史 ----------
// const searchAccInfo = ref({})
const searchAccInfo = ref({
  toUid: '',
  username: '',
  acId: '',
  acAccount: '',
  channelCode: '',
})
// 重置
const onAccReset = () => {
  searchAccInfo.value = {}
  getAccTableData()
}

const getAccTableData = async () => {
  console.log(searchAccInfo.value)
  const voRes = await getOrderAccOverview({...searchAccInfo.value});
  console.log(voRes.data)
  if (voRes.code === 0) {
    showCostOrderAccTitle.value = `订单核算(用户归属:${searchAccInfo.value.username})`
    costOrderAccTable.value = voRes.data.list;
  }
}

// 搜索
const onAccSubmit = async () => {
  getAccTableData()
}

const showCostOrderAccVisible = ref(false)
const showCostOrderAccTitle = ref()
let costOrderAccTable = ref([]);
const showCostOrderAcc = async (row) => {
  searchAccInfo.value.toUid = row.sysUser.ID
  searchAccInfo.value.username = row.sysUser.username
  costOrderAccTable.value = [];
  const voRes = await getOrderAccOverview({...searchAccInfo.value});
  console.log(voRes.data)
  if (voRes.code === 0) {
    showCostOrderAccTitle.value = `订单核算(用户归属:${searchAccInfo.value.username})`
    costOrderAccTable.value = voRes.data.list;
    showCostOrderAccVisible.value = true
  }
}
const closeCostOrderAcc = () => {
  showCostOrderAccVisible.value = false
  costOrderAccTable.value = [];
}
// ---------- 消费历史 ----------

// ---------- 消费结算 ----------
const showCostRechargeVisible = ref(false)
let costRechargeForm = ref({});

const showCostRecharge = async (row) => {
  const voRes = await getUserWalletCostOV({toUid: row.sysUser.ID});
  console.log(voRes.data)
  if (voRes.code === 0) {
    costRechargeForm.value = voRes.data[0];
    costRechargeForm.value.username = row.sysUser.username
    console.log(costRechargeForm.value)
    showCostRechargeVisible.value = true
  }
}
const closeCostRecharge = () => {
  showCostRechargeVisible.value = false
  costRechargeForm.value = {};
}
// ---------- 消费结算 ----------

// ---------- 充值划转 ----------
const showRechargeVisible = ref(false)
const rechargeForm = ref({
  toUid: 0,
  toUsername: '',
  recharge: 0,
  type: 2,
})
const showRecharge = async (row) => {
  showRechargeVisible.value = true
  rechargeForm.value.toUid = row.sysUser.ID
  rechargeForm.value.toUsername = row.sysUser.username
  rechargeForm.value.type = 2
  // console.log(rechargeForm.value)
}

// 积分划转、充值
const transferRecharge = async (type) => {
  rechargeForm.value.type = type;
  let req = {...rechargeForm.value}
  // console.log(req)
  ElMessageBox.confirm(
      '是否为此用户（充值）划转积分?',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  ).then(async () => {
    const res = await transferUserWallet(req)
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: res.msg,
      })
    } else {
      ElMessage({
        type: 'error',
        message: res.msg,
      })
    }
  })
  showRechargeVisible.value = false
  operateRechargeVisible.value = false
}

const clearRecharge = () => {
  showRechargeVisible.value = false
  rechargeForm.value = {
    toUid: 0,
    toUsername: '',
    recharge: 0,
    type: 2,
  }
}
// ---------- 充值划转 ----------

// ---------- 充值充值 ----------
const operateRechargeVisible = ref(false)
const showOperateRecharge = async (row) => {
  operateRechargeVisible.value = true
  rechargeForm.value.toUid = row.sysUser.ID
  rechargeForm.value.toUsername = row.sysUser.username
  rechargeForm.value.type = 1
}
const clearOperateRecharge = () => {
  operateRechargeVisible.value = false
  rechargeForm.value = {
    toUid: 0,
    toUsername: '',
    recharge: 0,
    type: 2,
  }
}
// ---------- 充值充值 ----------

// ---------- tree ----------
const generateTreeData = (list) => {
  const map = new Map();
  const result = [];

  list.forEach(item => {
    const {organization, organizationID, sysUserID, isAdmin, sysUser} = item;
    const {ID, authorityId} = sysUser;

    if (!map.has(organizationID)) {
      map.set(organizationID, {
        ...item,
        children: []
      });
      result.push(map.get(organizationID));
    }

    if (authorityId === 1001) {
      // map.get(organizationID).leader = sysUser;
    } else if (authorityId === 1002) {
      map.get(organizationID).children.push(item);
    }
  });

  return result;
};
// const list = [
//   // ... (你提供的数据)
// ];
const treeData = ref();
// ---------- tree ----------

// ---------- prod ----------
const vcpTableData = ref([])

const setOptions = () => {
  channelCodeOptions.value = []
  setChannelCodeOptions(vcpTableData.value, channelCodeOptions.value, false)
  console.log(...channelCodeOptions.value)
}
const setChannelCodeOptions = (ChannelCodeData, optionsData, disabled) => {
  ChannelCodeData &&
  ChannelCodeData.forEach(item => {
    if (item.children && item.children.length) {
      const option = {
        value: item.ID + '',
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
        value: item.ID + '',
        label: item.productName,
      }
      optionsData.push(option)
    }
  })
}
const channelCodeOptions = ref([])

const channelCodeProps = {
  expandTrigger: 'hover',
  checkStrictly: false,
  emitPath: false,
  multiple: true,
}

const rowId = ref()

function handleChange(value) {
  const selectedValues = value;
  const allValues = [];
  let retList = [];
  // 遍历已选中的节点
  for (let i = 0; i < selectedValues.length; i++) {
    const val = selectedValues[i];
    // 找到当前节点的所有父节点并添加到allValues中
    findAllParentValues(channelCodeOptions.value, val, allValues);
  }
  const uniqueValues = [...new Set(allValues)];
  // // 更新选中的节点值
  // channelCodeOptions.value.forEach(item => {
  //   if (uniqueValues.includes(item.value)) {
  //     item.checked = true;
  //   } else {
  //     item.checked = false;
  //   }
  // });
  console.log('selectedValues', selectedValues)
  console.log('uniqueValues', uniqueValues)

  let concat = uniqueValues.concat([...selectedValues]);

  for (let i = 0; i < concat.length; i++) {
    retList.push(Number(concat[i]));
  }
  userForm.value.codeProdIDS = retList;
  console.log('prod IDS', userForm.value.codeProdIDS)
}

function findAllParentValues(options, value, allValues) {
  options.forEach(item => {
    if (item.children && item.children.length > 0) {
      if (item.children.some(child => child.value === value)) {
        allValues.push(item.value);
        findAllParentValues(options, item.value, allValues);
      } else {
        findAllParentValues(item.children, value, allValues);
      }
    }
  });
}

// ---------- prod ----------
// ---------- org ----------

// ---------- org ----------


</script>

<style scoped lang="scss">
.qrcode-generator {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.org-top {
  padding-bottom: 20px;
}

.gva-organization-box {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  height: calc(100vh - 260px);

  &-left {
    box-sizing: border-box;
    padding: 20px;
    width: 260px;
    background: #fff;

    & > .toolbar {
      margin-bottom: 20px;
    }

    .tree-body-tree-node {
      padding-right: 20px;
      width: 100%;
      display: flex;
      justify-content: space-between;
      flex: 1;

      &-label {
        display: inline-block;
        max-width: 100px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }

    & > .tree-body {
      height: calc(100% - 52px);
      overflow: auto;

      &::-webkit-scrollbar {
        width: 2px;
        height: 2px;
      }
    }
  }

  &-right {
    box-sizing: border-box;
    padding: 20px;
    background: #fff;

    & > .toolbar {
      padding-bottom: 20px;
      display: flex;
      justify-content: space-between;

      .toolbar-search {
        display: flex;
        align-items: center;
        margin-bottom: 20px;

        .el-input {
          flex: 1;
        }

        .el-button {
          margin-left: 10px;
        }
      }
    }

    width: calc(100% - 270px);
  }
}

.acc-container {
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

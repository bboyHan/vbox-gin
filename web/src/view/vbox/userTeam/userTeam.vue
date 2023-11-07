<template>
  <div class="organization">
    <div class="gva-search-box org-top">
      成员管理
    </div>
    <div class="gva-organization-box">
      <div class="gva-organization-box-left">
        <div class="toolbar">
        </div>
      </div>
      <div class="gva-organization-box-right">
        <div class="toolbar">
          <div class="toolbar-search">
            <el-input v-model="userSearch.userName" placeholder="请输入要搜索的用户名" />
            <el-button type="primary" @click="getUserTable">搜索</el-button>
            <el-button type="primary" @click="addUser()">新增成员</el-button>
          </div>
          <div>
          </div>
        </div>
        <div class="table-body">
          <el-table :data="userTable" @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55" />
            <el-table-column prop="sysUser.userName" label="用户名" width="120" />
            <el-table-column label="操作列" min-width="220">
              <template #default="{row}">
                <el-button v-auth="btnAuth.rechargeBtn" link type="primary" icon="wallet" @click="showOperateRecharge(row)"> 充值 </el-button>
                <el-button link type="primary" icon="switch" @click="showRecharge(row)"> 积分划转 </el-button>
                <el-button type="primary" link icon="magic-stick" @click="resetPasswordFunc(row)"> 重置密码 </el-button>
                <el-button type="primary" link icon="lock" @click="getAuthCaptcha(row)"> 安全码 </el-button>
                <el-popover v-model="row.visible" placement="top" width="160">
                  <p>确定要删除此用户吗</p>
                  <div style="text-align: right; margin-top: 8px;">
                    <el-button type="primary" link @click="row.visible = false">取消</el-button>
                    <el-button type="primary" @click="deleteUserFunc(row)">确定</el-button>
                  </div>
                  <template #reference>
                    <el-button type="primary" link icon="delete">删除</el-button>
                  </template>
                </el-popover>
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

    <!--  新增成员  -->
    <el-dialog v-model="addUserDialog" title="新增成员" width="400px">
      <el-form v-model="userForm" label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="userForm.username" prop="username"/>
        </el-form-item>
        <el-form-item :minlength="6" label="新密码" prop="newPassword">
          <el-input v-model="userForm.newPassword" show-password />
        </el-form-item>
        <el-form-item :minlength="6" label="确认密码" prop="confirmPassword">
          <el-input v-model="userForm.confirmPassword" show-password />
        </el-form-item>
        <el-form-item label="账号开关" prop="enable">
          <el-switch v-model="userForm.enable" active-value="1" inactive-value="0" active-text="开启"
                     inactive-text="关闭" inline-prompt style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949;" width="60"></el-switch>
        </el-form-item>
        <el-form-item label="安全码" prop="status">
          <el-switch v-model="userForm.enableAuth" active-value="1" inactive-value="0" active-text="开启"
                     inactive-text="关闭" inline-prompt style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949;" width="60"></el-switch>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addUserClear">取消</el-button>
        <el-button type="primary" @click="addUserEnter">确定</el-button>
      </template>
    </el-dialog>

    <!-- 防爆验证码 -->
    <el-dialog v-model="showAuthCaptcha" title="重置安全码" width="360px" @close="clearAuthCaptcha">
      <el-form ref="modifyCapForm" :model="capModify" label-width="80px">
        <el-form-item label="用户ID" prop="toUid">
          <el-input v-model="capModify.ID" disabled />
        </el-form-item>
        <el-form-item label="登录密码" prop="password">
          <el-input v-model="capModify.password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showAuthCaptcha = false">取 消</el-button>
          <el-button type="primary" @click="resetAuthCaptcha">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 查看 -->
    <el-dialog v-model="showQRCode" title="安全码" width="300px" @close="closeAuthCaptcha">
      <div class="qrcode-generator">
        <div v-if="isNotSetting" style="margin-bottom: 20px">
          暂未设置安全码，请尽快设置！
        </div>
        <div v-else>
          <img :src="qrcodeUrl" alt="QR Code" style="height: 200px"/>
        </div>
        <el-button link type="primary" icon="lock" @click="resetShowAuthCaptcha"> 设置(或重置) </el-button>
      </div>
    </el-dialog>

    <!-- 积分划转 -->
    <el-dialog v-model="showRechargeVisible" title="积分划转" width="360px" @close="clearRecharge">
      <el-form :model="rechargeForm" label-width="80px">
        <el-form-item label="用户ID" prop="toUid">
          <el-input v-model="rechargeForm.toUid" disabled />
        </el-form-item>
        <el-form-item label="用户名" prop="toUsername">
          <el-input v-model="rechargeForm.toUsername" disabled />
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
    <el-dialog v-model="operateRechargeVisible" title="积分充值" width="360px" @close="clearOperateRecharge">
      <el-form :model="rechargeForm" label-width="80px">
        <el-form-item label="用户ID" prop="toUid">
          <el-input v-model="rechargeForm.toUid" disabled />
        </el-form-item>
        <el-form-item label="用户名" prop="toUsername">
          <el-input v-model="rechargeForm.toUsername" disabled />
        </el-form-item>
        <el-form-item label="充值积分" prop="recharge">
          <el-input v-model.number="rechargeForm.recharge" />
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

import { ref } from 'vue'
import { ElMessageBox, ElMessage } from 'element-plus'
export default {
  name: 'Organization',
}
</script>

<script setup>
import { findOrgUserListSelf } from '@/plugin/organization/api/organization';
import {deleteUser, selfRegister, resetCaptcha, resetPassword} from '@/api/user';
import {ElMessage, ElMessageBox} from "element-plus";
import {transferUserWallet} from "@/api/userWallet";
import { useBtnAuth } from '@/utils/btnAuth'

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
  enable: '1',
  enableAuth: '1',
})

// 增加新用户
const addUser = async() => {
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
const getUserTable = async() => {
  const res = await findOrgUserListSelf({ page: page.value, pageSize: pageSize.value, ...userSearch.value })
  if (res.code === 0) {
    page.value = res.data.page
    pageSize.value = res.data.pageSize
    total.value = res.data.total
    userTable.value = res.data.list
  }
}

// 组织用户搜索
const userSearch = ref({
  userName: '',
})

const deleteUserFunc = async(row) => {
  ElMessageBox.confirm('删除用户为高危操作，请核实是否删除？', '提示', {
    type: 'warning',
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    showClose: false,
  }).then(async() => {
    const res = await deleteUser({ id: row.sysUser.ID })
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


// 添加组织成员
const addUserEnter = async() => {
  userForm.value.enableAuth = Number(userForm.value.enableAuth)
  userForm.value.enable = Number(userForm.value.enable)
  const res = await selfRegister(userForm.value)
  if (res.code === 0) {
    ElMessage.success('添加成功')
    await getUserTable()
  }
  userForm.value = {
    username: '',
    confirmPassword: '',
    newPassword: '',
    enable: '1',
    enableAuth: '1',
  }
  addUserDialog.value = false
}

// 初始化方法
const init = async() => {
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
  ).then(async() => {
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
const resetShowAuthCaptcha = async() => {
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
          ElMessage.success('重置安全码成功！')
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

const closeAuthCaptcha = async() => {
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
  }else {
    isNotSetting.value = true;
    showQRCode.value = true;
  }
};
// ---------- 重置防爆码 end ----------

// ---------- 充值划转 ----------
const showRechargeVisible = ref(false)
const rechargeForm = ref({
  toUid: 0,
  toUsername: '',
  recharge: 0,
  type: 2,
})
const showRecharge = async(row) => {
  showRechargeVisible.value = true
  rechargeForm.value.toUid = row.sysUser.ID
  rechargeForm.value.toUsername = row.sysUser.userName
  rechargeForm.value.type = 2
  console.log(rechargeForm.value)
}

// 积分划转、充值
const transferRecharge = async(type) => {
  rechargeForm.value.type = type;
  let req = {...rechargeForm.value}
  console.log(req)
  ElMessageBox.confirm(
      '是否为此用户（充值）划转积分?',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  ).then(async() => {
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
const showOperateRecharge = async(row) => {
  operateRechargeVisible.value = true
  rechargeForm.value.toUid = row.sysUser.ID
  rechargeForm.value.toUsername = row.sysUser.userName
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

</script>

<style scoped lang="scss">
.qrcode-generator {
  display: flex;
  flex-direction: column;
  align-items: center;
}

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

<template>
  <div>
    <div class="grid grid-cols-12 w-full gap-2">
      <div class="col-span-3 h-full">
        <div class="w-full h-full bg-white px-4 py-8 rounded-lg shadow-lg box-border">
          <div class="user-card px-6 text-center bg-white shrink-0">
            <div class="flex justify-center">
              <Image
                v-model="userStore.userInfo.headerImg"
                file-type="image"
              />
            </div>
            <div class="py-6 text-center">
              <p
                v-if="!editFlag"
                class="text-3xl flex justify-center items-center gap-4"
              >
                {{ userStore.userInfo.nickName }}
              </p>
              <p class="text-gray-500 mt-2 text-md">有需求带上米，联系管理员</p>
            </div>
            <div class="w-full h-full text-left">
              <ul class="inline-block h-full w-full">
                <li class="info-list">
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
      <div class="col-span-9 ">
        <div class="bg-white h-full px-4 py-8 rounded-lg shadow-lg box-border">
          <el-tabs
            v-model="activeName"
            @tab-click="handleClick"
          >
            <el-tab-pane
              label="账号绑定"
              name="second"
            >
              <ul>
                <li class="borderd pt-2.5">
                  <p class="pb-2.5 text-xl text-gray-600">密码设置</p>
                  <p class="pb-2.5 text-lg text-gray-400">
                    修改个人密码
                    <a
                      href="javascript:void(0)"
                      class="float-right text-blue-400"
                      @click="showPassword = true"
                    >
                      <el-button type="primary" link icon="magic-stick"> 重置密码 </el-button>
                    </a>
                  </p>
                </li>
                <li class="borderd pt-2.5">
                  <p class="pb-2.5 text-xl text-gray-600">安全码设置</p>
                  <p class="pb-2.5 text-lg text-gray-400">
                    修改安全码
                    <a
                      href="javascript:void(0)"
                      class="float-right text-blue-400"
                      @click="getAuthCaptcha(userStore.userInfo)"
                    >
                      <el-button type="primary" link icon="lock"> 设置安全码 </el-button>
                    </a>
                  </p>
                </li>
              </ul>
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>
    </div>

    <el-dialog
      v-model="showPassword"
      title="修改密码"
      width="360px"
      @close="clearPassword"
    >
      <el-form
        ref="modifyPwdForm"
        :model="pwdModify"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item
          :minlength="6"
          label="原密码"
          prop="password"
        >
          <el-input
            v-model="pwdModify.password"
            show-password
          />
        </el-form-item>
        <el-form-item
          :minlength="6"
          label="新密码"
          prop="newPassword"
        >
          <el-input
            v-model="pwdModify.newPassword"
            show-password
          />
        </el-form-item>
        <el-form-item
          :minlength="6"
          label="确认密码"
          prop="confirmPassword"
        >
          <el-input
            v-model="pwdModify.confirmPassword"
            show-password
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button

            @click="showPassword = false"
          >取 消</el-button>
          <el-button

            type="primary"
            @click="savePassword"
          >确 定</el-button>
        </div>
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
  </div>
</template>

<script setup>
import { setSelfInfo, changePassword } from '@/api/user.js'
import { reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/pinia/modules/user'
import Image from '@/components/selectImage/Image.vue'
import QRCode from "qrcode";

defineOptions({
  name: 'Person',
})

const activeName = ref('second')
const rules = reactive({
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '最少6个字符', trigger: 'blur' },
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '最少6个字符', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请输入确认密码', trigger: 'blur' },
    { min: 6, message: '最少6个字符', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== pwdModify.value.newPassword) {
          callback(new Error('两次密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
})

const userStore = useUserStore()
const modifyPwdForm = ref(null)
const showPassword = ref(false)
const pwdModify = ref({})
const nickName = ref('')
const editFlag = ref(false)
const savePassword = async() => {
  modifyPwdForm.value.validate((valid) => {
    if (valid) {
      changePassword({
        password: pwdModify.value.password,
        newPassword: pwdModify.value.newPassword,
      }).then((res) => {
        if (res.code === 0) {
          ElMessage.success('修改密码成功！')
        }
        showPassword.value = false
      })
    } else {
      return false
    }
  })
}

const clearPassword = () => {
  pwdModify.value = {
    password: '',
    newPassword: '',
    confirmPassword: '',
  }
  modifyPwdForm.value.clearValidate()
}

watch(() => userStore.userInfo.headerImg, async(val) => {
  const res = await setSelfInfo({ headerImg: val })
  if (res.code === 0) {
    userStore.ResetUserInfo({ headerImg: val })
    ElMessage({
      type: 'success',
      message: '设置成功',
    })
  }
})

const openEdit = () => {
  nickName.value = userStore.userInfo.nickName
  editFlag.value = true
}

const closeEdit = () => {
  nickName.value = ''
  editFlag.value = false
}

const enterEdit = async() => {
  const res = await setSelfInfo({
    nickName: nickName.value
  })
  if (res.code === 0) {
    userStore.ResetUserInfo({ nickName: nickName.value })
    ElMessage({
      type: 'success',
      message: '设置成功',
    })
  }
  nickName.value = ''
  editFlag.value = false
}

const handleClick = (tab, event) => {
  console.log(tab, event)
}

const changePhoneFlag = ref(false)
const time = ref(0)
const phoneForm = reactive({
  phone: '',
  code: ''
})

const getCode = async() => {
  time.value = 60
  let timer = setInterval(() => {
    time.value--
    if (time.value <= 0) {
      clearInterval(timer)
      timer = null
    }
  }, 1000)
}

const closeChangePhone = () => {
  changePhoneFlag.value = false
  phoneForm.phone = ''
  phoneForm.code = ''
}

const changePhone = async() => {
  const res = await setSelfInfo({ phone: phoneForm.phone })
  if (res.code === 0) {
    ElMessage.success('修改成功')
    userStore.ResetUserInfo({ phone: phoneForm.phone })
    closeChangePhone()
  }
}

const changeEmailFlag = ref(false)
const emailTime = ref(0)
const emailForm = reactive({
  email: '',
  code: ''
})

const getEmailCode = async() => {
  emailTime.value = 60
  let timer = setInterval(() => {
    emailTime.value--
    if (emailTime.value <= 0) {
      clearInterval(timer)
      timer = null
    }
  }, 1000)
}

const closeChangeEmail = () => {
  changeEmailFlag.value = false
  emailForm.email = ''
  emailForm.code = ''
}

const changeEmail = async() => {
  const res = await setSelfInfo({ email: emailForm.email })
  if (res.code === 0) {
    ElMessage.success('修改成功')
    userStore.ResetUserInfo({ email: emailForm.email })
    closeChangeEmail()
  }
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
  console.log(row)
  let authCaptcha = row.authCaptcha;
  capModify.value = JSON.parse(JSON.stringify(row))
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
</script>

<style lang="scss">
.borderd {
  @apply border-b-2 border-solid border-gray-100 border-t-0 border-r-0 border-l-0;
    &:last-child{
      @apply border-b-0;
    }
 }

.info-list{
  @apply w-full whitespace-nowrap overflow-hidden text-ellipsis py-3 text-lg text-gray-700
}

</style>

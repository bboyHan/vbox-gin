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
            <el-tab-pane label="账号绑定" name="second">
              <ul>
                <li class="borderd pt-2.5">
                  <p class="pb-2.5 text-xl text-gray-600">密码生成器(开发中)</p>
                  <p class="pb-2.5 text-lg text-gray-400">
                    复杂密码生成器
                    <a href="javascript:void(0)" class="float-right text-blue-400" @click="genPassword = true">
                      <el-button type="primary" link icon="magic-stick"> 密码生成器(开发中) </el-button>
                    </a>
                  </p>
                </li>
                <li class="borderd pt-2.5">
                  <p class="pb-2.5 text-xl text-gray-600">密码设置</p>
                  <p class="pb-2.5 text-lg text-gray-400">
                    修改个人密码
                    <a href="javascript:void(0)" class="float-right text-blue-400" @click="showPassword = true">
                      <el-button type="primary" link icon="magic-stick"> 重置密码 </el-button>
                    </a>
                  </p>
                </li>
                <li class="borderd pt-2.5">
                  <p class="pb-2.5 text-xl text-gray-600">安全码设置</p>
                  <p class="pb-2.5 text-lg text-gray-400">
                    修改安全码
                    <a href="javascript:void(0)" class="float-right text-blue-400" @click="getAuthCaptcha(userStore.userInfo)">
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

    <el-dialog v-model="showPassword" title="修改密码" width="360px" @close="clearPassword">
      <el-form ref="modifyPwdForm" :model="pwdModify" :rules="rules" label-width="80px">
        <el-form-item :minlength="6" label="原密码" prop="password">
          <el-input v-model="pwdModify.password" show-password/>
        </el-form-item>
        <el-form-item :minlength="6" label="新密码" prop="newPassword">
          <el-input v-model="pwdModify.newPassword" show-password/>
        </el-form-item>
        <el-form-item :minlength="6" label="确认密码" prop="confirmPassword">
          <el-input v-model="pwdModify.confirmPassword" show-password/>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showPassword = false">取 消</el-button>
          <el-button type="primary" @click="savePassword">确 定</el-button>
        </div>
      </template>
    </el-dialog>
    <el-dialog v-model="genPassword" title="密码生成器" width="560px" @close="closeGenPassword">
      <div class="container">
        <div class="result">
          <div class="result__title field-title">生成的密码</div>
          <div class="result__info right">点击复制</div>
          <div class="result__info left">复制</div>
          <div class="result__viewbox" id="result">点击生成</div>
          <button id="copy-btn" style="--x: 0; --y: 0"><i class="far fa-copy"></i></button>
        </div>
        <div class="length range__slider" data-min="4" data-max="32">
          <div class="length__title field-title" data-length='0'>长度：</div>
          <input id="slider" type="range" min="4" max="32" value="16" />
        </div>

        <div class="settings">
          <span class="settings__title field-title">settings</span>
          <div class="setting">
            <input type="checkbox" id="uppercase" checked />
            <label for="uppercase">包含大写</label>
          </div>
          <div class="setting">
            <input type="checkbox" id="lowercase" checked />
            <label for="lowercase">包含小写</label>
          </div>
          <div class="setting">
            <input type="checkbox" id="number" checked />
            <label for="number">包括数字</label>
          </div>
          <div class="setting">
            <input type="checkbox" id="symbol" />
            <label for="symbol">包括符号</label>
          </div>
        </div>
        <button class="btn generate" id="generate">生成密码</button>
      </div>
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
const genPassword = ref(false)
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

const closeGenPassword = () => {
  genPassword.value = false
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
<style lang="css" scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  width: 100%;
  height: 100%;
  background-image: linear-gradient(to top, #209cff 0%, #68e0cf 100%);
}

button {
  border: 0;
  outline: 0;
}

.container {
  margin: 40px auto;
  width: 400px;
  height: 600px;
  padding: 10px 25px;
  background: #0a0e31;
  border-radius: 10px;
  box-shadow: 0 0 5px rgba(0, 0, 0, 0.45), 0 4px 8px rgba(0, 0, 0, 0.35), 0 8px 12px rgba(0, 0, 0, 0.15);
  font-family: "Montserrat";
}
.container h2.title {
  font-size: 1.75rem;
  margin: 10px -5px;
  margin-bottom: 30px;
  color: #fff;
}

.result {
  position: relative;
  width: 100%;
  height: 65px;
  overflow: hidden;
}
.result__info {
  position: absolute;
  bottom: 4px;
  color: #fff;
  font-size: 0.8rem;
  transition: all 150ms ease-in-out;
  transform: translateY(200%);
  opacity: 0;
}
.result__info.right {
  right: 8px;
}
.result__info.left {
  left: 8px;
}
.result__viewbox {
  width: 100%;
  height: 100%;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  color: #fff;
  text-align: center;
  line-height: 65px;
}
.result #copy-btn {
  position: absolute;
  top: var(--y);
  left: var(--x);
  width: 38px;
  height: 38px;
  background: #fff;
  border-radius: 50%;
  opacity: 0;
  transform: translate(-50%, -50%) scale(0);
  transition: all 350ms cubic-bezier(0.175, 0.885, 0.32, 1.275);
  cursor: pointer;
  z-index: 2;
}
.result #copy-btn:active {
  box-shadow: 0 0 0 200px rgba(255, 255, 255, 0.08);
}
.result:hover #copy-btn {
  opacity: 1;
  transform: translate(-50%, -50%) scale(1.35);
}

.field-title {
  position: absolute;
  top: -10px;
  left: 8px;
  transform: translateY(-50%);
  font-weight: 800;
  color: rgba(255, 255, 255, 0.5);
  text-transform: uppercase;
  font-size: 0.65rem;
  pointer-events: none;
  user-select: none;
}

.options {
  width: 100%;
  height: auto;
  margin: 50px 0;
}

.range__slider {
  position: relative;
  width: 100%;
  height: calc(65px - 10px);
  display: flex;
  justify-content: center;
  align-items: center;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  margin: 30px 0;
}
.range__slider::before, .range__slider::after {
  position: absolute;
  color: #fff;
  font-size: 0.9rem;
  font-weight: bold;
}
.range__slider::before {
  content: attr(data-min);
  left: 10px;
}
.range__slider::after {
  content: attr(data-max);
  right: 10px;
}
.range__slider .length__title::after {
  content: attr(data-length);
  position: absolute;
  right: -16px;
  font-variant-numeric: tabular-nums;
  color: #fff;
}

#slider {
  -webkit-appearance: none;
  width: calc(100% - (70px));
  height: 2px;
  border-radius: 5px;
  background: rgba(255, 255, 255, 0.314);
  outline: none;
  padding: 0;
  margin: 0;
  cursor: pointer;
}
#slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: white;
  cursor: pointer;
  transition: all 0.15s ease-in-out;
}
#slider::-webkit-slider-thumb:hover {
  background: #d4d4d4;
  transform: scale(1.2);
}
#slider::-moz-range-thumb {
  width: 20px;
  height: 20px;
  border: 0;
  border-radius: 50%;
  background: white;
  cursor: pointer;
  transition: background 0.15s ease-in-out;
}
#slider::-moz-range-thumb:hover {
  background: #d4d4d4;
}

.settings {
  position: relative;
  height: auto;
  widows: 100%;
  display: flex;
  flex-direction: column;
}
.settings .setting {
  position: relative;
  width: 100%;
  height: calc(65px - 10px);
  background: rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  display: flex;
  align-items: center;
  padding: 10px 25px;
  color: #fff;
  margin-bottom: 8px;
}
.settings .setting input {
  opacity: 0;
  position: absolute;
}
.settings .setting input + label {
  user-select: none;
}
.settings .setting input + label::before, .settings .setting input + label::after {
  content: "";
  position: absolute;
  transition: 150ms cubic-bezier(0.24, 0, 0.5, 1);
  transform: translateY(-50%);
  top: 50%;
  right: 10px;
  cursor: pointer;
}
.settings .setting input + label::before {
  height: 30px;
  width: 50px;
  border-radius: 30px;
  background: rgba(214, 214, 214, 0.434);
}
.settings .setting input + label::after {
  height: 24px;
  width: 24px;
  border-radius: 60px;
  right: 32px;
  background: #fff;
}
.settings .setting input:checked + label:before {
  background: #5d68e2;
  transition: all 150ms cubic-bezier(0, 0, 0, 0.1);
}
.settings .setting input:checked + label:after {
  right: 14px;
}
.settings .setting input:focus + label:before {
  box-shadow: 0 0 0 2px rgba(255, 255, 255, 0.75);
}

.btn.generate {
  user-select: none;
  position: relative;
  width: 100%;
  height: 50px;
  margin: 10px 0;
  border-radius: 8px;
  color: #fff;
  border: none;
  background-image: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  letter-spacing: 1px;
  font-weight: bold;
  text-transform: uppercase;
  cursor: pointer;
  transition: all 150ms ease;
}
.btn.generate:active {
  transform: translateY(-3%);
  box-shadow: 0 4px 8px rgba(255, 255, 255, 0.08);
}

.support {
  position: fixed;
  right: 10px;
  bottom: 10px;
  padding: 10px;
  display: flex;
}

a {
  margin: 0 20px;
  color: #fff;
  font-size: 2rem;
  transition: all 400ms ease;
}

a:hover {
  color: #222;
}
</style>
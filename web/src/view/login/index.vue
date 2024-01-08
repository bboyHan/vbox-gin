<template>
  <div id="userLayout" class="w-full h-full relative">
    <div
      class="rounded-lg flex items-center justify-evenly w-full h-full bg-white md:w-screen md:h-screen"
    >
      <div class="md:w-3/5 w-10/12 h-full flex items-center justify-evenly">
        <div class="oblique h-[130%] w-3/5 bg-white transform -rotate-12 absolute -ml-52" />
        <!-- 分割斜块 -->
        <div class="z-[999] pt-12 pb-10 md:w-96 w-full  rounded-lg flex flex-col justify-between box-border">
          <div>
            <div class="flex items-center justify-center">
              <img class="w-24" src="@/assets/logo.png" alt>
            </div>
            <div class="mb-9">
            </div>
            <el-form
              ref="loginForm"
              :model="loginFormData"
              :rules="rules"
              :validate-on-rule-change="false"
              @keyup.enter="submitForm"
            >
              <el-form-item prop="username" class="mb-6">
                <el-input v-model="loginFormData.username" size="large" placeholder="用户名" suffix-icon="user" />
              </el-form-item>
              <el-form-item prop="password" class="mb-6">
                <el-input
                  v-model="loginFormData.password"
                  show-password
                  size="large"
                  type="password"
                  placeholder="密码"
                />
              </el-form-item>
              <el-form-item v-if="loginFormData.openCaptcha" prop="captcha" class="mb-6">
                <div class="flex w-full justify-between">
                  <el-input v-model="loginFormData.captcha" placeholder="验证码" size="large" class="flex-1 mr-5" />
                  <div class="w-1/3 h-11 bg-[#c3d4f2] rounded">
                    <img v-if="picPath" class="w-full h-full" :src="picPath" alt="验证码" @click="loginVerify()">
                  </div>
                </div>
              </el-form-item>
              <el-form-item prop="authCaptcha" class="mb-6">
                <el-input v-model="loginFormData.authCaptcha" size="large" placeholder="安全码" />
              </el-form-item>
              <el-form-item class="mb-6">
                <el-button class="shadow shadow-blue-600 h-11 w-full" type="primary" size="large" @click="submitForm">登 录</el-button>
              </el-form-item>
            </el-form>
          </div>
        </div>
      </div>
      <div class="background-image">
        <img src="@/assets/login_bg.svg" alt="banner">
      </div>
    </div>
  </div>
</template>

<script setup>
import { captcha } from '@/api/user'
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/pinia/modules/user'

defineOptions({
  name: 'Login',
})

// 验证函数
const checkUsername = (rule, value, callback) => {
  if (value.length < 4) {
    return callback(new Error('请输入正确的用户名'))
  } else {
    callback()
  }
}
const checkPassword = (rule, value, callback) => {
  if (value.length < 6) {
    return callback(new Error('请输入正确的密码'))
  } else {
    callback()
  }
}

// 获取验证码
const loginVerify = () => {
  captcha({}).then(async(ele) => {
    rules.captcha.push({
      max: ele.data.captchaLength,
      min: ele.data.captchaLength,
      message: `请输入${ele.data.captchaLength}位验证码`,
      trigger: 'blur',
    })
    picPath.value = ele.data.picPath
    loginFormData.captchaId = ele.data.captchaId
    loginFormData.openCaptcha = ele.data.openCaptcha
  })
}
loginVerify()

// 登录相关操作
const loginForm = ref(null)
const picPath = ref('')
const loginFormData = reactive({
  username: '',
  password: '',
  authCaptcha: '',
  captcha: '',
  captchaId: '',
  openCaptcha: false,
})
const rules = reactive({
  username: [{ validator: checkUsername, trigger: 'blur' }],
  password: [{ validator: checkPassword, trigger: 'blur' }],
  captcha: [
    {
      message: '验证码格式不正确',
      trigger: 'blur',
    },
  ],
})

const userStore = useUserStore()
const login = async() => {
  return await userStore.LoginIn(loginFormData)
}
const submitForm = () => {
  loginForm.value.validate(async(v) => {
    if (v) {
      const flag = await login()
      if (!flag) {
        loginVerify()
      }
    } else {
      ElMessage({
        type: 'error',
        message: '请正确填写登录信息',
        showClose: true,
      })
      loginVerify()
      return false
    }
  })
}
</script>
<style scoped>
.background-image {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}
.background-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>

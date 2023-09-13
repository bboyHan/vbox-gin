import { login, getUserInfo, setSelfInfo, getOwnerUserList } from '@/api/user'
import { getVboxChannelProductAll } from '@/api/vboxChannelProduct'
import { jsonInBlacklist } from '@/api/jwt'
import router from '@/router/index'
import { ElLoading, ElMessage } from 'element-plus'
import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { useRouterStore } from './router'

export const useUserStore = defineStore('user', () => {
  const loadingInstance = ref(null)
  const usersItem = ref([{
    ID: '',
    userName: '',
  }])

  //channel product list
  const channelProductList = ref([{
    channelCode: '',
    productName: '',
    productId: '',
  }])

  const chanCodeMap = ref(new Map())
  const ownerUsersMap = ref(new Map())

  const setChannelProductList = (val) => {
    channelProductList.value = val
    for (const item of channelProductList.value) {
      chanCodeMap.value.set(item.channelCode + '', item.productName)
    }
    console.log(channelProductList.value)
    console.log(chanCodeMap.value)
  }

  /* 获取通道产品信息*/
  const GetChannelProductList = async() => {
    const res = await getVboxChannelProductAll()
    if (res.code === 0) {
      // console.log("---xxx---")
      // console.log(res.data)
      setChannelProductList(res.data.list)
    }
    return res
  }

  const LoadAllUser = async ()  => {

    const res = await getOwnerUserList({ page: 1, pageSize: 9999 })
    console.log('== res ==>' + JSON.stringify(res))
    if (res.code === 0) {
      // console.log('== res.data.marks==>' + JSON.stringify(res.data.list))
      usersItem.value = res.data.list
      for (const item of usersItem.value) {
        ownerUsersMap.value.set(item.ID + '', item.userName + '')
      }
    }
    return usersItem.value
  }

  const userInfo = ref({
    uuid: '',
    nickName: '',
    headerImg: '',
    authority: {},
    sideMode: 'dark',
    activeColor: 'var(--el-color-primary)',
    baseColor: '#fff'
  })
  const token = ref(window.localStorage.getItem('token') || '')
  const setUserInfo = (val) => {
    userInfo.value = val
  }

  const setToken = (val) => {
    token.value = val
  }

  const NeedInit = () => {
    token.value = ''
    window.localStorage.removeItem('token')
    localStorage.clear()
    router.push({ name: 'Init', replace: true })
  }

  const ResetUserInfo = (value = {}) => {
    userInfo.value = {
      ...userInfo.value,
      ...value
    }
  }
  /* 获取用户信息*/
  const GetUserInfo = async() => {
    const res = await getUserInfo()
    if (res.code === 0) {
      setUserInfo(res.data.userInfo)
    }
    return res
  }
  /* 登录*/
  const LoginIn = async(loginInfo) => {
    loadingInstance.value = ElLoading.service({
      fullscreen: true,
      text: '登录中，请稍候...',
    })
    try {
      const res = await login(loginInfo)
      if (res.code === 0) {
        setUserInfo(res.data.user)
        setToken(res.data.token)
        const routerStore = useRouterStore()
        await routerStore.SetAsyncRouter()
        const asyncRouters = routerStore.asyncRouters
        asyncRouters.forEach(asyncRouter => {
          router.addRoute(asyncRouter)
        })

        if (!router.hasRoute(userInfo.value.authority.defaultRouter)) {
          ElMessage.error("请联系管理员进行授权")
        } else {
          await router.replace({ name: userInfo.value.authority.defaultRouter })
        }

        loadingInstance.value.close()

        const isWin = ref(/windows/i.test(navigator.userAgent))
        if (isWin.value) {
          window.localStorage.setItem('osType', 'WIN')
        } else {
          window.localStorage.setItem('osType', 'MAC')
        }
        return true
      }
    } catch (e) {
      loadingInstance.value.close()
    }
    loadingInstance.value.close()
  }
  /* 登出*/
  const LoginOut = async() => {
    const res = await jsonInBlacklist()
    if (res.code === 0) {
      token.value = ''
      sessionStorage.clear()
      localStorage.clear()
      router.push({ name: 'Login', replace: true })
      window.location.reload()
    }
  }
  /* 清理数据 */
  const ClearStorage = async() => {
    token.value = ''
    sessionStorage.clear()
    localStorage.clear()
  }
  /* 设置侧边栏模式*/
  const changeSideMode = async(data) => {
    const res = await setSelfInfo({ sideMode: data })
    if (res.code === 0) {
      userInfo.value.sideMode = data
      ElMessage({
        type: 'success',
        message: '设置成功'
      })
    }
  }

  const mode = computed(() => userInfo.value.sideMode)
  const sideMode = computed(() => {
    if (userInfo.value.sideMode === 'dark') {
      return '#191a23'
    } else if (userInfo.value.sideMode === 'light') {
      return '#fff'
    } else {
      return userInfo.value.sideMode
    }
  })
  const baseColor = computed(() => {
    if (userInfo.value.sideMode === 'dark') {
      return '#fff'
    } else if (userInfo.value.sideMode === 'light') {
      return '#191a23'
    } else {
      return userInfo.value.baseColor
    }
  })
  const activeColor = computed(() => {
    return 'var(--el-color-primary)'
  })

  watch(() => token.value, () => {
    window.localStorage.setItem('token', token.value)
  })

  return {
    userInfo,
    token,
    channelProductList,
    chanCodeMap,
    ownerUsersMap,
    NeedInit,
    ResetUserInfo,
    GetUserInfo,
    LoadAllUser,
    GetChannelProductList,
    LoginIn,
    LoginOut,
    changeSideMode,
    mode,
    sideMode,
    setToken,
    baseColor,
    activeColor,
    loadingInstance,
    ClearStorage
  }
})

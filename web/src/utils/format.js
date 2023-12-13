import { formatTimeToStr } from '@/utils/date'
import { getDict } from '@/utils/dictionary'
import { codeToText } from 'element-china-area-data';

// 付款码状态
export const formatStatus2ETag = (status) => {
  if(status >= 500) {
    return 'danger'
  } else if(status >= 200 && status <= 299) {
    return 'success'
  } else if(status >= 300 && status <= 499) {
    return ''
  } else {
    return 'info'
  }
}

// 格式化时间为时分秒的形式(时分秒)
export const formatTime = (seconds) => {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const remainingSeconds = seconds % 60;
  return `${hours} 时 ${minutes} 分 ${remainingSeconds} 秒`;
};

// 字符串连接
export const formatJoin = (separator, ...str) => {
  return str.join(separator);
}

// 付款码状态
export const formatPayCodeStatus = (status) => {
  if(status === 2) {
    return '待使用'
  } else if(status === 1) {
    return '已使用'
  } else if(status === 3) {
    return '已失效'
  } else {
    return '-'
  }
}

// 付款码状态
export const formatRegionCode = (locationCode, isStrict) => {
  if (!locationCode){
    return '-'
  }
  if (isStrict){

  }else {
    return codeToText[locationCode.slice(0, 2)]
  }
}

// 付款码样式
export const formatPayCodeColor = (status) => {
  if (status === 3) {
    return '#606266'
  } else if(status === 2) {
    return '#a64406'
  } else if(status === 1) {
    return '#05811d'
  } else {
    return ''
  }
}

// 订单付款样式
export const formatPayedColor = (status, acId) => {
  if (status === 2 && acId === "") { //匹配账号未支付
    return '#a64406'
  } else if(status === 2 && acId !== "") { //未匹配账号的
    return '#606266'
  } else if(status === 1) {
    return '#05811d'
  } else if(status === 3) {
    return '#3118dc'
  } else {
    return ''
  }
}

export const formatPayed = (status, acId) => {
  console.log(status + " --- " +acId)
  if (status === 2 && acId !== "") {
    return '待支付'
  } else if(status === 2 && acId === "") {
    return '待取码'
  } else if(status === 1) {
    return '已支付'
  } else if(status === 3) {
    return '支付超时'
  } else {
    return '-'
  }
}

export const formatNotifyColor = (status) => {
  if (status === 2) {
    return '#606266'
  } else if(status === 1) {
    return '#05811d'
  } else {
    return ''
  }
}

export const formatNotify = (status) => {
  if (status === 2) {
    return '未通知'
  } else if(status === 1) {
    return '已通知'
  } else {
    return '-'
  }
}

export const formatBoolean = (bool) => {
  if (bool !== null) {
    return bool ? '是' : '否'
  } else {
    return ''
  }
}
export const formatDate = (time) => {
  if (time !== null && time !== '') {
    var date = new Date(time)
    return formatTimeToStr(date, 'yyyy-MM-dd hh:mm:ss')
  } else {
    return ''
  }
}

export const formatUtcTimestamp = (utcTimestamp) => {
  if (utcTimestamp !== null && utcTimestamp !== '') {
    const date = new Date(0) // 创建一个时间对象，并将其初始化为 1970 年 1 月 1 日 00:00:00 UTC 时间
    date.setUTCSeconds(utcTimestamp) // 将时间戳设置为 UTC 时间
    return formatTimeToStr(date, 'yyyy-MM-dd hh:mm:ss')
  } else {
    return ''
  }
}

export const filterDict = (value, options) => {
  const rowLabel = options && options.filter(item => item.value === value)
  return rowLabel && rowLabel[0] && rowLabel[0].label
}

export const getDictFunc = async(type) => {
  const dicts = await getDict(type)
  return dicts
}

const path = import.meta.env.VITE_BASE_PATH + ':' + import.meta.env.VITE_SERVER_PORT + '/'
export const ReturnArrImg = (arr) => {
  const imgArr = []
  if (arr instanceof Array) { // 如果是数组类型
    for (const arrKey in arr) {
      if (arr[arrKey].slice(0, 4) !== 'http') {
        imgArr.push(path + arr[arrKey])
      } else {
        imgArr.push(arr[arrKey])
      }
    }
  } else { // 如果不是数组类型
    if (arr.slice(0, 4) !== 'http') {
      imgArr.push(path + arr)
    } else {
      imgArr.push(arr)
    }
  }
  return imgArr
}

export const onDownloadFile = (url) => {
  window.open(path + url)
}

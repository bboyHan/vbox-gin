import { formatTimeToStr } from '@/utils/date'
import { getDict } from '@/utils/dictionary'

export const formatCallbackStatusColor = (value) => {
  if (value === 1) {
    return 'success'
  } else if(value === 2) {
    return 'info'
  }
}

export const formatCallbackOrderStatus = (value) => {
  if (value === 1) {
    return '通知成功'
  } else if(value === 2) {
    return '未通知'
  }
}

export const formatOrderStatusColor = (value) => {
  if (value === 1) {
    return 'success'
  } else if(value === 2) {
    return 'info'
  } else if(value === 3) {
    return ''
  } else if(value === 4) {
    return 'warning'
  }
}

export const formatOrderStatus = (value) => {
  if (value === 1) {
    return '支付成功'
  } else if(value === 2) {
    return '等待支付'
  } else if(value === 3) {
    return '支付超时'
  } else if(value === 4) {
    return '等待取码'
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
  let imgArr = []
  if (arr instanceof Array){ // 如果是数组类型
    for (let arrKey in arr) {
      if (arr[arrKey].slice(0, 4) !== 'http'){
        imgArr.push(path + arr[arrKey])
      }else {
        imgArr.push(arr[arrKey])
      }
    }
  }else { // 如果不是数组类型
    if (arr.slice(0, 4) !== 'http'){
      imgArr.push(path + arr)
    }else {
      imgArr.push(arr)
    }
  }
  return imgArr
}

export const onDownloadFile = (url) => {
  window.open(path + url)

}
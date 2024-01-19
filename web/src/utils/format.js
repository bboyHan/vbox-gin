import {formatTimeToStr} from '@/utils/date'
import {getDict} from '@/utils/dictionary'
import {codeToText} from 'element-china-area-data';

// 格式化运营商
export const formatOPDesc = (operator) => {
    let opMsg = '';
    if (operator === 'yidong') {
        opMsg = '运营商：移动'
    } else if (operator === 'liantong') {
        opMsg = '运营商：联通'
    } else if (operator === 'dianxin') {
        opMsg = '运营商：电信'
    } else if (operator === 'default') {
        opMsg = '运营商：默认'
    } else {
        opMsg = '运营商：未知'
    }
    return opMsg
}

// 格式化运营商（简化）
export const formatOPSimple = (operator) => {
    let opMsg = '';
    if (operator === 'yidong') {
        opMsg = '移动'
    } else if (operator === 'liantong') {
        opMsg = '联通'
    } else if (operator === 'dianxin') {
        opMsg = '电信'
    } else if (operator === 'default') {
        opMsg = '默认'
    } else {
        opMsg = '未知'
    }
    return opMsg
}

// 格式化金额
export const formatMoneyDesc = (money) => {
    return '金额：' + money + "元"
}

// 付款码状态
export const formatStatus2ETag = (status) => {
    if (status >= 500) {
        return 'danger'
    } else if (status >= 200 && status <= 299) {
        return 'success'
    } else if (status >= 300 && status <= 499) {
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
    if (status === 2) {
        return '待使用'
    } else if (status === 1) {
        return '已使用'
    } else if (status === 3) {
        return '已失效'
    } else if (status === 4) {
        return '冷却中'
    } else {
        return '-'
    }
}

// 付款码状态
export const formatRegionCode = (locationCode, isStrict) => {
    if (!locationCode) {
        return '-'
    }
    if (String(locationCode) === '10') {
        return '默认'
    }
    if (isStrict) {
        return codeToText[locationCode]
    } else {
        return codeToText[locationCode.slice(0, 2)]
    }
}

// 付款码样式
export const formatPayCodeColor = (status) => {
    if (status === 3) {
        return '#606266'
    } else if (status === 4) {
        return '#2b7fff'
    } else if (status === 2) {
        return '#a64406'
    } else if (status === 1) {
        return '#05811d'
    } else {
        return ''
    }
}

// 订单付款样式
export const formatPayedColor = (status, acId) => {
    if (status === 2 && acId === "") { //匹配账号未支付
        return '#a64406'
    } else if (status === 2 && acId !== "") { //未匹配账号的
        return '#606266'
    } else if (status === 1) {
        return '#05811d'
    } else if (status === 3) {
        return '#3118dc'
    } else if (status === 0) {
        return '#dc1818'
    } else {
        return ''
    }
}

export const formatPayed = (status, acId) => {
    console.log(status + " --- " + acId)
    if (status === 2 && acId !== "") {
        return '待支付'
    } else if (status === 2 && acId === "") {
        return '待取码'
    } else if (status === 1) {
        return '已支付'
    } else if (status === 3) {
        return '支付超时'
    } else if (status === 0) {
        return '匹配失败'
    } else {
        return '-'
    }
}

export const formatCDStatusColor = (status) => {
    if (status === 2) {
        return '#606266'
    } else if (status === 1) {
        return '#2b7fff'
    } else {
        return ''
    }
}

export const formatCDStatus = (status) => {
    if (status === 2) {
        return '冷却中'
    } else if (status === 1) {
        return '正常'
    } else {
        return '-'
    }
}

export const formatNotifyColor = (status) => {
    if (status === 2) {
        return '#606266'
    } else if (status === 1) {
        return '#05811d'
    } else {
        return ''
    }
}

export const formatNotify = (status) => {
    if (status === 2) {
        return '未通知'
    } else if (status === 1) {
        return '已通知'
    } else {
        return '-'
    }
}

// 补单
export const formatHandNotify = (status) => {
    if (status === 2) {
        return '默认'
    } else if (status === 1) {
        return '已补单'
    } else if (status === 3) {
        return '候补单'
    } else {
        return '-'
    }
}

// 补单
export const formatHandNotifyColor = (status) => {
    if (status === 2) {
        return '#606266'
    } else if (status === 1) {
        return '#05811d'
    } else if (status === 3) {
        return '#3118dc'
    } else {
        return ''
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

export const getDictFunc = async (type) => {
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

export function calculatePercentage(num, total) {
    if (isNaN(num) || isNaN(total)) {
        return 0;
    }

    num = parseFloat(num);
    total = parseFloat(total);

    const percentage = total <= 0 ? 0 : Math.round((num / total) * 10000) / 100.0;

    return percentage;
}

export function formatMoney(num, prefix) {
    if (isNaN(num)) {
        num = 0
    }
    if (!prefix) prefix = "￥";
    let money = String(num).replace(/\B(?=(\d{3})+(?!\d))/g, ',')

    return prefix + money;
}

//Body格式化，如果是json解析成json，如果不是json，就返回原值
export const fmtBody = (value) => {
    try {
        return JSON.parse(value)
    } catch (err) {
        return value
    }
}

// 缺省，如果内容大于200个字符，则后面后省略号代替
export const fmtSimpleBody = (value) => {
    try {
        if (value.length > 200) {
            return value.substring(0, 200) + '......'
        } else {
            return value
        }
    } catch (err) {
        return value
    }
}
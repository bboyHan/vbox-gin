
export const typeMap = {
  0: '原生',
  1: '引导',
  2: '预产',
}

export const payTypeMap = {
  'wechat': '微信',
  'alipay': '支付宝',
  'app': '三方APP',
}

export const mustUint = (rule, value, callback) => {
  if (!/^[0-9]*[1-9][0-9]*$/.test(value)) {
    return callback(new Error('请输入正整数'))
  }
  return callback()
}
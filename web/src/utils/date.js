// 对Date的扩展，将 Date 转化为指定格式的String
// 月(M)、日(d)、小时(h)、分(m)、秒(s)、季度(q) 可以用 1-2 个占位符，
// 年(y)可以用 1-4 个占位符，毫秒(S)只能用 1 个占位符(是 1-3 位的数字)
// (new Date()).Format("yyyy-MM-dd hh:mm:ss.S") ==> 2006-07-02 08:09:04.423
// (new Date()).Format("yyyy-M-d h:m:s.S")      ==> 2006-7-2 8:9:4.18
// eslint-disable-next-line no-extend-native
Date.prototype.Format = function(fmt) {
  var o = {
    'M+': this.getMonth() + 1, // 月份
    'd+': this.getDate(), // 日
    'h+': this.getHours(), // 小时
    'm+': this.getMinutes(), // 分
    's+': this.getSeconds(), // 秒
    'q+': Math.floor((this.getMonth() + 3) / 3), // 季度
    'S': this.getMilliseconds() // 毫秒
  }
  if (/(y+)/.test(fmt)) { fmt = fmt.replace(RegExp.$1, (this.getFullYear() + '').substr(4 - RegExp.$1.length)) }
  for (var k in o) {
    if (new RegExp('(' + k + ')').test(fmt)) { fmt = fmt.replace(RegExp.$1, (RegExp.$1.length === 1) ? (o[k]) : (('00' + o[k]).substr(('' + o[k]).length))) }
  }
  return fmt
}

export function formatTimeToStr(times, pattern) {
  var d = new Date(times).Format('yyyy-MM-dd hh:mm:ss')
  if (pattern) {
    d = new Date(times).Format(pattern)
  }
  return d.toLocaleString()
}


export function get5MinNearlyOneHour() {
  // 获取当前时间
  let now = new Date();

  // 计算当前时间的分钟数
  let currentMinute = now.getMinutes();

  // 计算最新的时间点
  let latestMinute = Math.floor(currentMinute / 5) * 5;

  // 创建一个空数组来存储时间点
  let timePoints = [];

  // 循环生成时间点
  for (let i = 0; i < 12; i++) { // 12 表示一个小时内的时间点数量
                                 // 计算每个时间点的小时和分钟数
    let hour = now.getHours();
    let minute = latestMinute - (i * 5);

    // 处理小时和分钟的边界情况
    if (minute < 0) {
      hour -= 1;
      minute += 60;
    }

    // 格式化小时和分钟，确保是两位数
    hour = hour.toString().padStart(2, '0');
    minute = minute.toString().padStart(2, '0');

    // 拼接时间点字符串并添加到数组中
    let timePoint = hour + ':' + minute;
    timePoints.unshift(timePoint); // 使用 unshift() 方法将时间点添加到数组的开头
  }

// 输出时间点数组
//   console.log(timePoints);
  return timePoints
}
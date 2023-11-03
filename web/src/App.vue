<template>
  <div id="app">
    <el-config-provider :locale="zhCn">
      <router-view />
    </el-config-provider>
  </div>
</template>

<script setup>
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
// element 2.3.8之前使用下面的语句
// import zhCn from 'element-plus/lib/locale/lang/zh-cn'

defineOptions({
  name: 'App'
})

</script>
<style lang="scss">
@tailwind base;
@tailwind components;
@tailwind utilities;
// 引入初始化样式
#app {
  background: #eee;
  height: 100vh;
  overflow: hidden;
  font-weight: 400 !important;
}
.el-button{
  font-weight: 400 !important;
}

.el-table__body, .el-table__header {
  border-collapse: collapse;
  border-bottom-width: 2px;
  border-bottom-style: solid;
}

.el-table__body th, .el-table__header th {
  border-bottom-width: 2px;
  border-bottom-style: solid;
}

.el-table__body tr, .el-table__header tr {
  border-bottom-width: 2px;
  border-bottom-style: solid;
}

.el-table__body td, .el-table__header td {
  border-right-width: 2px;
  border-right-style: solid;
}

// loading
.el-button .custom-loading .circular {
  margin-right: 6px;
  width: 18px;
  height: 18px;
  animation: loading-rotate 2s linear infinite;
}
.el-button .custom-loading .circular .path {
  animation: loading-dash 1.5s ease-in-out infinite;
  stroke-dasharray: 90, 150;
  stroke-dashoffset: 0;
  stroke-width: 2;
  stroke: var(--el-button-text-color);
  stroke-linecap: round;
}
</style>

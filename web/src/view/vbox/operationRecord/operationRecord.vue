<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="结果状态码">
          <el-input v-model="searchInfo.status" placeholder="搜索条件"/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit" >查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table ref="multipleTable" :data="tableData" style="width: 100%"  row-key="ID" @selection-change="handleSelectionChange">
        <el-table-column align="left" label="日期" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作人" width="140">
          <template #default="scope">
            <div>{{ scope.row.user.userName }}</div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态码" prop="status" width="120">
          <template #default="scope">
            <div>
              <el-tag effect="dark" :type="formatStatus2ETag(scope.row.status)">{{ scope.row.status }}</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="请求IP" prop="ip" width="120"/>
        <el-table-column align="left" label="处理记录" prop="path">
          <template #default="scope">
            <div>
              <el-popover v-if="scope.row.resp" placement="left-start" trigger="click">
                <div class="popover-box">
                  <pre>{{ fmtBody(scope.row.resp) }}</pre>
                </div>
                <template #reference>
                  <el-icon style="cursor: pointer;"><warning /></el-icon>
                </template>
              </el-popover>
              <span v-else>无</span>
              {{ fmtBody(scope.row.resp) }}
            </div>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { formatStatus2ETag } from '@/utils/format'
import {
  getSysOperationRecordSelf,
} from '@/api/sysOperationRecord' // 此处请自行替换地址
import { formatDate } from '@/utils/format'
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import {Warning} from "@element-plus/icons-vue";

defineOptions({
  name: 'SysOperationRecord'
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const onReset = () => {
  searchInfo.value = {}
}
// 条件搜索前端看此方法
const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  if (searchInfo.value.status === '') {
    searchInfo.value.status = null
  }
  getTableData()
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async() => {
  const table = await getSysOperationRecordSelf({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value,
  })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

const multipleSelection = ref([])
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

const fmtBody = (value) => {
  try {
    return JSON.parse(value)
  } catch (err) {
    return value
  }
}

</script>

<style lang="scss">
.table-expand {
  padding-left: 60px;
  font-size: 0;
  label {
    width: 90px;
    color: #99a9bf;
    .el-form-item {
      margin-right: 0;
      margin-bottom: 0;
      width: 50%;
    }
  }
}
.popover-box {
  background: #112435;
  color: #f08047;
  height: 300px;
  width: 420px;
  overflow: auto;
}
.popover-box::-webkit-scrollbar {
  display: none; /* Chrome Safari */
}
</style>

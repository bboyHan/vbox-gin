<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
        <el-form-item label="IP" prop="pay_ip">
          <el-input v-model="searchInfo.pay_ip" placeholder="输入需要查询的IP" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-organization-box">
      <div class="table-body">
        <el-table style="width: 100%" tooltip-effect="dark" :data="tableData">
<!--          国家|区域|省份|城市|ISP -->
          <el-table-column align="center" label="IP" prop="ip" width="120" />
          <el-table-column align="center" label="国家" prop="country" width="120" />
          <el-table-column align="center" label="区域" prop="region" width="120" />
          <el-table-column align="center" label="省份" prop="province" width="120" />
          <el-table-column align="center" label="城市" prop="city" width="120" />
          <el-table-column align="center" label="ISP" prop="isp" width="120" />
        </el-table>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ElMessage, ElMessageBox } from 'element-plus'
import {reactive, ref} from 'vue'
import {getVboxProxyList} from "@/api/vboxProxy";
import {queryIpRegion} from "@/api/payOrder";

const searchRule = reactive({
  createdAt: [
    { validator: (rule, value, callback) => {}, trigger: 'change' }
  ],
})
const elSearchFormRef = ref()
const tableData = ref([])
const searchInfo = ref({})

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    getTableData()
  })
}

const getTableData = async() => {
  const table = await queryIpRegion({ ...searchInfo.value })
  console.log(table)

  let resData = table.data;
  let country = resData.split("|")[0];
  let region = resData.split("|")[1];
  let province = resData.split("|")[2];
  let city = resData.split("|")[3];
  let isp = resData.split("|")[4];

  let tmpData = [{
    ip : searchInfo.value.pay_ip,
    country,
    region,
    province,
    city,
    isp,
  }]
  if (table.code === 0) {
    tableData.value = tmpData;
  }else {
    ElMessage({
      type: 'error',
      message: '请输入合规的ip地址'
    })
  }
}


</script>
<script>
export default {
name: 'Ip',
}
</script>

<style scoped lang="scss">
</style>
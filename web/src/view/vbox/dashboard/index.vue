<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
        <el-form-item label="通道ID" prop="cid">
          <el-cascader
            v-model="searchInfo.cid"
            :options="channelCodeOptions"
            :props="channelCodeProps"
            @change="handleChange"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">分析</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="data-center-box">
      <CenterCard title="当前通道账号">
        <template #action>
          <span class="gvaIcon-prompt" style="color: #999" />
        </template>
        <template #body>
          <ReclaimMileage :channel-code="searchInfo.cid" :acc-on="accOn" :acc-off="accOff" :acc-total="accTotal"/>
        </template>
      </CenterCard>
      <CenterCard title="订单统计">
        <template #action>
          <span class="gvaIcon-prompt" style="color: #999" />
        </template>
        <template #body>
          <Order :channel-code="searchInfo.cid"/>
        </template>
      </CenterCard>

      <CenterCard title="近1小时实时成单" style="grid-column-start: span 3;">
        <template #action>
          <span class="gvaIcon-prompt" style="color: #999" />
        </template>
        <template #body>
          <lineCharts :channel-code="searchInfo.cid"/>
        </template>
      </CenterCard>
    </div>
  </div>
</template>

<script setup>
import CenterCard from './dataCenterComponents/centerCard.vue'
import ChainRatio from './dataCenterComponents/chainRatio.vue'
import ReclaimMileage from './dataCenterComponents/ReclaimMileage.vue'
import RecoveryRate from './dataCenterComponents/RecoveryRate.vue'
import Order from './dataCenterComponents/order.vue'
import lineCharts from './dataCenterComponents/lineCharts.vue'
import part from './dataCenterComponents/part.vue'
import {reactive, ref, nextTick, defineEmits} from "vue";
import {getChannelProductSelf} from "@/api/channelProduct";
import {getChannelAccountList} from "@/api/channelAccount";

const searchInfo = ref({})
const elSearchFormRef = ref()

const setOptions = async () =>{
  channelCodeOptions.value = []
  setChannelCodeOptions(vcpTableData.value, channelCodeOptions.value, false)
}

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

const isUpd = ref(false)
// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    getTableData()
  })
}

const searchRule = reactive({})

//通道产品
const channelCodeOptions = ref([])
const vcpTableData = ref([])
//账号明细
const accTableData = ref([])

const channelCodeProps = {
  expandTrigger: 'hover',
  checkStrictly: false,
  emitPath: false,
}

const handleChange = (value) => {
  
  console.log('cid:',value)
  searchInfo.value.cid = value
  // getTableData()
}

const setChannelCodeOptions = (ChannelCodeData, optionsData, disabled) => {
  ChannelCodeData &&
  ChannelCodeData.forEach(item => {
    if (item.children && item.children.length) {
      const option = {
        value: item.channelCode + '',
        label: item.productName,
        children: []
      }
      setChannelCodeOptions(
          item.children,
          option.children,
      )
      optionsData.push(option)
    } else {
      const option = {
        value: item.channelCode + '',
        label: item.productName,
      }
      optionsData.push(option)
    }
  })
}

const accOn = ref()
const accOff = ref()
const accTotal = ref()

const setAccSwitchView = async () => {
  const re = accTableData.value
  let countOn = 0
  let countOff = 0
  for (let i = 0; i < re.length; i++) {
    if (re[i].status === 1){
      countOn++;
    }else {
      countOff++;
    }
  }

  accOn.value = countOn
  accOff.value = countOff
  accTotal.value = re.length
  console.log(re.length)

}

const getTableData = async() => {
  await nextTick()
  const vcpTable = await getChannelProductSelf({ page: 1, pageSize: 9999, ...searchInfo.value })
  const table = await getChannelAccountList({ page: 1, pageSize: 9999, ...searchInfo.value })
  if (table.code === 0) {
    accTableData.value = table.data.list
    await setAccSwitchView()
  }
  if (vcpTable.code === 0) {
    vcpTableData.value = vcpTable.data.list
    setOptions()
  }
}

getTableData()

</script>

<script>
export default {
  name: 'DataCenter',
}
</script>

<style lang="scss" scoped>

.data-center-box{
  width: 100%;
  display: grid;
  grid-template-columns: 2fr 4fr;
  column-gap: 10px;
}
</style>

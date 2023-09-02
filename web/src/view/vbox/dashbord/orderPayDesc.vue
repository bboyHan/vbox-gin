<template>
 
 <el-descriptions
    class="margin-top"
    title="核销信息"
    :column="3"
    :size="size"
    border
  >
    <template #extra>
      <!-- <el-button type="primary">Operation</el-button> -->
    </template>
    <el-descriptions-item>
      <template #label>
        <div class="cell-item">
          <el-icon :style="iconStyle">
            <user />
          </el-icon>
          账户名
        </div>
      </template>
      kooriookami
    </el-descriptions-item>
    <el-descriptions-item>
      <template #label>
        <div class="cell-item">
          <el-icon :style="iconStyle">
            <iphone />
          </el-icon>
          Telephone
        </div>
      </template>
      18100000000
    </el-descriptions-item>
    <el-descriptions-item>
      <template #label>
        <div class="cell-item">
          <el-icon :style="iconStyle">
            <location />
          </el-icon>
          Place
        </div>
      </template>
      Suzhou
    </el-descriptions-item>
    <el-descriptions-item>
      <template #label>
        <div class="cell-item">
          <el-icon :style="iconStyle">
            <tickets />
          </el-icon>
          Remarks
        </div>
      </template>
      <el-tag size="small">School</el-tag>
    </el-descriptions-item>
    <el-descriptions-item>
      <template #label>
        <div class="cell-item">
          <el-icon :style="iconStyle">
            <office-building />
          </el-icon>
          Address
        </div>
      </template>
      No.1188, Wuzhong Avenue, Wuzhong District, Suzhou, Jiangsu Province
    </el-descriptions-item>
  </el-descriptions>

</template>

<script>
export default {
  name: 'OrderPayDesc'
}
</script>

<script setup>
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive,watch ,computed} from 'vue'
import {
  Iphone,
  Location,
  OfficeBuilding,
  Tickets,
  User,
} from '@element-plus/icons-vue'


const size = ref('')
const iconStyle = computed(() => {
  const marginMap = {
    large: '8px',
    default: '6px',
    small: '4px',
  }
  return {
    marginRight: marginMap[size.value] || marginMap.default,
  }
})
const blockMargin = computed(() => {
  const marginMap = {
    large: '32px',
    default: '28px',
    small: '24px',
  }
  return {
    marginTop: marginMap[size.value] || marginMap.default,
  }
})

</script>

<style>

.el-descriptions {
  margin-top: 20px;
}
.cell-item {
  display: flex;
  align-items: center;
}
.margin-top {
  margin-top: 20px;
}
</style>
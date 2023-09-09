<template>
  <div class="page">

 
  <div class="gva-card-box">
    <div class="gva-card quick-entrance">
      <h3>当日交易实时概况</h3>
      <el-row>
        <el-col :span="6">
          <el-statistic  :value="usersOrderata.tInCome" >
            <template #title>
            <div style="font-size: 14px; align-items: center">今日收入</div>
            </template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>元(RMB)</span>
            <!-- <span class="green">
              元
              <el-icon>
                <CaretTop />
              </el-icon>
            </span> -->
          </div>
        </div>
        </el-col>
        <el-col :span="6">
          <el-statistic :value="usersOrderata.openChId">
            <template #title>
              <div style="font-size: 14px; display: inline-flex; align-items: center">
                通道概览
                <!-- <el-icon style="margin-left: 4px" :size="12">
                  <Male />
                </el-icon> -->
              </div>
            </template>
            <template #suffix>/{{ usersOrderata.chIdCnt }}</template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>开启个数/通道总数</span>
            <!-- <span class="green">
              
              <el-icon>
                <CaretTop />
              </el-icon>
            </span> -->
          </div>
        </div>
        </el-col>
        <el-col :span="6">
          <el-statistic  :value="usersOrderata.tOrderQuantify" >
            <template #title>
            <div style="font-size: 14px; align-items: center">今日订单量</div>
            </template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>单</span>
            <!-- <span class="green">
              {{ usersOrderata.tOkRate }} %)
            </span> -->
          </div>
        </div>
        </el-col>
        <el-col :span="6">
          <el-statistic  :value="usersOrderata.tOkOrderQuantify" >
            <template #title>
            <div style="font-size: 14px; align-items: center">今日成单量</div>
            </template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>单</span>
            <span class="green">
              (成功率 {{ usersOrderata.tOkRate }} %)
            </span>
          </div>
        </div>
        </el-col>
      </el-row>
    </div>
        
  </div>
  
  <div class="gva-card-box">
    <div class="gva-card quick-entrance">
      <h3>昨日交易概况</h3>
      <el-row>
        <el-col :span="6">
          <el-statistic  :value="usersOrderata.yInCome" >
            <template #title>
            <div style="font-size: 14px; align-items: center">昨日收入</div>
            </template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>元(RMB)</span>
            <!-- <span class="green">
              元
              <el-icon>
                <CaretTop />
              </el-icon>
            </span> -->
          </div>
        </div>
        </el-col>
        <el-col :span="6">
          <el-statistic :value="usersOrderata.openChId">
            <template #title>
              <div style="font-size: 14px; display: inline-flex; align-items: center">
                通道概览
                <!-- <el-icon style="margin-left: 4px" :size="12">
                  <Male />
                </el-icon> -->
              </div>
            </template>
            <template #suffix>/{{ usersOrderata.chIdCnt }}</template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>开启个数/通道总数</span>
            <!-- <span class="green">
              
              <el-icon>
                <CaretTop />
              </el-icon>
            </span> -->
          </div>
        </div>
        </el-col>
        <el-col :span="6">
          <el-statistic  :value="usersOrderata.yOrderQuantify" >
            <template #title>
            <div style="font-size: 14px; align-items: center">昨日订单量</div>
            </template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>单</span>
            <!-- <span class="green">
              {{ usersOrderata.tOkRate }} %)
            </span> -->
          </div>
        </div>
        </el-col>
        <el-col :span="6">
          <el-statistic  :value="usersOrderata.yOkOrderQuantify" >
            <template #title>
            <div style="font-size: 14px; align-items: center">昨日成单量</div>
            </template>
          </el-statistic>
          <div class="statistic-footer">
          <div class="footer-item">
            <span>单</span>
            <span class="green">
              (成功率 {{ usersOrderata.yOkRate }} %)
            </span>
          </div>
        </div>
        </el-col>
      </el-row>
    </div>
        
  </div>
  <div class="gva-card-box">

    <div class="gva-card quick-entrance">
      <div>图形</div>
    </div>
  </div>
</div>


  </template>

<script>
export default {
  name: 'homePageDesc'
}
</script>

<script setup>
import { ChatLineRound, Male } from '@element-plus/icons-vue'
import {
  getVboxUserPayOrderAnalysis,
  getHomePagePayOrderAnalysis,
  getVboxUserPayOrderAnalysisIncomeCharts,
  getSelectPayOrderAnalysisQuantifyCharts,
  getSelectPayOrderAnalysisChannelIncomeCharts,
  getSelectPayOrderAnalysisIncomeBarCharts
} from '@/api/vboxPayOrder'
import { ref} from 'vue'


const usersOrderata = ref({
          username: '',
          balance:0,
          chIdCnt:0,
          openChId:0,
          yOrderQuantify:0,
          yOkOrderQuantify:0,
          yOkRate:0,
          yInCome:0,
          tOrderQuantify:0,
          tOkOrderQuantify:0,
          tOkRate:0,
          tInCome:0,
        })

const getUsersOrderata = async() => {
  const table = await getHomePagePayOrderAnalysis()
  if (table.code === 0) {
    usersOrderata.value = table.data.resultData
    console.log("usersOrderata= " + JSON.stringify(usersOrderata.value))
  }
}
getUsersOrderata()
</script>

<style lang="scss" scoped>
.el-col {
  text-align: center;
}

.page {
    @apply p-0;
    .gva-card-box{
      @apply p-4;
      &+.gva-card-box{
        @apply pt-0;
      }
    }
    .gva-card {
      @apply box-border bg-white rounded h-auto px-6 py-8 overflow-hidden shadow-sm;
      .gva-card-title{
        @apply pb-5 border-t-0 border-l-0 border-r-0 border-b border-solid border-gray-100;
      }
    }
    .gva-top-card {
        @apply h-72 flex items-center justify-between text-gray-500;
        &-left {
          @apply h-full flex flex-col w-auto;
            &-title {
              @apply text-3xl text-gray-600;
            }
            &-dot {
              @apply mt-4 text-gray-600 text-lg;
            }
            &-item{
              +.gva-top-card-left-item{
                margin-top: 24px;
              }
              margin-top: 14px;
            }
        }
        &-right {
            height: 600px;
            width: 600px;
            margin-top: 28px;
        }
    }
     ::v-deep(.el-card__header){
          @apply p-0  border-gray-200;
        }
        .card-header{
          @apply pb-5 border-b border-solid border-gray-200 border-t-0 border-l-0 border-r-0;
        }
    .quick-entrance-items {
      @apply flex items-center justify-center text-center text-gray-800;
        .quick-entrance-item {
          @apply px-8 py-6 flex items-center flex-col transition-all duration-100 ease-in-out rounded-lg cursor-pointer;
          &:hover{
            @apply shadow-lg;
          }
            &-icon {
              @apply flex items-center h-16 w-16 rounded-lg justify-center mx-0 my-auto text-2xl;
            }
            p {
                @apply mt-2.5;
            }
        }
    }
}

.statistic-footer {
  // display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  font-size: 12px;
  color: var(--el-text-color-regular);
  margin-top: 16px;
}

.statistic-footer .footer-item {
  // display: flex;
  justify-content: space-between;
  align-items: center;
}

.statistic-footer .footer-item span:last-child {
  display: inline-flex;
  align-items: center;
  margin-left: 4px;
}

.custom-title {
  font-size: 12px; /* Adjust the font size as needed */
}

.green {
  color: var(--el-color-success);
}
.red {
  color: var(--el-color-error);
}
</style>
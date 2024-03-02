<template>
  <div>
    <div v-if="isMobile">手机页面内容</div>
    <div v-else>
      <div>
        <div class="gva-search-box">
          <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" :rules="searchRule" @keyup.enter="onSubmit"
                   label-width="auto" label-position="right">
            <el-form-item label="付方单号" prop="orderId">
              <el-input v-model="searchInfo.orderId" placeholder="搜索付方单号"/>
            </el-form-item>
            <el-form-item label="付方ID" prop="pAccount">
              <el-input v-model="searchInfo.pAccount" placeholder="搜索付方ID"/>
            </el-form-item>
            <el-form-item label="订单状态" prop="orderStatus">
              <el-select v-model="searchInfo.orderStatus" placeholder="选择状态" style="width: 100px">
                <el-option label="已支付" value="1"/>
                <el-option label="未支付" value="2"/>
                <el-option label="超时" value="3"/>
                <el-option label="失败" value="-1"/>
              </el-select>
            </el-form-item>
            <el-form-item label="回调状态" prop="cbStatus">
              <el-select v-model="searchInfo.cbStatus" placeholder="选择状态" style="width: 100px">
                <el-option label="已回调" value="1"/>
                <el-option label="未回调" value="2"/>
              </el-select>
            </el-form-item>
            <el-form-item label="补单状态" prop="handStatus">
              <el-select v-model="searchInfo.handStatus" placeholder="选择状态" style="width: 100px">
                <el-option label="已补单" value="1"/>
                <el-option label="默认" value="2"/>
              </el-select>
            </el-form-item>
            <el-form-item label="通道账号" prop="acAccount">
              <el-input v-model="searchInfo.acAccount" placeholder="搜索通道账号"/>
            </el-form-item>
            <el-form-item label="通道ID" prop="channelCode">
              <el-input v-model="searchInfo.channelCode" placeholder="搜索通道ID"/>
            </el-form-item>
            <el-form-item>
              <el-button icon="refresh" @click="onReset"></el-button>
              <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
              <el-button icon="refresh" @click="resetSimple(true)">简约版</el-button>
              <el-button icon="refresh" @click="resetSimple(false)">详情版</el-button>
            </el-form-item>
          </el-form>
        </div>
        <div class="gva-table-box">
          <!--   简约版   -->
          <el-table
              v-if="isSimple"
              ref="multipleTable"
              style="width: 100%"
              tooltip-effect="dark"
              :data="tableData"
              row-key="ID"
              border :table-layout="'fixed'"
          >
            <el-table-column align="center" label="通道ID" prop="channelCode" width="70"/>
            <el-table-column align="center" label="充值账号" prop="acAccount" width="200">
              <template #default="scope">
                <div v-if="isPendingAcc(scope.row)">
                  <el-button type="info" link @click="getAccDetails(scope.row)">
                    {{ scope.row.acAccount }}
                  </el-button>
                  <el-button type="primary" link @click="openOrderHisShow(scope.row)">
                    <el-icon style="margin-right: 1px">
                      <Search/>
                    </el-icon>
                  </el-button>
                </div>
                <div v-else-if="!isPendingAcc(scope.row) && scope.row.orderStatus === 0">
                  <el-button type="info" link class="table-button">库存不足</el-button>
                </div>
                <div v-else>
                  <el-button type="info" :loading-icon="Eleme" loading link class="table-button">匹配中</el-button>
                </div>
              </template>
            </el-table-column>
            <el-table-column align="center" label="订单ID" prop="orderId" width="260"/>
            <el-table-column align="center" label="金额" prop="money" width="120"/>
            <el-table-column align="center" label="订单状态" prop="orderStatus" width="120">
              <template #default="scope">
                <el-button style="width: 90px"
                           :color="formatPayedColor(scope.row.orderStatus, scope.row.acId, scope.row.platId)"
                           @click="openSubmitCard(scope.row)">
                  {{ formatPayed(scope.row.orderStatus, scope.row.acId, scope.row.platId) }}
                </el-button>
              </template>
            </el-table-column>
            <el-table-column align="center" label="回调状态" prop="cbStatus" width="120">
              <template #default="scope">
                <el-button style="width: 90px" :color="formatNotifyColor(scope.row.cbStatus)">
                  {{ formatNotify(scope.row.cbStatus) }}
                </el-button>
              </template>
            </el-table-column>
            <el-table-column align="center" label="创建时间" width="180">
              <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
            </el-table-column>
            <el-table-column align="left" label="操作" width="260">
              <template #default="scope">
                <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
                  <el-icon style="margin-right: 5px">
                    <InfoFilled/>
                  </el-icon>
                  详情
                </el-button>
                <el-button type="primary" link class="table-button" @click="notifyPayOrder(scope.row)">
                  <el-icon style="margin-right: 5px">
                    <Position/>
                  </el-icon>
                  补单
                </el-button>
                <el-button v-if="Number(scope.row.channelCode) === 1101" type="primary" link class="table-button"
                           @click="openSubmitCard(scope.row)">
                  <el-icon style="margin-right: 5px">
                    <Notification/>
                  </el-icon>
                  核对
                </el-button>
              </template>
            </el-table-column>
            <!--          <el-table-column align="center" label="平台ID" prop="platId" width="260"/>-->
          </el-table>

          <!--   详情版   -->
          <el-table v-else ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID"
                    border>
            <el-table-column align="center" label="通道编码" prop="channelCode" width="100"/>
            <el-table-column align="center" label="充值账号" prop="acAccount" width="200">
              <template #default="scope">
                <div v-if="isPendingAcc(scope.row)">
                  <el-button type="info" link @click="getAccDetails(scope.row)">
                    {{ scope.row.acAccount }}
                  </el-button>
                  <el-button type="primary" link @click="openOrderHisShow(scope.row)">
                    <el-icon style="margin-right: 1px">
                      <Search/>
                    </el-icon>
                  </el-button>
                </div>
                <div v-else-if="!isPendingAcc(scope.row) && scope.row.orderStatus === 0">
                  <el-button type="info" link class="table-button">库存不足</el-button>
                </div>
                <div v-else>
                  <el-button type="info" :loading-icon="Eleme" loading link class="table-button">匹配中</el-button>
                </div>
              </template>
            </el-table-column>
            <el-table-column align="center" label="订单ID" prop="orderId" width="260"/>
            <el-table-column align="center" label="金额" prop="money" width="80"/>
            <el-table-column align="center" label="订单状态" prop="orderStatus" width="120">
              <template #default="scope">
                <el-button style="width: 90px"
                           :color="formatPayedColor(scope.row.orderStatus, scope.row.acId, scope.row.platId)"
                           @click="openSubmitCard(scope.row)">
                  {{ formatPayed(scope.row.orderStatus, scope.row.acId, scope.row.platId) }}
                </el-button>
              </template>
            </el-table-column>
            <el-table-column align="center" label="回调状态" prop="cbStatus" width="120">
              <template #default="scope">
                <el-button style="width: 90px" :color="formatNotifyColor(scope.row.cbStatus)">
                  {{ formatNotify(scope.row.cbStatus) }}
                </el-button>
              </template>
            </el-table-column>
            <el-table-column align="center" label="创建时间" width="180">
              <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
            </el-table-column>
            <el-table-column align="center" label="回调时间" width="180">
              <template #default="scope">{{ formatDate(scope.row.cbTime) }}</template>
            </el-table-column>
            <el-table-column align="center" label="操作" width="180">
              <template #default="scope">
                <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
                  <el-icon style="margin-right: 5px">
                    <InfoFilled/>
                  </el-icon>
                  详情
                </el-button>
                <el-button type="primary" link class="table-button" @click="notifyPayOrder(scope.row)">
                  <el-icon style="margin-right: 5px">
                    <Position/>
                  </el-icon>
                  补单
                </el-button>
              </template>
            </el-table-column>
            <el-table-column align="center" label="补单状态" prop="handStatus" width="120">
              <template #default="scope">
                <el-button style="width: 90px" :color="formatHandNotifyColor(scope.row.handStatus)">
                  {{ formatHandNotify(scope.row.handStatus) }}
                </el-button>
              </template>
            </el-table-column>
            <el-table-column align="center" label="过期时间" width="180">
              <template #default="scope">{{ formatDate(scope.row.expTime) }}</template>
            </el-table-column>
            <el-table-column align="center" label="付方ID" prop="pAccount" width="160">
              <template #default="scope">
                {{ scope.row.pAccount }}
              </template>
            </el-table-column>
            <el-table-column align="center" label="单价积分" prop="unitPrice" width="120"/>
            <el-table-column align="center" label="平台ID" prop="platId" width="500"/>
            <el-table-column align="center" label="访客ip" prop="payIp" width="180"/>
            <el-table-column align="center" label="区域" prop="payRegion" width="240"/>
            <el-table-column align="center" label="客户端设备" prop="payDevice" width="120"/>
          </el-table>
          <div class="gva-pagination">
            <el-pagination
                layout="total, sizes, prev, pager, next, jumper"
                :current-page="page"
                :page-size="pageSize"
                :page-sizes="[10, 30, 50, 100]"
                :total="total"
                @current-change="handleCurrentChange"
                @size-change="handleSizeChange"
            />
          </div>
        </div>

        <!-- 订单查看详情 -->
        <el-dialog v-model="detailShow" style="width: 70%" :draggable="true" lock-scroll
                   :before-close="closeDetailShow"
                   title="查看详情"
                   destroy-on-close>
          <el-row :gutter="24">
            <el-col :span="14">
              <el-scrollbar height="550px">
                <el-descriptions :column="6" border>
                  <el-descriptions-item label="订单ID" :span="6">{{ formData.orderId }}</el-descriptions-item>
                  <el-descriptions-item label="平台ID" :span="6">{{ formData.platId }}</el-descriptions-item>
                  <el-descriptions-item label="付方ID" :span="6">{{ formData.pAccount }}</el-descriptions-item>
                  <el-descriptions-item label="账号ID" :span="3">{{ formData.acId }}</el-descriptions-item>
                  <el-descriptions-item label="通道账号" :span="3">{{ formData.acAccount }}</el-descriptions-item>
                  <el-descriptions-item label="金额" :span="3">{{ formData.money }}</el-descriptions-item>
                  <el-descriptions-item label="单价积分" :span="3">{{ formData.unitPrice }}</el-descriptions-item>
                  <el-descriptions-item label="通道编码" :span="3">{{ formData.channelCode }}</el-descriptions-item>
                  <el-descriptions-item label="平台id" :span="6">{{ formData.platId }}</el-descriptions-item>
                  <el-descriptions-item label="客户ip" :span="3">{{ formData.payIp }}</el-descriptions-item>
                  <el-descriptions-item label="客户端设备" :span="3">{{ formData.payDevice }}</el-descriptions-item>
                  <el-descriptions-item label="区域" :span="6">{{ formData.payRegion }}</el-descriptions-item>
                  <el-descriptions-item label="商铺ID" :span="3">{{
                      formData.ext.shop.productId
                    }}
                  </el-descriptions-item>
                  <el-descriptions-item label="商铺备注" :span="3">{{
                      formData.ext.shop.shopRemark
                    }}
                  </el-descriptions-item>
                  <el-descriptions-item label="商铺成率" :span="6">
                    <div>近一小时：成单 / 总成单 - {{ formData.ext.dv.x2 }} / {{ formData.ext.dv.x1 }}
                      成率：{{ calculatePercentage(formData.ext.dv.x2, formData.ext.dv.x1) }}%
                    </div>
                    <div>今日：成单 / 总成单 - {{ formData.ext.dv.x4 }} / {{ formData.ext.dv.x3 }}
                      成率：{{ calculatePercentage(formData.ext.dv.x4, formData.ext.dv.x3) }}%
                    </div>
                  </el-descriptions-item>
                  <el-descriptions-item label="订单状态" :span="3">{{
                      formatPayed(formData.orderStatus)
                    }}
                  </el-descriptions-item>
                  <el-descriptions-item label="回调状态" :span="3">{{
                      formatNotify(formData.cbStatus)
                    }}
                  </el-descriptions-item>
                  <el-descriptions-item label="回调时间" :span="3">{{
                      formatDate(formData.cbTime)
                    }}
                  </el-descriptions-item>
                  <el-descriptions-item label="过期时间" :span="3">{{
                      formatDate(formData.expTime)
                    }}
                  </el-descriptions-item>
                </el-descriptions>
              </el-scrollbar>
            </el-col>
            <el-col :span="10">
              <el-scrollbar height="550px">
                <el-timeline style="max-width: 600px">
                  <el-timeline-item
                      v-for="(activity, index) in stepData"
                      :key="index"
                      :timestamp="activity.CreatedAt"
                  >
                    {{ activity.resp }}
                  </el-timeline-item>
                </el-timeline>
              </el-scrollbar>
            </el-col>
          </el-row>

        </el-dialog>

        <!--  补单  -->
        <el-dialog
            v-model="dialogFormVisible"
            :before-close="closeDialog"
            :title="typeTitle" :draggable="true"
            destroy-on-close
            style="width: 450px"
        >
          <el-scrollbar height="100px">
            <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="120px">
              <el-form-item label="订单ID" prop="orderId">
                <el-input disabled v-model="formData.orderId" :clearable="true" placeholder="请输入"
                          style="width: 80%"/>
              </el-form-item>
              <el-form-item label="谷歌动态验证" prop="authCaptcha">
                <el-input v-model="formData.authCaptcha" :clearable="true" placeholder="请输入谷歌动态验证"
                          style="width: 80%"/>
              </el-form-item>
            </el-form>
          </el-scrollbar>
          <template #footer>
            <div class="dialog-footer">
              <el-button @click="closeDialog">取 消</el-button>
              <el-button type="primary" @click="enterDialog">确 定</el-button>
            </div>
          </template>
        </el-dialog>

        <!-- 账号查看详情 -->
        <el-dialog v-model="detailAccShow" style="width: 800px" lock-scroll :before-close="closeAccDetailShow"
                   :draggable="true"
                   title="查看详情" destroy-on-close>
          <el-scrollbar height="550px">
            <el-descriptions :column="6" border>
              <el-descriptions-item label="用户归属" :span="6">{{ formAccData.username }}</el-descriptions-item>
              <el-descriptions-item label="账户ID" :span="6">{{ formAccData.acId }}</el-descriptions-item>
              <el-descriptions-item label="账户备注" :span="6">{{ formAccData.acRemark }}</el-descriptions-item>
              <el-descriptions-item label="通道账户" :span="3">{{ formAccData.acAccount }}</el-descriptions-item>
              <el-descriptions-item label="账户密码" :span="3">{{ formAccData.acPwd }}</el-descriptions-item>
              <el-descriptions-item label="ck" :span="6">
                <el-input v-model="formData.token" type="textarea" readonly/>
              </el-descriptions-item>
              <el-descriptions-item label="通道id" :span="6">{{ formAccData.cid }}</el-descriptions-item>
              <el-descriptions-item label="笔数限制" :span="2">{{ formAccData.countLimit }}</el-descriptions-item>
              <el-descriptions-item label="日限额" :span="2">{{ formAccData.dailyLimit }}</el-descriptions-item>
              <el-descriptions-item label="总限额" :span="2">{{ formAccData.totalLimit }}</el-descriptions-item>
              <el-descriptions-item label="状态开关" :span="3">{{
                  formAccData.status === 0 ? '关闭' : '开启'
                }}
              </el-descriptions-item>
              <el-descriptions-item label="系统开关" :span="3">{{
                  formAccData.sysStatus === 0 ? '关闭' : '开启'
                }}
              </el-descriptions-item>
            </el-descriptions>
          </el-scrollbar>
        </el-dialog>

        <!-- 指定账户官方充值详情 -->
        <el-dialog v-model="orderHisVisible" style="width: 1100px" lock-scroll :before-close="closeOrderHisShow"
                   :draggable="true"
                   title="查看充值详情" destroy-on-close>
          <el-scrollbar height="550px">
            <el-table tooltip-effect="dark" :data="orderHisTableData" row-key="ID" style="width: 100%">
              <el-table-column align="left" label="充值类型" prop="ShowName" width="180"/>
              <el-table-column align="left" label="渠道" prop="PayChannel" width="100"/>
              <el-table-column align="left" label="上游订单" prop="SerialNo" width="380"/>
              <el-table-column align="left" label="充值账号" prop="ProvideID" width="120"/>
              <el-table-column align="left" label="金额" prop="PayAmt" width="120">
                <template #default="scope">
                  {{ Number(scope.row.PayAmt) / 100 }}
                </template>
              </el-table-column>
              <el-table-column align="left" label="充值时间" prop="PayTime" width="160">
                <template #default="scope">
                  {{ formatUtcTimestamp(scope.row.PayTime) }}
                </template>
              </el-table-column>
            </el-table>
          </el-scrollbar>
        </el-dialog>

        <!-- 查看充值详情 2000 -->
        <el-dialog v-model="orderHis2000Visible" style="width: 1100px" :draggable="true" lock-scroll
                   :before-close="closeOrderHis2000Show"
                   title="查看充值详情" destroy-on-close>
          <el-scrollbar height="550px">
            <el-descriptions :column="4" border style="background-color: #a5abb4">
              <el-descriptions-item label="名称">{{ orderHis2000Info.gameName }}</el-descriptions-item>
              <el-descriptions-item label="账号">{{ orderHis2000Info.account }}</el-descriptions-item>
              <el-descriptions-item label="区域">{{ orderHis2000Info.zoneName }}</el-descriptions-item>
              <el-descriptions-item label="积分">{{ orderHis2000Info.leftCoins }}</el-descriptions-item>
            </el-descriptions>
            <el-table tooltip-effect="dark" :data="orderHis2000List" row-key="ID" style="width: 100%">
              <el-table-column align="center" label="账号" prop="acAccount" width="220"/>
              <el-table-column align="center" label="订单ID" prop="orderId" width="230"/>
              <el-table-column align="center" label="金额" prop="money" width="90"/>
              <el-table-column align="center" label="首查积分" prop="hisBalance" width="90"/>
              <el-table-column align="center" label="首查时间" prop="nowTime" width="160">
                <template #default="scope">
                  {{ formatUtcTimestamp(scope.row.nowTime) }}
                </template>
              </el-table-column>
              <el-table-column align="center" label="核准积分" prop="nowBalance" width="90">
                <template #default="scope">
                  <div v-if="Number(scope.row.nowBalance) === 0">-</div>
                  <div v-else>{{ Number(scope.row.nowBalance) }}</div>
                </template>
              </el-table-column>
              <el-table-column align="center" label="核准时间" prop="checkTime" width="160">
                <template #default="scope">
                  {{ formatUtcTimestamp(scope.row.checkTime) }}
                </template>
              </el-table-column>
            </el-table>
          </el-scrollbar>
        </el-dialog>

        <!-- 查看充值详情 4000 -->
        <el-dialog v-model="orderHis4000Visible" style="width: 1100px" :draggable="true" lock-scroll
                   :before-close="closeOrderHis4000Show"
                   title="查看充值详情" destroy-on-close>
          <el-scrollbar height="550px">
            <el-table tooltip-effect="dark" :data="orderHis4000TableData" row-key="ID" style="width: 100%">
              <el-table-column align="left" label="充值类型" prop="payProductName" width="80"/>
              <el-table-column align="left" label="渠道" prop="appName" width="100"/>
              <el-table-column align="left" label="上游订单" prop="orderId" width="280"/>
              <el-table-column align="left" label="充值账号" prop="displayAccount" width="120"/>
              <el-table-column align="left" label="金额" prop="orderAmount" width="100">
                <template #default="scope">
                  {{ Number(scope.row.orderAmount) }}
                </template>
              </el-table-column>
              <el-table-column align="left" label="时间" prop="PayTime" width="160">
                <template #default="scope">
                  {{ formatUtcTimestamp(scope.row.timestampMs / 1000) }}
                </template>
              </el-table-column>
            </el-table>
          </el-scrollbar>
        </el-dialog>

        <!--  核对  -->
        <el-dialog v-model="dialogCardSubmitVisible" :before-close="closeSubmitCard" title="核对" :draggable="true"
                   destroy-on-close style="width: 850px">
          <el-row :gutter="10">
            <el-col :span="24">
              <el-descriptions column="6" border>
                <el-descriptions-item label="账户" :span="3">{{ cardSubmitForm.x1 }}</el-descriptions-item>
                <el-descriptions-item label="卡号" :span="3">{{ cardSubmitForm.x3 }}</el-descriptions-item>
                <el-descriptions-item label="金额" :span="3">{{ cardSubmitForm.x2 }}</el-descriptions-item>
                <el-descriptions-item label="密码" :span="3">{{ cardSubmitForm.x4 }}</el-descriptions-item>
                <el-descriptions-item label="充值页面" :span="6">
                  <el-button type="primary" @click="openWindowExt(extUrl)">新窗口打开</el-button>
                </el-descriptions-item>
              </el-descriptions>
            </el-col>
            <!--        <el-col :span="24">-->
            <!--          <div style="height: 1000px">-->
            <!--            <iframe :src="extUrl" width="100%" height="100%"></iframe>-->
            <!--          </div>-->
            <!--        </el-col>-->
          </el-row>
        </el-dialog>
      </div>
    </div>
  </div>
</template>

<script setup>
import {
  callback2Pa,
  findPayOrder,
  getPayOrderList
} from '@/api/payOrder'
import {
  findChannelAccount, queryAccOrderHis,
} from '@/api/channelAccount'
// 全量引入格式化工具 请按需保留
import {
  getDictFunc,
  formatDate,
  formatBoolean,
  filterDict,
  ReturnArrImg,
  onDownloadFile,
  formatNotify,
  formatPayed,
  formatPayedColor,
  formatNotifyColor,
  formatHandNotify,
  formatHandNotifyColor,
  formatUtcTimestamp, calculatePercentage, formatPayCodeStatus
} from '@/utils/format'
import {ElMessage, ElMessageBox} from 'element-plus'
import {ref, reactive, onMounted} from 'vue'
import {Eleme, InfoFilled, Notification, Position, Search} from "@element-plus/icons-vue";

defineOptions({
  name: 'PayOrder'
})

// 使用 ref 来存储设备类型
const isMobile = ref(false);

// 检测设备类型
const detectDeviceType = () => {
  const userAgent = navigator.userAgent.toLowerCase();
  // isMobile.value = /iphone|ipad|ipod|android|blackberry|mini|windows\sce|palm/i.test(userAgent);
};

// 在组件挂载后检测设备类型
onMounted(() => {
  detectDeviceType();
});

//页面简约切换
const isSimple = ref(true)

// 获取付方账号详情
const getPADetails = (paID) => {
  console.log("查当前详情付方账号：" + paID)
}

// 重置
const resetSimple = (status) => {
  isSimple.value = status
}

const stepData = ref([])

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  authCaptcha: '',
  orderId: '',
  pAccount: '',
  money: 0,
  unitPrice: 0,
  uid: 0,
  acId: '',
  channelCode: '',
  platId: '',
  payIp: '',
  payRegion: '',
  payDevice: '',
  resourceUrl: '',
  notifyUrl: '',
  orderStatus: 0,
  cbStatus: 0,
  handStatus: 0,
  codeUseStatus: false,
  asyncTime: new Date(),
  cbTime: new Date(),
})


// 验证规则
const rule = reactive({})

const searchRule = reactive({
  createdAt: [
    {
      validator: (rule, value, callback) => {
        if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
          callback(new Error('请填写结束日期'))
        } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
          callback(new Error('请填写开始日期'))
        } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
          callback(new Error('开始日期应当早于结束日期'))
        } else {
          callback()
        }
      }, trigger: 'change'
    }
  ],
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async (valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
    if (searchInfo.value.orderStatus === "") {
      searchInfo.value.orderStatus = null
    }
    if (searchInfo.value.cbStatus === "") {
      searchInfo.value.cbStatus = null
    }
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async () => {
  const table = await getPayOrderList({page: page.value, pageSize: pageSize.value, ...searchInfo.value})
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () => {
}

// 获取需要的字典 可能为空 按需保留
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')
const typeTitle = ref('')


// 弹窗控制标记
const dialogFormVisible = ref(false)


// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findPayOrder({ID: row.ID})
  if (res.code === 0) {
    formData.value = res.data.repayOrder
    stepData.value = res.data.repayOrder.ext.records
    for (let i = 0; i < stepData.value.length; i++) {
      stepData.value[i].CreatedAt = formatDate(stepData.value[i].CreatedAt)
    }
    openDetailShow()
  }
}

// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
    orderId: '',
    pAccount: '',
    money: 0,
    unitPrice: 0,
    uid: 0,
    acId: '',
    channelCode: '',
    platId: '',
    payIp: '',
    payRegion: '',
    payDevice: '',
    notifyUrl: '',
    orderStatus: false,
    cbStatus: false,
    cbTime: new Date(),
  }
}

// ---- 补单 ----
// 打开弹窗
const openDialog = () => {
  type.value = 'notify'
  dialogFormVisible.value = true
  typeTitle.value = '补单'
}

// 打开详情（补单使用）
const notifyPayOrder = async (row) => {
  // 打开弹窗
  const res = await findPayOrder({ID: row.ID})
  if (res.code === 0) {
    formData.value = res.data.repayOrder
    openDialog()
  }
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    authCaptcha: '',
    orderId: '',
    pAccount: '',
    money: 0,
    unitPrice: 0,
    uid: 0,
    acId: '',
    channelCode: '',
    platId: '',
    payIp: '',
    payRegion: '',
    payDevice: '',
    notifyUrl: '',
    orderStatus: false,
    cbStatus: false,
    cbTime: new Date(),
  }
}
// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    switch (type.value) {
      case 'notify':
        console.log(formData.value)
        let res = await callback2Pa(formData.value);
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: '回调成功'
          })
        }
        dialogFormVisible.value = false
    }
  })
}


// ------ 账户匹配 ------
const isPendingAcc = (row) => {
  console.log(row.adId)
  return row.acId !== "";
}
// ------ 账户匹配 ------

// ------ 账户详情 ------
const formAccData = ref({
  acId: '',
  acRemark: '',
  acAccount: '',
  acPwd: '',
  token: '',
  cid: '',
  countLimit: 0,
  dailyLimit: 0,
  totalLimit: 0,
  status: 0,
  sysStatus: 0,
  username: 0,
})
const detailAccShow = ref(false)
// 打开详情弹窗
const openAccDetailShow = () => {
  detailAccShow.value = true
}
// 打开详情
const getAccDetails = async (row) => {
  // 打开弹窗
  const res = await findChannelAccount({acId: row.acId})
  if (res.code === 0) {
    formAccData.value = res.data.revca
    openAccDetailShow()
  }
}
// 关闭详情弹窗
const closeAccDetailShow = () => {
  detailAccShow.value = false
  formAccData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    uid: 0,
  }
}
// ------ 账户详情 ------

// ------ 账户充值记录 ---------
const orderHisVisible = ref(false)
const orderHis2000Visible = ref(false)
const orderHis4000Visible = ref(false)
const orderHisTableData = ref([])
const orderHis4000TableData = ref([])
const orderHis2000Info = ref([])
const orderHis2000List = ref([])
const openOrderHisShow = async (row) => {
  // orderHisVisible.value = true
  let req = {...row}
  console.log(req)
  let cid = req.channelCode
  if (cid >= 1000 && cid <= 1099) {
    orderHisVisible.value = true;
  } else if (cid >= 1100 && cid <= 1199) {
    orderHisVisible.value = true;
  } else if (cid >= 1200 && cid <= 1299) {
    orderHisVisible.value = true;
  } else if (cid >= 3000 && cid <= 3099) {
    orderHisVisible.value = true;
  } else if (cid >= 2000 && cid <= 2099) {
    orderHis2000Visible.value = true;
  } else if (cid >= 4000 && cid <= 4099) {
    orderHis4000Visible.value = true;
  }
  await queryAccOrderHisFunc(req)
}
const closeOrderHisShow = () => {
  orderHisVisible.value = false
  orderHisTableData.value = []
}
const closeOrderHis2000Show = () => {
  orderHis2000Visible.value = false
  orderHis2000List.value = []
  orderHis2000Info.value = {}
}
const closeOrderHis4000Show = () => {
  orderHis4000Visible.value = false
  orderHis4000TableData.value = []
}
const queryAccOrderHisFunc = async (row) => {
  const req = {...row}
  const resAcc = await findChannelAccount({acId: req.acId})

  let res = await queryAccOrderHis(resAcc.data.revca)
  console.log(res.data)
  let cid = resAcc.data.revca.cid
  if (cid >= 3000 && cid <= 3099) {
    if (res.code === 0) {
      orderHisTableData.value = res.data.list.WaterList
    }
  } else if (cid >= 1000 && cid <= 1099) {
    if (res.code === 0) {
      orderHisTableData.value = res.data.list.WaterList
    }
  } else if (cid >= 1100 && cid <= 1199) {
    if (res.code === 0) {
      orderHisTableData.value = res.data.list.WaterList
    }
  } else if (cid >= 4000 && cid <= 4099) {
    if (res.code === 0) {
      orderHis4000TableData.value = res.data.list
    }
  } else if (cid >= 2000 && cid <= 2099) {
    if (res.code === 0) {
      orderHis2000Info.value = res.data.list.info
      orderHis2000List.value = res.data.list.list
    }
  }
}
// ------ 账户充值记录 ---------

// ------ 打开三方页面 ---------
const dialogCardSubmitVisible = ref(false)
const cardSubmitForm = ref({x1: '', x2: '', x3: '', x4: ''})
const extUrl = ref()

const openWindowExt = (extUrl) => {
  window.open(extUrl, '_blank')
}

const openSubmitCard = (row) => {
  cardSubmitForm.value = row
  console.log("row", row)
  let cid = Number(row.channelCode);
  let money = Number(row.money);
  if (cid === 1101) {
    console.log("露一手")
    dialogCardSubmitVisible.value = true;
    let c = row.platId.split("_");
    cardSubmitForm.value.x1 = row.acAccount
    cardSubmitForm.value.x2 = money
    cardSubmitForm.value.x3 = c[0]
    cardSubmitForm.value.x4 = c[1]

    switch (money) {
      case 10:
        extUrl.value = "https://www.junka.com/Official/JNetMapSup/CommonCharge.aspx?category=AAJSUPHNZG&product=AAJSUPHNZG010CZ"
        break
      case 20:
        extUrl.value = "https://www.junka.com/Official/JNetMapSup/CommonCharge.aspx?category=AAJSUPHNZG&product=AAJSUPHNZG020CZ"
        break
      case 30:
        extUrl.value = "https://www.junka.com/Official/JNetMapSup/CommonCharge.aspx?category=AAJSUPHNZG&product=AAJSUPHNZG030CZ"
        break
      case 50:
        extUrl.value = "https://www.junka.com/Official/JNetMapSup/CommonCharge.aspx?category=AAJSUPHNZG&product=AAJSUPHNZG050CZ"
        break
      case 100:
        extUrl.value = "https://www.junka.com/Official/JNetMapSup/CommonCharge.aspx?category=AAJSUPHNZG&product=AAJSUPHNZG100CZ"
        break
      case 200:
        extUrl.value = "https://www.junka.com/Official/JNetMapSup/CommonCharge.aspx?category=AAJSUPHNZG&product=AAJSUPHNZG200CZ"
        break
      default:
        extUrl.value = `https://www.junka.com/Official/JNetMapSup/CommonCharge.aspx?category=AAJSUPHNZG`
    }
  }
}

const closeSubmitCard = () => {
  dialogCardSubmitVisible.value = false
  cardSubmitForm.value = {}
}
// ------ 打开三方页面 ---------
</script>

<style>
/*.el-table__body, .el-table__header {
  border-collapse: collapse;
  border-bottom-width: 2px;
  border-bottom-style: solid;
}

.el-table__body tr, .el-table__header tr {
  border-bottom-width: 2px;
  border-bottom-style: solid;
}

.el-table__body td, .el-table__header th {
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
}*/
</style>

<template>
  <div>
    <div class="gva-search-box">
      <el-row :gutter="12">
        <el-col :span="24">
          <div class="gva-btn-list">
            <ul>
              <li><span style="color: blue;">【注意】</span>同类型通道，原则上不允许创建相同账号，<span style="color: red;">如在同一时间段产生同金额的订单，会对多个通道中的相同账号，均认定为“支付成功”</span>，请谨慎操作!
              </li>
              <li><span style="color: blue;">【规避方法】</span>1、自行核对避免同类型通道创建相同账号；2、引导类商铺管理，对不同通道的商铺金额进行隔离管理（如1003开启10元金额，1004则不开启10元金额）；如产生影响，后台正常计分，请自行保障风险!
              </li>
              <li><span style="color: blue;">【系统开关】</span>当选择通道账号开启时，状态开关栏中显示的<span
                  style="color: blue;">"系统开关"</span>如未能正常开启，请先自行查看【<span
                  style="color: blue;text-decoration: underline;"><b>账号详情</b></span>】核查原因！
              </li>
              <li><span style="color: blue;">【限额说明】</span>针对添加通道账号的限额设置，包括<span style="color: blue;">"总额限制"、"日额限制"、"进单限制"、"拉单限制"</span>，系统计算规则：当前账号进单的累计金额或累计笔数<span
                  style="color: blue;">大于或等于限制时</span>，系统进行自动关号。<span
                  style="color: red;text-decoration: underline;"><b>非精准控制，请知悉！</b></span></li>
            </ul>
          </div>
        </el-col>

        <el-col :span="24">
          <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" :rules="searchRule" @keyup.enter="onSubmit"
                   label-width="auto" label-position="right">
            <el-form-item label="通道账户" prop="acAccount">
              <el-input v-model="searchInfo.acAccount" placeholder="搜索通道账户"/>
            </el-form-item>
            <el-form-item label="账户备注" prop="acRemark">
              <el-input v-model="searchInfo.acRemark" placeholder="搜索备注"/>
            </el-form-item>
            <el-form-item label="账户ID" prop="acId">
              <el-input v-model.number="searchInfo.acId" placeholder="搜索账户ID"/>
            </el-form-item>
            <el-form-item label="通道ID" prop="cid">
              <el-input v-model.number="searchInfo.cid" placeholder="搜索通道ID"/>
            </el-form-item>
            <el-form-item label="开关状态" prop="status">
              <el-select v-model="searchInfo.status" placeholder="选择状态" style="width: 120px">
                <el-option label="已开启" value="1"/>
                <el-option label="已关闭" value="0"/>
              </el-select>
            </el-form-item>
            <el-form-item label="系统状态" prop="sysStatus">
              <el-select v-model="searchInfo.sysStatus" placeholder="选择系统状态" style="width: 120px">
                <el-option label="已开启" value="1"/>
                <el-option label="已关闭" value="0"/>
              </el-select>
            </el-form-item>
            <el-form-item label="控制策略" prop="sysStatus">
              <el-select v-model="searchInfo.ctlStatus" placeholder="选择控制策略" style="width: 120px">
                <el-option label="模糊控制" value="1"/>
                <el-option label="精准控制" value="2"/>
              </el-select>
            </el-form-item>
            <el-form-item label="归属用户" prop="username">
              <el-input v-model.number="searchInfo.username" placeholder="搜索归属用户"/>
            </el-form-item>
            <el-form-item>
              <el-button icon="refresh" @click="onReset"></el-button>
              <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
            </el-form-item>
          </el-form>
        </el-col>
      </el-row>

    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
        <el-popover v-model:visible="deleteVisible" placement="top" width="240">
          <p><span style="color: red;">注意：删除后，通道账户将无法使用(预产类资源也将清除)</span>，确定要删除吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button type="primary" link @click="deleteVisible = false">取消</el-button>
            <el-button type="primary" @click="onDelete">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length"
                       @click="deleteVisible = true">删除
            </el-button>
          </template>
        </el-popover>
        <el-popover v-model:visible="switchOnVisible" placement="top" width="240">
          <p><span style="color: red;">注意：开启后，通道账户将进入待使用状态(预产类资源也将启用)</span>，确定批量开启吗？
          </p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button type="primary" link @click="switchOnVisible = false">取消</el-button>
            <el-button type="primary" @click="onSwitchEnable">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="turn-off" style="margin-left: 10px;" :disabled="!multipleSelection.length"
                       @click="switchOnVisible = true">批量开启
            </el-button>
          </template>
        </el-popover>
        <el-popover v-model:visible="switchOffVisible" placement="top" width="240">
          <p><span style="color: red;">注意：关闭后，通道账户将无法使用(预产类资源也将禁用)</span>，确定批量关闭吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button type="primary" link @click="switchOffVisible = false">取消</el-button>
            <el-button type="primary" @click="onSwitchDisable">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="open" style="margin-left: 10px;" :disabled="!multipleSelection.length"
                       @click="switchOffVisible = true">批量关闭
            </el-button>
          </template>
        </el-popover>
        <el-button type="primary" icon="wallet" @click="showCostOrderAcc('')">核算</el-button>
        <el-row :gutter="12">
          <span v-for="item in countItem">
            <el-col :span="12">
              <el-popover trigger="hover" placement="right-end" width="450">
                <el-row>
                  <el-col v-for="ele in item.list" :span="24">
                  <el-button>
                    <el-col :span="24">
                      <div>
                      【通道ID：{{ ele.cid }}】<el-icon class="is-loading" style="margin-right: 2px"><Loading/> </el-icon>
                    已开启 <span style="color: red"><b>{{ ele.total }} </b></span> 个
                      </div>
                    </el-col>
                  </el-button>
                    <div style="margin: 5px">
                      <span v-for="em in ele.list">
                        <span style="padding: 5px">{{ em.money }}元(<span
                            style="color: red;">{{ em.unused }}</span>个)</span>
                      </span>
                    </div>
                  </el-col>
                </el-row>
                <div style="text-align: right; margin-top: 8px;">
                </div>
                <template #reference>
                  <el-button icon="search">池标识【{{ item.orgId }}】</el-button>
                </template>
              </el-popover>
            </el-col>
          </span>
        </el-row>
      </div>

      <el-table ref="multipleTable" tooltip-effect="dark" :data="tableData" row-key="ID" border resizable="true"
                @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55"/>
        <el-table-column align="left" label="ID" prop="acId" width="140">
          <template #default="scope">
            {{ scope.row.acId }}
            <el-popover placement="top" width="300">
              <p><span style="color: blue;">【核算查询】</span>查询该账号近日订单统计情况!</p>
              <template #reference>
                <el-button type="primary" link icon="wallet-filled" class="table-button" style="margin-left: 5px"
                           @click="showCostOrderAcc(scope.row)"></el-button>
              </template>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column align="center" label="通道ID" prop="cid" width="80"/>
        <el-table-column align="center" label="账户备注" prop="acRemark" width="180"/>
        <el-table-column align="center" label="通道账户" prop="acAccount" width="180"/>
        <el-table-column align="left" label="账户密钥" prop="acPwd" width="120"/>
        <el-table-column align="left" label="CK" prop="token" width="260">
          <template #default="scope">
            <el-input v-model="scope.row.token" :rows="3" readonly="readonly">
              <template #append>
                <el-button type="primary" link icon="edit" @click="updateCaTokenFunc(scope.row)"></el-button>
              </template>
            </el-input>
          </template>
        </el-table-column>
        <el-table-column align="center" label="控制策略" prop="ctlStatus" width="100">
          <template #default="scope">
            <el-switch v-model="scope.row.ctlStatus" inline-prompt :active-value="2" active-text="精准"
                       :inactive-value="1" inactive-text="模糊" size="large"
                       @change="()=>{switchCtlEnable(scope.row)}"/>
          </template>
        </el-table-column>
        <el-table-column align="center" label="状态 / 系统开关" prop="status" width="200">
          <template #default="scope">
            <el-row :gutter="12">
              <el-col :span="8">
                <el-popover trigger="hover" placement="top" width="240">
                  <p><span style="color: red;">注意</span>：操作后将通过系统审核，<span style="color: red;">审核通过后开启（或关闭）账号关联资源；</span><span
                      style="color: blue;">未通过系统审核请查看"操作日志"</span>核查原因，确定操作？</p>
                  <template #reference>
                    <el-switch v-model="scope.row.status" inline-prompt :active-value="1" active-text="开启"
                               :inactive-value="0" inactive-text="关闭" size="large"
                               @change="()=>{switchEnable(scope.row)}"/>
                  </template>
                </el-popover>
              </el-col>
              <el-col :span="8">
                <el-switch v-model="scope.row.sysStatus" inline-prompt :active-value="1" active-text="开启"
                           :inactive-value="0" inactive-text="关闭" size="large" disabled/>
              </el-col>
              <el-col :span="8">
                <el-button v-if="scope.row.sysStatus === 1" type="primary" circle link icon="info-filled"
                           @click="getDetails(scope.row)">
                </el-button>
                <el-button v-else type="danger" size="small" circle icon="info-filled" color="red"
                           @click="getDetails(scope.row)">
                </el-button>
              </el-col>
            </el-row>
          </template>
        </el-table-column>
        <el-table-column align="center" label="冷却状态" prop="cbStatus" width="120">
          <template #default="scope">
            <el-button style="width: 80px" :color="formatCDStatusColor(scope.row.cdStatus)">
              {{ formatCDStatus(scope.row.cdStatus) }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column align="center" label="查询" width="120">
          <template #default="scope">
            <el-row>
              <el-col :span="12">
                <el-button type="primary" link class="table-button" @click="openOrderHisShow(scope.row)">记录
                </el-button>
              </el-col>
              <el-col :span="12">
                <el-button type="primary" link class="table-button" @click="openOrderSysShow(scope.row)">成单
                </el-button>
              </el-col>
            </el-row>
          </template>
        </el-table-column>
        <el-table-column align="center" label="日限额" prop="dailyLimit" width="90"/>
        <el-table-column align="center" label="总限额" prop="totalLimit" width="90"/>
        <el-table-column align="center" label="日限进单" prop="dlyCntLimit" width="90"/>
        <el-table-column align="center" label="总限进单" prop="inCntLimit" width="90"/>
        <el-table-column align="center" label="总限拉单" prop="countLimit" width="90"/>
        <el-table-column align="center" label="操作" width="160">
          <template #default="scope">
            <el-row :gutter="12">
              <el-col :span="16">
                <span v-if="Number(scope.row.cid) === 3000">
                  <el-row>
                    <el-col :span="24">
                      <el-button type="primary" link class="table-button"
                                 @click="createByChannelPayCodeFunc(scope.row)">
                        产码
                      </el-button>
                      <el-button type="primary" link icon="info-filled" class="table-button"
                                 @click="openPayCodeOverviewShow(scope.row)"></el-button>
                    </el-col>
                  </el-row>
                </span>
                <span v-else>
                  <span>非预产通道</span>
                </span>
              </el-col>
              <el-col :span="8">
                <el-popover placement="top" width="240">
                  <p><span style="color: blue;">【通道转移】</span>（用于同类型通道切换通道编码使用，不允许跨类型迁移）<span
                      style="color: red;">注意：转移后，将进行旧资源池清理工作</span>，请谨慎操作!</p>
                  <template #reference>
                    <el-button type="primary" link icon="chrome-filled" class="table-button"
                               @click="transferAccFunc(scope.row)"></el-button>
                  </template>
                </el-popover>
              </el-col>
            </el-row>
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作" min-width="120">
          <template #default="scope">
            <el-row>
              <el-col :span="24">
                <el-button type="primary" link icon="info-filled" class="table-button"
                           @click="getDetails(scope.row)"></el-button>
                <el-button type="primary" link icon="edit" class="table-button"
                           @click="updateChannelAccountFunc(scope.row)"></el-button>
                <el-button type="warning" link icon="delete" @click="deleteRow(scope.row)"></el-button>
              </el-col>
            </el-row>
          </template>
        </el-table-column>
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

    <el-dialog v-model="dialogChanFormVisible" :before-close="closeChanDialog" :title="typeTitle" destroy-on-close
               style="width: 85%;min-width: 1000px">
      <el-scrollbar height="700px">
        <div>
          <div>
            <el-row :gutter="24">
              <el-col :span="24" :xs="24">
                <div v-for="(parent, index) in parentNodes" :key="index" class="card-container">
                  <el-col :span="24" :xs="24">
                    <div class="flex flex-wrap items-center justify-between"
                         style="margin-left: 10px;margin-bottom: -30px"><h2>
                      {{ parent.productName }}</h2></div>
                    <el-divider></el-divider>
                  </el-col>
                  <el-row :gutter="12">
                    <div v-if="parent.children && parent.children.length > 0"
                         style="flex-wrap: wrap;  justify-content: center;display: flex;">
                      <div v-for="(node, childIndex) in parent.children" :key="childIndex">
                        <el-col class="card" :span="24" :xs="24">
                          <div @click="handleProdClick(node)">
                            <CenterCard title="" :custom-style="accCustomStyle">
                              <template #action>
                                <span class="gvaIcon-prompt" style="color: #999"></span>
                              </template>
                              <template #body>
                                <div class="acc-container">
                                  <div class="indicator">
                                  <span>
                                    <div class="label">编码</div>
                                    <div class="value">{{ node.channelCode }}</div>
                                  </span>
                                    <span>
                                    <div class="label">名称</div>
                                    <div class="value">{{ node.productName }}</div>
                                  </span>
                                  </div>
                                </div>
                              </template>
                            </CenterCard>
                          </div>
                        </el-col>
                      </div>
                    </div>
                    <div v-else>
                      <el-col class="card" :span="24" :xs="24">
                        <div @click="handleProdClick(parent)">
                          <CenterCard title="" :custom-style="accCustomStyle">
                            <template #action>
                              <span class="gvaIcon-prompt" style="color: #999"></span>
                            </template>
                            <template #body>
                              <div class="acc-container">
                                <div class="indicator">
                                <span>
                                  <div class="label">编码</div>
                                  <div class="value">{{ parent.channelCode }}</div>
                                </span>
                                  <span>
                                  <div class="label">名称</div>
                                  <div class="value">{{ parent.productName }}</div>
                                </span>
                                </div>
                              </div>
                            </template>
                          </CenterCard>
                        </div>
                      </el-col>
                    </div>
                  </el-row>
                </div>
              </el-col>
            </el-row>
          </div>
        </div>
      </el-scrollbar>
    </el-dialog>

    <!--  创建 1000 -->
    <el-dialog v-model="dialog1000FormVisible" :before-close="close1000Dialog" :draggable="true" :title="typeTitle"
               destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="产码方式" prop="type">
              <el-button v-model="formData.type" readonly disabled>
                {{ formatProdType(Number(formData.type)) }}
              </el-button>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通道ID" prop="cid" disabled>
              <el-input v-model="formData.cid" readonly disabled></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="通道账户" prop="acAccount">
              <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入账号或扫码自动获取账户"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="账户备注" prop="acRemark">
              <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入备注"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="扫码授权" prop="token">
          <el-row :gutter="12">
            <el-col :span="24">
              <div v-if="imageQrCode">
                <img :src="imageQrCode" alt="qr" style="width: 180px;height: 180px"/>
              </div>
              <div v-else>
                <div style="width: 180px;height: 180px;font-size: 12px;color:grey;border: 1px solid grey">
                  请点击"获取二维码"
                </div>
              </div>
            </el-col>
            <el-col :span="24">
              <el-button type="primary" link @click="loginQr">重新获取二维码</el-button>
            </el-col>

            <el-col :span="18">
              <el-input v-model="imageQrStatusMsg" placeholder="二维码状态" style="color: red !important;"/>
            </el-col>

            <el-col :span="6">
              <el-button type="primary" @click="getQrCookie" style="width: 100%">获取CK</el-button>
            </el-col>
            <el-col :span="24">
              <el-input v-model="formData.token" type="textarea" :clearable="true"
                        placeholder="请输入CK或点击【获取CK】自动获取"/>
            </el-col>
          </el-row>
        </el-form-item>
        <el-row>
          <el-col :span="1"></el-col>
          <el-col :span="22">
            <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
          </el-col>
          <el-col :span="1"></el-col>
          <el-col :span="12">
            <el-form-item label="日限额" prop="dailyLimit">
              <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="总限额" prop="totalLimit">
              <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="日限进单" prop="dlyCntLimit">
              <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限进单" prop="inCntLimit">
              <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限拉单" prop="countLimit">
              <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="状态开关" prop="status">
              <el-switch v-model="formData.status" active-value="1" inactive-value="0" active-text="开启"
                         inactive-text="关闭"></el-switch>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="控制策略" prop="ctlStatus">
              <el-switch v-model="formData.ctlStatus" active-value="2" inactive-value="1" active-text="精准控制"
                         inactive-text="模糊控制"></el-switch>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="returnPreStep(1000)">上一步</el-button>
          <el-button @click="close1000Dialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  创建 1100 -->
    <el-dialog v-model="dialog1100FormVisible" :before-close="close1100Dialog" :draggable="true" :title="typeTitle"
               destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="产码方式" prop="type" disabled>
              <el-button v-model="formData.type" readonly disabled>
                {{ formatProdType(Number(formData.type)) }}
              </el-button>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通道ID" prop="cid" disabled>
              <el-input v-model="formData.cid" readonly disabled></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="通道账户" prop="acAccount">
              <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入账号"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="账户备注" prop="acRemark">
              <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入备注"/>
            </el-form-item>
          </el-col>
        </el-row>
        <!--        <el-form-item label="报文" prop="token">
                  <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入CK"/>
                </el-form-item>-->
        <el-form-item label="扫码授权" prop="token">
          <el-row :gutter="12">
            <el-col :span="24">
              <div v-if="imageQrCode">
                <img :src="imageQrCode" alt="qr" style="width: 180px;height: 180px"/>
              </div>
              <div v-else>
                <div style="width: 180px;height: 180px;font-size: 12px;color:grey;border: 1px solid grey">
                  请点击"获取二维码"
                </div>
              </div>
            </el-col>
            <el-col :span="24">
              <el-button type="primary" link @click="loginQr">重新获取二维码</el-button>
            </el-col>

            <el-col :span="18">
              <el-input v-model="imageQrStatusMsg" placeholder="二维码状态" style="color: red !important;"/>
            </el-col>

            <el-col :span="6">
              <el-button type="primary" @click="getQrCookie" style="width: 100%">获取CK</el-button>
            </el-col>
            <el-col :span="24">
              <el-input v-model="formData.token" type="textarea" :clearable="true"
                        placeholder="请输入CK或点击【获取CK】自动获取"/>
            </el-col>
          </el-row>
        </el-form-item>
        <el-row>
          <el-col :span="1"></el-col>
          <el-col :span="22">
            <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
          </el-col>
          <el-col :span="1"></el-col>
          <el-col :span="12">
            <el-form-item label="日限额" prop="dailyLimit">
              <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="总限额" prop="totalLimit">
              <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="日限进单" prop="dlyCntLimit">
              <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限进单" prop="inCntLimit">
              <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限拉单" prop="countLimit">
              <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="状态开关" prop="status">
              <el-switch v-model="formData.status" active-value="1" inactive-value="0" active-text="开启"
                         inactive-text="关闭"></el-switch>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="控制策略" prop="ctlStatus">
              <el-switch v-model="formData.ctlStatus" active-value="2" inactive-value="1" active-text="精准控制"
                         inactive-text="模糊控制"></el-switch>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="returnPreStep(1100)">上一步</el-button>
          <el-button @click="close1100Dialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  创建 1200 -->
    <el-dialog v-model="dialog1200FormVisible" :before-close="close1200Dialog" :draggable="true" :title="typeTitle"
               destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="产码方式" prop="type">
              <el-button v-model="formData.type" readonly disabled>
                {{ formatProdType(Number(formData.type)) }}
              </el-button>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通道ID" prop="cid" disabled>
              <el-input v-model="formData.cid" readonly disabled></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="通道账户" prop="acAccount">
              <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入账号"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="账户备注" prop="acRemark">
              <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入备注"/>
            </el-form-item>
          </el-col>
        </el-row>
        <!--        <el-form-item label="报文" prop="token">
                  <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入CK"/>
                </el-form-item>-->
        <el-form-item label="扫码授权" prop="token">
          <el-row :gutter="12">
            <el-col :span="24">
              <div v-if="imageQrCode">
                <img :src="imageQrCode" alt="qr" style="width: 180px;height: 180px"/>
              </div>
              <div v-else>
                <div style="width: 180px;height: 180px;font-size: 12px;color:grey;border: 1px solid grey">
                  请点击"获取二维码"
                </div>
              </div>
            </el-col>
            <el-col :span="24">
              <el-button type="primary" link @click="loginQr">重新获取二维码</el-button>
            </el-col>

            <el-col :span="18">
              <el-input v-model="imageQrStatusMsg" placeholder="二维码状态" style="color: red !important;"/>
            </el-col>

            <el-col :span="6">
              <el-button type="primary" @click="getQrCookie" style="width: 100%">获取CK</el-button>
            </el-col>
            <el-col :span="24">
              <el-input v-model="formData.token" type="textarea" :clearable="true"
                        placeholder="请输入CK或点击【获取CK】自动获取"/>
            </el-col>
          </el-row>
        </el-form-item>
        <el-row>
          <el-col :span="1"></el-col>
          <el-col :span="22">
            <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
          </el-col>
          <el-col :span="1"></el-col>
          <el-col :span="12">
            <el-form-item label="日限额" prop="dailyLimit">
              <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="总限额" prop="totalLimit">
              <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="日限进单" prop="dlyCntLimit">
              <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限进单" prop="inCntLimit">
              <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限拉单" prop="countLimit">
              <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="状态开关" prop="status">
              <el-switch v-model="formData.status" active-value="1" inactive-value="0" active-text="开启"
                         inactive-text="关闭"></el-switch>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="控制策略" prop="ctlStatus">
              <el-switch v-model="formData.ctlStatus" active-value="2" inactive-value="1" active-text="精准控制"
                         inactive-text="模糊控制"></el-switch>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="returnPreStep(1200)">上一步</el-button>
          <el-button @click="close1200Dialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  创建 2000 -->
    <el-dialog v-model="dialog2000FormVisible" :before-close="close2000Dialog" :draggable="true" :title="typeTitle"
               destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule2000" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="产码方式" prop="type" disabled>
              <el-button v-model="formData.type" readonly disabled>
                {{ formatProdType(Number(formData.type)) }}
              </el-button>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通道ID" prop="cid" disabled>
              <el-input v-model="formData.cid" readonly disabled></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="24">
            <el-form-item label="账户备注" prop="acRemark">
              <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入备注"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="报文" prop="token">
          <el-input v-model="formData.token" type="textarea" :clearable="true"
                    placeholder="输入示例：https://security.seasungame.com/security_extend_server/helper/balance/queryBalance?gameCode=jx3&account=aa123123&accountType=&zoneCode=z22&SN=98710641126&remark=&sign=36A360706FD189A2BF867D70F61117BE"/>
        </el-form-item>
        <el-row>
          <el-col :span="1"></el-col>
          <el-col :span="22">
            <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
          </el-col>
          <el-col :span="1"></el-col>
          <el-col :span="12">
            <el-form-item label="日限额" prop="dailyLimit">
              <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="总限额" prop="totalLimit">
              <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="日限进单" prop="dlyCntLimit">
              <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限进单" prop="inCntLimit">
              <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限拉单" prop="countLimit">
              <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="状态开关" prop="status">
              <el-switch v-model="formData.status" active-value="1" inactive-value="0" active-text="开启"
                         inactive-text="关闭"></el-switch>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="控制策略" prop="ctlStatus">
              <el-switch v-model="formData.ctlStatus" active-value="2" inactive-value="1" active-text="精准控制"
                         inactive-text="模糊控制"></el-switch>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="returnPreStep(1200)">上一步</el-button>
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  创建 3000 -->
    <el-dialog v-model="dialog3000FormVisible" :before-close="close3000Dialog" :draggable="true" :title="typeTitle"
               destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="产码方式" prop="type" disabled>
              <el-button v-model="formData.type" readonly disabled>
                {{ formatProdType(Number(formData.type)) }}
              </el-button>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通道ID" prop="cid" disabled>
              <el-input v-model="formData.cid" readonly disabled></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="通道账户" prop="acAccount">
              <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="账户备注" prop="acRemark">
              <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="扫码授权" prop="token">
          <el-row :gutter="12">
            <el-col :span="24">
              <div v-if="imageQrCode">
                <img :src="imageQrCode" alt="qr" style="width: 180px;height: 180px"/>
              </div>
              <div v-else>
                <div style="width: 180px;height: 180px;font-size: 12px;color:grey;border: 1px solid grey">
                  请点击"获取二维码"
                </div>
              </div>
            </el-col>
            <el-col :span="24">
              <el-button type="primary" link @click="loginQr">重新获取二维码</el-button>
            </el-col>

            <el-col :span="18">
              <el-input v-model="imageQrStatusMsg" placeholder="二维码状态" style="color: red !important;"/>
            </el-col>

            <el-col :span="6">
              <el-button type="primary" @click="getQrCookie" style="width: 100%">获取CK</el-button>
            </el-col>
            <el-col :span="24">
              <el-input v-model="formData.token" type="textarea" :clearable="true"
                        placeholder="请输入CK或点击【获取CK】自动获取"/>
            </el-col>
          </el-row>
        </el-form-item>
        <el-row>
          <el-col :span="1"></el-col>
          <el-col :span="22">
            <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
          </el-col>
          <el-col :span="1"></el-col>
          <el-col :span="12">
            <el-form-item label="日限额" prop="dailyLimit">
              <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="总限额" prop="totalLimit">
              <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="日限进单" prop="dlyCntLimit">
              <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限进单" prop="inCntLimit">
              <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限拉单" prop="countLimit">
              <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="状态开关" prop="status">
              <el-switch v-model="formData.status" active-value="1" inactive-value="0" active-text="开启"
                         inactive-text="关闭"></el-switch>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="控制策略" prop="ctlStatus">
              <el-switch v-model="formData.ctlStatus" active-value="2" inactive-value="1" active-text="精准控制"
                         inactive-text="模糊控制"></el-switch>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="returnPreStep(3000)">上一步</el-button>
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>


      <!--  创建 10000 -->
      <el-dialog v-model="dialog10000FormVisible" :before-close="close10000Dialog" :draggable="true" :title="typeTitle"
               destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="产码方式" prop="type" disabled>
              <el-button v-model="formData.type" readonly disabled>
                {{ formatProdType(Number(formData.type)) }}
              </el-button>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通道ID" prop="cid" disabled>
              <el-input v-model="formData.cid" readonly disabled></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="通道账户" prop="acAccount">
              <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="账户备注" prop="acRemark">
              <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
        </el-row>
        <!--        <el-form-item label="报文" prop="token">
                  <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入"/>
                </el-form-item>-->
        <el-form-item label="付款链接" prop="token">

              <el-input v-model="formData.token" type="textarea" :clearable="true"
                        placeholder="请输入付款链接"/>
        </el-form-item>
        <el-row>
          <el-col :span="1"></el-col>
          <el-col :span="22">
            <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
          </el-col>
          <el-col :span="1"></el-col>
          <el-col :span="12">
            <el-form-item label="日限额" prop="dailyLimit">
              <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="总限额" prop="totalLimit">
              <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="日限进单" prop="dlyCntLimit">
              <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限进单" prop="inCntLimit">
              <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限拉单" prop="countLimit">
              <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="状态开关" prop="status">
              <el-switch v-model="formData.status" active-value="1" inactive-value="0" active-text="开启"
                         inactive-text="关闭"></el-switch>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="控制策略" prop="ctlStatus">
              <el-switch v-model="formData.ctlStatus" active-value="2" inactive-value="1" active-text="精准控制"
                         inactive-text="模糊控制"></el-switch>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="returnPreStep(10000)">上一步</el-button>
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  创建 4000 -->
    <el-dialog v-model="dialog4000FormVisible" :before-close="close4000Dialog" :draggable="true" :title="typeTitle"
               destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="产码方式" prop="type" disabled>
              <el-button v-model="formData.type" readonly disabled>
                {{ formatProdType(Number(formData.type)) }}
              </el-button>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通道ID" prop="cid" disabled>
              <el-input v-model="formData.cid" readonly disabled></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="通道账户" prop="acAccount">
              <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="账户备注" prop="acRemark">
              <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="报文" prop="token">
          <el-input v-model="formData.token" type="textarea" :clearable="true"
                    placeholder="输入示例：https://yaoshi.sdo.com/apipool?system_deviceId=83241004891122-40caee54f545f373&sequence=7&isHttp=0&netFlag=WIFI&method=txz_bs_mixed.dqOrder.list&ticket=cFFHWkDhsIx112NhvcX09FFaDUm3p%2F4Pk146eTWCS7IdU34mtVVI8rgoCVXnJmTJ9kfSiHkT0BP8SXK1sdeLempgKsItWc2F3FPn3BMsa6stXomxFjDNyaOieJADp3NapOnjl9Qnh7n9zi%2BavTlWAxE45Y9R38iCZz6x98tLMu0%3D&txzDeviceId=861110048918892&sndaId=1122890350&maxCount=10&version=a.9.4.8&timestampMs=-1"/>
        </el-form-item>
        <el-row>
          <el-col :span="1"></el-col>
          <el-col :span="22">
            <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
          </el-col>
          <el-col :span="1"></el-col>
          <el-col :span="12">
            <el-form-item label="日限额" prop="dailyLimit">
              <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="总限额" prop="totalLimit">
              <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="日限进单" prop="dlyCntLimit">
              <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限进单" prop="inCntLimit">
              <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限拉单" prop="countLimit">
              <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="状态开关" prop="status">
              <el-switch v-model="formData.status" active-value="1" inactive-value="0" active-text="开启"
                         inactive-text="关闭"></el-switch>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="控制策略" prop="ctlStatus">
              <el-switch v-model="formData.ctlStatus" active-value="2" inactive-value="1" active-text="精准控制"
                         inactive-text="模糊控制"></el-switch>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="returnPreStep(4000)">上一步</el-button>
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  创建 5000 -->
    <el-dialog v-model="dialog5000FormVisible" :before-close="close5000Dialog" :draggable="true" :title="typeTitle"
               destroy-on-close overflow>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="产码方式" prop="type" disabled>
              <el-button v-model="formData.type" readonly disabled>
                {{ formatProdType(Number(formData.type)) }}
              </el-button>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通道ID" prop="cid" disabled>
              <el-input v-model="formData.cid" readonly disabled></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="账户" prop="acAccount">
              <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="备注" prop="acRemark">
              <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="CK" prop="token">
          <el-input v-model="formData.token" type="textarea" :clearable="true"
                    placeholder="输入CK"/>
        </el-form-item>
        <el-row>
          <el-col :span="1"></el-col>
          <el-col :span="22">
            <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
          </el-col>
          <el-col :span="1"></el-col>
          <el-col :span="12">
            <el-form-item label="日限额" prop="dailyLimit">
              <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="总限额" prop="totalLimit">
              <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="日限进单" prop="dlyCntLimit">
              <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限进单" prop="inCntLimit">
              <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限拉单" prop="countLimit">
              <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="状态开关" prop="status">
              <el-switch v-model="formData.status" active-value="1" inactive-value="0" active-text="开启"
                         inactive-text="关闭"></el-switch>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="控制策略" prop="ctlStatus">
              <el-switch v-model="formData.ctlStatus" active-value="2" inactive-value="1" active-text="精准控制"
                         inactive-text="模糊控制"></el-switch>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="returnPreStep(5000)">上一步</el-button>
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  创建 6000 -->
    <el-dialog v-model="dialog6000FormVisible" :before-close="close6000Dialog" :draggable="true" :title="typeTitle"
               destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="产码方式" prop="type" disabled>
              <el-button v-model="formData.type" readonly disabled>
                {{ formatProdType(Number(formData.type)) }}
              </el-button>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通道ID" prop="cid" disabled>
              <el-input v-model="formData.cid" readonly disabled></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="CK" prop="token">
          <el-input v-model="formData.token" type="textarea" :clearable="true" :rows="4"
                    placeholder="请输入CK" @input="handleTokenInput"/>
        </el-form-item>
        <el-row :gutter="24">
          <el-col :span="24">
            <el-form-item label="通道账户" prop="acAccount">
              <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入账户"/>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="备注" prop="acRemark">
              <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入备注"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="1"></el-col>
          <el-col :span="22">
            <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
          </el-col>
          <el-col :span="1"></el-col>
          <el-col :span="12">
            <el-form-item label="日限额" prop="dailyLimit">
              <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="总限额" prop="totalLimit">
              <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="日限进单" prop="dlyCntLimit">
              <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限进单" prop="inCntLimit">
              <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限拉单" prop="countLimit">
              <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row>
          <el-col :span="12">
            <el-form-item label="状态开关" prop="status">
              <el-switch v-model="formData.status" active-value="1" inactive-value="0" active-text="开启"
                         inactive-text="关闭"></el-switch>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="控制策略" prop="ctlStatus">
              <el-switch v-model="formData.ctlStatus" active-value="2" inactive-value="1" active-text="精准控制"
                         inactive-text="模糊控制"></el-switch>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="returnPreStep(6000)">上一步</el-button>
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  创建 8000 -->
    <el-dialog v-model="dialog8000FormVisible" :before-close="close8000Dialog" :draggable="true" :title="typeTitle"
               destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule8000" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="产码方式" prop="type" disabled>
              <el-button v-model="formData.type" readonly disabled>
                {{ formatProdType(Number(formData.type)) }}
              </el-button>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通道ID" prop="cid" disabled>
              <el-input v-model="formData.cid" readonly disabled></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="通道账户" prop="acAccount">
              <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入账号"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="账户备注" prop="acRemark">
              <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入备注"/>
            </el-form-item>
          </el-col>

        </el-row>
        <el-form-item label="CK" prop="token">
          <el-input v-model="formData.token" type="textarea" :clearable="true"
                    placeholder="输入示例：https://security.seasungame.com/security_extend_server/helper/balance/queryBalance?gameCode=jx3&account=aa123123&accountType=&zoneCode=z22&SN=98710641126&remark=&sign=36A360706FD189A2BF867D70F61117BE"/>
        </el-form-item>
        <el-row>
          <el-col :span="1"></el-col>
          <el-col :span="22">
            <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
          </el-col>
          <el-col :span="1"></el-col>
          <el-col :span="12">
            <el-form-item label="日限额" prop="dailyLimit">
              <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="总限额" prop="totalLimit">
              <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="日限进单" prop="dlyCntLimit">
              <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限进单" prop="inCntLimit">
              <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限拉单" prop="countLimit">
              <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="状态开关" prop="status">
          <el-switch v-model="formData.status" active-value="1" inactive-value="0" active-text="开启"
                     inactive-text="关闭"></el-switch>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="returnPreStep(8000)">上一步</el-button>
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  创建 9000 -->
    <el-dialog v-model="dialog9000FormVisible" :before-close="close9000Dialog" :draggable="true" :title="typeTitle"
               destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="产码方式" prop="type" disabled>
              <el-button v-model="formData.type" readonly disabled>
                {{ formatProdType(Number(formData.type)) }}
              </el-button>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通道ID" prop="cid" disabled>
              <el-input v-model="formData.cid" readonly disabled></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="24">
            <el-form-item label="账户备注" prop="acRemark">
              <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入备注"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="CK" prop="token">
          <el-input v-model="formData.token" type="textarea" :clearable="true"
                    placeholder="输入CK"/>
        </el-form-item>
        <el-row>
          <el-col :span="1"></el-col>
          <el-col :span="22">
            <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
          </el-col>
          <el-col :span="1"></el-col>
          <el-col :span="12">
            <el-form-item label="日限额" prop="dailyLimit">
              <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="总限额" prop="totalLimit">
              <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="日限进单" prop="dlyCntLimit">
              <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限进单" prop="inCntLimit">
              <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总限拉单" prop="countLimit">
              <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="状态开关" prop="status">
          <el-switch v-model="formData.status" active-value="1" inactive-value="0" active-text="开启"
                     inactive-text="关闭"></el-switch>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="returnPreStep(9000)">上一步</el-button>
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  修改 1000 -->
    <el-dialog v-model="dialogUpd1000FormVisible" :before-close="closeUpd1000Dialog" :draggable="true"
               :title="typeTitle" destroy-on-close>
      <el-scrollbar height="300px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
          <el-row>
            <el-col :span="24">
              <el-form-item label="通道ID" prop="cid">
                <el-input v-model="formData.cid" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="账户备注" prop="acRemark">
                <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="通道账户" prop="acAccount">
                <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <!--          <el-form-item label="报文" prop="token">
                      <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入"/>
                    </el-form-item>-->
          <el-form-item label="扫码授权" prop="token">
            <el-row :gutter="12">
              <el-col :span="24">
                <div v-if="imageQrCode">
                  <img :src="imageQrCode" alt="qr" style="width: 180px;height: 180px"/>
                </div>
                <div v-else>
                  <div style="width: 180px;height: 180px;font-size: 12px;color:grey;border: 1px solid grey">
                    请点击"获取二维码"
                  </div>
                </div>
              </el-col>
              <el-col :span="24">
                <el-button type="primary" link @click="loginQr">重新获取二维码</el-button>
              </el-col>

              <el-col :span="18">
                <el-input v-model="imageQrStatusMsg" placeholder="二维码状态" style="color: red !important;"/>
              </el-col>

              <el-col :span="6">
                <el-button type="primary" @click="getQrCookie" style="width: 100%">获取CK</el-button>
              </el-col>
              <el-col :span="24">
                <el-input v-model="formData.token" type="textarea" :clearable="true"
                          placeholder="请输入CK或点击【获取CK】自动获取"/>
              </el-col>
            </el-row>
          </el-form-item>
          <el-row>
            <el-col :span="1"></el-col>
            <el-col :span="22">
              <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
            </el-col>
            <el-col :span="1"></el-col>
            <el-col :span="12">
              <el-form-item label="日限额" prop="dailyLimit">
                <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="总限额" prop="totalLimit">
                <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="日限进单" prop="dlyCntLimit">
                <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限进单" prop="inCntLimit">
                <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限拉单" prop="countLimit">
                <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeUpd1000Dialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  修改 1100 -->
    <el-dialog v-model="dialogUpd1100FormVisible" :before-close="closeUpd1100Dialog" :draggable="true"


    <!--  修改 10000 -->
    <el-dialog v-model="dialogUpd10000FormVisible" :before-close="closeUpd10000Dialog" :draggable="true"
               :title="typeTitle" destroy-on-close>
      <el-scrollbar height="300px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
          <el-row>
            <el-col :span="24">
              <el-form-item label="通道ID" prop="cid">
                <el-input v-model="formData.cid" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="账户备注" prop="acRemark">
                <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="通道账户" prop="acAccount">
                <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <!--          <el-form-item label="报文" prop="token">
                      <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入"/>
                    </el-form-item>-->
          <el-form-item label="付款链接" prop="token">

                <el-input v-model="formData.token" type="textarea" :clearable="true"
                          placeholder="请输入付款链接"/>

          </el-form-item>
          <el-row>
            <el-col :span="1"></el-col>
            <el-col :span="22">
              <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
            </el-col>
            <el-col :span="1"></el-col>
            <el-col :span="12">
              <el-form-item label="日限额" prop="dailyLimit">
                <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="总限额" prop="totalLimit">
                <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="日限进单" prop="dlyCntLimit">
                <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限进单" prop="inCntLimit">
                <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限拉单" prop="countLimit">
                <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeUpd10000Dialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  修改 2000 -->
    <el-dialog v-model="dialogUpd2000FormVisible" :before-close="closeUpd2000Dialog" :draggable="true"
               :title="typeTitle" destroy-on-close>
      <el-scrollbar height="300px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
          <el-row>
            <el-col :span="24">
              <el-form-item label="通道ID" prop="cid">
                <el-input v-model="formData.cid" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="账户备注" prop="acRemark">
                <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="通道账户" prop="acAccount">
                <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <!--          <el-form-item label="报文" prop="token">
                      <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入"/>
                    </el-form-item>-->
          <el-form-item label="扫码授权" prop="token">
            <el-row :gutter="12">
              <el-col :span="24">
                <div v-if="imageQrCode">
                  <img :src="imageQrCode" alt="qr" style="width: 180px;height: 180px"/>
                </div>
                <div v-else>
                  <div style="width: 180px;height: 180px;font-size: 12px;color:grey;border: 1px solid grey">
                    请点击"获取二维码"
                  </div>
                </div>
              </el-col>
              <el-col :span="24">
                <el-button type="primary" link @click="loginQr">重新获取二维码</el-button>
              </el-col>

              <el-col :span="18">
                <el-input v-model="imageQrStatusMsg" placeholder="二维码状态" style="color: red !important;"/>
              </el-col>

              <el-col :span="6">
                <el-button type="primary" @click="getQrCookie" style="width: 100%">获取CK</el-button>
              </el-col>
              <el-col :span="24">
                <el-input v-model="formData.token" type="textarea" :clearable="true"
                          placeholder="请输入CK或点击【获取CK】自动获取"/>
              </el-col>
            </el-row>
          </el-form-item>
          <el-row>
            <el-col :span="1"></el-col>
            <el-col :span="22">
              <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
            </el-col>
            <el-col :span="1"></el-col>
            <el-col :span="12">
              <el-form-item label="日限额" prop="dailyLimit">
                <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="总限额" prop="totalLimit">
                <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="日限进单" prop="dlyCntLimit">
                <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限进单" prop="inCntLimit">
                <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限拉单" prop="countLimit">
                <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeUpd1100Dialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  修改 1200 -->
    <el-dialog v-model="dialogUpd1200FormVisible" :before-close="closeUpd1200Dialog" :draggable="true"
               :title="typeTitle" destroy-on-close>
      <el-scrollbar height="300px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
          <el-row>
            <el-col :span="24">
              <el-form-item label="通道ID" prop="cid">
                <el-input v-model="formData.cid" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="账户备注" prop="acRemark">
                <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="通道账户" prop="acAccount">
                <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <!--          <el-form-item label="报文" prop="token">
                      <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入"/>
                    </el-form-item>-->
          <el-form-item label="扫码授权" prop="token">
            <el-row :gutter="12">
              <el-col :span="24">
                <div v-if="imageQrCode">
                  <img :src="imageQrCode" alt="qr" style="width: 180px;height: 180px"/>
                </div>
                <div v-else>
                  <div style="width: 180px;height: 180px;font-size: 12px;color:grey;border: 1px solid grey">
                    请点击"获取二维码"
                  </div>
                </div>
              </el-col>
              <el-col :span="24">
                <el-button type="primary" link @click="loginQr">重新获取二维码</el-button>
              </el-col>

              <el-col :span="18">
                <el-input v-model="imageQrStatusMsg" placeholder="二维码状态" style="color: red !important;"/>
              </el-col>

              <el-col :span="6">
                <el-button type="primary" @click="getQrCookie" style="width: 100%">获取CK</el-button>
              </el-col>
              <el-col :span="24">
                <el-input v-model="formData.token" type="textarea" :clearable="true"
                          placeholder="请输入CK或点击【获取CK】自动获取"/>
              </el-col>
            </el-row>
          </el-form-item>
          <el-row>
            <el-col :span="1"></el-col>
            <el-col :span="22">
              <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
            </el-col>
            <el-col :span="1"></el-col>
            <el-col :span="12">
              <el-form-item label="日限额" prop="dailyLimit">
                <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="总限额" prop="totalLimit">
                <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="日限进单" prop="dlyCntLimit">
                <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限进单" prop="inCntLimit">
                <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限拉单" prop="countLimit">
                <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeUpd1200Dialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  修改 2000 -->
    <el-dialog v-model="dialogUpd2000FormVisible" :before-close="closeUpd2000Dialog" :draggable="true"
               :title="typeTitle" destroy-on-close>
      <el-scrollbar height="300px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule2000" label-width="80px">
          <el-row>
            <el-col :span="24">
              <el-form-item label="通道id" prop="cid">
                <el-input v-model="formData.cid" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="账户备注" prop="acRemark">
                <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="通道账户" prop="acAccount">
                <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <!--          <el-form-item label="报文" prop="token">
                      <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入"/>
                    </el-form-item>-->
          <el-row :gutter="12">
            <el-col :span="24">
              <el-form-item label="报文" prop="token">
                    <el-input v-model="formData.token" type="textarea" :clearable="true"
                              placeholder="输入报文链接"/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="1"></el-col>
            <el-col :span="22">
              <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
            </el-col>
            <el-col :span="1"></el-col>
            <el-col :span="12">
              <el-form-item label="日限额" prop="dailyLimit">
                <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="总限额" prop="totalLimit">
                <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="日限进单" prop="dlyCntLimit">
                <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限进单" prop="inCntLimit">
                <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限拉单" prop="countLimit">
                <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeUpd2000Dialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>
    <!--  修改 3000 -->
    <el-dialog v-model="dialogUpd3000FormVisible" :before-close="closeUpd3000Dialog" :draggable="true"
               :title="typeTitle" destroy-on-close>
      <el-scrollbar height="300px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
          <el-row>
            <el-col :span="24">
              <el-form-item label="通道ID" prop="cid">
                <el-input v-model="formData.cid" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="账户备注" prop="acRemark">
                <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="通道账户" prop="acAccount">
                <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <!--          <el-form-item label="报文" prop="token">
                      <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入"/>
                    </el-form-item>-->
          <el-form-item label="扫码授权" prop="token">
            <el-row :gutter="12">
              <el-col :span="24">
                <div v-if="imageQrCode">
                  <img :src="imageQrCode" alt="qr" style="width: 180px;height: 180px"/>
                </div>
                <div v-else>
                  <div style="width: 180px;height: 180px;font-size: 12px;color:grey;border: 1px solid grey">
                    请点击"获取二维码"
                  </div>
                </div>
              </el-col>
              <el-col :span="24">
                <el-button type="primary" link @click="loginQr">重新获取二维码</el-button>
              </el-col>

              <el-col :span="18">
                <el-input v-model="imageQrStatusMsg" placeholder="二维码状态" style="color: red !important;"/>
              </el-col>

              <el-col :span="6">
                <el-button type="primary" @click="getQrCookie" style="width: 100%">获取CK</el-button>
              </el-col>
              <el-col :span="24">
                <el-input v-model="formData.token" type="textarea" :clearable="true"
                          placeholder="请输入CK或点击【获取CK】自动获取"/>
              </el-col>
            </el-row>
          </el-form-item>
          <el-row>
            <el-col :span="1"></el-col>
            <el-col :span="22">
              <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
            </el-col>
            <el-col :span="1"></el-col>
            <el-col :span="12">
              <el-form-item label="日限额" prop="dailyLimit">
                <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="总限额" prop="totalLimit">
                <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="日限进单" prop="dlyCntLimit">
                <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限进单" prop="inCntLimit">
                <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限拉单" prop="countLimit">
                <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeUpd3000Dialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  修改 4000 -->
    <el-dialog v-model="dialogUpd4000FormVisible" :before-close="closeUpd4000Dialog" :draggable="true"
               :title="typeTitle" destroy-on-close>
      <el-scrollbar height="300px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
          <el-row>
            <el-col :span="24">
              <el-form-item label="通道ID" prop="cid">
                <el-input v-model="formData.cid" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="账户备注" prop="acRemark">
                <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="通道账户" prop="acAccount">
                <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="报文" prop="token">
            <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入"/>
          </el-form-item>
          <el-row>
            <el-col :span="1"></el-col>
            <el-col :span="22">
              <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
            </el-col>
            <el-col :span="1"></el-col>
            <el-col :span="12">
              <el-form-item label="日限额" prop="dailyLimit">
                <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="总限额" prop="totalLimit">
                <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="日限进单" prop="dlyCntLimit">
                <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限进单" prop="inCntLimit">
                <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限拉单" prop="countLimit">
                <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeUpd4000Dialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  修改 5000 -->
    <el-dialog v-model="dialogUpd5000FormVisible" :before-close="closeUpd5000Dialog" :draggable="true"
               :title="typeTitle" destroy-on-close>
      <el-scrollbar height="300px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
          <el-row>
            <el-col :span="24">
              <el-form-item label="通道ID" prop="cid">
                <el-input v-model="formData.cid" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="账户" prop="acAccount">
                <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入账户" disabled/>
              </el-form-item>
            </el-col>

            <el-col :span="12">
              <el-form-item label="备注" prop="acRemark">
                <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入备注"/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="CK" prop="token">
            <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入CK"/>
          </el-form-item>
          <el-row>
            <el-col :span="1"></el-col>
            <el-col :span="22">
              <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
            </el-col>
            <el-col :span="1"></el-col>
            <el-col :span="12">
              <el-form-item label="日限额" prop="dailyLimit">
                <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="总限额" prop="totalLimit">
                <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="日限进单" prop="dlyCntLimit">
                <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限进单" prop="inCntLimit">
                <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限拉单" prop="countLimit">
                <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeUpd5000Dialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  修改 6000 -->
    <el-dialog v-model="dialogUpd6000FormVisible" :before-close="closeUpd6000Dialog" :draggable="true"
               :title="typeTitle" destroy-on-close>
      <el-scrollbar height="300px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
          <el-row>
            <el-col :span="24">
              <el-form-item label="通道ID" prop="cid">
                <el-input v-model="formData.cid" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="账户备注" prop="acRemark">
                <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="通道账户" prop="acAccount">
                <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="CK" prop="token">
            <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入"/>
          </el-form-item>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeUpd6000Dialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  修改 8000 -->
    <el-dialog v-model="dialogUpd8000FormVisible" :before-close="closeUpd8000Dialog" :draggable="true"
               :title="typeTitle" destroy-on-close>
      <el-scrollbar height="300px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule8000" label-width="80px">
          <el-row>
            <el-col :span="24">
              <el-form-item label="通道id" prop="cid">
                <el-input v-model="formData.cid" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="账户备注" prop="acRemark">
                <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="通道账户" prop="acAccount">
                <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="12">
            <el-col :span="24">
              <el-form-item label="CK" prop="token">
                <el-input v-model="formData.token" type="textarea" :clearable="true"
                          placeholder="输入CK"/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="1"></el-col>
            <el-col :span="22">
              <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
            </el-col>
            <el-col :span="1"></el-col>
            <el-col :span="12">
              <el-form-item label="日限额" prop="dailyLimit">
                <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="总限额" prop="totalLimit">
                <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="日限进单" prop="dlyCntLimit">
                <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限进单" prop="inCntLimit">
                <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限拉单" prop="countLimit">
                <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeUpd8000Dialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog v-model="dialogUpd9000FormVisible" :before-close="closeUpd9000Dialog" :draggable="true"
               :title="typeTitle" destroy-on-close>
      <el-scrollbar height="300px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule2000" label-width="80px">
          <el-row>
            <el-col :span="24">
              <el-form-item label="通道id" prop="cid">
                <el-input v-model="formData.cid" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item label="账户备注" prop="acRemark">
                <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="通道账户" prop="acAccount">
                <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="12">
            <el-col :span="24">
              <el-form-item label="CK" prop="token">
                <el-input v-model="formData.token" type="textarea" :clearable="true"
                          placeholder="输入CK"/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="1"></el-col>
            <el-col :span="22">
              <warning-bar title="注：默认0，则无限额控制。【日限额/总限额】为金额限制,【进单限数/拉单限数】为笔数限制"/>
            </el-col>
            <el-col :span="1"></el-col>
            <el-col :span="12">
              <el-form-item label="日限额" prop="dailyLimit">
                <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="总限额" prop="totalLimit">
                <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="日限进单" prop="dlyCntLimit">
                <el-input v-model.number="formData.dlyCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限进单" prop="inCntLimit">
                <el-input v-model.number="formData.inCntLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限拉单" prop="countLimit">
                <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeUpd9000Dialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  CK  -->
    <el-dialog v-model="dialogTokenFormVisible" :before-close="closeUpdTokenDialog" :draggable="true" title="变更CK"
               destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="token" prop="token">
          <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeUpdTokenDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 查看详情 -->
    <el-dialog v-model="detailShow" width="60%" lock-scroll :draggable="true" :before-close="closeDetailShow"
               title="查看详情" overflow
               destroy-on-close>
      <div v-if="stepDataShow">
        <el-row :gutter="24">
          <el-col :span="14">
            <el-scrollbar height="550px">
              <el-descriptions :column="6" border>
                <el-descriptions-item label="用户归属" :span="6">{{ formData.username }}</el-descriptions-item>
                <el-descriptions-item label="账户ID" :span="6">{{ formData.acId }}</el-descriptions-item>
                <el-descriptions-item label="账户备注" :span="6">{{ formData.acRemark }}</el-descriptions-item>
                <el-descriptions-item label="通道账户" :span="3">{{ formData.acAccount }}</el-descriptions-item>
                <el-descriptions-item label="账户密码" :span="3">{{ formData.acPwd }}</el-descriptions-item>
                <el-descriptions-item label="ck" :span="6">
                  <el-input v-model="formData.token" type="textarea" readonly/>
                </el-descriptions-item>
                <el-descriptions-item label="通道id" :span="6">{{ formData.cid }}</el-descriptions-item>
                <el-descriptions-item label="总限拉单" :span="2">{{ formData.countLimit }}</el-descriptions-item>
                <el-descriptions-item label="日限进单" :span="2">{{ formData.dlyCntLimit }}</el-descriptions-item>
                <el-descriptions-item label="总限进单" :span="2">{{ formData.inCntLimit }}</el-descriptions-item>
                <el-descriptions-item label="日限额" :span="3">{{ formData.dailyLimit }}</el-descriptions-item>
                <el-descriptions-item label="总限额" :span="3">{{ formData.totalLimit }}</el-descriptions-item>
                <el-descriptions-item label="状态开关" :span="2">{{ formData.status === 0 ? '关闭' : '开启' }}
                </el-descriptions-item>
                <el-descriptions-item label="系统开关" :span="2">{{ formData.sysStatus === 0 ? '关闭' : '开启' }}
                </el-descriptions-item>
                <el-descriptions-item label="限额策略" :span="2">{{ formData.status === 1 ? '模糊控制' : '精准控制' }}
                </el-descriptions-item>
              </el-descriptions>
            </el-scrollbar>
          </el-col>
          <el-col :span="10">
            <el-scrollbar height="550px">
              <el-timeline style="max-width: 600px; margin-top: 20px">
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
      </div>

      <div v-else>
        <el-row :gutter="24">
          <el-col :span="24">
            <el-scrollbar>
              <el-descriptions :column="6" border>
                <el-descriptions-item label="用户归属" :span="6">{{ formData.username }}</el-descriptions-item>
                <el-descriptions-item label="账户ID" :span="6">{{ formData.acId }}</el-descriptions-item>
                <el-descriptions-item label="账户备注" :span="6">{{ formData.acRemark }}</el-descriptions-item>
                <el-descriptions-item label="通道账户" :span="3">{{ formData.acAccount }}</el-descriptions-item>
                <el-descriptions-item label="账户密码" :span="3">{{ formData.acPwd }}</el-descriptions-item>
                <el-descriptions-item label="ck" :span="6">
                  <el-input v-model="formData.token" type="textarea" readonly/>
                </el-descriptions-item>
                <el-descriptions-item label="通道id" :span="6">{{ formData.cid }}</el-descriptions-item>
                <el-descriptions-item label="总限拉单" :span="2">{{ formData.countLimit }}</el-descriptions-item>
                <el-descriptions-item label="日限进单" :span="2">{{ formData.dlyCntLimit }}</el-descriptions-item>
                <el-descriptions-item label="总限进单" :span="2">{{ formData.inCntLimit }}</el-descriptions-item>
                <el-descriptions-item label="日限额" :span="3">{{ formData.dailyLimit }}</el-descriptions-item>
                <el-descriptions-item label="总限额" :span="3">{{ formData.totalLimit }}</el-descriptions-item>
                <el-descriptions-item label="状态开关" :span="2">{{ formData.status === 0 ? '关闭' : '开启' }}
                </el-descriptions-item>
                <el-descriptions-item label="系统开关" :span="2">{{ formData.sysStatus === 0 ? '关闭' : '开启' }}
                </el-descriptions-item>
                <el-descriptions-item label="限额策略" :span="2">{{ formData.status === 1 ? '模糊控制' : '精准控制' }}
                </el-descriptions-item>
              </el-descriptions>
            </el-scrollbar>
          </el-col>
        </el-row>

      </div>

    </el-dialog>

    <!-- 查看充值详情 1000 3000 -->
    <el-dialog v-model="orderHisVisible" style="width: 1100px" :draggable="true" lock-scroll
               :before-close="closeOrderHisShow"
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

    <!-- 查看充值详情 5000 -->
    <el-dialog v-model="orderHis5000Visible" style="width: 1100px" :draggable="true" lock-scroll
               :before-close="closeOrderHis5000Show"
               title="查看充值详情" destroy-on-close>
      <el-scrollbar height="550px">
        <el-table tooltip-effect="dark" :data="orderHis5000TableData" row-key="ID" style="width: 100%">
          <el-table-column align="left" label="付款人" prop="buyer" width="160"/>
          <el-table-column align="left" label="上游商品" prop="skuTitle"/>
          <el-table-column align="left" label="金额" prop="money" width="80"/>
          <el-table-column align="left" label="订单状态" prop="orderStatus" width="100"/>
          <el-table-column align="left" label="时间" prop="createTime" width="160">
            <template #default="scope">
              {{ formatDate(scope.row.createTime) }}
            </template>
          </el-table-column>
        </el-table>
      </el-scrollbar>
    </el-dialog>

    <!-- 查看充值详情 8000 -->
    <el-dialog v-model="orderHis8000Visible" style="width: 1100px" :draggable="true" lock-scroll
               :before-close="closeOrderHis8000Show"
               title="查看充值详情" destroy-on-close>
      <el-scrollbar height="550px">
        <el-descriptions :column="4" border style="background-color: #a5abb4">
          <el-descriptions-item label="用户ID">{{ orderHis8000Info.userId }}</el-descriptions-item>
          <el-descriptions-item label="积分">{{ orderHis8000Info.diamond }}</el-descriptions-item>
        </el-descriptions>
        <el-table tooltip-effect="dark" :data="orderHis8000List" row-key="ID" style="width: 100%">
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

    <!-- 查看充值详情 9000 -->
    <el-dialog v-model="orderHis9000Visible" style="width: 1100px" :draggable="true" lock-scroll
               :before-close="closeOrderHis9000Show"
               title="查看充值详情" destroy-on-close>
      <el-scrollbar height="550px">
        <el-descriptions :column="4" border style="background-color: #a5abb4">
          <el-descriptions-item label="用户">{{ orderHis9000Info.account }}</el-descriptions-item>
          <el-descriptions-item label="寄售点数">{{ orderHis9000Info.jsBalance }}</el-descriptions-item>
        </el-descriptions>
        <el-table tooltip-effect="dark" :data="orderHis9000List" row-key="ID" style="width: 100%">
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

    <!-- 查看产码详情 -->
    <el-dialog v-model="payCodeOverviewVisible" style="width: 1100px" lock-scroll :draggable="true"
               :before-close="closePayCodeOverviewShow" title="查看产码详情" destroy-on-close>
      <el-scrollbar height="550px">
        <div class="region-card-container">
          <div v-for="pcData in Object.entries(payCodeMap)" style="width: 100%">
            <!--            <div>￥{{ pcData[0].x1 }}，{{ formatOPSimple(pcData.x2) }}, {{ codeToText[pcData.x3] }}({{ pcData.x4 }})</div>-->
            <div>{{ formatOPSimple(pcData[0]) }}</div>
            <el-divider></el-divider>
            <span v-for="pcDetail in pcData[1]" style="padding: 10px">
              <el-badge :value="pcDetail.x4">
                <div v-if="formatRegionCode(pcDetail.x3, false)">
                  <el-button>{{ pcDetail.x1 }}元</el-button>
                </div>
                <div v-else>
                  <el-button>{{ formatRegionCode(pcDetail.x3, false) }} | {{ pcDetail.x1 }}元</el-button>
                </div>
              </el-badge>
            </span>
          </div>
        </div>
      </el-scrollbar>
    </el-dialog>

    <!-- 产码-->
    <el-dialog width="60%" v-model="pcDialogFormVisible" :before-close="closePcDialog" :title="typeTitle"
               :draggable="true"
               destroy-on-close>
      <el-form :model="pcFormData" label-position="right" ref="pcElFormRef" :rules="pcRule" label-width="100px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="通道" prop="cid">
              <el-input v-model="pcFormData.cid" :clearable="true" placeholder="请输入" disabled/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通道账户ID" prop="acId">
              <el-input v-model="pcFormData.acId" :clearable="true" placeholder="请输入" disabled/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通道账户" prop="acAccount">
              <el-input v-model="pcFormData.acAccount" :clearable="true" placeholder="请输入" disabled/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="账户备注" prop="acRemark">
              <el-input v-model="pcFormData.acRemark" :clearable="true" placeholder="请输入" disabled/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="过期时间" prop="expTime">
          <el-row>
            <el-col>
              <el-input-number v-model="numHours" size="small"
                               :parser="(value) => value.replace(/￥\s?|(,*)/g, '')"
                               controls-position="right" @change="handleChangeH" :min="0">
              </el-input-number>
              <span> 小 时 </span>
              <el-input-number v-model="numMinutes" size="small"
                               :parser="(value) => value.replace(/￥\s?|(,*)/g, '')"
                               controls-position="right" @change="handleChangeM" :min="0">
              </el-input-number>
              <span> 分 钟 </span>
              <el-input-number v-model="numSeconds" size="small"
                               :parser="(value) => value.replace(/￥\s?|(,*)/g, '')"
                               controls-position="right" @change="handleChangeS" :min="0">
              </el-input-number>
              <span> 秒 </span>
            </el-col>
            <el-col>
              <el-button link type="primary" @click="default7Day">7天</el-button>
              <el-button link type="primary" @click="default1Day">1天</el-button>
              <el-button link type="primary" @click="default2Hour">2小时</el-button>
              <el-button link type="primary" @click="default1Hour">1小时</el-button>
              <el-button link type="primary" @click="default10Minute">10分钟</el-button>
              <el-button link type="primary" @click="default0Second">重置</el-button>
            </el-col>
          </el-row>
        </el-form-item>

        <el-card style="{width: 100% !important}" shadow="never">
          <template #header>
            <div class="card-header">
              <span>明细</span>
            </div>
          </template>
          <div>
            <el-table :data="pcFormData.list" style="width: 100%">
              <el-table-column label="报文" prop="imgBaseStr" style="width: 100%">
                <template #default="scope">
                  <el-input :rows="2" type="textarea" v-if="activeUpdIndex === scope.$index"
                            v-model="scope.row.imgBaseStr" :required="true"></el-input>
                  <el-input :rows="2" type="textarea" disabled v-model="scope.row.imgBaseStr" readonly
                            v-else></el-input>
                </template>
              </el-table-column>
              <el-table-column label="金额（元）" prop="money" width="120px">
                <template #default="scope">
                  <el-input v-if="activeUpdIndex === scope.$index"
                            v-model.number="scope.row.money"
                            placeholder="输入金额"
                            :formatter="(value) => `￥ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')"
                            :parser="(value) => value.replace(/￥\s?|(,*)/g, '')" :required="true">
                  </el-input>
                  <span v-else>￥{{ scope.row.money }}</span>
                </template>
              </el-table-column>
              <!--                <el-table-column label="运营商" prop="operator" width="120px">
                                <template #default="scope">
                                  <el-select :required="true" v-model="scope.row.operator" placeholder="请选择通信商" filterable style="width: 100%">
                                    <el-option v-for="item in operators" :key="item.value" :label="item.label" :value="item.value"/>
                                  </el-select>
                                </template>
                              </el-table-column>
                              <el-table-column label="地区" prop="locList" width="120px">
                                <template #default="scope">
                                  <el-cascader
                                      :change-on-select="true"
                                      style="width:100%"
                                      :options="regionOptions"
                                      v-model="scope.row.locList"
                                      @change="chge"
                                      placeholder="选择地区"
                                      filterable
                                      :props="{checkStrictly: false}" :rules="cascaderRules"
                                  >
                                  </el-cascader>
                                </template>
                              </el-table-column>-->
              <el-table-column align="right" width="200">
                <template #header>
                  <el-button type="primary" @click="handleAdd2Upd">
                    <Plus style="width:1em; height:1em;"/>
                  </el-button>
                </template>
                <template #default="scope">
                  <div v-if="activeUpdIndex === scope.$index">
                    <el-button type="primary" @click="handleSave2Upd()"><Select style="width:1em; height:1em;"/>
                    </el-button>
                  </div>
                  <div v-else>
                    <el-button type="success" @click="handleEdit2Upd(scope.$index)">
                      <Edit style="width:1em; height:1em;"/>
                    </el-button>
                    <el-popconfirm @confirm="handleDelete2Upd(scope.$index)" width="220" confirm-button-text="Yes"
                                   cancel-button-text="No, Thanks" :icon="InfoFilled" icon-color="#626AEF"
                                   title="确定要删除该商品吗？">
                      <template #reference>
                        <el-button type="danger">
                          <Delete style="width:1em; height:1em;"/>
                        </el-button>
                      </template>
                    </el-popconfirm>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-card>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closePcDialog">取 消</el-button>
          <el-button type="primary" @click="enterPcDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 获取订单acc统计数据 -->
    <el-dialog v-model="showCostOrderAccVisible" :title="showCostOrderAccTitle" :draggable="true" width="1000px"
               @close="closeCostOrderAcc">
      <div class="gva-search-box">
        <el-form :inline="true" :model="searchAccInfo" class="demo-form-inline" @keyup.enter="onAccSubmit">
          <el-form-item label="账户名" prop="acId">
            <el-input v-model="searchAccInfo.acAccount" placeholder="搜索通道账户"/>
          </el-form-item>
          <el-form-item label="账户备注" prop="acRemark">
            <el-input v-model="searchAccInfo.acRemark" placeholder="搜索账户备注"/>
          </el-form-item>
          <el-form-item label="账户ID" prop="acAccount">
            <el-input v-model="searchAccInfo.acId" placeholder="搜索通道账户ID"/>
          </el-form-item>
          <el-form-item label="通道ID" prop="cid">
            <el-input v-model.number="searchAccInfo.channelCode" placeholder="搜索通道ID"/>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="search" @click="onAccSubmit">查询</el-button>
            <el-button icon="refresh" @click="onAccReset">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
      <div class="gva-table-box">
        <el-scrollbar>
          <el-table ref="multipleTable" tooltip-effect="dark" :data="costOrderAccTable" border resizable="true"
                    height="420"
                    show-summary>
            <el-table-column align="center" label="通道ID" width="80">
              <template #default="{row}">
                {{ String(row.channelCode) }}
              </template>
            </el-table-column>
            <el-table-column align="center" label="账号ID" width="90">
              <template #default="{row}">
                {{ String(row.acId) }}
              </template>
            </el-table-column>
            <el-table-column align="center" label="通道账号" width="160">
              <template #default="{row}">
                {{ String(row.acAccount) }}
              </template>
            </el-table-column>
            <el-table-column align="center" sortable label="3日前" prop="x1" width="120"/>
            <el-table-column align="center" sortable label="2日前" prop="x2" width="120"/>
            <el-table-column align="center" sortable label="昨日" prop="x3" width="120"/>
            <el-table-column align="center" sortable label="今日" prop="x4" width="120"/>
            <el-table-column align="center" sortable label="总充值" prop="x0" width="120"/>
          </el-table>
        </el-scrollbar>
      </div>
    </el-dialog>

    <!-- 查询指定账户订单 -->
    <el-dialog v-model="orderSysVisible" style="width: 1100px" lock-scroll :draggable="true" title="查看系统充值详情"
               destroy-on-close>
      <div class="gva-search-box">
        <el-form :inline="true" class="demo-form-inline">
          <el-form-item>
            <el-button icon="refresh" @click="resetSimple(true)">简约版</el-button>
            <el-button icon="refresh" @click="resetSimple(false)">详情版</el-button>
          </el-form-item>
        </el-form>
      </div>
      <div class="gva-table-box">
        <!--   简约版   -->
        <el-table v-if="isSimple" ref="multipleTable" style="width: 100%" tooltip-effect="dark"
                  :data="orderSysTableData" row-key="ID" border>
          <el-table-column align="center" label="账号" prop="acAccount" width="180"/>
          <el-table-column align="center" label="订单ID" prop="orderId" width="230"/>
          <el-table-column align="center" label="金额" prop="money" width="120"/>
          <el-table-column align="center" label="订单状态" prop="orderStatus" width="120">
            <template #default="scope">
              <el-button style="width: 90px"
                         :color="formatPayedColor(scope.row.orderStatus, scope.row.acId, scope.row.platId)">
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
              <el-button type="primary" link class="table-button" @click="getSysDetails(scope.row)">
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
        </el-table>

        <!--   详情版   -->
        <el-table v-else ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="orderSysTableData"
                  row-key="ID" border @selection-change="handleSelectionChange">
          <el-table-column type="selection" width="55"/>
          <el-table-column align="left" label="付方ID" prop="pAccount" width="160"/>
          <el-table-column align="left" label="单价积分" prop="unitPrice" width="120"/>
          <el-table-column align="left" label="用户ID" prop="CreatedBy" width="120"/>
          <el-table-column align="left" label="通道编码" prop="channelCode" width="120"/>
          <el-table-column align="left" label="平台id" prop="platId" width="320"/>
          <el-table-column align="left" label="访客ip" prop="payIp" width="180"/>
          <el-table-column align="left" label="区域" prop="payRegion" width="240"/>
          <el-table-column align="left" label="客户端设备" prop="payDevice" width="120"/>
          <el-table-column align="left" label="账号ID" prop="acId" width="180"/>
          <el-table-column align="left" label="订单ID" prop="orderId" width="230"/>
          <el-table-column align="left" label="金额" prop="money" width="120"/>
          <el-table-column align="left" label="订单状态" prop="orderStatus" width="120">
            <template #default="scope">
              <el-button style="width: 90px"
                         :color="formatPayedColor(scope.row.orderStatus, scope.row.acId, scope.row.platId)">
                {{ formatPayed(scope.row.orderStatus, scope.row.acId, scope.row.platId) }}
              </el-button>
            </template>
          </el-table-column>
          <el-table-column align="left" label="回调状态" prop="cbStatus" width="120">
            <template #default="scope">
              <el-button style="width: 90px" :color="formatNotifyColor(scope.row.cbStatus)">
                {{ formatNotify(scope.row.cbStatus) }}
              </el-button>
            </template>
          </el-table-column>
          <el-table-column align="left" label="补单状态" prop="handStatus" width="120">
            <template #default="scope">
              <el-button style="width: 90px" :color="formatHandNotifyColor(scope.row.handStatus)">
                {{ formatHandNotify(scope.row.handStatus) }}
              </el-button>
            </template>
          </el-table-column>
          <el-table-column align="left" label="创建时间" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
          </el-table-column>
          <el-table-column align="left" label="回调时间" width="180">
            <template #default="scope">{{ formatDate(scope.row.cbTime) }}</template>
          </el-table-column>
          <el-table-column align="left" label="操作" width="280">
            <template #default="scope">
              <el-button type="primary" link class="table-button" @click="getSysDetails(scope.row)">
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
        </el-table>
        <div class="gva-pagination">
          <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="sysPage"
              :page-size="sysPageSize"
              :page-sizes="[10, 30, 50, 100]"
              :total="sysTotal"
              @current-change="handleSysCurrentChange"
              @size-change="handleSysSizeChange"
          />
        </div>
      </div>
    </el-dialog>

    <!-- 订单查看详情 -->
    <el-dialog v-model="sysDetailShow" style="width: 800px" lock-scroll :before-close="closeSysDetailShow"
               title="查看详情" :draggable="true"
               destroy-on-close>
      <el-descriptions column="1" border>
        <el-descriptions-item label="订单ID">{{ sysFormData.orderId }}</el-descriptions-item>
        <el-descriptions-item label="付方ID">{{ sysFormData.pAccount }}</el-descriptions-item>
        <el-descriptions-item label="金额">{{ sysFormData.money }}</el-descriptions-item>
        <el-descriptions-item label="单价积分">{{ sysFormData.unitPrice }}</el-descriptions-item>
        <el-descriptions-item label="ID">{{ sysFormData.acId }}</el-descriptions-item>
        <el-descriptions-item label="账号">{{ sysFormData.acAccount }}</el-descriptions-item>
        <el-descriptions-item label="通道编码">{{ sysFormData.channelCode }}</el-descriptions-item>
        <el-descriptions-item label="平台ID">{{ sysFormData.platId }}</el-descriptions-item>
        <el-descriptions-item label="客户ip">{{ sysFormData.payIp }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ sysFormData.payRegion }}</el-descriptions-item>
        <el-descriptions-item label="客户端设备">{{ sysFormData.payDevice }}</el-descriptions-item>
        <el-descriptions-item label="订单状态">{{ formatBoolean(sysFormData.orderStatus) }}</el-descriptions-item>
        <el-descriptions-item label="回调状态">{{ formatBoolean(sysFormData.cbStatus) }}</el-descriptions-item>
        <el-descriptions-item label="回调时间">{{ formatDate(sysFormData.cbTime) }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <!--  补单  -->
    <el-dialog
        v-model="dialogSysFormVisible"
        :before-close="closeSysDialog"
        :title="typeTitle" :draggable="true"
        destroy-on-close
        style="width: 450px"
    >
      <el-scrollbar height="100px">
        <el-form :model="sysFormData" label-position="right" ref="elSysFormRef" label-width="120px">
          <el-form-item label="订单ID" prop="authCaptcha">
            <el-input disabled v-model="sysFormData.orderId" :clearable="true" placeholder="请输入" style="width: 80%"/>
          </el-form-item>
          <el-form-item label="谷歌动态验证" prop="authCaptcha">
            <el-input v-model="sysFormData.authCaptcha" :clearable="true" placeholder="请输入谷歌动态验证"
                      style="width: 80%"/>
          </el-form-item>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeSysDialog">取 消</el-button>
          <el-button type="primary" @click="enterSysDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  通道转移 -->
    <el-dialog v-model="dialogTransferFormVisible" :before-close="closeDialog" :draggable="true" :title="typeTitle"
               destroy-on-close style="width: 600px">
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-row>
          <el-col :span="24">
            <el-form-item label="当前通道" prop="cid" disabled>
              <el-input v-model="formData.cid" readonly disabled></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="通道账户" prop="acAccount">
              <el-input v-model="formData.acAccount" readonly disabled/>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="账户备注" prop="acRemark">
              <el-input v-model="formData.acRemark" readonly disabled/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="24">
            <el-form-item label="转移通道" prop="cid" :required="true">
              <el-cascader v-model="formData.cid" :options="channelCodeOptions" :props="channelCodeProps" @change=""
                           style="width: 100%"/>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createChannelAccount,
  deleteChannelAccount,
  deleteChannelAccountByIds,
  updateChannelAccount,
  findChannelAccount,
  getChannelAccountList,
  queryAccOrderHis,
  countAcc,
  switchEnableCA,
  switchEnableCAByIds,
  transferChannelForAcc,
  loginByQr,
  loginQrStatusCheck,
} from '@/api/channelAccount'
import {
  getChannelProductSelf
} from '@/api/channelProduct'
import {codeToText, regionData} from 'element-china-area-data';
import {useRouter} from 'vue-router'

// 全量引入格式化工具 请按需保留
import {
  getDictFunc,
  formatDate,
  formatBoolean,
  filterDict,
  ReturnArrImg,
  onDownloadFile,
  formatUtcTimestamp,
  formatJoin,
  formatMoneyDesc,
  formatOPDesc,
  formatOPSimple,
  formatPayedColor,
  formatNotifyColor,
  formatPayed,
  formatNotify,
  formatHandNotifyColor,
  formatHandNotify,
  formatCDStatusColor,
  formatCDStatus,
  formatRegionCode, formatProdType, findValuesContainingString
} from '@/utils/format'
import {ElMessage, ElMessageBox} from 'element-plus'
import {ref, reactive, nextTick} from 'vue'
import WarningBar from "@/components/warningBar/warningBar.vue";
import {Delete, Edit, Eleme, InfoFilled, Loading, Plus, Position, Search, Select} from "@element-plus/icons-vue";
import {
  batchCreateChannelPayCode,
  createChannelPayCode,
  findChannelPayCode,
  getPayCodeOverviewByChanAcc
} from "@/api/channelPayCode";
import dayjs from "dayjs";
import utcPlugin from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';
import provinces from "@/assets/json/provinces.json";
import {callback2Pa, findPayOrder, getOrderAccOverview, getPayOrderList} from "@/api/payOrder";
import CenterCard from "@/view/vbox/dashboard/dataCenterComponents/centerCard.vue";

defineOptions({
  name: 'ChannelAccount'
})

// 注册插件
dayjs.extend(utcPlugin);
dayjs.extend(timezone);

const countItem = ref([])
const imageQrCode = ref()
const imageQrSig = ref()
const imageQrStatusMsg = ref()
const loginQr = async () => {
  imageQrCode.value = ''
  imageQrSig.value = ''
  let res = await loginByQr()
  if (res.code === 0) {
    console.log(res.data.img)
    imageQrCode.value = res.data.img.QrImg;
    imageQrSig.value = res.data.img.QrSig;
  }
}
const getQrCookie = async () => {
  let res = await loginQrStatusCheck({sig: imageQrSig.value})
  if (res.code === 0) {
    console.log(res.data.ret)
    formData.value.token = 'openid=' + res.data.ret.OpenID + ';openkey=' + res.data.ret.OpenKey
    formData.value.acAccount = res.data.ret.qq
    imageQrStatusMsg.value = "获取成功，自动填入账号与CK信息"

  } else if (res.code === 7) {
    imageQrStatusMsg.value = res.msg
  }
}

const queryAccOrderHisFunc = async (row, cid) => {
  let res = await queryAccOrderHis(row)
  // console.log(res.data)
  if (cid >= 1000 && cid <= 1099) {
    if (res.code === 0) {
      orderHisTableData.value = res.data.list.WaterList
    }
  } else if (cid >= 1100 && cid <= 1199) {
    if (res.code === 0) {
      orderHisTableData.value = res.data.list.WaterList
    }
  } else if (cid >= 1200 && cid <= 1299) {
    if (res.code === 0) {
      orderHisTableData.value = res.data.list.WaterList
    }
  } else if (cid >= 2000 && cid <= 2099) {
    if (res.code === 0) {
      orderHis2000Info.value = res.data.list.info
      orderHis2000List.value = res.data.list.list
    }
  } else if (cid >= 3000 && cid <= 3099) {
    if (res.code === 0) {
      orderHisTableData.value = res.data.list.WaterList
    }
  } else if (cid >= 4000 && cid <= 4099) {
    if (res.code === 0) {
      orderHis4000TableData.value = res.data.list
    }
  } else if (cid >= 5000 && cid <= 5099) {
    if (res.code === 0) {
      orderHis5000TableData.value = res.data.list
    }
  } else if (cid >= 8000 && cid <= 8099) {
    if (res.code === 0) {
       console.log(JSON.stringify(res.data))
      orderHis8000Info.value = res.data.list.info
      orderHis8000List.value = res.data.list.list
    }
  } else if (cid >= 9000 && cid <= 9099) {
    if (res.code === 0) {
      orderHis9000Info.value = res.data.list.info
      orderHis9000List.value = res.data.list.list
    }
  }
}
// 系统查单
const sysPage = ref(1)
const sysTotal = ref(0)
const sysPageSize = ref(10)

// 系统查单分页 -----------
//页面简约切换
const isSimple = ref(true)
// 重置
const resetSimple = (status) => {
  isSimple.value = status
}
const handleSysSizeChange = (val) => {
  sysPageSize.value = val
  pageSize.value = val
  queryAccOrderSysFunc(req.value)
}

// 系统查单修改页面容量
const handleSysCurrentChange = (val) => {
  sysPage.value = val
  page.value = val
  queryAccOrderSysFunc(req.value)
}
const req = ref()
const queryAccOrderSysFunc = async (row) => {
  req.value = {...row}
  console.log(req)

  let res = await getPayOrderList({page: page.value, pageSize: pageSize.value, acId: req.value.acId, orderStatus: 1})
  console.log(res.data)
  if (res.code === 0) {
    orderSysTableData.value = res.data.list
    sysTotal.value = res.data.total
    sysPage.value = res.data.page
    sysPageSize.value = res.data.pageSize
  }
}
// 系统查单 打开详情
const getSysDetails = async (row) => {
  console.log(row)
  // 打开弹窗
  const res = await findPayOrder({ID: row.ID})
  if (res.code === 0) {
    sysFormData.value = res.data.repayOrder
    openSysDetailShow()
  }
}
const sysDetailShow = ref(false)
// 关闭详情弹窗
const closeSysDetailShow = () => {
  sysDetailShow.value = false
  sysFormData.value = {
    orderId: '',
    pAccount: '',
    money: 0,
    unitPrice: 0,
    uid: 0,
    acId: '',
    acAccount: '',
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
// 系统查单打开详情弹窗
const openSysDetailShow = () => {
  sysDetailShow.value = true
}
const sysFormData = ref({
  authCaptcha: '',
  orderId: '',
  pAccount: '',
  money: 0,
  unitPrice: 0,
  uid: 0,
  acId: '',
  acAccount: '',
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
//补单
// 打开弹窗
const openSysDialog = () => {
  type.value = 'notify'
  dialogSysFormVisible.value = true
  typeTitle.value = '补单'
}

// 打开详情（补单使用）
const notifyPayOrder = async (row) => {
  // 打开弹窗
  const res = await findPayOrder({ID: row.ID})
  if (res.code === 0) {
    sysFormData.value = res.data.repayOrder
    openSysDialog()
  }
}

const getPayCodeOverviewByChanAccFunc = async (row) => {
  const req = {...row}
  console.log(req)

  let res = await getPayCodeOverviewByChanAcc(req)
  console.log(res.data)
  if (res.code === 0) {
    if (res.data.list.length === 0) {
      ElMessage({
        type: 'error',
        message: '该账号暂无可用的预产记录'
      })
    }
    payCodeTableData.value = res.data.list;
    payCodeMap.value = payCodeTableData.value.reduce((acc, cur) => {
      const {x2, ...rest} = cur;
      acc[x2] = acc[x2] || [];
      acc[x2].push(rest);
      return acc;
    }, {});
  }
}
// queryAccOrderHisFunc()

//countAcc
const countAccFunc = async () => {
  let res = await countAcc()
  if (res.code === 0) {
    countItem.value = res.data.list
  }
  console.log(res.data);
}

//通道产品

const channelCodeOptions = ref([])
const vcpTableData = ref([])

const channelCodeProps = {
  expandTrigger: 'hover',
  checkStrictly: false,
  emitPath: false,
}

const handleChange = (value) => {
  getOptionData()
  console.log(value)
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

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  acId: '',
  acRemark: '',
  acAccount: '',
  acPwd: '',
  token: '',
  cid: '',
  countLimit: 0,
  dailyLimit: 0,
  dlyCntLimit: 0,
  inCntLimit: 0,
  totalLimit: 0,
  type: 0,
  status: 0,
  sysStatus: 0,
  ctlStatus: 0,
  cbStatus: 0,
  username: '',
})

// 验证规则
const rule = reactive({
  cid: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
    {
      whitespace: true,
      message: '不能只输入空格',
      trigger: ['input', 'blur'],
    }],
})
// 验证规则
const rule2000 = reactive({
  acRemark: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  }, {
    whitespace: true,
    message: '不能只输入空格',
    trigger: ['input', 'blur'],
  }],
  cid: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
    {
      whitespace: true,
      message: '不能只输入空格',
      trigger: ['input', 'blur'],
    }],
})
// 验证规则
const rule8000 = reactive({
  acRemark: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  }, {
    whitespace: true,
    message: '不能只输入空格',
    trigger: ['input', 'blur'],
  }],
  cid: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
    {
      whitespace: true,
      message: '不能只输入空格',
      trigger: ['input', 'blur'],
    }],
  acAccount: [{required: true, message: '', trigger: ['blur'],}],
})
// 验证规则
const pcRule = reactive({
  acAccount: [{required: true, message: '', trigger: ['blur'],}],
  platId: [{required: true, message: '请选择', trigger: ['blur'],}],
  cid: [{required: true, message: '请选择', trigger: ['blur'],}],
  acId: [{required: true, message: '请选择', trigger: ['input', 'blur'],}],
  expTime: [{required: true, validator: validateTimeLimit, trigger: 'blur',},],
  operator: [{required: true, message: '', trigger: ['blur'],}],
  location: [{required: true, message: '请选择省或者省市', trigger: ['input', 'blur'],}],
  money: [{validator: checkMoney, trigger: 'blur'}]
})

function validateTimeLimit(rule, value, callback) {
  if (numHours.value === 0 && numMinutes.value === 0 && numSeconds.value === 0) {
    callback(new Error('过期时间填写不能都为 0'));
  } else {
    callback();
  }
}

function checkMoney(rule, value, callback) {
  if (Number(value) <= 0) {
    callback(new Error('请输入正确的金额'));
  } else {
    callback();
  }
}

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
const elForm2000Ref = ref()
const elSearchFormRef = ref()
const elSysFormRef = ref()

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
  const table = await getChannelAccountList({page: page.value, pageSize: pageSize.value, ...searchInfo.value})
  const vcpTable = await getChannelProductSelf({page: 1, pageSize: 999, ...searchInfo.value})

  await countAccFunc()

  if (table.code === 0) {
    tableData.value = table.data.list
    vcpTableData.value = vcpTable.data.list
    //card select
    parentNodes.value = getParentNodes(vcpTableData.value);

    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
    // setOptions()
    // console.log(provinces)
    setRegionOptions(provinces, regionOptions.value, false)
  }
}

// 根据不同的产品类型切换 Option
const payType = ref(0)

const getOptionData = async () => {
  const vcpTable = await getChannelProductSelf({page: 1, pageSize: 999, type: formData.value.type})
  await countAccFunc()

  if (vcpTable.code === 0) {
    vcpTableData.value = vcpTable.data.list
    setOptions()
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () => {
  channelCodeOptions.value = []
  setChannelCodeOptions(vcpTableData.value, channelCodeOptions.value, false)
}

const setRegionOptions = (ChannelCodeData, optionsData, disabled) => {
  ChannelCodeData &&
  ChannelCodeData.forEach(item => {
    if (item.children && item.children.length) {
      const option = {
        value: item.code + '',
        label: item.name,
        children: []
      }
      setRegionOptions(
          item.children,
          option.children,
      )
      optionsData.push(option)
    } else {
      const option = {
        value: item.code + '',
        label: item.name,
      }
      optionsData.push(option)
    }
  })
}
// 获取需要的字典 可能为空 按需保留

// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    deleteChannelAccountFunc(row)
  })
}

// 批量删除控制标记
const deleteVisible = ref(false)

// 多选删除
const onDelete = async () => {
  const ids = []
  if (multipleSelection.value.length === 0) {
    ElMessage({
      type: 'warning',
      message: '请选择要删除的数据'
    })
    return
  }
  multipleSelection.value &&
  multipleSelection.value.map(item => {
    ids.push(item.ID)
  })
  const res = await deleteChannelAccountByIds({ids})
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '删除成功'
    })
    if (tableData.value.length === ids.length && page.value > 1) {
      page.value--
    }
    deleteVisible.value = false
    getTableData()
  }
}

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')
const typeTitle = ref('')

// 更新行
// ca 更新
const dialogUpd9000FormVisible = ref(false)
const dialogUpd8000FormVisible = ref(false)
const dialogUpd6000FormVisible = ref(false)
const dialogUpd5000FormVisible = ref(false)
const dialogUpd4000FormVisible = ref(false)
const dialogUpd3000FormVisible = ref(false)
const dialogUpd2000FormVisible = ref(false)
const dialogUpd1200FormVisible = ref(false)
const dialogUpd1100FormVisible = ref(false)
const dialogUpd1000FormVisible = ref(false)
const updateChannelAccountFunc = async (row) => {
  const res = await findChannelAccount({ID: row.ID})
  type.value = 'update'
  typeTitle.value = '修改'
  if (res.code === 0) {
    formData.value = res.data.revca
    let cid = Number(res.data.revca.cid)
    if (cid >= 9000 && cid <= 9099) {
      dialogUpd9000FormVisible.value = true
    } else if (cid >= 8000 && cid <= 8099) {
      dialogUpd8000FormVisible.value = true
    } else if (cid >= 6000 && cid <= 6099) {
      dialogUpd6000FormVisible.value = true
    } else if (cid >= 5000 && cid <= 5099) {
      dialogUpd5000FormVisible.value = true
    } else if (cid >= 4000 && cid <= 4099) {
      dialogUpd4000FormVisible.value = true
    } else if (cid >= 3000 && cid <= 3099) {
      dialogUpd3000FormVisible.value = true
    } else if (cid >= 2000 && cid <= 2099) {
      dialogUpd2000FormVisible.value = true
    } else if (cid >= 1000 && cid <= 1099) {
      dialogUpd1000FormVisible.value = true
    } else if (cid >= 1100 && cid <= 1199) {
      dialogUpd1100FormVisible.value = true
    } else if (cid >= 1200 && cid <= 1299) {
      dialogUpd1200FormVisible.value = true
    }
  }
}

// 删除行
const deleteChannelAccountFunc = async (row) => {
  const res = await deleteChannelAccount({ID: row.ID})
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '删除成功'
    })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}

// 弹窗控制标记
const dialog1000FormVisible = ref(false)
const dialog1100FormVisible = ref(false)
const dialog1200FormVisible = ref(false)
const dialog2000FormVisible = ref(false)
const dialog3000FormVisible = ref(false)
const dialog4000FormVisible = ref(false)
const dialog5000FormVisible = ref(false)
const dialog6000FormVisible = ref(false)
const dialog8000FormVisible = ref(false)
const dialog9000FormVisible = ref(false)
const dialog10000FormVisible = ref(false)

// 系统查单补单
const dialogSysFormVisible = ref(false)
// 上一步
const returnPreStep = (cid) => {
  dialogChanFormVisible.value = true
  if (cid >= 1000 && cid <= 1099) {
    dialog1000FormVisible.value = false
  } else if (cid >= 1100 && cid <= 1199) {
    dialog1100FormVisible.value = false
  } else if (cid >= 1200 && cid <= 1299) {
    dialog1200FormVisible.value = false
  } else if (cid >= 2000 && cid <= 2099) {
    dialog2000FormVisible.value = false
  } else if (cid >= 3000 && cid <= 3099) {
    dialog3000FormVisible.value = false
  } else if (cid >= 4000 && cid <= 4099) {
    dialog4000FormVisible.value = false
  } else if (cid >= 5000 && cid <= 5099) {
    dialog5000FormVisible.value = false
  } else if (cid >= 6000 && cid <= 6099) {
    dialog6000FormVisible.value = false
  } else if (cid >= 8000 && cid <= 8099) {
    dialog8000FormVisible.value = false
  } else if (cid >= 9000 && cid <= 9099) {
    dialog9000FormVisible.value = false
  } else if (cid >= 10000 && cid <= 10099) {
    dialog10000FormVisible.value = false
  }

}

const closeSysDialog = () => {
  dialogSysFormVisible.value = false
  sysFormData.value = {
    authCaptcha: '',
    orderId: '',
    pAccount: '',
    money: 0,
    unitPrice: 0,
    uid: 0,
    acId: '',
    acAccount: '',
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
// 系统单 弹窗确定
const enterSysDialog = async () => {
  elSysFormRef.value?.validate(async (valid) => {
    if (!valid) return
    switch (type.value) {
      case 'notify':
        console.log(sysFormData.value);
        let res = await callback2Pa(sysFormData.value);
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: '回调成功'
          })
        }
        dialogSysFormVisible.value = false;
        break;
    }
  })
}
// 查看详情控制标记
const detailShow = ref(false)

// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}

const stepData = ref([])
const stepDataShow = ref(false)

// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findChannelAccount({ID: row.ID})
  if (res.code === 0) {
    formData.value = res.data.revca
    stepData.value = res.data.revca.ext.records

    if (stepData.value !== null) {
      for (let i = 0; i < stepData.value.length; i++) {
        stepData.value[i].CreatedAt = formatDate(stepData.value[i].CreatedAt)
      }
      stepDataShow.value = true
    }
    openDetailShow()
  }
}

// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}

// 打开弹窗
const openDialog = () => {
  type.value = 'create'
  typeTitle.value = '创建'
  channelCodeOptions.value = []
  dialogChanFormVisible.value = true
}

// 关闭弹窗
const close1000Dialog = () => {
  dialog1000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}
// 关闭弹窗
const close1100Dialog = () => {
  dialog1100FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    ctlStatus: 0,
    sysStatus: 0,
    uid: 0,
  }
}
// 关闭弹窗
const close1200Dialog = () => {
  dialog1200FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}
// 关闭弹窗
const close2000Dialog = () => {
  dialog2000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}
// 关闭弹窗
const close3000Dialog = () => {
  dialog3000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}


// 关闭弹窗
const close10000Dialog = () => {
  dialog10000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}
// 关闭弹窗
const close4000Dialog = () => {
  dialog4000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dailyLimit: 0,
    dlyCntLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}
// 关闭弹窗
const close5000Dialog = () => {
  dialog5000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}


// 关闭弹窗
const close6000Dialog = () => {
  dialog6000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    dlyCntLimit: 0,
    inCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}

// 关闭弹窗
const close8000Dialog = () => {
  dialog8000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}

// 关闭弹窗
const close9000Dialog = () => {
  dialog9000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}

const closeUpd9000Dialog = () => {
  dialogUpd9000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}

const closeUpd8000Dialog = () => {
  dialogUpd8000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}

const closeUpd6000Dialog = () => {
  dialogUpd6000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}

const closeUpd5000Dialog = () => {
  dialogUpd5000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}

const closeUpd4000Dialog = () => {
  dialogUpd4000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    ctlStatus: 0,
    sysStatus: 0,
    uid: 0,
  }
}
const closeUpd3000Dialog = () => {
  dialogUpd3000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}
const closeUpd2000Dialog = () => {
  dialogUpd2000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}
const closeUpd1000Dialog = () => {
  dialogUpd1000FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}
const closeUpd1100Dialog = () => {
  dialogUpd1100FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}
const closeUpd1200Dialog = () => {
  dialogUpd1200FormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}
const closeUpdTokenDialog = () => {
  dialogTokenFormVisible.value = false
  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}
const closeDialog = () => {
  dialogTokenFormVisible.value = false
  dialog1000FormVisible.value = false
  dialog1100FormVisible.value = false
  dialog1200FormVisible.value = false
  dialog2000FormVisible.value = false
  dialog3000FormVisible.value = false
  dialog4000FormVisible.value = false
  dialog5000FormVisible.value = false
  dialog6000FormVisible.value = false
  dialog8000FormVisible.value = false
  dialog9000FormVisible.value = false
  dialog10000FormVisible.value = false
  dialogUpd1000FormVisible.value = false
  dialogUpd1100FormVisible.value = false
  dialogUpd1200FormVisible.value = false
  dialogUpd2000FormVisible.value = false
  dialogUpd3000FormVisible.value = false
  dialogUpd4000FormVisible.value = false
  dialogUpd5000FormVisible.value = false
  dialogUpd6000FormVisible.value = false
  dialogUpd8000FormVisible.value = false
  dialogUpd9000FormVisible.value = false
  dialogUpd10000FormVisible.value = false
  dialogTransferFormVisible.value = false

  formData.value = {
    acId: '',
    acRemark: '',
    acAccount: '',
    acPwd: '',
    cid: '',
    countLimit: 0,
    inCntLimit: 0,
    dlyCntLimit: 0,
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    ctlStatus: 0,
    uid: 0,
  }
}
// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    formData.value.status = Number(formData.value.status)
    formData.value.ctlStatus = Number(formData.value.ctlStatus)
    let res
    switch (type.value) {
      case 'create':
        res = await createChannelAccount(formData.value)
        break
      case 'update':
        res = await updateChannelAccount(formData.value)
        break
      case 'transfer':
        console.log("我tm要转移了", formData.value)
        res = await transferChannelForAcc(formData.value)
        break
      default:
        res = await createChannelAccount(formData.value)
        break
    }
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '创建/更改成功'
      })
      closeDialog()
      getTableData()
    }
  });
}

// 精准/模糊控制
const caCtlInfo = ref({
  status: 1,
  id: '',
})

const switchCtlEnable = async (row) => {
  console.log(row)
  caCtlInfo.value = JSON.parse(JSON.stringify(row))
  await nextTick()
  const req = {
    ...caCtlInfo.value
  }
  const res = await updateChannelAccount(req)
  if (res.code === 0) {
    ElMessage({type: 'success', message: `开启${req.status === 1 ? '模糊控制' : '精准控制'}成功`})
    await getTableData()
  }
}

// 通道账号开关（批量）
const switchOnVisible = ref(false)
const switchOffVisible = ref(false)
// 通道账号开关
const caInfo = ref({
  status: 1,
  id: '',
})
// 批量ca data
const switchData = ref({
  ids: [],
  status: 0,
})

const switchEnable = async (row) => {
  console.log(row)
  caInfo.value = JSON.parse(JSON.stringify(row))
  await nextTick()
  const req = {
    ...caInfo.value
  }
  const res = await switchEnableCA(req)
  if (res.code === 0) {
    ElMessage({type: 'success', message: `${req.status === 0 ? '禁用' : '启用'}成功`})
    await getTableData()
  }
}

// 批量开启
const onSwitchEnable = async () => {
  const ids = []
  if (multipleSelection.value.length === 0) {
    ElMessage({
      type: 'warning',
      message: '请选择要开启的数据'
    })
    return
  }

  multipleSelection.value &&
  multipleSelection.value.map(item => {
    ids.push(item.ID)
  })
  switchData.value.ids = ids
  switchData.value.status = 1
  const req = {
    ...switchData.value
  }
  const res = await switchEnableCAByIds(req)
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '开启成功'
    })
    if (tableData.value.length === ids.length && page.value > 1) {
      page.value--
    }
    switchOnVisible.value = false
    getTableData()
  }
}

//批量关闭
const onSwitchDisable = async () => {
  const ids = []
  if (multipleSelection.value.length === 0) {
    ElMessage({
      type: 'warning',
      message: '请选择要关闭的数据'
    })
    return
  }
  multipleSelection.value &&
  multipleSelection.value.map(item => {
    ids.push(item.ID)
  })
  switchData.value.ids = ids
  switchData.value.status = 0
  const req = {
    ...switchData.value
  }
  const res = await switchEnableCAByIds(req)
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '关闭成功'
    })
    if (tableData.value.length === ids.length && page.value > 1) {
      page.value--
    }
    switchOffVisible.value = false
    getTableData()
  }
}

// 通道账号token更新
const dialogTokenFormVisible = ref(false)

const caTokenInfo = ref({
  token: '',
  id: '',
})

const updateCaTokenFunc = async (row) => {
  const res = await findChannelAccount({ID: row.ID})
  type.value = 'update'
  typeTitle.value = '修改'
  if (res.code === 0) {
    formData.value = res.data.revca
    dialogTokenFormVisible.value = true
  }
}

const updTokenInfo = async (row) => {
  console.log(row)
  caTokenInfo.value = JSON.parse(JSON.stringify(row))
  await nextTick()
  const req = {
    ...caTokenInfo.value
  }
  const res = await updateTokenInfoFunc(req)
  if (res.code === 0) {
    ElMessage({type: 'success', message: `更新CK${req.code === 0 ? '成功' : '失败'}`})
    await getTableData()
  }
}

const updateTokenInfoFunc = async (row) => {
  const res = await findChannelAccount({ID: row.ID})
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.revca
    dialogTokenFormVisible.value = true
  }
}

// 充值记录查询
const orderHisVisible = ref(false)
const orderHis2000Visible = ref(false)
const orderHis4000Visible = ref(false)
const orderHis5000Visible = ref(false)
const orderHis8000Visible = ref(false)
const orderHis9000Visible = ref(false)
// 系统记录查询
const orderSysVisible = ref(false)
// 产码统计查询
const payCodeOverviewVisible = ref(false)
const channelCode = ref("")
const orderHisTableData = ref([])
const orderHis4000TableData = ref([])
const orderHis5000TableData = ref([])
const orderHis8000Info = ref([])
const orderHis8000List = ref([])
const orderHis9000Info = ref([])
const orderHis9000List = ref([])
const orderHis2000Info = ref([])
const orderHis2000List = ref([])
const orderSysTableData = ref([])
const payCodeTableData = ref([])
const payCodeMap = ref({})
const closeOrderSysShow = () => {
  orderSysVisible.value = false
  orderSysTableData.value = []
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
const closeOrderHis5000Show = () => {
  orderHis5000Visible.value = false
  orderHis5000TableData.value = []
}
const closeOrderHis8000Show = () => {
  orderHis8000Visible.value = false
  orderHis8000List.value = []
  orderHis8000Info.value = {}
}
const closeOrderHis9000Show = () => {
  orderHis9000Visible.value = false
  orderHis9000List.value = []
  orderHis9000Info.value = {}
}
const closePayCodeOverviewShow = () => {
  payCodeOverviewVisible.value = false
  payCodeTableData.value = []
}
const openOrderHisShow = async (row) => {
  let req = {...row}
  console.log(req)
  let cid = req.cid
  if (cid >= 3000 && cid <= 3099) {
    orderHisVisible.value = true;
  } else if (cid >= 1000 && cid <= 1099) {
    orderHisVisible.value = true;
  } else if (cid >= 1100 && cid <= 1199) {
    orderHisVisible.value = true;
  } else if (cid >= 1200 && cid <= 1299) {
    orderHisVisible.value = true;
  } else if (cid >= 2000 && cid <= 2099) {
    orderHis2000Visible.value = true;
  } else if (cid >= 4000 && cid <= 4099) {
    orderHis4000Visible.value = true;
  } else if (cid >= 5000 && cid <= 5099) {
    orderHis5000Visible.value = true;
  } else if (cid >= 8000 && cid <= 8099) {
    orderHis8000Visible.value = true;
  } else if (cid >= 9000 && cid <= 9099) {
    orderHis9000Visible.value = true;
  }

  await queryAccOrderHisFunc(req, cid)
}

const openOrderSysShow = async (row) => {
  orderSysVisible.value = true
  let req = {...row}
  console.log(req)
  await queryAccOrderSysFunc(req)
}
const openPayCodeOverviewShow = async (row) => {
  payCodeOverviewVisible.value = true
  let req = {...row}
  console.log(req)
  await getPayCodeOverviewByChanAccFunc({acId: req.acId, codeStatus: 2})
}
const router = useRouter()
const goWorkLog = async () => {
  router.push({name: 'vbRecord', replace: true})
}

//  产码 ---------------------

const numHours = ref(0)
const numMinutes = ref(0)
const numSeconds = ref(0)
const activeNames = ref(['1'])
const cardBackgroundColor = (remaining) => {
  // if (remaining === 0) return '#d30404'; // 红色
  if (remaining === 0) return '#909399'; // 红色
  else if (remaining >= 1 && remaining <= 3) return '#d30404'; // 橙色
  else if (remaining >= 4 && remaining <= 10) return '#ec7b0b'; // 黄色
  else return '#4a9b22'; // 绿色
};
const default7Day = async () => {
  numHours.value = 24 * 7
  numMinutes.value = 0
  numSeconds.value = 0
}

const default1Day = async () => {
  numHours.value = 24
  numMinutes.value = 0
  numSeconds.value = 0
}

const default2Hour = async () => {
  numHours.value = 2
  numMinutes.value = 0
  numSeconds.value = 0
}

const default1Hour = async () => {
  numHours.value = 1
  numMinutes.value = 0
  numSeconds.value = 0
}

const default10Minute = async () => {
  numHours.value = 0
  numMinutes.value = 10
  numSeconds.value = 0
}

const default0Second = async () => {
  numHours.value = 0
  numMinutes.value = 0
  numSeconds.value = 0
}
const pcElFormRef = ref()

const pcFormData = ref({
  cid: '',
  acId: '',
  acAccount: '',
  acRemark: '',
  type: 2,
  codeStatus: 0,
  expTime: '',
  list: [
    {
      // operator: '',
      // location: '',
      // locList: '',
      imgBaseStr: '',
      imgContent: '',
      address: '',
      money: 0,
      status: 0
    }
  ]
})
const pcDialogFormVisible = ref(false)

const createByChannelPayCodeFunc = async (row) => {
  const res = await findChannelAccount({ID: row.ID})
  let vca = res.data.revca
  if (vca.sysStatus !== 1) {
    ElMessage({
      type: 'error',
      message: '该通道账号未经过系统审核开启，不允许创建产码，请核查账号情况'
    })
    return
  }
  type.value = 'createPc'
  typeTitle.value = '添加产码'
  if (res.code === 0) {
    pcFormData.value = {
      cid: vca.cid,
      acId: vca.acId,
      acAccount: vca.acAccount,
      acRemark: vca.acRemark,
      type: 2,
      expTime: '',
      list: [
        {
          // operator: '',
          // location: '',
          // locList: '',
          imgBaseStr: '',
          imgContent: '',
          address: '',
          money: 0,
          status: 0
        }
      ]
    }
    pcDialogFormVisible.value = true
  }
}
let activeUpdIndex = ref(-1);
// 新增行
const handleAdd2Upd = function () {
  let item = {
    // operator: '',
    // location: '',
    imgBaseStr: '',
    money: 0,
  };
  pcFormData.value.list.push(item);
  activeUpdIndex.value = pcFormData.value.list.length - 1;
};
// 编辑行
const handleEdit2Upd = (index) => {
  activeUpdIndex.value = index;
};
// 保存行
const handleSave2Upd = () => {
  let create = {...pcFormData.value}
  let newList = []
  let ele = pcFormData.value.list[activeUpdIndex.value];
  // ele.location = ele.locList[0]
  pcFormData.value.list[activeUpdIndex.value] = ele
  newList.push(ele)
  create.list = newList
  console.log(pcFormData.value.list)
  activeUpdIndex.value = -1;
};

// 删除行
const handleDelete2Upd = function (index) {
  let ele = pcFormData.value.list[index];
  console.log(ele)
  let id = ele.id;
  if (id) {
    console.log("有id，要删库 -> id: " + id)
  } else {
    console.log("没id的临时数据，随便删")
  }
  pcFormData.value.list.splice(index, 1);
};
const getIntervalTime = async () => {
  const now = new Date()
  let expirationTime = new Date(now.getTime() + numHours.value * 60 * 60 * 1000)
  expirationTime = new Date(expirationTime.getTime() + numMinutes.value * 60 * 1000)
  expirationTime = new Date(expirationTime.getTime() + numSeconds.value * 1000)
  let intervalTime = dayjs(expirationTime).tz('Asia/Shanghai');
  console.log('intervalTime', intervalTime)
  // pcFormData.value.expTime = intervalTime.format('YYYY-MM-DD HH:mm:ss')
  pcFormData.value.expTime = new Date(intervalTime)
  // console.log('expTime', intervalTime)
  return expirationTime
}
const img_base_str = ref('')
const pcFileList = ref([]);

const pcImgList = ref([]);
const pcDialogImgVisible = ref(false);
const pcDialogImageUrls = ref("");

const uploadImgToBase64 = (file) => {
  // 核心方法，将图片转成base64字符串形式
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = function () {
      // 图片转base64完成后返回reader对象
      resolve(reader);
    };
    reader.onerror = reject;
  });
}

const getFiles = async (file, fileList) => {
  const isLt2M = file.size / 1024 / 1024 < 2;
  if (isLt2M) {
    try {
      const data = await uploadImgToBase64(file.raw);
      // img_base_str.value = data.result;
      pcImgList.value.push(data.result);

    } catch (error) {
      console.error(error);
    }
  } else {
    ElMessage({
      type: 'error',
      message: '上传图片大小不能超过 2MB!'
    })
  }


  console.log("file-111", JSON.stringify(file));
  console.log("fileList-111", JSON.stringify(fileList));
  console.log("list-111", JSON.stringify(pcImgList));
  // formData.value.imgBaseStr=img_base_str.value
}
const handlePicRemoves = (file, fileList) => {
  let hideUploadEdit = fileList.length
  if (hideUploadEdit >= 1) {
    img_base_str.value = "";
  }
};
const handlePicPreviews = (file) => {
  console.log('file=' + file.url);
  pcDialogImageUrls.value = file.url;
  pcDialogImgVisible.value = true;
}
// ------------获取省市 -------
const selectedCity = ref([]);
const regionOptions = ref([])
const optionsRegion = regionData;
const chge = () => {
  const lastElement = selectedCity.value[selectedCity.value.length - 1]
  pcFormData.value.location = lastElement
  console.log(selectedCity);
};

const handleChangeH = (value) => {
  console.log('h:', value)
  numHours.value = value
}
const operators = [
  {
    value: 'dianxin',
    label: '电信',
  },
  {
    value: 'yidong',
    label: '移动',
  },
  {
    value: 'liantong',
    label: '联通',
  }
]
const handleChangeM = (value) => {
  console.log('m:', value)
  numMinutes.value = value
}

const handleChangeS = (value) => {
  console.log('s:', value)
  numSeconds.value = value
}
const closePcDialog = () => {
  pcDialogFormVisible.value = false
  pcFormData.value = {
    cid: '',
    acId: '',
    acAccount: '',
    acRemark: '',
    type: 2,
    codeStatus: 0,
    expTime: '',
    list: [
      {
        // operator: '',
        // location: '',
        // locList: '',
        imgBaseStr: '',
        imgContent: '',
        address: '',
        money: 0,
        status: 0
      }
    ]
  }
  pcImgList.value = []
}
// 弹窗确定
const enterPcDialog = async () => {
  await getIntervalTime()
  pcElFormRef.value?.validate(async (valid) => {
        // console.log('formData' + JSON.stringify(formData.value))
        if (!valid) return
        let res
        switch (type.value) {
          case 'createPc':
            let flag = false;
            let td = pcFormData.value;
            if (!td.list || td.list.length === 0) {
              ElMessage({
                type: 'error',
                message: '至少上传一个码'
              })
              flag = true;
              break
            } else {
              for (let i = 0; i < pcFormData.value.list.length; i++) {
                let item = pcFormData.value.list[i]
                if (item.money <= 0) {
                  ElMessage({
                    type: 'error',
                    message: '金额需大于0'
                  })
                  flag = true;
                  break
                }
                if (!item.imgBaseStr) {
                  ElMessage({
                    type: 'error',
                    message: '传入正确报文'
                  })
                  flag = true;
                  break
                }
                /*if(!item.location){
                  ElMessage({
                    type: 'error',
                    message: '传入正确地区'
                  })
                  flag = true;
                  break
                }
                if (!item.operator) {
                  ElMessage({
                    type: 'error',
                    message: '传入正确运营商'
                  })
                  flag = true;
                  break
                }*/
              }
            }
            if (flag) {
              return
            } else {
              pcFormData.value.type = 2
              console.log(pcFormData.value)
              res = await batchCreateChannelPayCode(pcFormData.value)
            }
            break
          default:
            break
        }
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: '创建/更改成功'
          })
          closePcDialog()
          getTableData()
        }
      }
  )
}
//  产码 ---------------------

// 产码选择通道
// 弹窗控制标记
const dialogChanFormVisible = ref(false)

// 关闭弹窗
const closeChanDialog = () => {
  dialogChanFormVisible.value = false
}

const elChanFormRef = ref()
/*const enterChanDialog = async () => {
  elChanFormRef.value?.validate(async (valid) => {
    if (!valid) return
    console.log(formData.value.cid)
    dialogChanFormVisible.value = false
    let channelCode = Number(formData.value.cid);
    let type = Number(formData.value.type);

    if (channelCode >= 3000 && channelCode < 3099) {
      dialog3000FormVisible.value = true
    } else if (channelCode >= 1000 && channelCode < 1099) {
      dialog1000FormVisible.value = true
    } else if (channelCode >= 1100 && channelCode < 1199) {
      dialog1100FormVisible.value = true
    } else if (channelCode >= 2000 && channelCode < 2099) {
      dialog2000FormVisible.value = true
    } else if (channelCode >= 4000 && channelCode < 4099) {
      dialog4000FormVisible.value = true
    }
  });
}*/


// ---------- 消费历史 ----------
// const searchAccInfo = ref({})
const searchAccInfo = ref({
  toUid: '',
  username: '',
  acId: '',
  acAccount: '',
  channelCode: '',
})
// 重置
const onAccReset = () => {
  searchAccInfo.value = {}
  getAccTableData()
}

const getAccTableData = async () => {
  console.log(searchAccInfo.value)
  const voRes = await getOrderAccOverview({...searchAccInfo.value});
  console.log(voRes.data)
  if (voRes.code === 0) {
    showCostOrderAccTitle.value = `订单核算(用户归属:${searchAccInfo.value.username})`
    costOrderAccTable.value = voRes.data.list;
  }
}

// 搜索
const onAccSubmit = async () => {
  getAccTableData()
}

const showCostOrderAccVisible = ref(false)
const showCostOrderAccTitle = ref()
let costOrderAccTable = ref([]);
const showCostOrderAcc = async (row) => {
  if (row != '') {
    searchAccInfo.value.toUid = row.CreatedBy
    const res = await findChannelAccount({ID: row.ID})
    if (res.code === 0) {
      formData.value = res.data.revca
    }
    searchAccInfo.value.username = formData.value.username
    searchAccInfo.value.acAccount = row.acAccount
    searchAccInfo.value.acId = row.acId
  }
  costOrderAccTable.value = [];
  const voRes = await getOrderAccOverview({...searchAccInfo.value});
  console.log(voRes.data)
  if (voRes.code === 0) {
    if (row != '') {
      showCostOrderAccTitle.value = `订单核算(用户归属:${searchAccInfo.value.username})`
    } else {
      showCostOrderAccTitle.value = `订单核算`
    }
    costOrderAccTable.value = voRes.data.list;
    showCostOrderAccVisible.value = true
    searchAccInfo.value = {}
  }
}
const closeCostOrderAcc = () => {
  showCostOrderAccVisible.value = false
  costOrderAccTable.value = [];
}
// ---------- 消费历史 ---------

// ---------- 通道卡片 ---------
// 获取所有父节点
// const parentNodes = getParentNodes(list);
const parentNodes = ref([])

function getParentNodes(nodes) {
  const result = [];
  for (const node of nodes) {
    if (node.parentId === "0") {
      result.push(node);
    }
  }
  return result;
}

function shouldShowParent(node) {
  return !(node.children && node.children.length > 0);
}

function handleProdClick(node) {
  console.log("handleProdClick" + node)
  formData.value.cid = node.channelCode

  dialogChanFormVisible.value = false
  let channelCode = Number(formData.value.cid);
  console.log(formData.value)

  if (channelCode >= 3000 && channelCode < 3099) {
    formData.value.type = 2
    dialog3000FormVisible.value = true
    loginQr()

  } else if (channelCode >= 1000 && channelCode < 1099) {
    formData.value.type = 1
    dialog1000FormVisible.value = true
    loginQr()
  } else if (channelCode >= 1100 && channelCode < 1199) {
    formData.value.type = 1
    dialog1100FormVisible.value = true
    loginQr()

  } else if (channelCode >= 1200 && channelCode < 1299) {
    formData.value.type = 1
    dialog1200FormVisible.value = true
    loginQr()

  } else if (channelCode >= 2000 && channelCode < 2099) {
    formData.value.type = 1
    dialog2000FormVisible.value = true
  } else if (channelCode >= 4000 && channelCode < 4099) {
    formData.value.type = 1
    dialog4000FormVisible.value = true
  } else if (channelCode >= 5000 && channelCode < 5099) {
    formData.value.type = 1
    dialog5000FormVisible.value = true
  } else if (channelCode >= 6000 && channelCode < 6099) {
    formData.value.type = 4
    dialog6000FormVisible.value = true
  } else if (channelCode >= 8000 && channelCode < 8099) {
    formData.value.type = 1
    dialog8000FormVisible.value = true
  } else if (channelCode >= 9000 && channelCode < 9099) {
    formData.value.type = 1
    dialog9000FormVisible.value = true
  }   else if (channelCode >= 10000 && channelCode < 10099) {
    formData.value.type = 2
    dialog10000FormVisible.value = true
  }
}

const accCustomStyle = ref({
  background: 'linear-gradient(to right, #2ecc71, #3498db)',
  color: '#FFF',
  height: '120px',
})
// ---------- 通道卡片 ---------

// ---------- 通道转移 ---------
const dialogTransferFormVisible = ref(false)

const transferAccFunc = async (row) => {
  const res = await findChannelAccount({ID: row.ID})
  typeTitle.value = '通道转移'
  type.value = 'transfer'
  const vcpTable = await getChannelProductSelf({page: 1, pageSize: 999, type: formData.value.type})

  if (vcpTable.code === 0) {
    let retList = vcpTable.data.list
    // console.log(retList)
    let cidNum = Number(row.cid);
    let prodList = []
    let cidFilter;
    if (cidNum >= 3000 && cidNum < 3099) {
      cidFilter = 3000
    } else if (cidNum >= 1000 && cidNum < 1099) {
      cidFilter = 1000
    } else if (cidNum >= 1100 && cidNum < 1199) {
      cidFilter = 1000
    } else if (cidNum >= 1200 && cidNum < 1299) {
      cidFilter = 1200
    } else if (cidNum >= 2000 && cidNum < 2099) {
      cidFilter = 2000
    } else if (cidNum >= 4000 && cidNum < 4099) {
      cidFilter = 4000
    }
    for (let i = 0; i < retList.length; i++) {
      if (cidFilter === Number(retList[i].channelCode)) {
        prodList.push(retList[i])
      }
    }
    console.log(prodList)
    vcpTableData.value = prodList
    channelCodeOptions.value = []
    setTransferOptions(vcpTableData.value, channelCodeOptions.value, cidNum)
    console.log(channelCodeOptions.value)
  }

  if (res.code === 0) {
    formData.value = res.data.revca
    console.log("转移前", formData.value)
    dialogTransferFormVisible.value = true
  }
}


const setTransferOptions = (ChannelCodeData, optionsData, cidNum) => {
  ChannelCodeData &&
  ChannelCodeData.forEach(item => {
    let flag = false
    if (Number(item.channelCode) === cidNum) {
      flag = true
    }
    if (item.children && item.children.length) {
      const option = {
        value: item.channelCode + '',
        label: item.productName,
        disabled: flag,
        children: []
      }
      setTransferOptions(
          item.children,
          option.children,
          cidNum
      )
      optionsData.push(option)
    } else {
      const option = {
        value: item.channelCode + '',
        label: item.productName,
        disabled: flag,
      }
      optionsData.push(option)
    }
  })
}
// ---------- 通道转移 ---------
// ---------- 卡密 ---------
const handleTokenInput = () => {
  const token = formData.value.token;
  const acAccount = formData.value.acAccount
  if (!acAccount) {
    const pinValue = findValuesContainingString(token, "pin");
    if (pinValue) {
      formData.value.acAccount = pinValue;
    }
  }
}
// ---------- 卡密 ---------
</script>

<style lang="scss" scoped>
.region-card-container {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  width: 100%;
}

.region-card {
  margin: 10px;
  width: 250px;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  text-align: center;
  color: white;
  position: relative;
  transition: transform 0.3s;
}

.region-card:hover {
  transform: translateY(-5px);
}

.region-tag {
  position: absolute;
  top: 10px;
  left: 10px;
  color: white;
  padding: 5px 10px;
  border-radius: 5px;
}

.region-name {
  margin: 0;
  font-size: 14px;
}

.region-code {
  margin: 0;
  font-size: 12px;
}


.region-title {
  margin-bottom: 10px;
}

.region-title h2 {
  font-size: 20px;
  color: white;
}

.region-title p {
  font-size: 14px;
}

.region-business-data {
  display: flex;
  justify-content: space-around;
  margin-top: 15px;
}

.region-data-item {
  flex: 1;
}

.region-label {
  font-size: 14px;
}

.region-value {
  padding-top: 5px;
  font-size: 18px;
  font-weight: bold;
}

.tab {
  margin-bottom: 20px;
}

.tab h2 {
  cursor: pointer;
  padding: 10px;
  background-color: #ccc;
  margin: 0;
}

.card-container {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
}

.card {
  width: 260px;
}

.indicator {
  display: flex;
  justify-content: space-around; // 使子元素水平居中展开
  padding: 15px;
  border-radius: 8px; // 添加圆角
}

.indicator span {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 10px; // 调整间距

  &:not(:last-child) {
    border-right: 2px solid #fff; // 白色边框
    margin-right: 10px; // 调整间距
  }
}

.acc-container {
  color: #FFFFFF;
}

.label {
  color: #F5F5F5;
  font-size: 16px;
}

.value {
  color: #FFFFFF;
  font-size: 16px;
  font-weight: bold;
  margin-top: 15px; // 调整间距
}

.scrolling-text {
  height: 30px; /* 设置显示区域的高度 */
  overflow: hidden; /* 隐藏超出显示区域的内容 */
  position: relative; /* 设置为相对定位，以便在其中添加绝对定位的子元素 */
}

.scrolling-text ul {
  list-style-type: none; /* 移除列表默认样式 */
  padding: 0;
  margin: 0;
  animation: scroll-text 4s linear infinite; /* 使用动画实现滚动效果，20s表示滚动完成需要的时间，可根据需要调整 */
}

@keyframes scroll-text {
  0% {
    transform: translateY(0); /* 初始位置在顶部 */
  }
  100% {
    transform: translateY(-100%); /* 最终位置在顶部的上方，根据行数和行高进行计算 */
  }
}

</style>

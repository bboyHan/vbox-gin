<template>
  <div>
    <div class="gva-table-box">
      <el-row :gutter="16">
        <div v-for="item in countItem">
          <el-col :span="16">
            通道ID： {{ item.cid }}
            <el-button color="#05411d">
              <el-icon class="is-loading">
                <Loading/>
              </el-icon>
              已开启数量： {{ item.total }}
            </el-button>
          </el-col>
        </div>
      </el-row>
    </div>
    <el-divider/>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule"
               @keyup.enter="onSubmit">
        <el-form-item label="通道账户" prop="acAccount">
          <el-input v-model="searchInfo.acAccount" placeholder="搜索条件"/>
        </el-form-item>
        <el-form-item label="账户备注" prop="acRemark">
          <el-input v-model="searchInfo.acRemark" placeholder="搜索条件"/>
        </el-form-item>
        <el-form-item label="账户id" prop="acId">
          <el-input v-model.number="searchInfo.acId" placeholder="搜索条件"/>
        </el-form-item>
        <el-form-item label="通道id" prop="cid">
          <el-cascader
              v-model="searchInfo.cid"
              :options="channelCodeOptions"
              :props="channelCodeProps"
              @change="handleChange"
              style="width: 100%"
              placeholder="选择通道"
          />
        </el-form-item>
        <el-form-item label="开关状态" prop="status">
          <el-select v-model="searchInfo.status" placeholder="选择状态">
            <el-option label="已开启" value="1"/>
            <el-option label="已关闭" value="0"/>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
        <el-popover v-model:visible="deleteVisible" placement="top" width="160">
          <p>确定要删除吗？</p>
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
        <el-popover v-model:visible="switchOnVisible" placement="top" width="160">
          <p>确定批量开启吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button type="primary" link @click="switchOnVisible = false">取消</el-button>
            <el-button type="primary" @click="onSwitchEnable">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length"
                       @click="switchOnVisible = true">批量开启
            </el-button>
          </template>
        </el-popover>
        <el-popover v-model:visible="switchOffVisible" placement="top" width="160">
          <p>确定批量关闭吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button type="primary" link @click="switchOffVisible = false">取消</el-button>
            <el-button type="primary" @click="onSwitchDisable">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length"
                       @click="switchOffVisible = true">批量关闭
            </el-button>
          </template>
        </el-popover>
      </div>
      <el-table ref="multipleTable" tooltip-effect="dark" :data="tableData" row-key="ID" border resizable="true"
                @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55"/>
        <el-table-column align="left" label="通道id" prop="cid" width="80"/>
        <el-table-column align="left" label="账户备注" prop="acRemark" width="120"/>
        <el-table-column align="left" label="通道账户" prop="acAccount" width="120"/>
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
        <el-table-column align="center" label="日限额" prop="dailyLimit" width="120"/>
        <el-table-column align="center" label="总限额" prop="totalLimit" width="120"/>
        <el-table-column align="center" label="笔数限额" prop="countLimit" width="120"/>
        <el-table-column align="left" label="状态 / 系统开关" prop="status" width="140">
          <template #default="scope">
            <el-row :gutter="12">
              <el-col :span="12">
                <el-switch v-model="scope.row.status" inline-prompt :active-value="1" active-text="开启"
                           :inactive-value="0" inactive-text="关闭" size="large"
                           @change="()=>{switchEnable(scope.row)}"/>
              </el-col>
              <el-col :span="12">
                <el-switch v-model="scope.row.sysStatus" inline-prompt :active-value="1" active-text="开启"
                           :inactive-value="0" inactive-text="关闭" size="large" disabled/>
              </el-col>
            </el-row>
          </template>
        </el-table-column>
        <el-table-column align="center" label="查询" width="180">
          <template #default="scope">
            <el-row :gutter="12">
              <el-col :span="12">
                <el-button type="primary" link class="table-button" @click="openOrderHisShow(scope.row)">充值记录</el-button>
              </el-col>
              <el-col :span="12">
                <el-button type="primary" link class="table-button" @click="openOrderSysShow(scope.row)">系统记录</el-button>
              </el-col>
            </el-row>
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作" width="120">
          <template #default="scope">
            <div v-if="Number(scope.row.cid) === 3000">
              <el-row>
                <el-col :span="24">
                  <el-button type="primary" link class="table-button" @click="createByChannelPayCodeFunc(scope.row)">
                    产码
                  </el-button>
                  <el-button type="primary" link icon="info-filled" class="table-button"
                             @click="openPayCodeOverviewShow(scope.row)"></el-button>
                </el-col>
              </el-row>
            </div>
            <div v-else>
              <span>非预产通道</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作" width="120">
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

    <!--  创建  -->
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :title="typeTitle" destroy-on-close>
      <el-scrollbar height="500px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
          <el-form-item label="账户备注" prop="acRemark">
            <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
          </el-form-item>
          <el-row>
            <el-col :span="12">
              <el-form-item label="通道账户" prop="acAccount">
                <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="账户密码" prop="acPwd">
                <el-input v-model="formData.acPwd" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="产码方式" prop="type">
                <el-radio-group v-model="formData.type" @change="handleChange">
                  <el-radio label="1">
                    <template #default><span>引导</span></template>
                  </el-radio>
                  <el-radio label="2">
                    <template #default><span>预产</span></template>
                  </el-radio>
                  <el-radio label="3">
                    <template #default><span>原生</span></template>
                  </el-radio>
                </el-radio-group>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="通道id" prop="cid">
                <el-cascader v-model="formData.cid" :options="channelCodeOptions" :props="channelCodeProps" @change=""
                             style="width: 100%"/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="token" prop="token">
            <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入"/>
          </el-form-item>
          <el-row>
            <el-col :span="6"></el-col>
            <el-col :span="12">
              <warning-bar title="注：默认0，则无限额控制"/>
            </el-col>
            <el-col :span="6"></el-col>
            <el-col :span="8">
              <el-form-item label="日限额" prop="dailyLimit">
                <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限额" prop="totalLimit">
                <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="笔数限额" prop="countLimit">
                <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="状态开关" prop="status">
            <el-switch v-model="formData.status" active-value="1" inactive-value="0" active-text="开启"
                       inactive-text="关闭"></el-switch>
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

    <!--  修改  -->
    <el-dialog v-model="dialogUpdFormVisible" :before-close="closeDialog" :title="typeTitle" destroy-on-close>
      <el-scrollbar height="500px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
          <el-form-item label="账户备注" prop="acRemark">
            <el-input v-model="formData.acRemark" :clearable="true" placeholder="请输入"/>
          </el-form-item>
          <el-row>
            <el-col :span="12">
              <el-form-item label="通道账户" prop="acAccount">
                <el-input v-model="formData.acAccount" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="账户密码" prop="acPwd">
                <el-input v-model="formData.acPwd" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item label="通道id" prop="cid">
                <el-input v-model="formData.cid" :clearable="true" placeholder="请输入" disabled/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="token" prop="token">
            <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入"/>
          </el-form-item>
          <el-row>
            <el-col :span="6"></el-col>
            <el-col :span="12">
              <warning-bar title="注：默认0，则无限额控制"/>
            </el-col>
            <el-col :span="6"></el-col>
            <el-col :span="8">
              <el-form-item label="日限额" prop="dailyLimit">
                <el-input v-model.number="formData.dailyLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="总限额" prop="totalLimit">
                <el-input v-model.number="formData.totalLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="笔数限额" prop="countLimit">
                <el-input v-model.number="formData.countLimit" :clearable="true" placeholder="请输入"/>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="状态开关" prop="status">
            <el-switch v-model="formData.status" active-value="1" inactive-value="0" active-text="开启"
                       inactive-text="关闭"></el-switch>
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

    <!--  CK  -->
    <el-dialog v-model="dialogTokenFormVisible" :before-close="closeDialog" :title="变更CK"
               destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="token" prop="token">
          <el-input v-model="formData.token" type="textarea" :clearable="true" placeholder="请输入"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 查看详情 -->
    <el-dialog v-model="detailShow" style="width: 800px" lock-scroll :before-close="closeDetailShow" title="查看详情"
               destroy-on-close>
      <el-scrollbar height="550px">
        <el-descriptions :column="6" border>
          <el-descriptions-item label="用户id" :span="6">{{ formData.CreatedBy }}</el-descriptions-item>
          <el-descriptions-item label="账户ID" :span="6">{{ formData.acId }}</el-descriptions-item>
          <el-descriptions-item label="账户备注" :span="6">{{ formData.acRemark }}</el-descriptions-item>
          <el-descriptions-item label="通道账户" :span="3">{{ formData.acAccount }}</el-descriptions-item>
          <el-descriptions-item label="账户密码" :span="3">{{ formData.acPwd }}</el-descriptions-item>
          <el-descriptions-item label="ck" :span="6">
            <el-input v-model="formData.token" type="textarea" readonly/>
          </el-descriptions-item>
          <el-descriptions-item label="通道id" :span="6">{{ formData.cid }}</el-descriptions-item>
          <el-descriptions-item label="笔数限制" :span="2">{{ formData.countLimit }}</el-descriptions-item>
          <el-descriptions-item label="日限额" :span="2">{{ formData.dailyLimit }}</el-descriptions-item>
          <el-descriptions-item label="总限额" :span="2">{{ formData.totalLimit }}</el-descriptions-item>
          <el-descriptions-item label="状态开关" :span="3">{{
              formData.status === 0 ? '关闭' : '开启'
            }}
          </el-descriptions-item>
          <el-descriptions-item label="系统开关" :span="3">{{
              formData.sysStatus === 0 ? '关闭' : '开启'
            }}
          </el-descriptions-item>
        </el-descriptions>
      </el-scrollbar>
    </el-dialog>

    <!-- 查看充值详情 -->
    <el-dialog v-model="orderHisVisible" style="width: 1100px" lock-scroll :before-close="closeOrderHisShow"
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

    <!-- 查看产码详情 -->
    <el-dialog v-model="payCodeOverviewVisible" style="width: 1100px" lock-scroll
               :before-close="closePayCodeOverviewShow" title="查看产码详情" destroy-on-close>
      <el-scrollbar height="550px">
        <div class="region-card-container">
          <div v-for="pcData in Object.entries(payCodeMap)" style="width: 100%">
            <!--            <div>￥{{ pcData[0].x1 }}，{{ formatOPSimple(pcData.x2) }}, {{ codeToText[pcData.x3] }}({{ pcData.x4 }})</div>-->
            <div>{{ formatOPSimple(pcData[0]) }}</div>
            <el-divider></el-divider>
            <span v-for="pcDetail in pcData[1]" style="padding: 10px">
              <el-badge :value="pcDetail.x4">
                <el-button>{{ codeToText[pcDetail.x3] }} | {{ pcDetail.x1 }}元</el-button>
              </el-badge>
            </span>
          </div>
        </div>
      </el-scrollbar>
    </el-dialog>

    <!-- 产码-->
    <el-dialog width="60%" v-model="pcDialogFormVisible" :before-close="closePcDialog" :title="typeTitle"
               destroy-on-close>
      <el-scrollbar height="450px">
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
                <el-table-column label="运营商" prop="operator" width="120px">
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
                </el-table-column>
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
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closePcDialog">取 消</el-button>
          <el-button type="primary" @click="enterPcDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 查询指定账户订单 -->
    <el-dialog v-model="orderSysVisible" style="width: 1100px" lock-scroll :before-close="closeOrderSysShow"
               title="查看系统充值详情" destroy-on-close>
      <el-scrollbar height="450px">
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
          <el-table
              v-if="isSimple"
              ref="multipleTable"
              style="width: 100%"
              tooltip-effect="dark"
              :data="orderSysTableData"
              row-key="ID"
              border
              @selection-change="handleSelectionChange"
          >
            <el-table-column align="center" label="账号ID" prop="acId" width="180"/>
            <el-table-column align="center" label="订单ID" prop="orderId" width="220"/>
            <el-table-column align="center" label="金额" prop="money" width="120"/>
            <el-table-column align="center" label="订单状态" prop="orderStatus" width="120">
              <template #default="scope">
                <el-button style="width: 90px" :color="formatPayedColor(scope.row.orderStatus, scope.row.acId)">
                  {{ formatPayed(scope.row.orderStatus, scope.row.acId) }}
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
            <el-table-column align="left" label="操作" width="240">
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
          <el-table
              v-else
              ref="multipleTable"
              style="width: 100%"
              tooltip-effect="dark"
              :data="orderSysTableData"
              row-key="ID"
              border
              @selection-change="handleSelectionChange"
          >
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
            <el-table-column align="left" label="订单ID" prop="orderId" width="220"/>
            <el-table-column align="left" label="金额" prop="money" width="120"/>
            <el-table-column align="left" label="订单状态" prop="orderStatus" width="120">
              <template #default="scope">
                <el-button style="width: 90px" :color="formatPayedColor(scope.row.orderStatus, scope.row.acId)">
                  {{ formatPayed(scope.row.orderStatus, scope.row.acId) }}
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
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closePcDialog">取 消</el-button>
          <el-button type="primary" @click="enterPcDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 订单查看详情 -->
    <el-dialog v-model="sysDetailShow" style="width: 800px" lock-scroll :before-close="closeSysDetailShow" title="查看详情"
               destroy-on-close>
      <el-scrollbar height="550px">
        <el-descriptions column="1" border>
          <el-descriptions-item label="订单ID">{{ sysFormData.orderId }}</el-descriptions-item>
          <el-descriptions-item label="付方ID">{{ sysFormData.pAccount }}</el-descriptions-item>
          <el-descriptions-item label="金额">{{ sysFormData.money }}</el-descriptions-item>
          <el-descriptions-item label="单价积分">{{ sysFormData.unitPrice }}</el-descriptions-item>
          <el-descriptions-item label="账号ID">{{ sysFormData.acId }}</el-descriptions-item>
          <el-descriptions-item label="通道编码">{{ sysFormData.channelCode }}</el-descriptions-item>
          <el-descriptions-item label="平台ID">{{ sysFormData.platId }}</el-descriptions-item>
          <el-descriptions-item label="客户ip">{{ sysFormData.payIp }}</el-descriptions-item>
          <el-descriptions-item label="区域">{{ sysFormData.payRegion }}</el-descriptions-item>
          <el-descriptions-item label="客户端设备">{{ sysFormData.payDevice }}</el-descriptions-item>
          <el-descriptions-item label="订单状态">{{ formatBoolean(sysFormData.orderStatus) }}</el-descriptions-item>
          <el-descriptions-item label="回调状态">{{ formatBoolean(sysFormData.cbStatus) }}</el-descriptions-item>
          <el-descriptions-item label="回调时间">{{ formatDate(sysFormData.cbTime) }}</el-descriptions-item>
        </el-descriptions>
      </el-scrollbar>
    </el-dialog>

    <!--  补单  -->
    <el-dialog
        v-model="dialogSysFormVisible"
        :before-close="closeSysDialog"
        :title="typeTitle"
        destroy-on-close
        style="width: 450px"
    >
      <el-scrollbar height="100px">
        <el-form :model="sysFormData" label-position="right" ref="elSysFormRef" label-width="120px">
          <el-form-item label="订单ID" prop="authCaptcha">
            <el-input disabled v-model="sysFormData.orderId" :clearable="true" placeholder="请输入" style="width: 80%"/>
          </el-form-item>
          <el-form-item label="安全码" prop="authCaptcha">
            <el-input v-model="sysFormData.authCaptcha" :clearable="true" placeholder="请输入安全码" style="width: 80%"/>
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
} from '@/api/channelAccount'
import {
  getChannelProductSelf
} from '@/api/channelProduct'
import {codeToText, regionData} from 'element-china-area-data';

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
  formatPayed, formatNotify, formatHandNotifyColor, formatHandNotify
} from '@/utils/format'
import {ElMessage, ElMessageBox} from 'element-plus'
import {ref, reactive, nextTick} from 'vue'
import WarningBar from "@/components/warningBar/warningBar.vue";
import {Delete, Edit, Eleme, InfoFilled, Loading, Plus, Position, Select} from "@element-plus/icons-vue";
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
import {callback2Pa, findPayOrder, getPayOrderList} from "@/api/payOrder";

defineOptions({
  name: 'ChannelAccount'
})

// 注册插件
dayjs.extend(utcPlugin);
dayjs.extend(timezone);

const countItem = ref([])
const queryAccOrderHisFunc = async (row) => {
  const req = {...row}
  console.log(req)

  let res = await queryAccOrderHis(req)
  console.log(res.data)
  if (res.code === 0) {
    orderHisTableData.value = res.data.list.WaterList
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
  queryAccOrderSysFunc()
}

// 系统查单修改页面容量
const handleSysCurrentChange = (val) => {
  sysPage.value = val
  queryAccOrderSysFunc()
}
const queryAccOrderSysFunc = async (row) => {
  const req = {...row}
  console.log(req)

  let res = await getPayOrderList({page: page.value, pageSize: pageSize.value, acId: req.acId})
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
  totalLimit: 0,
  type: 0,
  status: 0,
  sysStatus: 0,
  uid: 0,
})

// 验证规则
const rule = reactive({
  acAccount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
    {
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
  type: [{
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
const cascaderRules = reactive({
  acAccount: [
    { required: true, message: '请选择', trigger: ['input', 'blur', 'change'], }
  ],
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
      setChannelCodeOptions(
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
const dialogUpdFormVisible = ref(false)
const updateChannelAccountFunc = async (row) => {
  const res = await findChannelAccount({ID: row.ID})
  type.value = 'update'
  typeTitle.value = '修改'
  if (res.code === 0) {
    formData.value = res.data.revca
    dialogUpdFormVisible.value = true
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
const dialogFormVisible = ref(false)

// 系统查单补单
const dialogSysFormVisible = ref(false)
// 关闭弹窗
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

// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findChannelAccount({ID: row.ID})
  if (res.code === 0) {
    formData.value = res.data.revca
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
    dailyLimit: 0,
    totalLimit: 0,
    status: 0,
    sysStatus: 0,
    uid: 0,
  }
}

// 打开弹窗
const openDialog = () => {
  type.value = 'create'
  typeTitle.value = '创建'
  dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  dialogUpdFormVisible.value = false
  dialogTokenFormVisible.value = false
  formData.value = {
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
// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    formData.value.status = Number(formData.value.status)
    let res
    switch (type.value) {
      case 'create':
        res = await createChannelAccount(formData.value)
        break
      case 'update':
        res = await updateChannelAccount(formData.value)
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
  })
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
// 系统记录查询
const orderSysVisible = ref(false)
// 产码统计查询
const payCodeOverviewVisible = ref(false)
const channelCode = ref("")
const orderHisTableData = ref([])
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
const closePayCodeOverviewShow = () => {
  payCodeOverviewVisible.value = false
  payCodeTableData.value = []
}
const openOrderHisShow = async (row) => {
  orderHisVisible.value = true
  let req = {...row}
  console.log(req)
  await queryAccOrderHisFunc(req)
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
  await getPayCodeOverviewByChanAccFunc({acId:req.acId, codeStatus: 2})
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
      operator: '',
      location: '',
      locList: '',
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
          operator: '',
          location: '',
          locList: '',
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
    operator: '',
    location: '',
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
  ele.location = ele.locList[0]
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
        operator: '',
        location: '',
        locList: '',
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
            }else {
              for (let i = 0; i < pcFormData.value.list.length; i++) {
                let item = pcFormData.value.list[i]
                if(item.money <= 0) {
                  ElMessage({
                    type: 'error',
                    message: '金额需大于0'
                  })
                  flag = true;
                  break
                }
                if (!item.imgBaseStr){
                  ElMessage({
                    type: 'error',
                    message: '传入正确报文'
                  })
                  flag = true;
                  break
                }
                if(!item.location){
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
                }
              }
            }
            if (flag) {
              return
            }else {
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

</script>

<style scoped>
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
</style>

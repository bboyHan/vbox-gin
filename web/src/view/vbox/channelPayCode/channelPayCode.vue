<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule"
               @keyup.enter="onSubmit">
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
        <el-form-item label="通道账户名" prop="acAccount">
          <el-input v-model="searchInfo.acAccount" placeholder="搜索条件"/>
        </el-form-item>
        <el-form-item label="地区" prop="location">
          <el-cascader
              style="width:100%"
              :options="regionOptions"
              v-model="searchInfo.location"
              @change="chge"
              placeholder="地区"
          >
          </el-cascader>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="searchInfo.codeStatus" placeholder="选择状态">
            <el-option label="已使用" value="1"/>
            <el-option label="待使用" value="2"/>
            <el-option label="已失效" value="3"/>
          </el-select>
        </el-form-item>
        <el-form-item label="运营商" prop="status">
          <el-select v-model="searchInfo.operator" placeholder="选择ISP">
            <el-option label="移动" value="yidong"/>
            <el-option label="联通" value="liantong"/>
            <el-option label="电信" value="dianxin"/>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button icon="refresh" @click="onReset"></el-button>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button type="primary" icon="info-filled" class="table-button" @click="openPayCodeOverviewShow">概览
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-search-box">
      <div class="gva-btn-list">
<!--        <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>-->
        <el-button type="primary" icon="plus" @click="openBatchDialog">新增</el-button>
        <el-popover v-model:visible="deleteVisible" :disabled="!multipleSelection.length" placement="top" width="160">
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
      </div>
    </div>
    <div class="gva-table-box">
      <el-table
          ref="multipleTable"
          style="width: 100%"
          tooltip-effect="dark"
          :data="tableData"
          row-key="ID"
          @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55"/>
        <el-table-column align="left" label="创建日期" width="160">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="通道ID" prop="cid" width="80"/>
        <el-table-column align="left" label="通道账户名" prop="acAccount" width="140"/>
        <el-table-column align="left" label="备注" prop="acRemark" width="160"/>
        <el-table-column align="left" label="平台ID" prop="platId" width="360"/>
<!--        <el-table-column align="left" label="过期时间" prop="expTime" width="160">
          <template #default="scope">{{ formatDate(scope.row.expTime) }}</template>
        </el-table-column>-->
        <el-table-column align="left" label="剩余时间" prop="expTime" width="140">
          <template #default="scope">
            <span v-if="countdowns[scope.$index] > 0">{{ formatTime(countdowns[scope.$index]) }} </span>
            <span v-else-if="countdowns[scope.$index] <= 0 && scope.row.codeStatus === 1"> 0 /（已使用）</span>
            <span v-else> -1 /（已过期）</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="运营商" prop="operator" width="70">
          <template #default="{ row }">
            {{ getOperatorChinese(row.operator) }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="地区" prop="location" width="180">
          <template #default="{ row }">
            {{ regionCode2Text(row.location, 2) }} | {{ regionCode2Text(row.location, 4) }} |
            {{ regionCode2Text(row.location, 6) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" prop="codeStatus" width="90">
          <template #default="scope">
            <el-button style="width: 60px" :color="formatPayCodeColor(scope.row.codeStatus)">
              {{ formatPayCodeStatus(scope.row.codeStatus) }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column align="left" label="金额" prop="money" width="80">
        </el-table-column>
<!--        <el-table-column align="left" label="付款码" prop="imgContent" width="120">
          <template #default="{ row }">
            <div v-if="!dialogImageShow[row.ID]">
              <el-button link icon="search" @click="toggleDialog(row.ID, row.imgContent)">预览</el-button>
            </div>
            <div v-else>
              <el-button link icon="search" @click="toggleDialog(row.ID, row.imgContent)">取消预览</el-button>
              <el-image :src="row.imgBaseStr" fit="contain" class="thumbnail-image"/>
            </div>
          </template>
        </el-table-column>-->
        <el-table-column align="left" label="操作" width="120">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
              <el-icon style="margin-right: 5px">
                <InfoFilled/>
              </el-icon>
            </el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)"></el-button>
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

    <div class="gva-search-box">

      <div class="region-card-container">
        <!-- 查看产码详情 -->
        <el-dialog v-model="payCodeOverviewVisible" style="width: 1100px" lock-scroll
                   :before-close="closePayCodeOverviewShow" title="查看产码详情" destroy-on-close>
          <el-scrollbar height="550px">
            <div class="gva-search-box">
              <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline"
                       :rules="searchRule"
                       @keyup.enter="onSubmit">
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
                <el-form-item label="通道账户名" prop="acAccount">
                  <el-input v-model="searchInfo.acAccount" placeholder="搜索条件"/>
                </el-form-item>
                <el-form-item label="金额" prop="acAccount">
                  <el-input v-model="searchInfo.money" placeholder="搜索条件"/>
                </el-form-item>
                <el-form-item label="地区" prop="location">
                  <el-cascader
                      style="width:100%"
                      :options="regionOptions"
                      v-model="searchInfo.location"
                      @change="chge"
                      placeholder="地区"
                  >
                  </el-cascader>
                </el-form-item>
                <el-form-item label="状态" prop="status">
                  <el-select v-model="searchInfo.codeStatus" placeholder="选择状态">
                    <el-option label="已使用" value="1"/>
                    <el-option label="待使用" value="2"/>
                    <el-option label="已失效" value="3"/>
                  </el-select>
                </el-form-item>
                <el-form-item label="运营商" prop="status">
                  <el-select v-model="searchInfo.operator" placeholder="选择ISP">
                    <el-option label="移动" value="yidong"/>
                    <el-option label="联通" value="liantong"/>
                    <el-option label="电信" value="dianxin"/>
                  </el-select>
                </el-form-item>
                <el-form-item>
                  <el-button icon="refresh" @click="onReset"></el-button>
                  <el-button type="primary" class="table-button" @click="openPayCodeOverviewShow">查看</el-button>
                </el-form-item>
              </el-form>
            </div>
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

      </div>
      <!--      <el-collapse v-model="activeNames" @change="">
              <el-collapse-item title="预产统计视图（点击可收缩）" name="1">
                <div class="region-card-container">
                  <div
                      v-for="province in allProvinces"
                      :key="province.code"
                      class="region-card"
                      :style="{ backgroundColor: cardBackgroundColor(province.remaining) }"
                  >
                    <div @click="openDialog">
                      <div class="region-tag">
                        <p class="region-code">地区编码：{{ province.code }}</p>
                      </div>
                      <div class="region-title">
                        <h2>{{ province.name }}</h2>
                      </div>
                      <div class="region-business-data">
                        <div class="region-data-item">
                          <div class="region-label">待使用数</div>
                          <div class="region-value">{{ province.total }}</div>
                        </div>
                        <div class="region-data-item">
                          <div class="region-label">冷却中</div>
                          <div class="region-value">{{ province.used }}</div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </el-collapse-item>
            </el-collapse>-->

    </div>

    <!--  批量创建  -->
    <el-dialog width="60%" v-model="dialogBatchFormVisible" :before-close="closeBatchDialog" :title="typeTitle"
               destroy-on-close>
      <el-scrollbar height="500px">
        <el-form :model="batchFormData" label-position="right" ref="elBatchFormRef" :rules="batchRule"
                 label-width="80px">
          <el-form-item label="产码方式" prop="type">
            <el-button v-model="batchFormData.type" type="primary">预产</el-button>
          </el-form-item>
          <el-form-item label="通道" prop="cid">
            <el-cascader
                v-model="batchFormData.cid"
                :options="channelCodeOptions"
                :props="channelCodeProps"
                @change="handleChange"
                style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="通道账户" prop="acId">
            <el-select
                v-model="batchFormData.acId"
                placeholder="请选择通道账号"
                filterable
                clearable
                style="width: 100%"
                @change="handleAccChange"
            >
              <el-option
                  v-for="item in accList"
                  :key="item.acAccount"
                  :label="formatJoin(' -- 备注： ', item.acAccount, item.acRemark)"
                  :value="item.acId"
              />
            </el-select>
          </el-form-item>
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
              <el-table :data="batchFormData.list" style="width: 100%">
                <el-table-column label="报文" prop="imgBaseStr" style="width: 100%">
                  <template #default="scope">
                    <el-input :rows="2" type="textarea" v-if="activeUpdIndex === scope.$index"
                              v-model="scope.row.imgBaseStr"></el-input>
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
                              :parser="(value) => value.replace(/￥\s?|(,*)/g, '')">
                    </el-input>
                    <span v-else>￥{{ scope.row.money }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="运营商" prop="operator" width="120px">
                  <template #default="scope">
                    <el-select v-model="scope.row.operator" placeholder="请选择通信商" filterable style="width: 100%">
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
                        :props="{checkStrictly: false}"
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
          <el-button @click="closeBatchDialog">取 消</el-button>
          <el-button type="primary" @click="enterBatchDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!--  传码  -->
<!--    <el-dialog width="40%" v-model="dialogFormVisible" :before-close="closeDialog" :title="typeTitle" destroy-on-close>
      <el-scrollbar height="450px">
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
          <el-form-item label="产码方式" prop="type">
            <el-button v-model="formData.type" type="primary">预产</el-button>
          </el-form-item>
          <el-form-item label="通道" prop="cid">
            <el-cascader
                v-model="formData.cid"
                :options="channelCodeOptions"
                :props="channelCodeProps"
                @change="handleChange"
                style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="通道账户" prop="acId">
            <el-select
                v-model="formData.acId"
                placeholder="请选择通道账号"
                filterable
                clearable
                style="width: 100%"
                @change="handleAccChange"
            >
              <el-option
                  v-for="item in accList"
                  :key="item.acAccount"
                  :label="formatJoin(' &#45;&#45; 备注： ', item.acAccount, item.acRemark)"
                  :value="item.acId"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="金额" prop="money">
            <el-input v-model.number="formData.money"
                      placeholder="输入金额"
                      :formatter="(value) => `￥ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')"
                      :parser="(value) => value.replace(/￥\s?|(,*)/g, '')">
            </el-input>
          </el-form-item>
          <el-form-item label="平台ID" prop="platId">
            <el-input v-model="formData.platId" placeholder="输入平台ID"></el-input>
          </el-form-item>
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
          <el-row>
            <el-col :span="12">
              <el-form-item label="运营商" prop="operator">
                <el-select v-model="formData.operator" placeholder="请选择通信商" filterable style="width: 100%">
                  <el-option v-for="item in operators" :key="item.value" :label="item.label" :value="item.value"/>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="地区" prop="location">
                <el-cascader
                    :change-on-select="true"
                    style="width:100%"
                    :options="regionOptions"
                    v-model="selectedCity"
                    @change="chge"
                    placeholder="选择地区"
                    filterable
                    :props="{checkStrictly: false}"
                >
                </el-cascader>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="图片上传">
            <el-upload
                class="avatar-uploader"
                action=""
                :on-change="getFiles"
                :on-remove="handlePicRemoves"
                :on-preview="handlePicPreviews"
                v-model="lists"
                :limit="8"
                list-type="picture-card"
                :file-list="fileList"
                :auto-upload="false"
                accept="image/png, image/gif, image/jpg, image/jpeg"
                drag
                multiple
            >
              &lt;!&ndash; 图标 &ndash;&gt;
              <el-icon style="font-size: 25px;">
                <Plus/>
              </el-icon>
              <template #tip>
                <div class="el-upload__tip">
                  拖拽或点击上传
                </div>
              </template>
            </el-upload>
            <el-dialog v-model="dialogVisibles" title="预览" destroy-on-close>
              <img :src="dialogImageUrs" style="display: block;max-width: 500px;margin: 0 auto;height: 500px;" alt=""/>
            </el-dialog>
          </el-form-item>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>-->

    <el-dialog v-model="detailShow" style="width: 800px" lock-scroll :before-close="closeDetailShow" title="查看详情"
               destroy-on-close>
      <el-scrollbar height="550px">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="通道id" :span="1">
            {{ formData.cid }}
          </el-descriptions-item>
          <el-descriptions-item label="通道账户名" :span="1">
            {{ formData.acAccount }}
          </el-descriptions-item>
          <el-descriptions-item label="金额" :span="1">
            ￥{{ formData.money }}
          </el-descriptions-item>
          <el-descriptions-item label="过期时间" :span="1">
            {{ formatDate(formData.expTime) }}
          </el-descriptions-item>
          <el-descriptions-item label="运营商" :span="1">
            {{ getOperatorChinese(formData.operator) }}
          </el-descriptions-item>
          <el-descriptions-item label="地区" :span="1">
            {{ codeToText[formData.location] }}
          </el-descriptions-item>
          <el-descriptions-item label="状态" :span="2">
            {{ formatPayCodeStatus(formData.codeStatus) }}
          </el-descriptions-item>
          <el-descriptions-item label="平台ID" :span="2">
            {{ formData.platId }}
          </el-descriptions-item>
          <el-descriptions-item label="预处理" :span="2">
            <img :src="qrcodeUrl" alt="QR Code" style="height: 200px"/>
          </el-descriptions-item>
<!--          <el-descriptions-item label="付款码">
            <el-image :src="formData.imgBaseStr" fit="contain" class="thumbnail-image"/>
          </el-descriptions-item>-->
        </el-descriptions>
      </el-scrollbar>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  batchCreateChannelPayCode,
  createChannelPayCode,
  deleteChannelPayCode,
  deleteChannelPayCodeByIds,
  findChannelPayCode,
  getChannelPayCodeList,
  getPayCodeOverview
} from '@/api/channelPayCode'
import {getChannelProductSelf} from '@/api/channelProduct'
import {getChannelAccountList} from '@/api/channelAccount'

// 全量引入格式化工具 请按需保留
import {
  formatDate,
  formatJoin,
  formatOPSimple,
  formatPayCodeColor,
  formatPayCodeStatus,
  formatTime,
} from '@/utils/format'
import {ElMessage, ElMessageBox} from 'element-plus'
import {onMounted, reactive, ref} from 'vue'
import {codeToText, regionData} from 'element-china-area-data';
import {Delete, Edit, InfoFilled, Plus, Select} from '@element-plus/icons-vue';
import dayjs from 'dayjs';
import utcPlugin from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';
import provinces from '@/assets/json/provinces.json'
import QRCode from "qrcode";

defineOptions({
  name: 'ChannelPayCode'
})

// 批量 ------------------
let activeUpdIndex = ref(-1);
// 新增行
const handleAdd2Upd = function () {
  let item = {
    operator: '',
    location: '',
    imgBaseStr: '',
    money: 0,
  };
  batchFormData.value.list.push(item);
  activeUpdIndex.value = batchFormData.value.list.length - 1;
};
// 编辑行
const handleEdit2Upd = (index) => {
  activeUpdIndex.value = index;
};
// 保存行
const handleSave2Upd = () => {
  let create = {...batchFormData.value}
  let newList = []
  let ele = batchFormData.value.list[activeUpdIndex.value];
  ele.location = ele.locList[0]
  batchFormData.value.list[activeUpdIndex.value] = ele
  newList.push(ele)
  create.list = newList
  console.log(batchFormData.value.list)
  activeUpdIndex.value = -1;
};

// 删除行
const handleDelete2Upd = function (index) {
  let ele = batchFormData.value.list[index];
  console.log(ele)
  let id = ele.id;
  if (id) {
    console.log("有id，要删库 -> id: " + id)
  } else {
    console.log("没id的临时数据，随便删")
  }
  batchFormData.value.list.splice(index, 1);
};
// 批量 ------------------

// 注册插件
dayjs.extend(utcPlugin);
dayjs.extend(timezone);

const payCodeTableData = ref([])
const payCodeOverviewVisible = ref(false)
const closePayCodeOverviewShow = () => {
  payCodeOverviewVisible.value = false
  payCodeTableData.value = []
}
const openPayCodeOverviewShow = async () => {
  payCodeOverviewVisible.value = true
  let req = {...searchInfo.value}
  console.log(req)
  await getPayCodeOverviewByChanAccFunc(req)
}
const getPayCodeOverviewByChanAccFunc = async (row) => {
  const req = {...row}
  console.log(req)

  let res = await getPayCodeOverview(req)
  console.log(res.data)
  if (res.code === 0) {
    payCodeTableData.value = res.data.list
    payCodeMap.value = payCodeTableData.value.reduce((acc, cur) => {
      const {x2, ...rest} = cur;
      acc[x2] = acc[x2] || [];
      acc[x2].push(rest);
      return acc;
    }, {})
  }
}

const pcData = ref([]);
const payCodeMap = ref({});

// 缩略图
const dialogImageShow = ref({})
const qrcodeUrl = ref('');

const toggleDialog = (id, content) => {
  console.log(id)
  dialogImageShow.value[id] = !dialogImageShow.value[id];
};

const batchFormData = ref({
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

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  cid: '',
  acId: '',
  acAccount: '',
  acRemark: '',
  expTime: '',
  operator: '',
  location: '',
  imgBaseStr: '',
  imgContent: '',
  imgShowContent: '',
  mid: '',
  type: 2,
  codeStatus: 0,
  money: 0,
})

// 验证规则
const batchRule = reactive({
  acAccount: [{required: true, message: '', trigger: ['blur'],}],
  cid: [{required: true, message: '请选择', trigger: ['blur'],}],
  acId: [{required: true, message: '请选择', trigger: ['input', 'blur'],}],
  expTime: [{required: true, validator: validateTimeLimit, trigger: 'blur',},],
})
// 验证规则
const rule = reactive({
  acAccount: [{required: true, message: '', trigger: ['blur'],}],
  cid: [{required: true, message: '请选择', trigger: ['blur'],}],
  platId: [{required: true, message: '请选择', trigger: ['blur'],}],
  acId: [{required: true, message: '请选择', trigger: ['input', 'blur'],}],
  expTime: [{required: true, validator: validateTimeLimit, trigger: 'blur',},],
  operator: [{required: true, message: '', trigger: ['blur'],}],
  location: [{required: true, message: '请选择地区', trigger: ['input', 'blur'],}],
  money: [{validator: checkMoney, trigger: 'blur'}]
})

function checkMoney(rule, value, callback) {
  if (Number(value) <= 0) {
    callback(new Error('请输入正确的金额'));
  } else {
    callback();
  }
}

function checkExpirationTime(rule, value, callback) {
  if (new Date(value).getTime() <= new Date().getTime()) {
    callback(new Error('过期时刻要大于当前时间'));
  } else {
    callback();
  }
}

function validateTimeLimit(rule, value, callback) {
  if (numHours.value === 0 && numMinutes.value === 0 && numSeconds.value === 0) {
    callback(new Error('过期时间填写不能都为 0'));
  } else {
    callback();
  }
}


function regionCode2Text(content, index) {
  if (content.length < index) {
    return ' - '
  } else {
    return codeToText[content.slice(0, index)]
  }
}

function getOperatorChinese(operator) {
  const operatorMap = {
    liantong: '联通',
    yidong: '移动',
    dianxin: '电信'
  };
  return operatorMap[operator] || operator;
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
// --  获取过期时间


const numHours = ref(0)
const numMinutes = ref(0)
const numSeconds = ref(0)

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

const getIntervalTime = async () => {
  const now = new Date()
  let expirationTime = new Date(now.getTime() + numHours.value * 60 * 60 * 1000)
  expirationTime = new Date(expirationTime.getTime() + numMinutes.value * 60 * 1000)
  expirationTime = new Date(expirationTime.getTime() + numSeconds.value * 1000)
  let intervalTime = dayjs(expirationTime).tz('Asia/Shanghai');
  console.log('intervalTime', intervalTime)
  // let expTime = intervalTime.format('YYYY-MM-DD HH:mm:ss')
  formData.value.expTime = new Date(intervalTime)
  batchFormData.value.expTime = new Date(intervalTime)
  // console.log('expTime', intervalTime)
  return expirationTime
}


const handleChangeH = (value) => {
  console.log('h:', value)
  numHours.value = value
}

const handleChangeM = (value) => {
  console.log('m:', value)
  numMinutes.value = value
}

const handleChangeS = (value) => {
  console.log('s:', value)
  numSeconds.value = value
}


// -----------上传图片--------------
const img_base_str = ref('')
const dialogImageUrs = ref("");
const fileList = ref([]);

const dialogVisibles = ref(false);
const lists = ref([]);

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
      lists.value.push(data.result);

    } catch (error) {
      console.error(error);
    }
  } else {
    ElMessage({
      type: 'error',
      message: '上传图片大小不能超过 2MB!'
    })
  }

  // formData.value.imgBaseStr=img_base_str.value
}
// const getFiles = async (file, fileList) => {
//   const isLt2M = file.size / 1024 / 1024 < 2;
//   if (isLt2M) {
//     try {
//       const data = await uploadImgToBase64(file.raw);
//       img_base_str.value = data.result;
//       lists.value= data.result;

//     } catch (error) {
//       console.error(error);
//     }
//   } else {
//     ElMessage({
//                   type: 'error',
//                   message: '上传图片大小不能超过 2MB!'
//                 })
//   }
//   // console.log("fileList-111", fileList);
//   // console.log("file-111", file);
//   formData.value.imgBaseStr=img_base_str.value
// }

const handlePicRemoves = (file, fileList) => {
  let hideUploadEdit = fileList.length
  if (hideUploadEdit >= 1) {
    img_base_str.value = "";
  }
};

const handlePicPreviews = (file) => {
  dialogImageUrs.value = file.url;
  dialogVisibles.value = true;
}


// ------------获取省市 -------
const selectedCity = ref([]);
const optionsRegion = regionData;
const chge = () => {
  const lastElement = selectedCity.value[selectedCity.value.length - 1]
  formData.value.location = lastElement
  console.log(selectedCity);
};


// --------- 获取通信商 -----------
const operatorSelect = ref('')
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

// ------- 获取通道账号 -------
const accList = ref([])
const acIdList = ref([])
const sysUserAcId = ref('')
const selectCid = ref('')
const handleAccChange = (value) => {
  // console.log(value)
  // getACCChannelAccountByAcid()
  getALlChannelAccount()

}
// 获取唯一通道账号
const getACCChannelAccountByAcid = async () => {
  const res = await getChannelAccountList({acId: formData.value.acId, page: 1, pageSize: 999})
  acIdList.value = res.data.list
  // console.log(JSON.stringify(accList))
  formData.value.acAccount = acIdList.value[0].acAccount
  formData.value.acRemark = acIdList.value[0].acRemark
  console.log(JSON.stringify(formData.value))
  return res
}


// 获取通道账号
const getALlChannelAccount = async () => {
  const res = await getChannelAccountList({cid: formData.value.cid, page: 1, pageSize: 999})
  accList.value = res.data.list
}

// -------------- 同一通道产品的归集 ------------------------

// region
const regionOptions = ref([])

//通道产品
const channelCodeOptions = ref([])
const vcpTableData = ref([])

const channelCodeProps = {
  expandTrigger: 'hover',
  checkStrictly: false,
  emitPath: false,
}

const handleChange = (value) => {
  console.log(value)
  getALlChannelAccount()
}

const handleOptChange = async (value) => {
  const vcpTable = await getChannelProductSelf({page: 1, pageSize: 999, type: formData.value.type})

  if (vcpTable.code === 0) {
    vcpTableData.value = vcpTable.data.list
    setOptions()
  }
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

const elFormRef = ref()
const elBatchFormRef = ref()
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
  console.log(page.value)
  getTableData()
}

// 查询
const getTableData = async () => {
  const table = await getChannelPayCodeList({page: page.value, pageSize: pageSize.value, ...searchInfo.value})
  const vcpTable = await getChannelProductSelf({page: 1, pageSize: 999, type: 2})
  vcpTableData.value = vcpTable.data.list
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
  setOptions()
  // getALlChannelAccount()

}
// setRegionOptions(provinces, regionOptions.value, false)

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () => {
  channelCodeOptions.value = []
  setChannelCodeOptions(vcpTableData.value, channelCodeOptions.value, false)
  // console.log(provinces)
  setRegionOptions(provinces, regionOptions.value, false)
}


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
    deleteChannelPayCodeFunc(row)
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
  const res = await deleteChannelPayCodeByIds({ids})
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
const updateChannelPayCodeFunc = async (row) => {
  const res = await findChannelPayCode({ID: row.ID})
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.rechannelPayCode
    dialogFormVisible.value = true
  }
}

// 删除行
const deleteChannelPayCodeFunc = async (row) => {
  const res = await deleteChannelPayCode({ID: row.ID})
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
const dialogBatchFormVisible = ref(false)

// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findChannelPayCode({ID: row.ID})
  if (res.code === 0) {
    formData.value = res.data.rechannelPayCode
    const content = formData.value.imgContent
    if (content) {
      QRCode.toDataURL(content)
          .then((dataUrl) => {
            qrcodeUrl.value = dataUrl
          })
          .catch((error) => {
            console.error('Failed to generate QR code:', error);
          });
    }
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
    cid: '',
    acId: '',
    acAccount: '',
    acRemark: '',
    expTime: '',
    operator: '',
    location: '',
    imgBaseStr: '',
    mid: '',
    codeStatus: 0,
    money: 0,
  }
}

// 打开弹窗
const openRegionDialog = (province) => {
  type.value = 'create'
  typeTitle.value = '创建'
  selectedCity.value = [];
  dialogFormVisible.value = true
  // getALlChannelAccount()
  formData.value = {
    cid: '',
    acId: '',
    acAccount: '',
    acRemark: '',
    expTime: '',
    operator: '',
    location: '',
    imgBaseStr: '',
    platId: '',
    mid: '',
    type: 2,
    codeStatus: 0,
    money: 0,
  }
}

// 打开弹窗
const openBatchDialog = () => {
  type.value = 'create'
  typeTitle.value = '批量创建'
  selectedCity.value = [];
  dialogBatchFormVisible.value = true
  // getALlChannelAccount()
  batchFormData.value = {
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
}
// 打开弹窗
const openDialog = () => {
  type.value = 'create'
  typeTitle.value = '创建'
  selectedCity.value = [];
  dialogFormVisible.value = true
  // getALlChannelAccount()
  formData.value = {
    cid: '',
    acId: '',
    acAccount: '',
    acRemark: '',
    expTime: '',
    operator: '',
    location: '',
    imgBaseStr: '',
    mid: '',
    type: 2,
    codeStatus: 0,
    money: 0,
  }
}

// 关闭弹窗
const closeBatchDialog = () => {
  dialogBatchFormVisible.value = false
  batchFormData.value = {
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
  lists.value = []
}
// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    cid: '',
    acId: '',
    acAccount: '',
    acRemark: '',
    expTime: '',
    operator: '',
    location: '',
    imgBaseStr: '',
    mid: '',
    type: 2,
    codeStatus: 0,
    money: 0,
  }
  lists.value = []
}
// 弹窗确定
const enterBatchDialog = async () => {
  await getIntervalTime()
  elBatchFormRef.value?.validate(async (valid) => {
        if (!valid) {
          return
        }
        console.log(batchFormData.value)
        let res
        switch (type.value) {
          case 'batchCreate':
            batchFormData.value.type = 2
            res = await batchCreateChannelPayCode(batchFormData.value)
            break
          default:
            res = await batchCreateChannelPayCode(batchFormData.value)
            break
        }
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: '创建成功'
          })
          closeBatchDialog()
          getTableData()
        }
      }
  )
}

const enterDialog = async () => {
  const accInfo = await getACCChannelAccountByAcid()
  await getIntervalTime()
  console.log('accInfo ' + JSON.stringify(accInfo.data.list))
  console.log('formData pre' + JSON.stringify(formData.value))

  elFormRef.value?.validate(async (valid) => {
        if (!valid) {
          return
        }

        if (lists.value.length < 1) {
          ElMessage({
            showClose: true,
            message: "请上传至少一个资源",
            type: 'error'
          })
          return
        }
        formData.value.money = Number(formData.value.money)

        if (formData.value.money < 1) {
          ElMessage({
            showClose: true,
            message: "请输入正确的金额",
            type: 'error'
          })
          return
        }
        let res
        switch (type.value) {
          case 'create':
            // console.log(">>>>>>" + lists.value.length)
            for (let i = 0; i < lists.value.length; i++) {

              // console.log('formData lists.value[i] ' + i + '  ' + JSON.stringify(lists.value[i]))
              formData.value.imgBaseStr = lists.value[i]
              // console.log('formData after ' + i + '  ' + JSON.stringify(formData.value))
              res = await createChannelPayCode(formData.value)
            }
            break
          default:
            res = await createChannelPayCode(formData.value)
            break
        }
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: '创建成功'
          })
          closeDialog()
          getTableData()
        }
      }
  )
}

// 倒计时数组
const countdowns = ref([]);

// 计算倒计时
const calculateCountdown = () => {
  setInterval(() => {
    const currentTime = new Date();
    tableData.value.forEach((item, index) => {
      const expTime = new Date(item.expTime);
      const timeDiffInSeconds = (expTime - currentTime) / 1000;
      countdowns.value[index] = timeDiffInSeconds > 0 ? Math.floor(timeDiffInSeconds) : -1;
    });
  }, 1000);
};

const activeNames = ref(['1'])
const allProvinces = ref([]);
const provincesD = ref([
  {name: '北京', code: '11', total: 1000, used: 500, remaining: 500},
  {name: '上海', code: '31', total: 800, used: 300, remaining: 0},
  {name: '河北', code: '13', total: 800, used: 300, remaining: 3},
  {name: '山西', code: '14', total: 800, used: 300, remaining: 8},
  {name: '河南', code: '41', total: 800, used: 300, remaining: 12},
])

const cardBackgroundColor = (remaining) => {
  // if (remaining === 0) return '#d30404'; // 红色
  if (remaining === 0) return '#909399'; // 红色
  else if (remaining >= 1 && remaining <= 3) return '#d30404'; // 橙色
  else if (remaining >= 4 && remaining <= 10) return '#ec7b0b'; // 黄色
  else return '#4a9b22'; // 绿色
};

onMounted(() => {
  calculateCountdown();

  const provincesData = provinces;
  allProvinces.value = provincesData.map((obj) => {
    const backendProvince = provincesD.value.find((province) => {
      // console.log('json code:' + obj.code + ' 后端数据:' + province.code)
      let pc = String(province.code).slice(0, 2);
      let bc = String(obj.code).slice(0, 2);
      // console.log('转换后- json code:' + pc + ' 后端数据:' + bc)
      return pc === bc;
    });
    return backendProvince || {code: obj.code, name: obj.name, total: 0, used: 0, remaining: 0};
  });
});


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

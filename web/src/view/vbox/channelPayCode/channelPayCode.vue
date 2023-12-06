<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
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
          <el-input v-model="searchInfo.acAccount" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="地区" prop="location">
          <el-cascader
              style="width:100%"
              :options="optionsRegion"
              v-model="searchInfo.location"
              @change="chge"
              placeholder="省 / 市 / 区"
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
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
        <el-popover v-model:visible="deleteVisible" :disabled="!multipleSelection.length" placement="top" width="160">
          <p>确定要删除吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button type="primary" link @click="deleteVisible = false">取消</el-button>
            <el-button type="primary" @click="onDelete">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="deleteVisible = true">删除</el-button>
          </template>
        </el-popover>
      </div>
      <el-table
          ref="multipleTable"
          style="width: 100%"
          tooltip-effect="dark"
          :data="tableData"
          row-key="ID"
          @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="创建日期" width="160">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="通道ID" prop="cid" width="80" />
        <el-table-column align="left" label="通道账户名" prop="acAccount" width="100" />
        <el-table-column align="left" label="备注" prop="acRemark" width="100" />
        <el-table-column align="left" label="过期时间" prop="expTime" width="160" />
        <el-table-column align="left" label="剩余时间" prop="expTime" width="140">
          <template #default="scope">
            <span v-if="countdowns[scope.$index] > 0">{{ formatTime(countdowns[scope.$index]) }} </span>
            <span v-else>-1 （已过期）</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="运营商" prop="operator" width="70" >
          <template #default="{ row }">
          {{ getOperatorChinese(row.operator) }}
        </template>
        </el-table-column>
        <el-table-column align="center" label="地区" prop="location" width="180">
          <template #default="{ row }">
            {{ codeToText[row.location.slice(0, 2)] }} | {{ codeToText[row.location.slice(0, 4)] }} | {{ codeToText[row.location.slice(0, 6)] }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" prop="codeStatus" width="90" >
          <template #default="scope">
            <el-button style="width: 60px" :color="formatPayCodeColor(scope.row.codeStatus)">
              {{ formatPayCodeStatus(scope.row.codeStatus) }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column align="left" label="金额" prop="money" width="80" >
        </el-table-column>
        <el-table-column align="left" label="付款码" prop="imgBaseStr" width="120" >
          <template #default="{ row }">
            <div v-if="!dialogImageShow[row.ID]">
              <el-button link icon="search" @click="toggleDialog(row.ID)" >预览</el-button>
            </div>
            <div v-else>
              <el-button link icon="search" @click="toggleDialog(row.ID)">取消预览</el-button>
              <el-image :src="row.imgBaseStr" fit="contain" class="thumbnail-image"/>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" >
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
              <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>
              详情
            </el-button>
            <el-button type="primary" link icon="edit" class="table-button" @click="updateChannelPayCodeFunc(scope.row)">变更</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
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
    <el-dialog width="30%" v-model="dialogFormVisible" :before-close="closeDialog" :title="type==='create'?'添加':'修改'" destroy-on-close>
      <el-scrollbar height="450px" >
        <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
          <el-form-item label="产码方式"  prop="type" >
<!--            <el-radio-group v-model="formData.type" @change="handleOptChange">
              <el-radio label="2" ><template #default><span>预产</span></template></el-radio>
            </el-radio-group>-->
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
          <el-form-item label="通道账户"  prop="acId" >
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
                  :label="formatJoin(' -- 备注： ', item.acAccount, item.acRemark)"
                  :value="item.acId"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="金额"  prop="money" >
              <el-input v-model.number="formData.money"
                        placeholder="输入金额"
                        :formatter="(value) => `￥ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')"
                        :parser="(value) => value.replace(/￥\s?|(,*)/g, '')">
              </el-input>
          </el-form-item>
          <el-form-item label="过期时间"  prop="expTime" >
            <!-- <el-date-picker
                v-model="formData.expTime"
                type="datetime"
                value-format="YYYY-MM-DD HH:mm:ss"
                style="width: 100%"
                align="center"
            >
            </el-date-picker> -->
            <el-input-number v-model="numHours" size="small" controls-position="right" @change="handleChangeH" :min="0">
            </el-input-number> 
            <span> 小时</span>
            <el-input-number v-model="numMinutes" size="small" controls-position="right" @change="handleChangeM" :min="0">
            </el-input-number> 
            <span> 分钟</span>
            <el-input-number v-model="numSeconds" size="small" controls-position="right" @change="handleChangeS" :min="0">
            </el-input-number>
            <span> 秒</span>
          </el-form-item>
          <el-form-item label="运营商"  prop="operator" >
            <el-select
                v-model="formData.operator"
                placeholder="请选择通信商"
                filterable
                style="width: 100%"
            >
              <el-option
                  v-for="item in operators"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="地区"  prop="location" >
            <el-cascader
                :change-on-select="true"
                style="width:100%"
                :options="optionsRegion"
                v-model="selectedCity"
                @change="chge"
                placeholder="省 / 市 / 区"
                filterable
                :props="{checkStrictly: true}"
            >
            </el-cascader>
          </el-form-item>
          <el-form-item label="图片上传" >
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
            >
              <!-- 图标 -->
              <el-icon style="font-size: 25px;"><Plus /></el-icon>
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
    </el-dialog>

    <el-dialog v-model="detailShow" style="width: 800px" lock-scroll :before-close="closeDetailShow" title="查看详情" destroy-on-close>
      <el-scrollbar height="550px">
        <el-descriptions column="1" border>
          <el-descriptions-item label="通道id">
            {{ formData.cid }}
          </el-descriptions-item>
          <el-descriptions-item label="通道账户名">
            {{ formData.acAccount }}
          </el-descriptions-item>
          <el-descriptions-item label="过期时间">
            {{ formData.expTime }}
          </el-descriptions-item>
          <el-descriptions-item label="运营商">
            {{ getOperatorChinese(formData.operator) }}
            <!-- <template #default="{ row }">
            {{ getOperatorChinese(row.operator) }}
            </template> -->
          </el-descriptions-item>
          <el-descriptions-item label="地区">
            {{ formData.location }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            {{ formatPayCodeStatus(formData.codeStatus) }}
          </el-descriptions-item>
          <el-descriptions-item label="付款码" >
            <el-image :src="formData.imgBaseStr" fit="contain" class="thumbnail-image"/>
          </el-descriptions-item>
        </el-descriptions>
      </el-scrollbar>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createChannelPayCode,
  deleteChannelPayCode,
  deleteChannelPayCodeByIds,
  updateChannelPayCode,
  findChannelPayCode,
  getChannelPayCodeList
} from '@/api/channelPayCode'
import {
  getChannelProductSelf
} from '@/api/channelProduct'
import {
  getChannelAccountList
} from '@/api/channelAccount'

// 全量引入格式化工具 请按需保留
import {
  getDictFunc,
  formatDate,
  formatBoolean,
  filterDict,
  ReturnArrImg,
  onDownloadFile,
  formatPayCodeStatus,
  formatJoin,
  formatPayCodeColor,
  formatTime,
} from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted } from 'vue'
import { codeToText, regionData } from 'element-china-area-data';
import { InfoFilled, Plus } from '@element-plus/icons-vue';
import dayjs from 'dayjs';
import utcPlugin from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';

defineOptions({
    name: 'ChannelPayCode'
})

// 注册插件
dayjs.extend(utcPlugin);
dayjs.extend(timezone);

// 缩略图
const dialogImageShow = ref({})
const toggleDialog = (id) => {
  console.log(id)
  dialogImageShow.value[id] = !dialogImageShow.value[id];
};

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
  mid: '',
  type: 2,
  codeStatus: 0,
  money: 0,
})

const accFormData = ref({
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
  uid: 0,
})

// 验证规则
const rule = reactive({
  acAccount: [{
    required: true,
    message: '',
    trigger: [ 'blur'],
  }
  ],
  cid: [{
    required: true,
    message: '请选择',
    trigger: [ 'blur'],
  }
  ],
  acId: [{
    required: true,
    message: '请选择',
    trigger: ['input', 'blur'],
  }
  ],
  // expTime: [{
  //   required: true,
  //   message: '过期时刻要大于当前时间',
  //   trigger: [ 'blur'],
  // },
  // { 
  //     validator: checkExpirationTime, 
  //     trigger: 'blur' 
  //   }
  // ],
  expTime: [
    {
      required: true,
      validator: validateTimeLimit,
      trigger: 'blur',
    },
  ],
  operator: [{
    required: true,
    message: '',
    trigger: ['blur'],
  }
  ],
  location: [{
    required: true,
    message: '请选择省或者省市',
    trigger: ['input', 'blur'],
  }
  ],
  // imgBaseStr: [{
  //   required: true,
  //   message: '',
  //   trigger: ['input', 'blur'],
  // }
  // ],
  money: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  }
  ]
  
})

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
    { validator: (rule, value, callback) => {
      if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
        callback(new Error('请填写结束日期'))
      } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
        callback(new Error('请填写开始日期'))
      } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
        callback(new Error('开始日期应当早于结束日期'))
      } else {
        callback()
      }
    }, trigger: 'change' }
  ],
})
// --  获取过期时间


const numHours = ref(0)
const numMinutes = ref(10)
const numSeconds = ref(0)

const getIntervalTime = async() => {
  const now = new Date()
  let expirationTime = new Date(now.getTime() + numHours.value * 60 * 60 * 1000)
  expirationTime = new Date(expirationTime.getTime() + numMinutes.value * 60 * 1000)
  expirationTime = new Date(expirationTime.getTime() + numSeconds.value * 1000)
  let intervalTime = dayjs(expirationTime).tz('Asia/Shanghai');
  console.log('intervalTime', intervalTime)
  formData.value.expTime = intervalTime.format('YYYY-MM-DD HH:mm:ss')
  console.log('expTime', intervalTime)
  return expirationTime
}


const handleChangeH = (value) => {
  console.log('h:',value)
  numHours.value = value
}

const handleChangeM = (value) => {
  console.log('m:',value)
  numMinutes.value = value
}

const handleChangeS = (value) => {
  console.log('s:',value)
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

    const getFiles = async (file ,fileList) => {
        const isLt2M = file.size / 1024 / 1024 < 2;
        if (isLt2M) {
          try {
            const data = await uploadImgToBase64(file.raw);
            // img_base_str.value = data.result;
            lists.value.push(data.result) ;

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
        console.log("list-111", JSON.stringify(lists));
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
  if (hideUploadEdit >= 1){
    img_base_str.value = "";
  } 
};

const handlePicPreviews = (file) => {
  console.log('file=' + file.url);
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
const getACCChannelAccountByAcid = async() => {
  const res = await getChannelAccountList({ acId: formData.value.acId ,page: 1, pageSize: 999})
  acIdList.value = res.data.list
  total.value = res.data.total
  // console.log(JSON.stringify(accList))
  formData.value.acAccount = acIdList.value[0].acAccount
  formData.value.acRemark = acIdList.value[0].acRemark
  console.log(JSON.stringify(formData.value))
  return res
}


// 获取通道账号
const getALlChannelAccount = async() => {
  const res = await getChannelAccountList({ cid: formData.value.cid ,page: 1, pageSize: 999})
  accList.value = res.data.list
  total.value = res.data.total
}

// -------------- 同一通道产品的归集 ------------------------


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
  elSearchFormRef.value?.validate(async(valid) => {
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
const getTableData = async() => {
  const table = await getChannelPayCodeList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  const vcpTable = await getChannelProductSelf({ page: 1, pageSize: 999, type: 2 })
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

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
  channelCodeOptions.value = []
  setChannelCodeOptions(vcpTableData.value, channelCodeOptions.value, false)
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
const onDelete = async() => {
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
      const res = await deleteChannelPayCodeByIds({ ids })
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

// 更新行
const updateChannelPayCodeFunc = async(row) => {
    const res = await findChannelPayCode({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data.rechannelPayCode
        dialogFormVisible.value = true
    }
}



// 更新行
const createByChannelPayCodeFunc = async(row) => {
    const res = await findChannelPayCode({ ID: row.ID })
    type.value = 'create'
    if (res.code === 0) {
        // formData.value = res.data.rechannelPayCode
        selectedCity.value = [];
        formData.value = {
          cid: res.data.rechannelPayCode.cid,
          acId: res.data.rechannelPayCode.acId,
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
          numHours.value = 0
          numMinutes.value = 10
          numSeconds.value = 0
          dialogFormVisible.value = true
    }
}

// 删除行
const deleteChannelPayCodeFunc = async (row) => {
    const res = await deleteChannelPayCode({ ID: row.ID })
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


// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findChannelPayCode({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.rechannelPayCode
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
const openDialog = () => {
    type.value = 'create'
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
          numHours.value = 0
          numMinutes.value = 10
          numSeconds.value = 0
    
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
    lists.value =[]
    numHours.value = 0
    numMinutes.value = 10
    numSeconds.value = 0
}
// 弹窗确定
const enterDialog = async () => {
  const accInfo = await getACCChannelAccountByAcid()
  getIntervalTime()
  console.log('accInfo ' + JSON.stringify(accInfo.data.list))
  console.log('formData pre' + JSON.stringify(formData.value))
  
  elFormRef.value?.validate( async (valid) => {
        // console.log('formData' + JSON.stringify(formData.value))
        if (!valid) return

        formData.value.money = Number(formData.value.money)

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
          case 'update':
            res = await updateChannelPayCode(formData.value)
            break
          default:
            res = await createChannelPayCode(formData.value)
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

onMounted(() => {
  calculateCountdown();
});
</script>

<style>

</style>

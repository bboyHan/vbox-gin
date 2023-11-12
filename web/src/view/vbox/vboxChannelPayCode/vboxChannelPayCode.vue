<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
      <el-form-item label="创建日期" prop="createdAt">
      <template #label>
        <span>
          创建日期
          <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
            <el-icon><QuestionFilled /></el-icon>
          </el-tooltip>
        </span>
      </template>
      <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始日期" :disabled-date="time=> searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"></el-date-picker>
       —
      <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束日期" :disabled-date="time=> searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"></el-date-picker>
      </el-form-item>
        <el-form-item label="通道id" prop="cid">
         <el-input v-model="searchInfo.cid" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="通道账户名" prop="acAccount">
         <el-input v-model="searchInfo.acAccount" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="省市" prop="location">
         <el-input v-model="searchInfo.location" placeholder="搜索条件" />

        </el-form-item>
        <!-- <el-form-item label="用户id" prop="uid">
            
             <el-input v-model.number="searchInfo.uid" placeholder="搜索条件" />

        </el-form-item> -->
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
        <el-table-column align="left" label="日期" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="通道ID" prop="cid" width="120" />
        <el-table-column align="left" label="通道账户名" prop="acAccount" width="120" />
        <el-table-column align="left" label="过期时间" prop="timeLimit" width="120" />
        <el-table-column align="left" label="运营商" prop="operator" width="120" />
        <el-table-column align="left" label="省市" prop="location" width="120" />
        <el-table-column align="left" label="付款码" prop="imgBaseStr" width="120" >
          <template #default="{ row }">
            <el-image :src="row.imgBaseStr" fit="contain" class="thumbnail-image"/>
          </template>
        </el-table-column>

        
        <el-table-column align="left" label="操作">
            <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
                <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>
                查看详情
            </el-button>
            <el-button type="primary" link icon="edit" class="table-button" @click="updateVboxChannelPayCodeFunc(scope.row)">变更</el-button>
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
      <el-scrollbar height="600px" >
          <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="120px">
            <el-form-item label="通道:" prop="cid">
                <el-cascader
                  v-model="formData.cid"
                  :options="channelCodeOptions"
                  :props="channelCodeProps"
                  @change="handleChange"
                  style="width: 100%"
                />
              </el-form-item>
            <el-form-item label="通道账户:"  prop="acAccount" >
              <el-select
                v-model="formData.acAccount"
                placeholder="请选择通道账号"
                filterable
                style="width: 100%"
              >
                <el-option
                  v-for="item in accList"
                  :key="item.acAccount"
                  :label="item.acAccount"
                  :value="item.acAccount"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="过期时间:"  prop="timeLimit" >
              <el-date-picker
              v-model="formData.timeLimit"
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
              style="width: 100%"
              align="center"
              >
               </el-date-picker>
            </el-form-item>
            <el-form-item label="运营商:"  prop="operator" >
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
            <el-form-item label="省市:"  prop="selectedCity" >
              <el-cascader 
                style="width:100%"
                size="small"
                :options="optionsRegion"
                v-model="selectedCity"
                @change="chge"
                placeholder="省 / 市 / 区"
                > 
              </el-cascader>
            </el-form-item>
            <el-form-item label="图片上传:"  prop="img_base_str" >
               <el-upload
							class="avatar-uploader"
							action=""
							:on-change="getFiles" 
							:on-remove="handlePicRemoves" 
							:on-preview="handlePicPreviews" 
							v-model="img_base_str"
							:limit="1" 
							list-type="picture-card" 
							:file-list="filelists" 
							:auto-upload="false" 
							accept="image/png, image/gif, image/jpg, image/jpeg"
						>
						<!-- 图标 -->
							<el-icon
								style="font-size: 25px;"
								><Plus /></el-icon>
 
						</el-upload>
            <el-dialog
							v-model="dialogVisibles"
							title="预览"
							destroy-on-close
						>
							<img
								:src="dialogImageUrs"
								style="
									display: block;
									max-width: 500px;
									margin: 0 auto;
									height: 500px;
								"
							/>
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
                        {{ formData.timeLimit }}
                </el-descriptions-item>
                <el-descriptions-item label="运营商">
                        {{ formData.operator }}
                </el-descriptions-item>
                <el-descriptions-item label="省市">
                        {{ formData.location }}
                </el-descriptions-item>
                <el-descriptions-item label="图片base64编码" >
                    <el-image :src="formData.imgBaseStr" fit="contain" class="thumbnail-image"/>
                </el-descriptions-item>
        </el-descriptions>
      </el-scrollbar>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createVboxChannelPayCode,
  deleteVboxChannelPayCode,
  deleteVboxChannelPayCodeByIds,
  updateVboxChannelPayCode,
  findVboxChannelPayCode,
  getVboxChannelPayCodeList
} from '@/api/vboxChannelPayCode'
import {
  getChannelProductSelf
} from '@/api/channelProduct'
import {
  getChannelAccountList
} from '@/api/channelAccount'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import {regionData} from 'element-china-area-data';

defineOptions({
    name: 'VboxChannelPayCode'
})

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        cid: '',
        acAccount: '',
        timeLimit: '',
        operator: '',
        location: '',
        imgBaseStr: '',
        uid: 0,
        })


// 验证规则
const rule = reactive({
})

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


// -----------上传图片--------------
const img_base_str = ref('')
const dialogImageUrs = ref("");
const filelists = ref([]);

const dialogVisibles = ref(false);
const lists = ref("");

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
      img_base_str.value = data.result;
      lists.value= data.result;

    } catch (error) {
      console.error(error);
    }
  } else {
    ElMessage({
                  type: 'error',
                  message: '上传图片大小不能超过 2MB!'
                })
  }
  console.log("fileList-111", fileList);
  console.log("file-111", file);
  formData.value.imgBaseStr=img_base_str.value
}

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
  console.log(optionsRegion.value);
};

// --------- 获取通信商 -----------
const operatorSelect = ref('')
const operators = [
  {
    value: '电信',
    label: '电信',
  },
  {
    value: '移动',
    label: '移动',
  },
  {
    value: '联通',
    label: '联通',
  }
]

// ------- 获取通道账号 -------
const accList = ref([])
const sysUserAccs = ref('')
const endTime = ref('')

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
  const table = await getVboxChannelPayCodeList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  const vcpTable = await getChannelProductSelf({ page: 1, pageSize: 999, ...searchInfo.value })
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
            deleteVboxChannelPayCodeFunc(row)
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
      const res = await deleteVboxChannelPayCodeByIds({ ids })
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
const updateVboxChannelPayCodeFunc = async(row) => {
    const res = await findVboxChannelPayCode({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data.revboxChannelPayCode
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteVboxChannelPayCodeFunc = async (row) => {
    const res = await deleteVboxChannelPayCode({ ID: row.ID })
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
  const res = await findVboxChannelPayCode({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.revboxChannelPayCode
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
          cid: '',
          acAccount: '',
          timeLimit: '',
          operator: '',
          location: '',
          imgBaseStr: '',
          uid: 0,
          }
}


// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
    // getALlChannelAccount()
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        cid: '',
        acAccount: '',
        timeLimit: '',
        operator: '',
        location: '',
        imgBaseStr: '',
        uid: 0,
        }
}
// 弹窗确定
const enterDialog = async () => {
     

      console.log('formData' + JSON.stringify(formData.value))
     elFormRef.value?.validate( async (valid) => {
   
          // console.log('formData' + JSON.stringify(formData.value))
             if (!valid) return
              let res
              switch (type.value) {
                case 'create':
                  res = await createVboxChannelPayCode(formData.value)
                  break
                case 'update':
                  res = await updateVboxChannelPayCode(formData.value)
                  break
                default:
                  res = await createVboxChannelPayCode(formData.value)
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

</script>

<style>

</style>

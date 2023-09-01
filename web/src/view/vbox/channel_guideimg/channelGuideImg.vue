<template>
  <div>

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
        <el-table-column align="left" label="图片" prop="img_base_str" width="200"  >
          <template #default="{ row }">
            <!-- {{ formatValue(row.img_base_str) }} -->
            <!-- <div class="cell-content">{{ row.img_base_str }} -->
              <el-image :src="row.img_base_str" fit="contain" class="thumbnail-image"/>
              
               <!-- <CustomPic pic-type="file" :pic-src="row.img_base_str" preview/> -->
            <!-- </div> -->
          </template>
        </el-table-column>
        
        <el-table-column align="left" label="日期" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="通道id" prop="c_channel_id" width="80" />
        <el-table-column align="left" label="文件base 64编码" prop="img_base_str" width="200" >
          <template #default="{ row }">
             <!-- {{ formatValue(row.img_base_str) }}  -->
            <div class="cell-content">{{ row.img_base_str }}
            </div>
          </template>
        </el-table-column> 
        <el-table-column align="left" label="顺序" prop="img_num" width="80" />
        <el-table-column align="left" label="文件名" prop="file_name" width="120" />
        <!-- <el-table-column align="left" label="图片地址" prop="url" width="120" /> -->
        <el-table-column align="left" label="文件标签" prop="tag" width="120" />
        <!-- <el-table-column align="left" label="编号" prop="key" width="120" /> -->
        <el-table-column align="left" label="操作">
            <template #default="scope">
            <el-button type="primary" link icon="edit" class="table-button" @click="updateChannel_guideimgFunc(scope.row)">变更</el-button>
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

    
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :title="type==='create'?'添加':'修改'" destroy-on-close>
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="130px">
        <el-form-item label="照片" prop="img_base_str">
							 <el-upload
							class="avatar-uploader"
							action=""
							:on-change="getFilesj" 
							:on-remove="handlePicRemovesj" 
							:on-preview="handlePicPreviewsj" 
							v-model="formData.img_base_str"
							:limit="1" 
							list-type="picture-card" 
							:file-list="filelistsj" 
							:auto-upload="false" 
							accept="image/png, image/gif, image/jpg, image/jpeg"
						>
						<!-- 图标 -->
							<el-icon
								style="font-size: 25px;"
								><Plus /></el-icon>
 
						</el-upload>
						<el-dialog
							v-model="dialogVisiblesj"
							title="预览"
							destroy-on-close
						>
							<img
								:src="dialogImageUrsj"
								style="
									display: block;
									max-width: 500px;
									margin: 0 auto;
									height: 500px;
								"
							/>
						</el-dialog>
						</el-form-item>
        <el-form-item label="通道id:"  prop="c_channel_id" >
          <el-input v-model="formData.c_channel_id" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="文件base64编码:"  prop="img_base_str" >
          <el-input type="textarea"  v-model="formData.img_base_str" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="图片顺序:"  prop="img_num" >
          <el-input v-model.number="formData.img_num" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="文件名:"  prop="file_name" >
          <el-input v-model="formData.file_name" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <!-- <el-form-item label="图片地址:"  prop="url" >
          <el-input v-model="formData.url" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="文件标签:"  prop="tag" >
          <el-input v-model="formData.tag" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="编号:"  prop="key" >
          <el-input v-model="formData.key" :clearable="true"  placeholder="请输入" />
        </el-form-item> -->
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

<script>
export default {
  name: 'Channel_guideimg'
}
</script>

<script setup>
import {
  createChannel_guideimg,
  deleteChannel_guideimg,
  deleteChannel_guideimgByIds,
  updateChannel_guideimg,
  findChannel_guideimg,
  getChannel_guideimgList
} from '@/api/channelGuideImg'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive,watch } from 'vue'
import SparkMD5 from 'spark-md5'
import {
  findFile,
  breakpointContinueFinish,
  removeChunk,
  breakpointContinue
} from '@/api/breakpoint'

const dialogVisiblesj = ref(false);
const dialogImageUrsj = ref("");
const filelistsj = ref([]);
const listst = ref("");


// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        c_channel_id: '',
        img_base_str: '',
        img_num: 0,
        file_name: '',
        url: '',
        tag: '',
        key: '',
        })

// 验证规则
const rule = reactive({
})



const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})



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
  const table = await getChannel_guideimgList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
const setOptions = async () =>{
}

// 获取需要的字典 可能为空 按需保留
setOptions()


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
            deleteChannel_guideimgFunc(row)
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
      const res = await deleteChannel_guideimgByIds({ ids })
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
const updateChannel_guideimgFunc = async(row) => {
    const res = await findChannel_guideimg({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data.rechGuideImg
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteChannel_guideimgFunc = async (row) => {
    const res = await deleteChannel_guideimg({ ID: row.ID })
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

// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        c_channel_id: '',
        img_base_str: '',
        img_num: 0,
        file_name: '',
        url: '',
        tag: '',
        key: '',
        }
}
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              let res
              switch (type.value) {
                case 'create':
                  res = await createChannel_guideimg(formData.value)
                  break
                case 'update':
                  res = await updateChannel_guideimg(formData.value)
                  break
                default:
                  res = await createChannel_guideimg(formData.value)
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
const getFilesj = async (file, fileList) => {
  const isLt2M = file.size / 1024 / 1024 < 2;
  if (isLt2M) {
    try {
      const data = await uploadImgToBase64(file.raw);
      formData.value.img_base_str = data.result;
      listst.value= data.result;
      // 接口请求
      /*
      const response = await fetch('接口地址');
      const test = await response.json();
      if (test.ret == 200) {
        this.$message.success("识别成功");
        window.localStorage.setItem('userImg', this.form.idCardBack);
        this.form.idCardValidity = test.data.end_date;
        console.log("yy.data.address--1", test.data.address);
      } else {
        this.$message.error("身份证识别错误");
      }
      console.log("test", test);
      */
      // this.form.idCardBack= file.url;
      // console.log("file-curl-1",this.form.idCardBack);

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
}

const handlePicRemovesj = (file, fileList) => {
  let hideUploadEdit = fileList.length
  if (hideUploadEdit >= 1){
    formData.value.img_base_str = "";
  } 
  
};

const handlePicPreviewsj = (file) => {
  console.log('file=' + file.url);
  dialogImageUrsj.value = file.url;
  dialogVisiblesj.value = true;
}

const formatValue = (value) => {
  const maxLength = 20; // 设置需要截取的最大长度
  if (value.length <= maxLength) {
    return value;
  } else {
    return value.slice(0, maxLength) + '...';
  }
};

</script>


<style lang='scss' scoped>


.avatar-uploader {
  width: 150px;
  height: 150px !important;
  overflow: hidden;
}
.cell-content {
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  user-select: all;
}

// class="thumbnail-image"
.thumbnail-image {
  max-width: 100px; /* 调整图片最大宽度 */
  min-height: 200px;
}
</style>


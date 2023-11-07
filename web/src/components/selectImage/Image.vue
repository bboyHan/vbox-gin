<template>
  <div
    v-if="!multiple"
    class="update-image"
    :style="{
      'background-image': `url(${getUrl(modelValue)})`,
      'position': 'relative',
    }"
  >
    <el-icon
      v-if="isVideoExt(modelValue || '')"
      :size="32"
      class="video video-icon"
      style=""
    ><VideoPlay /></el-icon>
    <video
      v-if="isVideoExt(modelValue || '')"
      class="avatar video-avatar video"
      muted
      preload="metadata"
      style=""
    >
      <source :src="getUrl(modelValue) + '#t=1'">
    </video>
  </div>
  <div
    v-else
    class="multiple-img"
  >
    <div
      v-for="(item, index) in multipleValue"
      :key="index"
      class="update-image"
      :style="{
        'background-image': `url(${getUrl(item)})`,
        'position': 'relative',
      }"
    >
    </div>
  </div>
</template>

<script setup>

import { getUrl, isVideoExt } from '@/utils/image'
import { onMounted, ref } from 'vue'
import { getFileList, editFileName } from '@/api/fileUploadAndDownload'
import { ElMessage, ElMessageBox } from 'element-plus'

const imageUrl = ref('')
const imageCommon = ref('')

const search = ref({})
const page = ref(1)
const total = ref(0)
const pageSize = ref(20)

const props = defineProps({
  modelValue: {
    type: [String, Array],
    default: ''
  },
  multiple: {
    type: Boolean,
    default: false
  },
  fileType: {
    type: String,
    default: ''
  }
})

const multipleValue = ref([])

onMounted(() => {
  if (props.multiple) {
    multipleValue.value = props.modelValue
  }
})

const emits = defineEmits(['update:modelValue'])

const deleteImg = (index) => {
  multipleValue.value.splice(index, 1)
  emits('update:modelValue', multipleValue.value)
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getImageList()
}

const handleCurrentChange = (val) => {
  page.value = val
  getImageList()
}
const editFileNameFunc = async(row) => {
  ElMessageBox.prompt('请输入文件名或者备注', '编辑', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputPattern: /\S/,
    inputErrorMessage: '不能为空',
    inputValue: row.name
  }).then(async({ value }) => {
    row.name = value
    // console.log(row)
    const res = await editFileName(row)
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '编辑成功!',
      })
      getImageList()
    }
  }).catch(() => {
    ElMessage({
      type: 'info',
      message: '取消修改'
    })
  })
}

const drawer = ref(false)
const picList = ref([])

const imageTypeList = ['png', 'jpg', 'jpge', 'gif', 'bmp', 'webp']
const videoTyteList = ['mp4', 'avi', 'rmvb', 'rm', 'asf', 'divx', 'mpg', 'mpeg', 'mpe', 'wmv', 'mkv', 'vob']

const listObj = {
  image: imageTypeList,
  video: videoTyteList
}

const chooseImg = (url) => {
  console.log(url)
  if (props.fileType) {
    const typeSuccess = listObj[props.fileType].some(item => {
      if (url.includes(item)) {
        return true
      }
    })
    if (!typeSuccess) {
      ElMessage({
        type: 'error',
        message: '当前类型不支持使用'
      })
      return
    }
  }
  if (props.multiple) {
    multipleValue.value.push(url)
    emits('update:modelValue', multipleValue.value)
  } else {
    emits('update:modelValue', url)
  }
  drawer.value = false
}
const openChooseImg = async() => {
  if (props.modelValue && !props.multiple) {
    emits('update:modelValue', '')
    return
  }
  await getImageList()
  drawer.value = true
}

const getImageList = async() => {
  const res = await getFileList({ page: page.value, pageSize: pageSize.value, ...search.value })
  if (res.code === 0) {
    picList.value = res.data.list
    total.value = res.data.total
    page.value = res.data.page
    pageSize.value = res.data.pageSize
  }
}

</script>

<style scoped lang="scss">

.multiple-img{
  display: flex;
  gap:8px;
  width: 100%;
  flex-wrap: wrap;
}

.add-image{
  width: 120px;
  height: 120px;
  line-height: 120px;
  display: flex;
  justify-content: center;
  border-radius: 20px;
  border: 1px dashed #ccc;
  background-size: cover;
  cursor: pointer;
}

.update-image {
  cursor: pointer;
  width: 120px;
  height: 120px;
  line-height: 120px;
  display: flex;
  justify-content: center;
  border-radius: 20px;
  border: 1px dashed #ccc;
   background-repeat: no-repeat;
   background-size: cover;
   position: relative;
  &:hover {
    color: #fff;
    background: linear-gradient(
            to bottom,
            rgba(255, 255, 255, 0.15) 0%,
            rgba(0, 0, 0, 0.15) 100%
    ),
    radial-gradient(
            at top center,
            rgba(255, 255, 255, 0.4) 0%,
            rgba(0, 0, 0, 0.4) 120%
    )
    #989898;
    background-blend-mode: multiply, multiply;
    background-size: cover;
    .update {
      color: #fff;
    }
    .video {
      opacity: 0.2;
    }
  }

  .video-icon {
    position: absolute; left: calc(50% - 16px); top: calc(50% - 16px);
  }
  video {
    object-fit: cover; max-width:100%; border-radius: 20px;
  }
  .update {
    height: 120px;
    width: 120px;
    text-align: center;
    color: transparent;
    position: absolute;
  }
}

.upload-btn-media-library {
  margin-left: 20px;
}

.media {
  display: flex;
  flex-wrap: wrap;

  .media-box {
    width: 120px;
    margin-left: 20px;

    .img-title {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      line-height: 36px;
      text-align: center;
      cursor: pointer;
    }

    .header-img-box-list {
      width: 120px;
      height: 120px;
      border: 1px dashed #ccc;
      border-radius: 8px;
      text-align: center;
      line-height: 120px;
      cursor: pointer;
      overflow: hidden;
      .el-image__inner {
        max-width: 120px;
        max-height: 120px;
        vertical-align: middle;
        width: unset;
        height: unset;
      }

      .el-image {
        position: relative;
      }
      .video-icon {
        position: absolute; left: calc(50% - 16px); top: calc(50% - 16px);
      }
      video {
        object-fit: cover; max-width:100%; min-height: 100%; border-radius: 8px;
      }
    }
  }
}
</style>

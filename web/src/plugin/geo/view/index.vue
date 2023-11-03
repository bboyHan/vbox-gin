<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="handleAdd({...geo,level:-1})">新增省会</el-button>
      </div>
      <el-table
        ref="mainTable"
        :data="treeData"
        style="width: 100%"
        row-key="code"
        border
        lazy
        :load="load"
        :indent="8"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
      >
        <el-table-column prop="id" label="编号" width="140" />
        <el-table-column prop="name" label="城市名称" min-width="160" />
        <el-table-column prop="level" label="城市等级" width="100">
          <template #default="scope">
            <el-tag>{{ levelMap[scope.row.level] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="code" label="城市编码" width="120" />
        <el-table-column prop="geocode" label="地理编码" width="100" />
        <el-table-column prop="latitude" label="经度" width="100" />
        <el-table-column prop="longitude" label="纬度" width="100" />
        <el-table-column prop="sort" label="排序" width="100" />
        <el-table-column width="220" label="操作">
          <template #default="scope">
            <el-button type="primary" link @click="handleEdit(scope.row)">编辑</el-button>
            <el-button v-if="scope.row.level!=2" type="primary" link @click="handleAdd(scope.row)">新增子地区</el-button>
            <el-button type="primary" link @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="geoDialog" title="城市管理">
      <el-form v-model="geo" label-width="120px">
        <el-form-item label="城市名称">
          <el-input v-model="geo.name" />
        </el-form-item>
        <el-form-item label="城市编码">
          <el-input v-model="geo.code" :disabled="geo.id !== 0" />
        </el-form-item>
        <el-form-item label="地理编码">
          <el-input v-model="geo.geocode" />
        </el-form-item>
        <el-form-item label="经度">
          <el-input v-model="geo.latitude" />
        </el-form-item>
        <el-form-item label="纬度">
          <el-input v-model="geo.longitude" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input v-model="geo.sort" />

        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="clear">取消</el-button>
        <el-button type="primary" @click="enter">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { getGeos, getGeo, editGeo, createGeo, deleteGeo } from '@/plugin/geo/api/geo.js'
import { ElMessageBox, ElMessage } from 'element-plus'

const loadMap = new Map()
const mainTable = ref()// table的ref

const treeData = ref([

])

const baseGeo = {
  id: 0,
  name: '',
  level: 0,
  code: '0',
  geocode: '',
  latitude: '',
  longitude: '',
  sort: 0
}

const levelMap = {
  0: '省',
  1: '市',
  2: '区',
  3: '街道',
  4: '社区'
}

const tempObj = {}

const geoDialog = ref(false)

const handleAdd = async(row) => {
  tempObj.row = row
  geo.value = {
    ...baseGeo, level: row.level, parentCode: row.code
  }
  geoDialog.value = true
}

const enter = async() => {
  if (geo.value.id) {
    const res = await editGeo(geo.value)
    if (res.code === 0) {
      reload(tempObj.row.parentCode)
      ElMessage.success('编辑成功')
      clear()
    }
  } else {
    const res = await createGeo(geo.value)
    if (res.code === 0) {
      reload(tempObj.row.code)
      ElMessage.success('添加成功')
      clear()
    }
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确认删除该地区吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
    const res = await deleteGeo(row)
    if (res.code === 0) {
      reload(row.parentCode)
      ElMessage.success('删除成功')
    }
  })
}

const clear = () => {
  geoDialog.value = false
  tempObj.row = {}
}

const geo = ref({ ...baseGeo })

const handleEdit = async(row) => {
  tempObj.row = row
  const res = await getGeo({ id: row.id, level: row.level })
  if (res.code === 0) {
    geo.value = res.data
    geoDialog.value = true
  }
}

const getTreeData = async(level, code) => {
  const res = await getGeos({ level, code })
  if (level !== 1) {
    res.data.forEach(element => {
      element.hasChildren = true
    })
  }
  return res.data
}

const init = async() => {
  treeData.value = await getTreeData('', '0')
}

init()

const reload = (code) => {
  code = code || 0
  if (loadMap.get(code)) {
    const { row, treeNode, resolve } = loadMap.get(code)
    mainTable.value.store.states.lazyTreeNodeMap[code] = []// 清空节点的数据
    load(row, treeNode, resolve)
  }
}

const load = async(row, treeNode, resolve) => {
  loadMap.set(row.code, { row, treeNode, resolve })
  resolve(await getTreeData(row.level, row.code) || [])
}

</script>

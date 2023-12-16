<template>
    <div>
      <div class="gva-search-box">
        <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
          <el-form-item label="地区" >
            <!-- <span></span> -->
            <el-cascader
                :change-on-select="true"
                style="width:100%"
                :options="optionsRegion"
                v-model="selectedCity"
                @change="chge"
                placeholder="选择地区"
                filterable
                :props="{checkStrictly: true}"
            >
            </el-cascader>
          </el-form-item>
          <el-form-item label="状态" >
            <el-select v-model="searchInfo.codeStatus" placeholder="选择状态" @change="changeCodeStatus" clearable>
              <el-option label="已使用" value="1"/>
              <el-option label="待使用" value="2"/>
              <el-option label="已失效" value="3"/>
            </el-select>
          </el-form-item>
          <el-form-item label="运营商" >
            <el-select v-model="searchInfo.operator" placeholder="选择ISP" @change="changeCodeOprater" clearable>
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
      <div class="container">
        <div id="myecharts" ref="map" ></div>
        <div class="table-container">
          <el-table
              :data="tableData"
              :default-sort="{ prop: 'codeNums', order: 'descending' }"
          >
            <el-table-column align="left" label="排名" prop="order" width="100" />
            <el-table-column align="left" label="下层区域" prop="location" width="120" />
            <el-table-column align="left" label="产码数" prop="codeNums" width="100" />
            <el-table-column align="left" label="占比" prop="ratio" width="120" :formatter="formatPercentage"/>
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
      </div>
    </div>
</template>
  
<script setup>
import {ref, onMounted, reactive} from 'vue';
import * as echarts from 'echarts';
import { mapData } from '@/utils/china';
import { codeToText, regionData } from 'element-china-area-data';
import {
  getBdaChIndexDList
} from '@/api/bdaChIndexD'

import {
  getChannelPayCodeStatisByLocation
} from '@/api/channelPayCode'


const location = ref('');
// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(30)
const tableData = ref([])
const searchInfo = ref({})

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

const elSearchFormRef = ref()

// // 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
    getTableData()
  })
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
  const table = await getChannelPayCodeStatisByLocation({ page: page.value, pageSize: pageSize.value, 
    location: location.value,operator:searchInfo.value.operator, codeStatus:searchInfo.value.codeStatus})
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}



getTableData()


const formatPercentage = (row, column, cellValue) => {
  return (cellValue * 100).toFixed(2) + '%';
};

// ---地图悬浮显示
const provicesData = ref([

])
const getProviceTableData = async() => {
  const table = await getChannelPayCodeStatisByLocation({ page: page.value, pageSize: pageSize.value, location: location.value,
    operator:searchInfo.value.operator, codeStatus:searchInfo.value.codeStatus })
  if (table.code === 0) {
    provicesData.value = table.data.list
    // console.log(JSON.stringify(provicesData.value ))
  }
}


getProviceTableData()
// ==== map ===

function selectRegionByCodeNums(nums) {
  if (nums === 0) {
    return '#ff0000'; // 红色
  } else if (nums < 5) {
    return '#ffa500'; // 橙色
  } else {
    return '#0000ff'; // 蓝色
  }
}
  
const map = ref('山西省');
const myChart = ref(null)
onMounted(() => {
  myChart.value = echarts.init(map.value);
  echarts.registerMap('china', mapData);

  let option = {
    tooltip: {
      // 鼠标移到图里面不弹出浮动提示框
      show: true,
      position: 'top',
      backgroundColor: 'rgba(0, 119, 239, 0.8)', //提示框背景色
      borderColor: '#057DD9', //边框颜色
      borderWidth: 2,
      textStyle: {
        color: '#ffffff', //文字颜色
      },
      // 自定义提示框自动调用函数
      formatter: function (params) {
        // console.log('params', params)
        // console.log('params json', JSON.stringify(params))
        // console.log('params name', params.name)
        let provinceName = params.name
        let adcode = '';
        let codenums = '';
        if (provinceName) {
          // 查找当前省份对应的 adcode
          let feature = mapData.features.find(item => item.properties.name === provinceName);
          adcode = feature?.properties.adcode;
          // console.log('params adcode', adcode)
          // console.log('params provicesData', JSON.stringify(provicesData.value))
          const filteredData = provicesData.value.filter(item => item.location === provinceName);
          codenums = filteredData.length > 0 ? filteredData.map(item => item.codeNums) : 0;
          // console.log('params codenums', codenums)
        }
        // 将adcode信息添加到弹框中进行展示
        return `省份: ${provinceName} <br/> adcode: ${adcode} <br/> 产码数: ${codenums}`;
      },
    },
    geo: {
      type: 'map',
      map: 'china',
      label: {
        show: true,
        fontSize: 8
      },
      // 移入状态高亮情况下的样式
      emphasis: {
        label: {
          show: true, //移入的时候显示文本
          color: "white",
          textBorderColor: selectedBorderColor,
          textBorderWidth: 1,
        },
      },// 普通状态下的地图省份样式
      itemStyle: {
        normal: {
          // areaColor: areaColor,
          // 径向渐变，前三个参数分别是圆心 x, y 和半径，取值同线性渐变
          areaColor: {
            type: "radial",
            x: 0.5,
            y: 0.5,
            r: 1.5,
            colorStops: [
              {
                offset: 0,
                color: areaColor[0], // 0% 处的颜色
              },
              {
                offset: 1,
                color: areaColor[1], // 100% 处的颜色
              },
            ],
            global: false, // 缺省为 false
          },
          borderColor,
          borderWidth: 1,
        },
        emphasis: {
          // areaColor: emphasisAreaColor, //鼠标指上时的颜色
          // 径向渐变，前三个参数分别是圆心 x, y 和半径，取值同线性渐变
          areaColor: {
            type: "radial",
            x: 0.5,
            y: 0.5,
            r: 1.5,
            colorStops: [
              {
                offset: 0,
                color: selectedBorderColor, // 0% 处的颜色
              },
              {
                offset: 1,
                color: selectedBorderColor, // 100% 处的颜色
              },
            ],
            global: false, // 缺省为 false
          },
          borderColor: selectedBorderColor,
          borderWidth: 2,
        },
      },
      // 被选中状态下的地图省份样式
      select: {
        label: {
          show: false, //1.国内业务布局：地图上省份的文字 最好不显示 和设计稿一致
          color: "white",
          textBorderColor: selectedBorderColor,
          textBorderWidth: 1,
        },
        itemStyle: {
          // areaColor: areaColor,
          // 径向渐变，前三个参数分别是圆心 x, y 和半径，取值同线性渐变
          areaColor: {
            type: "radial",
            x: 0.5,
            y: 0.5,
            r: 1.5,
            colorStops: [
              {
                offset: 0,
                color: selectedAreaColor[0], // 0% 处的颜色
              },
              {
                offset: 1,
                color: selectedAreaColor[1], // 100% 处的颜色
              },
            ],
            global: false, // 缺省为 false
          },
          borderColor: selectedBorderColor,
          borderWidth: 2,
        },
      },

    },
    series: [
      {
        name: '山西',
        type: 'map',
        coordinateSystem: 'geo',
        geoIndex: 0,
        animationDuration: 1200,
        rippleEffect: {
          brushType: 'stroke',
        },
        label: {
          // 显示地图省名称
          normal: {
            color: '#389dff',
            formatter: '{b}',
            position: [-12, -1],
            show: true,
          },
          emphasis: {
            show: true,
            color: '#fff',
          },
        },

        itemStyle: {
          normal: {
            areaColor: '#0d0059',
            borderColor: '#389dff',
            borderWidth: 0.5,
          },
          emphasis: {
            areaColor: '#59D5F5',
            shadowOffsetX: 8,
            shadowOffsetY: 8,
            shadowBlur: 5,
            borderWidth: 0,
            shadowColor: '#0074BC',
          },
        },
        data: [
          { name: '广东', value: [113.88308, 22.55329], weidu: 116.397128, jingdu: 112.98626, orderNum: 1, clientNum: 2 },
        ]
      },
    ],
  };

  myChart.value.setOption(option);

  myChart.value.on('mouseover', function (params) {
      // 获取当前省份名称
    let provinceName = params.name;
    // 查找当前省份对应的 adcode
    let feature = mapData.features.find(item => item.properties.name === provinceName);
    let adcode = feature?.properties.adcode;
      // 这里可以根据需要自定义弹出提示框、弹出层等来展示统计数据
    console.log('Province:', provinceName, 'Adcode:', adcode);
  });
});

// ------------获取省市 -------
const selectedCity = ref([]);

const optionsRegion = regionData;
const chge = () => {

  const lastElement = selectedCity.value[selectedCity.value.length - 1]
  location.value= lastElement
//   console.log(selectedCity);
  console.log('location:', location.value);
//   showProvince('山西省')
    getTableData()
};

// ----- 获取状态

const changeCodeStatus = () => {
  console.log('changeCodeStatus');
  getTableData()
};

const changeCodeOprater = () => {
  console.log('changeCodeOprater');
  getTableData()
};

// ---------   -----

const showProvince = name => {
    myChart.value &&
    myChart.value.dispatchAction({
      type: 'geoselected',
      name,
    });
};

// series ■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
const dotColor = "white"; //地图上的标记点颜色
// const dotColor = "#087fd8", //地图上的标记点颜色
const lineColor = ["#067DD6", "#067dd6"]; //地图省边框线、迁徙轨迹线条颜色（线性渐变：上，下）
const selectedBorderColor = "#067DD6"; //被选中地图省份边框颜色
const borderColor = "#A4D9FF"; //地图省份边框颜色
const selectedAreaColor = ["#1D99F5", "#1D99F5"]; //被选中地图省份填充色（径向渐变：内，外）
const areaColor = ["#d6edff", "#d6edff", '#ff0000', '#ffa500', '#0000ff']; //地图省份填充色（径向渐变：内，外）
const emphasisAreaColor = ["#1D99F5 ", "#1D99F5 "]; //移入地图省份时的填充色（径向渐变：内，外）
</script>
  
<style>
#myecharts {
  width: 1400px;
  height: 1100px;
  margin-left: -100px;
}

.container {
  display: flex;
}

.table-container {
  width: 30%; /* 设置宽度占比 */
}
</style>



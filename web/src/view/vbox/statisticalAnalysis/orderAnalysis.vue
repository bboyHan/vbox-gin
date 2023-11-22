
<template>
    <div>
        <div class="gva-btn-list">
            <span>省市：</span>
            <el-cascader
                :change-on-select="true"
                style="width:15%"
                :options="optionsRegion"
                v-model="selectedCity"
                @change="chge"
                placeholder="省 / 市 / 区 "
                filterable
                :props="{checkStrictly: true}"
            >
            </el-cascader>
        </div> 
        <div class="container">
            <div id="myecharts" ref="map" ></div>


            <div class="table-container">
                
                <el-table
                :data="tableData"
                :default-sort="{ prop: 'codeNums', order: 'descending' }"
                >
                <el-table-column align="left" label="排名" prop="order" width="100" />
                <el-table-column align="left" label="区域" prop="location" width="120" />
                <el-table-column align="left" label="产码数" prop="codeNums" width="100" />
                <el-table-column align="left" label="占比" prop="ratio" width="120" />
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
  import { ref, onMounted } from 'vue';
  import * as echarts from 'echarts';
  import { mapData } from '@/utils/china';
  import { codeToText, regionData } from 'element-china-area-data';
  import {
  getBdaChIndexDList
} from '@/api/bdaChIndexD'

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])


// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// // 搜索
// const onSubmit = () => {
//   elSearchFormRef.value?.validate(async(valid) => {
//     if (!valid) return
//     page.value = 1
//     pageSize.value = 10
//     getTableData()
//   })
// }

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
  const table = await getBdaChIndexDList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ==== map ===

  
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
                console.log('params', params)
                console.log('params json', JSON.stringify(params))
                console.log('params name', params.name)
                let provinceName = params.name
                let adcode = ''; 
                if (provinceName) {
                    // 查找当前省份对应的 adcode
                    let feature = mapData.features.find(item => item.properties.name === provinceName);
                    adcode = feature?.properties.adcode;
                    console.log('params adcode', adcode)
                }
                // 将adcode信息添加到弹框中进行展示
                return `${provinceName} <br/> adcode: ${adcode}`; 
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

const location = ref('');
const optionsRegion = regionData;
const chge = () => {

  const lastElement = selectedCity.value[selectedCity.value.length - 1]
  location.value= lastElement
  console.log(selectedCity);
  console.log('location:', location.value);
  showProvince('山西省')
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
const areaColor = ["#d6edff", "#d6edff"]; //地图省份填充色（径向渐变：内，外）
const emphasisAreaColor = ["#1D99F5 ", "#1D99F5 "]; //移入地图省份时的填充色（径向渐变：内，外）
  </script>
  
  <style>
  #myecharts {
    width: 1400px;
    height: 1100px;
    margin-left: -200px;
  }

  .container {
  display: flex;
}

.table-container {
  width: 30%; /* 设置宽度占比 */
}
  </style>



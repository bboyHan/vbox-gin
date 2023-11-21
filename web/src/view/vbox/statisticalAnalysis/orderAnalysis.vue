<template>
    <div id="myecharts" ref="map"></div>
  </template>
  
  <script>
  import * as echarts from 'echarts';
  import { mapData } from '@/utils/china';
  
  export default {
    name: 'Middle',
    mounted() {
      let myChart = echarts.init(this.$refs.map);
      echarts.registerMap('china', mapData);
  
      let option = {
        geo: {
          type: 'map',
          map: 'china',
          label: {
            show: true,
            fontSize: 8
          }
        }
      };
  
      myChart.setOption(option);
  
      myChart.on('mouseover', function (params) {
        // 获取当前省份名称
        let provinceName = params.name;
        // 查找当前省份对应的 adcode
        // let adcode;
        // console.log('mapData',JSON.stringify(mapData.features))
        // Object.keys(mapData.features).forEach(ad => {

        //     // console.log('mapData.features[ad]',mapData.features[ad].properties.name)
        //     if (mapData.features[ad].properties.name === provinceName) {
        //         adcode = ad;
        //     }
        // });
        // 查找当前省份对应的 feature 对象
        let feature = mapData.features.find(item => item.properties.name === provinceName);

        // console.log('feature',JSON.stringify(feature))
        // 获取当前省份对应的 adcode
        let adcode = feature?.properties.adcode;

        console.log('Province:', provinceName, 'Adcode:', adcode);
        // 这里可以根据需要自定义弹出提示框、弹出层等来展示统计数据
      });
    }
  };
  </script>
  
  <style>
  #myecharts {
    width: 1500px;
    height: 1200px;
    margin-left: -200px; /* 将整体左移 */
  }
  </style>
  
<template>
  <div class="data-center-card" :style="customStyle">
    <div class="data-center-card-head">
      <div class="data-center-card-head-title">{{ title }}</div>
      <slot name="action" />
    </div>
    <div class="data-center-card-body">
      <slot name="body" />
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  title: String,
  channelCode: String,
  customStyle: {},
})


const customStyle = props.customStyle
// console.log(customStyle)
</script>

<style lang="scss" scoped>
.data-center-card:hover {
  transform: scale(1.02); /* 当鼠标悬停在卡片上时，稍微放大卡片 */
}

.data-center-card{
  /* 如果要使用自定义颜色，请替换上面的颜色值 */ /* 设置卡片背景颜色 */
  border: 1px solid #ccc; /* 设置卡片边框 */
  border-radius: 10px; /* 设置卡片边框的圆角 */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); /* 设置卡片的阴影效果 */
  transition: transform 0.3s ease; /* 添加过渡效果 */
  //min-width: 250px;
  flex: auto;
  padding: 8px;
  margin-bottom: 12px;
  &-head{
    display: flex;
    align-items: center;
    justify-content: space-between;

    &-title{
      font-size: 14px;
      font-weight: 500;
      padding-left: 10px;
      position: relative;
      &::before{
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        width: 4px;
        height: 100%;
        border-radius: 4px;
        background-color: #1890ff;
      }
    }
  }
  &-body{
    padding-top: 10px;
  }
}
</style>

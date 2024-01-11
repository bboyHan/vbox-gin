
<template>
  <div class="mt-2">
  </div>
</template>
<!--
<script setup>
import { ref, unref } from 'vue';
import { MarkDownActionType, MarkdownViewer } from '/@/components/Markdown';

const markDownRef = ref(null);
const valueRef = ref(`
### 0. 对接信息
`);

function toggleTheme() {
  const markDown = unref(markDownRef);
  if (!markDown) return;
  const vditor = markDown.getVditor();
  vditor.setTheme('dark', 'dark', 'dracula');
}

function handleChange(v) {
  valueRef.value = v;
}

function clearValue() {
  valueRef.value = '';
}

export { valueRef, toggleTheme, markDownRef, handleChange, clearValue };
</script>
-->

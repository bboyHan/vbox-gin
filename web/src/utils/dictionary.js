import { useDictionaryStore } from '@/pinia/modules/dictionary'
//  获取字典方法 使用示例 getDict('sex').then(res)  或者 async函数下 const res = await getDict('sex')
export const getDict = async(type) => {
  const dictionaryStore = useDictionaryStore()
  await dictionaryStore.getDictionary(type)
  return dictionaryStore.dictionaryMap[type]
}

//  字典文字展示方法
export const showDictLabel = (dict, code) => {
  if (!dict) {
    return ''
  }
  const dictMap = {}
  dict.forEach(item => {
    dictMap[item.value] = item.label
  })
  return dictMap[code]
}

// 获取数据中指定key的value集合
function getCollection2Key(list, key) {
  return list.map(function(item) {
    return item[key];
  });
}
import FileSaver from "file-saver";
import XLSX from "xlsx";

export const useExcle = (name, tableName) => {
    //name表示生成excel的文件名     tableName表示表格的id
    var sel = XLSX.utils.table_to_book(document.querySelector(tableName))
    var selIn = XLSX.write(sel, { bookType: 'xlsx', bookSST: true, type: 'array' })
    try {
      FileSaver.saveAs(new Blob([selIn], { type: 'application/octet-stream' }), name)
    } catch (e) {
      if (typeof console !== 'undefined') console.log(e, selIn) 
    }
    return selIn
}
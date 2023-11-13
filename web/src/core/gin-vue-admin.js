/*
 * gin-vue-admin web框架组
 *
 * */
// 加载网站配置文件夹
import { register } from './global'

export default {
  install: (app) => {
    register(app)
    // 默认自动化文档地址:http://127.0.0.1:${import.meta.env.VITE_SERVER_PORT}/swagger/index.html
    console.log(`
       默认前端文件运行地址:http://127.0.0.1:${import.meta.env.VITE_CLI_PORT}
    `)
  }
}

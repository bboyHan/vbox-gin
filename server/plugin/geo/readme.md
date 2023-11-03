

## GVA 城市管理插件
#### 开发者：Mr.奇淼

### 示例

![图片](https://qmplusimg.henrongyi.top/mini_geo1.jpg)

![图片](https://qmplusimg.henrongyi.top/mini_geo2.jpg)

![图片](https://qmplusimg.henrongyi.top/mini_geo3.jpg)

#### 1. 前往GVA主程序下的initialize/plugin.go 下注册插件
        
    PluginInit(PublicGroup, geo.CreateGeoPlug())
        // 前端已有对应调用逻辑，无需任何改动，只要填充配置皆可使用

### 2. 配置说明

#### 2-1 全局配置结构体说明
    
    无

### 3. 可直接调用的接口
    获取子区域信息列表： /geo/getGeos [get]
    入参示例 level=1&code=11
    
    获取当前区域信息： /geo/getGeo [get]
    入参示例 level=1&id=1
    
    修改区域信息： /geo/editGeo [put]
    入参示例 {   
                id: 1,
                name: '',
                level: 0,
                geocode: '',
                latitude: '',
                longitude: '',
                sort: 0
            }
    
    创建区域信息： /geo/createGeo [post]
    入参示例 {
                name: '',
                level: 0,
                code: '0',
                geocode: '',
                latitude: '',
                longitude: '',
                sort: 0
            }
    
    删除区域信息： /geo/deleteGeo [delete]
    入参示例 level=1&id=1
### 数据库
    请到 https://github.com/piexlmax/geoDB 下载数据文件 并将表导入自己目前的开发库
    执行如下sql 添加 菜单 
    INSERT INTO `gvaplug`.`sys_base_menus`(`id`, `created_at`, `updated_at`, `deleted_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `keep_alive`, `default_menu`, `title`, `icon`, `close_tab`) VALUES (31, '2022-06-28 10:58:43.347', '2022-06-28 10:58:43.347', NULL, 0, '28', 'geo', 'geo', 0, 'plugin/geo/view/index.vue', 0, 0, 0, '城市管理', 'aim', 0);
    到角色管理下分配菜单即可使用




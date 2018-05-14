# h5client

H5 测试客户端

## 环境搭建

安装以下工具：

-   nodejs

    点击官方页面下载安装：[nodejs官网](https://nodejs.org)

-   cnpm

    ```dos
    npm install cnpm -g --registry=https://registry.npm.taobao.org
    ```

## 更新依赖

```dos
cnpm install
```

或者

```dos
npm install
```

## 运行

  1. 拷贝`webpack.config.js.sample`，并重命名为`webpack.config.js`
  
  1. 运行下列命令：

  - NodeJS服务器版
    
      ```dos
      npm run web
      ```
    

  - 桌面版

      `webpack.config.js`文件中，devServer.open = false

      ```dos
      npm run web
      npm run app
      ```

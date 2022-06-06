// 1.引入 electron
const { app, BrowserWindow } = require('electron')
// 定义一个窗口
let win = null
// 2.引入自定义的菜单
require('./menu')

// 3.监听ready
app.on('ready', function () {
  win = new BrowserWindow({
    width: 800,
    height: 600,
  })
  // 打开控制台
  win.webContents.openDevTools()
  win.loadFile('./index.html')
  // 4.监听窗口关闭事件
  win.on('close', () => {
    win = null
  })
})
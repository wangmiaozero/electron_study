import { createProtocol } from 'vue-cli-plugin-electron-builder/lib'
import { app, BrowserWindow, ipcMain, Tray, Menu, screen } from 'electron'

const path = require('path')
const iconPath = path.join(__static, 'icon.png')

let mainWindow
let tray
let remindWindow

app.on('ready', async () => {
  mainWindow = new BrowserWindow({
    frame: false,
    resizable: true,
    width: 800,
    height: 600,
    icon: iconPath,
    webPreferences:{
      backgroundThrottling: false,
      nodeIntegration:true,
      contextIsolation: false
    }
  })
  if (process.env.WEBPACK_DEV_SERVER_URL) {
     // 开发环境下，访问dev server的地址
    await mainWindow.loadURL(process.env.WEBPACK_DEV_SERVER_URL + '/main.html')
  } else {
    createProtocol('app')
   // 生产环境下，使用`file://`协议加载main.html
    mainWindow.loadURL(`file://${__dirname}/main.html`)
  }
  mainWindow.removeMenu()
  setTray ()
})

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) createWindow()
})

ipcMain.on('mainWindow:close', () => {
  mainWindow.hide()
})

ipcMain.on('remindWindow:close', () => {
  remindWindow.close()
})

ipcMain.on('setTaskTimer', (event, time, task) => {
  const now = new Date()
  const date = new Date()
  date.setHours(time.slice(0,2), time.slice(3),0)
  const timeout = date.getTime() - now.getTime()
  setTimeout(() => {
    createRemindWindow(task)
  }, timeout)
})

function setTray () {
  tray = new Tray(iconPath)
  tray.setToolTip('Tasky')
  tray.on('click', () => {
    if(mainWindow.isVisible()){
      mainWindow.hide()
    }else{
      mainWindow.show()
    }
  })
  tray.on('right-click', () => {
    const menuConfig = Menu.buildFromTemplate([
      {
        label: 'Quit',
        click: () => app.quit()
      }
    ])
    tray.popUpContextMenu(menuConfig)
  })

}

function createRemindWindow (task) {
  if(remindWindow) remindWindow.close()
  //创建提醒窗口
  remindWindow = new BrowserWindow({
    height: 450,
    width: 360,
    resizable: true,
    frame: false,
    icon: iconPath,
    show: false,
    webPreferences:{
      nodeIntegration:true,
      contextIsolation: false,
    }
  })
  remindWindow.removeMenu()
  //获取屏幕尺寸
  const size = screen.getPrimaryDisplay().workAreaSize
  //获取托盘位置的y坐标（windows在右下角，Mac在右上角）
  const { y } = tray.getBounds()
  //获取窗口的宽高
  const { height, width } = remindWindow.getBounds()
  //计算窗口的y坐标
  const yPosition = process.platform === 'darwin' ? y : y - height;
  //setBounds设置窗口的位置
  remindWindow.setBounds({
    x: size.width - width,//x坐标为屏幕宽度 - 窗口宽度
    y: yPosition,
    height,
    width
  })
    //当有多个应用时，提醒窗口始终处于最上层
  remindWindow.setAlwaysOnTop(true)

  if (process.env.WEBPACK_DEV_SERVER_URL) {
    remindWindow.loadURL(process.env.WEBPACK_DEV_SERVER_URL + '/remind.html')
  } else {
    createProtocol('app')
    remindWindow.loadURL(`file://${__dirname}/remind.html`)
  }

  remindWindow.webContents.on('did-finish-load', () => {
    remindWindow.webContents.send('setTask', task)
  })

  remindWindow.show()
  // 提醒窗口会在一段时间后关闭，可以通过remindWindow.close()来关闭窗口。
  // 当窗口关闭后，我们可以设置remindWindow = null来回收分配给该渲染进程的资源。
  remindWindow.on('closed', () => { remindWindow = null })
  setTimeout( () => {
    remindWindow && remindWindow.close()
  }, 50 * 1000)
}


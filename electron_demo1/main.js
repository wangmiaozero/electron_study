const { app,BrowserView,BrowserWindow,Menu,ipcMain } = require('electron')
const nativeImage = require('electron').nativeImage;

var mainWindow = null;
var overlay = nativeImage.createFromPath('./assets/img/logo.ico');
const mainWindowURL = "https://www.wangmiaozero.cn"


function createWindow() {
  console.log("createWindow......");
  Menu.setApplicationMenu(null)
  mainWindow = new BrowserWindow({
    width: 1024,
    height: 768,
    minWidth: 1024,
    minHeight: 768,
    webPreferences: { nodeIntegration: true },
    useContentSize: true,
    maximizable: true, //支持最大化
    show: false,   //为了让初始化窗口显示无闪烁，先关闭显示，等待加载完成后再显示。
    backgroundColor:'#fff',//设置背景色
    //icon: "./assets/img/logo.ico"
    icon: "https://www.wangmiaozero.cn/image/logo/favicon.ico"
    // fullscreen: true,
    //frame: false,//隐藏标题栏
  })
  // Emitted when the window is closed.
  mainWindow.on('closed', function () {
    console.log("closed......");
    mainWindow = null;
  })
  // mainWindow.once('ready-to-show', () => {    //开启这个以后点击了图片虽然没有白屏，但是会感觉不到点了没点
  mainWindow.maximize();    //打开时最大化打开，不是全屏，保留状态栏
  // })
  mainWindow.setOverlayIcon(overlay, "全栈 | 兔子零-汪苗的个人网站-全栈修炼")
  app.setName("全栈 | 兔子零-汪苗的个人网站-全栈修炼")
  //mainWindow.setAutoHideMenuBar(true);//自动隐藏菜单
  mainWindow.loadURL(mainWindowURL);
  // app.commandLine.appendSwitch("--disable-http-cache")   禁用缓存
  // mainWindow.webContents.openDevTools({mode:'bottom'});
  mainWindow.show();
  //mainWindow.webContents.openDevTools()
}


app.on('ready', () => {
  console.log("ready......");
  createWindow()
  /*  let BrowserView = electron.BrowserView
   let view = new BrowserView();
   mainWindow.setBrowserView(view);
   view.setBounds({x:0,y:0,width:1024,height:768,useContentSize: true,})
   view.webContents.loadURL('https://www.wangmiaozero.cn'); */
})

// Quit when all windows are closed.
app.on('window-all-closed', function () {
  // On OS X it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q
  if (process.platform !== 'darwin') {
    app.quit();
  }
})

app.on('activate', () => {
  // 在macOS上，当单击dock图标并且没有其他窗口打开时，
  // 通常在应用程序中重新创建一个窗口。
  if (BrowserWindow.getAllWindows().length === 0) {
      createWindow()
  }
})
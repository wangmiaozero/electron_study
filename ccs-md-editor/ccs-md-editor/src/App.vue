<template>
  <div>
    <div class = "WinTitle">
      <div class = "Logo">
        <img src = "./assets/icon/logo.svg" style = "height: 26px;"/>
        <img src = "./assets/icon/icon.svg" style = "margin-left: 6px; height: 20px;"/>
      </div>
      <div class = "WinTitleControl"></div>
      <div class = "Control">
        <div class = "Button" @click = "BtnMinimize"><img src = "./assets/icon/minimize.svg" style = "height: 12px;"/></div>
        <div class = "Button" @click = "BtnMaximize"><img src = "./assets/icon/maximize.svg" style = "height: 12px;"/></div>
        <div class = "Button" @click = "BtnWinClose"><img src = "./assets/icon/close.svg" style = "height: 12px;"/></div>
      </div>
    </div>
    <div class = "WinContent"><router-view></router-view></div>
  </div>
</template>

<script>
const ipcRenderer = window.electron.ipcRenderer
export default {
  methods: {
    BtnMinimize: function () { ipcRenderer.send('min') },
    BtnMaximize: function () { ipcRenderer.send('max') },
    BtnWinClose: function () { ipcRenderer.send('close') }
  }
}
</script>

<style>
@import './assets/css/ccs-mde-global.css';
#app {
  font-family: 'Consolas', '微软雅黑', monospace, 'Microsoft YaHei', 'hw-iconfont';
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: var(--ccs-mde-primary-font-color);
}

.WinTitle {
  position: relative;
  width: 100vw;
  height: var(--ccs-mde-winTitle-height);
  background: var(--ccs-mde-winTitle-color);
  box-shadow: 0px 2px 4px #00000030;
  z-index: 4000;
}

.WinContent {
  width: 100vw;
  height: calc(100vh - var(--ccs-mde-winTitle-height));
}

.Logo {
  height: 100%;
  display: flex;
  align-items: center;
  padding-left: 16px;
}

.WinTitleControl {
  position: absolute;
  top: 0;
  width: calc(100vw - 160px);
  height: 100%;
  /* 设置标题栏可拖动 */
  -webkit-app-region: drag;
}

.Control {
  position: absolute;
  top: 0;
  right: 0;
  height: 100%;
  display: flex;
  align-items: center;
  padding-right: 2px;
}

.Button {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: var(--ccs-mde-winTitle-color);
  box-shadow:  3px 3px 6px #1a1a1b,
              -3px -3px 6px #303031;
  margin-right: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  user-select: none;
  transition: background .3s;
}

.Button:hover {
  background: #282828;
}

.Button:active {
  border-radius: 8px;
  background: #252526;
  box-shadow: inset 3px 3px 6px #1a1a1b,
              inset -3px -3px 6px #303031;
}

.v3-menus {
  background: #1b1b1b !important;
  padding: unset !important;
  border-radius: 6px !important;
  overflow: hidden !important;
  min-width: 240px !important;
  box-shadow: 0 1px 3px #101010 !important;
}

.v3-menus-item {
  padding: 12px 30px 12px 24px !important;
  font-size: 14px !important;
  line-height: unset !important;
}

.v3-menus-active {
  background: #232323 !important;
  color: var(--ccs-mde-primary-font-color) !important;
}

.v3-menus-icon {
  position: relative;
}

.MenusIcon {
  height: 20px;
  position: absolute;
  top: -10px;
  left: -5px;
}

.v3-menus-tip {
  font-size: 14px !important;
  color: #606266 !important;
}

.v3-menus-divided {
  border-color: #252525 !important;
}

.vditor--dark {
  --panel-background-color: var(--ccs-mde-background-color);
  --toolbar-background-color: var(--ccs-mde-toolbar-color);
  --textarea-background-color: var(--ccs-mde-background-color);
  --border-color: #1a1a1a;
}

.vditor-img__img {
  background: var(--ccs-mde-img-background-color);
}

.vditor-img__img > img {
  border-radius: 10px;
  box-shadow: 0 0 20px #080808;
}
</style>

<template>
  <div class = "FileListArea">
    <div style = "width: 100%;">
      <div class = "SearchBox">
        <input v-model = "filter.search" placeholder = "搜索笔记" @input = "GetFileList"/>
      </div>
    </div>
    <div class = "FileItemList" data-type = "fileList" @contextmenu = "MenuRightClick">
      <div v-if = "file.list.length === 0" class = "FileEmpty">
        <img src = "../assets/icon/empty.svg" style = "position: absolute; width: 70px; top: 200px; opacity: .2;">
      </div>
      <el-scrollbar height = "100%">
        <div class = "FileItemCard" v-for = "item in file.list" :key = "item.id">
          <div class = "FileItem" :class = "item.id === file.current ? 'FileItem_Activate' : ''"
            :data-id = "item.id" data-type = "fileItem" @contextmenu = "MenuRightClick" @click = "FileItemClick">
            <div style = "font-weight: 900; position: relative;">
              <img src = "../assets/icon/markdown.svg" style = "width: 22px;">
              <span style = "padding-left: 6px;">{{ item.title }}</span>
            </div>
            <div v-if = "item.summary" style = "font-size: 14px; padding-top: 10px; color: #707070; word-wrap: break-word;">
              <span>{{ item.summary }}</span>
            </div>
            <div style = "font-size: 14px; padding-top: 10px; color: #454545;">
              <span>{{ item.updateDateTime }}</span>
            </div>
          </div>
        </div>
      </el-scrollbar>
    </div>
    <!-- 右键菜单 -->
    <vue3-menus v-model:open = "menus.state" :event = "menus.value" :menus = "menus.list">
    </vue3-menus>
  </div>
</template>

<script>
import { ElScrollbar } from 'element-plus'
import { Vue3Menus } from 'vue3-menus'
export default {
  components: { ElScrollbar, Vue3Menus },
  data () {
    return {
      filter: { search: '' },
      menus: { value: {}, state: false, list: []
      },
      file: { current: this.$route.query.id || '', list: [] }
    }
  },
  methods: {
    MenuRightClick: function (e) {      
      this.menus.state = false
      this.menus.list = []

      const noteType = e.currentTarget.dataset.type
      const noteId = e.currentTarget.dataset.id || ''

      const _this = this

      let menus = [
        {
          label: "新建笔记",
          icon: '<img src = "/icon/new.svg" class = "MenusIcon">',
          tip: 'Ctrl + N',
          divided: true,
          allow: ['fileItem', 'fileList'],
          click: () => { _this.$emit('NewNote', null) }
        },
        {
          label: "保存",
          icon: '<img src = "/icon/save.svg" class = "MenusIcon">',
          allow: ['fileItem', 'id'],
          click: () => { _this.$emit('SaveNote', noteId) }
        },
        {
          label: "打开",
          icon: '<img src = "/icon/open.svg" class = "MenusIcon">',
          allow: ['fileItem'],
          click: () => { _this.$router.push('?id=' + noteId) }
        },
        {
          label: "关闭笔记",
          icon: '<img src = "/icon/close.svg" class = "MenusIcon">',
          allow: ['fileItem', 'id'],
          click: () => { _this.$router.push('/') }
        },
        {
          label: "删除",
          icon: '<img src = "/icon/delete.svg" class = "MenusIcon">',
          divided: true,
          allow: ['fileItem'],
          click: () => { _this.$emit('RemoveNote', noteId) }
        },
        {
          label: "清空搜索",
          icon: '<img src = "/icon/clean.svg" class = "MenusIcon">',
          allow: ['fileItem', 'fileList'],
          click: () => {
            _this.filter.search = ''
            _this.GetFileList()
          }
        },
        {
          label: "刷新",
          icon: '<img src = "/icon/refresh.svg" class = "MenusIcon">',
          tip: 'Ctrl + R',
          allow: ['fileItem', 'fileList'],
          click: () => { _this.GetFileList() }
        }
      ]
      menus.forEach(element => {
        if (element.allow.indexOf(noteType) > -1) {
          if (element.allow.indexOf('id') === -1) {
            this.menus.list.push(element)
          } else {
            if (noteId === this.file.current) {
              this.menus.list.push(element)
            }
          }
        }
      })
      
      this.$nextTick(() => {
        this.menus.value = e
        this.menus.state = true
      })
      e.preventDefault()
      e.stopPropagation()
    },
    FileItemClick: function (e) {
      const id = e.currentTarget.dataset.id || ''
      this.$router.push('?id=' + id)
    },
    GetFileList: function () {
      const url = 'http://localhost:41220/Note/List'
      const data = this.filter
      this.$http.post(url, data).then((res) => {
        this.file.list = res.data
      })
    }
  },
  mounted () {
    this.GetFileList()
  },
  watch: {
    $route: function (url) {
      this.GetFileList()
      this.file.current = url.query.id || ''
    }
  }
}
</script>

<style lang = "css" scoped>
  .el-scrollbar {
    --el-scrollbar-opacity: 1;
    --el-scrollbar-background-color: #00000020;
    --el-scrollbar-hover-opacity: 1;
    --el-scrollbar-hover-background-color: #00000030;
  }

  .FileListArea {
    position: relative;
    width: 100%;
    height: calc(100% - 14px - 10px);
  }

  .SearchBox {
    padding: 10px 20px;
    margin: 14px 10px;
    border-radius: 100px;
    background: #222222;
    box-shadow: inset 4px 4px 6px #1c1c1c,
                inset -4px -4px 6px #282828;
  }

  .SearchBox input {
    width: 100%;
    margin: 0;
    padding: 0;
    background: unset;
    border: unset;
    color: var(--ccs-mde-primary-font-color);
    font-size: 14px;
    text-align: center;
  }

  .SearchBox input:focus-visible {
    outline: unset;
  }

  .FileItemList {
    position: relative;
    height: calc(100% - 53px);
  }

  .FileItemCard {
    padding: 0px 10px 15px 10px;
    user-select: none;
  }

  .FileItemCard:last-child {
    padding: 0px 10px 0px 10px;
  }

  .FileItem {
    --padding: 16px;
    padding: var(--padding);
    width: calc(100% - var(--padding) * 2);
    border-radius: 12px;
    background: #262626;
    box-shadow: 5px 5px 16px #1a1a1a,
                5px 5px 16px #2f2f2f;
    overflow: hidden;
    transition: background .3s;
  }

  .FileItem:hover {
    background: #282828;
  }

  .FileItem:active {
    background: #262626;
  }

  .FileItem_Activate {
    box-shadow: inset 4px 4px 12px #1a1a1a,
                inset -4px -4px 12px #2f2f2f;
  }

  .FileEmpty {
    position: absolute;
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    z-index: 0;
  }
</style>
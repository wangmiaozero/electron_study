<template>
  <div class = "Main">
    <div class = "FileList" :style = "{ width: control.layout.fileList.width }" @mouseup = "SplitLineMouseUp">
      <FileList @NewNote = "NewNote" @SaveNote = "SaveNote" @RemoveNote = "RemoveNote" ref = "FileList"></FileList>
    </div>
    <div class = "SplitLine" @mousedown = "SplitLineMouseDown" @mouseup = "SplitLineMouseUp"></div>
    <div class = "Editor" @mouseup = "SplitLineMouseUp"
      :style = "{ width: 'calc(100% - var(--ccs-mde-splitLine-width) - ' + control.layout.fileList.width + ')' }">
      <Editor v-if = "editor.state" :id = "editor.id" ref = "Editor" @NewNote = "NewNote" @GetNoteList = "GetNoteList"></Editor>
      <div v-else style = "width: 100%; height: 100%; display: flex; justify-content: center; align-items: center;">
        <img src = "../assets/icon/icon.svg" style = "height: 50px; opacity: .05;"/>
      </div>
    </div>  
  </div>
</template>

<script>
import FileList from '../components/FileList.vue'
import Editor from '../components/Editor.vue'
// import { ElMessage } from 'element-plus'
export default {
  components: { FileList, Editor },
  data () {
    return {
      control: {
        layout: {
          fileList: { width: '300px' },
          splitLine: { state: false }
        }
      },
      editor: { id: '', state: false }
    }
  },
  methods: {
    SplitLineMouseDown: function () {
      this.control.layout.splitLine.state = true
    },
    SplitLineMouseUp: function () {
      this.control.layout.splitLine.state = false
    },
    NewNote: function () {
      const url = 'http://localhost:41220/Note/New'
      this.$http.get(url).then((res) => {
        const id = res.data
        this.$router.push('?id=' + id)
      })
    },
    SaveNote: function (id) {
      this.$refs.Editor.SaveNote(id)
    },
    RemoveNote: function (id) {
      const url = 'http://localhost:41220/Note/Remove'
      const data = { id }
      this.$http.post(url, data).then((res) => {
        const noteId = res.data || ''
        if (noteId === this.$route.query.id && noteId != '') {
          this.$router.push('/')
        }
        this.GetNoteList()
      })
    },
    GetNoteList: function () {
      this.$refs.FileList.GetFileList()
    }
  },
  mounted () {
    const id = this.$route.query.id || ''
    if (id != '') {
      this.editor = { id, state: true }
    }
    document.onmousemove = (e) => {
      // 可移动的区间为 300px ~ 50%的屏幕
      const minLocX = 300
      const maxLocX = document.body.clientWidth * .5

      let fileListWidth = parseFloat(this.control.layout.fileList.width)
      if (this.control.layout.splitLine.state) {
        fileListWidth = e.clientX
        if (fileListWidth < minLocX) {
          fileListWidth = minLocX
        } else if (fileListWidth > maxLocX){
          fileListWidth = maxLocX
        }
        this.control.layout.fileList.width = fileListWidth.toString() + 'px'
      }
    }
  },
  watch: {
    $route: function (url) {
      const id = url.query.id || ''
      this.editor = { id, state: false }
      this.$nextTick(() => {
        if (id != '') { this.editor.state = true }
      })
    }
  }
}
</script>

<style lang = "css" scoped>
  .Main {
    width: 100%;
    height: 100%;
    display: flex;
  }

  .FileList {
    height: 100%;
    background: var(--ccs-mde-fileList-color);
  }

  .SplitLine {
    width: var(--ccs-mde-splitLine-width);
    height: 100%;
    cursor: col-resize;
  }

  .Editor {
    height: 100%;
  }
</style>

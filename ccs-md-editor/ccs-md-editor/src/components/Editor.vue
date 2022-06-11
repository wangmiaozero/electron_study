<template>
  <div class = "Editor">
    <div style = "width: 100%;">
      <div class = "Title">
        <input v-model = "form.title" placeholder = "无标题笔记" @input = "() => { SaveNote(id) }"/>
      </div>
    </div>
    <div id = "NoteEditor" ref = "NoteEditor"></div>
    <div class = "NoteCount">
      <span>输入字数: {{ noteCount }}</span>
    </div>
  </div>
</template>

<script>
import Vditor from 'vditor'
const ipcRenderer = window.electron.ipcRenderer

export default {
  props: {
    id: { type: String }
  },
  data () {
    return {
      form: {
        title: '',
        content: ''
      },
      noteEditor: {},
      noteCount: 0
    }
  },
  methods: {
    InitEditor: function () {
      const _this = this
      const url = 'http://localhost:41220/Note/UploadImage'
      this.noteEditor = new Vditor('NoteEditor', {
        cdn: '/js/vditor',
        width: 'calc(100% - 2px)',
        height: '100%',
        theme: 'dark',
        placeholder: '快些记录下你的灵感吧',
        outline: { enable: true },
        toolbar: [
          'headings', 'bold', 'italic', 'strike', 'link',
          '|',
          'list', 'ordered-list', 'check', 'outdent', 'indent',
          '|',
          'quote', 'line', 'code', 'inline-code', 'insert-before', 'insert-after',
          '|',
          'upload', 'table',
          '|',
          'undo', 'redo',
          '|',
          'edit-mode',
          {
            hotkey: '⌘N',
            name: 'newNote',
            tipPosition: 's',
            tip: '新建笔记',
            className: 'right',
            icon: '<svg t="1642927183068" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="2159" width="15" height="15"><path d="M1015.42418895 306.72863086a39.93795454 39.93795454 0 0 0-10.0734527-19.1614588L717.92885638 3.1563774a39.88320813 39.88320813 0 0 0-36.10566298-10.40193423H152.96693232a148.11807594 148.11807594 0 0 0-104.23833499 42.70267937A145.38072483 145.38072483 0 0 0 5.56056787 138.65526451v747.57062501a145.38072483 145.38072483 0 0 0 43.16802946 103.19814096A148.11807594 148.11807594 0 0 0 152.96693232 1032.12671086h718.0346047a148.06332851 148.06332851 0 0 0 104.21096228-42.70268038A145.38072483 145.38072483 0 0 0 1018.38052876 886.22588952V320.96285765a36.81737448 36.81737448 0 0 0-2.95633981-14.23422679zM723.59517257 122.77862705l163.61148445 161.77745879h-163.61148445V122.77862705z m221.09586084 763.44726247a73.4978814 73.4978814 0 0 1-73.68949639 73.08727771H152.96693232a73.4978814 73.4978814 0 0 1-73.71686908-73.08727771V138.65526451A73.25151898 73.25151898 0 0 1 152.96693232 65.8417209h496.91137117v255.12113675a36.54363937 36.54363937 0 0 0 36.8721219 36.40677079h257.940608v528.85626108zM521.44178387 311.92959797a36.84474819 36.84474819 0 0 0-73.68949536 0v182.30759314H263.50117587a36.40677181 36.40677181 0 1 0 0 72.8135426h184.25111264v182.30759315a36.84474819 36.84474819 0 0 0 73.68949536 0v-182.30759315h184.25111162a36.40677181 36.40677181 0 1 0 0-72.8135426h-184.25111162v-182.30759314z m0 0" p-id="2160" fill="#d1d2d5"></path></svg>',
            click () { _this.$emit('NewNote') }
          },
          {
            hotkey: '⌘S',
            name: 'saveNote',
            tipPosition: 's',
            tip: '保存笔记',
            className: 'right',
            icon: '<svg t="1640486568771" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="1925" width="15" height="15"><path d="M782.24893527-2.53674976H55.84411208c-34.31933527 0-62.30668986 27.86071497-62.30668986 62.18005025v901.42076908c0 34.31933527 27.98735459 62.18005025 62.30668986 62.18005024H960.30423189c34.31933527 0 62.30668986-27.86071497 62.30668985-62.18005024V236.93875943L782.24893527-2.53674976z m-505.79861643 50.65584542h429.05501063v306.34122512c0 16.33651014-13.29715942 29.50702995-29.76030918 29.50702995H306.21062801c-16.33651014 0-29.76030918-13.29715942-29.76030917-29.50702995V48.11909566z m6.83853913 921.42982801V639.2728116c0-5.19222415 4.30574687-9.37133139 9.49797101-9.37133141h460.20835556c5.19222415 0 9.49797101 4.17910725 9.49797103 9.37133141v330.27611207H283.28885797z m688.66621836-8.61149372c0 6.20534106-5.31886377 11.52420483-11.65084444 11.52420483H813.27564058V639.2728116c0-33.05293913-26.97423768-60.02717682-60.15381643-60.02717682H292.91346861c-33.17957875 0-60.15381643 26.97423768-60.15381643 60.02717682v333.3154628H55.84411208c-6.33198068 0-11.65084444-5.31886377-11.65084444-11.52420483V59.64330049c0-6.20534106 5.31886377-11.52420483 11.65084444-11.52420483H225.79447343v306.34122512c0 44.19722513 36.09228986 80.16287536 80.41615458 80.16287536h369.53439228c44.32386473 0 80.41615459-35.96565024 80.41615459-80.16287536V48.11909566h5.06558454l210.72831691 209.96847921v702.84985508z" fill="#b9b9b9" p-id="1926"></path><path d="M672.07247149 704.49221256H371.17674975c-13.93035748 0-25.3279227 11.39756522-25.3279227 25.32792271s11.39756522 25.3279227 25.3279227 25.3279227h300.89572174c13.93035748 0 25.3279227-11.39756522 25.32792272-25.3279227s-11.2709256-25.3279227-25.32792272-25.32792271zM672.07247149 832.52486183H371.17674975c-13.93035748 0-25.3279227 11.39756522-25.3279227 25.32792271s11.39756522 25.3279227 25.3279227 25.32792271h300.89572174c13.93035748 0 25.3279227-11.39756522 25.32792272-25.32792271s-11.2709256-25.3279227-25.32792272-25.32792271z" fill="#b9b9b9" p-id="1927"></path></svg>',
            click () { _this.SaveNote(_this.id) }
          },
          {
            name: 'more',
              toolbar: [
                'both',
                  // 'export',          // 导出
                  'outline',         // 大纲
                  'preview',         // 预览
              ],
          }
        ],
        counter: { enable: true, after(length) { _this.noteCount = length } },
        after: () => {
          this.noteEditor.setTheme('dark', 'dark', 'native', '/js/vditor/dist/css/content-theme')
          this.noteEditor.setValue(this.form.content)
        },
        input: () => {
          this.form.content = this.noteEditor.getValue()
          this.SaveNote(this.id)
        },
        upload: {
          accept:'image/jpg, image/jpeg, image/png, image/gif',
          url: url,
          multiple: false,
          fieldName: 'imageFile',
          linkToImgUrl: url
        },
        preview: {
          parse: () => { this.LinkFail() }
        }
      })
    },
    GetNote: function () {
      return new Promise((resolve) => {
        this.form = { title: '', content: '' }

        const url = 'http://localhost:41220/Note/Get'
        const data = { id: this.id }
        this.$http.post(url, data).then((res) => {
          const { title, content } = res.data
          this.form = { title, content }
          resolve()
        })
      })
    },
    SaveNote: function (id) {
      const url = 'http://localhost:41220/Note/Save'
      const { title, content } = this.form
      const data = { id, title, content }
      this.$http.post(url, data).then(() => {
        this.$emit('GetNoteList')
      })
    },
    LinkFail: function () {
      var linkEventList = [...document.getElementsByClassName('vditor-ir__link')]
      linkEventList.push(...document.getElementsByTagName('a'))
      if (linkEventList.length > 0) {
        for (let index = 0; index < linkEventList.length; index ++) {
          linkEventList[index].onclick = (e) => {
            e.stopPropagation()

            let url = ''
            if (e.currentTarget.tagName === 'A') {
              // a 标签跳转
              url = e.currentTarget.href
            } else if (e.currentTarget.tagName === 'SPAN') {
              // span 标签跳转
              const childNodes = e.currentTarget.parentElement.childNodes
              childNodes.forEach(element => {
                if (element.className === 'vditor-ir__marker vditor-ir__marker--link') {
                  url = element.innerText
                }
              })
            }

            url = url || ''
            if (url !== '') { ipcRenderer.send('openUrl', url) }
            return false
          }
        }
      }
    }
  },
  mounted () {
    this.GetNote().then(() => { this.InitEditor() })
  },
  updated () { this.LinkFail() }
}
</script>

<style lang = "css" scoped>
  .Editor {
    width: 100%;
    height: calc(100% - 10px - 10px) !important;
  }

  .Title {
    padding: 10px 20px;
    margin: 10px 10px;
  }

  .Title input {
    width: 100%;
    margin: 0;
    padding: 0;
    background: unset;
    border: unset;
    color: var(--ccs-mde-primary-font-color);
    font-size: 20px;
    font-weight: 900;
    text-align: center;
  }

  .Title input:focus-visible {
    outline: unset;
  }

  #NoteEditor {
    height: calc(100% - 48px - 10px - 30px);
  }

  .vditor {
    border: unset;
    border-radius: unset;
  }

  .vditor-reset:focus-visible {
    outline: unset;
  }

  #NoteEditor >>> .vditor-toolbar {
    padding-left: unset !important;
    border-bottom: 0;
    display: flex;
    justify-content: center;
    padding-bottom: 10px;
  }

  #NoteEditor >>> .vditor-outline {
    width: 300px;
    margin-right: 30px;
  }

  #NoteEditor >>> .vditor-outline__title {
    display: none;
  }

  #NoteEditor >>> .vditor-outline li > span {
    padding: 2px 10px;
  }

  #NoteEditor >>> .vditor-outline li > span > span {
    color: var(--ccs-mde-unimportant-font-color);
    padding: 2px 8px;
    border-radius: 50px;
  }

  #NoteEditor >>> .vditor-outline li > span > span:hover {
    font-weight: 900;
    color: var(--ccs-mde-primary-font-color);
  }

  #NoteEditor >>> .vditor-preview {
    border-left: unset;
  }

  #NoteEditor >>> .vditor-counter {
    display: none;
  }

  .NoteCount {
    --padding: 10px;
    width: calc(100% - var(--padding));
    text-align: right;
    padding-top: var(--padding);
    padding-right: var(--padding);
    line-height: 20px;
    font-size: 14px;
    color: var(--ccs-mde-unimportant-font-color);
  }
</style>
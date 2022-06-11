package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/nu7hatch/gouuid"
)

// sqlite3
var db = new(sql.DB)

// 数据文件路径
var dataFolderPath string
var imageFolderPath string

// 设置数据目录
func setDataFolderPath() {
	// 数据目录
	userFolderPath, _ := os.UserHomeDir()
	dataFolderPath = fmt.Sprintf("%s/CCSMdEditor", userFolderPath)
	_, err := os.Stat(dataFolderPath)
	if err != nil {
		// 当文件夹不存在时, 自动创建文件夹
		os.Mkdir(dataFolderPath, os.ModePerm)
	}

	// 图片目录
	imageFolderPath := fmt.Sprintf("%s/Image", dataFolderPath)
	if err != nil {
		// 当文件夹不存在时, 自动创建文件夹
		os.Mkdir(imageFolderPath, os.ModePerm)
	}
}

// 全局数据库指针
func openDB() *sql.DB {
	// 打开数据库, 不存在则创建
	db, _ = sql.Open("sqlite3", dataFolderPath+"/note.db")
	return db
}

// 初始化数据库
func initDB() {
	statement := `
		CREATE TABLE IF NOT EXISTS T_NOTE (
			id CHAR(36) PRIMARY KEY,
			title VARCHAR(256) NULL,
			content TEXT NULL,
			state INT NULL,
			updateDateTime DATETIME NULL
		);
	`
	db.Exec(statement)
}

// 接口实现
// 读取笔记列表
func getNoteList(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	type Filter struct {
		Search string `json:"search" bson:"search"`
	}

	var filter Filter
	json.Unmarshal(body, &filter)
	filterSearch := filter.Search

	statement := fmt.Sprintf(`
		SELECT
			id,
			title,
			SUBSTRING(content, 1, 50) AS summary,
			strftime("%%Y-%%m-%%d %%H:%%M:%%f", updateDateTime) AS updateDateTime
		FROM T_NOTE
		WHERE (title LIKE '%%%s%%' OR content LIKE '%%%s%%') AND state = 1
		ORDER BY updateDateTime DESC
	`, filterSearch, filterSearch)

	noteList := []map[string]interface{}{}
	rows, _ := db.Query(statement)
	for rows.Next() {
		var id string
		var title string
		var summary string
		var updateDateTime string
		rows.Scan(&id, &title, &summary, &updateDateTime)
		if len(summary) > 50 {
			summary += "..."
		}
		noteList = append(noteList, map[string]interface{}{
			"id":             id,
			"title":          title,
			"summary":        summary,
			"updateDateTime": updateDateTime,
		})
	}

	resJson, _ := json.Marshal(noteList)
	fmt.Fprint(w, string(resJson))
}

// 创建笔记
func newNote(w http.ResponseWriter, r *http.Request) {
	// 创建笔记 id
	_id, _ := uuid.NewV4()
	id := _id.String()
	title := "无标题笔记"
	content := ""
	state := "1"
	updateDateTime := time.Now().Format("2006-01-02 15:04:05.000")

	// 写入笔记数据
	stmt, _ := db.Prepare("INSERT INTO T_NOTE(id, title, content, state, updateDateTime) VALUES(?, ?, ?, ?, ?)")
	stmt.Exec(id, title, content, state, updateDateTime)

	fmt.Fprint(w, id)
}

// 读取笔记
func getNote(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	type NoteId struct {
		Id string `json:"id" bson:"id"`
	}

	var noteId NoteId
	json.Unmarshal(body, &noteId)
	id := noteId.Id

	statement := fmt.Sprintf(`
		SELECT title, content FROM T_NOTE WHERE id = '%s' AND state = 1 LIMIT 1
	`, id)

	rows, _ := db.Query(statement)
	title := ""
	content := ""
	for rows.Next() {
		var _title string
		var _content string
		rows.Scan(&_title, &_content)
		title = _title
		content = _content
	}

	resJson, _ := json.Marshal(map[string]interface{}{
		"id":      id,
		"title":   title,
		"content": content,
	})
	fmt.Fprint(w, string(resJson))
}

// 保存笔记
func saveNote(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	type Note struct {
		Id      string `json:"id" bson:"id"`
		Title   string `json:"title" bson:"title"`
		Content string `json:"content" bson:"content"`
	}

	var note Note
	json.Unmarshal(body, &note)
	id := note.Id
	title := note.Title
	if title == "" {
		title = "无标题笔记"
	}
	content := note.Content
	updateDateTime := time.Now().Format("2006-01-02 15:04:05.000")

	stmt, _ := db.Prepare("UPDATE T_NOTE SET title=?, content=?, updateDateTime=? WHERE id=?")
	stmt.Exec(title, content, updateDateTime, id)

	fmt.Fprint(w, id)
}

// 删除笔记
func removeNote(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	type NoteId struct {
		Id string `json:"id" bson:"id"`
	}

	var noteId NoteId
	json.Unmarshal(body, &noteId)
	id := noteId.Id
	updateDateTime := time.Now().Format("2006-01-02 15:04:05.000")

	stmt, _ := db.Prepare("UPDATE T_NOTE SET updateDateTime=?, state=0 WHERE id=?")
	stmt.Exec(updateDateTime, id)

	fmt.Fprint(w, id)
}

// 上传图片
func uploadImage(w http.ResponseWriter, r *http.Request) {
	f, fh, err := r.FormFile("imageFile")

	type ImageResData struct {
		Code int                    `json:"code" bson:"code"`
		Msg  string                 `json:"msg" bson:"msg"`
		Data map[string]interface{} `json:"data" bson:"data"`
	}

	var resData ImageResData
	if err != nil {
		resData.Code = 400
		resData.Msg = "上传失败, 该图片可能是外链, 建议图片保存后再上传"
		resJson, _ := json.Marshal(resData)
		fmt.Fprintln(w, string(resJson))
		return
	}

	// 检查文件格式
	fileExt := strings.ToLower(path.Ext(fh.Filename))
	if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
		resData.Code = 400
		resData.Msg = "上传失败! 只允许 png, jpg, gif, jpeg 文件"
		resJson, _ := json.Marshal(resData)
		fmt.Fprintln(w, string(resJson))
		return
	}

	// 根据文件内容生成 md5 作为文件名, 避免上传重复文件
	// 图片将会上传到 /Static/Upload/Note 相应月份的文件夹中
	md5 := md5.New()
	fObj, _ := fh.Open()
	fByte, _ := ioutil.ReadAll(fObj)
	md5.Write(fByte)
	fileName := hex.EncodeToString(md5.Sum(nil))
	fileName += fileExt

	// 创建文件夹
	folderName := time.Now().Format("200601")
	folderName = fmt.Sprintf("Image/%s", folderName)
	folderPath := dataFolderPath + "/" + folderName
	filePath := folderPath + "/" + fileName
	url := filePath

	_, err = os.Stat(folderPath)
	if err != nil {
		// 当文件夹不存在时, 自动创建文件夹
		os.Mkdir(folderPath, os.ModePerm)
	}

	// 保存图片
	saveFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		resData.Code = 400
		resData.Msg = "上传失败!"
		resJson, _ := json.Marshal(resData)
		fmt.Fprintln(w, string(resJson))
		return
	}

	io.Copy(saveFile, f)
	defer f.Close()
	defer saveFile.Close()

	resData.Code = 200
	resData.Msg = "上传成功"
	resData.Data = map[string]interface{}{
		"errFiles": []string{},
		"succMap":  map[string]string{fileName: url},
	}
	resJson, _ := json.Marshal(resData)
	fmt.Fprintln(w, string(resJson))
}

// 程序入口
func main() {
	// 创建数据文件夹
	setDataFolderPath()

	// 初始化数据库
	db = openDB()
	initDB()

	// 接口
	http.HandleFunc("/Note/List", getNoteList)
	http.HandleFunc("/Note/New", newNote)
	http.HandleFunc("/Note/Get", getNote)
	http.HandleFunc("/Note/Save", saveNote)
	http.HandleFunc("/Note/Remove", removeNote)
	http.HandleFunc("/Note/UploadImage", uploadImage)

	// 绑定指定端口
	http.ListenAndServe(":41220", nil)
}

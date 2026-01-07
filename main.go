package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"lab1/TreeAdapter"
	"lab1/common"
	"lab1/editor"
	"lab1/log"
	"lab1/statistics"
	"lab1/storage"
	"lab1/workspace"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
// 计时器绑定
var timeStatistics = &statistics.Statistics{}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}
func main() {
	// 1. 初始化依赖组件
	fileStorage := storage.NewLocalStorage("./workspace_state.json") // 状态存储路径
	logModule := log.NewLogModule()

	// 2. 初始化工作区
	ws := workspace.NewWorkspace("./workspace_state.json")

	// 3. 日志模块订阅工作区事件（观察者模式）
	ws.RegisterObserver(logModule)

	//计时器绑定
	timeStatistics = statistics.NewStatistics()
	ws.RegisterObserver(timeStatistics)
	//timeStatistics.GetFormattedDuration()
	//日志模块订阅编辑器事件

	// 4. 从本地存储恢复上次工作区状态（备忘录模式）
	if err := restoreWorkspaceState(ws, fileStorage); err != nil {
		fmt.Printf("恢复工作区失败，使用新状态: %v\n", err)
	} else {
		fmt.Println("工作区已恢复上次状态")
	}

	startInteractiveLoop(ws)
}

// 修复后的 restoreWorkspaceState 函数
func restoreWorkspaceState(ws *workspace.Workspace, storage *storage.LocalStorage) error {
	// 调用 Workspace 的 RestoreState 方法，传入编辑器工厂函数
	// 工厂函数复用之前定义的 editor.EditorFactory（需确保已导入 editor 包）
	fmt.Println("restoreWorkspaceState")
	return ws.RestoreState(editor.EditorFactory)
}

// 启动用户交互循环
func startInteractiveLoop(ws *workspace.Workspace) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("编辑器启动完成，支持指令: load/save/close/undo/exit....")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		handleCommand(ws, input, true)
		//fmt.Printf("[debug]active_file: %s\n", ws.GetActiveEditor().GetFilePath())
		activeEditor := ws.GetActiveEditor()
		if activeEditor == nil {
			fmt.Println("[debug]active_file: 无激活的编辑器/文件")
		} else {
			fmt.Printf("[debug]active_file: %s\n", activeEditor.GetFilePath())
		}
	}
}

// 处理用户指令
func handleCommand(ws *workspace.Workspace, input string, debug bool) {
	input = strings.TrimSpace(input)
	parts := strings.SplitN(input, " ", 10)
	if len(parts) == 0 {
		fmt.Println("无效指令")
		return
	}
	cmd := parts[0]
	switch cmd {
	case "load": //完成
		_load(ws, parts, false)
	case "save": //完成
		_Save(ws, input, debug, parts)
	case "close": //完成
		_close(ws, parts)
	case "init": //完成
		_init(ws, input)
	case "undo":
		_undo(ws)
	case "redo":
		_redo(ws)
	case "editor-list": //完成
		_EditorList(ws)
	case "edit": //完成
		_edit(ws, parts)
	case "exit":
		_exit(ws)
	case "dir-tree": //完成
		//_dirTree(ws, parts)
		_dirTreeV2(ws, parts)
	case "append":
		_append(ws, parts)
	case "insert":
		_insert(ws, parts)
	case "show":
		_show(ws, parts)
	case "delete":
		// 先尝试按XML指令处理，若不是XML编辑器，再按文本指令处理
		if !_xmlDelete(ws, parts) {
			_delete(ws, parts)
		}
	case "replace":
		_replace(ws, parts)
	case "log-on":
		_LogOn(ws, parts)
	case "log-off":
		_LogOff(ws, parts)
	case "log-show":
		_LogShow(ws, parts)
	//case :
	case "insert-before":
		_insertBefore(ws, parts)
	case "append-child":
		_appendChild(ws, parts)
	case "edit-id":
		_editId(ws, parts)
	case "edit-text":
		_editText(ws, parts)
	case "xml-tree":
		//_xmlTree(ws, parts)
		_xmlTreeV2(ws, parts)
	case "spell-check":
		_spellCheck(ws, parts)
	default:
		fmt.Println("未知指令，支持: load/save/close/undo/exit")
	}
}


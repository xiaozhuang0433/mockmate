package watcher

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher 文件监听器
type Watcher struct {
	watcher   *fsnotify.Watcher
	callbacks map[string][]func()
	debounce  map[string]time.Time
}

// NewWatcher 创建新的文件监听器
func NewWatcher() (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		watcher:   w,
		callbacks: make(map[string][]func()),
		debounce:  make(map[string]time.Time),
	}, nil
}

// WatchDir 监听目录
func (w *Watcher) WatchDir(dir string, callback func()) error {
	// 规范化路径
	dir = filepath.Clean(dir)

	// 监听目录
	err := w.watcher.Add(dir)
	if err != nil {
		return err
	}

	// 记录回调函数
	w.callbacks[dir] = append(w.callbacks[dir], callback)

	log.Printf("开始监听目录: %s", dir)
	return nil
}

// WatchFile 监听文件
func (w *Watcher) WatchFile(file string, callback func()) error {
	// 规范化路径
	file = filepath.Clean(file)

	// 监听文件所在目录
	dir := filepath.Dir(file)
	err := w.watcher.Add(dir)
	if err != nil {
		return err
	}

	// 记录文件特定的回调
	w.callbacks[file] = append(w.callbacks[file], callback)

	log.Printf("开始监听文件: %s", file)
	return nil
}

// Start 启动监听
func (w *Watcher) Start() {
	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}

				log.Printf("检测到文件事件: %s, 操作: %s", event.Name, event.Op)

				// 只处理写入和创建事件
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					w.handleEvent(event.Name)
				}

			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				log.Printf("监听错误: %v", err)
			}
		}
	}()
}

// handleEvent 处理文件事件
func (w *Watcher) handleEvent(filename string) {
	filename = filepath.Clean(filename)

	// 防抖：如果同一文件在 500ms 内重复触发，忽略
	if lastTime, ok := w.debounce[filename]; ok {
		if time.Since(lastTime) < 500*time.Millisecond {
			return
		}
	}
	w.debounce[filename] = time.Now()

	// 查找直接匹配的文件回调
	if callbacks, ok := w.callbacks[filename]; ok {
		log.Printf("触发文件回调: %s", filename)
		for _, cb := range callbacks {
			go cb()
		}
		return
	}

	// 查找目录回调
	dir := filepath.Dir(filename)
	if callbacks, ok := w.callbacks[dir]; ok {
		log.Printf("触发目录回调: %s (文件: %s)", dir, filename)
		for _, cb := range callbacks {
			go cb()
		}
		return
	}

	// 尝试检查是否是监听目录下的文件
	for watchedDir, callbacks := range w.callbacks {
		if strings.HasPrefix(filename, watchedDir+string(filepath.Separator)) {
			log.Printf("触发目录回调: %s (文件: %s)", watchedDir, filename)
			for _, cb := range callbacks {
				go cb()
			}
			return
		}
	}

	log.Printf("未找到匹配的回调: %s", filename)
}

// Close 关闭监听器
func (w *Watcher) Close() error {
	return w.watcher.Close()
}

package main

import (
	"context"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var mainwin *ui.Window

var msgEntry *ui.MultilineEntry

func myUi() ui.Control {

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	input1 := ui.NewEntry()
	input2 := ui.NewEntry()
	input2.SetText("10")

	btn := ui.NewButton("开始")
	btn.OnClicked(func(this *ui.Button) {
		sendMsg(msgEntry, "开始启动------------------------------------")
		getUrlData(input1, input2)
		//this.Disable()
	})

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	entryForm.Append("网址", input1, false)
	entryForm.Append("并发", input2, false)
	entryForm.Append("", btn, true)


	hbox.Append(entryForm, true)

	// 底部消息框
	hbox1 := ui.NewHorizontalBox()
	hbox1.SetPadded(true)
	msgEntry = ui.NewMultilineEntry()
	msgEntry.SetReadOnly(true)
	hbox1.Append(msgEntry, true)

	vbox.Append(hbox, false)
	vbox.Append(hbox1, true)

	return vbox
}

func setupUI() {
	mainwin = ui.NewWindow("网站压测", 450, 580, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	//tab := ui.NewTab()
	mainwin.SetChild(myUi())
	mainwin.SetMargined(true)

	mainwin.Show()
}



func main() {
	ui.Main(setupUI)
}

func sendMsg(box *ui.MultilineEntry, msg string) {
	timeStr := time.Now().Format("01-02 15:04:05")
	box.Append(timeStr + " " + msg + "\n")
}

func getUrlData(input1, input2 *ui.Entry) bool {

	url := input1.Text()
	curr := input2.Text()

	if url == "" {
		sendMsg(msgEntry, "网址不能为空")
		return false
	}

	if !strings.Contains(url, "http") {
		sendMsg(msgEntry, "网址格式不对")
		return false
	}

	validUrl := regexp.MustCompile(`(http|ftp|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?`)

	if !validUrl.MatchString(url) {
		sendMsg(msgEntry, "网址格式不对.")
		return false
	}


	if curr == "" {
		sendMsg(msgEntry, "并发数量不能为空")
		return false
	}

	qty, err := strconv.Atoi(curr)
	if err != nil {
		sendMsg(msgEntry, "必须是整数")
		return false
	}

	if qty <= 0 || qty > 100 {
		sendMsg(msgEntry, "并发数量在1~100之间 ")
		return false
	}



	ctx, _ := context.WithTimeout(context.TODO(), 20*time.Second)
	// 并发请求http get
	for i := 1; i <= qty; i++ {
		go getEachUrl(ctx, i, url)
	}
	return true

}

// 请求
func getEachUrl(ctx context.Context, i int, url string) {

	defer func() {
		if err := recover(); err != nil {
			sendMsg(msgEntry, "错误退出，第"+strconv.Itoa(i)+"次")
		}
	}()

	str := "正在发起第" + strconv.Itoa(i) + "次请求"
	ui.QueueMain(func() {
		sendMsg(msgEntry, str)
	})
	// http.get
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		str = "第" + strconv.Itoa(i) + "次请求无响应"
		return
	}
	defer resp.Body.Close()


	for {
		select {
		case <-ctx.Done():
			str = "第" + strconv.Itoa(i) + "次请求超时退出"
			ui.QueueMain(func() {
				sendMsg(msgEntry, str)
			})
			return
		default:

			val := time.Since(start)
			str = "第" + strconv.Itoa(i) + "次请求耗时" + val.String()
			ui.QueueMain(func() {
				sendMsg(msgEntry, str)
			})

			time.Sleep(200 * time.Millisecond)
			return
		}
	}


}
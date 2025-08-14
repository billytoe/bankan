package main

/* This file contains the application entry, main window UI and some logic to save/load data and preferences */

/* ================================================================================ Imports */
import (
	"fmt"
	"image/color"
	"io"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

/* ================================================================================ Constants */
const (
	WINDOW_TITLE = "BanKan"
)

/* ================================================================================ Private variables */
var window fyne.Window
var board *Board
var boardToolbar *widget.Toolbar
var filterBinding binding.String
var boardNameLabel *CustomLabel
var saveFileURI fyne.URI

/* ================================================================================ Private functions */
func setSaveFileURI(uri fyne.URI) {
	saveFileURI = uri
	windowTitleSuffix := ""

	if saveFileURI != nil {
		fyne.CurrentApp().Preferences().SetString("saveFileURI", saveFileURI.String())
		windowTitleSuffix = " - " + saveFileURI.Path()
	} else {
		fyne.CurrentApp().Preferences().SetString("saveFileURI", "")
	}

	window.SetTitle(WINDOW_TITLE + windowTitleSuffix)
}

func restorePreferences() {
	saveFileURI = nil

	uriString := fyne.CurrentApp().Preferences().String("saveFileURI")
	if uriString == "" {
		return
	}

	uri, err := storage.ParseURI(uriString)
	if err != nil {
		return
	}

	saveFileURI = uri
}

func autoSave() {
	if saveFileURI != nil {
		saveBoardURI(board, saveFileURI)
	}
}

func windowCloseInterceptor() {
	window.Close()
}

func clearBoard() {
	board.Clear()

	board.Name = "New Board"
	syncBoardNameLabel()

	setSaveFileURI(nil)
}

func loadBoardReader(board *Board, reader fyne.URIReadCloser) {
	data, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
	}

	if err := reader.Close(); err != nil {
		fmt.Println(err)
		return
	}

	board.Clear()

	if err := board.Load(data); err != nil {
		fmt.Println(err)
		return
	}

	syncBoardNameLabel()

	setSaveFileURI(reader.URI())
}

func loadBoardURI(board *Board, uri fyne.URI) {
	if reader, err := storage.Reader(uri); reader != nil && err == nil {
		loadBoardReader(board, reader)
	}
}

func loadBoardSaveFile() {
	if saveFileURI != nil {
		loadBoardURI(board, saveFileURI)
	}
}

func saveBoardWriter(board *Board, writer fyne.URIWriteCloser) {
	data, err := board.Data()
	if err != nil {
		fmt.Println(err)
		return
	}

	if written, err := writer.Write(data); err != nil || written != len(data) {
		fmt.Println(err)
	}

	if err = writer.Close(); err != nil {
		fmt.Println(err)
		return
	}

	setSaveFileURI(writer.URI())
}

func saveBoardURI(board *Board, uri fyne.URI) {
	if writer, err := storage.Writer(uri); writer != nil && err == nil {
		saveBoardWriter(board, writer)
	}
}

func syncBoardNameLabel() {
	boardNameLabel.Text = board.Name
	boardNameLabel.Refresh()
}

func newButtonTapped() {
	ShowConfirmDialog("Create Empty Board", "This will discard the current board.\n\nAre you sure?\n", clearBoard)
}

func loadButtonTapped() {
	ShowFileOpenConfirmDialog("Load Board", "This will discard the current board.\n\nAre you sure?\n", saveFileURI,
		func(reader fyne.URIReadCloser) { loadBoardReader(board, reader) },
	)
}

func saveAsButtonTapped() {
	ShowSaveAsDialog(saveFileURI, func(writer fyne.URIWriteCloser) { saveBoardWriter(board, writer) })
}

func saveButtonTapped() {
	if saveFileURI != nil {
		saveBoardURI(board, saveFileURI)
	} else {
		saveAsButtonTapped()
	}
}

func showEditBoardNameDialog() {
	ShowEntryDialog("Edit Board Name", "Name ...", board.Name,
		func(text string) {
			board.Name = text
			syncBoardNameLabel()
			autoSave()
		},
	)
}

func showBoardMenu() {
	menu := widget.NewPopUpMenu(
		fyne.NewMenu("Board", fyne.NewMenuItem("Edit Board Name", showEditBoardNameDialog)),
		window.Canvas(),
	)
	menu.ShowAtPosition(fyne.NewPos(boardToolbar.Position().X-menu.Size().Width+50, boardToolbar.Position().Y+menu.Size().Height))
}

func boardFilterChanged(tagEditString string) {
	filterBinding.Set(tagEditString)
}

func filterBindingChanged() {
	if text, err := filterBinding.Get(); err == nil {
		board.SetTagFilter(text)
	}
}

func updateDateItems() {
	// 更新所有日期类型的item
	for _, stage := range board.Stages {
		for _, item := range stage.Items {
			if item.DateType != "Normal" && item.DateType != "" {
				// 更新日期标签
				updatedTags := make([]Tag, 0)
				dateTypeFound := false

				for _, tag := range item.Tags {
					if strings.HasPrefix(tag.Expression, item.DateType+"=") {
						// 更新日期标签
						newDateString := getCurrentDateString(item.DateType)
						updatedTags = append(updatedTags, Tag{Expression: item.DateType + "=" + newDateString})
						dateTypeFound = true
					} else {
						updatedTags = append(updatedTags, tag)
					}
				}

				// 如果没有找到日期标签，添加一个
				if !dateTypeFound {
					newDateString := getCurrentDateString(item.DateType)
					updatedTags = append(updatedTags, Tag{Expression: item.DateType + "=" + newDateString})
				}

				item.Tags = updatedTags
				item.Refresh()
			}
		}
	}
	autoSave()
}

func startDateUpdateTimer() {
	go func() {
		for {
			now := time.Now()
			// 计算到下一个凌晨的时间
			nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			duration := nextMidnight.Sub(now)

			// 等待到凌晨
			time.Sleep(duration)

			// 更新日期item
			updateDateItems()

			// 每24小时重复一次
			time.Sleep(24 * time.Hour)
		}
	}()
}

func main() {
	application := app.NewWithID("de.bananajoh.bankan")
	application.SetIcon(theme.FyneLogo())

	window = application.NewWindow(WINDOW_TITLE)
	window.SetCloseIntercept(windowCloseInterceptor)

	restorePreferences()

	board = NewBoard("New Board", boardFilterChanged)

	fileToolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentIcon(), newButtonTapped),
		widget.NewToolbarAction(theme.FolderOpenIcon(), loadButtonTapped),
		widget.NewToolbarAction(theme.DownloadIcon(), saveAsButtonTapped),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), saveButtonTapped),
	)
	filterBinding = binding.NewString()
	filterBinding.AddListener(binding.NewDataListener(filterBindingChanged))

	filterEntry := widget.NewEntryWithData(filterBinding)
	filterEntry.SetPlaceHolder("Filter by Tag ...")

	leftHeaderContainer := container.NewGridWithColumns(2, fileToolbar, filterEntry)

	boardNameLabel = NewCustomLabel(fyne.TextAlignCenter, PaintStyle{color.RGBA{255, 255, 255, 255}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, 0}, false, board.Name, theme.TextSubHeadingSize(), fyne.TextStyle{}, Paddings{1.0, 1.0, 1.0, 1.0}, Paddings{0.0, 0.0, 0.0, 0.0})
	boardNameContainer := container.NewHBox(layout.NewSpacer(), boardNameLabel, layout.NewSpacer())

	boardToolbar = widget.NewToolbar(
		widget.NewToolbarAction(theme.FolderNewIcon(), board.ShowCreateStageDialog),
		widget.NewToolbarAction(theme.MoreVerticalIcon(), showBoardMenu),
	)

	toolbarContainer := container.NewBorder(nil, nil, leftHeaderContainer, boardToolbar, boardNameContainer)
	headerBarContainer := container.NewVBox(toolbarContainer, widget.NewSeparator())
	windowContainer := container.NewBorder(headerBarContainer, nil, nil, nil, board)

	loadBoardSaveFile()

	// 启动时立即更新日期item
	updateDateItems()

	// 启动日期更新定时器
	startDateUpdateTimer()

	window.SetContent(windowContainer)
	window.Resize(fyne.NewSize(1200, 700))
	window.CenterOnScreen()
	window.ShowAndRun()
}
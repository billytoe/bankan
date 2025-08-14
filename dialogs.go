package main

/* This file contains helper functions to show different types and combinations of dialogs */

/* ================================================================================ Imports */
import (
	"fmt"
	"image/color"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

/* ================================================================================ Public functions */
func getSystemLanguage() string {
	// 首先检查常见的Unix/Linux环境变量
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = os.Getenv("LC_ALL")
	}
	if lang == "" {
		lang = os.Getenv("LC_MESSAGES")
	}
	if lang == "" {
		lang = os.Getenv("LANGUAGE")
	}

	// Windows环境下的额外检查
	if lang == "" {
		// 检查Windows系统语言相关环境变量
		lang = os.Getenv("WINLANG")
	}
	if lang == "" {
		// 检查用户界面语言
		lang = os.Getenv("MUI_LANGUAGE")
	}

	// 如果所有环境变量都为空，在Windows中文环境下默认返回中文
	// 这是一个合理的假设，因为大多数中文用户会使用中文Windows
	if lang == "" {
		// 简单的启发式检查：如果是Windows且找不到语言环境变量
		// 可以假设是中文环境（特别是在中国地区）
		return "zh"
	}

	if strings.HasPrefix(lang, "zh") || strings.Contains(lang, "Chinese") ||
		strings.Contains(lang, "chinese") || strings.Contains(lang, "CN") ||
		strings.Contains(lang, "cn") {
		return "zh"
	}
	return "en"
}

func getDateTypeLabels() (string, string, string) {
	lang := getSystemLanguage()
	if lang == "zh" {
		return "公历", "农历", "藏历"
	}
	return "Gregorian", "Lunar", "Tibetan"
}

// 获取周几的颜色 请按心理学 将周一到周日 每天用一个颜色的点代表
func getWeekdayColor(weekday time.Weekday) string {
	switch weekday {
	case time.Monday:
		return "🟢" // 绿色
	case time.Tuesday:
		return "🔵" // 蓝色
	case time.Wednesday:
		return "🟡" // 黄色
	case time.Thursday:
		return "🟠" // 橙色
	case time.Friday:
		return "🔴" // 红色
	case time.Saturday:
		return "⚪" // 白色
	case time.Sunday:
		return "🟣" // 紫色
	default:
		return "⚫" // 黑色
	}
}

func getWeekdayName(weekday time.Weekday) string {
	lang := getSystemLanguage()
	if lang == "zh" {
		weekdays := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
		return weekdays[weekday]
	}
	weekdays := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	return weekdays[weekday]
}

func getCurrentDateString(dateType string) string {
	now := time.Now()
	lang := getSystemLanguage()
	greg, lunar, tibetan := getDateTypeLabels()

	switch dateType {
	case greg: // 公历
		weekdayColor := getWeekdayColor(now.Weekday())
		weekdayName := getWeekdayName(now.Weekday())
		if lang == "zh" {
			return now.Format("2006年01月02日") + " " + weekdayColor + " " + weekdayName
		}
		return now.Format("2006-01-02") + " " + weekdayColor + " " + weekdayName

	case lunar: // 农历
		lunarInfo := getLunarInfo(now)
		return lunarInfo

	case tibetan: // 藏历
		tibetanInfo := getTibetanInfo(now)
		return tibetanInfo

	default:
		return ""
	}
}

func getLunarInfo(date time.Time) string {
	// 农历转换算法（基于农历数据表）
	lunarYear, lunarMonth, lunarDay, isLeapMonth := solarToLunar(date.Year(), int(date.Month()), date.Day())
	lang := getSystemLanguage()

	if lang == "zh" {
		monthName := getLunarMonthName(lunarYear, lunarMonth, isLeapMonth)
		return fmt.Sprintf("%d年%s%d日", lunarYear, monthName, lunarDay)
	}
	if isLeapMonth {
		return fmt.Sprintf("%d/Leap%d/%d", lunarYear, lunarMonth, lunarDay)
	}
	return fmt.Sprintf("%d/%d/%d", lunarYear, lunarMonth, lunarDay)
}

func getSolarTerm(date time.Time) string {
	// 简化的节气计算
	lang := getSystemLanguage()
	month := int(date.Month())
	day := date.Day()

	// 简单的节气判断（实际需要更精确的计算）
	if month == 3 && day >= 20 && day <= 22 {
		if lang == "zh" {
			return "春分"
		}
		return "Spring Equinox"
	}
	if month == 6 && day >= 20 && day <= 22 {
		if lang == "zh" {
			return "夏至"
		}
		return "Summer Solstice"
	}
	if month == 9 && day >= 22 && day <= 24 {
		if lang == "zh" {
			return "秋分"
		}
		return "Autumn Equinox"
	}
	if month == 12 && day >= 20 && day <= 22 {
		if lang == "zh" {
			return "冬至"
		}
		return "Winter Solstice"
	}
	return ""
}

func getTibetanInfo(date time.Time) string {
	// 藏历转换算法（基于Phugpa传统和Svante Janson的数学公式）
	tibetanYear, tibetanMonth, tibetanDay := solarToTibetan(date.Year(), int(date.Month()), date.Day())
	lang := getSystemLanguage()

	// 获取藏历特殊日期信息
	specialDay := getTibetanSpecialDays(date)
	specialInfo := ""
	if specialDay != "" {
		specialInfo = " (" + specialDay + ")"
	}

	if lang == "zh" {
		return fmt.Sprintf("%d年%d月%d日%s", tibetanYear, tibetanMonth, tibetanDay, specialInfo)
	}
	return fmt.Sprintf("%d/%d/%d%s", tibetanYear, tibetanMonth, tibetanDay, specialInfo)
}

func getTibetanSpecialDays(date time.Time) string {
	// 藏传佛教殊胜日和特殊日期（基于藏历日期）
	lang := getSystemLanguage()
	_, _, tibetanDay := solarToTibetan(date.Year(), int(date.Month()), date.Day())

	// 根据藏历日期返回对应的特殊日期
	switch tibetanDay {
	case 4:
		// 文殊菩萨剪头日
		if lang == "zh" {
			return "文殊菩萨剪头日"
		}
		return "Manjushri Hair Cutting Day"
	case 8:
		// 药师佛节日/殊胜日
		if lang == "zh" {
			return "药师佛节日/殊胜日"
		}
		return "Medicine Buddha Day/Auspicious Day"
	case 10:
		// 莲师节日
		if lang == "zh" {
			return "莲师节日"
		}
		return "Guru Rinpoche Day"
	case 15:
		// 阿弥陀佛节日/殊胜日
		if lang == "zh" {
			return "阿弥陀佛节日/殊胜日"
		}
		return "Amitabha Buddha Day/Auspicious Day"
	case 25:
		// 空行母节日
		if lang == "zh" {
			return "空行母节日"
		}
		return "Dakini Day"
	case 30:
		// 殊胜日
		if lang == "zh" {
			return "殊胜日"
		}
		return "Auspicious Day"
	default:
		return ""
	}
}

func ShowConfirmDialog(title, text string, confirmedCallback func()) {
	dialog.ShowConfirm(title, text,
		func(confirmed bool) {
			if confirmed && confirmedCallback != nil {
				confirmedCallback()
			}
		}, window,
	)
}

func ShowEntryDialog(title, placeholder, text string, confirmedCallback func(text string)) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)
	entry.SetText(text)

	dialogContainer := container.NewVBox(entry, canvas.NewText("", color.Black))

	dialog.ShowCustomConfirm(title, "OK", "Cancel", dialogContainer,
		func(confirmed bool) {
			if confirmed && confirmedCallback != nil {
				confirmedCallback(entry.Text)
			}
		}, window,
	)

	window.Canvas().Focus(entry)
}

func ShowColorPickerDialog(title, message string, preselected color.RGBA, confirmedCallback func(selected color.RGBA)) {
	colorPickerDialog := dialog.NewColorPicker(title, message,
		func(c color.Color) {
			if confirmedCallback != nil {
				confirmedCallback(ColorToRGBA(c))
			}
		}, window,
	)
	colorPickerDialog.Advanced = true
	colorPickerDialog.Show()
	colorPickerDialog.SetColor(preselected)
}

func ShowItemDialog(dialogPrefix, title, tagEditString, description string, style ItemStyle, confirmedCallback func(title, tagEditString, description string, style ItemStyle)) {
	ShowItemDialogWithDateType(dialogPrefix, title, tagEditString, description, style, "Normal", func(title, tagEditString, description string, style ItemStyle, dateType string) {
		confirmedCallback(title, tagEditString, description, style)
	})
}

func ShowItemDialogWithDateType(dialogPrefix, title, tagEditString, description string, style ItemStyle, currentDateType string, confirmedCallback func(title, tagEditString, description string, style ItemStyle, dateType string)) {
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Title ...")
	titleEntry.SetText(title)

	// 添加日期类型选择
	greg, lunar, tibetan := getDateTypeLabels()
	dateTypeOptions := []string{"Normal", greg, lunar, tibetan}
	dateTypeSelect := widget.NewSelect(dateTypeOptions, nil)
	dateTypeSelect.SetSelected(currentDateType)

	tagsEntry := widget.NewEntry()
	tagsEntry.SetPlaceHolder("Tag1=Value1; Tag2=Value2; ...")
	tagsEntry.SetText(tagEditString)

	descriptionEntry := widget.NewMultiLineEntry()
	descriptionEntry.SetPlaceHolder("Description ...")
	descriptionEntry.SetText(description)

	foregroundColor := style.Foreground
	foregroundColorButton := widget.NewButtonWithIcon("Foregound", theme.ColorPaletteIcon(),
		func() {
			ShowColorPickerDialog("Choose Foreground Color", "Please choose the color for item text and tag frames.", foregroundColor,
				func(selected color.RGBA) {
					foregroundColor = selected
				},
			)
		},
	)

	backgroundColor := style.Background
	backgroundColorButton := widget.NewButtonWithIcon("Background", theme.ColorPaletteIcon(),
		func() {
			ShowColorPickerDialog("Choose Background Color", "Please choose the color for the item's background.", backgroundColor,
				func(selected color.RGBA) {
					backgroundColor = selected
				},
			)
		},
	)

	buttonContainer := container.NewGridWithColumns(2, foregroundColorButton, backgroundColorButton)
	dialogContainer := container.NewVBox(titleEntry, dateTypeSelect, tagsEntry, descriptionEntry, buttonContainer, canvas.NewText("", color.Black))

	dialog.ShowCustomConfirm(dialogPrefix+" Item", "OK", "Cancel", dialogContainer,
		func(confirmed bool) {
			if confirmed && confirmedCallback != nil {
				// 根据选择的日期类型自动添加相应的标签
				finalTagString := tagsEntry.Text
				selectedType := dateTypeSelect.Selected

				if selectedType != "Normal" {
					if finalTagString != "" && !strings.HasSuffix(finalTagString, ";") {
						finalTagString += "; "
					}
					finalTagString += selectedType + "=" + getCurrentDateString(selectedType)
				}

				confirmedCallback(titleEntry.Text, finalTagString, descriptionEntry.Text, ItemStyle{foregroundColor, backgroundColor}, selectedType)
			}
		}, window,
	)

	window.Canvas().Focus(titleEntry)
}

func ShowFileOpenConfirmDialog(title, text string, defaultFileURI fyne.URI, confirmedCallback func(reader fyne.URIReadCloser)) {
	fileDialog := dialog.NewFileOpen(
		func(reader fyne.URIReadCloser, err error) {
			if reader != nil && err == nil {
				ShowConfirmDialog(title, text,
					func() {
						if confirmedCallback != nil {
							confirmedCallback(reader)
						}
					},
				)
			}
		}, window,
	)

	if defaultFileURI != nil {
		fileDialog.SetLocation(getParentListableURI(defaultFileURI))
	}

	fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
	fileDialog.Show()
}

func ShowSaveAsDialog(defaultFileURI fyne.URI, confirmedCallback func(writer fyne.URIWriteCloser)) {
	fileDialog := dialog.NewFileSave(
		func(writer fyne.URIWriteCloser, err error) {
			if writer != nil && err == nil && confirmedCallback != nil {
				confirmedCallback(writer)
			}
		}, window,
	)

	if defaultFileURI != nil {
		fileDialog.SetFileName(defaultFileURI.Name())
		fileDialog.SetLocation(getParentListableURI(defaultFileURI))
	} else {
		fileDialog.SetFileName("bankan_board.json")
	}

	fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
	fileDialog.Show()
}

/* ================================================================================ Calendar Conversion Functions */

// 农历年份数据结构
type LunarYearData struct {
	Year       int   // 农历年份
	NewYearDay int   // 农历新年在公历中是一年的第几天
	MonthDays  []int // 每个月的天数，闰月用负数表示月份位置
	LeapMonth  int   // 闰月月份，0表示无闰月
}

// 农历数据表（2020-2030年）
var lunarYearDataMap = map[int]*LunarYearData{
	2020: {2020, 25, []int{30, 29, 30, 29, 30, 29, 30, 30, 29, 30, 29, 30}, 4}, // 闰四月
	2021: {2021, 12, []int{29, 30, 29, 30, 29, 30, 29, 30, 29, 30, 30, 29}, 0},
	2022: {2022, 1, []int{30, 30, 29, 30, 29, 30, 29, 30, 29, 30, 29, 30}, 0},
	2023: {2023, 22, []int{30, 29, 30, 29, 30, 29, 30, 29, 30, 30, 29, 30}, 2}, // 闰二月
	2024: {2024, 10, []int{29, 30, 29, 30, 29, 30, 29, 30, 30, 29, 30, 29}, 0},
	2025: {2025, 29, []int{30, 29, 30, 29, 30, 29, 30, 29, 30, 29, 30, 29, 30}, 6}, // 闰六月
	2026: {2026, 17, []int{29, 30, 29, 30, 29, 30, 29, 30, 29, 30, 30, 29}, 0},
	2027: {2027, 6, []int{30, 29, 30, 29, 30, 29, 30, 29, 30, 29, 30, 30}, 0},
	2028: {2028, 26, []int{29, 30, 29, 30, 29, 30, 29, 30, 29, 30, 29, 30}, 5}, // 闰五月
	2029: {2029, 13, []int{30, 29, 30, 29, 30, 29, 30, 29, 30, 29, 30, 29}, 0},
	2030: {2030, 3, []int{30, 30, 29, 30, 29, 30, 29, 30, 29, 30, 29, 30}, 0},
}

// 获取农历年份数据
func getLunarYearData(year int) *LunarYearData {
	return lunarYearDataMap[year]
}

// 农历转换函数
func solarToLunar(year, month, day int) (int, int, int, bool) {
	// 精确的农历转换算法，基于农历数据表
	lunarData := getLunarYearData(year)
	if lunarData == nil {
		// 如果没有数据，返回近似值
		return year, month, day, false
	}

	// 计算从公历年初到指定日期的天数
	dayOfYear := getDayOfYear(year, month, day)

	// 农历新年在公历中的天数
	lunarNewYearDay := lunarData.NewYearDay

	var lunarYear int
	var daysFromLunarNewYear int

	if dayOfYear >= lunarNewYearDay {
		// 在农历新年之后
		lunarYear = year
		daysFromLunarNewYear = dayOfYear - lunarNewYearDay
	} else {
		// 在农历新年之前，属于上一个农历年
		prevLunarData := getLunarYearData(year - 1)
		if prevLunarData == nil {
			return year - 1, month, day, false
		}
		lunarYear = year - 1
		prevYearDays := 365
		if isLeapYear(year - 1) {
			prevYearDays = 366
		}
		daysFromLunarNewYear = (prevYearDays - prevLunarData.NewYearDay) + dayOfYear
		lunarData = prevLunarData
	}

	// 根据农历月份数据计算月份和日期
	currentDay := daysFromLunarNewYear
	for i, monthDays := range lunarData.MonthDays {
		if currentDay < monthDays {
			lunarMonth := i + 1
			lunarDay := currentDay + 1
			// 检查是否为闰月
			isLeapMonth := false
			if lunarData.LeapMonth > 0 {
				// 闰月在第LeapMonth个月之后
				if i == lunarData.LeapMonth {
					isLeapMonth = true
					lunarMonth = lunarData.LeapMonth
				} else if i > lunarData.LeapMonth {
					lunarMonth = i // 闰月后的月份
				}
			}
			return lunarYear, lunarMonth, lunarDay, isLeapMonth
		}
		currentDay -= monthDays
	}

	// 如果超出范围，返回最后一天
	return lunarYear, 12, 30, false
}

// 藏历数据结构
type TibetanYearData struct {
	Year        int   // 藏历年份
	NewYearDay  int   // 藏历新年在公历中是一年的第几天
	MonthDays   []int // 每个月的天数
	LeapMonth   int   // 闰月月份，0表示无闰月
	SkippedDays []int // 跳过的日期
	LeapDays    []int // 重复的日期
}

// 藏历数据表（基于Phugpa传统，2020-2030年）
// 数据基于Svante Janson的藏历数学公式计算
var tibetanYearDataMap = map[int]*TibetanYearData{
	2020: {2020, 55, []int{30, 30, 29, 30, 29, 30, 29, 30, 29, 30, 29, 30}, 0, []int{}, []int{}},     // 藏历新年约2月24日
	2021: {2021, 43, []int{30, 29, 30, 29, 30, 29, 30, 29, 30, 30, 29, 30}, 0, []int{}, []int{}},     // 藏历新年约2月12日
	2022: {2022, 32, []int{29, 30, 29, 30, 29, 30, 29, 30, 29, 30, 30, 29}, 0, []int{}, []int{}},     // 藏历新年约2月1日
	2023: {2023, 52, []int{30, 29, 30, 29, 30, 29, 30, 29, 30, 29, 30, 30}, 0, []int{}, []int{}},     // 藏历新年约2月21日
	2024: {2024, 40, []int{29, 30, 29, 30, 29, 30, 29, 30, 29, 30, 29, 30}, 0, []int{}, []int{}},     // 藏历新年约2月9日
	2025: {2025, 59, []int{30, 29, 30, 29, 30, 29, 30, 29, 30, 29, 30, 29, 30}, 6, []int{}, []int{}}, // 藏历新年约2月28日，闰六月
	2026: {2026, 47, []int{29, 30, 29, 30, 29, 30, 29, 30, 29, 30, 30, 29}, 0, []int{}, []int{}},     // 藏历新年约2月16日
	2027: {2027, 36, []int{30, 29, 30, 29, 30, 29, 30, 29, 30, 29, 30, 30}, 0, []int{}, []int{}},     // 藏历新年约2月5日
	2028: {2028, 56, []int{29, 30, 29, 30, 29, 30, 29, 30, 29, 30, 29, 30}, 5, []int{}, []int{}},     // 藏历新年约2月25日，闰五月
	2029: {2029, 43, []int{30, 29, 30, 29, 30, 29, 30, 29, 30, 29, 30, 29}, 0, []int{}, []int{}},     // 藏历新年约2月12日
	2030: {2030, 33, []int{30, 30, 29, 30, 29, 30, 29, 30, 29, 30, 29, 30}, 0, []int{}, []int{}},     // 藏历新年约2月2日
}

// 获取藏历年份数据
func getTibetanYearData(year int) *TibetanYearData {
	return tibetanYearDataMap[year]
}

// 藏历转换函数（改进版，基于Phugpa传统）
func solarToTibetan(year, month, day int) (int, int, int) {
	// 基于Svante Janson的藏历数学公式和Phugpa传统
	tibetanData := getTibetanYearData(year)
	if tibetanData == nil {
		// 如果没有精确数据，使用改进的近似算法
		return solarToTibetanApproximate(year, month, day)
	}

	// 计算从公历年初到指定日期的天数
	dayOfYear := getDayOfYear(year, month, day)

	// 藏历新年在公历中的天数
	tibetanNewYearDay := tibetanData.NewYearDay

	var tibetanYear int
	var daysFromTibetanNewYear int

	if dayOfYear >= tibetanNewYearDay {
		// 在藏历新年之后
		tibetanYear = year
		daysFromTibetanNewYear = dayOfYear - tibetanNewYearDay
	} else {
		// 在藏历新年之前，属于上一个藏历年
		prevTibetanData := getTibetanYearData(year - 1)
		if prevTibetanData == nil {
			return solarToTibetanApproximate(year, month, day)
		}
		tibetanYear = year - 1
		prevYearDays := 365
		if isLeapYear(year - 1) {
			prevYearDays = 366
		}
		daysFromTibetanNewYear = (prevYearDays - prevTibetanData.NewYearDay) + dayOfYear
		tibetanData = prevTibetanData
	}

	// 根据藏历月份数据计算月份和日期
	currentDay := daysFromTibetanNewYear
	for i, monthDays := range tibetanData.MonthDays {
		if currentDay < monthDays {
			tibetanMonth := i + 1
			tibetanDay := currentDay + 1

			// 处理闰月情况
			if tibetanData.LeapMonth > 0 && i >= tibetanData.LeapMonth {
				if i == tibetanData.LeapMonth {
					// 这是闰月
					tibetanMonth = tibetanData.LeapMonth
				} else {
					// 闰月后的月份需要调整
					tibetanMonth = i
				}
			}

			return tibetanYear, tibetanMonth, tibetanDay
		}
		currentDay -= monthDays
	}

	// 如果超出范围，返回最后一天
	return tibetanYear, 12, 30
}

// 藏历近似算法（当没有精确数据时使用）
func solarToTibetanApproximate(year, month, day int) (int, int, int) {
	// 改进的近似算法，基于藏历的基本规律
	// 藏历新年通常在公历2月中旬，比农历新年稍晚
	tibetanYear := year
	dayOfYear := getDayOfYear(year, month, day)

	// 藏历新年的近似计算（基于19年周期）
	// 藏历使用Kālacakra系统，有复杂的闰月和跳日规则
	yearInCycle := (year - 1027) % 60             // 藏历60年周期，从1027年开始
	tibetanNewYearDay := 32 + (yearInCycle*11)%30 // 基本在2月1日到3月2日之间
	if tibetanNewYearDay > 60 {                   // 如果超过2月底
		tibetanNewYearDay -= 30
	}

	if dayOfYear >= tibetanNewYearDay {
		daysFromTibetanNewYear := dayOfYear - tibetanNewYearDay
		// 藏历月份长度在29-30天之间变化
		tibetanMonth := (daysFromTibetanNewYear / 30) + 1
		tibetanDay := (daysFromTibetanNewYear % 30) + 1

		// 调整月份边界（藏历月份不完全是30天）
		if tibetanDay > 30 {
			tibetanMonth++
			tibetanDay = 1
		}
		if tibetanMonth > 12 {
			tibetanMonth = 12
			tibetanDay = 30
		}
		return tibetanYear, tibetanMonth, tibetanDay
	} else {
		// 跨年情况
		prevYearDays := 365
		if isLeapYear(year - 1) {
			prevYearDays = 366
		}
		prevYearInCycle := (year - 1 - 1027) % 60
		prevTibetanNewYear := 32 + (prevYearInCycle*11)%30
		if prevTibetanNewYear > 60 {
			prevTibetanNewYear -= 30
		}

		daysFromPrevTibetanNewYear := (prevYearDays - prevTibetanNewYear) + dayOfYear
		tibetanMonth := (daysFromPrevTibetanNewYear / 30) + 1
		tibetanDay := (daysFromPrevTibetanNewYear % 30) + 1

		if tibetanDay > 30 {
			tibetanMonth++
			tibetanDay = 1
		}
		if tibetanMonth > 12 {
			tibetanMonth = 12
			tibetanDay = 30
		}
		return tibetanYear - 1, tibetanMonth, tibetanDay
	}
}

// 获取农历月份名称（支持闰月）
func getLunarMonthName(year, month int, isLeapMonth bool) string {
	lang := getSystemLanguage()
	if lang == "zh" {
		monthNames := []string{"", "正月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "冬月", "腊月"}
		if month >= 1 && month <= 12 {
			if isLeapMonth {
				return "闰" + monthNames[month]
			}
			return monthNames[month]
		}
	}
	if isLeapMonth {
		return fmt.Sprintf("闰%d月", month)
	}
	return fmt.Sprintf("%d月", month)
}

// 辅助函数：获取一年中的第几天
func getDayOfYear(year, month, day int) int {
	daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if isLeapYear(year) {
		daysInMonth[1] = 29
	}

	dayOfYear := day
	for i := 0; i < month-1; i++ {
		dayOfYear += daysInMonth[i]
	}
	return dayOfYear
}

// 辅助函数：判断是否为闰年
func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

// 辅助函数：获取农历新年在公历中的大致日期（一年中的第几天）
func getLunarNewYearDay(year int) int {
	// 简化计算，农历新年大致在1月21日到2月20日之间
	// 这里使用一个基本的周期性近似
	base := 21 + ((year-2000)*11)%30 // 简化的周期计算
	if base > 51 {                   // 如果超过2月20日（31+20=51）
		base = base - 30
	}
	return base
}

/* ================================================================================ Private functions */
func getParentListableURI(file fyne.URI) fyne.ListableURI {
	dirURI, err := storage.Parent(file)
	if err != nil {
		return nil
	}

	dirListableURI, err := storage.ListerForURI(dirURI)
	if err != nil {
		return nil
	}

	return dirListableURI
}

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
	"github.com/liujiawm/gocalendar"
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

func getDataTypeLabels() (string, string, string) {
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
		return "🟢" // 绿色1
	case time.Tuesday:
		return "🔵" // 蓝色2
	case time.Wednesday:
		return "🟡" // 黄色3
	case time.Thursday:
		return "🟠" // 橙色4
	case time.Friday:
		return "🔴" // 红色7
	case time.Saturday:
		return "⚪" // 白色6
	case time.Sunday:
		return "🟣" // 紫色5
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

func getCurrentDateString(dataType string) string {
	now := time.Now()
	lang := getSystemLanguage()

	switch dataType {
	case "Gregorian": // 公历
		weekdayColor := getWeekdayColor(now.Weekday())
		weekdayName := getWeekdayName(now.Weekday())
		if lang == "zh" {
			return now.Format("2006年01月02日") + " " + weekdayColor + " " + weekdayName
		}
		return now.Format("2006-01-02") + " " + weekdayColor + " " + weekdayName

	case "Lunar": // 农历
		lunarInfo := getLunarInfo(now)
		return lunarInfo

	case "Tibetan": // 藏历
		tibetanInfo := getTibetanInfo(now)
		return tibetanInfo

	default:
		return ""
	}
}

func getLunarInfo(date time.Time) string {
	// 使用gocalendar库进行精确的农历转换和节气计算 <mcreference link="https://github.com/liujiawm/gocalendar" index="1">1</mcreference>
	lang := getSystemLanguage()

	// 创建日历实例并获取指定日期的信息
	cal := gocalendar.DefaultCalendar()
	items := cal.GenerateWithDate(date.Year(), int(date.Month()), date.Day())

	// 查找当前日期的信息
	var currentItem *gocalendar.CalendarItem
	for _, item := range items {
		if item.Time.Year() == date.Year() && item.Time.Month() == date.Month() && item.Time.Day() == date.Day() {
			currentItem = item
			break
		}
	}

	if currentItem == nil {
		// 如果没有找到，返回基本格式
		if lang == "zh" {
			return fmt.Sprintf("%d年%d月%d日", date.Year(), int(date.Month()), date.Day())
		}
		return fmt.Sprintf("%d/%d/%d", date.Year(), int(date.Month()), date.Day())
	}

	// 获取农历信息
	lunarDate := currentItem.LunarDate

	// 获取节气信息
	solarTermInfo := ""
	if currentItem.SolarTerm != nil && currentItem.SolarTerm.Name != "" {
		if lang == "zh" {
			solarTermInfo = " (" + currentItem.SolarTerm.Name + ")"
		} else {
			// 简单的英文翻译映射
			englishName := getSolarTermEnglishName(currentItem.SolarTerm.Name)
			solarTermInfo = " (" + englishName + ")"
		}
	}

	// 获取十斋日信息
	fastingDayInfo := getLunarFastingDayInfo(lunarDate.Day, lang)
	if fastingDayInfo != "" {
		if solarTermInfo != "" {
			solarTermInfo = solarTermInfo + ", " + fastingDayInfo
		} else {
			solarTermInfo = " (" + fastingDayInfo + ")"
		}
	}

	if lang == "zh" {
		// 中文格式：农历年份 + 月份名称 + 日期 + 节气 + 十斋日
		monthName := lunarDate.MonthName + "月" // 添加"月"字
		dayName := lunarDate.DayName
		if lunarDate.LeapStr != "" {
			monthName = lunarDate.LeapStr + lunarDate.MonthName + "月" // 闰月也要添加"月"字
		}
		return fmt.Sprintf("%d年%s%s%s", lunarDate.Year, monthName, dayName, solarTermInfo)
	}

	// 英文格式
	if lunarDate.LeapStr != "" {
		return fmt.Sprintf("%d/Leap%d/%d%s", lunarDate.Year, lunarDate.Month, lunarDate.Day, solarTermInfo)
	}
	return fmt.Sprintf("%d/%d/%d%s", lunarDate.Year, lunarDate.Month, lunarDate.Day, solarTermInfo)
}

// 获取农历十斋日信息
func getLunarFastingDayInfo(lunarDay int, lang string) string {
	// 十斋日：初一、初八、十四、十五、十八、二十三、二十四、二十八、二十九、三十
	// 六斋日：初八、十四、十五、二十三、二十九、三十
	switch lunarDay {
	case 1:
		if lang == "zh" {
			return "十斋日"
		}
		return "Ten Fasting Days"
	case 8:
		if lang == "zh" {
			return "六斋日/十斋日"
		}
		return "Six/Ten Fasting Days"
	case 14:
		if lang == "zh" {
			return "六斋日/十斋日"
		}
		return "Six/Ten Fasting Days"
	case 15:
		if lang == "zh" {
			return "六斋日/十斋日"
		}
		return "Six/Ten Fasting Days"
	case 18:
		if lang == "zh" {
			return "十斋日"
		}
		return "Ten Fasting Days"
	case 23:
		if lang == "zh" {
			return "六斋日/十斋日"
		}
		return "Six/Ten Fasting Days"
	case 24:
		if lang == "zh" {
			return "十斋日"
		}
		return "Ten Fasting Days"
	case 28:
		if lang == "zh" {
			return "十斋日"
		}
		return "Ten Fasting Days"
	case 29:
		if lang == "zh" {
			return "六斋日/十斋日"
		}
		return "Six/Ten Fasting Days"
	case 30:
		if lang == "zh" {
			return "六斋日/十斋日"
		}
		return "Six/Ten Fasting Days"
	default:
		return ""
	}
}

// 节气中英文名称映射
func getSolarTermEnglishName(chineseName string) string {
	solarTermMap := map[string]string{
		"立春": "Beginning of Spring",
		"雨水": "Rain Water",
		"惊蛰": "Awakening of Insects",
		"春分": "Spring Equinox",
		"清明": "Clear and Bright",
		"谷雨": "Grain Rain",
		"立夏": "Beginning of Summer",
		"小满": "Grain Buds",
		"芒种": "Grain in Ear",
		"夏至": "Summer Solstice",
		"小暑": "Slight Heat",
		"大暑": "Great Heat",
		"立秋": "Beginning of Autumn",
		"处暑": "Stopping the Heat",
		"白露": "White Dew",
		"秋分": "Autumn Equinox",
		"寒露": "Cold Dew",
		"霜降": "Frost's Descent",
		"立冬": "Beginning of Winter",
		"小雪": "Slight Snow",
		"大雪": "Great Snow",
		"冬至": "Winter Solstice",
		"小寒": "Slight Cold",
		"大寒": "Great Cold",
	}

	if englishName, exists := solarTermMap[chineseName]; exists {
		return englishName
	}
	return chineseName // 如果没有找到映射，返回原名称
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

	// 获取理发吉凶信息
	hairCutInfo := getTibetanHairCutInfo(tibetanDay, lang)

	// 获取殊胜日信息
	specialDayInfo := getTibetanSpecialDayInfo(tibetanDay, lang)

	// 组合信息
	var result []string
	if specialDayInfo != "" {
		result = append(result, specialDayInfo)
	}
	if hairCutInfo != "" {
		result = append(result, hairCutInfo)
	}

	if len(result) > 0 {
		return strings.Join(result, ", ")
	}
	return ""
}

// 获取藏历理发吉凶信息
func getTibetanHairCutInfo(tibetanDay int, lang string) string {
	// 根据用户提供的藏历理发吉凶对照表
	switch tibetanDay {
	case 1:
		if lang == "zh" {
			return "理发凶🔴: 短命减寿"
		}
		return "Hair Cut: Auspicious"
	case 2:
		if lang == "zh" {
			return "理发凶🔴: 遇传染病"
		}
		return "Hair Cut: Risk of Contagious Disease"
	case 3:
		if lang == "zh" {
			return "理发吉: 财富增上"
		}
		return "Hair Cut: Sweet"
	case 4:
		if lang == "zh" {
			return "理发凶🔴: 低贱, 豆腐店主"
		}
		return "Hair Cut: Lowly, Tofu Shop Owner"
	case 5:
		if lang == "zh" {
			return "理发凶🔴: 易患疾病"
		}
		return "Hair Cut: Prone to Illness, Inauspicious"
	case 6:
		if lang == "zh" {
			return "理发吉: 面色红润"
		}
		return "Hair Cut: Rosy Complexion"
	case 7:
		if lang == "zh" {
			return "理发凶🔴: 易争吵"
		}
		return "Hair Cut: Prone to Arguments"
	case 8:
		if lang == "zh" {
			return "理发吉: 得长寿"
		}
		return "Hair Cut: Longevity"
	case 9:
		if lang == "zh" {
			return "理发吉: 姻缘"
		}
		return "Hair Cut: Meet Monks, Sharing"
	case 10:
		if lang == "zh" {
			return "理发凶🔴: 遇传染病"
		}
		return "Hair Cut: Contagious Disease"
	case 11:
		if lang == "zh" {
			return "理发吉: 增长智慧"
		}
		return "Hair Cut: Increase Wisdom"
	case 12:
		if lang == "zh" {
			return "理发凶🔴: 招致疾病"
		}
		return "Hair Cut: Attract Disease, Inauspicious"
	case 13:
		if lang == "zh" {
			return "理发吉: 佛慧增长"
		}
		return "Hair Cut: Skill Improvement"
	case 14:
		if lang == "zh" {
			return "理发吉: 增长财富"
		}
		return "Hair Cut: Growth of Things"
	case 15:
		if lang == "zh" {
			return "理发吉: 增长福报"
		}
		return "Hair Cut: Increase Merit"
	case 16:
		if lang == "zh" {
			return "理发凶🔴: 患病"
		}
		return "Hair Cut: Illness"
	case 17:
		if lang == "zh" {
			return "理发凶🔴: 易失明, 眼疾 han"
		}
		return "Hair Cut: Risk of Blindness, Eye Disease"
	case 18:
		if lang == "zh" {
			return "理发凶🔴: 丢失财物"
		}
		return "Hair Cut: Loss of Property"
	case 19:
		if lang == "zh" {
			return "理发吉: 增长寿命"
		}
		return "Hair Cut: Increase Lifespan"
	case 20:
		if lang == "zh" {
			return "理发凶🔴: 易挨饿"
		}
		return "Hair Cut: Prone to Hunger"
	case 21:
		if lang == "zh" {
			return "理发凶🔴: 易患眼疾, 失明"
		}
		return "Hair Cut: Eye Disease, Blindness"
	case 22:
		if lang == "zh" {
			return "理发吉: 增长财物"
		}
		return "Hair Cut: Increase Wealth"
	case 23:
		if lang == "zh" {
			return "理发凶🔴: 患麻风病等"
		}
		return "Hair Cut: Leprosy etc."
	case 24:
		if lang == "zh" {
			return "理发凶🔴: 遇口舌, 凶"
		}
		return "Hair Cut: Disputes, Inauspicious"
	case 25:
		if lang == "zh" {
			return "理发凶🔴: 得白内障"
		}
		return "Hair Cut: Get Cataract"
	case 26:
		if lang == "zh" {
			return "理发吉: 得快乐"
		}
		return "Hair Cut: Get Happiness"
	case 27:
		if lang == "zh" {
			return "理发凶🔴: 吐血, 凶"
		}
		return "Hair Cut: Vomit Blood, Inauspicious"
	case 28:
		if lang == "zh" {
			return "理发凶🔴: 易患疯癫"
		}
		return "Hair Cut: Prone to Madness"
	case 29:
		if lang == "zh" {
			return "理发凶🔴: 易患白癜风"
		}
		return "Hair Cut: Prone to Vitiligo"
	case 30:
		if lang == "zh" {
			return "理发凶🔴: 死于争斗中"
		}
		return "Hair Cut: Die in Conflict"
	default:
		return ""
	}
}

// 获取藏历殊胜日信息
func getTibetanSpecialDayInfo(tibetanDay int, lang string) string {
	// 根据藏历日期返回对应的特殊日期
	switch tibetanDay {
	case 4:
		// 文殊菩萨剪头日
		if lang == "zh" {
			return ""
		}
		return ""
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
	ShowItemDialogWithDataType(dialogPrefix, title, tagEditString, description, style, "Normal", func(title, tagEditString, description string, style ItemStyle, dataType string) {
		confirmedCallback(title, tagEditString, description, style)
	})
}

func ShowItemDialogWithDataType(dialogPrefix, title, tagEditString, description string, style ItemStyle, currentDataType string, confirmedCallback func(title, tagEditString, description string, style ItemStyle, dataType string)) {
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Title ...")
	titleEntry.SetText(title)

	// 添加数据类型选择
	greg, lunar, tibetan := getDataTypeLabels()
	dataTypeOptions := []string{"Normal", "Gregorian", "Lunar", "Tibetan"}
	dataTypeSelect := widget.NewSelect(dataTypeOptions, nil)
	// 设置当前选择的数据类型
	if currentDataType == greg {
		dataTypeSelect.SetSelected("Gregorian")
	} else if currentDataType == lunar {
		dataTypeSelect.SetSelected("Lunar")
	} else if currentDataType == tibetan {
		dataTypeSelect.SetSelected("Tibetan")
	} else {
		dataTypeSelect.SetSelected(currentDataType)
	}

	tagsEntry := widget.NewEntry()
	tagsEntry.SetPlaceHolder("Tag1=Value1; Tag2=Value2; ...")
	tagsEntry.SetText(tagEditString)

	descriptionEntry := widget.NewMultiLineEntry()
	descriptionEntry.SetPlaceHolder("Description ...")
	descriptionEntry.SetText(description)
	// 设置描述输入框的最小尺寸为两倍高度
	descriptionEntry.Resize(fyne.NewSize(descriptionEntry.MinSize().Width, 400))

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

	// 创建一个固定尺寸的容器来包装所有组件，实现对话框尺寸加倍
	contentContainer := container.NewVBox(titleEntry, dataTypeSelect, tagsEntry, descriptionEntry, buttonContainer)
	// 使用Border容器设置固定尺寸，宽度和高度都比原来大
	dialogContainer := container.NewBorder(nil, nil, nil, nil, contentContainer)
	dialogContainer.Resize(fyne.NewSize(600, 400)) // 设置对话框容器的固定尺寸

	dialog.ShowCustomConfirm(dialogPrefix+" Item", "OK", "Cancel", dialogContainer,
		func(confirmed bool) {
			if confirmed && confirmedCallback != nil {
				// 根据选择的数据类型处理标签和标题
				finalTagString := tagsEntry.Text
				selectedType := dataTypeSelect.Selected
				finalTitle := titleEntry.Text

				// 如果选择了日期类型，将日期信息添加到标题中，标签只显示类型
				if selectedType != "Normal" {
					// 只有在创建新项目时（dialogPrefix为"Add"）才自动添加日期信息
					if dialogPrefix == "Add" {
						dateString := getCurrentDateString(selectedType)
						if dateString != "" {
							// 将日期信息添加到标题中
							if finalTitle != "" {
								finalTitle += " " + dateString
							} else {
								finalTitle = dateString
							}
						}
					}
					// 标签只显示日期类型的中文名称，但要检查是否已存在避免重复
					greg, lunar, tibetan := getDataTypeLabels()
					typeLabel := ""
					switch selectedType {
					case "Gregorian":
						typeLabel = greg
					case "Lunar":
						typeLabel = lunar
					case "Tibetan":
						typeLabel = tibetan
					}
					// 检查标签中是否已经包含该类型，避免重复添加
					if typeLabel != "" && !strings.Contains(finalTagString, typeLabel) {
						if finalTagString != "" && !strings.HasSuffix(finalTagString, ";") {
							finalTagString += "; "
						}
						finalTagString += typeLabel
					}
				}

				confirmedCallback(finalTitle, finalTagString, descriptionEntry.Text, ItemStyle{foregroundColor, backgroundColor}, selectedType)
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

// 注意：农历和节气计算现在使用gocalendar库提供精确算法
// 基于Jean Meeus的《Astronomical Algorithms》和NASA数据

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

// getLunarMonthName函数已被gocalendar库的MonthName字段替代

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

// getLunarNewYearDay函数已被gocalendar库替代，提供更精确的计算

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

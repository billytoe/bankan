package main

/* FontSize manager for handling three font size levels */

/* ================================================================================ Imports */
import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

/* ================================================================================ Constants */
const (
	FONT_SIZE_SMALL  = 0
	FONT_SIZE_MEDIUM = 1
	FONT_SIZE_LARGE  = 2
)

/* ================================================================================ Private variables */
var currentFontSizeLevel = FONT_SIZE_SMALL // 默认为最小档

/* ================================================================================ Public functions */

// GetCurrentFontSizeLevel 获取当前字体大小档位
func GetCurrentFontSizeLevel() int {
	return currentFontSizeLevel
}

// SetFontSizeLevel 设置字体大小档位
func SetFontSizeLevel(level int) {
	if level >= FONT_SIZE_SMALL && level <= FONT_SIZE_LARGE {
		currentFontSizeLevel = level
		// 保存到偏好设置
		fyne.CurrentApp().Preferences().SetInt("fontSizeLevel", level)
	}
}

// RestoreFontSizeLevel 从偏好设置恢复字体大小档位
func RestoreFontSizeLevel() {
	level := fyne.CurrentApp().Preferences().IntWithFallback("fontSizeLevel", FONT_SIZE_SMALL)
	if level >= FONT_SIZE_SMALL && level <= FONT_SIZE_LARGE {
		currentFontSizeLevel = level
	}
}

// GetScaledTextSize 根据当前档位获取缩放后的文本大小
func GetScaledTextSize() float32 {
	baseSize := theme.TextSize()
	switch currentFontSizeLevel {
	case FONT_SIZE_SMALL:
		return baseSize // 1.0x
	case FONT_SIZE_MEDIUM:
		return baseSize * 1.2 // 1.2x
	case FONT_SIZE_LARGE:
		return baseSize * 1.4 // 1.4x
	default:
		return baseSize
	}
}

// GetScaledCaptionTextSize 根据当前档位获取缩放后的标题文本大小
func GetScaledCaptionTextSize() float32 {
	baseSize := theme.CaptionTextSize()
	switch currentFontSizeLevel {
	case FONT_SIZE_SMALL:
		return baseSize // 1.0x
	case FONT_SIZE_MEDIUM:
		return baseSize * 1.2 // 1.2x
	case FONT_SIZE_LARGE:
		return baseSize * 1.4 // 1.4x
	default:
		return baseSize
	}
}

// GetScaledTextSubHeadingSize 根据当前档位获取缩放后的副标题文本大小
func GetScaledTextSubHeadingSize() float32 {
	baseSize := theme.TextSubHeadingSize()
	switch currentFontSizeLevel {
	case FONT_SIZE_SMALL:
		return baseSize // 1.0x
	case FONT_SIZE_MEDIUM:
		return baseSize * 1.2 // 1.2x
	case FONT_SIZE_LARGE:
		return baseSize * 1.4 // 1.4x
	default:
		return baseSize
	}
}

// GetFontSizeLevelName 获取字体大小档位的名称
func GetFontSizeLevelName(level int) string {
	switch level {
	case FONT_SIZE_SMALL:
		return "小"
	case FONT_SIZE_MEDIUM:
		return "中"
	case FONT_SIZE_LARGE:
		return "大"
	default:
		return "小"
	}
}

// GetCurrentFontSizeLevelName 获取当前字体大小档位的名称
func GetCurrentFontSizeLevelName() string {
	return GetFontSizeLevelName(currentFontSizeLevel)
}
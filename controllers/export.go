package controllers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"DLM_backend/services"
	"DLM_backend/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// ExportInspection 导出点检记录为Excel文件
func ExportInspection(c *gin.Context) {
	// 检查用户角色是否为管理员
	claims, exists := c.Get("claims")
	if !exists {
		utils.UnauthorizedResponse(c, "token claims not found")
		return
	}
	role := claims.(jwt.MapClaims)["role"].(string)
	if role != "admin" {
		utils.UnauthorizedResponse(c, "only admin can export records")
		return
	}

	// 创建Excel文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	// 创建工作表
	sheetName := "点检记录"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		utils.ServerErrorResponse(c, "failed to create sheet")
		return
	}
	f.SetActiveSheet(index)

	// 设置表头
	headers := []string{
		"ID", "单位", "仓号", "挡粮门位置", "保管责任人",
		"检查时间", "变形和裂痕情况", "变形和裂痕说明",
		"闭合情况", "闭合说明", "栓销状况", "栓销说明",
		"主体墙状况", "主体墙说明", "仓门地基状况", "地基说明",
		"安全绳装置", "安全绳说明", "补充说明", "责任人签名",
		"联系电话", "图片列表",
	}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheetName, cell, header)
	}

	// 获取所有记录
	records, err := services.GetInspectionRecords()
	if err != nil {
		utils.ServerErrorResponse(c, "failed to get records")
		return
	}

	// 写入数据
	for i, record := range records {
		row := i + 2 // 第一行是表头
		f.SetCellValue(sheetName, "A"+strconv.Itoa(row), record.ID)
		f.SetCellValue(sheetName, "B"+strconv.Itoa(row), record.Unit)
		f.SetCellValue(sheetName, "C"+strconv.Itoa(row), record.WarehouseNumber)
		f.SetCellValue(sheetName, "D"+strconv.Itoa(row), record.GrainDoorPosition)
		f.SetCellValue(sheetName, "E"+strconv.Itoa(row), record.Caretaker)
		f.SetCellValue(sheetName, "F"+strconv.Itoa(row), record.InspectionTime.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, "G"+strconv.Itoa(row), record.DeformationCrack)
		f.SetCellValue(sheetName, "H"+strconv.Itoa(row), record.DeformationCrackDescription)
		f.SetCellValue(sheetName, "I"+strconv.Itoa(row), record.ClosureStatus)
		f.SetCellValue(sheetName, "J"+strconv.Itoa(row), record.ClosureDescription)

		// 处理栓销状态
		// 'normal': '正常',
		// 'loose': '松动',
		// 'deformed': '变形',
		// 'missing': '缺失',
		var pinStatus []string
		if err := json.Unmarshal([]byte(record.PinStatus), &pinStatus); err == nil {
			pinStatusStr := ""
			for _, status := range pinStatus {
				switch status {
				case "normal":
					pinStatusStr += "正常"
				case "loose":
					pinStatusStr += "松动"
				case "deformed":
					pinStatusStr += "变形"
				case "missing":
					pinStatusStr += "缺失"
				}
				pinStatusStr += ","
			}
			// 去掉最后一个逗号
			pinStatusStr = pinStatusStr[:len(pinStatusStr)-1]
			f.SetCellValue(sheetName, "K"+strconv.Itoa(row), pinStatusStr)
		} else {
			// 如果解析失败，使用原始字符串
			f.SetCellValue(sheetName, "K"+strconv.Itoa(row), string(record.PinStatus))
		}

		f.SetCellValue(sheetName, "L"+strconv.Itoa(row), record.PinDescription)

		// 处理主体墙状态
		// 'normal': '正常',
		// 'damaged': '破损',
		// 'cracked': '有裂缝',

		var mainWallStatus []string
		if err := json.Unmarshal([]byte(record.MainWallStatus), &mainWallStatus); err == nil {
			mainWallStatusStr := ""
			for _, status := range mainWallStatus {
				switch status {
				case "normal":
					mainWallStatusStr += "正常"
				case "damaged":
					mainWallStatusStr += "破损"
				case "cracked":
					mainWallStatusStr += "有裂缝"
				}
				mainWallStatusStr += ","
			}
			// 去掉最后一个逗号
			mainWallStatusStr = mainWallStatusStr[:len(mainWallStatusStr)-1]
			f.SetCellValue(sheetName, "M"+strconv.Itoa(row), mainWallStatusStr)
		} else {
			// 如果解析失败，使用原始字符串
			f.SetCellValue(sheetName, "M"+strconv.Itoa(row), string(record.MainWallStatus))
		}

		f.SetCellValue(sheetName, "N"+strconv.Itoa(row), record.MainWallDescription)

		// 处理仓门地基状态
		// 'normal': '正常',
		// 'frozen': '冻胀',
		// 'sinking': '下沉',
		// 'collapsed': '塌陷'
		var foundationStatus []string
		if err := json.Unmarshal([]byte(record.WarehouseFoundation), &foundationStatus); err == nil {
			foundationStatusStr := ""
			for _, status := range foundationStatus {
				switch status {
				case "normal":
					foundationStatusStr += "正常"
				case "frozen":
					foundationStatusStr += "冻胀"
				case "sinking":
					foundationStatusStr += "下沉"
				case "collapsed":
					foundationStatusStr += "塌陷"
				}
				foundationStatusStr += ","
			}
			// 去掉最后一个逗号
			foundationStatusStr = foundationStatusStr[:len(foundationStatusStr)-1]
			f.SetCellValue(sheetName, "O"+strconv.Itoa(row), foundationStatusStr)
		} else {
			// 如果解析失败，使用原始字符串
			f.SetCellValue(sheetName, "O"+strconv.Itoa(row), string(record.WarehouseFoundation))
		}

		f.SetCellValue(sheetName, "P"+strconv.Itoa(row), record.WarehouseFoundationDescription)
		f.SetCellValue(sheetName, "Q"+strconv.Itoa(row), record.SafetyRopeInstalled)
		f.SetCellValue(sheetName, "R"+strconv.Itoa(row), record.SafetyRopeDescription)
		f.SetCellValue(sheetName, "S"+strconv.Itoa(row), record.Remarks)
		f.SetCellValue(sheetName, "T"+strconv.Itoa(row), record.Signature)
		f.SetCellValue(sheetName, "U"+strconv.Itoa(row), record.ContactNumber)

		// 处理图片路径
		var imagePaths []string
		if err := json.Unmarshal([]byte(record.Images), &imagePaths); err == nil {
			imagePathsStr := ""
			// 动态判断协议 (HTTP/HTTPS)
			scheme := "http"
			// 检查 X-Forwarded-Proto 头（当使用反向代理如Nginx时）
			if c.GetHeader("X-Forwarded-Proto") == "https" {
				scheme = "https"
			} else if c.Request.TLS != nil {
				// 如果是直接的TLS连接
				scheme = "https"
			}

			// 拼接前缀，创建完整URL
			for _, path := range imagePaths {
				if path == "" {
					continue
				}
				imagePathsStr += fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, path) + ","
			}
			// 去掉最后一个逗号
			if len(imagePathsStr) > 0 {
				imagePathsStr = imagePathsStr[:len(imagePathsStr)-1]
			}

			f.SetCellValue(sheetName, "V"+strconv.Itoa(row), imagePathsStr)
		} else {
			// 如果解析失败，使用原始字符串
			f.SetCellValue(sheetName, "V"+strconv.Itoa(row), string(record.Images))
		}
	}

	// 调整列宽
	for i := 0; i < len(headers); i++ {
		col := string(rune('A' + i))
		f.SetColWidth(sheetName, col, col, 15)
	}

	// 创建导出目录
	exportDir := "exports"
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		utils.ServerErrorResponse(c, "failed to create export directory")
		return
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("inspection_records_%s.xlsx", timestamp)
	filepath := filepath.Join(exportDir, filename)

	// 保存文件
	if err := f.SaveAs(filepath); err != nil {
		utils.ServerErrorResponse(c, "failed to save excel file")
		return
	}

	// 返回文件下载地址
	downloadURL := fmt.Sprintf("/exports/%s", filename)
	utils.SuccessResponse(c, gin.H{
		"url":      downloadURL,
		"filename": filename,
	})
}

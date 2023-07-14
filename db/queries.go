package db

import (
	"time"

	"gorm.io/gorm"
)

func DbSendVoiceMail(DBConn *gorm.DB, message Invoice) (bool, uint) {

	result := DBConn.Create(&message)
	handleError(result.Error, "error has happened while processing dbmSendVoiceMail query")
	if result.RowsAffected == 0 {
		return false, 0
	}

	return true, message.ID
}

func DbMarkAsSeen(DBConn *gorm.DB, messageid uint) (bool, uint) {
	var invoice Invoice

	result := DBConn.First(&invoice, messageid)
	handleError(result.Error, "error has happened while processing dbMarkAsSeen query")
	if result.RowsAffected == 0 {
		return false, 0
	}

	return true, messageid
}

func DbReadMessageTime(DBConn *gorm.DB, messageid uint) (bool, bool, time.Time) {
	var invoice Invoice

	result := DBConn.First(&invoice, messageid)
	handleError(result.Error, "error has happened while processing dbReadMessageTime query")

	Read := false
	if result.RowsAffected == 0 {
		return false, false, time.Now()
	}
	if invoice.CreatedAt != invoice.UpdatedAt {
		Read = true
	}

	return true, Read, invoice.UpdatedAt
}

func DbReadUnSeenReceived(DBConn *gorm.DB, toUserId uint, skip int, limit int) []Invoice {
	var invoices []Invoice

	result := DBConn.Where("ToID = ? and Seen = false", toUserId).Offset(skip).Limit(limit).Find(&invoices)
	handleError(result.Error, "error has happened while processing dbReadUnSeenReceived query")

	return invoices
}

func DbReadAllReceived(DBConn *gorm.DB, toUserId uint, skip int, limit int) []Invoice {
	var invoices []Invoice

	result := DBConn.Where("ToID = ?", toUserId).Offset(skip).Limit(limit).Find(&invoices)
	handleError(result.Error, "error has happened while processing dbReadAllReceived query")

	return invoices
}

func DbReadAllSent(DBConn *gorm.DB, fromUserId uint, skip int, limit int) []Invoice {
	var invoices []Invoice

	result := DBConn.Where("FromID = ?", fromUserId).Offset(skip).Limit(limit).Find(&invoices)
	handleError(result.Error, "error has happened while processing dbReadAllSent query")

	return invoices
}

package main

import (
	"time"

	"gorm.io/gorm"
)

func dbmSendVoiceMail(DBConn *gorm.DB, message Invoice) (bool, uint) {

	result := DBConn.Create(&message)
	handleError(result.Error, "error has happened while processing dbmSendVoiceMail query")
	if result.RowsAffected == 0 {
		return false, 0
	}

	return true, message.ID
}

func dbMarkAsSeen(DBConn *gorm.DB, messageid uint) (bool, uint) {
	var invoice Invoice

	result := DBConn.First(&invoice, messageid)
	handleError(result.Error, "error has happened while processing dbMarkAsSeen query")
	if result.RowsAffected == 0 {
		return false, 0
	}

	return true, messageid
}

func dbReadMessageTime(DBConn *gorm.DB, messageid uint) (bool, bool, time.time) {
	var invoice Invoice

	result := DBConn.First(&invoice, messageid)
	handleError(result.Error, "error has happened while processing dbReadMessageTime query")
    
	var Read := false
	if result.RowsAffected == 0 {
		return false, false, time.Now()
	}
	if invoice.CreatedAt != invoice.UpdatedAt {
		Read = true
	}
	

	return true, Read, invoice.UpdatedAt
}

func dbReadUnSeenReceived(DBConn *gorm.DB, toUserId uint, skip int, limit int) []Invoice {
	var invoices []Invoice

	result := DBConn.Where("ToID = ? and Seen = false", toUserId).Offset(skip).Limit(limit).Find(&invoices)
	handleError(result.Error, "error has happened while processing dbReadUnSeenReceived query")

	return invoices
}

func dbReadAllReceived(DBConn *gorm.DB, toUserId uint, skip int, limit int) []Invoice {
	var invoices []Invoice

	result := DBConn.Where("ToID = ?", toUserId).Offset(skip).Limit(limit).Find(&invoices)
	handleError(result.Error, "error has happened while processing dbReadAllReceived query")

	return invoices
}

func dbReadAllSent(DBConn *gorm.DB, fromUserId uint, skip int, limit int) []Invoice {
	var invoices []Invoice

	result := DBConn.Where("FromID = ?", fromUserId).Offset(skip).Limit(limit).Find(&invoices)
	handleError(result.Error, "error has happened while processing dbReadAllSent query")

	return invoices
}

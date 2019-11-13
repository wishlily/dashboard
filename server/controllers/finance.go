package controllers

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/wishlily/dashboard/server/account"
)

var (
	// Finance url /api/finance http func
	Finance finance
)

type finance struct {
	sync.Mutex
}

// Init data load
func (f *finance) Init() error {
	return account.Init()
}

// Records finance GET method
//   /api/finance/record?start=yy-mm-dd hh:mm:ss&end=yy-mm-dd hh:mm:ss
func (f *finance) Records(c *gin.Context) {
	type records struct {
		Start string `form:"start" binding:"required"`
		End   string `form:"end" binding:"required"`
	}

	f.Lock()
	defer f.Unlock()

	var param records
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad parameter request!"})
		return
	}
	rs, err := account.Records(param.Start, param.End)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rs)
}

// Record finance POST method
//   /api/finance/record
//     type: add/del/chg
func (f *finance) Record(c *gin.Context) {
	type record struct {
		Type string         `json:"type" binding:"required"`
		Data account.Record `json:"data" binding:"required"`
	}

	f.Lock()
	defer f.Unlock()

	var rd record
	if err := c.ShouldBind(&rd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad parameter request!"})
		return
	}
	if strings.EqualFold(rd.Type, "add") {
		if err := rd.Data.Add(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error()})
			return
		}
	} else if strings.EqualFold(rd.Type, "del") {
		if err := rd.Data.Del(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error()})
			return
		}
	} else if strings.EqualFold(rd.Type, "chg") {
		if err := rd.Data.Chg(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad parameter request!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// Accounts finance GET method
//   /api/finance/account
//     list: account name list
func (f *finance) Accounts(c *gin.Context) {
	f.Lock()
	defer f.Unlock()

	accts, err := account.Accounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error()})
		return
	}
	if _, ok := c.GetQuery("list"); ok {
		// name list
		var list []string
		for _, a := range accts {
			if len(a.Account.ID) > 0 {
				list = append(list, a.Account.ID)
			}
		}
		c.JSON(http.StatusOK, list)
		return
	}
	c.JSON(http.StatusOK, accts)
}

// Account finance POST method
//   /api/finance/account
func (f *finance) Account(c *gin.Context) {
	type acct struct {
		Type string          `json:"type" binding:"required"`
		Data account.Account `json:"data" binding:"required"`
	}

	f.Lock()
	defer f.Unlock()

	var at acct
	if err := c.ShouldBind(&at); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad parameter request!"})
		return
	}
	if strings.EqualFold(at.Type, "add") ||
		strings.EqualFold(at.Type, "chg") {
		if err := at.Data.Update(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error()})
			return
		}
	} else if strings.EqualFold(at.Type, "del") {
		var a account.Account
		a.Account.ID = at.Data.Account.ID
		a.Debit.ID = at.Data.Debit.ID
		if err := a.Update(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad parameter request!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

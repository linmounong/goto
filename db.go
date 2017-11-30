package main

import (
	"errors"
	"flag"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	db         *gorm.DB
	dbFilename = flag.String("db", "goto.db", "filename of the SQLite file")

	predefinedLinks = []GotoInfo{
		{Key: "edit", Url: "/h"},
		{Key: "ee", Url: "/h"},
		{Key: "search", Url: "/h/s/"},
		{Key: "ss", Url: "/h/s/"},
	}
)

type GotoInfo struct {
	Key      string `gorm:"primary_key"`
	Url      string `gorm:"index"`
	Owner    string `gorm:"index"`
	UseCount int    `gorm:"index"`
}

func DbOpen() {
	if db != nil {
		log.Fatal("openDb() was called")
	}
	var err error
	db, err = gorm.Open("sqlite3", *dbFilename)
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&GotoInfo{})
	for _, info := range predefinedLinks {
		if DbFind(info.Key) == nil {
			DbCreate(&info)
		}
	}
}

func DbClose() {
	if db == nil {
		log.Fatal("openDb() never called")
	}
	db.Close()
	db = nil
}

func DbCreate(info *GotoInfo) error {
	db.Create(info)
	return nil
}

func DbFind(key string) *GotoInfo {
	var results []GotoInfo
	db.Where("key = ?", key).First(&results)
	if len(results) == 0 {
		return nil
	}
	return &results[0]
}

func DbUpdateOrCreate(username, key, url string) (*GotoInfo, error) {
	info := DbFind(key)
	if info == nil {
		info = &GotoInfo{
			Key:   key,
			Url:   url,
			Owner: username,
		}
		if err := DbCreate(info); err == nil {
			return info, nil
		} else {
			return nil, err
		}
	}
	if info.Owner != username {
		if info.Owner == "" {
			db.Model(info).UpdateColumn("owner", username)
		} else {
			return nil, errors.New("not your link")
		}
	}
	if info.Url != url {
		db.Model(info).UpdateColumn("url", url)
	}
	return info, nil
}

func DbIncr(key string) {
	info := DbFind(key)
	if info != nil {
		db.Model(info).UpdateColumn("use_count", info.UseCount+1)
	}
}

func DbRemove(key string) {
	if key != "" {
		db.Delete(&GotoInfo{Key: key})
	}
}

func DbSearch(keyPattern, urlPattern, username string, page, per int) []GotoInfo {
	if page <= 0 || per <= 0 {
		return []GotoInfo{}
	}
	skip := (page - 1) * per
	s := db.Model(&GotoInfo{}).Offset(skip).Limit(per).Order("use_count desc")
	if keyPattern != "" {
		s = s.Where("key LIKE ?", "%"+keyPattern+"%")
	}
	if urlPattern != "" {
		s = s.Where("url LIKE ?", "%"+urlPattern+"%")
	}
	if username != "" {
		s = s.Where("owner = ?", username)
	}
	var results []GotoInfo
	s.Find(&results)
	return results
}

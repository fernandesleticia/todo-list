package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

var db, _ = gorm.Open("mysql", "root:root@/agenda?charset=utf8&parseTime=True&loc=Local")

type AgendaItemModel struct {
	Id          int `gorm:primary_key`
	Description string
	Done        bool
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("All good with Agenda API")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `I'am alive`)
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	log.WithFields(log.Fields{"description": description}).Info("Add a new item")
	agenda_item := &AgendaItemModel{Description: description, Done: false}
	db.Create(&agenda_item)
	result := db.Last(&agenda_item)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Value)
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	has_item := GetItemByID(id)
	if has_item == false {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated": false, "error": Recorde Not Found}`)
	} else {
		done, _ := strconv.ParseBool(r.FormValue("done"))
		log.WithFields(log.Fields{"Id": id, "Done": done}).Info("Updating item")
		item := &AgendaItemModel{}
		db.First(&item, id)
		item.Done = done
		db.Save(&item)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated": true}`)
	}
}

func GetItemByID(Id int) bool {
	item := &AgendaItemModel{}
	result := db.First(&item, Id)
	if result.Error != nil {
		log.Warn("Item not found in database")
		return false

	}
	return true
}

func main() {
	defer db.Close()

	db.Debug().DropTableIfExists(&AgendaItemModel{})
	db.Debug().AutoMigrate(&AgendaItemModel{})

	log.Info("Starting agenda")
	router := mux.NewRouter()
	router.HandleFunc("/healthz", Healthz).Methods("GET")
	router.HandleFunc("/item", CreateItem).Methods("POST")
	router.HandleFunc("/update/{id}", UpdateItem).Methods("POST")
	http.ListenAndServe(":8000", router)
}

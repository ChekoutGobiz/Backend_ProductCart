package config

import (
	"os"

	"github.com/gocroot/helper/atdb"
)

var MongoString string = os.Getenv("mongodb+srv://dewidesember20:x6Wl5XF3ZNE1QcnS@cluster0.gkatt.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")

var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "jajankuy",
}

var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)

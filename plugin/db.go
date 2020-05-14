package mapper

import r "gopkg.in/rethinkdb/rethinkdb-go.v6"

var MapperDB = r.DBCreate("mapper")

var POITable = r.DB("mapper").TableCreate(
	"poi",
	r.TableCreateOpts{
		PrimaryKey: "id",
	},
)

var BoundsTable = r.DB("mapper").TableCreate(
	"bounds",
	r.TableCreateOpts{
		PrimaryKey: "id",
	},
)

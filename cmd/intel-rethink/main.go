package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/superloach/go-intel"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const defVers = "a8ca614df70e09516b36f060ef0304464e29dc75"

var (
	// process options
	conc     = flag.Int("conc", 100, "number of concurrent jobs")
	timeout  = flag.Duration("timeout", time.Second*5, "max request timeout")
	logint   = flag.Int("logint", 25, "interval to log progress")
	maxtries = flag.Int("maxtries", 10, "maximum times to retry on error")

	// network options
	base   = flag.String("base", "intel.ingress.com", "intel site url")
	secure = flag.Bool("secure", true, "use https")
	ua     = flag.String("ua", "Foo Bar Browser", "user agent")
	vers   = flag.String("version", defVers, "internal intel version")

	// auth options
	csrf   = flag.String("csrf", "", "csrf token")
	sessid = flag.String("sessid", "", "google(?) session id")

	// location options
	lat  = flag.Float64("lat", 0.0, "latitude")
	lng  = flag.Float64("lng", 0.0, "longitude")
	zoom = flag.Int("zoom", 17, "zoom level 0-17")

	// db options
	dburl   = flag.String("dburl", "", "rethink url")
	dbname  = flag.String("dbname", "", "rethink db name")
	dbtable = flag.String("dbtable", "", "rethink table name")
)

func dedup(orig []string) []string {
	keys := make(map[string]struct{})
	list := make([]string, 0)
	for _, item := range orig {
		if _, ok := keys[item]; !ok {
			keys[item] = struct{}{}
			list = append(list, item)
		}
	}
	return list
}

func main() {
	flag.Parse()

	if *csrf == "" {
		panic("please provide a csrf token")
	}
	if *sessid == "" {
		panic("please provide a sessid")
	}

	client, err := intel.NewClient()
	if err != nil {
		panic(err)
	}

	client.Client.Timeout = *timeout
	client.MaxTries = *maxtries

	client.Base = *base
	client.Secure = *secure
	client.UA = *ua
	client.Version = *vers

	client.CSRF = *csrf
	client.SessID = *sessid

	if *dbname == "" {
		panic("please provide a db name")
	}
	if *dbtable == "" {
		panic("please provide a db table")
	}

	db, err := r.Connect(r.ConnectOpts{
		Address: *dburl,
		Timeout: *timeout,
	})
	if err != nil {
		panic(err)
	}

	tileKey := intel.TileKey(*lat, *lng, *zoom)
	fmt.Println("tile key:", tileKey)

	portalIDs, err := client.PortalIDs([]string{tileKey})
	if err != nil {
		panic(err)
	}
	portalIDs = dedup(portalIDs)
	l := len(portalIDs)

	jobs := make(chan struct{}, *conc)
	done := make(chan struct{})

	for _, portalID := range portalIDs {
		go func(guid string) {
			jobs <- struct{}{}

			portal, err := client.GetPortal(guid)
			if err != nil {
				<-jobs
				l--
				fmt.Println(err)
				return
			}

			_, err = r.DB(*dbname).Table(*dbtable).Insert(portal).Run(db)
			if err != nil {
				<-jobs
				l--
				fmt.Println(err)
				return
			}

			done <- <-jobs
		}(portalID)
	}

	i := 0
	for range done {
		i++
		if i%*logint == 0 {
			fmt.Printf("%d/%d (%.1f%%)\n", i, l, float64(i)/float64(l)*100)
		}
		if i == l {
			close(done)
		}
	}

	ol := len(portalIDs)
	fmt.Printf(
		"lost %d of %d (%.1f%%)\n",
		ol-l, ol, float64(ol-l)/float64(ol)*100,
	)
}

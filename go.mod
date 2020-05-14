module github.com/superloach/mapper

go 1.14

require (
	github.com/Necroforger/dgwidgets v0.0.0-20190131052008-56c8c1ca33e0
	github.com/bwmarrin/discordgo v0.20.3
	github.com/go-snart/bot v0.0.0-00010101000000-000000000000
	github.com/go-snart/db v0.0.0-00010101000000-000000000000
	github.com/go-snart/plugin-admin v0.0.0-00010101000000-000000000000
	github.com/go-snart/plugin-help v0.0.0-00010101000000-000000000000
	github.com/go-snart/route v0.0.0-00010101000000-000000000000
	github.com/namsral/flag v1.7.4-pre
	github.com/paul-mannino/go-fuzzywuzzy v0.0.0-20200127021948-54652b135d0e
	github.com/superloach/minori v0.0.0-20200401022729-31f6f02808bc
	gopkg.in/rethinkdb/rethinkdb-go.v6 v6.2.1
)

replace (
	github.com/go-snart/bot => ../../go-snart/bot
	github.com/go-snart/db => ../../go-snart/db
	github.com/go-snart/plugin-admin => ../../go-snart/plugin-admin
	github.com/go-snart/plugin-help => ../../go-snart/plugin-help
	github.com/go-snart/route => ../../go-snart/route
)

module juzhao

go 1.14

require (
	github.com/controller v0.0.0-00010101000000-000000000000 // indirect
	github.com/flosch/pongo2 v0.0.0-20200518135938-dfb43dbdc22a // indirect
	github.com/gin-gonic/gin v1.6.3 // indirect
	github.com/juju/errors v0.0.0-20200330140219-3fe23663418f // indirect
	github.com/middleware v0.0.0-00010101000000-000000000000 // indirect
	github.com/model v0.0.0-00010101000000-000000000000 // indirect
	github.com/mongo v0.0.0-00010101000000-000000000000
	github.com/ratelimit v0.0.0-00010101000000-000000000000 // indirect
	github.com/router v0.0.0-00010101000000-000000000000
	github.com/service v0.0.0-00010101000000-000000000000 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/ini.v1 v1.56.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace (
	github.com/controller => ./controller
	github.com/middleware => ./middleware
	github.com/model => ./model
	github.com/mongo => ./mongo
	github.com/ratelimit => ./ratelimit
	github.com/router => ./router
	github.com/service => ./service
)

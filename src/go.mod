module juzhao

go 1.14

require (
	github.com/controller v0.0.0-00010101000000-000000000000 // indirect
	github.com/flosch/pongo2 v0.0.0-20200518135938-dfb43dbdc22a // indirect
	github.com/gin-gonic/gin v1.6.3 // indirect
	github.com/juju/errors v0.0.0-20200330140219-3fe23663418f // indirect
	github.com/middleware v0.0.0-00010101000000-000000000000 // indirect
	github.com/router v0.0.0-00010101000000-000000000000
	github.com/smartystreets/goconvey v1.6.4 // indirect
	gopkg.in/ini.v1 v1.56.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace (
	github.com/controller => ./controller
	github.com/middleware => ./middleware
	github.com/router => ./router
)

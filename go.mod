module github.com/thetreep/covidtracker

go 1.12

require (
	github.com/google/uuid v1.1.1
	github.com/graphql-go/graphql v0.7.9
	github.com/graphql-go/handler v0.2.3
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.6.0
	github.com/thetreep/toolbox v0.0.0-20200526130145-8ad9a40150e2
	go.mongodb.org/mongo-driver v1.3.3
	golang.org/x/net v0.0.0-20200520182314-0ba52f642ac2
)

replace github.com/thetreep/toolbox => ../toolbox

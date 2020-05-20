module github.com/thetreep/covidtracker

go 1.12

require (
	github.com/google/uuid v1.1.1
	github.com/graphql-go/graphql v0.7.9
	github.com/graphql-go/handler v0.2.3
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/thetreep/toolbox v0.0.0-20200519135852-f98c1ce884cc // indirect
	go.mongodb.org/mongo-driver v1.3.3
)

replace github.com/thetreep/toolbox => ../toolbox

module github.com/silinternational/rest-data-archiver

go 1.22

replace github.com/silinternational/rest-data-archiver => ./

require (
	github.com/aws/aws-lambda-go v1.38.0
	github.com/aws/aws-sdk-go v1.51.9
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

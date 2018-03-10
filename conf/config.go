package conf

const (
	// ElasticURL to connect to, elasticsearch provides a relative urlpath within docker
	ElasticURL = "http://linkesdb:9200"
	// ElasticURL = "http://127.0.0.1:9200"
	// ElasticURL   = "elasticsearch://linkesdb:9200" // no es node found
	indexName = "yo"
	docType   = "log"
	appName   = "myApp"
)

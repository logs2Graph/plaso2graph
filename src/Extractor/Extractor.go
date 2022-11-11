package Extractor

import (
	"log"
)

// Contains generic functions for extractors and the Extract function which calls the correct extractor
// depending on the extractor specified in the args map

func Extract(data []interface{}, args map[string]interface{}) {
	if args["extractor"] == nil {
		log.Fatal("No extractor specified")
	}

	var extractor = args["extractor"].(string)
	switch extractor {
	case "neo4j":
		Neo4jExtract(data, args)
		break
	case "json":
		JsonExtract(data, args)
		break
	case "xml":
		XmlExtract(data, args)
		break
	case "csv":
		CsvExtract(data, args)
		break
	}
}

func InitializeExtractor(args map[string]interface{}) map[string]interface{} {
	if args["extractor"] == nil {
		log.Fatal("No extractor specified")
	}

	var extractor = args["extractor"].(string)

	switch extractor {
	case "neo4j":
		return InitializeNeo4jExtractor(args)
	case "json":
		return InitializeJsonExtractor(args)
	case "xml":
		return InitializeXmlExtractor(args)
	case "csv":
		return InitializeCsvExtractor(args)
	}
	return args
}

func PostProcessing(args map[string]interface{}) {
	if args["extractor"] == nil {
		log.Fatal("No extractor specified")
	}

	var extractor = args["extractor"].(string)

	switch extractor {
	case "neo4j":
		Neo4jPostProcessing(args)
		break
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

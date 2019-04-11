package helper

import (
	"strings"
)

// Helper
func GenerateInsertFields(fieldsArr []string, fieldPrefix string, fieldException []string) (string) {
	var result []string
	isException := false
	for _,field := range fieldsArr {
		isException = false
		for _,field2 := range fieldException {
			if field == field2 { isException = true }
		}
		if !isException { result = append(result, fieldPrefix + field) } 
	}
	return strings.Join(result, ",")
}

func GenerateUpdateFields(fieldsArr []string, fieldException []string) (string) {
	var result []string
	isException := false
	for _,field := range fieldsArr {
		isException = false
		for _,field2 := range fieldException {
			if field == field2 { isException = true }
		}
		if !isException { result = append(result, field + " = :" + field) } 
	}
	return strings.Join(result, ",")
}

func GenerateUpdateFieldsPrimaryKeyConstraint(primaryKeyArr []string) (string) {
	var result []string
	for _,field := range primaryKeyArr {
		result = append(result, field + " = :" + field) 
	}
	return strings.Join(result, " AND ")
}

package utils

import (
	"strconv"
	"strings"

	"github.com/playedu/playedu-go/pkg/constants"
)

func GetPageParams(pageStr, sizeStr string) (int, int) {
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)

	if page <= 0 {
		page = constants.DefaultPage
	}
	if size <= 0 {
		size = constants.DefaultPageSize
	}
	if size > constants.MaxPageSize {
		size = constants.MaxPageSize
	}

	return page, size
}

func GetOffset(page, size int) int {
	return (page - 1) * size
}

func InArray(needle interface{}, haystack interface{}) bool {
	switch v := haystack.(type) {
	case []int:
		n, ok := needle.(int)
		if !ok {
			return false
		}
		for _, item := range v {
			if item == n {
				return true
			}
		}
	case []string:
		n, ok := needle.(string)
		if !ok {
			return false
		}
		for _, item := range v {
			if item == n {
				return true
			}
		}
	}
	return false
}

func StringToIntSlice(str, separator string) []int {
	if str == "" {
		return []int{}
	}
	
	parts := strings.Split(str, separator)
	result := make([]int, 0, len(parts))
	
	for _, part := range parts {
		if num, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
			result = append(result, num)
		}
	}
	
	return result
}

func RemoveDuplicates(slice []int) []int {
	keys := make(map[int]bool)
	result := []int{}

	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			result = append(result, entry)
		}
	}

	return result
}

func CalculateProgress(finished, total int) int {
	if total == 0 {
		return 0
	}
	return (finished * 100) / total
}

func IsFinished(finishedDuration, totalDuration int) bool {
	if totalDuration == 0 {
		return false
	}
	progress := (finishedDuration * 100) / totalDuration
	return progress >= constants.FinishThreshold
}

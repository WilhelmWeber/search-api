package controllers

import (
	"net/http"
	"strconv"

	"github.com/WilhelmWeber/search-api/src/libs"
	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	var page_num int //複数ページ表示用
	var show_annotations [][]libs.Annotation

	const BASE_URI string = "http://localhost:3000"
	q := c.Query("q")
	motivation := c.Query("motivation")
	date := c.Query("date")
	user := c.Query("userId")
	page := c.Query("page")

	queryPath := BASE_URI + c.Request.URL.Path

	client := libs.DBConnect()
	annotations := libs.GetAnnotations(q, motivation, user, date, client)

	page_num = (len(annotations) / 10)
	if len(annotations)%page_num == 0 {
		page_num -= 1
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "query param of 'page' must be numerical",
		})
	}
	if pageInt > (page_num + 1) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not Found",
		})
	}

	if page_num > 1 {
		k := 0
		for i := 0; i < page_num; i++ {
			tmp := show_annotations[i]
			for j := 0; j < 10; j++ {
				tmp = append(tmp, annotations[k])
				k += 1
				if k == (len(annotations) - 1) {
					break
				}
			}
		}
	} else {
		show_annotations[0] = annotations
	}

	within := gin.H{
		"@type": "sc:Layer",
		"total": len(annotations),
		"first": queryPath + "&page=1",
		"last":  queryPath + "&page=" + strconv.Itoa(page_num+1),
	}

	var resources []gin.H

	for _, anno := range show_annotations[pageInt-1] {
		resource := gin.H{
			"@type": "cnt:ContentAsText",
			"chars": anno.Chars,
		}
		elem := gin.H{
			"@id":        BASE_URI + "/presentation/" + anno.Manifest_id + "/annolist.json",
			"@type":      "oa:Annotation",
			"motivation": anno.Motivation,
			"resource":   resource,
			"on":         anno.On,
		}
		resources = append(resources, elem)
	}

	c.JSON(http.StatusOK, gin.H{
		"@context":   "http://iiif.io/api/presentation/2/context.json",
		"@id":        queryPath + "&page=" + page,
		"@type":      "sc:AnnotationList",
		"within":     within,
		"startIndex": strconv.Itoa((pageInt - 1) * 10),
		"resources":  resources,
	})
}

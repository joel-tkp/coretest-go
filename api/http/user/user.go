package user

import (
	"net/http"
    "github.com/julienschmidt/httprouter"
	"strconv" // String Converter
	"time" // time lib

	"References/coretest/api"
	"References/coretest/api/http/base"
	"References/coretest/service/user"
)

var userService api.UserService

func Init(service api.UserService) {
	userService = service
}

func Create(w http.ResponseWriter, r *http.Request, qs httprouter.Params) {
	start := time.Now()
	u := userService.Create(r.FormValue("name"), r.FormValue("email"), true, "testIdempotencyKey")
	// return response as JSON
	base.WriteSuccessResponse(w, base.Success, "User " + u.Name + " created!", start, u)
}

func Detail(w http.ResponseWriter, r *http.Request, qs httprouter.Params) {
	start := time.Now()
	objId,_ := strconv.ParseInt(qs.ByName("id"), 0, 64)
	u, err := userService.Get(objId)
	if err != nil {
		base.WriteSuccessResponse(w, base.Error, "User not found!", start, u)
	}
	// return response as JSON
	base.WriteSuccessResponse(w, base.Success, "User " + u.Name + " retrieved!", start, u)
}

func List(w http.ResponseWriter, r *http.Request, qs httprouter.Params) {
	start := time.Now()
	orderBy := r.URL.Query().Get("orderBy")
	orderDir := r.URL.Query().Get("orderDir")
	paginated := r.URL.Query().Get("paginated")

	var objectList []user.User
	if paginated == "no" {
		objectList, _ = userService.AllList(orderBy, orderDir)
	} else {
		page,_ := strconv.ParseInt(r.URL.Query().Get("page"), 0, 64)
		objectList, _ = userService.PaginatedList(10, int32(page), orderBy, orderDir)
	}
	// return response as JSON
	base.WriteSuccessResponse(w, base.Success, "User list retrieved!", start, objectList)
}

func Update(w http.ResponseWriter, r *http.Request, qs httprouter.Params) {
	start := time.Now()
	objId,_ := strconv.ParseInt(qs.ByName("id"), 0, 64)
	u, err := userService.Get(objId)
	if err != nil {
		base.WriteSuccessResponse(w, base.Error, "User not found!", start, u)
	}
	u = userService.Update(u.ID, r.FormValue("name"), r.FormValue("email"), true, "testIdempotencyKey")
	// return response as JSON
	base.WriteSuccessResponse(w, base.Success, "User " + u.Name + " updated!", start, u)
}

func Delete(w http.ResponseWriter, r *http.Request, qs httprouter.Params) {
	start := time.Now()
	objId,_ := strconv.ParseInt(qs.ByName("id"), 0, 64)
	u, err := userService.Get(objId)
	if err != nil {
		base.WriteSuccessResponse(w, base.Error, "User not found!", start, u)
	}
	userService.Delete(u.ID)
	// return response as JSON
	base.WriteSuccessResponse(w, base.Success, "User " + u.Name + " created!", start, u)
}

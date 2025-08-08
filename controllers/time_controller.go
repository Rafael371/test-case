package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"test-case/utils"
)

func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	//get param
	hStr := r.URL.Query().Get("hour")
	mStr := r.URL.Query().Get("minute")
	sStr := r.URL.Query().Get("second")

	//convert to int
	h, err1 := strconv.Atoi(hStr)
	m, err2 := strconv.Atoi(mStr)
	s, err3 := strconv.Atoi(sStr)

	//handle input null
	if err1 != nil || err2 != nil || err3 != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	//handle 00:00:00
	if h == 0 && m == 0 && s == 0 {
		h = 24
	}

	//get total second on earth time
	total := utils.EarthToTotalSeconds(h,m,s)

	//convert to Roketin hour
	ah, am, as := utils.ConvertToRoketinHour(total)

	fmt.Fprintf(w, "On Earth: %02d:%02d:%02d → On ABC Planet: %02d:%02d:%02d\n", h, m, s, ah, am, as)
}
package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"test-case/utils"
)

func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	hStr := r.URL.Query().Get("hour")
	mStr := r.URL.Query().Get("minute")
	sStr := r.URL.Query().Get("second")

	h, err1 := strconv.Atoi(hStr)
	m, err2 := strconv.Atoi(mStr)
	s, err3 := strconv.Atoi(sStr)

	if err1 != nil || err2 != nil || err3 != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	if h == 0 && m == 0 && s == 0 {
		h = 24
	}

	total := utils.EarthToTotalSeconds(h,m,s)
	ah, am, as := utils.TotalSecondsToRoketin(total)

	fmt.Fprintf(w, "On Earth: %02d:%02d:%02d → On ABC Planet: %02d:%02d:%02d\n", h, m, s, ah, am, as)
}
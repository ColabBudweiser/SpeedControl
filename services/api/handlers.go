package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	da "speedСontrol/dataAccess"
	m "speedСontrol/models"
)

func HandleInputMessage(w http.ResponseWriter, r *http.Request) {
	var q m.SpeedInfoInputQuery

	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if res := ValidateInputParams(q, w); res == false {
		return
	}

	msg := m.NewSpeedControlMsg(q.DateTime, q.Number, q.Speed)
	_, _ = fmt.Fprintf(w, "%s %s %s - информация получена", q.DateTime, q.Number, q.Speed)
	da.SaveSpeedControlInfo(*msg)
}

func GetInfoByDate(w http.ResponseWriter, r *http.Request) {
	d := r.URL.Query().Get("date")
	s := r.URL.Query().Get("speed")

	i := m.NewInfoByDateQuery(d, s)
	if res := ValidateInputParams(i, w); res == false {
		return
	}

	msg := m.NewByDateQuery(d, s)
	da.WriteInfoByDate(*msg, w)
}

func GetExtremesByDate(w http.ResponseWriter, r *http.Request) {
	d := r.URL.Query().Get("date")

	i := m.NewExtremesByDateQuery(d)
	if res := ValidateInputParams(i, w); res == false {
		return
	}

	msg := m.NewExtremesMessage(d)
	da.WriteExtremesByDate(*msg, w)
}

func ValidateInputParams(i m.ValidationType, w http.ResponseWriter) bool {
	if err := i.Validate(); err != nil {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintf(w, "Неверный формат входных данных")
		return false
	}
	return true
}

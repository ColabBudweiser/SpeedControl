package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"speedСontrol/services/api"
	c "speedСontrol/services/config"
	"time"
)

func GetHttpServer() *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/speedControlMessages", api.HandleInputMessage)
	router.HandleFunc("/speedControlsByDate", checkAccessTime(api.GetInfoByDate))
	router.HandleFunc("/extremesByDate", checkAccessTime(api.GetExtremesByDate))
	http.Handle("/", router)

	return &http.Server{
		Addr:    fmt.Sprintf("%s%s", c.Conf.Host, c.Conf.ListenPort),
		Handler: nil,
	}
}

func checkAccessTime(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		now := time.Now()
		h := now.Hour()
		m := now.Minute()

		var ok bool
		if c.Conf.Start.Hour < h && h < c.Conf.Finish.Hour {
			ok = true
		} else if c.Conf.Start.Hour == h {
			if c.Conf.Start.Minute <= m {
				ok = true
			}
		} else if c.Conf.Finish.Hour == h {
			if c.Conf.Finish.Minute >= m {
				ok = true
			}
		} else {
			ok = false
		}

		if !ok {
			_, _ = fmt.Fprintf(res, "Время доступа к запросам на выборку данных - %s-%s",
				c.Conf.QueryStartTime,
				c.Conf.QueryFinishTime)
			res.WriteHeader(http.StatusForbidden)
			return
		}
		next(res, req)
	}
}

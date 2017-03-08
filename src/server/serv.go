package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/CanonicalLtd/UCWifiConnect/netman"
)

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func para(s string) string {
	return fmt.Sprintf("<p>%s</p>", s)
}

func form(SSIDs []netman.SSID) string {
	form_ := "<form>"
	for _, s := range SSIDs {
		line := fmt.Sprintf("<input type='radio' name='essid' value='%s' checked>%s<br>", s.Ssid, s.Ssid)
		fmt.Println(line)
		form_ = form_ + line 
	}
	form_ = form_ + "</form>"
	return form_
}

func handler(w http.ResponseWriter, r *http.Request) {
	d, _ := os.Getwd()
	//fmt.Fprintf(w, "<p>URL path: %s</p>", r.URL.Path[1:])
	fmt.Fprintf(w, "<p>pwd: %s", d)
	ssids, _, _ := netman.Ssids()
	ssids_form := form(ssids)
	fmt.Fprintf(w, ssids_form)
	//for _, s := range getWifi() {
	//	fmt.Fprintf(w, para(s))
	//}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

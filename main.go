package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/tidwall/gjson"
)

func generateHeader(name string) string {
	return fmt.Sprintf(`

//------------------------------------------------
//                                              
//                   %s
//                                              
//------------------------------------------------`, name)
}

func buybindGen(key string, item string) string {
	return fmt.Sprintf(
		`
bind %s "buy %s"`, key, item)

}

func keybindGen(key string, action string) string {
	return fmt.Sprintf(
		`
bind %s "%s"`, key, action)

}

func dropbindGen(key string) string {
	return fmt.Sprintf(
		`
alias "dropbind" "buy ak47; buy m4a1; slot 1; drop; say_team WEAPON DROPPED"
bind "%s" "dropbind"`, key)
}

func jumpbindGen(key string) string {
	return fmt.Sprintf(
		`
alias "+jumpthrow" "+jump; -attack"
alias "-jumpthrow" "-jump"
bind "%s" "+jumpthrow"`, key)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func generateAutoexec(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	check(err)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(body))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(body))

	log.Printf("BODY: %q", rdr1)
	req.Body = rdr2
}

func main() {
	var autoexec string

	f, err := os.Create("/tmp/autoexec.cfg")
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(autoexec)

	srcJSON := `{"binds": {"buybinds": {"ump45": "kp_slash"},"miscbinds": {"use weapon_c4": "c","use weapon_c4;drop": "z","drop": "l","jumpthrow": "capslock"}}}`

	miscbinds := gjson.Get(srcJSON, "binds.miscbinds").Map()
	buybinds := gjson.Get(srcJSON, "binds.buybinds").Map()

	autoexec += generateHeader("Misc Binds")
	for k, v := range miscbinds {
		var tmp string
		switch k {
		case "jumpthrow":
			tmp = jumpbindGen(v.String())
		case "drop":
			tmp = dropbindGen(v.String())
		default:
			tmp = keybindGen(v.String(), k)
		}
		autoexec += tmp
	}

	autoexec += generateHeader("Buybinds")
	for k, v := range buybinds {
		autoexec += buybindGen(v.String(), k)
	}

	w.WriteString(autoexec)
	w.WriteString(fmt.Sprintf("\n\n\n\necho \"Autoexec.cfg loaded\" \nhost_writeconfig"))
	w.Flush()

	http.HandleFunc("/generateAutoexec", generateAutoexec)
	http.ListenAndServe(":8081", nil)
}

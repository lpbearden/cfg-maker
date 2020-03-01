package main

import (
	"bufio"
	"fmt"
	"os"
)

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

func main() {
	fmt.Println("generating keybinds")
	autoexec := buybindGen("kp_slash", "ump45") + dropbindGen("l") + jumpbindGen("capslock")

	f, err := os.Create("/tmp/autoexec.cfg")
	check(err)

	defer f.Close()

	w := bufio.NewWriter(f)
	n4, err := w.WriteString(autoexec)
	fmt.Printf("wrote %d bytes\n", n4)

	w.Flush()
}

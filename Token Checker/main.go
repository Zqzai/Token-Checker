package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	cyan      = "\033[36m"
	red       = "\033[31m"
	blue      = "\033[34m"
	purple    = "\033[35m"
	green     = "\033[32m"
	yellow    = "\033[33m"
	orange    = "\033[33m"
	reset     = "\033[39m\u001b[24m"
	underline = "\u001b[4m"
	locked    = 0
	valid     = 0
	invalid   = 0
)

var tokenschecked = 0

func checktoken(token string) {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", "https://discord.com/api/v9/users/@me/affinities/guilds", nil)
	request.Header.Set("authorization", token)
	resp, _ := client.Do(request)
	if resp.StatusCode == 200 {
		fmt.Println("(" + green + underline + "VALID" + reset + "): " + token)
		valid++
		f, err := os.OpenFile("valid.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(token + "\n"); err != nil {
			fmt.Println(err)
		}
	} else if resp.StatusCode == 403 {
		fmt.Println("(" + orange + underline + "LOCKED" + reset + "): " + token)
		locked++
		f, err := os.OpenFile("locked.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(token + "\n"); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("(" + red + underline + "INVALID" + reset + "): " + token)
		invalid++
		f, err := os.OpenFile("invalid.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(token + "\n"); err != nil {
			fmt.Println(err)
		}
	}
	tokenschecked++
}

func main() {
	file, err := os.Open("tokens.txt")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	f, err := os.OpenFile("valid.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	f, err = os.OpenFile("locked.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	f, err = os.OpenFile("invalid.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	var tokens []string
	var token string

	for {
		_, err := fmt.Fscanln(file, &token)
		if err != nil {
			break
		}
		tokens = append(tokens, token)
	}

	starttime := time.Now()

	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range tokens {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	duplicates := len(tokens) - len(list)

	for _, token := range list {
		go checktoken(token)
	}

	totaltokens := len(list)

	for {
		if tokenschecked == totaltokens {
			finishtime := time.Now()
			elapsedseconds := finishtime.Sub(starttime).Seconds()
			elapsedsecondss := fmt.Sprintf("%.2f", elapsedseconds)

			fmt.Println("(" + purple + underline + "TOTAL" + reset + "): " + fmt.Sprint(totaltokens) + " (" + green + underline + "VALID" + reset + "): " + fmt.Sprint(valid) + " (" + orange + underline + "LOCKED" + reset + "): " + fmt.Sprint(locked) + " (" + red + underline + "INVALID" + reset + "): " + fmt.Sprint(invalid) + " (" + cyan + underline + "DUPLICATES" + reset + "): " + fmt.Sprint(duplicates) + " (" + blue + underline + "TIME" + reset + "): " + fmt.Sprint(elapsedsecondss) + "s")
			fmt.Scanln()
			break
		}
	}
}

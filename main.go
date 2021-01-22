package main

import (
    "fyne.io/fyne"
    "fyne.io/fyne/app"
    "fyne.io/fyne/widget"
    "encoding/json"
	"io/ioutil"
    "net/http"
    "fmt"
)

type tokenResponse struct {
    Token string
}

type LogsListItem struct {
	User_id		int 	`json:"user_id"`
	Email		string	`json:"email"`
	Date_time	string	`json:"date_time"`
	Ip_address	string	`json:"ip_address"`
}

type LogsList struct {
	List []LogsListItem `json:"list"`
}

func main() {
	myApp := app.New()
    myWindow := myApp.NewWindow("Login")
    myWindow.Resize(fyne.NewSize(200,200))
    login := widget.NewEntry()
    password := widget.NewEntry()
    loginInput := widget.NewLabel("Login")
    passwordInput := widget.NewLabel("Password")
	myWindow.SetContent(widget.NewVBox(
        loginInput,
        login,
        passwordInput,
        password,
		widget.NewButton("Authorize", func() {
            result := authorize(login.Text, password.Text)
            myWindow1 := myApp.NewWindow("Token")
            tokenText := widget.NewLabel("Your token - " + result)
            logsLabel := widget.NewLabel("Entry tour city:")
            logsEntry := widget.NewEntry()
            vBox := widget.NewVBox()
            myWindow1.SetContent(widget.NewVBox(
                tokenText,
                logsLabel,
                logsEntry,
                widget.NewButton("Get logs", func() {
                    logs := getLogs(logsEntry.Text, result)
                    for i := 0; i < len(logs.List); i++ {
                        widget.NewLabel(fmt.Sprintf(""))
                    }
                }),
            ))
            myWindow1.Resize(fyne.NewSize(500,100))
            myWindow1.Show()
        }),
    ))
    
    myWindow.ShowAndRun()
}

func authorize(login string, password string) string {
    resp, err := http.Get("http://localhost:9999/token?email=" + login +"&password=" + password)
	if err != nil {
		panic(err.Error())
    }
    body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err.Error())
    }
    var token = new(tokenResponse)
    err = json.Unmarshal(body, &token)
    if err != nil {
		panic(err.Error())
    }
    return token.Token
}

func getLogs(city string, token string) *LogsList {
    client := &http.Client{}
    req, _ := http.NewRequest("GET", "http://localhost:9999/logs/" + city, nil)
    req.Header.Set("Authorization","Bearer " + token)
    resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
    }
    body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err.Error())
    }
    var logs = new(LogsList)
    err = json.Unmarshal(body, &logs)
    if err != nil {
		panic(err.Error())
    }
    fmt.Print(logs.List)
    return logs
}
package main

import (
    "fyne.io/fyne"
    "fyne.io/fyne/app"
    "fyne.io/fyne/widget"
    "encoding/json"
	"io/ioutil"
    "net/http"
    "fmt"
    "fyne.io/fyne/data/validation"
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
    login.Validator = validation.NewRegexp(`^[A-Za-z0-9_-@]+$`, "password can only contain letters, numbers, '_', and '-'")
    password := widget.NewPasswordEntry()
    password.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "password can only contain letters, numbers, '_', and '-'")
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
            logsLabel := widget.NewLabel("Entry your city:")
            logsEntry := widget.NewEntry()
            vBox := widget.NewVBox()
            myWindow1.SetContent(widget.NewVBox(
                tokenText,
                logsLabel,
                logsEntry,
                widget.NewButton("Get logs", func() {
                    logs := getLogs(logsEntry.Text, result)
                    for i := 0; i < len(logs.List); i++ {
                        var logsMainGroup = widget.NewVBox(
                            widget.NewLabel(fmt.Sprintf("user_id : %d ", logs.List[i].User_id)),
                            widget.NewLabel(fmt.Sprintf("Email : %s ", logs.List[i].Email)),
                            widget.NewLabel(fmt.Sprintf("ip address : %s ", logs.List[i].Ip_address)),
                        )
                        vBox.Append(widget.NewGroup(logs.List[i].Date_time))
                        vBox.Append(widget.NewHBox(logsMainGroup))
                    }
                    vBox.Append(widget.NewButton("Закрыть", func() {
                        myApp.Quit()
                     }))
                     myWindow1.SetContent(vBox)
                     myWindow1.Resize(fyne.NewSize(500,700))
                }),
            ))
            myWindow1.Resize(fyne.NewSize(500,100))
            myWindow1.Show()
        }),
    ))
    myWindow.ShowAndRun()
}

func authorize(login string, password string) string {
    resp, err := http.Get("http://localhost:8080/token?email=" + login +"&password=" + password)
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
    req, _ := http.NewRequest("GET", "http://localhost:8080/logs/" + city, nil)
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
    return logs
}
# INTRO
This project is made for "Idee i Informatyka" course, which is on MIMUW faculty during masters degree.
Main Code Contributors:
- Kacper Rzetelski
- Franciszek Chnatów
- Michał Lange

# GOAL
Our goal is to create Minimal Working Project (MVP) of an app that will help people in terms of preventing phishing. This app can sent emails to choosen addresses that will contain a dangerous link. Later we can check if a person both clicked the mail and interacted with the malicious content of the webpage. 

# Project consist of several parts
## Golang
To use our app you need to install 'go'. You can do this using this page:
```
https://go.dev/doc/install
```
## Malicious Webpage
To start the server with malicious website you should follow instructions below. Starting from main project directory write this in command line:

```
go run ./cmd/iii/main.go server --loglevel=4
```
Default port is 8080.
The malicious webpage should be available at, but do not go there now:
```
http://localhost:8080/
```
## Sending email
Now to send an email, that we prepared use line below. 
```
go run ./cmd/iii/main.go sender --from=a.schubert@uw.edu.pl --to=test@test.com --username=9a346bcdbc771c --password=0d55fba525a381 --host=smtp.mailtrap.io --port=2525 --loglevel=4
```
 Most email platforms use very inteligent algorythms to filter out malicious emails, so for simplicity we will use "Mail Trap" platform for sending and receiving emails.


## Recieving email
To see the mail visit https://mailtrap.io/signin and login to our account (login and password was send by mail). Next go to "Email Testing" > "My Inbox". 
You should see a new email about (fake) registration deadline. This email was made to look just like emails from university.







Be carefull it looks just like the login page to USOS.


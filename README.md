# INTRO
This project is made for "Idee i Informatyka" course, which is on MIMUW faculty during masters degree.
Main Code Contributors:
- Kacper Rzetelski
- Franciszek Hnatów
- Michał Lange

# GOAL
Our goal is to create Minimal Working Project (MVP) of an app that will help people in terms of preventing phishing. This app can sent emails to chosen addresses that will contain a dangerous link. Later we can check if a person both clicked the mail and interacted with the malicious content of the webpage.

# Instructions to test the app
## Golang
To use our app you need to install 'go'. You can do this using this page:
```
https://go.dev/doc/install
```
## Malicious Webpage
To start the server with malicious website you should follow the instructions below. Starting from main project directory type the following in command line:

```
go run ./cmd/iii/main.go server --loglevel=4
```
If you want to specify address of the server, do this first to understand possible flags and arguments:
```
go run ./cmd/iii/main.go server --help
```
Closing this terminal will probably result in shutdown of the server. For the next step use a new terminal window. By default the server runs on:
```
http://localhost:8080/
```
But first we need to send the malicious email.
## Sending email
In order to send the prepared malicious email, type the line below into a new command line.
```
go run ./cmd/iii/main.go sender --from=a.schubert@uw.edu.pl --to=test@test.com --username=9a346bcdbc771c --password=0d55fba525a381 --host=smtp.mailtrap.io --port=2525 --loglevel=4
```
 Most email platforms use very sophisticated algorithms to filter out malicious emails, so for simplicity we will use "Mail Trap" platform for sending and receiving emails. 

 *At your own risk you can use other platforms, but you will need to supply correct argument values yourself.

## Receiving email
To see the mail visit https://mailtrap.io/signin and login to our account (login and password was sent by mail). Next go to "Email Testing" > "My Inbox".
You should see a new email about (fake) registration deadline. This email was made to look just like emails from university.
## Stats
Type the line below to command line to see status of this test:
```
curl localhost:8080/stats  -H "Content-Type: application-json" -s
```
It should print:
```
{"test@test.com":{"seen":true,"clicked":false,"logged_in":false}}
```
This is the storage of our server. This particular message means you saw the content of the email (when you opened it on the platform). Now you can click the link in mail and see fake USOS (be careful, it looks just like the login page). If you do so, the value "clicked" will change to true. If you log in to (fake) USOS your credentials will be shown in server terminal, the value "logged_in" will change to true, and you will be redirected to standard USOS login page.

## Closing
Thanks for testing our solution you can now quit by pressing "Ctrl + C" in terminal with the server.
# Description
App uses SQLite3 - DB was used via [sqlite browser](https://sqlitebrowser.org/dl/) to create and maintain <br />
Data to backend sends through url.

>Example of sending data to subscribe email: <br />http://localhost:8080/subscribe?email= **someMail**

Sends daily emails to subscribed e-mails with the latest exchange rate
## For running locally
1. Go have to be installed - [link](https://go.dev/doc/install)
2. gcc have to be installed - [example](https://code.visualstudio.com/docs/cpp/config-mingw#_installing-the-mingww64-toolchain)
3. To run app use command: ```go run ./main/.```

>[!IMPORTANT] 
> To being able to send emails in app should be added following values:
> * line 145: variable ***apw*** - app password for gmail account, generation example - [link](https://mailmeteor.com/blog/gmail-smtp-settings) 
> * line 146: variable ***sendingEmail*** - gmail account used for sending emails for which the password is specified in line 145
# Endpoints:
App runs on ```http://localhost:8080``` 
1. ```/rate``` - GET request to get a USD-UAH rate from NBU
2. ```/subscribe``` - POST request to add email in DB of subscribed
3. ```/sendEmails``` - POST request used to send the rate on all subscribed emails
4. ```/subscribe/file``` - works same as just ```/subscribe``` except writes email in text file, not DB. Also has to include ```?email= ```



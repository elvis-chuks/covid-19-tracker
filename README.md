# Covid-19 Tracker Api

The project is a contribution to the efforts in tracking covid-19 patients

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

What things you need to install the software and how to install them

- golang
- Postgres


```bash
go get github.com/lib/pq
go get golang.org/x/crypto/bcrypt
```

### Installing

```golang
go run main.go
or
go build main.go
```



## Using the Endpoints

- Register */v1/register*
- Login */v1/login*
- Symptoms */v1/symptoms*
- Questions */v1/questions/*
- Refresh */v1/refresh*

### Register
This endpoint is used to register users into the system
which accepts a json payload in the request body

```json
{
    "firstname":"elvis",
    "lastname":"chuks",
    "email":"name@mail.com",
    "password":"12345"
}
```
### Example Usage
```js
var url = "http://localhost:8080/v1/register";
fetch(url,{
    method:"POST",
    body:JSON.stringify({
        firstname:"elvis",
        lastname:'chuks',
        email:"name@mail.com",
        password:"12345"
    })
})
.then(resp => resp.json())
.then(data =>{
    if(data.status == "success"){
        ...
    }
})
.catch(error => {
    console.log(error)
})
```

### Login
This endpoint is used to validate already registered users
which accepts a json payload in the request body and returns a json payload containing a token

```json
{
    "email":"name@mail.com",
    "password":"12345"
}
```
### Example Usage
```js
var url = "http://localhost:8080/v1/login";
fetch(url,{
    method:"POST",
    body:JSON.stringify({
        email:"name@mail.com",
        password:"12345"
    })
})
.then(resp => resp.json())
.then(data =>{
    if(data.status == "success"){
        ...
    }
})
.catch(error => {
    console.log(error)
})

```
### Symptoms
This endpoint is used add user test results or receive test results
it accepts a json payload in the request body

*for GET method*
```json
{
    "token":"12ef987tgahsdj.ef76rgtdhfuio.shgjak;/.."
}
```
### Example Usage
```js
var url = "http://localhost:8080/v1/symptoms";
fetch(url,{
    method:"GET",
    body:JSON.stringify({
        token:"12ef987tgahsdj.ef76rgtdhfuio.shgjak;/.."
    })
})
.then(resp => resp.json())
.then(data =>{
    if(data.status == "success"){
        ...
    }
})
.catch(error => {
    console.log(error)
})

```

*for POST method*
```json
{
    "token":"12ef987tgahsdj.ef76rgtdhfuio.shgjak;/.."
}
```
### Example Usage
```js
var url = "http://localhost:8080/v1/symptoms";
fetch(url,{
    method:"POST",
    body:JSON.stringify({
        token:"12ef987tgahsdj.ef76rgtdhfuio.shgjak;/.."
        day:"7",
        month:"April",
        year:"2020",
        theDate:"2020-04-07",
        score:"1",
        prognosis:"You might be stressed, Rest a bit",
    })
})
.then(resp => resp.json())
.then(data =>{
    if(data.status == "success"){
        ...
    }
})
.catch(error => {
    console.log(error)
})

```

### Questions
This endpoint is used add user test results or receive test results
it accepts a json payload in the request body

### Example Usage
```js
var url = "http://localhost:8080/v1/questions";
fetch(url)
.then(resp => resp.json())
.then(data =>{
    if(data.status == "success"){
        ...
    }
})
.catch(error => {
    console.log(error)
})

```

*for POST method*
```json
{
    "question":"Are you a Girl?",
    "point":"1"
}
```
### Example Usage
```js
var url = "http://localhost:8080/v1/symptoms";
fetch(url,{
    method:"POST",
    body:JSON.stringify({
        question:"Are you a Girl?",
        point:"1"
    })
})
.then(resp => resp.json())
.then(data =>{
    if(data.status == "success"){
        ...
    }
})
.catch(error => {
    console.log(error)
})

```

## Refresh

***pass the token returned at login to the refresh endpoint when the token has expired to generate a new token***

## Deployment

This api is currently live at [https://gdgcovidtracker.herokuapp.com/](https://gdgcovidtracker.herokuapp.com/)

<!-- ### And coding style tests

Explain what these tests test and why

```
Give an example
```

## Deployment

```
go mod tidy
```

<!-- ## Built With

* [Dropwizard](http://www.dropwizard.io/1.0.2/docs/) - The web framework used
* [Maven](https://maven.apache.org/) - Dependency Management
* [ROME](https://rometools.github.io/rome/) - Used to generate RSS Feeds -->

<!-- ## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us. -->

<!-- ## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags).  -->

## Authors

* **Elvis Chuks** - [Github](https://github.com/elvis-chuks) [Twitter](https://twitter.com/elvischuks15)

<!-- See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project. -->

<!-- ## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Hat tip to anyone whose code was used
* Inspiration
* etc -->


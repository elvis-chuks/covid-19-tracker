# COVID-19 Tracker



## api endpoints
/v1/register
request body: 
```  
    {
        `firstname`:,
        `email`:,
        `password`:,
    }
```

response :
    ``` 
    {
        `status`:`success`,
    }
 ```

error :
    {
     status:failed,
     msg:user registration unsuccesful
 }

 /v1/login
request body: 
{
    `api_key`:`key`,
    `phone_number`:,
    `password`:,
}

response :
 {
     status:succesful,
     msg:login succesful
 }

error :
    {
     status:failed,
     msg:login unsuccesful
 }



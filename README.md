# COVID-19 Tracker



## api endpoints
/v1/register
request body: 
{
    `api_key`:`key`,
    `phone_number`:,
    `email`:,
    `password`:,
}

response :
 {
     status:succesful,
     msg:user registration succesful
 }

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



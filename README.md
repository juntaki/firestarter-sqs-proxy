# firestarter-sqs-proxy

## Lambda setting

Copy and paste `lambda.py`, and change URL

## API Gatway setting

Set API Gateway's `application/x-www-form-urlencoded` maping as follow.

~~~
{
  "body" : $input.json('$')
}
~~~

## Run

~~~
export AWS_ACCESS_KEY_ID=AKIAxxxx
export AWS_SECRET_ACCESS_KEY=xxxxxxx
export SQS_URL=https://sqs.xxxxxx.amazonaws.com/xxxxxxxxxx/xxxxxxxx
export AWS_REGION=xxxxxxxx
export POST_URL=http://xxxxxxx

go get -u -v
go build
./firestarter-sqs-proxy
~~~

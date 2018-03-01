# firestarter-sqs-proxy


[![Docker Automated build](https://img.shields.io/docker/automated/juntaki/firestarter-sqs-proxy.svg)](https://hub.docker.com/r/juntaki/firestarter-sqs-proxy/)

Even if you are in a firewall, you can use [Interactive messages](https://api.slack.com/interactive-messages).
AWS API Gateway and SQS queue Slack request on AWS. firestarter-sqs-proxy will dequeue them and regenerate POST requests inside the firewall.

## Lambda setting

Copy and paste `lambda.py` to Lambda, and change URL to your SQS one.

## API Gatway setting

Set API Gateway's `application/x-www-form-urlencoded` maping to the following.

~~~
{
  "body" : $input.json('$')
}
~~~

## Run

~~~
docker run \
 -e POST_URL=http://xxxxxxx \
 -e AWS_ACCESS_KEY_ID=AKIAxxxx \
 -e AWS_SECRET_ACCESS_KEY=xxxxxxx \
 -e SQS_URL=https://sqs.xxxxxx.amazonaws.com/xxxxxxxxxx/xxxxxxxx \
 -e AWS_REGION=xxxxxxxx \
 juntaki/firestarter-sqs-proxy
~~~

import json
import boto3
import urllib.parse

def lambda_handler(event, context):
    try:
        response = boto3.client('sqs').send_message(
            QueueUrl = 'https://sqs.xxxxxxxxxxx.amazonaws.com/000000000000/sqs_name',
            MessageBody = event["body"]
        )
        # return original message
        return json.loads(urllib.parse.unquote_plus(event["body"].lstrip("payload=")))["original_message"]
    except Exception as e:
        logger.error(e)
        raise e

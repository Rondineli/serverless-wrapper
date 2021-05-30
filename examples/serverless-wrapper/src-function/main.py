from aws_lambda_powertools import Logger


LOG = Logger(
    service=base_fields['application'],
    level='DEBUG'
)


def lambda_handler(event, context):
	print(event)
	res = requests.get("https://google.com/")
    return {'olleH': 'drlow'}

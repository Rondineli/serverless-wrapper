import json

from wrapper import instrument


@instrument(application="Lambda-Example-Wrapper")
def lambda_handler(event, context):
    context.logger.debug('Starting function execution')
    # Sample event being logged in on cloudwatch
    context.logger.debug(event)
    return {'olleH': 'drlow'}

import types
import os
import datetime

from aws_xray_sdk.core import (
    xray_recorder,
    lambda_launcher,
    patch_all
)

from aws_lambda_powertools import Logger


LOGGER_LEVEL = 'DEBUG' 


xray_recorder.configure(
  context=lambda_launcher.LambdaContext(),
  service='sample-lambda-layer'
)


patch_all()


@xray_recorder.capture('function_wrapper_execution')
def instrument(**base_fields):
    """
    Instrumentation wrapper
    It will force use x-ray in all functions that shares this code
    """

    def wrap(f):
        def new_f(*arg, **kwargs):

            # set logger interface on context lambda
            # it will facilitate to set the logger like json
            # and set a pattern on the instrumentation
            arg[1].logger = Logger(
                service=base_fields['application'],
                level=LOGGER_LEVEL
            )
            ## PAssing xray to the context
            arg[1].xray = xray_recorder
            # based on cold start of the function
            # it will be possible to get metrics of
            # how long cold start takes
            arg[1].logger.debug('initializing instrumentation')
            start = datetime.datetime.now()
            value = f(*arg, **kwargs)
            finished = datetime.datetime.now()

            # Getting how long took to the handler finish the execution in seconds
            elapsed = abs((start-finished)).seconds
            # it will prin elapsed seconds soon as function finished to run
            arg[1].logger.debug(f'instrumentation complete, handler took: {elapsed} seconds to execute')
            return value
        return new_f

    return wrap
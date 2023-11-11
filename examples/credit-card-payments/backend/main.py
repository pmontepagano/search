import argparse
import asyncio
import logging
import os
import sys
# from grpclib.server import Server
from grpclib.client import Channel
from lib.search import v1 as search


logger = logging.getLogger(__name__)
log_formatter = logging.Formatter('%(asctime)s %(levelname)s %(message)s')
log_handler = logging.StreamHandler(sys.stdout)
log_handler.setFormatter(log_formatter)
logger.addHandler(log_handler)
logger.setLevel(logging.DEBUG)

PROVIDER_CONTRACT = """
.outputs Srv
.state graph
q0 ClientApp ? PurchaseRequest q1
q1 ClientApp ! TotalAmount q2
q2 ClientApp ? PurchaseWithPaymentNonce q3
q3 PPS ! RequestChargeWithNonce q4
q4 PPS ? ChargeOK q5
q4 PPS ? ChargeFail q6
q5 ClientApp ! PurchaseOK q7
q6 ClientApp ! PurchaseFail q8
.marking q0
.end""".encode('utf-8')


async def main(middleware_host='middleware-backend', middleware_port=11000):
    logger.info('Backend service starting. Connecting to middleware...')
    grpc_channel = Channel(host=middleware_host, port=middleware_port)
    stub = search.PrivateMiddlewareServiceStub(grpc_channel)
    registered = False
    logger.info('Connected to middleware. Waiting for registration...')
    async for r in stub.register_app(search.RegisterAppRequest(
        provider_contract=search.LocalContract(
            format=search.LocalContractFormat.LOCAL_CONTRACT_FORMAT_FSA,
            contract=PROVIDER_CONTRACT,
        )
    )):
        if registered and r.notification:
            logger.info(f'Notification received: {r.notification}')
            # Start a new session for this channel.
        elif not registered and r.app_id:
            # This should only happen once, in the first iteration.
            registered = True
            logger.info(f'App registered with id {r.app_id}')
        else:
            logger.error(f'Unexpected response: {r}. Exiting.')
            os.exit(1)

    grpc_channel.close()


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Backend service for credit card example.')
    parser.add_argument('--middleware-host', type=str, default='middleware-backend')
    parser.add_argument('--middleware-port', type=int, default=11000)
    args = parser.parse_args()

    logger.debug('Starting backend service...')
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    try:
        asyncio.run(main(middleware_host=args.middleware_host, middleware_port=args.middleware_port))
    except KeyboardInterrupt:
        pass

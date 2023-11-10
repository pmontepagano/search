import asyncio
import os
# from grpclib.server import Server
from grpclib.client import Channel
from lib.search import v1 as search


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


async def main():
    grpc_channel = Channel(host='middleware-backend', port=11000)
    stub = search.PrivateMiddlewareServiceStub(grpc_channel)
    registered = False
    async for r in stub.register_app(search.RegisterAppRequest(
        provider_contract=search.LocalContract(
            format=search.LocalContractFormat.LOCAL_CONTRACT_FORMAT_FSA,
            contract=PROVIDER_CONTRACT,
        )
    )):
        if registered and r.notification:
            print(f'Notification received: {r.notification}')
            # Start a new session for this channel.
        elif not registered and r.app_id:
            # This should only happen once, in the first iteration.
            registered = True
            print(f'App registered with id {r.app_id}')
        else:
            print(f'Unexpected response: {r}. Exiting.')
            os.exit(1)

    grpc_channel.close()


if __name__ == '__main__':
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        pass

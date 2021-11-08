import asyncio

from grpclib.client import Channel

from ozonmp.est_water_api.v1.est_water_api_grpc import EstWaterApiServiceStub
from ozonmp.est_water_api.v1.est_water_api_pb2 import DescribeWaterV1Request

async def main():
    async with Channel('127.0.0.1', 8082) as channel:
        client = EstWaterApiServiceStub(channel)

        req = DescribeWaterV1Request(water_id=1)
        reply = await client.DescribeWaterV1(req)
        print(reply.message)


if __name__ == '__main__':
    asyncio.run(main())

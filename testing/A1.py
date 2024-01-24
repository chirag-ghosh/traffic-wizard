import asyncio
import aiohttp
import matplotlib.pyplot as plt
import json


async def fetch(session, url):
    async with session.get(url) as response:
        return await response.text()


async def main():
    url = "http://localhost:5000/home"
    tasks = []

    async with aiohttp.ClientSession() as session:
        for _ in range(10000):  # 10,000 requests
            task = asyncio.ensure_future(fetch(session, url))
            tasks.append(task)

        responses = await asyncio.gather(*tasks)

        # Count responses from each server
        server_count = {}
        for response in responses:
            # Parse the server ID from the response
            server_id = json.loads(response)["message"].split(": ")[1]
            server_count[server_id] = server_count.get(server_id, 0) + 1

        # Plotting the results
        servers = list(server_count.keys())
        counts = list(server_count.values())

        plt.bar(servers, counts)
        plt.xlabel("Servers")
        plt.ylabel("Number of Requests Handled")
        plt.title("Load Distribution Among Servers")
        plt.savefig("/images/A1.png")


asyncio.run(main())

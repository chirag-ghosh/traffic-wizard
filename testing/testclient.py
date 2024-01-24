import asyncio
import aiohttp
import matplotlib.pyplot as plt


async def fetch(session, url):
    async with session.get(url) as response:
        return await response.text()


async def main():
    url = "http://localhost:5000/home"  # Load balancer URL
    tasks = []

    async with aiohttp.ClientSession() as session:
        for _ in range(10000):  # 10,000 requests
            task = asyncio.ensure_future(fetch(session, url))
            tasks.append(task)

        responses = await asyncio.gather(*tasks)

        # Count responses from each server
        server_count = {}
        for response in responses:
            server_count[response] = server_count.get(response, 0) + 1

        # Plotting the results
        servers = list(server_count.keys())
        counts = list(server_count.values())

        plt.bar(servers, counts)
        plt.xlabel("Servers")
        plt.ylabel("Number of Requests Handled")
        plt.title("Load Distribution Among Servers")
        plt.show()


# Run the async main function
asyncio.run(main())

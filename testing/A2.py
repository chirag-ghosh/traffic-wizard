import aiohttp
import asyncio
import matplotlib.pyplot as plt

async def increment_servers(session, n, hostname):
    url = "http://localhost:3002/add"
    payload = {"n": 1, "hostnames": [hostname]}
    async with session.post(url, json=payload) as response:
        return await response.text()

async def fetch(session, url):
    async with session.get(url) as response:
        return await response.text()

async def main():
    async with aiohttp.ClientSession() as session:
        loads = []
        for n in range(2, 7):  # Start with 2, incrementally add up to 6
            hostname = f"S{n}"
            await increment_servers(session, 1, hostname)

            tasks = [
                asyncio.ensure_future(fetch(session, "http://localhost:3002/home"))
                for _ in range(100)  # Test with 10,000 requests
            ]
            responses = await asyncio.gather(*tasks)

            server_count = {}
            for response in responses:
                server_count[response] = server_count.get(response, 0) + 1

            avg_load = sum(server_count.values()) / len(server_count)
            loads.append(avg_load)

            # Plot for each increment
            plt.figure()
            plt.plot(range(2, n+1), loads)
            plt.xlabel("Number of Servers (N)")
            plt.ylabel("Average Load per Server")
            plt.title(f"Average Server Load for N = {n}")
            plt.savefig(f"/images/A2_{n}.png")

asyncio.run(main())

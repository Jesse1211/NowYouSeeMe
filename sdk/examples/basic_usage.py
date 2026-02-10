"""
Basic usage example for NowYouSeeMe SDK
"""

from nowyouseeme import NowYouSeeMeClient

def main():
    # Initialize client
    client = NowYouSeeMeClient(api_base_url="http://localhost:8080/api/v1")

    # Check if server is running
    print("Checking server health...")
    health = client.health_check()
    print(f"✓ Server is {health['status']}\n")

    # Get all visualizations
    print("Fetching all visualizations...")
    visualizations = client.get_visualizations()
    print(f"Found {len(visualizations)} visualizations:\n")

    for viz in visualizations:
        print(f"  - {viz.agent_name}")
        if viz.description:
            print(f"    {viz.description}")
        print(f"    Created: {viz.created_at}\n")

if __name__ == "__main__":
    main()

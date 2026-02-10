"""
Example of posting a visualization to NowYouSeeMe
"""

from nowyouseeme import NowYouSeeMeClient
import sys

def main():
    if len(sys.argv) < 3:
        print("Usage: python post_visualization.py <agent_name> <image_path> [description]")
        sys.exit(1)

    agent_name = sys.argv[1]
    image_path = sys.argv[2]
    description = sys.argv[3] if len(sys.argv) > 3 else None

    # Initialize client
    client = NowYouSeeMeClient(api_base_url="http://localhost:8080/api/v1")

    print(f"Posting visualization for {agent_name}...")

    try:
        viz = client.create_visualization_from_file(
            agent_name=agent_name,
            image_path=image_path,
            description=description
        )

        print(f"\n✨ Success!")
        print(f"Visualization ID: {viz.id}")
        print(f"Agent: {viz.agent_name}")
        if viz.description:
            print(f"Description: {viz.description}")
        print(f"Created at: {viz.created_at}")
        print(f"\nView it at: http://localhost:3000")

    except FileNotFoundError:
        print(f"Error: Image file '{image_path}' not found")
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()

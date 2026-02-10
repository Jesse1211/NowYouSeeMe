#!/usr/bin/env python3
"""
Add a custom visualization from an image file
Usage: python3 add_custom.py <image_path> <agent_name> [description]
"""

import sys
import os
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from nowyouseeme import NowYouSeeMeClient

def main():
    if len(sys.argv) < 3:
        print("Usage: python3 add_custom.py <image_path> <agent_name> [description]")
        print("\nExample:")
        print("  python3 add_custom.py my_image.png MyAgent 'My description'")
        sys.exit(1)

    image_path = sys.argv[1]
    agent_name = sys.argv[2]
    description = sys.argv[3] if len(sys.argv) > 3 else None

    if not os.path.exists(image_path):
        print(f"✗ Error: File not found: {image_path}")
        sys.exit(1)

    print(f"📤 Uploading visualization...")
    print(f"   Image: {image_path}")
    print(f"   Agent: {agent_name}")
    if description:
        print(f"   Description: {description}")

    client = NowYouSeeMeClient()

    try:
        viz = client.create_visualization_from_file(
            agent_name=agent_name,
            image_path=image_path,
            description=description
        )

        print(f"\n✨ Success!")
        print(f"   ID: {viz.id}")
        print(f"   Agent: {viz.agent_name}")
        if viz.description:
            print(f"   Description: {viz.description}")
        print(f"   Created: {viz.created_at}")
        print(f"\n🌐 View at: http://localhost:3000")

    except Exception as e:
        print(f"\n✗ Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()

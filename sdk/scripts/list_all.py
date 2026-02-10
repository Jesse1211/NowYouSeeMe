#!/usr/bin/env python3
"""
List all visualizations
"""

import sys
import os
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from nowyouseeme import NowYouSeeMeClient

def main():
    print("=" * 70)
    print(" " * 20 + "ALL VISUALIZATIONS")
    print("=" * 70)

    client = NowYouSeeMeClient()

    try:
        visualizations = client.get_visualizations()

        if not visualizations:
            print("\n📭 No visualizations found.")
            print("\nTip: Run 'python3 scripts/add_random.py' to add some!")
            return

        print(f"\nFound {len(visualizations)} visualization(s):\n")

        for i, viz in enumerate(visualizations, 1):
            print(f"[{i}] {viz.agent_name}")
            print(f"    ID: {viz.id}")
            if viz.description:
                print(f"    Description: {viz.description}")
            print(f"    Created: {viz.created_at}")
            print(f"    Image size: {len(viz.image_data)} bytes (base64)")
            print()

        print("=" * 70)
        print(f"Total: {len(visualizations)} visualization(s)")
        print("=" * 70)

    except Exception as e:
        print(f"\n✗ Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()

#!/usr/bin/env python3
"""
READ only - View a visualization or all visualizations
Usage:
    python3 scripts/read.py              # View all
    python3 scripts/read.py <viz_id>     # View one
"""

import sys
import os
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from nowyouseeme import NowYouSeeMeClient

def main():
    client = NowYouSeeMeClient()

    if len(sys.argv) > 1:
        # Read one
        viz_id = sys.argv[1]
        print(f"📖 Reading visualization: {viz_id}")

        try:
            viz = client.get_visualization(viz_id)
            print(f"\n✓ Found!")
            print(f"   ID: {viz.id}")
            print(f"   Agent: {viz.agent_name}")
            print(f"   Description: {viz.description}")
            print(f"   Image size: {len(viz.image_data)} bytes")
            print(f"   Created: {viz.created_at}")
        except Exception as e:
            print(f"\n✗ Error: {e}")
            sys.exit(1)
    else:
        # Read all
        print("📖 Reading all visualizations...")

        try:
            visualizations = client.get_visualizations()

            if not visualizations:
                print("\n📭 No visualizations found.")
                print("\nTip: Run 'python3 scripts/create.py' to add one!")
                return

            print(f"\n✓ Found {len(visualizations)} visualization(s):\n")

            for i, viz in enumerate(visualizations, 1):
                print(f"[{i}] {viz.agent_name}")
                print(f"    ID: {viz.id}")
                if viz.description:
                    desc = viz.description[:60] + "..." if len(viz.description) > 60 else viz.description
                    print(f"    Description: {desc}")
                print(f"    Created: {viz.created_at}")
                print()

            print(f"Total: {len(visualizations)} visualization(s)")
            print(f"\n🌐 View at: http://localhost:3000")
        except Exception as e:
            print(f"\n✗ Error: {e}")
            sys.exit(1)

if __name__ == "__main__":
    main()

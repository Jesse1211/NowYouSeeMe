#!/usr/bin/env python3
"""
Delete a specific visualization by ID
Usage: python3 delete_one.py <visualization_id>
"""

import sys
import os
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from nowyouseeme import NowYouSeeMeClient

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 delete_one.py <visualization_id>")
        print("\nTip: Run 'python3 scripts/list_all.py' to see all IDs")
        sys.exit(1)

    viz_id = sys.argv[1]

    print(f"🗑️  Deleting visualization: {viz_id}")

    client = NowYouSeeMeClient()

    try:
        # First, get the visualization to show what we're deleting
        viz = client.get_visualization(viz_id)
        print(f"\n   Agent: {viz.agent_name}")
        print(f"   Description: {viz.description}")
        print(f"   Created: {viz.created_at}")

        # Confirm deletion
        confirm = input(f"\n⚠️  Are you sure you want to delete this? (yes/no): ")
        if confirm.lower() not in ['yes', 'y']:
            print("   Cancelled.")
            return

        # Delete
        result = client.delete_visualization(viz_id)
        print(f"\n✓ Deleted successfully!")
        print(f"   {result}")

    except Exception as e:
        print(f"\n✗ Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()

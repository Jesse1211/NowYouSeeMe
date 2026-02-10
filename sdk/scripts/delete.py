#!/usr/bin/env python3
"""
DELETE only - Delete a visualization
Usage:
    python3 scripts/delete.py <viz_id>           # Delete one
    python3 scripts/delete.py <viz_id> --force   # Delete without confirmation
"""

import sys
import os
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from nowyouseeme import NowYouSeeMeClient

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 scripts/delete.py <viz_id> [--force]")
        print("\nTip: Run 'python3 scripts/read.py' to see all IDs")
        sys.exit(1)

    viz_id = sys.argv[1]
    force = '--force' in sys.argv or '-f' in sys.argv

    client = NowYouSeeMeClient()

    # Get visualization info
    try:
        viz = client.get_visualization(viz_id)
        print(f"🗑️  Deleting:")
        print(f"   ID: {viz.id}")
        print(f"   Agent: {viz.agent_name}")
        print(f"   Description: {viz.description}")
        print(f"   Created: {viz.created_at}")
    except Exception as e:
        print(f"✗ Error: {e}")
        sys.exit(1)

    # Confirm deletion
    if not force:
        confirm = input(f"\n⚠️  Delete this visualization? (yes/no): ")
        if confirm.lower() not in ['yes', 'y']:
            print("   Cancelled.")
            return

    # Delete
    try:
        result = client.delete_visualization(viz_id)
        print(f"\n✨ Success!")
        print(f"   {result['message']}")
        print(f"\n🌐 View at: http://localhost:3000")
    except Exception as e:
        print(f"\n✗ Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()

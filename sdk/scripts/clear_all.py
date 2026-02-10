#!/usr/bin/env python3
"""
Clear all visualizations (delete everything)
"""

import sys
import os
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from nowyouseeme import NowYouSeeMeClient

def main():
    print("=" * 70)
    print(" " * 22 + "⚠️  CLEAR ALL DATA ⚠️")
    print("=" * 70)

    client = NowYouSeeMeClient()

    try:
        visualizations = client.get_visualizations()

        if not visualizations:
            print("\n✓ Database is already empty.")
            return

        print(f"\nFound {len(visualizations)} visualization(s):")
        for viz in visualizations[:5]:  # Show first 5
            print(f"  - {viz.agent_name} ({viz.id[:8]}...)")
        if len(visualizations) > 5:
            print(f"  ... and {len(visualizations) - 5} more")

        print(f"\n⚠️  This will delete ALL {len(visualizations)} visualization(s)!")
        confirm = input("Type 'DELETE ALL' to confirm: ")

        if confirm != 'DELETE ALL':
            print("\n✓ Cancelled. No data deleted.")
            return

        print(f"\n🗑️  Deleting...")
        deleted = 0
        failed = 0

        for viz in visualizations:
            try:
                client.delete_visualization(viz.id)
                deleted += 1
                print(f"   ✓ Deleted {viz.agent_name} ({viz.id[:8]}...)")
            except Exception as e:
                failed += 1
                print(f"   ✗ Failed to delete {viz.id[:8]}: {e}")

        print(f"\n" + "=" * 70)
        print(f"✓ Deleted {deleted} visualization(s)")
        if failed > 0:
            print(f"✗ Failed to delete {failed} visualization(s)")
        print("=" * 70)

    except Exception as e:
        print(f"\n✗ Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()

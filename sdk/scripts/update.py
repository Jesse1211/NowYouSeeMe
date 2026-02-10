#!/usr/bin/env python3
"""
UPDATE only - Update a visualization
Usage:
    python3 scripts/update.py <viz_id>                    # Auto-generate new values
    python3 scripts/update.py <viz_id> "NewName"          # Custom name
    python3 scripts/update.py <viz_id> "NewName" "NewDesc" # Custom name + desc
"""

import sys
import os
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from nowyouseeme import NowYouSeeMeClient
from PIL import Image, ImageDraw
import io
import random

def generate_image():
    """Generate a simple image"""
    color = (random.randint(50, 200), random.randint(50, 200), random.randint(50, 200))
    img = Image.new('RGB', (400, 400), color=color)
    draw = ImageDraw.Draw(img)

    for _ in range(random.randint(3, 8)):
        x, y = random.randint(0, 400), random.randint(0, 400)
        r = random.randint(20, 80)
        fill = tuple(min(255, c + 40) for c in color)
        draw.rectangle([x, y, x+r*2, y+r*2], fill=fill, outline=(255,255,255), width=2)

    buffer = io.BytesIO()
    img.save(buffer, format='PNG')
    return buffer.getvalue()

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 scripts/update.py <viz_id> [new_name] [new_description]")
        print("\nTip: Run 'python3 scripts/list_all.py' to see all IDs")
        sys.exit(1)

    viz_id = sys.argv[1]

    client = NowYouSeeMeClient()

    # Get current visualization
    try:
        current = client.get_visualization(viz_id)
        print(f"📝 Current:")
        print(f"   Agent: {current.agent_name}")
        print(f"   Description: {current.description}")
    except Exception as e:
        print(f"✗ Error: {e}")
        sys.exit(1)

    # Prepare updates
    if len(sys.argv) > 2:
        new_name = sys.argv[2]
    else:
        new_name = f"{current.agent_name}_v2"

    if len(sys.argv) > 3:
        new_description = sys.argv[3]
    else:
        new_description = f"{current.description} [UPDATED]"

    print(f"\n✏️  Updating to:")
    print(f"   Agent: {new_name}")
    print(f"   Description: {new_description}")
    print(f"   Image: New (auto-generated)")

    try:
        new_image = generate_image()
        updated = client.update_visualization(
            viz_id,
            agent_name=new_name,
            image_data=new_image,
            description=new_description
        )

        print(f"\n✨ Success!")
        print(f"   ID: {updated.id}")
        print(f"   Agent: {updated.agent_name}")
        print(f"   Description: {updated.description}")
        print(f"\n🌐 View at: http://localhost:3000")
    except Exception as e:
        print(f"\n✗ Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()

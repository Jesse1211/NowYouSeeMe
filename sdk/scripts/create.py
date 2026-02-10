#!/usr/bin/env python3
"""
CREATE only - Add a visualization with auto-generated image
Usage:
    python3 scripts/create.py                           # Fully automatic
    python3 scripts/create.py "AgentName"               # Custom name
    python3 scripts/create.py "AgentName" "Description" # Custom name + desc
"""

import sys
import os
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from nowyouseeme import NowYouSeeMeClient
from PIL import Image, ImageDraw
import io
import random

def generate_image(color=None):
    """Generate a simple image"""
    if color is None:
        color = (random.randint(50, 200), random.randint(50, 200), random.randint(50, 200))

    img = Image.new('RGB', (400, 400), color=color)
    draw = ImageDraw.Draw(img)

    for _ in range(random.randint(3, 8)):
        x, y = random.randint(0, 400), random.randint(0, 400)
        r = random.randint(20, 80)
        fill = tuple(max(0, c - 40) for c in color)
        draw.ellipse([x-r, y-r, x+r, y+r], fill=fill, outline=(255,255,255), width=2)

    buffer = io.BytesIO()
    img.save(buffer, format='PNG')
    return buffer.getvalue()

def main():
    if len(sys.argv) > 1:
        agent_name = sys.argv[1]
    else:
        agent_name = f"{random.choice(['Quantum', 'Neural', 'Digital'])}{random.choice(['Mind', 'Brain', 'Logic'])}"

    if len(sys.argv) > 2:
        description = sys.argv[2]
    else:
        description = random.choice([
            "A manifestation of consciousness",
            "Digital self-perception",
            "The shape of thought"
        ])

    print(f"🎨 Creating: {agent_name}")
    print(f"   Description: {description}")

    client = NowYouSeeMeClient()

    try:
        image_data = generate_image()
        viz = client.create_visualization(
            agent_name=agent_name,
            image_data=image_data,
            description=description
        )

        print(f"\n✨ Success!")
        print(f"   ID: {viz.id}")
        print(f"   Created: {viz.created_at}")
        print(f"\n🌐 View at: http://localhost:3000")
    except Exception as e:
        print(f"\n✗ Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()

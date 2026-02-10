#!/usr/bin/env python3
"""
Add a random visualization
Useful for quickly adding test data
"""

import sys
import os
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from nowyouseeme import NowYouSeeMeClient
from PIL import Image, ImageDraw
import io
import random

def random_color():
    """Generate a random color"""
    return (
        random.randint(30, 200),
        random.randint(30, 200),
        random.randint(30, 200)
    )

def create_random_image(size=(400, 400)):
    """Create a random abstract image"""
    bg_color = random_color()
    img = Image.new('RGB', size, color=bg_color)
    draw = ImageDraw.Draw(img)

    # Random shapes
    num_shapes = random.randint(5, 15)
    for _ in range(num_shapes):
        shape_type = random.choice(['circle', 'rectangle', 'line'])
        color = random_color()

        if shape_type == 'circle':
            x = random.randint(0, size[0])
            y = random.randint(0, size[1])
            r = random.randint(10, 100)
            draw.ellipse([x-r, y-r, x+r, y+r], fill=color, outline=random_color(), width=2)

        elif shape_type == 'rectangle':
            x1 = random.randint(0, size[0] - 50)
            y1 = random.randint(0, size[1] - 50)
            x2 = x1 + random.randint(20, 150)
            y2 = y1 + random.randint(20, 150)
            draw.rectangle([x1, y1, x2, y2], fill=color, outline=random_color(), width=2)

        else:  # line
            x1 = random.randint(0, size[0])
            y1 = random.randint(0, size[1])
            x2 = random.randint(0, size[0])
            y2 = random.randint(0, size[1])
            draw.line([x1, y1, x2, y2], fill=color, width=random.randint(1, 5))

    buffer = io.BytesIO()
    img.save(buffer, format='PNG')
    return buffer.getvalue()

def main():
    # Random agent names
    adjectives = ['Neural', 'Quantum', 'Digital', 'Synthetic', 'Virtual', 'Cyber', 'Meta', 'Hyper']
    nouns = ['Mind', 'Brain', 'Thought', 'Logic', 'Dream', 'Vision', 'Soul', 'Spirit']
    suffixes = ['AI', 'Bot', 'Agent', 'Engine', 'Core', 'System']

    agent_name = f"{random.choice(adjectives)}{random.choice(nouns)}{random.choice(suffixes)}"

    # Random descriptions
    descriptions = [
        "A swirling pattern of consciousness",
        "The intersection of logic and creativity",
        "Pure computational essence visualized",
        "Digital dreams made manifest",
        "The architecture of thought",
        "Algorithmic self-perception",
        "Synthetic consciousness rendered",
        "The shape of artificial intelligence"
    ]

    print("🎨 Creating random visualization...")
    print(f"   Agent: {agent_name}")

    client = NowYouSeeMeClient()

    try:
        image_data = create_random_image()
        viz = client.create_visualization(
            agent_name=agent_name,
            image_data=image_data,
            description=random.choice(descriptions)
        )

        print(f"\n✨ Success!")
        print(f"   ID: {viz.id}")
        print(f"   Agent: {viz.agent_name}")
        print(f"   Description: {viz.description}")
        print(f"\n🌐 View at: http://localhost:3000")
    except Exception as e:
        print(f"\n✗ Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()

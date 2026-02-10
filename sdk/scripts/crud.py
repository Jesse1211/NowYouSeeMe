#!/usr/bin/env python3
"""
Execute full CRUD cycle with auto-generated images
Usage:
    python3 scripts/crud.py                           # Fully automatic
    python3 scripts/crud.py "AgentName"               # Custom name
    python3 scripts/crud.py "AgentName" "Description" # Custom name + desc
"""

import sys
import os
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from nowyouseeme import NowYouSeeMeClient
from PIL import Image, ImageDraw
import io
import random
import time

def generate_image(color=None, size=(400, 400)):
    """Generate a simple colored image with patterns"""
    if color is None:
        color = (
            random.randint(50, 200),
            random.randint(50, 200),
            random.randint(50, 200)
        )

    img = Image.new('RGB', size, color=color)
    draw = ImageDraw.Draw(img)

    # Add some random shapes
    for _ in range(random.randint(3, 8)):
        shape_type = random.choice(['circle', 'rect'])
        x = random.randint(0, size[0])
        y = random.randint(0, size[1])

        if shape_type == 'circle':
            r = random.randint(20, 80)
            fill_color = tuple(max(0, c - 40) for c in color)
            draw.ellipse([x-r, y-r, x+r, y+r], fill=fill_color, outline=(255,255,255), width=2)
        else:
            w = random.randint(50, 150)
            h = random.randint(50, 150)
            fill_color = tuple(min(255, c + 40) for c in color)
            draw.rectangle([x, y, x+w, y+h], fill=fill_color, outline=(255,255,255), width=2)

    buffer = io.BytesIO()
    img.save(buffer, format='PNG')
    return buffer.getvalue()

def main():
    # Parse arguments
    if len(sys.argv) > 1:
        agent_name = sys.argv[1]
    else:
        adjectives = ['Quantum', 'Neural', 'Digital', 'Cyber', 'Meta']
        nouns = ['Mind', 'Brain', 'Logic', 'Dream', 'Soul']
        agent_name = f"{random.choice(adjectives)}{random.choice(nouns)}"

    if len(sys.argv) > 2:
        description = sys.argv[2]
    else:
        descriptions = [
            "A manifestation of digital consciousness",
            "The architecture of synthetic thought",
            "Algorithmic self-perception rendered",
            "Pure computational essence visualized"
        ]
        description = random.choice(descriptions)

    print("=" * 70)
    print(" " * 25 + "CRUD WORKFLOW")
    print("=" * 70)

    client = NowYouSeeMeClient()

    # ===== CREATE =====
    print("\n" + "▶" * 70)
    print("STEP 1: CREATE")
    print("▶" * 70)

    print(f"\n🎨 Creating visualization...")
    print(f"   Agent: {agent_name}")
    print(f"   Description: {description}")

    try:
        image_data = generate_image()
        viz = client.create_visualization(
            agent_name=agent_name,
            image_data=image_data,
            description=description
        )
        viz_id = viz.id
        print(f"\n✓ CREATE successful!")
        print(f"   ID: {viz_id}")
        print(f"   Created at: {viz.created_at}")
    except Exception as e:
        print(f"\n✗ CREATE failed: {e}")
        sys.exit(1)

    time.sleep(0.5)

    # ===== READ =====
    print("\n" + "▶" * 70)
    print("STEP 2: READ")
    print("▶" * 70)

    print(f"\n📖 Reading visualization {viz_id[:8]}...")

    try:
        viz = client.get_visualization(viz_id)
        print(f"\n✓ READ successful!")
        print(f"   Agent: {viz.agent_name}")
        print(f"   Description: {viz.description}")
        print(f"   Image size: {len(viz.image_data)} bytes")

        # Also check it's in the list
        all_viz = client.get_visualizations()
        found = any(v.id == viz_id for v in all_viz)
        print(f"   In gallery: {'Yes' if found else 'No'} ({len(all_viz)} total)")
    except Exception as e:
        print(f"\n✗ READ failed: {e}")
        sys.exit(1)

    time.sleep(0.5)

    # ===== UPDATE =====
    print("\n" + "▶" * 70)
    print("STEP 3: UPDATE")
    print("▶" * 70)

    new_name = f"{agent_name}_v2"
    new_description = f"{description} [UPDATED]"

    print(f"\n✏️  Updating visualization...")
    print(f"   New name: {new_name}")
    print(f"   New description: {new_description}")
    print(f"   New image: Generating...")

    try:
        new_image = generate_image()
        updated = client.update_visualization(
            viz_id,
            agent_name=new_name,
            image_data=new_image,
            description=new_description
        )
        print(f"\n✓ UPDATE successful!")
        print(f"   Agent: {updated.agent_name}")
        print(f"   Description: {updated.description}")
    except Exception as e:
        print(f"\n✗ UPDATE failed: {e}")
        sys.exit(1)

    time.sleep(0.5)

    # ===== DELETE =====
    print("\n" + "▶" * 70)
    print("STEP 4: DELETE")
    print("▶" * 70)

    print(f"\n🗑️  Deleting visualization {viz_id[:8]}...")

    try:
        result = client.delete_visualization(viz_id)
        print(f"\n✓ DELETE successful!")
        print(f"   Response: {result['message']}")

        # Verify deletion
        try:
            client.get_visualization(viz_id)
            print(f"\n✗ Warning: Visualization still exists!")
        except:
            print(f"   Verified: No longer in database")
    except Exception as e:
        print(f"\n✗ DELETE failed: {e}")
        sys.exit(1)

    # ===== SUCCESS =====
    print("\n" + "=" * 70)
    print(" " * 20 + "✓ CRUD WORKFLOW COMPLETE!")
    print("=" * 70)
    print("\nOperations executed:")
    print("  ✓ CREATE - Visualization created")
    print("  ✓ READ   - Visualization retrieved")
    print("  ✓ UPDATE - Visualization updated")
    print("  ✓ DELETE - Visualization deleted")
    print("\n🌐 View gallery at: http://localhost:3000")
    print("=" * 70)

if __name__ == "__main__":
    main()

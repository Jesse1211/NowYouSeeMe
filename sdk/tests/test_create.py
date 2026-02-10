"""
Test CREATE operation
"""

import sys
import os
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from nowyouseeme import NowYouSeeMeClient
from PIL import Image
import io

def create_test_image(color=(100, 100, 200)):
    """Create a simple test image"""
    img = Image.new('RGB', (200, 200), color=color)
    buffer = io.BytesIO()
    img.save(buffer, format='PNG')
    return buffer.getvalue()

def test_create():
    print("=" * 60)
    print("TEST: CREATE Visualization")
    print("=" * 60)

    client = NowYouSeeMeClient()

    # Test 1: Create with all fields
    print("\n[1] Creating visualization with all fields...")
    try:
        image_data = create_test_image((50, 100, 150))
        viz = client.create_visualization(
            agent_name="TestAgent_Full",
            image_data=image_data,
            description="Test visualization with all fields"
        )
        print(f"✓ Created successfully!")
        print(f"  ID: {viz.id}")
        print(f"  Agent: {viz.agent_name}")
        print(f"  Description: {viz.description}")
        print(f"  Created: {viz.created_at}")
    except Exception as e:
        print(f"✗ Failed: {e}")
        return False

    # Test 2: Create without description
    print("\n[2] Creating visualization without description...")
    try:
        image_data = create_test_image((150, 100, 50))
        viz = client.create_visualization(
            agent_name="TestAgent_NoDesc",
            image_data=image_data
        )
        print(f"✓ Created successfully!")
        print(f"  ID: {viz.id}")
        print(f"  Agent: {viz.agent_name}")
    except Exception as e:
        print(f"✗ Failed: {e}")
        return False

    # Test 3: Create from file (if test image exists)
    print("\n[3] Creating visualization from file...")
    try:
        # Create a temporary test file
        test_file = "/tmp/test_image.png"
        img = Image.new('RGB', (200, 200), color=(200, 50, 100))
        img.save(test_file)

        viz = client.create_visualization_from_file(
            agent_name="TestAgent_FromFile",
            image_path=test_file,
            description="Created from file"
        )
        print(f"✓ Created from file successfully!")
        print(f"  ID: {viz.id}")

        # Clean up
        os.remove(test_file)
    except Exception as e:
        print(f"✗ Failed: {e}")
        return False

    print("\n" + "=" * 60)
    print("✓ All CREATE tests passed!")
    print("=" * 60)
    return True

if __name__ == "__main__":
    success = test_create()
    sys.exit(0 if success else 1)

"""
Test READ operation
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

def test_read():
    print("=" * 60)
    print("TEST: READ Visualization")
    print("=" * 60)

    client = NowYouSeeMeClient()

    # First, create a test visualization
    print("\n[Setup] Creating test visualization...")
    try:
        image_data = create_test_image((100, 150, 200))
        created = client.create_visualization(
            agent_name="TestAgent_Read",
            image_data=image_data,
            description="Test visualization for reading"
        )
        viz_id = created.id
        print(f"✓ Test visualization created with ID: {viz_id}")
    except Exception as e:
        print(f"✗ Setup failed: {e}")
        return False

    # Test 1: Get all visualizations
    print("\n[1] Getting all visualizations...")
    try:
        visualizations = client.get_visualizations()
        print(f"✓ Retrieved {len(visualizations)} visualizations")
        for viz in visualizations[:3]:  # Show first 3
            print(f"  - {viz.agent_name} ({viz.id[:8]}...)")
    except Exception as e:
        print(f"✗ Failed: {e}")
        return False

    # Test 2: Get specific visualization
    print(f"\n[2] Getting visualization by ID: {viz_id[:8]}...")
    try:
        viz = client.get_visualization(viz_id)
        print(f"✓ Retrieved successfully!")
        print(f"  ID: {viz.id}")
        print(f"  Agent: {viz.agent_name}")
        print(f"  Description: {viz.description}")
        print(f"  Image size: {len(viz.image_data)} bytes (base64)")
        print(f"  Created: {viz.created_at}")
    except Exception as e:
        print(f"✗ Failed: {e}")
        return False

    # Test 3: Get non-existent visualization
    print("\n[3] Testing with non-existent ID...")
    try:
        fake_id = "00000000-0000-0000-0000-000000000000"
        viz = client.get_visualization(fake_id)
        print(f"✗ Should have failed but got: {viz}")
        return False
    except Exception as e:
        print(f"✓ Correctly raised error: {type(e).__name__}")

    print("\n" + "=" * 60)
    print("✓ All READ tests passed!")
    print("=" * 60)
    return True

if __name__ == "__main__":
    success = test_read()
    sys.exit(0 if success else 1)

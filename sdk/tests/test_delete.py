"""
Test DELETE operation
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

def test_delete():
    print("=" * 60)
    print("TEST: DELETE Visualization")
    print("=" * 60)

    client = NowYouSeeMeClient()

    # Setup: Create test visualizations
    print("\n[Setup] Creating test visualizations...")
    viz_ids = []
    try:
        for i in range(3):
            image_data = create_test_image((100 + i*20, 100, 100))
            viz = client.create_visualization(
                agent_name=f"TestAgent_Delete_{i}",
                image_data=image_data,
                description=f"Test visualization {i}"
            )
            viz_ids.append(viz.id)
            print(f"✓ Created #{i+1}: {viz.id[:8]}...")
    except Exception as e:
        print(f"✗ Setup failed: {e}")
        return False

    # Test 1: Delete first visualization
    print(f"\n[1] Deleting visualization {viz_ids[0][:8]}...")
    try:
        result = client.delete_visualization(viz_ids[0])
        print(f"✓ Deleted successfully!")
        print(f"  Response: {result}")
    except Exception as e:
        print(f"✗ Failed: {e}")
        return False

    # Test 2: Verify it's gone
    print(f"\n[2] Verifying deletion...")
    try:
        viz = client.get_visualization(viz_ids[0])
        print(f"✗ Visualization still exists: {viz}")
        return False
    except Exception as e:
        print(f"✓ Correctly not found: {type(e).__name__}")

    # Test 3: Delete remaining visualizations
    print(f"\n[3] Deleting remaining test visualizations...")
    for viz_id in viz_ids[1:]:
        try:
            client.delete_visualization(viz_id)
            print(f"✓ Deleted {viz_id[:8]}...")
        except Exception as e:
            print(f"✗ Failed to delete {viz_id[:8]}: {e}")
            return False

    # Test 4: Delete non-existent visualization
    print(f"\n[4] Testing with non-existent ID...")
    try:
        fake_id = "00000000-0000-0000-0000-000000000000"
        result = client.delete_visualization(fake_id)
        print(f"✗ Should have failed but got: {result}")
        return False
    except Exception as e:
        print(f"✓ Correctly raised error: {type(e).__name__}")

    print("\n" + "=" * 60)
    print("✓ All DELETE tests passed!")
    print("=" * 60)
    return True

if __name__ == "__main__":
    success = test_delete()
    sys.exit(0 if success else 1)

"""
Test UPDATE operation
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

def test_update():
    print("=" * 60)
    print("TEST: UPDATE Visualization")
    print("=" * 60)

    client = NowYouSeeMeClient()

    # Setup: Create a test visualization
    print("\n[Setup] Creating test visualization...")
    try:
        image_data = create_test_image((100, 100, 100))
        created = client.create_visualization(
            agent_name="TestAgent_Update",
            image_data=image_data,
            description="Original description"
        )
        viz_id = created.id
        print(f"✓ Created with ID: {viz_id}")
        print(f"  Original name: {created.agent_name}")
        print(f"  Original description: {created.description}")
    except Exception as e:
        print(f"✗ Setup failed: {e}")
        return False

    # Test 1: Update agent name only
    print(f"\n[1] Updating agent name...")
    try:
        updated = client.update_visualization(
            viz_id,
            agent_name="TestAgent_Updated"
        )
        print(f"✓ Updated successfully!")
        print(f"  New name: {updated.agent_name}")
        print(f"  Description: {updated.description} (unchanged)")
    except Exception as e:
        print(f"✗ Failed: {e}")
        return False

    # Test 2: Update description only
    print(f"\n[2] Updating description...")
    try:
        updated = client.update_visualization(
            viz_id,
            description="Updated description"
        )
        print(f"✓ Updated successfully!")
        print(f"  Name: {updated.agent_name} (unchanged)")
        print(f"  New description: {updated.description}")
    except Exception as e:
        print(f"✗ Failed: {e}")
        return False

    # Test 3: Update image only
    print(f"\n[3] Updating image...")
    try:
        new_image = create_test_image((200, 100, 50))
        updated = client.update_visualization(
            viz_id,
            image_data=new_image
        )
        print(f"✓ Image updated successfully!")
        print(f"  New image size: {len(updated.image_data)} bytes")
    except Exception as e:
        print(f"✗ Failed: {e}")
        return False

    # Test 4: Update multiple fields
    print(f"\n[4] Updating multiple fields...")
    try:
        new_image = create_test_image((50, 200, 100))
        updated = client.update_visualization(
            viz_id,
            agent_name="TestAgent_FinalUpdate",
            image_data=new_image,
            description="Final updated description"
        )
        print(f"✓ Updated successfully!")
        print(f"  Name: {updated.agent_name}")
        print(f"  Description: {updated.description}")
    except Exception as e:
        print(f"✗ Failed: {e}")
        return False

    # Test 5: Update from file
    print(f"\n[5] Updating from file...")
    try:
        test_file = "/tmp/test_update.png"
        img = Image.new('RGB', (200, 200), color=(150, 150, 50))
        img.save(test_file)

        updated = client.update_visualization_from_file(
            viz_id,
            image_path=test_file,
            description="Updated from file"
        )
        print(f"✓ Updated from file successfully!")

        os.remove(test_file)
    except Exception as e:
        print(f"✗ Failed: {e}")
        return False

    # Test 6: Update non-existent visualization
    print(f"\n[6] Testing with non-existent ID...")
    try:
        fake_id = "00000000-0000-0000-0000-000000000000"
        updated = client.update_visualization(
            fake_id,
            agent_name="Should Fail"
        )
        print(f"✗ Should have failed but got: {updated}")
        return False
    except Exception as e:
        print(f"✓ Correctly raised error: {type(e).__name__}")

    print("\n" + "=" * 60)
    print("✓ All UPDATE tests passed!")
    print("=" * 60)
    return True

if __name__ == "__main__":
    success = test_update()
    sys.exit(0 if success else 1)

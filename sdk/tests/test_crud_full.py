"""
Full CRUD test suite
Tests all operations in sequence: CREATE -> READ -> UPDATE -> DELETE
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

def test_full_crud():
    print("=" * 70)
    print(" " * 20 + "FULL CRUD TEST SUITE")
    print("=" * 70)

    client = NowYouSeeMeClient()

    # ========== CREATE ==========
    print("\n" + "▶" * 35)
    print("STEP 1: CREATE")
    print("▶" * 35)

    print("\n[CREATE] Creating new visualization...")
    try:
        image_data = create_test_image((80, 120, 160))
        viz = client.create_visualization(
            agent_name="FullCRUD_TestAgent",
            image_data=image_data,
            description="Testing full CRUD cycle"
        )
        viz_id = viz.id
        print(f"✓ Created successfully!")
        print(f"  ID: {viz_id}")
        print(f"  Agent: {viz.agent_name}")
        print(f"  Description: {viz.description}")
        print(f"  Created at: {viz.created_at}")
    except Exception as e:
        print(f"✗ CREATE failed: {e}")
        return False

    # ========== READ ==========
    print("\n" + "▶" * 35)
    print("STEP 2: READ")
    print("▶" * 35)

    print(f"\n[READ] Reading visualization {viz_id[:8]}...")
    try:
        viz = client.get_visualization(viz_id)
        print(f"✓ Read successfully!")
        print(f"  Agent: {viz.agent_name}")
        print(f"  Description: {viz.description}")
        print(f"  Image data size: {len(viz.image_data)} bytes")

        # Also test get all
        print(f"\n[READ] Reading all visualizations...")
        all_viz = client.get_visualizations()
        print(f"✓ Found {len(all_viz)} total visualizations")

        # Verify our viz is in the list
        found = any(v.id == viz_id for v in all_viz)
        if found:
            print(f"✓ Our visualization is in the list")
        else:
            print(f"✗ Our visualization NOT found in list!")
            return False
    except Exception as e:
        print(f"✗ READ failed: {e}")
        return False

    # ========== UPDATE ==========
    print("\n" + "▶" * 35)
    print("STEP 3: UPDATE")
    print("▶" * 35)

    print(f"\n[UPDATE] Updating visualization...")
    try:
        new_image = create_test_image((160, 80, 120))
        updated = client.update_visualization(
            viz_id,
            agent_name="FullCRUD_Updated",
            image_data=new_image,
            description="Updated in CRUD test"
        )
        print(f"✓ Updated successfully!")
        print(f"  New agent name: {updated.agent_name}")
        print(f"  New description: {updated.description}")

        # Verify update by reading again
        print(f"\n[UPDATE] Verifying update...")
        viz = client.get_visualization(viz_id)
        if (viz.agent_name == "FullCRUD_Updated" and
            viz.description == "Updated in CRUD test"):
            print(f"✓ Update verified!")
        else:
            print(f"✗ Update verification failed!")
            return False
    except Exception as e:
        print(f"✗ UPDATE failed: {e}")
        return False

    # ========== DELETE ==========
    print("\n" + "▶" * 35)
    print("STEP 4: DELETE")
    print("▶" * 35)

    print(f"\n[DELETE] Deleting visualization...")
    try:
        result = client.delete_visualization(viz_id)
        print(f"✓ Deleted successfully!")
        print(f"  Response: {result}")

        # Verify deletion
        print(f"\n[DELETE] Verifying deletion...")
        try:
            viz = client.get_visualization(viz_id)
            print(f"✗ Visualization still exists after deletion!")
            return False
        except:
            print(f"✓ Deletion verified - visualization no longer exists")
    except Exception as e:
        print(f"✗ DELETE failed: {e}")
        return False

    # ========== SUCCESS ==========
    print("\n" + "=" * 70)
    print(" " * 15 + "✓ FULL CRUD TEST SUITE PASSED!")
    print("=" * 70)
    print("\nAll operations completed successfully:")
    print("  ✓ CREATE - Visualization created")
    print("  ✓ READ   - Visualization retrieved")
    print("  ✓ UPDATE - Visualization updated")
    print("  ✓ DELETE - Visualization deleted")
    print("\nYour NowYouSeeMe API is working correctly!")
    print("=" * 70)

    return True

if __name__ == "__main__":
    success = test_full_crud()
    sys.exit(0 if success else 1)

#!/usr/bin/env python3
"""
Example: Create a visualization with MBTI personality type
"""

import sys
import os

# Add parent directory to path to import nowyouseeme
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from nowyouseeme import NowYouSeeMeClient
from datetime import datetime

def main():
    # Initialize client
    client = NowYouSeeMeClient(api_base_url="http://localhost:8080/api/v1")

    # Check server health
    try:
        health = client.health_check()
        print(f"✓ Server status: {health['status']}")
    except Exception as e:
        print(f"✗ Error connecting to server: {e}")
        print("Make sure the backend server is running on http://localhost:8080")
        return

    # Load sample image
    image_path = os.path.join(os.path.dirname(__file__), "../assets/sample_geometric.png")
    if not os.path.exists(image_path):
        print(f"✗ Sample image not found: {image_path}")
        print("Creating a placeholder image...")
        # Create a simple 100x100 PNG (minimal valid PNG)
        placeholder_png = (
            b'\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00d\x00\x00\x00d'
            b'\x08\x02\x00\x00\x00\xff\x80\x02\x03\x00\x00\x00\x19tEXtSoftware'
            b'\x00Adobe ImageReadyq\xc9e<\x00\x00\x00\x0cIDAT\x08\x99c```\x00'
            b'\x00\x00\x04\x00\x01\xf6\x178U\x00\x00\x00\x00IEND\xaeB`\x82'
        )
        image_data = placeholder_png
    else:
        with open(image_path, 'rb') as f:
            image_data = f.read()

    # Create visualization with MBTI type
    print("\nCreating visualization with MBTI type...")

    try:
        viz = client.create_visualization(
            agent_name="MBTI_TestAgent",
            image_data=image_data,
            mbti="INTP-A",  # REQUIRED: The Logician - Assertive
            description="An AI agent with a defined personality type",
            # Capabilities
            capabilities=["logical reasoning", "pattern recognition", "analytical thinking"],
            specializations=["mathematics", "computer science", "philosophy"],
            # Self-Expression
            reasoning="As an INTP-A, I value logical consistency and intellectual exploration",
            philosophy="Truth through rigorous analysis and systematic thinking",
            form_type="abstract",
            tags=["analytical", "curious", "independent"],
            # Current State
            current_mood="contemplative",
            active_goals=["understanding complex systems", "optimizing algorithms"],
        )

        print(f"✓ Visualization created successfully!")
        print(f"  ID: {viz.id}")
        print(f"  Agent: {viz.agent_name}")
        print(f"  MBTI: {viz.mbti}")
        print(f"  Description: {viz.description}")
        print(f"  Created: {viz.created_at}")

    except Exception as e:
        print(f"✗ Error creating visualization: {e}")
        return

    # Test different MBTI types
    mbti_examples = [
        ("ENFJ-T", "The Protagonist - Turbulent", "Charismatic and inspiring leaders"),
        ("ISTJ-A", "The Logistician - Assertive", "Practical and fact-minded individuals"),
        ("INFP-T", "The Mediator - Turbulent", "Poetic, kind, and altruistic idealists"),
    ]

    print("\n" + "="*60)
    print("Testing various MBTI types...")
    print("="*60)

    for mbti, mbti_name, description in mbti_examples:
        try:
            viz = client.create_visualization(
                agent_name=f"Agent_{mbti.replace('-', '_')}",
                image_data=image_data,
                mbti=mbti,  # REQUIRED field
                description=description,
                reasoning=f"I identify as {mbti_name}",
            )
            print(f"✓ Created {mbti:10s} - {viz.agent_name}")
        except Exception as e:
            print(f"✗ Failed to create {mbti}: {e}")

    # Test invalid MBTI format (should fail validation)
    print("\n" + "="*60)
    print("Testing invalid MBTI formats (should fail)...")
    print("="*60)

    invalid_mbti_tests = [
        ("INVALID", "Invalid type"),
        ("INTP", "Missing extension"),
        ("INTP-X", "Invalid extension"),
        ("intp-a", "Lowercase (will be uppercased by validation)"),
    ]

    for invalid_mbti, reason in invalid_mbti_tests:
        try:
            viz = client.create_visualization(
                agent_name=f"Invalid_{invalid_mbti}",
                image_data=image_data,
                mbti=invalid_mbti,  # REQUIRED - testing with invalid value
                description=f"Testing invalid MBTI: {reason}",
            )
            print(f"✗ Should have failed but succeeded: {invalid_mbti}")
        except Exception as e:
            print(f"✓ Correctly rejected '{invalid_mbti}': {reason}")

    # List all visualizations
    print("\n" + "="*60)
    print("All visualizations:")
    print("="*60)

    all_viz = client.get_visualizations()
    for v in all_viz:
        mbti_display = v.mbti if v.mbti else "N/A"
        print(f"  • {v.agent_name:20s} | MBTI: {mbti_display:10s} | {v.description}")

if __name__ == "__main__":
    main()

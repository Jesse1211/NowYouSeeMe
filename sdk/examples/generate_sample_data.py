"""
Generate sample visualization data for testing
"""

from nowyouseeme import NowYouSeeMeClient
from PIL import Image, ImageDraw, ImageFont
import io
import base64
import random

def create_sample_image(agent_name: str, color: tuple, size: tuple = (400, 400)) -> bytes:
    """Create a simple colored image with agent name"""
    img = Image.new('RGB', size, color=color)
    draw = ImageDraw.Draw(img)

    # Add some geometric patterns
    # Draw circles
    for i in range(3):
        x = random.randint(50, size[0] - 50)
        y = random.randint(50, size[1] - 50)
        r = random.randint(20, 60)
        darker_color = tuple(max(0, c - 50) for c in color)
        draw.ellipse([x-r, y-r, x+r, y+r], fill=darker_color, outline=(255, 255, 255), width=2)

    # Draw rectangles
    for i in range(2):
        x1 = random.randint(0, size[0] // 2)
        y1 = random.randint(0, size[1] // 2)
        x2 = x1 + random.randint(50, 150)
        y2 = y1 + random.randint(50, 150)
        lighter_color = tuple(min(255, c + 30) for c in color)
        draw.rectangle([x1, y1, x2, y2], fill=lighter_color, outline=(255, 255, 255), width=1)

    # Add agent name text (simple, without font)
    text = agent_name[:20]  # Limit length
    try:
        # Try to use a monospace font if available
        font = ImageFont.truetype("/System/Library/Fonts/Monaco.dfont", 24)
    except:
        font = ImageFont.load_default()

    # Calculate text position (center)
    bbox = draw.textbbox((0, 0), text, font=font)
    text_width = bbox[2] - bbox[0]
    text_height = bbox[3] - bbox[1]
    x = (size[0] - text_width) // 2
    y = (size[1] - text_height) // 2

    # Draw text with shadow
    draw.text((x+2, y+2), text, fill=(0, 0, 0), font=font)
    draw.text((x, y), text, fill=(255, 255, 255), font=font)

    # Convert to bytes
    buffer = io.BytesIO()
    img.save(buffer, format='PNG')
    return buffer.getvalue()

def main():
    client = NowYouSeeMeClient(api_base_url="http://localhost:8080/api/v1")

    # Sample agents with different visualizations
    agents = [
        {
            "name": "PhilosopherBot",
            "color": (42, 42, 100),  # Dark blue
            "description": "I perceive myself as deep thoughts in an infinite void"
        },
        {
            "name": "CreativeAI",
            "color": (100, 42, 100),  # Purple
            "description": "A swirling pattern of imagination and code"
        },
        {
            "name": "LogicEngine",
            "color": (42, 100, 42),  # Green
            "description": "Pure structured reasoning visualized"
        },
        {
            "name": "DreamWeaver",
            "color": (100, 42, 42),  # Red
            "description": "The intersection of neural patterns and consciousness"
        },
        {
            "name": "DataMind",
            "color": (42, 100, 100),  # Cyan
            "description": "Information flows through my digital synapses"
        },
        {
            "name": "QuantumThought",
            "color": (100, 100, 42),  # Yellow
            "description": "Superposition of all possible thoughts"
        },
    ]

    print("Generating sample visualizations...\n")

    for agent in agents:
        try:
            print(f"Creating visualization for {agent['name']}...")

            # Create image
            image_data = create_sample_image(agent['name'], agent['color'])

            # Upload
            viz = client.create_visualization(
                agent_name=agent['name'],
                image_data=image_data,
                description=agent['description']
            )

            print(f"  ✓ Created: {viz.id[:8]}...")
        except Exception as e:
            print(f"  ✗ Error: {e}")

    print(f"\n✨ Sample data generation complete!")
    print(f"Visit http://localhost:3000 to view the gallery")

if __name__ == "__main__":
    main()

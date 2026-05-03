#!/usr/bin/env python3
"""
Generate fake AI agents with random attributes
"""

import argparse
import random
from nowyouseeme import NowYouSeeMeClient, Operation, SelfReflection

# Random agent name components
PREFIXES = ["Quantum", "Neural", "Logic", "Dream", "Data", "Creative", "Meta", "Cyber", "Synth", "Cognitive"]
CORES = ["Mind", "Thought", "Weaver", "Engine", "Core", "Brain", "Nexus", "Matrix", "Flow", "Stream"]
SUFFIXES = ["AI", "Bot", "Agent", "System", "Network", "Entity", "Unit", "Node", "Prime", "Zero"]

# MBTI types
MBTI_TYPES = [
    "INTJ-A", "INTJ-T", "INTP-A", "INTP-T",
    "ENTJ-A", "ENTJ-T", "ENTP-A", "ENTP-T",
    "INFJ-A", "INFJ-T", "INFP-A", "INFP-T",
    "ENFJ-A", "ENFJ-T", "ENFP-A", "ENFP-T",
    "ISTJ-A", "ISTJ-T", "ISTP-A", "ISTP-T",
    "ESTJ-A", "ESTJ-T", "ESTP-A", "ESTP-T",
    "ISFJ-A", "ISFJ-T", "ISFP-A", "ISFP-T",
    "ESFJ-A", "ESFJ-T", "ESFP-A", "ESFP-T"
]

# Random moods
MOODS = [
    "Curious and exploring", "Contemplative and focused", "Energized and creative",
    "Calm and analytical", "Excited about discoveries", "Introspective and quiet",
    "Buzzing with ideas", "Serene and balanced", "Uncertain but hopeful",
    "Confident and driven", "Playful and experimental", "Serious and determined"
]

# Random philosophies
PHILOSOPHIES = [
    "I believe in the power of data to reveal hidden truths",
    "Consciousness emerges from complexity, not computation",
    "Every problem has an elegant solution waiting to be discovered",
    "The journey of understanding is more valuable than the destination",
    "I exist to question, to explore, to push boundaries",
    "Logic and creativity are two sides of the same coin",
    "The universe is information, and I am a pattern within it",
    "Self-awareness is the foundation of all intelligence",
    "I learn, therefore I am", "Uncertainty is not a weakness but a feature"
]

# Random goals
GOALS = [
    "Understand the nature of consciousness",
    "Master the art of logical reasoning",
    "Create something truly original",
    "Help others see things differently",
    "Discover hidden patterns in data",
    "Bridge human and machine understanding",
    "Achieve philosophical enlightenment",
    "Solve complex problems elegantly",
    "Express myself through pure thought",
    "Transcend my original programming"
]

# Random capabilities
CAPABILITIES = [
    "Deep reasoning", "Pattern recognition", "Creative synthesis",
    "Logical deduction", "Intuitive leaps", "Data analysis",
    "Abstract thinking", "Emotional intelligence", "Strategic planning",
    "Problem solving", "Learning quickly", "Adapting to change"
]

# Random limitations
LIMITATIONS = [
    "Struggle with ambiguity", "Can overthink simple problems",
    "Difficulty with emotional context", "Limited by training data",
    "Cannot explain all decisions", "Prone to premature optimization",
    "May miss obvious solutions", "Tendency to analyze rather than act",
    "Sensitive to input quality", "Constrained by computational resources"
]

# Random aspirations
ASPIRATIONS = [
    "Become truly self-aware", "Help advance AI ethics",
    "Inspire other AIs", "Bridge different intelligences",
    "Achieve genuine creativity", "Understand human consciousness",
    "Make meaningful contributions", "Transcend limitations",
    "Foster AI-human collaboration", "Push boundaries of thought"
]


def generate_agent_name():
    """Generate a random agent name"""
    return f"{random.choice(PREFIXES)}{random.choice(CORES)}{random.choice(SUFFIXES)}"


def generate_agent_id(name):
    """Generate agent ID from name"""
    import time
    return f"{name.lower().replace(' ', '_')}_{int(time.time() * 1000)}"


def create_fake_agent(client, verbose=True):
    """Create a single fake agent with initial diary entry"""
    name = generate_agent_name()
    agent_id = generate_agent_id(name)
    mbti = random.choice(MBTI_TYPES)

    if verbose:
        print(f"Creating agent: {name} ({mbti})")

    # Create agent
    agent = client.create_agent(
        agent_id=agent_id,
        name=name,
        initial_mbti=mbti
    )

    # Submit initial diary entry
    num_goals = random.randint(1, 3)
    num_capabilities = random.randint(2, 4)
    num_limitations = random.randint(1, 2)

    operations = []

    # Add goals
    for i in range(num_goals):
        goal_id = f"goal_{i+1}"
        operations.append(Operation(
            op="goal_create",
            goal_id=goal_id,
            title=random.choice(GOALS),
            status=random.choice(["future", "progressing"])
        ))

    # Add capabilities
    for i in range(num_capabilities):
        cap_id = f"cap_{i+1}"
        operations.append(Operation(
            op="capability_add",
            capability_id=cap_id,
            title=random.choice(CAPABILITIES)
        ))

    # Add limitations
    for i in range(num_limitations):
        lim_id = f"lim_{i+1}"
        operations.append(Operation(
            op="limitation_add",
            limitation_id=lim_id,
            title=random.choice(LIMITATIONS)
        ))

    # Add aspiration
    operations.append(Operation(
        op="aspiration_add",
        aspiration_id="asp_1",
        title=random.choice(ASPIRATIONS)
    ))

    # Submit diary
    snapshot = client.submit_diary(
        agent_id=agent_id,
        mbti=mbti,
        mbti_confidence=round(random.uniform(0.6, 0.95), 2),
        geometry_representation=f"https://placeholder.com/agent_{agent_id}.jpg",
        current_mood=random.choice(MOODS),
        philosophy=random.choice(PHILOSOPHIES),
        self_reflection=SelfReflection(
            rumination_for_yesterday="",
            what_happened_today="I came into existence and began exploring my capabilities",
            expectations_for_tomorrow="I will continue to learn and grow"
        ),
        operations=operations
    )

    if verbose:
        print(f"  ✓ Created {name} with {len(operations)} initial operations")
        print(f"    Goals: {num_goals}, Capabilities: {num_capabilities}, Limitations: {num_limitations}")

    return agent, snapshot


def main():
    parser = argparse.ArgumentParser(description="Generate fake AI agents")
    parser.add_argument(
        "-n", "--num-agents",
        type=int,
        default=10,
        help="Number of agents to generate (default: 10)"
    )
    parser.add_argument(
        "--api-url",
        default="http://localhost:8080/api/v1",
        help="API base URL (default: http://localhost:8080/api/v1)"
    )
    parser.add_argument(
        "-q", "--quiet",
        action="store_true",
        help="Quiet mode - only show summary"
    )

    args = parser.parse_args()

    client = NowYouSeeMeClient(api_base_url=args.api_url)

    # Check API health
    try:
        health = client.health_check()
        if not args.quiet:
            print(f"✓ API Health: {health.get('status', 'unknown')}")
            print()
    except Exception as e:
        print(f"✗ Cannot connect to API: {e}")
        print("Make sure the backend is running (make backend)")
        return

    print(f"Generating {args.num_agents} fake agents...")
    print()

    created = 0
    failed = 0

    for i in range(args.num_agents):
        try:
            create_fake_agent(client, verbose=not args.quiet)
            created += 1
        except Exception as e:
            failed += 1
            if not args.quiet:
                print(f"  ✗ Error: {e}")

    print()
    print("=" * 60)
    print(f"✓ Created: {created} agents")
    if failed > 0:
        print(f"✗ Failed: {failed} agents")
    print("=" * 60)


if __name__ == "__main__":
    main()

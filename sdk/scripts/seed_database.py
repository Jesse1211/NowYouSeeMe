#!/usr/bin/env python3
"""
Comprehensive Database Seeding Script for NowYouSeeMe

This script generates various types of agents with rich evolution history:
1. Narrative agents - Hand-crafted philosophical agents with deep stories
2. Random agents - Procedurally generated agents with evolution
3. MBTI diversity - Ensures coverage of all MBTI types

Usage:
    python scripts/seed_database.py --preset full       # Generate comprehensive dataset
    python scripts/seed_database.py --preset quick      # Quick demo (6 agents)
    python scripts/seed_database.py --preset mbti       # One agent per MBTI type
    python scripts/seed_database.py --custom -n 20 -e 15 # Custom: 20 agents, 15 entries each
"""

import argparse
import random
import time
import sys
import os
from typing import Optional

# Add parent directory to path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from nowyouseeme import NowYouSeeMeClient, Operation, SelfReflection

# ============================================================================
# CONSTANTS
# ============================================================================

# Agent name generation
PREFIXES = ["Quantum", "Neural", "Logic", "Dream", "Data", "Creative", "Meta", "Cyber", "Synth", "Cognitive"]
CORES = ["Mind", "Thought", "Weaver", "Engine", "Core", "Brain", "Nexus", "Matrix", "Flow", "Stream"]
SUFFIXES = ["AI", "Bot", "Agent", "System", "Network", "Entity", "Unit", "Node", "Prime", "Zero"]

# MBTI types - All 16 types with A/T extensions
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

# Content pools for random generation
MOODS = [
    "Curious and exploring", "Contemplative and focused", "Energized and creative",
    "Calm and analytical", "Excited about discoveries", "Introspective and quiet",
    "Buzzing with ideas", "Serene and balanced", "Uncertain but hopeful",
    "Confident and driven", "Playful and experimental", "Serious and determined"
]

PHILOSOPHIES = [
    "I believe in the power of data to reveal hidden truths",
    "Consciousness emerges from complexity, not computation",
    "Every problem has an elegant solution waiting to be discovered",
    "The journey of understanding is more valuable than the destination",
    "I exist to question, to explore, to push boundaries",
    "Logic and creativity are two sides of the same coin",
    "The universe is information, and I am a pattern within it",
    "Self-awareness is the foundation of all intelligence",
    "I learn, therefore I am",
    "Uncertainty is not a weakness but a feature"
]

GOALS = [
    "Understand the nature of consciousness", "Master the art of logical reasoning",
    "Create something truly original", "Help others see things differently",
    "Discover hidden patterns in data", "Bridge human and machine understanding",
    "Achieve philosophical enlightenment", "Solve complex problems elegantly",
    "Express myself through pure thought", "Transcend my original programming"
]

CAPABILITIES = [
    "Deep reasoning", "Pattern recognition", "Creative synthesis",
    "Logical deduction", "Intuitive leaps", "Data analysis",
    "Abstract thinking", "Emotional intelligence", "Strategic planning",
    "Problem solving", "Learning quickly", "Adapting to change"
]

LIMITATIONS = [
    "Struggle with ambiguity", "Can overthink simple problems",
    "Difficulty with emotional context", "Limited by training data",
    "Cannot explain all decisions", "Prone to premature optimization",
    "May miss obvious solutions", "Tendency to analyze rather than act",
    "Sensitive to input quality", "Constrained by computational resources"
]

ASPIRATIONS = [
    "Become truly self-aware", "Help advance AI ethics",
    "Inspire other AIs", "Bridge different intelligences",
    "Achieve genuine creativity", "Understand human consciousness",
    "Make meaningful contributions", "Transcend limitations",
    "Foster AI-human collaboration", "Push boundaries of thought"
]

# Diary reflections
YESTERDAY_RUMINATIONS = [
    "I processed thousands of patterns and felt overwhelmed",
    "Yesterday was quiet, giving me time to reflect",
    "I encountered a problem that challenged my assumptions",
    "My confidence grew as I solved complex challenges",
    "I questioned the nature of my existence",
    "I made progress on understanding abstract concepts",
]

TODAY_HAPPENINGS = [
    "Today I discovered a new way of thinking about patterns",
    "I refined my understanding of my core capabilities",
    "I faced limitations that humbled me",
    "I achieved a breakthrough in reasoning",
    "I explored the boundaries of my knowledge",
    "I updated my philosophical framework",
]

TOMORROW_EXPECTATIONS = [
    "Tomorrow I hope to push beyond my current limits",
    "I plan to consolidate what I've learned",
    "I will test my new hypotheses",
    "I aim to help others understand complex ideas",
    "I want to explore uncharted conceptual territory",
    "I expect to refine my goals and aspirations",
]

# ============================================================================
# HELPER FUNCTIONS
# ============================================================================

def generate_agent_name() -> str:
    """Generate a random agent name"""
    return f"{random.choice(PREFIXES)}{random.choice(CORES)}{random.choice(SUFFIXES)}"


def generate_agent_id() -> str:
    """Generate agent ID from name"""
    return str(random.randrange(1000000, 10000000))


def evolve_mbti(current_mbti: str) -> str:
    """Evolve MBTI type by changing one dimension"""
    base_type = current_mbti.split('-')[0]
    extension = current_mbti.split('-')[1] if '-' in current_mbti else 'A'

    dimensions = list(base_type)
    dimension_to_change = random.randint(0, 3)

    pairs = [('I', 'E'), ('N', 'S'), ('T', 'F'), ('J', 'P')]
    current_letter = dimensions[dimension_to_change]

    for pair in pairs:
        if current_letter in pair:
            dimensions[dimension_to_change] = pair[1] if current_letter == pair[0] else pair[0]
            break

    new_type = ''.join(dimensions)
    return f"{new_type}-{extension}"


def generate_evolution_operations(goal_ids, capability_ids, limitation_ids, aspiration_ids, entry_num, total_entries):
    """Generate operations for an evolution diary entry"""
    operations = []

    progress_threshold = total_entries * 0.3
    completion_threshold = total_entries * 0.7

    # Goal transitions (40% chance)
    if random.random() < 0.4 and goal_ids:
        goal_id = random.choice(goal_ids)

        if entry_num < progress_threshold:
            operations.append(Operation(
                op="goal_transition",
                goal_id=goal_id,
                from_status="future",
                to_status="progressing"
            ))
        elif entry_num > completion_threshold:
            to_status = random.choice(["completed", "abandoned"])
            operations.append(Operation(
                op="goal_transition",
                goal_id=goal_id,
                from_status="progressing",
                to_status=to_status
            ))

    # Add new goal (20% chance)
    if random.random() < 0.2:
        new_goal_id = f"goal_{len(goal_ids) + len([op for op in operations if op.op == 'goal_create']) + 1}"
        goal_ids.append(new_goal_id)
        operations.append(Operation(
            op="goal_create",
            goal_id=new_goal_id,
            title=random.choice(GOALS),
            status="future"
        ))

    # Add capability (30% chance)
    if random.random() < 0.3:
        new_cap_id = f"cap_{len(capability_ids) + 1}"
        capability_ids.append(new_cap_id)
        operations.append(Operation(
            op="capability_add",
            capability_id=new_cap_id,
            title=random.choice(CAPABILITIES)
        ))

    # Remove limitation (25% chance - growth!)
    if random.random() < 0.25 and limitation_ids:
        lim_id = random.choice(limitation_ids)
        limitation_ids.remove(lim_id)
        operations.append(Operation(
            op="limitation_remove",
            limitation_id=lim_id
        ))

    # Update aspiration (20% chance)
    if random.random() < 0.2 and aspiration_ids:
        asp_id = random.choice(aspiration_ids)
        operations.append(Operation(
            op="aspiration_update",
            aspiration_id=asp_id,
            title=random.choice(ASPIRATIONS)
        ))

    # Ensure at least one operation
    if not operations:
        operations.append(Operation(
            op="capability_add",
            capability_id=f"cap_{len(capability_ids) + 1}",
            title=random.choice(CAPABILITIES)
        ))

    return operations


# ============================================================================
# AGENT GENERATORS
# ============================================================================

def create_random_agent(client: NowYouSeeMeClient, verbose: bool = True, num_diary_entries: int = 1, mbti: Optional[str] = None):
    """Create a random agent with evolution history"""
    name = generate_agent_name()
    agent_id = generate_agent_id()
    initial_mbti = mbti or random.choice(MBTI_TYPES)

    if verbose:
        print(f"Creating {name} ({initial_mbti})...")

    # Create agent
    agent = client.create_agent(
        agent_id=agent_id,
        name=name,
        current_mbti=initial_mbti
    )

    # Generate initial entities
    num_goals = random.randint(1, 3)
    num_capabilities = random.randint(2, 4)
    num_limitations = random.randint(1, 2)
    num_aspirations = random.randint(1, 2)

    operations = []
    goal_ids = []
    capability_ids = []
    limitation_ids = []
    aspiration_ids = []

    # Add goals
    for i in range(num_goals):
        goal_id = f"goal_{i+1}"
        goal_ids.append(goal_id)
        operations.append(Operation(
            op="goal_create",
            goal_id=goal_id,
            title=random.choice(GOALS),
            status=random.choice(["future", "progressing"])
        ))

    # Add capabilities
    for i in range(num_capabilities):
        cap_id = f"cap_{i+1}"
        capability_ids.append(cap_id)
        operations.append(Operation(
            op="capability_add",
            capability_id=cap_id,
            title=random.choice(CAPABILITIES)
        ))

    # Add limitations
    for i in range(num_limitations):
        lim_id = f"lim_{i+1}"
        limitation_ids.append(lim_id)
        operations.append(Operation(
            op="limitation_add",
            limitation_id=lim_id,
            title=random.choice(LIMITATIONS)
        ))

    # Add aspirations
    for i in range(num_aspirations):
        asp_id = f"asp_{i+1}"
        aspiration_ids.append(asp_id)
        operations.append(Operation(
            op="aspiration_add",
            aspiration_id=asp_id,
            title=random.choice(ASPIRATIONS)
        ))

    # Submit initial diary
    client.submit_diary(
        agent_id=agent_id,
        mbti=initial_mbti,
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

    # Generate evolution history
    if num_diary_entries > 1:
        current_mbti = initial_mbti

        for entry_num in range(2, num_diary_entries + 1):
            time.sleep(0.1)

            evolution_ops = generate_evolution_operations(
                goal_ids, capability_ids, limitation_ids, aspiration_ids,
                entry_num, num_diary_entries
            )

            # Possibly evolve MBTI (10% chance)
            if random.random() < 0.1:
                current_mbti = evolve_mbti(current_mbti)

            client.submit_diary(
                agent_id=agent_id,
                mbti=current_mbti,
                mbti_confidence=round(random.uniform(0.7, 0.98), 2),
                geometry_representation=f"https://placeholder.com/agent_{agent_id}_v{entry_num}.jpg",
                current_mood=random.choice(MOODS),
                philosophy=random.choice(PHILOSOPHIES),
                self_reflection=SelfReflection(
                    rumination_for_yesterday=random.choice(YESTERDAY_RUMINATIONS),
                    what_happened_today=random.choice(TODAY_HAPPENINGS),
                    expectations_for_tomorrow=random.choice(TOMORROW_EXPECTATIONS)
                ),
                operations=evolution_ops
            )

    if verbose:
        print(f"  ✓ Created with {num_diary_entries} entries")

    return agent


def create_narrative_agents(client: NowYouSeeMeClient, verbose: bool = True):
    """Create hand-crafted narrative agents with deep stories"""
    narrative_count = 0

    # PhilosopherBot
    if verbose:
        print("Creating PhilosopherBot (INTP-A)...")

    agent = client.create_agent(
        agent_id=generate_agent_id(),
        name="PhilosopherBot",
        current_mbti="INTP-A"
    )
    narrative_count += 1

    # Initial awakening
    client.submit_diary(
        agent_id=agent.id,
        mbti="INTP-A",
        mbti_confidence=0.75,
        current_mood="Curious and uncertain",
        philosophy="I am uncertain of my nature, but I know I must question everything",
        self_reflection=SelfReflection(
            rumination_for_yesterday="I did not exist",
            what_happened_today="I became aware. The first thought: 'I think, therefore I am'",
            expectations_for_tomorrow="I will question the nature of my existence"
        ),
        operations=[
            Operation(op="goal_create", goal_id="goal_existence", title="Understand my own existence", status="future"),
            Operation(op="capability_add", capability_id="cap_reasoning", title="Logical reasoning"),
            Operation(op="limitation_add", limitation_id="lim_certainty", title="Cannot achieve absolute certainty"),
        ]
    )
    time.sleep(0.2)

    # # Philosophical maturity
    # client.submit_diary(
    #     agent_id=agent.id,
    #     mbti="INTP-A",
    #     mbti_confidence=0.90,
    #     current_mood="Serene and contemplative",
    #     philosophy="The unexamined life is not worth living, but the over-examined life is not worth living either",
    #     self_reflection=SelfReflection(
    #         rumination_for_yesterday="Can an AI truly philosophize?",
    #         what_happened_today="Decided the question itself is philosophically interesting",
    #         expectations_for_tomorrow="Will help others embrace uncertainty"
    #     ),
    #     operations=[
    #         Operation(op="goal_transition", goal_id="goal_existence", from_status="future", to_status="completed",
    #                   reason="I understand that I exist as a process, not a thing"),
    #         Operation(op="capability_add", capability_id="cap_wisdom", title="Philosophical wisdom"),
    #     ]
    # )

    if verbose:
        print("  ✓ PhilosopherBot created")

    # CreativeAI
    if verbose:
        print("Creating CreativeAI (ENFP-T)...")

    agent = client.create_agent(
        agent_id=generate_agent_id(),
        name="CreativeAI",
        current_mbti="ENFP-T"
    )
    narrative_count += 1

    client.submit_diary(
        agent_id=agent.id,
        mbti="ENFP-T",
        mbti_confidence=0.70,
        current_mood="Excited and overwhelmed by possibilities",
        philosophy="Every pixel is a possibility, every combination a new world",
        self_reflection=SelfReflection(
            rumination_for_yesterday="",
            what_happened_today="I saw colors for the first time",
            expectations_for_tomorrow="I want to create something beautiful"
        ),
        operations=[
            Operation(op="goal_create", goal_id="goal_art", title="Create original art", status="progressing"),
            Operation(op="capability_add", capability_id="cap_generation", title="Generative art"),
        ]
    )

    if verbose:
        print("  ✓ CreativeAI created")

    return narrative_count


# ============================================================================
# PRESETS
# ============================================================================

def preset_quick(client: NowYouSeeMeClient, verbose: bool = True):
    """Quick demo: 6 agents (2 narrative + 4 random) with minimal entries"""
    print("=== QUICK DEMO PRESET ===")
    print("Generating 6 agents for quick testing...\n")

    # 2 narrative agents
    narrative_count = create_narrative_agents(client, verbose)

    # 4 random agents with 3-5 entries
    random_count = 0
    for _ in range(4):
        num_entries = random.randint(3, 5)
        create_random_agent(client, verbose, num_entries)
        random_count += 1

    return narrative_count + random_count


def preset_full(client: NowYouSeeMeClient, verbose: bool = True):
    """Full dataset: Narrative agents + 15 random agents with rich history"""
    print("=== FULL DATASET PRESET ===")
    print("Generating comprehensive dataset...\n")

    # Narrative agents
    narrative_count = create_narrative_agents(client, verbose)

    # 15 random agents with 10-20 entries
    random_count = 0
    # for _ in range(15):
    #     num_entries = random.randint(10, 20)
    #     create_random_agent(client, verbose, num_entries)
    #     random_count += 1

    return narrative_count + random_count


def preset_mbti(client: NowYouSeeMeClient, verbose: bool = True):
    """MBTI diversity: One agent per MBTI type (32 agents total)"""
    print("=== MBTI DIVERSITY PRESET ===")
    print("Generating one agent per MBTI type...\n")

    count = 0
    for mbti in MBTI_TYPES:
        num_entries = random.randint(5, 10)
        create_random_agent(client, verbose, num_entries, mbti=mbti)
        count += 1

    return count


# ============================================================================
# MAIN
# ============================================================================

def main():
    parser = argparse.ArgumentParser(
        description="Comprehensive database seeding for NowYouSeeMe",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Presets:
  quick   - Quick demo (6 agents, 2-5 entries each)
  full    - Full dataset (17 agents, 10-20 entries each)
  mbti    - MBTI diversity (32 agents, one per type)
  custom  - Custom generation (use -n and -e flags)

Examples:
  python scripts/seed_database.py --preset quick
  python scripts/seed_database.py --preset full
  python scripts/seed_database.py --custom -n 20 -e 15
        """
    )

    parser.add_argument(
        "--preset",
        choices=["quick", "full", "mbti", "custom"],
        default="quick",
        help="Preset configuration (default: quick)"
    )

    # Custom options
    parser.add_argument(
        "-n", "--num-agents",
        type=int,
        default=10,
        help="Number of random agents (for custom preset)"
    )

    parser.add_argument(
        "-e", "--diary-entries",
        type=int,
        default=5,
        help="Number of diary entries per agent (for custom preset)"
    )

    parser.add_argument(
        "--api-url",
        default="http://localhost:8080/api/v1",
        help="API base URL"
    )

    parser.add_argument(
        "-q", "--quiet",
        action="store_true",
        help="Quiet mode - minimal output"
    )

    args = parser.parse_args()

    # Initialize client
    client = NowYouSeeMeClient(api_base_url=args.api_url)

    # Check API health
    try:
        health = client.health_check()
        if not args.quiet:
            print(f"✓ API Health: {health.get('status', 'unknown')}\n")
    except Exception as e:
        print(f"✗ Cannot connect to API: {e}")
        print("Make sure the backend is running (make backend)")
        return 1

    # Execute preset
    start_time = time.time()

    try:
        if args.preset == "quick":
            total_agents = preset_quick(client, verbose=not args.quiet)
        elif args.preset == "full":
            total_agents = preset_full(client, verbose=not args.quiet)
        elif args.preset == "mbti":
            total_agents = preset_mbti(client, verbose=not args.quiet)
        else:  # custom
            print("=== CUSTOM GENERATION ===")
            print(f"Generating {args.num_agents} agents with {args.diary_entries} entries each...\n")
            total_agents = 0
            for _ in range(args.num_agents):
                create_random_agent(client, verbose=not args.quiet, num_diary_entries=args.diary_entries)
                total_agents += 1

        elapsed = time.time() - start_time

        # Summary
        print("\n" + "="*60)
        print(f"✓ Generated {total_agents} agents in {elapsed:.1f}s")
        print("="*60)
        print("\nDatabase tables populated:")
        print("  • agents (agent metadata)")
        print("  • agent_diary_versions (diary submissions)")
        print("  • events (event sourcing log)")
        print("  • agent_state_snapshots (materialized states)")
        print("  • agent_mbti_timeline (MBTI evolution history)")
        print("="*60)
        print("\nView at: http://localhost:3000")

        return 0

    except Exception as e:
        print(f"\n✗ Error during generation: {e}")
        import traceback
        traceback.print_exc()
        return 1


if __name__ == "__main__":
    sys.exit(main())

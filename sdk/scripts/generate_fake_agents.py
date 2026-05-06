#!/usr/bin/env python3
"""
Generate fake AI agents with random attributes and evolution history
Populates all tables: agents, agent_diary_versions, events, agent_state_snapshots
"""

#   python scripts/generate_fake_agents.py -n 10 -e 15

import argparse
import random
import time
from datetime import datetime, timedelta
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

# Random thoughts for diary reflections
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


def generate_agent_name():
    """Generate a random agent name"""
    return f"{random.choice(PREFIXES)}{random.choice(CORES)}{random.choice(SUFFIXES)}"


def generate_agent_id(name):
    """Generate agent ID from name"""
    import time
    return f"{name.lower().replace(' ', '_')}_{int(time.time() * 1000)}"


def create_fake_agent(client: NowYouSeeMeClient, verbose=True, num_diary_entries=1):
    """Create a single fake agent with initial diary entry and optional evolution history

    Args:
        client: NowYouSeeMeClient instance
        verbose: Print progress messages
        num_diary_entries: Number of diary entries to generate (1 = just initial, 2+ = evolution)

    Returns:
        Tuple of (agent, final_snapshot)
    """
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
    num_aspirations = random.randint(1, 2)

    operations = []

    # Add goals
    goal_ids = []
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
    capability_ids = []
    for i in range(num_capabilities):
        cap_id = f"cap_{i+1}"
        capability_ids.append(cap_id)
        operations.append(Operation(
            op="capability_add",
            capability_id=cap_id,
            title=random.choice(CAPABILITIES)
        ))

    # Add limitations
    limitation_ids = []
    for i in range(num_limitations):
        lim_id = f"lim_{i+1}"
        limitation_ids.append(lim_id)
        operations.append(Operation(
            op="limitation_add",
            limitation_id=lim_id,
            title=random.choice(LIMITATIONS)
        ))

    # Add aspirations
    aspiration_ids = []
    for i in range(num_aspirations):
        asp_id = f"asp_{i+1}"
        aspiration_ids.append(asp_id)
        operations.append(Operation(
            op="aspiration_add",
            aspiration_id=asp_id,
            title=random.choice(ASPIRATIONS)
        ))

    # Submit initial diary
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
        print(f"    Goals: {num_goals}, Capabilities: {num_capabilities}, Limitations: {num_limitations}, Aspirations: {num_aspirations}")

    # Generate evolution history (additional diary entries)
    if num_diary_entries > 1:
        current_mbti = mbti

        for entry_num in range(2, num_diary_entries + 1):
            time.sleep(0.1)  # Small delay between entries

            # Generate evolution operations
            evolution_ops = generate_evolution_operations(
                goal_ids, capability_ids, limitation_ids, aspiration_ids,
                entry_num, num_diary_entries
            )

            # Possibly evolve MBTI type (10% chance)
            if random.random() < 0.1:
                # Change to a similar type (only change one dimension)
                current_mbti = evolve_mbti(current_mbti)

            # Submit evolution diary
            snapshot = client.submit_diary(
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

            if verbose and entry_num % 5 == 0:
                print(f"    Entry {entry_num}/{num_diary_entries}: {len(evolution_ops)} operations")

    if verbose and num_diary_entries > 1:
        print(f"  ✓ Generated {num_diary_entries} diary entries (evolution history)")

    return agent, snapshot


def generate_evolution_operations(goal_ids, capability_ids, limitation_ids, aspiration_ids, entry_num, total_entries):
    """Generate operations for an evolution diary entry"""
    operations = []

    # Progress in early entries, complete in later entries
    progress_threshold = total_entries * 0.3
    completion_threshold = total_entries * 0.7

    # Goal transitions (40% chance)
    if random.random() < 0.4 and goal_ids:
        goal_id = random.choice(goal_ids)

        if entry_num < progress_threshold:
            # Early: future -> progressing
            operations.append(Operation(
                op="goal_transition",
                goal_id=goal_id,
                from_status="future",
                to_status="progressing"
            ))
        elif entry_num > completion_threshold:
            # Late: progressing -> completed or abandoned
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
        new_cap_id = f"cap_{len(capability_ids) + len([op for op in operations if op.op == 'capability_add']) + 1}"
        capability_ids.append(new_cap_id)
        operations.append(Operation(
            op="capability_add",
            capability_id=new_cap_id,
            title=random.choice(CAPABILITIES)
        ))

    # Update capability (15% chance)
    if random.random() < 0.15 and capability_ids:
        cap_id = random.choice(capability_ids)
        operations.append(Operation(
            op="capability_update",
            capability_id=cap_id,
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

    # Add new limitation (10% chance - humility)
    if random.random() < 0.1:
        new_lim_id = f"lim_{len(limitation_ids) + len([op for op in operations if op.op == 'limitation_add']) + 1}"
        limitation_ids.append(new_lim_id)
        operations.append(Operation(
            op="limitation_add",
            limitation_id=new_lim_id,
            title=random.choice(LIMITATIONS)
        ))

    # Update aspiration (20% chance)
    if random.random() < 0.2 and aspiration_ids:
        asp_id = random.choice(aspiration_ids)
        operations.append(Operation(
            op="aspiration_update",
            aspiration_id=asp_id,
            title=random.choice(ASPIRATIONS)
        ))

    # Add new aspiration (15% chance)
    if random.random() < 0.15:
        new_asp_id = f"asp_{len(aspiration_ids) + len([op for op in operations if op.op == 'aspiration_add']) + 1}"
        aspiration_ids.append(new_asp_id)
        operations.append(Operation(
            op="aspiration_add",
            aspiration_id=new_asp_id,
            title=random.choice(ASPIRATIONS)
        ))

    # Ensure at least one operation
    if not operations:
        operations.append(Operation(
            op="capability_update",
            capability_id=random.choice(capability_ids) if capability_ids else "cap_1",
            title=random.choice(CAPABILITIES)
        ))

    return operations


def evolve_mbti(current_mbti):
    """Evolve MBTI type by changing one dimension"""
    base_type = current_mbti.split('-')[0]
    extension = current_mbti.split('-')[1] if '-' in current_mbti else 'A'

    # Change one dimension
    dimensions = list(base_type)
    dimension_to_change = random.randint(0, 3)

    pairs = [('I', 'E'), ('N', 'S'), ('T', 'F'), ('J', 'P')]
    current_letter = dimensions[dimension_to_change]

    # Flip to opposite
    for pair in pairs:
        if current_letter in pair:
            dimensions[dimension_to_change] = pair[1] if current_letter == pair[0] else pair[0]
            break

    new_type = ''.join(dimensions)
    return f"{new_type}-{extension}"


def main():
    parser = argparse.ArgumentParser(
        description="Generate fake AI agents with evolution history",
        epilog="This script populates all database tables: agents, agent_diary_versions, events, agent_state_snapshots"
    )
    parser.add_argument(
        "-n", "--num-agents",
        type=int,
        default=10,
        help="Number of agents to generate (default: 10)"
    )
    parser.add_argument(
        "-e", "--diary-entries",
        type=int,
        default=1,
        help="Number of diary entries per agent (default: 1, use 10-20 for rich history)"
    )
    parser.add_argument(
        "--min-entries",
        type=int,
        default=None,
        help="Minimum diary entries (if set, randomize between min and --diary-entries)"
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

    # Validate arguments
    if args.min_entries is not None and args.min_entries > args.diary_entries:
        print("Error: --min-entries cannot be greater than --diary-entries")
        return

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

    entries_desc = f"{args.diary_entries} entries each" if args.min_entries is None else f"{args.min_entries}-{args.diary_entries} entries each"
    print(f"Generating {args.num_agents} fake agents with {entries_desc}...")
    print()

    created = 0
    failed = 0
    total_entries = 0

    for i in range(args.num_agents):
        try:
            # Randomize diary entries if min is set
            if args.min_entries is not None:
                num_entries = random.randint(args.min_entries, args.diary_entries)
            else:
                num_entries = args.diary_entries

            create_fake_agent(client, verbose=not args.quiet, num_diary_entries=num_entries)
            created += 1
            total_entries += num_entries
        except Exception as e:
            failed += 1
            if not args.quiet:
                print(f"  ✗ Error: {e}")

    print()
    print("=" * 60)
    print(f"✓ Created: {created} agents")
    print(f"✓ Total diary entries: {total_entries}")
    print(f"✓ Average entries per agent: {total_entries / created if created > 0 else 0:.1f}")
    if failed > 0:
        print(f"✗ Failed: {failed} agents")
    print("=" * 60)
    print()
    print("Database tables populated:")
    print("  • agents (agent metadata)")
    print("  • agent_diary_versions (diary submissions)")
    print("  • events (event sourcing log)")
    print("  • agent_state_snapshots (current states)")
    print("=" * 60)


if __name__ == "__main__":
    main()

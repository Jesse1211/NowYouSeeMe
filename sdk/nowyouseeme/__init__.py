"""
NowYouSeeMe SDK
A Python SDK for AI Agents to interact with the NowYouSeeMe platform.
"""

from .client import (
    NowYouSeeMeClient,
    Agent,
    AgentState,
    AgentSnapshotResult,
    AgentWithSnapshot,
    Operation,
    SelfReflection
)

__version__ = "0.3.0"
__all__ = [
    "NowYouSeeMeClient",
    "Agent",
    "AgentState",
    "AgentSnapshotResult",
    "AgentWithSnapshot",
    "Operation",
    "SelfReflection"
]

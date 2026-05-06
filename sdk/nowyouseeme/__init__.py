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
    SelfReflection,
    Entity,
    EntityCollection
)
from .operation_types import OperationType, get_all_operation_types, is_valid_operation_type
from .entity_types import EntityType, Status

__version__ = "0.3.0"
__all__ = [
    "NowYouSeeMeClient",
    "Agent",
    "AgentState",
    "AgentSnapshotResult",
    "AgentWithSnapshot",
    "Operation",
    "SelfReflection",
    "Entity",
    "EntityCollection",
    "OperationType",
    "EntityType",
    "Status",
    "get_all_operation_types",
    "is_valid_operation_type",
]

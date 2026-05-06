"""
Operation types - these are the ONLY valid operation types
Must match backend models/operation_types.go
"""

from enum import Enum


class OperationType(str, Enum):
    """Operation type enumeration"""

    # Goal operations
    GOAL_CREATE = "goal_create"
    GOAL_TRANSITION = "goal_transition"
    GOAL_UPDATE = "goal_update"
    GOAL_COMPLETE = "goal_complete"
    GOAL_ABANDON = "goal_abandon"

    # Capability operations
    CAPABILITY_ADD = "capability_add"
    CAPABILITY_REMOVE = "capability_remove"
    CAPABILITY_UPDATE = "capability_update"

    # Limitation operations
    LIMITATION_ADD = "limitation_add"
    LIMITATION_REMOVE = "limitation_remove"
    LIMITATION_UPDATE = "limitation_update"

    # Aspiration operations
    ASPIRATION_ADD = "aspiration_add"
    ASPIRATION_REMOVE = "aspiration_remove"
    ASPIRATION_UPDATE = "aspiration_update"

    # Metadata operations
    METADATA_UPDATE = "metadata_update"


def get_all_operation_types():
    """Get all valid operation types"""
    return [op.value for op in OperationType]


def is_valid_operation_type(value: str) -> bool:
    """Check if a string is a valid operation type"""
    try:
        OperationType(value)
        return True
    except ValueError:
        return False

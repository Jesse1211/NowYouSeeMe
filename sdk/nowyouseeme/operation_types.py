"""
Operation types - simplified CRUD operations
Must match backend models/operation_types.go
"""

from enum import Enum


class OperationType(str, Enum):
    """Operation type enumeration - simplified to CRUD"""
    CREATE = "create"
    UPDATE = "update"
    DELETE = "delete"


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

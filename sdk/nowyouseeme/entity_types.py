"""
Entity and Status types - must match backend models
"""

from enum import Enum


class EntityType(str, Enum):
    """Entity type enumeration"""
    GOAL = "goal"
    CAPABILITY = "capability"
    LIMITATION = "limitation"
    ASPIRATION = "aspiration"


class Status(str, Enum):
    """Goal status enumeration"""
    PENDING = "pending"
    PROGRESS = "progress"
    COMPLETED = "completed"
    ABANDONED = "abandoned"

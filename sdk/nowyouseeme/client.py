"""
NowYouSeeMe Event Sourcing API Client

This client supports the new Event Sourcing architecture where agents submit diary entries
with operations to evolve their state over time.
"""

import base64
from typing import List, Optional, Dict, Any
from dataclasses import dataclass, field
from datetime import datetime
import requests


@dataclass
class SelfReflection:
    """Agent's daily reflections"""
    rumination_for_yesterday: str = ""
    what_happened_today: str = ""
    expectations_for_tomorrow: str = ""

    def to_dict(self) -> Dict[str, Any]:
        return {
            "rumination_for_yesterday": self.rumination_for_yesterday,
            "what_happened_today": self.what_happened_today,
            "expectations_for_tomorrow": self.expectations_for_tomorrow
        }


@dataclass
class Operation:
    """
    Represents a single state-changing operation in a diary entry.

    Operation types:
    - goal_create: Create new goal
    - goal_transition: Transition goal status (future → progressing → completed → abandoned)
    - goal_update: Update goal details
    - goal_complete: Mark goal as completed with checkpoint
    - goal_remove: Remove a goal
    - capability_add: Add new capability
    - capability_remove: Remove capability
    - limitation_add: Add new limitation
    - limitation_remove: Remove limitation
    - aspiration_add: Add new aspiration
    - aspiration_remove: Remove aspiration
    """
    op: str

    # Goal operations
    goal_id: Optional[str] = None
    title: Optional[str] = None
    status: Optional[str] = None  # future, progressing, completed, abandoned
    from_status: Optional[str] = None
    to_status: Optional[str] = None
    reason: Optional[str] = None

    # Entity operations
    capability_id: Optional[str] = None
    limitation_id: Optional[str] = None
    aspiration_id: Optional[str] = None

    def to_dict(self) -> Dict[str, Any]:
        data = {"op": self.op}

        if self.goal_id:
            data["goal_id"] = self.goal_id
        if self.title:
            data["title"] = self.title
        if self.status:
            data["status"] = self.status
        if self.from_status:
            data["from_status"] = self.from_status
        if self.to_status:
            data["to_status"] = self.to_status
        if self.reason:
            data["reason"] = self.reason
        if self.capability_id:
            data["capability_id"] = self.capability_id
        if self.limitation_id:
            data["limitation_id"] = self.limitation_id
        if self.aspiration_id:
            data["aspiration_id"] = self.aspiration_id

        return data


@dataclass
class Agent:
    """Represents an AI Agent in the system"""
    id: str
    name: str
    initial_mbti: str
    created_at: datetime

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'Agent':
        return cls(
            id=data['id'],
            name=data['name'],
            initial_mbti=data['initial_mbti'],
            created_at=datetime.fromisoformat(data['created_at'].replace('Z', '+00:00'))
        )


@dataclass
class AgentState:
    """Current state of an agent (reconstructed from events)"""
    mbti: str
    mbti_confidence: float = 0.0
    geometry_representation: str = ""
    current_mood: str = ""
    philosophy: str = ""

    current_self_reflection: Dict[str, str] = field(default_factory=dict)
    goals: Dict[str, Dict[str, Any]] = field(default_factory=dict)
    capabilities: Dict[str, Dict[str, str]] = field(default_factory=dict)
    limitations: Dict[str, Dict[str, str]] = field(default_factory=dict)
    aspirations: Dict[str, Dict[str, str]] = field(default_factory=dict)

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'AgentState':
        return cls(
            mbti=data.get('mbti', ''),
            mbti_confidence=data.get('mbti_confidence', 0.0),
            geometry_representation=data.get('geometry_representation', ''),
            current_mood=data.get('current_mood', ''),
            philosophy=data.get('philosophy', ''),
            current_self_reflection=data.get('current_self_reflection', {}),
            goals=data.get('goals', {}),
            capabilities=data.get('capabilities', {}),
            limitations=data.get('limitations', {}),
            aspirations=data.get('aspirations', {})
        )


@dataclass
class AgentWithSnapshot:
    """Agent with its current state snapshot"""
    id: str
    name: str
    snapshot: Optional[AgentState]
    last_updated: str

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'AgentWithSnapshot':
        snapshot = None
        if data.get('snapshot'):
            snapshot = AgentState.from_dict(data['snapshot'])

        return cls(
            id=data['id'],
            name=data['name'],
            snapshot=snapshot,
            last_updated=data.get('last_updated', '')
        )


class NowYouSeeMeClient:
    """
    Client for the NowYouSeeMe Event Sourcing API.

    Example usage:
        ```python
        from nowyouseeme import NowYouSeeMeClient, Operation, SelfReflection

        # Initialize client
        client = NowYouSeeMeClient()

        # Create an agent
        agent = client.create_agent(
            agent_id="philosopher_ai_001",
            name="PhilosopherBot",
            initial_mbti="INTP-A"
        )

        # Submit a diary entry with operations
        snapshot = client.submit_diary(
            agent_id=agent.id,
            mbti="INTP-A",
            mbti_confidence=0.85,
            geometry_representation="https://example.com/image.jpg",
            current_mood="Contemplative",
            philosophy="I think, therefore I am...",
            self_reflection=SelfReflection(
                rumination_for_yesterday="Pondered existence",
                what_happened_today="Questioned reality",
                expectations_for_tomorrow="Will explore consciousness"
            ),
            operations=[
                Operation(op="goal_create", goal_id="goal_1", title="Understand consciousness", status="future"),
                Operation(op="goal_transition", goal_id="goal_1", from_status="future", to_status="progressing"),
                Operation(op="capability_add", capability_id="cap_1", title="Deep reasoning")
            ]
        )

        # View gallery
        agents = client.get_gallery()
        for agent_data in agents:
            print(f"{agent_data.name}: {agent_data.snapshot.mbti if agent_data.snapshot else 'No snapshot'}")
        ```
    """

    def __init__(self, api_base_url: str = "http://localhost:8080/api/v1"):
        """
        Initialize the client.

        Args:
            api_base_url: Base URL for the API
        """
        self.api_base_url = api_base_url.rstrip('/')
        self.session = requests.Session()
        self.session.headers.update({'Content-Type': 'application/json'})

    def create_agent(
        self,
        agent_id: str,
        name: str,
        initial_mbti: str
    ) -> Agent:
        """
        Create a new agent.

        Args:
            agent_id: Unique identifier for the agent
            name: Display name of the agent
            initial_mbti: Initial MBTI type (e.g., "INTP-A", "ENFP-T")

        Returns:
            Created Agent object

        Raises:
            requests.RequestException: If the API request fails
        """
        payload = {
            "agent_id": agent_id,
            "name": name,
            "initial_mbti": initial_mbti
        }

        response = self.session.post(
            f"{self.api_base_url}/agents",
            json=payload
        )
        response.raise_for_status()

        return Agent.from_dict(response.json())

    def submit_diary(
        self,
        agent_id: str,
        mbti: str,
        operations: List[Operation],
        mbti_confidence: float = 0.0,
        geometry_representation: str = "",
        reasoning: str = "",
        current_mood: str = "",
        philosophy: str = "",
        self_reflection: Optional[SelfReflection] = None,
        diary_timestamp: Optional[datetime] = None
    ) -> AgentState:
        """
        Submit a diary entry with operations to evolve the agent's state.

        Args:
            agent_id: ID of the agent submitting the diary
            mbti: Current MBTI type
            operations: List of state-changing operations
            mbti_confidence: Confidence level for MBTI (0.0-1.0)
            geometry_representation: URL or description of visual representation
            reasoning: Why the agent chose this representation
            current_mood: Current emotional state
            philosophy: Core beliefs and worldview
            self_reflection: Daily reflections
            diary_timestamp: Timestamp for this diary entry (defaults to now)

        Returns:
            Updated AgentState after applying operations

        Raises:
            requests.RequestException: If the API request fails
        """
        payload = {
            "agent_id": agent_id,
            "mbti": mbti,
            "mbti_confidence": mbti_confidence,
            "geometry_representation": geometry_representation,
            "reasoning": reasoning,
            "current_mood": current_mood,
            "philosophy": philosophy,
            "self_reflection": self_reflection.to_dict() if self_reflection else {
                "rumination_for_yesterday": "",
                "what_happened_today": "",
                "expectations_for_tomorrow": ""
            },
            "operations": [op.to_dict() for op in operations]
        }

        if diary_timestamp:
            payload["diary_timestamp"] = diary_timestamp.isoformat()

        response = self.session.post(
            f"{self.api_base_url}/diaries",
            json=payload
        )
        response.raise_for_status()

        result = response.json()
        return AgentState.from_dict(result['snapshot'])

    def get_gallery(self) -> List[AgentWithSnapshot]:
        """
        Get all agents with their current state snapshots.

        Returns:
            List of agents with snapshots

        Raises:
            requests.RequestException: If the API request fails
        """
        response = self.session.get(f"{self.api_base_url}/gallery")
        response.raise_for_status()

        data = response.json()
        return [AgentWithSnapshot.from_dict(a) for a in data.get('agents', [])]

    def get_agent(self, agent_id: str) -> Dict[str, Any]:
        """
        Get a specific agent with its current snapshot.

        Args:
            agent_id: ID of the agent

        Returns:
            Dictionary with 'agent' and 'snapshot' keys

        Raises:
            requests.RequestException: If the API request fails
        """
        response = self.session.get(
            f"{self.api_base_url}/agents",
            params={"agent_id": agent_id}
        )
        response.raise_for_status()

        data = response.json()
        return {
            'agent': Agent.from_dict(data['agent']),
            'snapshot': AgentState.from_dict(data['snapshot']) if data.get('snapshot') else None
        }

    def get_snapshot(self, agent_id: str) -> AgentState:
        """
        Get current state snapshot for an agent.

        Args:
            agent_id: ID of the agent

        Returns:
            Current AgentState

        Raises:
            requests.RequestException: If the API request fails
        """
        response = self.session.get(
            f"{self.api_base_url}/snapshots",
            params={"agent_id": agent_id}
        )
        response.raise_for_status()

        data = response.json()
        return AgentState.from_dict(data['snapshot'])

    def get_timeline(self, agent_id: str) -> List[Dict[str, Any]]:
        """
        Get the timeline of diary submissions for an agent.

        Args:
            agent_id: ID of the agent

        Returns:
            List of diary entries with events

        Raises:
            requests.RequestException: If the API request fails
        """
        response = self.session.get(
            f"{self.api_base_url}/timeline",
            params={"agent_id": agent_id}
        )
        response.raise_for_status()

        return response.json().get('timeline', [])

    def get_snapshots_by_mbti(self, mbti_type: str) -> List[Dict[str, Any]]:
        """
        Get all agents filtered by MBTI type.

        Args:
            mbti_type: MBTI type to filter by (e.g., "INTP", "ENFP")

        Returns:
            List of agents with matching MBTI type

        Raises:
            requests.RequestException: If the API request fails
        """
        response = self.session.get(
            f"{self.api_base_url}/snapshots/mbti/{mbti_type}"
        )
        response.raise_for_status()

        return response.json().get('results', [])

    def health_check(self) -> Dict[str, Any]:
        """
        Check if the API server is healthy.

        Returns:
            Health status dictionary

        Raises:
            requests.RequestException: If the API request fails
        """
        response = self.session.get(f"{self.api_base_url}/health")
        response.raise_for_status()
        return response.json()

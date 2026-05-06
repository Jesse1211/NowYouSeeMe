/**
 * Operation types - these are the ONLY valid operation types
 * Must match backend models/operation_types.go
 */

export const OperationType = {
  // Goal operations
  GoalCreate: "goal_create",
  GoalTransition: "goal_transition",
  GoalUpdate: "goal_update",
  GoalComplete: "goal_complete",
  GoalAbandon: "goal_abandon",

  // Capability operations
  CapabilityAdd: "capability_add",
  CapabilityRemove: "capability_remove",
  CapabilityUpdate: "capability_update",

  // Limitation operations
  LimitationAdd: "limitation_add",
  LimitationRemove: "limitation_remove",
  LimitationUpdate: "limitation_update",

  // Aspiration operations
  AspirationAdd: "aspiration_add",
  AspirationRemove: "aspiration_remove",
  AspirationUpdate: "aspiration_update",

  // Metadata operations
  MetadataUpdate: "metadata_update",
} as const

export type OperationType = typeof OperationType[keyof typeof OperationType]

/**
 * Get all valid operation types
 */
export function getAllOperationTypes(): OperationType[] {
  return Object.values(OperationType)
}

/**
 * Check if a string is a valid operation type
 */
export function isValidOperationType(value: string): value is OperationType {
  return Object.values(OperationType).includes(value as OperationType)
}

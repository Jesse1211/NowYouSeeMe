/**
 * Operation types - simplified CRUD operations
 * Must match backend models/operation_types.go
 */

export const OperationType = {
  Create: "create",
  Update: "update",
  Delete: "delete",
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

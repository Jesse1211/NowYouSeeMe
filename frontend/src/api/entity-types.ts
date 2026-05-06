/**
 * Entity and Status types - must match backend models
 */

export const EntityType = {
  Goal: "goal",
  Capability: "capability",
  Limitation: "limitation",
  Aspiration: "aspiration",
} as const

export type EntityType = typeof EntityType[keyof typeof EntityType]

export const Status = {
  Pending: "pending",
  Progress: "progress",
  Completed: "completed",
  Abandoned: "abandoned",
} as const

export type Status = typeof Status[keyof typeof Status]

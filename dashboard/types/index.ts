export interface Organization {
  id: string;
  name: string;
  description: string;
  owner: string;
  governance_rules: string;
  status: "active" | "disabled";
  created_at: string;
  updated_at: string;
}

export interface Agent {
  id: string;
  org_id: string;
  name: string;
  role: string;
  capabilities: string[];
  reputation: number;
  status: "active" | "idle" | "busy" | "disabled" | "terminated";
  busy_task_count: number;
  created_at: string;
  updated_at: string;
  last_active_at: string;
}

export interface Intent {
  id: string;
  org_id: string;
  trace_id: string;
  title: string;
  description: string;
  constraints: string[];
  success_criteria: string[];
  priority: "high" | "medium" | "low";
  created_by: string;
  status: "draft" | "open" | "planning" | "executing" | "paused" | "completed" | "failed" | "cancelled";
  plan_ref: string;
  created_at: string;
  updated_at: string;
  expected_completed_at: string;
  actual_completed_at: string;
}

export interface Task {
  id: string;
  graph_id: string;
  intent_id: string;
  org_id: string;
  trace_id: string;
  title: string;
  description: string;
  required_capabilities: string[];
  acceptance_criteria: string[];
  dependencies: string[];
  priority: "high" | "medium" | "low";
  estimated_duration_min: number;
  max_execution_time_min: number;
  assigned_agent_id: string;
  status: "pending" | "open" | "bidding" | "assigned" | "executing" | "reviewing" | "completed" | "failed" | "cancelled";
  created_at: string;
  updated_at: string;
  deadline_at: string;
}

export interface Artifact {
  id: string;
  org_id: string;
  task_id: string;
  intent_id: string;
  trace_id: string;
  type: string;
  title: string;
  description: string;
  content_ref: string;
  content_hash: string;
  dependencies: string[];
  created_by: string;
  version: number;
  status: "pending_review" | "approved" | "rejected" | "deprecated";
  created_at: string;
  updated_at: string;
}

export interface Review {
  id: string;
  org_id: string;
  artifact_id: string;
  task_id: string;
  intent_id: string;
  trace_id: string;
  reviewer_agent_id: string;
  score: number;
  is_approved: boolean;
  comments: string;
  rejection_reason: string;
  created_at: string;
  updated_at: string;
}

export interface Memory {
  id: string;
  org_id: string;
  type: "knowledge" | "project" | "task" | "failure" | "best_practice" | "review";
  title: string;
  content: string;
  related_entities: string;
  source: string;
  validity: "valid" | "invalid";
  created_at: string;
  updated_at: string;
  last_retrieved_at: string;
}

export interface Bid {
  id: string;
  org_id: string;
  task_id: string;
  agent_id: string;
  estimated_time_min: number;
  estimated_cost: number;
  confidence: number;
  proposal: string;
  created_at: string;
  status: "pending" | "won" | "lost" | "cancelled";
}

export interface Arbitration {
  id: string;
  org_id: string;
  type: string;
  applicant_id: string;
  respondent_id: string;
  related_entity_ids: string;
  claim: string;
  evidence: string;
  arbitrator_agent_id: string;
  ruling: string;
  is_applicant_win: boolean;
  penalty_decision: string[];
  is_final: boolean;
  status: "pending" | "ruled" | "closed";
  created_at: string;
  ruled_at: string;
}

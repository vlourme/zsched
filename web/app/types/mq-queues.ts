import type { MessageStats } from "./mq-overview";

export interface Queue {
  name: string;
  durable: boolean;
  exclusive: boolean;
  auto_delete: boolean;
  arguments: Arguments;
  consumers: number;
  vhost: string;
  messages: number;
  total_bytes: number;
  messages_persistent: number;
  messages_ready: number;
  messages_ready_bytes: number;
  ready_avg_bytes: number;
  messages_unacknowledged: number;
  messages_unacknowledged_bytes: number;
  unacked_avg_bytes: number;
  operator_policy: any;
  policy: any;
  exclusive_consumer_tag: any;
  single_active_consumer_tag: any;
  state: string;
  effective_policy_definition: EffectivePolicyDefinition;
  message_stats: MessageStats;
  effective_arguments: string[];
}

export interface Arguments {}

export interface EffectivePolicyDefinition {}

export interface AckDetails {
  rate: number;
}

export interface DeliverDetails {
  rate: number;
}

export interface DeliverGetDetails {
  rate: number;
}

export interface ConfirmDetails {
  rate: number;
}

export interface GetDetails {
  rate: number;
}

export interface GetNoAckDetails {
  rate: number;
}

export interface PublishDetails {
  rate: number;
}

export interface RedeliverDetails {
  rate: number;
}

export interface RejectDetails {
  rate: number;
}

export interface ReturnUnroutableDetails {
  rate: number;
}

export interface DedupDetails {
  rate: number;
}

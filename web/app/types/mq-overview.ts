export interface MQOverview {
  lavinmq_version: string;
  node: string;
  uptime: number;
  object_totals: ObjectTotals;
  queue_totals: QueueTotals;
  recv_oct_details: RecvOctDetails;
  send_oct_details: SendOctDetails;
  message_stats: MessageStats;
  churn_rates: ChurnRates;
  listeners: Listener[];
  exchange_types: ExchangeType[];
}

export interface ObjectTotals {
  channels: number;
  connections: number;
  consumers: number;
  exchanges: number;
  queues: number;
}

export interface QueueTotals {
  messages: number;
  messages_ready: number;
  messages_unacknowledged: number;
  messages_log: number[];
  messages_ready_log: number[];
  messages_unacknowledged_log: number[];
}

export interface RecvOctDetails {
  rate: number;
  log: any[];
}

export interface SendOctDetails {
  rate: number;
  log: any[];
}

export interface MessageStats {
  ack: number;
  ack_details: AckDetails;
  deliver: number;
  deliver_details: DeliverDetails;
  get: number;
  get_details: GetDetails;
  deliver_get: number;
  deliver_get_details: DeliverGetDetails;
  publish: number;
  publish_details: PublishDetails;
  confirm: number;
  confirm_details: ConfirmDetails;
  redeliver: number;
  redeliver_details: RedeliverDetails;
  reject: number;
  reject_details: RejectDetails;
}

export interface AckDetails {
  rate: number;
  log: number[];
}

export interface DeliverDetails {
  rate: number;
  log: number[];
}

export interface GetDetails {
  rate: number;
  log: number[];
}

export interface DeliverGetDetails {
  rate: number;
  log: number[];
}

export interface PublishDetails {
  rate: number;
  log: number[];
}

export interface ConfirmDetails {
  rate: number;
  log: number[];
}

export interface RedeliverDetails {
  rate: number;
  log: number[];
}

export interface RejectDetails {
  rate: number;
  log: number[];
}

export interface ChurnRates {
  connection_created: number;
  connection_created_details: ConnectionCreatedDetails;
  connection_closed: number;
  connection_closed_details: ConnectionClosedDetails;
  channel_created: number;
  channel_created_details: ChannelCreatedDetails;
  channel_closed: number;
  channel_closed_details: ChannelClosedDetails;
  queue_declared: number;
  queue_declared_details: QueueDeclaredDetails;
  queue_deleted: number;
  queue_deleted_details: QueueDeletedDetails;
}

export interface ConnectionCreatedDetails {
  rate: number;
}

export interface ConnectionClosedDetails {
  rate: number;
}

export interface ChannelCreatedDetails {
  rate: number;
}

export interface ChannelClosedDetails {
  rate: number;
}

export interface QueueDeclaredDetails {
  rate: number;
}

export interface QueueDeletedDetails {
  rate: number;
}

export interface Listener {
  ip_address: string;
  protocol: string;
  port: number;
}

export interface ExchangeType {
  name: string;
}

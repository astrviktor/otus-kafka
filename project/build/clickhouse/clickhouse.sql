CREATE TABLE IF NOT EXISTS jobs_kafka (
  id String,
  status String,
  timestamp String
) ENGINE = Kafka SETTINGS
            kafka_broker_list = 'kafka:9191',
            kafka_topic_list = 'clickhouse-topic',
            kafka_group_name = 'clickhouse',
            kafka_format = 'CSV',
            kafka_num_consumers = 1,
            kafka_skip_broken_messages = 10;


CREATE TABLE IF NOT EXISTS jobs (
  id String,
  status String,
  timestamp DateTime
) ENGINE = MergeTree()
ORDER BY timestamp;

CREATE MATERIALIZED VIEW IF NOT EXISTS jobs_mv
TO jobs
(
  id String,
  status String,
  timestamp DateTime
)
AS
SELECT
  id,
  status,
  parseDateTimeBestEffortOrNull(timestamp) AS timestamp
FROM jobs_kafka;
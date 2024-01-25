CREATE STREAM jobs_stream (
  id VARCHAR KEY,
  status VARCHAR,
  create_date VARCHAR,
  finish_date VARCHAR
) WITH (
  KAFKA_TOPIC = 'processor-topic',
  VALUE_FORMAT = 'JSON'
);

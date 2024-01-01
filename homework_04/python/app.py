from random import random

from datetime import timedelta
import faust

app = faust.App('windowing', broker='kafka://localhost:9092')


class Model(faust.Record, serializer='json'):
    random: float


TOPIC = 'tumbling_topic'
TABLE = 'tumbling_table'

tumbling_topic = app.topic(TOPIC, key_type=str, value_type=Model)
tumbling_table = app.Table(TABLE, partitions=1, default=int).tumbling(size=timedelta(minutes=5))


@app.agent(tumbling_topic)
async def print_windowed_events(stream):
    async for key, value in stream.items():
        tumbling_table[str(key)] += 1
        value = tumbling_table[str(key)]

        print(f'-- New Event accepted from {TOPIC} (every 10 secs) --')
        print(f'{TABLE} - key: {str(key)}, value: {value.value()}')


@app.timer(10.0, on_leader=True)
async def publish_every_1secs():
    rnd = round(random(), 2)
    msg = Model(random=rnd)

    if rnd <= 0.7:
        key = '0'
    else:
        key = '1'

    await tumbling_topic.send(key=key, value=msg)

if __name__ == '__main__':
    app.main()

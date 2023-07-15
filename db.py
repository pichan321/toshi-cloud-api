import os
import psycopg2

conn = psycopg2.connect("postgresql://pichan:OGdtBNEQGFcGS818wuLbxA@pichan-2902.g8z.cockroachlabs.cloud:26257/toshi-cloud?sslmode=verify-full")

try:
    with conn.cursor() as cur:
        cur.execute("show tables")
        res = cur.fetchall()
        conn.commit()
        print(res)
finally:
    conn.close()
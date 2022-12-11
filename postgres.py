
import psycopg2
import uuid
#connectionString = "postgresql://grubhub_user:kGigHiXuy9zWQIj-00BAQg@grubhub-2748.g8z.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full"
#postgresql://pichan:OGdtBNEQGFcGS818wuLbxA@pichan-2902.g8z.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full


try:
    # Connect to an existing database
    connection = psycopg2.connect(user="pichan",
                                  password="OGdtBNEQGFcGS818wuLbxA",
                                  host="pichan-2902.g8z.cockroachlabs.cloud",
                                  port="26257",
                                  database="toshi-cloud")

    # Create a cursor to perform database operations
    cursor = connection.cursor()
    # Print PostgreSQL details
    print("PostgreSQL server information")
    print(connection.get_dsn_parameters(), "\n")
    # Executing a SQL query

    cursor.execute("select * from accounts")
    # Fetch result
    result = cursor.fetchall()
    print(result)
    connection.commit()



except (Exception, Error) as error:
    print("Error while connecting to PostgreSQL", error)
finally:
    if (connection):
        cursor.close()
        connection.close()
        print("PostgreSQL connection is closed")


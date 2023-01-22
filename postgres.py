
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
    #cursor.execute("drop table accounts")
    #cursor.execute("create table accounts (uuid varchar(300) not null primary key, username varchar(300), password varchar(500), email varchar(300), token varchar(300))")

    #cursor.execute("create table files (uuid varchar(300) not null primary key, name varchar(500), size varchar(100), size_mb float, uploaded_date varchar(200), account_uuid varchar(300), bucket_uuid varchar(300), foreign key (account_uuid) references accounts(uuid), foreign key (bucket_uuid) references buckets(uuid))")
    # Fetch result
    cursor.execute("select count(*) from accounts where token = 2b79b781a452a41b8a13e0ebf5cfc4a6 and password = a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3")
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


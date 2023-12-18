import psycopg2
import random
import string
import datetime
from eth_account import Account
import secrets

# Connect to Postgres
conn = psycopg2.connect(
    host="localhost",
    port="5432",
    database="freezer",
    user="root",
    password="pass"
)
cursor = conn.cursor()


# Generate random uint256 value
def generate_uint256():
    return random.randint(0, 115792089237316195423570985008687907853269984665640564039457584007913129639935)


# Helper function to generate random string
def generate_random_string(length):
    letters = string.ascii_letters + string.digits
    return ''.join(random.choice(letters) for i in range(length))


# Helper function to generate random eth address
def generate_eth_address():
    priv = secrets.token_hex(32)
    private_key = "0x" + priv
    acct = Account.from_key(private_key)
    return acct.address


# Insert test data into pending_coin_distributions
def generate_pending_distributions(_cursor, _num_records):
    now = datetime.datetime.now()
    day = datetime.date.today()
    for i in range(_num_records):
        created_at = now - datetime.timedelta(minutes=i)
        internal_id = i
        iceflakes = random.randint(100, 1000)
        user_id = generate_random_string(10)
        eth_address = generate_eth_address()

        _cursor.execute("""
            INSERT INTO pending_coin_distributions (created_at, day, internal_id, iceflakes,
            user_id, eth_address)
            VALUES (%s, %s, %s, %s, %s, %s)
        """, (created_at, day, internal_id, iceflakes, user_id, eth_address))


# Generate data for coin_distributions_by_earner
def generate_coin_distributions_by_earner(_cursor, _num_records):
    now = datetime.date.today()
    for i in range(_num_records):
        internal_id = i
        username = generate_random_string(10)
        referred_by_username = generate_random_string(10)
        user_id = generate_random_string(10)
        earner_user_id = generate_random_string(10)
        eth_address = generate_eth_address()
        balance = random.randint(100, 10000)

        _cursor.execute("""
            INSERT INTO coin_distributions_by_earner (
                created_at, day, internal_id, balance,
                username, referred_by_username,
                user_id, earner_user_id, eth_address
            )
            VALUES (NOW(), %s, %s, %s, %s, %s, %s, %s, %s)
        """, (now, internal_id, balance, username, referred_by_username, user_id, earner_user_id, eth_address))


# Generate data for coin_distributions_pending_review
def generate_coin_distributions_pending_review(_cursor, _num_records):
    now = datetime.date.today()

    for i in range(_num_records):
        internal_id = i
        username = generate_random_string(10)

        # Calculate sum of balance
        _cursor.execute("SELECT SUM(balance) AS total_balance FROM coin_distributions_by_earner")
        rows = _cursor.fetchone()
        ice = rows[0]

        referred_by_username = generate_random_string(10)
        user_id = generate_random_string(10)
        eth_address = generate_eth_address()

        _cursor.execute("""
            INSERT INTO coin_distributions_pending_review (
                created_at, day, internal_id, ice, 
                username, referred_by_username,
                user_id, eth_address
            ) 
            VALUES (NOW(), %s, %s, %s, %s, %s, %s, %s)
        """, (now, internal_id, ice, username, referred_by_username, user_id, eth_address))


# Generate data for reviewed_coin_distributions
def generate_reviewed_coin_distributions(_cursor, _num_records):
    now = datetime.date.today()

    for i in range(_num_records):
        internal_id = i
        username = generate_random_string(10)
        referred_by_username = generate_random_string(10)
        user_id = generate_random_string(10)
        eth_address = generate_eth_address()
        iceflakes = generate_uint256()
        reviewer_user_id = generate_random_string(10)
        decision = random.choice(["approved", "rejected"])

        _cursor.execute("""
            INSERT INTO reviewed_coin_distributions (
                reviewed_at, created_at, day, review_day, ice, internal_id, iceflakes,
                username, referred_by_username, user_id,
                eth_address, reviewer_user_id, decision
            )
            VALUES (NOW(), NOW(), %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        """, (
            now, now, 1, internal_id, iceflakes, username, referred_by_username, user_id, eth_address, reviewer_user_id,
            decision))


num_records = 1700

generate_pending_distributions(cursor, num_records)
conn.commit()
print(f"{num_records} records inserted into pending_coin_distributions")

generate_coin_distributions_by_earner(cursor, num_records)
conn.commit()
print(f"{num_records} records inserted into coin_distributions_by_earner")

generate_coin_distributions_pending_review(cursor, num_records)
conn.commit()
print(f"{num_records} records inserted into coin_distributions_pending_review")

generate_reviewed_coin_distributions(cursor, num_records)
conn.commit()
print(f"{num_records} records inserted into reviewed_coin_distributions")

# Close connection
cursor.close()
conn.close()

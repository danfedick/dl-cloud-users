import json
import random
from faker import Faker

fake = Faker()

job_types = ['manager', 'product-manager' , 'executive', 'devops-engineer', 'devsecops-engineer', 'developer', "customer-success"]


users = []

for _ in range(50):
    user = {
        "name": fake.name(),
        "userid": fake.uuid4(),
        "address": fake.address().replace("\n", ", "),
        "phone": fake.phone_number(),
        "user_agent": fake.user_agent(),
        "company": "OpenGovCo",
        "email": fake.email(),
        "team": random.choice(job_types),
        "location": fake.city(),
        "credit_card": fake.credit_card_number(card_type='mastercard'),
        "social_security": fake.ssn(),
    }
    users.append(user)

with open("users.json", "w") as f:
    json.dump(users, f, indent=2)


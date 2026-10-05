from datetime import datetime, timedelta
import math
import random
from faker import Faker
import requests

# Configuration
# -------------------------------------------------------------------
BASE_URL = "http://localhost:50053"  

# User & Role Distribution
TOTAL_USERS = 120
OWNER_RATIO = 1 / 10  
PROPERTIES_PER_OWNER = 5 
RESERVATIONS_PER_GUEST = 10  
DEFAULT_USER_PASSWORD = "Test1212"  

# Property Randomization Parameters
MAX_LOCATION_RADIUS_KM = 60  
MIN_PROPERTY_PRICE = 50  
MAX_PROPERTY_PRICE = 500  
MIN_PROPERTY_IMAGES = 2  
MAX_PROPERTY_IMAGES = 5  

# Reservation Timing Parameters
MIN_STAY_DAYS = 2  
MAX_STAY_DAYS = 15  
BUFFER_BETWEEN_RESERVATIONS_DAYS = 5  
INITIAL_AVAILABILITY_OFFSET_DAYS = 2  

# Base coordinates: (lat, lon, city_name, country_name)
BASE_COORDINATES = [
    (44.4323, 26.1063, "Bucharest", "Romania"),  
    (46.7712, 23.6236, "Cluj-Napoca", "Romania"),  
    (45.7537, 21.2257, "Timișoara", "Romania"),  
    (47.1585, 27.6014, "Iași", "Romania"),  
    (44.1792, 28.6498, "Constanța", "Romania"),  
    (45.6580, 25.6012, "Brașov", "Romania"),  
    (44.3302, 23.7949, "Craiova", "Romania"),  
    (45.7983, 24.1256, "Sibiu", "Romania"),  
    (47.0515, 21.9402, "Oradea", "Romania"),  
    (45.4353, 28.0080, "Galați", "Romania"),  
]

SAMPLE_IMAGE_URLS = [
    "https://images.unsplash.com/photo-1570129477492-45c003edd2be",  
    "https://images.unsplash.com/photo-1600585154340-be6161a56a0c",  
    "https://images.unsplash.com/photo-1600596542815-ffad4c1539a9",  
    "https://images.unsplash.com/photo-1512917774080-9991f1c4c750",  
    "https://images.unsplash.com/photo-1613977257363-707ba9348227",  
    "https://images.unsplash.com/photo-1580587771525-78b9dba3b914",  
    "https://images.unsplash.com/photo-1513694203232-719a280e022f",  
    "https://images.unsplash.com/photo-1564013799919-ab600027ffc6",
    "https://images.unsplash.com/photo-1572120360610-d971b9d7767c",
    "https://images.unsplash.com/photo-1583608205776-bfd35f0d9f83",
    "https://images.unsplash.com/photo-1598228723793-52759bba239c",
    "https://images.unsplash.com/photo-1600585154526-990dced4db0d",
    "https://images.unsplash.com/photo-1600607687939-ce8a6c25118c",
    "https://images.unsplash.com/photo-1600566753376-12c8ab7fb75b",
    "https://images.unsplash.com/photo-1600573472591-ee6b68d14c68",
    "https://images.unsplash.com/photo-1600585152220-90363fe7e115",
    "https://images.unsplash.com/photo-1600566752355-35792bedcfea",
]

fake = Faker()


def get_random_offset_coords(
    base_lat, base_lon, max_km=MAX_LOCATION_RADIUS_KM
):
    """Generates random coordinates within max_km radius of a base coordinate."""
    r = math.sqrt(random.random()) * max_km  
    theta = random.uniform(0, 2 * math.pi)  

    dx = r * math.cos(theta)  
    dy = r * math.sin(theta)  

    delta_lat = dy / 111.0  
    delta_lon = dx / (111.0 * math.cos(math.radians(base_lat)))  

    new_lat = round(base_lat + delta_lat, 6)  
    new_lon = round(base_lon + delta_lon, 6)  
    return new_lat, new_lon


def generate_random_address(city_name, country_name):
    return f"{fake.street_address()}, {city_name}, {country_name}"


def create_users(num_users):
    """Creates a list of users via POST /user."""
    users = []
    print(f"--- 1. Creating {num_users} Users ---")  

    for i in range(num_users):  
        email = fake.unique.email()  
        username = fake.unique.user_name()  

        payload = {
            "user": {
                "first_name": fake.first_name(),  
                "last_name": fake.last_name(),  
                "user_name": username,  
                "email": email,  
                "password": DEFAULT_USER_PASSWORD,  
            }
        }

        response = requests.post(f"{BASE_URL}/user", json=payload)  
        if response.status_code in [200, 201]:  
            user_data = response.json().get("user", {})  
            user_id = user_data.get("id")  
            print(f"Created user: {username} (ID: {user_id})")  

            users.append(
                {
                    "id": user_id,
                    "email": email,
                    "password": DEFAULT_USER_PASSWORD,
                    "token": None,
                }
            )  
        else:
            print(
                f"Failed to create user {username}: {response.status_code} - {response.text}"
            )  

    return users


def login_users(users):
    """Logs in each user via POST /login to obtain auth tokens."""
    print("\n--- 2. Logging in Users ---")  
    for user in users:  
        payload = {"email": user["email"], "password": user["password"]}  

        response = requests.post(f"{BASE_URL}/login", json=payload)  
        if response.status_code == 200:  
            data = response.json()  
            token = data.get("JWT")  

            user["token"] = token  
            print(f"Successfully logged in user ID: {user['id']}, token: {token}")  
        else:
            print(
                f"Login failed for user ID {user['id']}: {response.status_code} - {response.text}"
            )  


def create_properties(owners):
    """Owners create properties via POST /property."""
    print("\n--- 3. Creating Properties (Owners) ---")  
    created_properties = []  

    for owner in owners:  
        headers = (
            {"Authorization": f"Bearer {owner['token']}"}
            if owner["token"]
            else {}
        )  

        for _ in range(PROPERTIES_PER_OWNER):  
            base_lat, base_lon, city_name, country_name = random.choice(BASE_COORDINATES)
            lat, lon = get_random_offset_coords(
                base_lat, base_lon, max_km=MAX_LOCATION_RADIUS_KM
            )  

            address = generate_random_address(city_name, country_name)

            num_images = random.randint(MIN_PROPERTY_IMAGES, MAX_PROPERTY_IMAGES)  
            selected_images = random.sample(SAMPLE_IMAGE_URLS, k=num_images)  

            price = random.randint(MIN_PROPERTY_PRICE, MAX_PROPERTY_PRICE)  

            payload = {
                "user_id": owner["id"],  
                "name": f"Pensiunea {fake.word().capitalize()}",  
                "description": fake.paragraph(),  
                "address": address,
                "price": price,  
                "location": {
                    "lat": lat,  
                    "long": lon,  
                },
                "image_urls": selected_images,  
            }

            response = requests.post(
                f"{BASE_URL}/property", json=payload, headers=headers
            )  
            if response.status_code in [200, 201]:  
                prop_data = response.json().get("property", {})  
                prop_id = prop_data.get("id")  
                print(
                    f"Owner {owner['id']} created Property '{payload['name']}' "
                    f"at [{lat}, {lon}] ({address}) (ID: {prop_id})"
                )  

                created_properties.append(
                    {
                        "id": prop_id,
                        "owner_id": owner["id"],
                        "next_available_date": datetime.now()
                        + timedelta(days=INITIAL_AVAILABILITY_OFFSET_DAYS),
                    }
                )  
            else:
                print(
                    f"Failed to create property for owner {owner['id']}: {response.status_code} - {response.text}"
                )  

    return created_properties


def create_reservations(guests, properties):
    """Guests create non-overlapping reservations via POST /reservation."""
    print("\n--- 4. Creating Non-Overlapping Reservations (Guests) ---")  
    if not properties:  
        print("No properties available to reserve!")  
        return

    prop_index = 0  

    for guest in guests:  
        headers = (
            {"Authorization": f"Bearer {guest['token']}"}
            if guest["token"]
            else {}
        )  

        for _ in range(RESERVATIONS_PER_GUEST):  
            target_property = properties[prop_index % len(properties)]  
            prop_index += 1  

            check_in = target_property["next_available_date"]  
            stay_duration = random.randint(MIN_STAY_DAYS, MAX_STAY_DAYS)  
            check_out = check_in + timedelta(days=stay_duration)  

            target_property["next_available_date"] = check_out + timedelta(
                days=BUFFER_BETWEEN_RESERVATIONS_DAYS
            )  

            payload = {
                "reservation": {
                    "property_id": target_property["id"],  
                    "user_id": guest["id"],  
                    "owner_id": target_property["owner_id"],  
                    "check_in_date": check_in.strftime("%Y-%m-%d"),  
                    "check_out_date": check_out.strftime("%Y-%m-%d"),  
                    "status": "RESERVATION_STATUS_CONFIRMED",  
                }
            }

            response = requests.post(
                f"{BASE_URL}/reservation", json=payload, headers=headers
            )  
            if response.status_code in [200, 201]:  
                res_data = response.json().get("reservation", {})  
                res_id = res_data.get("id")  
                print(
                    f"Guest {guest['id']} reserved Property {target_property['id']} "
                    f"from {payload['reservation']['check_in_date']} to {payload['reservation']['check_out_date']} (ID: {res_id})"
                )  
            else:
                print(
                    f"Failed to create reservation for guest {guest['id']}: {response.status_code} - {response.text}"
                )  


def main():
    users = create_users(TOTAL_USERS)  
    if not users:  
        print("User creation failed. Exiting.")  
        return

    login_users(users)  

    num_owners = max(1, int(len(users) * OWNER_RATIO))  
    owners = users[:num_owners]  
    guests = users[num_owners:]  

    print(f"\n--- Split: {len(owners)} Owner(s) / {len(guests)} Guest(s) ---")  

    properties = create_properties(owners)  
    create_reservations(guests, properties)  


if __name__ == "__main__":
    main()  
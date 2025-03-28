#!/usr/bin/env python3
import os
import subprocess
import sys
import requests

from time import sleep

LOGIN_URL = "http://127.0.0.1:8000/api/login?clientId=6551a3584403d5264584&responseType=code&redirectUri=http%3A%2F%2Flubricant-core.lubricant.svc.cluster.local%3A8080%2Fapi%2Fv1%2Fsignin&type=code&scope=read&state=casdoor"

CORE_API_BASE_URL = "http://127.0.0.1:8080"
CALLBACK_URL = CORE_API_BASE_URL+"/api/v1/signin"
USER_INFO_URL = CORE_API_BASE_URL+"/api/v1/user/info"
CREATE_GATEWAY_URL = CORE_API_BASE_URL+"/api/v1/gateway/internal/gateway"

COOKIE_FILE = "cookie.txt"

headers = {"Content-Type": "application/json"}
login_data = {
    "application": "application_lubricant",
    "organization": "built-in",
    "username": "admin",
    "autoSignin": True,
    "password": "123",
    "signinMethod": "Password",
    "type": "code"
}
create_gateway_data={
    "host":"", # set to empty means do not bind host Information.
    "description":"test_gateway",
    "username":"username", # Host username
    "password": "password", # Host password
    "tls_config": {
        "enable": False,
        "skip_verify": False,
        "from_file": False,
        "key": "",
        "cert": "",
        "ca": ""
    },
}

def check_pod_status(pod_name, namespace='lubricant'):
    try:
        result = subprocess.run(
            ["kubectl", "get", "pods", "-n", namespace, pod_name, "-o", "jsonpath='{.status.conditions[?(@.type==\"Ready\")].status}'"],
            capture_output=True,
            text=True,
            check=True
        )
        status = result.stdout.strip().strip("'")
        return "Running" if status == "True" else "Not Running"
    except subprocess.CalledProcessError as e:
        print(f"Error: Failed to get Pod status")
        print(f"Output: {e.output}")
        print(f"Error: {e.stderr}")
        sys.exit(1)
def login_and_get_session():
    print("Logging in...")

    try:
        response = requests.post(LOGIN_URL, headers=headers, json=login_data)
        response.raise_for_status()
        da = response.json()
        code = da.get("data")
        print(f"Received code: {code}")
    except Exception as e:
        print("Error: Failed to get code")
        print(f"Response: {response.text if 'response' in locals() else str(e)}")
        sys.exit(1)

    if not code:
        print("Error: Failed to get code")
        print(f"Response: {response.text}")
        sys.exit(1)

    print("Getting cookie...")
    try:
        session = requests.Session()
        callback_response = session.get(f"{CALLBACK_URL}?code={code}&state=casdoor")
        callback_response.raise_for_status()

        # Save cookies to file
        with open(COOKIE_FILE, 'w') as f:
            for cookie in session.cookies:
                f.write(f"{cookie.name}\t{cookie.value}\n")

        msg = callback_response.json().get("msg")
        if msg != "success":
            print(f"Error: Login failed, msg={msg}")
            print(f"Response: {callback_response.text}")
            sys.exit(1)

    except Exception as e:
        print("Error: Failed to get cookie request failed")
        print(f"Response: {callback_response.text if 'callback_response' in locals() else str(e)}")
        sys.exit(1)

    if not os.path.exists(COOKIE_FILE):
        print("Error: Failed to get cookie")
        print(f"File {COOKIE_FILE} does not exist")
        sys.exit(1)

    print(f"Cookie saved to {COOKIE_FILE}")
    with open(COOKIE_FILE, 'r') as f:
        print(f.read())

    return session

def get_user_info(session):
    print("Getting user info...")
    try:
        user_info_response = session.get(USER_INFO_URL)
        user_info_response.raise_for_status()

        msg = user_info_response.json().get("msg")
        if msg != "success":
            print(f"Error: Failed to get user info, msg={msg}")
            print(f"Response: {user_info_response.text}")
            sys.exit(1)

        print(f"User info: {user_info_response.text}")

    except Exception as e:
        print("Error: Failed to get user info")
        print(f"Response: {user_info_response.text if 'user_info_response' in locals() else str(e)}")
        sys.exit(1)

def test_create_gateway(session, gateway_id):
    print("API:Creating gateway...")
    try:
        create_gateway_response = session.post(CREATE_GATEWAY_URL + f"?gateway-id={gateway_id}", headers=headers, json=create_gateway_data)
        create_gateway_response.raise_for_status()

        msg = create_gateway_response.json().get("msg")
        if msg != "success":
            print(f"Error: Failed to create gateway, msg={msg}")
            print(f"Response: {create_gateway_response.text}")
            sys.exit(1)

    except Exception as e:
        print("Error: Failed to create gateway")
        print(f"Response: {create_gateway_response.text if 'create_gateway_response' in locals() else str(e)}")
        sys.exit(1)

    print("kubernetes: Deploy Gateway")
    os.system("kubectl apply -f deployment/tester/gateway.yaml")
    sleep(5)

    pod_status = check_pod_status(gateway_id)
    if pod_status == "Running":
        print("Gateway Deploy Success")
    else:
        print("Gateway Deploy Failed")
        print(f"Pod Status: {pod_status}")
        sys.exit(1)

def test_uncreated_gateway(session, gateway_id):
    print("kubernetes: Test Uncreated Gateway -- Gateway Status should be error.")
    try:
        create_gateway_response = session.post(CREATE_GATEWAY_URL + f"?gateway-id={gateway_id + "-1"}", json=create_gateway_data)
        create_gateway_response.raise_for_status()
    except Exception as e:
        print("Expected error occurred while creating uncreated gateway")
        print(f"Response: {create_gateway_response.text if 'create_gateway_response' in locals() else str(e)}")

    os.system("kubectl scale statefulset lubricant-gateway --replicas=3 -n lubricant")
    sleep(10)

    # os.system("kubectl get pods -n lubricant")
    subprocess.run(["kubectl", "get", "pods", "-n", "lubricant"],stdout=sys.stdout, stderr=sys.stderr)
    pod_status1 = check_pod_status("lubricant-gateway-1")
    pod_status2 = check_pod_status("lubricant-gateway-2")
    if pod_status1 != "Running" or pod_status2 == "Running":
        print("Test Failed:")
        print(f"Pod Status: {pod_status1} {pod_status2} Except: Running,and Error(or CrashLoopBackOff)")
        sys.exit(1)
    else:
        print("Gateway Deploy Success")
def main():
    session = login_and_get_session()
    get_user_info(session)
    print("Begin Test:")
    test_create_gateway(session, "lubricant-gateway-0")
    test_uncreated_gateway(session, "lubricant-gateway")

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: core-test.py <pod_name>")
        sys.exit(1)
    main()

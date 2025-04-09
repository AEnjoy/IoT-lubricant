#!/usr/bin/env python3
import base64
import json
import os
import subprocess
import sys
import time
from time import sleep

import requests
import yaml

LOGIN_URL = "http://127.0.0.1/casdoor-service/api/login?clientId=6551a3584403d5264584&responseType=code&redirectUri=http%3A%2F%2Flubricant-core.lubricant.svc.cluster.local%3A8080%2Fapi%2Fv1%2Fsignin&type=code&scope=read&state=casdoor"

CORE_API_BASE_URL = "http://127.0.0.1/lubricant-service"
CALLBACK_URL = CORE_API_BASE_URL + "/api/v1/signin"
USER_INFO_URL = CORE_API_BASE_URL + "/api/v1/user/info"
QUERY_TASK_STATUS_URL = CORE_API_BASE_URL + "/api/v1/task/query"
CREATE_GATEWAY_URL = CORE_API_BASE_URL + "/api/v1/gateway/add"
ADD_AGENT_URL = CORE_API_BASE_URL  # +"/api/v1/gateway/{ 0 }/agent/internal/add"
SET_AGENT_URL = CORE_API_BASE_URL + "/api/v1/agent/set"
AGENT_OPERATOR_URL = CORE_API_BASE_URL + "/api/v1/agent/operator"

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
create_gateway_data = {
    "host": "",  # set to empty means do not bind host Information.
    "description": "test_gateway",
    "username": "username",  # Host username
    "password": "password",  # Host password
    "tls_config": {
        "enable": False,
        "skip_verify": False,
        "from_file": False,
        "key": "",
        "cert": "",
        "ca": ""
    },
}
add_agent_data = {
    "description": "agent",
    "gather_cycle": 1,
    "report_cycle": 5,
    "address": "lubricant-agent.lubricant.svc.cluster.local:5436",
    "data_compress_algorithm": "default",
    "enable_stream_ability": False,
    "open_api_doc": "",
    "enable_conf": ""
}

set_agent_data = {}
userId = ""


def encode_file_detailed(filename, encoding='utf-8'):
    try:
        with open(filename, 'r', encoding=encoding) as file:
            text = file.read()
        text_bytes = text.encode(encoding)
        encoded_data = base64.b64encode(text_bytes)
        encoded_string = encoded_data.decode('utf-8')

        return encoded_string
    except Exception as e:
        print(f"Error: {str(e)}")
        sys.exit(1)


def read_yaml_file(file_path):
    try:
        with open(file_path, 'r', encoding='utf-8') as file:
            data = yaml.safe_load(file)
        return data
    except FileNotFoundError:
        print(f"Error：File {file_path} Not Found。")
        sys.exit(1)
    except yaml.YAMLError as e:
        print(f"Error：Failed to parse YAML：{e}")
        sys.exit(1)


def write_yaml_file(data, file_path):
    try:
        with open(file_path, 'w', encoding='utf-8') as file:
            yaml.dump(data, file, allow_unicode=True)
    except Exception as e:
        print(f"Error：Failed to write YAML：{e}")
def modify_user_id(data, new_uuid):
    try:
        data['spec']['template']['spec']['containers'][0]['env'][1]['value'] = new_uuid
        return data
    except (KeyError, IndexError, TypeError):
        print("Error: Index")
        sys.exit(1)

def check_pod_status(pod_name, namespace='lubricant'):
    try:
        result = subprocess.run(
            ["kubectl", "get", "pods", "-n", namespace, pod_name, "-o",
             "jsonpath='{.status.conditions[?(@.type==\"Ready\")].status}'"],
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


def get_nested_value(data, keys):
    if not keys:
        return data
    if isinstance(data, dict) and keys[0] in data:
        return get_nested_value(data[keys[0]], keys[1:])
    return None


def set_nested_value(data, keys, value):
    if not keys:
        return value
    if isinstance(data, dict) and keys[0] in data:
        data[keys[0]] = set_nested_value(data[keys[0]], keys[1:], value)
        return data
    return None


def get_value_by_path(data, path):
    keys = path.split('.')
    return get_nested_value(data, keys)


def set_value_by_path(data, path, value):
    keys = path.split('.')
    return set_nested_value(data, keys, value)


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

def modify_user_id(data, new_uuid):
    try:
        data['spec']['template']['spec']['containers'][0]['env'][1]['value'] = new_uuid
        return data
    except (KeyError, IndexError, TypeError):
        print("Error：can't find the specified key or index in the data.")
        sys.exit(1)
def get_user_info(session):
    global userId
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
        userId = user_info_response.json().get("data").get("id")
    except Exception as e:
        print("Error: Failed to get user info")
        print(f"Response: {user_info_response.text if 'user_info_response' in locals() else str(e)}")
        sys.exit(1)


def get_task_status(session, task_id):
    print("API:Getting task status...")
    try:
        get_task_status_response = session.get(QUERY_TASK_STATUS_URL + f"?taskId={task_id}", headers=headers)
        get_task_status_response.raise_for_status()
        msg = get_task_status_response.json().get("msg")
        if msg != "success":
            print(f"Error: Failed to get task status, msg={msg}")
            print(f"Response: {get_task_status_response.text}")
            sys.exit(1)
        return get_task_status_response.json().get("data").get("status")
    except Exception as e:
        print("Error: Failed to get task status")
        print(f"Response: {get_task_status_response.text if 'get_task_status_response' in locals() else str(e)}")
        sys.exit(1)


def test_create_gateway(session, gateway_id):
    print("API:Creating gateway...")
    try:
        create_gateway_response = session.post(CREATE_GATEWAY_URL + f"?gateway-id={gateway_id}", headers=headers,
                                               json=create_gateway_data)
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

    y = read_yaml_file("deployment/tester/gateway.yaml")
    y = modify_user_id(y, userId)
    write_yaml_file(y,"deployment/tester/gateway.yaml")

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
        create_gateway_response = session.post(CREATE_GATEWAY_URL + f"?gateway-id={gateway_id + "-1"}",
                                               json=create_gateway_data)
        create_gateway_response.raise_for_status()
    except Exception as e:
        print("Expected error occurred while creating uncreated gateway")
        print(f"Response: {create_gateway_response.text if 'create_gateway_response' in locals() else str(e)}")

    os.system("kubectl scale statefulset lubricant-gateway --replicas=3 -n lubricant")
    sleep(10)

    # os.system("kubectl get pods -n lubricant")
    subprocess.run(["kubectl", "get", "pods", "-n", "lubricant"], stdout=sys.stdout, stderr=sys.stderr)
    pod_status1 = check_pod_status("lubricant-gateway-1")
    pod_status2 = check_pod_status("lubricant-gateway-2")
    if pod_status1 != "Running" or pod_status2 == "Running":
        print("Test Failed:")
        print(f"Pod Status: {pod_status1} {pod_status2} Except: Running,and Error(or CrashLoopBackOff)")
        sys.exit(1)
    else:
        print("Gateway Deploy Success")


def test_add_agent(session, gateway_id):
    print("API:Adding agent...")
    agent_id = ""
    task_id = ""
    try:
        add_agent_response = session.post(ADD_AGENT_URL + f"/api/v1/gateway/{gateway_id}/agent/internal/add",
                                          headers=headers, json=add_agent_data)
        add_agent_response.raise_for_status()
        msg = add_agent_response.json().get("msg")
        if msg != "success":
            print(f"Error: Failed to add agent(return value failed), msg={msg}")
            print(f"Response: {add_agent_response.text}")
            sys.exit(1)
        agent_id = add_agent_response.json().get("data").get("agent_id")
        task_id = add_agent_response.json().get("data").get("task_id")
    except Exception as e:
        print("Error: Failed to add agent(request failed)")
        print(f"Response: {add_agent_response.text if 'add_agent_response' in locals() else str(e)}")
        sys.exit(1)
    time.sleep(3)
    if get_task_status(session, task_id) != "completed":
        print("Error: Failed to add agent (async task failed)")
        print(f"Response: {add_agent_response.text if 'add_agent_response' in locals() else str(e)}")
        sys.exit(1)
    return agent_id


def test_set_agent(session, gateway_id, agent_id):
    print("API:Setting agent...")
    global set_agent_data
    task_id = ""
    try:
        with open('test/request/set_agent_data.json', 'r', encoding='utf-8') as file:
            set_agent_data = json.load(file)
            set_agent_data['agentID'] = agent_id
            set_agent_data['gatewayID'] = gateway_id
            set_agent_data['dataSource']['originalFile'] = encode_file_detailed("test/mock_driver/clock/api.json")
            set_agent_data['dataSource']['enableFile'] = encode_file_detailed("test/mock_driver/clock/api.json.enable")
    except Exception as e:
        print("Error: Failed to load set_agent_data.json")
        print(f"Error: {str(e)}")
        sys.exit(1)
    try:
        set_agent_response = session.post(SET_AGENT_URL + f'?gateway-id={gateway_id}', headers=headers,
                                          json=set_agent_data)
        set_agent_response.raise_for_status()
        msg = set_agent_response.json().get("msg")
        if msg != "success":
            print(f"Error: Failed to set agent, msg={msg}")
            print(f"Response: {set_agent_response.text}")
            sys.exit(1)
        task_id = set_agent_response.json().get("data").get("taskId")
    except Exception as e:
        print("Error: Failed to set agent")
        print(f"Response: {set_agent_response.text if 'set_agent_response' in locals() else str(e)}")
        sys.exit(1)
    time.sleep(3)
    if get_task_status(session, task_id) != "completed":
        print("Error: Failed to set agent (async task failed)")
        print(f"Response: {set_agent_response.text if 'set_agent_response' in locals() else str(e)}")
        sys.exit(1)


def test_agent_operator(session, gateway_id, agent_id):
    print("API:Operating agent...")
    print("StartGather:")

    def make_url(operator):
        url = AGENT_OPERATOR_URL + f"?agent-id={agent_id}&gateway-id={gateway_id}&operator={operator}"
        return url

    def get_gather_status() -> bool:
        url = make_url("get-gather-status")
        try:
            get_gather_status_response = session.get(url, headers=headers)
            get_gather_status_response.raise_for_status()
            msg = get_gather_status_response.json().get("msg")
            if msg != "success":
                print(f"Error: Failed to get gather status, msg={msg}")
                print(f"Response: {get_gather_status_response.text}")
                sys.exit(1)
            return get_gather_status_response.json().get("data")
        except Exception as e:
            print("Error: Failed to get gather status")
            print(
                f"Response: {get_gather_status_response.text if 'get_gather_status_response' in locals() else str(e)}")
            sys.exit(1)

    task_id = ""

    # start_gather
    try:
        start_gather_response = session.get(make_url("start-gather"))
        start_gather_response.raise_for_status()
        msg = start_gather_response.json().get("msg")
        if msg != "success":
            print(f"Error: Failed to start gather, msg={msg}")
            print(f"Response: {start_gather_response.text}")
            sys.exit(1)
        task_id = start_gather_response.json().get("data").get("taskId")
    except Exception as e:
        print("Error: Failed to start gather")
        print(f"Response: {start_gather_response.text if 'start_gather_response' in locals() else str(e)}")
        sys.exit(1)
    time.sleep(3)
    if get_task_status(session, task_id) != "completed":
        print("Error: Failed to start gather (async task failed)")
        print(f"Response: {start_gather_response.text if 'start_gather_response' in locals() else str(e)}")
        sys.exit(1)
    if not get_gather_status():
        print("Error: Failed to start gather (gather status is not true)")
        print(f"Response: {start_gather_response.text if 'start_gather_response' in locals() else str(e)}")
        sys.exit(1)
    # stop_gather
    try:
        stop_gather_response = session.get(make_url("stop-gather"), headers=headers)
        stop_gather_response.raise_for_status()
        msg = stop_gather_response.json().get("msg")
        if msg != "success":
            print(f"Error: Failed to stop gather, msg={msg}")
            print(f"Response: {stop_gather_response.text}")
            sys.exit(1)
        task_id = stop_gather_response.json().get("data").get("taskId")
    except Exception as e:
        print("Error: Failed to stop gather")
        print(f"Response: {stop_gather_response.text if 'stop_gather_response' in locals() else str(e)}")
        sys.exit(1)
    time.sleep(3)
    if get_task_status(session, task_id) != "completed":
        print("Error: Failed to stop gather (async task failed)")
        print(f"Response: {stop_gather_response.text if 'stop_gather_response' in locals() else str(e)}")
        sys.exit(1)
    if get_gather_status():
        print("Error: Failed to stop gather (gather status is not false)")
        print(f"Response: {stop_gather_response.text if 'stop_gather_response' in locals() else str(e)}")
        sys.exit(1)


def main():
    session = login_and_get_session()
    get_user_info(session)
    print("Begin Test:")
    test_create_gateway(session, "lubricant-gateway-0")
    test_uncreated_gateway(session, "lubricant-gateway")
    agent_id = test_add_agent(session, "lubricant-gateway-0")
    test_set_agent(session, "lubricant-gateway-0", agent_id)
    test_agent_operator(session, "lubricant-gateway-0", agent_id)


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: core-test.py <pod_name>")
        sys.exit(1)
    main()

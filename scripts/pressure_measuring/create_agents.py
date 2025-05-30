#!/usr/bin/env python3
import base64
import json
import os
import sys
import time

import requests

LOGIN_URL = "http://127.0.0.1/casdoor-service/api/login?clientId=6551a3584403d5264584&responseType=code&redirectUri=http%3A%2F%2Flubricant-core.lubricant.svc.cluster.local%3A8080%2Fapi%2Fv1%2Fsignin&type=code&scope=read&state=casdoor"

CORE_API_BASE_URL = "http://127.0.0.1/lubricant-service"
CALLBACK_URL = CORE_API_BASE_URL + "/api/v1/signin"
ADD_AGENT_URL = CORE_API_BASE_URL
QUERY_TASK_STATUS_URL = CORE_API_BASE_URL + "/api/v1/task/query"
SET_AGENT_URL = CORE_API_BASE_URL + "/api/v1/agent/set"

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
add_agent_data = {
    "description": "agent",
    "gather_cycle": 1,
    "report_cycle": 5,
    "address": "lubricant-agent-0.lubricant-agent.lubricant.svc.cluster.local:5436",
    "data_compress_algorithm": "default",
    "enable_stream_ability": False,
    "open_api_doc": "",
    "enable_conf": ""
}

userId = ""
set_agent_data = {}

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

def test_set_agent(session, gateway_id, agent_id):
    print("API:Setting agent...")
    global set_agent_data
    task_id = ""
    try:
        with open('scripts/test/request/set_agent_data.json', 'r', encoding='utf-8') as file:
            set_agent_data = json.load(file)
            set_agent_data['agentID'] = agent_id
            set_agent_data['gatewayID'] = gateway_id
            set_agent_data['dataSource']['originalFile'] = encode_file_detailed("scripts/test/mock_driver/clock/api.json")
            set_agent_data['dataSource']['enableFile'] = encode_file_detailed("scripts/test/mock_driver/clock/api.json.enable")
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

def main():
    if len(sys.argv) != 4:
        print("Usage: python create_agents.py <gateway_id> <agent_number> <project_id>")
        sys.exit(1)

    gateway_id = sys.argv[1]
    agent_number = int(sys.argv[2])
    if agent_number <= 0:
        print("Error: agent_number must be greater than 0")
        sys.exit(1)
    add_agent_data['project_id'] = sys.argv[3]
    session = login_and_get_session()
    print("API:Login successful")

    for i in range(agent_number):
        print(f"API:Creating agent {i + 1}/{agent_number}...")
        agent_id = test_add_agent(session, gateway_id)
        print(f"API:Agent {i + 1}/{agent_number} created successfully, agent_id={agent_id}")
        # test_set_agent(session, gateway_id, agent_id) 在测试环境里,不需要设置agent 仅开发环境需要
        with open("agent_id.txt", "a") as f:
            f.write(f"{agent_id}\n")

if __name__ == "__main__":
    main()

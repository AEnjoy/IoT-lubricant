import sys
import time

import requests


def main(core_url):
    print("Waiting for service health check to return OK...")
    
    for i in range(30):
        try:
            response = requests.get(f"{core_url}/health")
            status_code = response.status_code
            if status_code == 200 and '"status": "OK"' in response.text:
                print("Service is healthy.")
                return 0
        except requests.RequestException:
            pass
        
        time.sleep(1)
    
    print("Service health check failed after waiting for 1 minute.")
    return 1

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python health_check.py <CORE_URL>")
        sys.exit(1)
    
    core_url = f"http://{sys.argv[1]}"
    exit_code = main(core_url)
    sys.exit(exit_code)

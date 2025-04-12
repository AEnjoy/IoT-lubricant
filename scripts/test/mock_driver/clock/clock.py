#!/usr/bin/env python3
from flask import Flask, request, jsonify
import time

app = Flask(__name__)

# 初始化模拟时间，设置为 2024-01-01 20:00:00
mock_time = time.mktime(time.strptime("2024-01-01 20:00:00", "%Y-%m-%d %H:%M:%S"))

# 每秒更新一次模拟时间
def update_mock_time():
    global mock_time
    mock_time += 1

# 获取当前模拟时间并返回 JSON 格式
@app.route('/api/v1/get/time', methods=['GET'])
def get_time():
    return jsonify({'time': time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(mock_time))})

# 设置模拟时间
@app.route('/api/v1/set/time', methods=['POST', 'GET'])
def set_time():
    global mock_time
    if request.method == 'POST':
        data = request.get_json()
        new_time = data.get('time')
    else:
        new_time = request.args.get('time')
    try:
        mock_time = time.mktime(time.strptime(new_time, "%Y-%m-%d %H:%M:%S"))
        return jsonify({'message': 'Time set successfully'})
    except ValueError:
        return jsonify({'error': 'Invalid time format'}), 400

if __name__ == '__main__':
    # 在后台线程中不断更新模拟时间
    import threading
    def run_update():
        while True:
            update_mock_time()
            time.sleep(1)
    thread = threading.Thread(target=run_update)
    thread.daemon = True
    thread.start()
    app.run(debug=True,port=80)
